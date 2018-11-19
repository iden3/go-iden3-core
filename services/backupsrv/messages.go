package backupsrv

import (
	"github.com/iden3/go-iden3/services/claimsrv"
)

type BackupData struct {
	IdAddrHex       string
	Data            string
	DataSignature   string
	Type            string
	KSign           string
	ProofOfKSignHex claimsrv.ProofOfClaimHex
	RelayAddr       string
	Timestamp       uint64
}
