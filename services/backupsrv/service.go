package backupsrv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/mongosrv"
	"github.com/iden3/go-iden3/utils"
	"gopkg.in/mgo.v2/bson"
)

type Service interface {
	GetPoWDifficulty() int
	GetLastVersion(idaddr common.Address) (uint64, error)
	Save(idaddr common.Address, saveBackupMsg BackupData) (uint64, error)
	RecoverAll(idaddr common.Address) ([]BackupData, error)
	RecoverSinceVersion(idaddr common.Address, version uint64) ([]BackupData, error)
	RecoverByType(idaddr common.Address, dataType string) ([]BackupData, error)
	RecoverSinceVersionByType(idaddr common.Address, version uint64, dataType string) ([]BackupData, error)
}

type ServiceImpl struct {
	mongodb       mongosrv.Service
	powDifficulty int
}

func New(mongoservice mongosrv.Service, powDiff int) *ServiceImpl {
	return &ServiceImpl{mongoservice, powDiff}
}

// GetPoWDifficulty returns the configured Proof-of-Work difficulty, setted in the config file of the backupserver
func (bs *ServiceImpl) GetPoWDifficulty() int {
	return bs.powDifficulty
}

func (bs *ServiceImpl) GetLastVersion(idaddr common.Address) (uint64, error) {
	var currVer BackupData
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "ksign": "currentversion"}).One(&currVer)
	return currVer.Version, err
}

// Save verifies the proofs for auth, and stores the data packet in the database
func (bs *ServiceImpl) Save(idaddr common.Address, m BackupData) (uint64, error) {
	// check PoW
	b, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}
	hash := utils.HashBytes(b)
	if !utils.CheckPoW(hash, bs.GetPoWDifficulty()) {
		return 0, errors.New("PoW not passed")
	}

	// check ksignClaim proof (in user identity tree and in the relay tree)
	proofOfKSign, err := m.ProofOfKSignHex.Unhex()
	if err != nil {
		return 0, err
	}
	kSign := common.HexToAddress(m.KSign)
	relayAddr := common.HexToAddress(m.RelayAddr)
	verified := claimsrv.CheckProofOfClaim(relayAddr, proofOfKSign, 140)
	if !verified {
		return 0, errors.New("ProofOfKSign can not be verified")
	}

	// check saveBackupMsg.KSign match with authorizedksign from the ProofOfKSign, Leaf[64:84] is where is placed the KeyToAuthorize (KSign authorized) in the Claim data
	if !bytes.Equal(kSign.Bytes(), proofOfKSign.ClaimProof.Leaf[64:84]) {
		return 0, errors.New("KSign not equal to the ProofOfKSign.ClaimProof.Leaf[KeyToAuthorize]")
	}

	// check idaddr match with setRootClaim from the proofOfKSign, Leaf[64:84] is where is placed the idaddr in the SetRootClaim
	if !bytes.Equal(idaddr.Bytes(), proofOfKSign.SetRootClaimProof.Leaf[64:84]) {
		return 0, errors.New("idaddr don't match with the idaddr from the ProofOfKSign.SetRootClaimProof.Leaf[EthID]")
	}

	// verify data signature
	sigBytes, err := common3.HexToBytes(m.DataSignature)
	if err != nil {
		return 0, err
	}
	sigBytes[64] -= 27
	msgHash := utils.EthHash([]byte(m.Data))
	verified = utils.VerifySig(kSign, sigBytes, msgHash[:])
	if !verified {
		return 0, errors.New("signature of the data can not be verified")
	}

	// check version (check that the current version is == lastversion+1)
	var aux BackupData
	err = bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "version": m.Version}).One(&aux)
	if err == nil {
		// if data exists, the given version is not valid
		return m.Version, errors.New("given version not valid")
	}

	// store in database
	err = bs.mongodb.GetCollections()["data"].Insert(m)
	if err != nil {
		return 0, err
	}

	// TODO store in leveldb instead of mongodb. key: idaddr+version, value: type+dataencrypted
	// the currentVersion will be stored as key: idaddr+"currver", value: currver
	currVer := BackupData{
		IdAddrHex:       idaddr.Hex(),
		Data:            "",
		DataSignature:   "",
		Type:            "",
		KSign:           "currentversion",
		ProofOfKSignHex: claimsrv.ProofOfClaimHex{},
		RelayAddr:       "",
		Version:         m.Version,
		Nonce:           0,
	}
	err = bs.mongodb.GetCollections()["data"].Update(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "ksign": "currentversion"}, currVer)
	return m.Version, nil
}

// RecoverAll returns all the data packets stored by an idaddr
func (bs *ServiceImpl) RecoverAll(idaddr common.Address) ([]BackupData, error) {

	// TODO auth verifications
	// check ksignClaim proof (in user identity tree and in the relay tree)

	// get from database
	var dataBackups []BackupData
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex())}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}

// RecoverSinceVersion returns all the data packets stored by an idaddr since after the version specified in the parameter
func (bs *ServiceImpl) RecoverSinceVersion(idaddr common.Address, version uint64) ([]BackupData, error) {
	color.Blue("version")
	fmt.Println(version)

	// TODO auth verifications

	// get from database
	var dataBackups []BackupData
	// get data with version greather or equal to 'version'
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "version": bson.M{"$gt": version}}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}

// RecoverByType returns all the data packets stored by an idaddr with the requested type
func (bs *ServiceImpl) RecoverByType(idaddr common.Address, dataType string) ([]BackupData, error) {

	// TODO auth verifications

	// get from database
	var dataBackups []BackupData
	// get data by type
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "type": dataType}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}

// RecoverSinceVersionByType returns all the data packets stored by an idaddr with the requested type since after the version specified in the parameter
func (bs *ServiceImpl) RecoverSinceVersionByType(idaddr common.Address, version uint64, dataType string) ([]BackupData, error) {

	// TODO auth verifications

	// get from database
	var dataBackups []BackupData
	// get data by type
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "version": bson.M{"$gt": version}, "type": dataType}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}
