package backupsrv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/services/claimsrv"
)

type BackupData struct {
	IdAddrHex       string
	Data            string
	DataSignature   string
	Type            string
	KSign           common.Address
	ProofOfKSignHex claimsrv.ProofOfClaimHex
	RelayAddr       common.Address
	Timestamp       uint64
}
