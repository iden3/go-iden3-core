package genericserver

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/adminsrv"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/identitysrv"
	"github.com/iden3/go-iden3/services/rootsrv"
	"github.com/iden3/go-iden3/services/signsrv"
	log "github.com/sirupsen/logrus"
)

var (
	dbMerkletreePrefix = []byte{0}
	dbIdentityPrefix   = []byte{1}
)

const (
	passwdPrefix = "passwd:"
	filePrefix   = "file:"
)

func Assert(msg string, err error) {
	if err != nil {
		log.Error(msg, " ", err.Error())
		os.Exit(1)
	}
}

var Claimservice claimsrv.Service
var Rootservice rootsrv.Service

var Idservice identitysrv.Service

var Adminservice adminsrv.Service

func LoadKeyStore() (*keystore.KeyStore, accounts.Account) {

	var err error
	var passwd string

	// Load keystore
	ks := keystore.NewKeyStore(C.KeyStore.Path, keystore.StandardScryptN, keystore.StandardScryptP)

	// Password can be prefixed by two options
	//   file: <path to file containing the password>
	//   passwd: raw password
	// if is not prefixed by any of those, file: is used
	if strings.HasPrefix(C.KeyStore.Password, passwdPrefix) {
		passwd = C.KeyStore.Password[len(passwdPrefix):]
	} else {
		filename := C.KeyStore.Password
		if strings.HasPrefix(filename, filePrefix) {
			filename = C.KeyStore.Password[len(filePrefix):]
		}
		passwdbytes, err := ioutil.ReadFile(filename)
		Assert("Cannot read password ", err)
		passwd = string(passwdbytes)
	}

	acc, err := ks.Find(accounts.Account{
		Address: common.HexToAddress(C.KeyStore.Address),
	})
	Assert("Cannot find keystore account", err)

	Assert("Cannot unlock account", ks.Unlock(acc, string(passwd)))
	log.WithField("acc", acc.Address.Hex()).Info("Keystore and account unlocked successfully")

	return ks, acc
}

func LoadWeb3(ks *keystore.KeyStore, acc *accounts.Account) *eth.Web3Client {
	// Create geth client
	url := C.Web3.Url
	hidden := strings.HasPrefix(url, "hidden:")
	if hidden {
		url = url[len("hidden:"):]
	}
	web3cli, err := eth.NewWeb3Client(url, ks, acc)
	Assert("Cannot open connection to web3 ", err)
	if hidden {
		log.WithField("url", "(hidden)").Info("Connection to web3 server opened")
	} else {
		log.WithField("url", C.Web3.Url).Info("Connection to web3 server opened")
	}
	return web3cli
}

func LoadStorage() db.Storage {
	// Open database
	storage, err := db.NewLevelDbStorage(C.Storage.Path, false)
	Assert("Cannot open storage", err)
	log.WithField("path", C.Storage.Path).Info("Storage opened")
	return storage
}

func LoadMerkele(storage db.Storage) *merkletree.MerkleTree {
	mtstorage := storage.WithPrefix(dbMerkletreePrefix)
	mt, err := merkletree.NewMerkleTree(mtstorage, 140)
	Assert("Cannot open merkle tree", err)
	log.WithField("hash", mt.RootKey().Hex()).Info("Current root")

	return mt
}

func LoadContract(client eth.Client, jsonabifile string, address *string) *eth.Contract {
	abiFile, err := os.Open(jsonabifile)
	Assert("Cannot read contract "+jsonabifile, err)

	abi, code, err := eth.UnmarshallSolcAbiJson(abiFile)
	Assert("Cannot parse contract "+jsonabifile, err)

	var addrPtr *common.Address
	if address != nil && len(strings.TrimSpace(*address)) > 0 {
		addr := common.HexToAddress(strings.TrimSpace(*address))
		addrPtr = &addr
	}
	return eth.NewContract(client, abi, code, addrPtr)
}

func LoadRootsService(client *eth.Web3Client) rootsrv.Service {
	return rootsrv.New(LoadContract(
		client,
		C.Contracts.RootCommits.JsonABI,
		&C.Contracts.RootCommits.Address,
	))
}

func LoadIdService(client *eth.Web3Client, claimservice claimsrv.Service, storage db.Storage) identitysrv.Service {

	idstorage := storage.WithPrefix(dbIdentityPrefix)

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

	return identitysrv.New(deployerContract, implContract, proxyContract, claimservice, idstorage)
}

func LoadClaimService(mt *merkletree.MerkleTree, rootservice rootsrv.Service, ks *keystore.KeyStore, acc accounts.Account) claimsrv.Service {
	log.WithField("idAddr", C.IdAddrRaw).Info("Running claim service")
	return claimsrv.New(C.IdAddr, mt, rootservice, signsrv.New(ks, acc))
}

func LoadAdminService(mt *merkletree.MerkleTree, rootservice rootsrv.Service, claimservice claimsrv.Service) adminsrv.Service {
	return adminsrv.New(mt, rootservice, claimservice)
}
