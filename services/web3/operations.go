package web3srv

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	client  *ethclient.Client
	key     *ecdsa.PrivateKey
	Address common.Address
)

// Open connects to Geth node and sets the key and Address
func Open(gethURL string, privK string) error {
	// geth set up
	var err error

	client, err = ethclient.Dial(gethURL)
	if err != nil {
		return err
	}
	key, err = crypto.HexToECDSA(privK)
	if err != nil {
		return err
	}
	Address = crypto.PubkeyToAddress(key.PublicKey)

	return nil
}

// GetBalance returns the current balance for a given address
func GetBalance(addr common.Address) (*big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// AddRoot executes the smart contract call to add a root to an identity
func AddRoot(root32 [32]byte, contractAddress common.Address) error {
	// update authServer merkle root in Identities smart contract
	instance, err := NewIdentitiesContract(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), Address)
	if err != nil {
		return err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(key)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	_, err = instance.SetRoot(auth, root32)
	return err
}

// GetRoot gets the current root of a given identity, from the smart contract of roots
func GetRoot(contractAddress common.Address) ([32]byte, error) {
	// this verification will be done in the smart contract
	// check contract.getAtBlock(block).root == root. TODO Now is getting at the current block, needs at the specified block
	// contractAddress := common.HexToAddress(contractAddressHex)
	instance, err := NewIdentitiesContract(contractAddress, client)
	if err != nil {
		return [32]byte{}, nil
	}
	return instance.GetRoot(nil, Address)
}
