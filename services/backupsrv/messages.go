package backupsrv

import (
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/utils"
)

type BackupData struct {
	IdAddrHex       string                       `json:"idaddrhex"`
	Data            string                       `json:"data"`
	DataSignature   string                       `json:"datasignature"`
	Type            string                       `json:"type"`
	KSignPk         *utils.PublicKey             `json:"ksignpk" binding:"required"`
	ProofOfKSignHex claimsrv.ProofOfClaimUserHex `json:"proofofksignhex"`
	RelayAddr       string                       `json:"relayaddr"`
	Version         uint64                       `json:"version"`
	Nonce           uint                         `json:"nonce"`
}

// IncrementNonce implements the method for the PoWData interface
func (bd BackupData) IncrementNonce() utils.PoWData {
	bd.Nonce++
	return bd
}
