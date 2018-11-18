package backupsrv

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/mongosrv"
	"github.com/iden3/go-iden3/utils"
	"gopkg.in/mgo.v2/bson"
)

type Service interface {
	Save(idaddr common.Address, saveBackupMsg BackupData) (uint64, error)
	RecoverAll(idaddr common.Address) ([]BackupData, error)
	RecoverByTimestamp(idaddr common.Address, timestamp uint64) ([]BackupData, error)
	RecoverByType(idaddr common.Address, dataType string) ([]BackupData, error)
}
type ServiceImpl struct {
	mongodb mongosrv.Service
}

func New(mongoservice mongosrv.Service) *ServiceImpl {
	return &ServiceImpl{mongoservice}
}

func (bs *ServiceImpl) Save(idaddr common.Address, m BackupData) (uint64, error) {
	// check ksignClaim proof (in user identity tree and in the relay tree)
	proofOfKSign, err := m.ProofOfKSignHex.Unhex()
	if err != nil {
		return 0, err
	}
	verified := claimsrv.CheckProofOfClaim(m.RelayAddr, proofOfKSign, 140)
	if !verified {
		return 0, errors.New("ProofOfKSign can not be verified")
	}

	// check saveBackupMsg.KSign match with authorizedksign from the ProofOfKSign, Leaf[64:84] is where is placed the KeyToAuthorize (KSign authorized) in the Claim data
	if !bytes.Equal(m.KSign.Bytes(), proofOfKSign.ClaimProof.Leaf[64:84]) {
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
	verified = utils.VerifySig(m.KSign, sigBytes, msgHash[:])
	if !verified {
		return 0, errors.New("signature of the data can not be verified")
	}

	// add timestamp in unixtime
	m.Timestamp = uint64(time.Now().Unix())

	// store in database
	err = bs.mongodb.GetCollections()["data"].Insert(m)
	if err != nil {
		return 0, err
	}

	return m.Timestamp, nil
}

func (bs *ServiceImpl) RecoverAll(idaddr common.Address) ([]BackupData, error) {

	// check ksignClaim proof (in user identity tree and in the relay tree)

	// check data signature

	// get from database
	var dataBackups []BackupData
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex())}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}

func (bs *ServiceImpl) RecoverByTimestamp(idaddr common.Address, timestamp uint64) ([]BackupData, error) {

	// [...] verifications

	// get from database
	var dataBackups []BackupData
	// get data with timestamp greather or equal to 'timestamp'
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "timestamp": bson.M{"$gte": timestamp}}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}
func (bs *ServiceImpl) RecoverByType(idaddr common.Address, dataType string) ([]BackupData, error) {

	// [...] verifications

	// get from database
	var dataBackups []BackupData
	// get data with timestamp greather or equal to 'timestamp'
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "type": dataType}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}
