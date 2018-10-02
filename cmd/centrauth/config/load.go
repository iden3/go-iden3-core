package config

import (
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
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

func LoadClaimService(mt *merkletree.MerkleTree, rootsrv rootsrv.Service, ks *keystore.KeyStore, acc accounts.Account) claimsrv.Service {
	return claimsrv.New(mt, rootsrv, signsrv.New(ks, acc))
}
