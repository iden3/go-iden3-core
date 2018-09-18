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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"


	log "github.com/sirupsen/logrus"
)

var (
	// ErrReceiptStatusFailed when recieving a failed transaction
	errReceiptStatusFailed = errors.New("receipt status is failed")
	// ErrReceiptNotRecieved when unable to retrieve a transaction
	errReceiptNotRecieved = errors.New("receipt not available")
)

type Client interface {
	SendRawTx(tx []byte);
	SendTransactionSync(to *common.Address, value *big.Int, gasLimit uint64, calldata []byte) (*types.Transaction, *types.Receipt, error)
	Call(to *common.Address, value *big.Int, calldata []byte) ([]byte, error)
	Sign(data ...[]byte) ([3][32]byte, error)
	NetworkID() (*big.Int,error)
	CodeAt(account common.Address) ([]byte, error)
	BalanceInfo() (string, error)
}

// Web3Client defines a connection to a client via websockets
type Web3Client struct {
	mutex    	   *sync.Mutex
	rpcclient      *rpc.Client
	ethclient      *ethclient.Client
	account        *accounts.Account
	ks             *keystore.KeyStore
	receiptTimeout time.Duration
	maxGasPrice    uint64
}

// NewWeb3Client creates a client, using a keystore and an account for transactions
func NewClientWithURL(rpcURL string, ks *keystore.KeyStore, account *accounts.Account) (*Web3Client, error) {

	rpcClient, err := rpc.DialContext(context.TODO(), rpcURL)
	if err != nil {
		return nil, err
	}

	return &Web3Client{
		rpcclient:      rpcClient,
		ethclient:      ethclient.NewClient(rpcClient),
		ks:             ks,
		account:        account,
		receiptTimeout: 120 * time.Second,
	}, nil
}

// NewWeb3Client creates a client, using a keystore and an account for transactions
func NewWeb3Client(client *ethclient.Client, ks *keystore.KeyStore, account *accounts.Account) *Web3Client {

	return &Web3Client{
		ethclient:         client,
		ks:             ks,
		account:        account,
		receiptTimeout: 120 * time.Second,
		maxGasPrice:    4000000000,
	}
}

// BalanceInfo retieves information about the default account
func (c *Web3Client) BalanceInfo() (string, error) {

	ctx := context.TODO()
	balance, err := c.ethclient.BalanceAt(ctx, c.account.Address, nil)
	if err != nil {

		return "", err
	}
	return balance.String(), nil
}


// SendTransactionSync executes a contract method and wait it finalizes
func (c *Web3Client) SendTransactionSync(to *common.Address, value *big.Int, gasLimit uint64, calldata []byte) (*types.Transaction, *types.Receipt, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error
	var tx *types.Transaction
	var receipt *types.Receipt

	ctx := context.TODO()

	if value == nil {
		value = big.NewInt(0)
	}

	network, err := c.ethclient.NetworkID(ctx)
	if err != nil {
		return nil, nil, err
	}

	gasPrice, err := c.ethclient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}

	if gasPrice.Uint64() > c.maxGasPrice {
		return nil, nil, fmt.Errorf("Max gas price reached %v > %v", gasPrice, c.maxGasPrice)
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
			log.Errorf("WEB3 Failed EstimateGas from=%v to=%v value=%v data=%v",
				callmsg.From.Hex(), callmsg.To.Hex(),
				callmsg.Value, hex.EncodeToString(callmsg.Data),
			)
			return nil, nil, err
		}
	}

	nonce, err := c.ethclient.NonceAt(ctx, c.account.Address, nil)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	log.WithFields(log.Fields{
		"tx":       tx.Hash().Hex(),
		"gasprice": fmt.Sprintf("%.2f Gwei", float64(tx.GasPrice().Uint64())/1000000000.0),
	}).Info("WEB3 Sending transaction")
	if err = c.ethclient.SendTransaction(ctx, tx); err != nil {
		return nil, nil, err
	}

	receipt, err = c.waitRecipt(tx.Hash())

	return tx, receipt, err
}
func (c *Web3Client) waitRecipt(txid common.Hash) (*types.Receipt, error) {
	var err error
	var receipt *types.Receipt

	start := time.Now()
	for receipt == nil && time.Now().Sub(start) < c.receiptTimeout {
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

// Call an constant method
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

// Do a web3 signature
func (c *Web3Client) Sign(data ...[]byte) ([3][32]byte, error) {
	web3SignaturePrefix := []byte("\x19Ethereum Signed Message:\n32")

	hash := crypto.Keccak256(data...)
	prefixedHash := crypto.Keccak256(web3SignaturePrefix, hash)

	var ret [3][32]byte

	// The produced signature is in the [R || S || V] format where V is 0 or 1.
	sig, err := c.ks.SignHash(*c.account, prefixedHash)
	if err != nil {
		return ret, err
	}

	// We need to convert it to the format []uint256 = {v,r,s} format
	ret[0][31] = sig[64] + 27
	copy(ret[1][:], sig[0:32])
	copy(ret[2][:], sig[32:64])
	return ret, nil
}

func (c *Web3Client) NetworkID() (*big.Int,error) {
	return c.ethclient.NetworkID(context.TODO())
}

func (c *Web3Client) CodeAt(account common.Address) ([]byte, error) {
	return c.ethclient.CodeAt(context.TODO(),account,nil)
}

func (c *Web3Client) SendRawTxSync(data []byte) (*types.Receipt,error) {
	txid := crypto.Keccak256Hash(data)
	err := c.rpcclient.CallContext(context.TODO(), nil, "eth_sendRawTransaction", common.ToHex(data))
	if err != nil {
		return nil, err
	}
	return c.waitRecipt(txid)
}

