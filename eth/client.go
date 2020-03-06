package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	log "github.com/sirupsen/logrus"
)

var (
	ErrAccountNil = fmt.Errorf("Authorized calls can't be made when the account is nil")
	// ErrReceiptStatusFailed when receiving a failed transaction
	ErrReceiptStatusFailed = fmt.Errorf("receipt status is failed")
	// ErrReceiptNotRecieved when unable to retrieve a transaction
	ErrReceiptNotRecieved = fmt.Errorf("receipt not available")
)

// Client is an ethereum client to call Smart Contract methods.
type Client struct {
	client         *ethclient.Client
	account        *accounts.Account
	ks             *ethkeystore.KeyStore
	ReceiptTimeout time.Duration
}

// NewClient creates a Client instance.  The account is not mandatory (it can
// be nil).  If the account is nil, CallAuth will fail with ErrAccountNil.
func NewClient(client *ethclient.Client, account *accounts.Account, ks *ethkeystore.KeyStore) *Client {
	return &Client{client: client, account: account, ks: ks, ReceiptTimeout: 60 * time.Second}
}

// BalanceAt retieves information about the default account
func (c *Client) BalanceAt(addr common.Address) (*big.Int, error) {
	return c.client.BalanceAt(context.TODO(), addr, nil)
}

// Account returns the underlying ethereum account
func (c *Client) Account() *accounts.Account {
	return c.account
}

// CallAuth performs a Smart Contract method call that requires authorization.
// This call requires a valid account with Ether that can be spend during the
// call.
func (c *Client) CallAuth(fn func(*ethclient.Client, *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	if c.account == nil {
		return nil, ErrAccountNil
	}
	nonce, err := c.client.PendingNonceAt(context.Background(), c.account.Address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyStoreTransactor(c.ks, *c.account)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	return fn(c.client, auth)
}

// Call performs a read only Smart Contract method call.
func (c *Client) Call(fn func(*ethclient.Client) error) error {
	return fn(c.client)
}

// WaitReceipt will block until a transaction is confirmed.  Internally it
// polls the state every 200 milliseconds.
func (c *Client) WaitReceipt(tx *types.Transaction) (*types.Receipt, error) {
	var err error
	var receipt *types.Receipt

	txid := tx.Hash()
	log.WithField("tx", txid.Hex()).Debug("Waiting for receipt")

	start := time.Now()
	for receipt == nil && time.Since(start) < c.ReceiptTimeout {
		receipt, err = c.client.TransactionReceipt(context.TODO(), txid)
		if receipt == nil {
			time.Sleep(200 * time.Millisecond)
		}
	}

	if receipt != nil && receipt.Status == types.ReceiptStatusFailed {
		log.WithField("tx", txid.Hex()).Error("WEB3 Failed transaction receipt")
		return receipt, ErrReceiptStatusFailed
	}

	if receipt == nil {
		log.WithField("tx", txid.Hex()).Error("WEB3 Failed transaction")
		return receipt, ErrReceiptNotRecieved
	}
	log.WithField("tx", txid.Hex()).Debug("WEB3 Success transaction")

	return receipt, err
}
