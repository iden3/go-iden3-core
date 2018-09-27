package config

import (
	"strings"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/identitysrv"
	"github.com/iden3/go-iden3/services/rootsrv"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/signsrv"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func assert(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func LoadKeyStore() (*keystore.KeyStore, accounts.Account) {
	// Load keystore
	ks := keystore.NewKeyStore(C.KeyStore.Path, keystore.StandardScryptN, keystore.StandardScryptP)
	passwd, err := ioutil.ReadFile(C.KeyStore.Password)
	assert(err)

	acc, err := ks.Find(accounts.Account{
		Address: common.HexToAddress(C.KeyStore.Address),
	})
	assert(err)

	assert(ks.Unlock(acc, string(passwd)))
	log.WithField("acc", acc.Address.Hex()).Info("Keystore and account unlocked sucessfully")

	return ks, acc
}

func LoadWeb3(ks *keystore.KeyStore, acc *accounts.Account) *eth.Web3Client {
	// Create geth client
	web3cli, err := eth.NewWeb3Client(C.Web3.Url, ks, acc)
	assert(err)
	log.WithField("url", C.Web3.Url).Info("Connection to web3 server opened")
	return web3cli
}

func LoadMerkele() *merkletree.MerkleTree {
	// Open database
	storage, err := merkletree.NewLevelDbStorage(C.Storage.Path)
	assert(err)
	defer storage.Close()

	mt, err := merkletree.New(storage, 140)
	assert(err)
	log.WithField("path", C.Web3.Url).Info("Database opened")
	log.WithField("hash", mt.Root().Hex()).Info("Current root")

	return mt
}

func loadContract(client eth.Client, jsonabifile string,  address *string) *eth.Contract {
	abiFile, err := os.Open(jsonabifile)
	assert(err)
	abi, code,err  := eth.UnmarshallSolcAbiJson(abiFile)
	assert(err)
	var addrPtr *common.Address
	if address == nil || len(strings.TrimSpace(*address))>0 {
		addr := common.HexToAddress(strings.TrimSpace(*address))
		addrPtr = &addr
	}
	return eth.NewContract(client, abi, code,addrPtr)	
}

func LoadIdService(client *eth.Web3Client) identitysrv.Service {
	deployerContract := loadContract(
		client,
		C.Contracts.Iden3Deployer.JsonABI,
		&C.Contracts.Iden3Deployer.Address)

	implContract := loadContract(
		client,
		C.Contracts.Iden3Impl.JsonABI,
		&C.Contracts.Iden3Impl.Address)

	proxyContract := loadContract(
		client,
		C.Contracts.Iden3Proxy.JsonABI,
		nil)

	return identitysrv.New(deployerContract,implContract,proxyContract)
} 

func LoadRootsService(client *eth.Web3Client) rootsrv.Service {	
	return rootsrv.New(loadContract(
		client,
		C.Contracts.RootCommits.JsonABI,
		&C.Contracts.RootCommits.Address,
	))
}

func LoadClaimService(mt *merkletree.MerkleTree, rootsrv rootsrv.Service, ks *keystore.KeyStore, acc accounts.Account) claimsrv.Service {
	return claimsrv.New(mt,rootsrv,signsrv.New(ks,acc))
}
