package backupsrv

import (
	"github.com/iden3/go-iden3/utils"
	"github.com/iden3/go-iden3/services/claimsrv"
)

type BackupData struct {
	IdAddrHex       string	`json:"idaddrhex"`
	Data            string	`json:"data"`
	DataSignature   string	`json:"datasignature"`
	Type            string	`json:"type"`
	KSign           string	`json:"ksign"`
	ProofOfKSignHex claimsrv.ProofOfClaimHex	`json:"proofofksignhex"`
	RelayAddr       string	`json:"relayaddr"`
	Timestamp       uint64	`json:"timestamp"`
	Nonce           int	`json:"nonce"`
}

// IncrementNonce implements the method for the PoWData interface
func (bd BackupData) IncrementNonce() utils.PoWData {
	bd.Nonce++
	return bd
}