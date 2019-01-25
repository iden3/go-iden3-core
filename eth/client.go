package eth

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/iden3/go-iden3/utils"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrReceiptStatusFailed when receiving a failed transaction
	errReceiptStatusFailed = errors.New("receipt status is failed")
	// ErrReceiptNotRecieved when unable to retrieve a transaction
	errReceiptNotRecieved = errors.New("receipt not available")
)

type Client interface {
	WaitReceipt(txid common.Hash) (*types.Receipt, error)
	SendTransaction(to *common.Address, value *big.Int, gasLimit uint64, calldata []byte) (*types.Transaction, error)
	Call(to *common.Address, value *big.Int, calldata []byte) ([]byte, error)
	Sign(data ...[]byte) ([3][32]byte, error)
	NetworkID() (*big.Int, error)
	CodeAt(addr common.Address) ([]byte, error)
	BalanceAt(addr common.Address) (*big.Int, error)
}

// Web3Client defines a connection to a client via websockets
type Web3Client struct {
	mutex          *sync.Mutex
	rpcclient      *rpc.Client
	ethclient      *ethclient.Client
	account        *accounts.Account
	ks             *keystore.KeyStore
	ReceiptTimeout time.Duration
	MaxGasPrice    uint64
}

// NewWeb3Client creates a client, using a keystore and an account for transactions
func NewWeb3Client(rpcURL string, ks *keystore.KeyStore, account *accounts.Account) (*Web3Client, error) {

	rpcClient, err := rpc.DialContext(context.TODO(), rpcURL)
	if err != nil {
		return nil, err
	}

	return &Web3Client{
		mutex:          &sync.Mutex{},
		rpcclient:      rpcClient,
		ethclient:      ethclient.NewClient(rpcClient),
		ks:             ks,
		account:        account,
		ReceiptTimeout: 120 * time.Second,
		MaxGasPrice:    4000000000,
	}, nil
}

// BalanceAt retieves information about the default account
func (c *Web3Client) BalanceAt(addr common.Address) (*big.Int, error) {
	return c.ethclient.BalanceAt(context.TODO(), addr, nil)
}

// SendTransaction executes a contract method and wait it finalizes
func (c *Web3Client) SendTransaction(to *common.Address, value *big.Int, gasLimit uint64, calldata []byte) (*types.Transaction, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error
	var tx *types.Transaction

	ctx := context.TODO()

	if value == nil {
		value = big.NewInt(0)
	}

	network, err := c.ethclient.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.ethclient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	if c.MaxGasPrice > 0 && gasPrice.Uint64() > c.MaxGasPrice {
		return nil, fmt.Errorf("Max gas price reached %v > %v", gasPrice, c.MaxGasPrice)
	}

	callmsg := ethereum.CallMsg{
		From:  c.account.Address,
		To:    to,
		Value: value,
		Data:  calldata,
	}

	if gasLimit == 0 {
		gasLimit, err = c.ethclient.EstimateGas(ctx, callmsg)
		if err != nil {
			sendto := "nil"
			if callmsg.To != nil {
				sendto = callmsg.To.Hex()
			}
			log.Errorf("WEB3 Failed EstimateGas from=%v to=%v value=%v data=%v",
				callmsg.From.Hex(), sendto,
				callmsg.Value, hex.EncodeToString(callmsg.Data),
			)
			return nil, err
		}
	}

	nonce, err := c.ethclient.NonceAt(ctx, c.account.Address, nil)
	if err != nil {
		return nil, err
	}

	if to == nil {
		tx = types.NewContractCreation(
			nonce,    // nonce int64
			value,    // amount *big.Int
			gasLimit, // gasLimit *big.Int
			gasPrice, // gasPrice *big.Int
			calldata, // data []byte
		)
	} else {
		tx = types.NewTransaction(
			nonce,    // nonce int64
			*to,      // to common.Address
			value,    // amount *big.Int
			gasLimit, // gasLimit *big.Int
			gasPrice, // gasPrice *big.Int
			calldata, // data []byte
		)
	}

	if tx, err = c.ks.SignTx(*c.account, tx, network); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"tx":       tx.Hash().Hex(),
		"gasprice": fmt.Sprintf("%.2f Gwei", float64(tx.GasPrice().Uint64())/1000000000.0),
	}).Info("WEB3 Sending transaction")
	if err = c.ethclient.SendTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, err
}

func (c *Web3Client) WaitReceipt(txid common.Hash) (*types.Receipt, error) {
	var err error
	var receipt *types.Receipt

	log.WithField("tx", txid.Hex()).Info("Waiting for receipt")

	start := time.Now()
	for receipt == nil && time.Now().Sub(start) < c.ReceiptTimeout {
		receipt, err = c.ethclient.TransactionReceipt(context.TODO(), txid)
		if receipt == nil {
			time.Sleep(200 * time.Millisecond)
		}
	}

	if receipt != nil && receipt.Status == types.ReceiptStatusFailed {
		log.WithField("tx", txid.Hex()).Error("WEB3 Failed transaction receipt")
		return receipt, errReceiptStatusFailed
	}

	if receipt == nil {
		log.WithField("tx", txid.Hex()).Error("WEB3 Failed transaction")
		return receipt, errReceiptNotRecieved
	}
	log.WithField("tx", txid.Hex()).Debug("WEB3 Success transaction")

	return receipt, err
}

// Call a constant method
func (c *Web3Client) Call(to *common.Address, value *big.Int, calldata []byte) ([]byte, error) {

	ctx := context.TODO()

	msg := ethereum.CallMsg{
		From:  c.account.Address,
		To:    to,
		Value: value,
		Data:  calldata,
	}

	return c.ethclient.CallContract(ctx, msg, nil)
}

// Sign does a web3 signature
func (c *Web3Client) Sign(data ...[]byte) ([3][32]byte, error) {

	var ret [3][32]byte

	// The produced signature is in the [R || S || V] format where V is 0 or 1.
	sig, err := utils.SignEthMsg(c.ks, *c.account, data[0])
	if err != nil {
		return ret, err
	}

	// We need to convert it to the format []uint256 = {v,r,s} format
	ret[0][31] = sig[64]
	copy(ret[1][:], sig[0:32])
	copy(ret[2][:], sig[32:64])
	return ret, nil
}

func (c *Web3Client) NetworkID() (*big.Int, error) {
	return c.ethclient.NetworkID(context.TODO())
}

func (c *Web3Client) CodeAt(addr common.Address) ([]byte, error) {
	return c.ethclient.CodeAt(context.TODO(), addr, nil)
}

func (c *Web3Client) Account() *accounts.Account {
	return c.account
}

func (c *Web3Client) KeyStore() *keystore.KeyStore {
	return c.ks
}
