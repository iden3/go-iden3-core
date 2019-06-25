package eth

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	log "github.com/sirupsen/logrus"
)

type Client2 struct {
	client         *ethclient.Client
	account        *accounts.Account
	ks             *ethkeystore.KeyStore
	ReceiptTimeout time.Duration
}

func NewClient2(client *ethclient.Client, account *accounts.Account, ks *ethkeystore.KeyStore) *Client2 {
	return &Client2{client: client, account: account, ks: ks, ReceiptTimeout: 120 * time.Second}
}

func (c *Client2) CallAuth(fn func(*ethclient.Client, *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
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

func (c *Client2) Call(fn func(*ethclient.Client) error) error {
	return fn(c.client)
}

func (c *Client2) WaitReceipt(tx *types.Transaction) (*types.Receipt, error) {
	var err error
	var receipt *types.Receipt

	txid := tx.Hash()
	log.WithField("tx", txid.Hex()).Info("Waiting for receipt")

	start := time.Now()
	for receipt == nil && time.Now().Sub(start) < c.ReceiptTimeout {
		receipt, err = c.client.TransactionReceipt(context.TODO(), txid)
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
