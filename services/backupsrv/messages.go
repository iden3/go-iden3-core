package backupsrv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/services/claimsrv"
)

type SaveBackupMsg struct {
	Data          string
	DataSignature string
	KSign         common.Address
	ProofOfKSign  claimsrv.ProofOfClaim
	RelayAddr     common.Address
	Timestamp     uint64
}