package backupsrv

import (
	"bytes"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/mongosrv"
	"github.com/iden3/go-iden3/utils"
)

type Service interface {
	SaveBackup(idaddr common.Address, saveBackupMsg SaveBackupMsg) error
	RecoverBackup(idaddr common.Address) error
}
type ServiceImpl struct {
	mongodb mongosrv.Service
}

func New(mongoservice mongosrv.Service) *ServiceImpl {
	return &ServiceImpl{mongoservice}
}

func (bs *ServiceImpl) SaveBackup(idaddr common.Address, m SaveBackupMsg) error {

	// check ksignClaim proof (in user identity tree and in the relay tree)
	verified := claimsrv.CheckProofOfClaim(m.RelayAddr, m.ProofOfKSign, 140)
	if !verified {
		return errors.New("ProofOfKSign can not be verified")
	}

	// check saveBackupMsg.KSign match with authorizedksign from the ProofOfKSign, Leaf[64:84] is where is placed the KeyToAuthorize (KSign authorized) in the Claim data
	if !bytes.Equal(m.KSign.Bytes(), m.ProofOfKSign.ClaimProof.Leaf[64:84]) {
		return errors.New("KSign not equal to the ProofOfKSign.ClaimProof.Leaf[KeyToAuthorize]")
	}

	// check idaddr match with setRootClaim from the proofOfKSign, Leaf[64:84] is where is placed the idaddr in the SetRootClaim
	if !bytes.Equal(idaddr.Bytes(), m.ProofOfKSign.SetRootClaimProof.Leaf[64:84]) {
		return errors.New("idaddr don't match with the idaddr from the ProofOfKSign.SetRootClaimProof.Leaf[EthID]")
	}

	// verify data signature
	sigBytes, err := common3.HexToBytes(m.DataSignature)
	if err != nil {
		return err
	}
	sigBytes[64] -= 27
	msgHash := utils.EthHash([]byte(m.Data))
	verified = utils.VerifySig(m.KSign, sigBytes, msgHash[:])
	if !verified {
		return errors.New("signature of the data can not be verified")
	}

	// add timestamp in unixtime
	m.Timestamp = uint64(time.Now().Unix())

	// store in database

	return nil
}

func (bs *ServiceImpl) RecoverBackup(idaddr common.Address) error {
	// check data signature

	// check ksignClaim proof (in user identity tree and in the relay tree)

	// store in database

	return nil
}
