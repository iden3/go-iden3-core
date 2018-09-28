package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/identitysrv"
	"github.com/iden3/go-iden3/services/rootsrv"
	"github.com/iden3/go-iden3/services/signsrv"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func assert(msg string, err error) {
	if err != nil {
		log.Error(msg, " ", err.Error())
		os.Exit(1)
	}
}

func LoadKeyStore() (*keystore.KeyStore, accounts.Account) {
	// Load keystore
	ks := keystore.NewKeyStore(C.KeyStore.Path, keystore.StandardScryptN, keystore.StandardScryptP)
	passwd, err := ioutil.ReadFile(C.KeyStore.Password)
	assert("Cannot read password", err)

	acc, err := ks.Find(accounts.Account{
		Address: common.HexToAddress(C.KeyStore.Address),
	})
	assert("Cannot find ksystore account", err)

	assert("Cannot unlock account", ks.Unlock(acc, string(passwd)))
	log.WithField("acc", acc.Address.Hex()).Info("Keystore and account unlocked sucessfully")

	return ks, acc
}

func LoadWeb3(ks *keystore.KeyStore, acc *accounts.Account) *eth.Web3Client {
	// Create geth client
	web3cli, err := eth.NewWeb3Client(C.Web3.Url, ks, acc)
	assert("Cannot open connection to web3", err)
	log.WithField("url", C.Web3.Url).Info("Connection to web3 server opened")
	return web3cli
}

func LoadMerkele() *merkletree.MerkleTree {
	// Open database
	storage, err := merkletree.NewLevelDbStorage(C.Storage.Path)
	assert("Cannot open database", err)

	mt, err := merkletree.New(storage, 140)
	assert("Cannot open merkle tree", err)
	log.WithField("path", C.Web3.Url).Info("Database opened")
	log.WithField("hash", mt.Root().Hex()).Info("Current root")

	return mt
}

func LoadContract(client eth.Client, jsonabifile string, address *string) *eth.Contract {
	abiFile, err := os.Open(jsonabifile)
	assert("Cannot read contract "+jsonabifile, err)
	abi, code, err := eth.UnmarshallSolcAbiJson(abiFile)
	assert("Cannot parse contract "+jsonabifile, err)
	var addrPtr *common.Address
	if address != nil && len(strings.TrimSpace(*address)) > 0 {
		addr := common.HexToAddress(strings.TrimSpace(*address))
		addrPtr = &addr
	}
	return eth.NewContract(client, abi, code, addrPtr)
}

func LoadIdService(client *eth.Web3Client) identitysrv.Service {
	deployerContract := LoadContract(
		client,
		C.Contracts.Iden3Deployer.JsonABI,
		&C.Contracts.Iden3Deployer.Address)

	implContract := LoadContract(
		client,
		C.Contracts.Iden3Impl.JsonABI,
		&C.Contracts.Iden3Impl.Address)

	proxyContract := LoadContract(
		client,
		C.Contracts.Iden3Proxy.JsonABI,
		nil)

	return identitysrv.New(deployerContract, implContract, proxyContract)
}

func LoadRootsService(client *eth.Web3Client) rootsrv.Service {
	return rootsrv.New(LoadContract(
		client,
		C.Contracts.RootCommits.JsonABI,
		&C.Contracts.RootCommits.Address,
	))
}

func LoadClaimService(mt *merkletree.MerkleTree, rootsrv rootsrv.Service, ks *keystore.KeyStore, acc accounts.Account) claimsrv.Service {
	return claimsrv.New(mt, rootsrv, signsrv.New(ks, acc))
}
