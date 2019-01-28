package backupsrv

import (
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/utils"
)

type BackupDataMsg struct {
	IdAddrHex       string                       `json:"idaddrhex" binding:"required"`
	Data            string                       `json:"data" binding:"required"`
	DataSignature   string                       `json:"datasignature" binding:"required"`
	Type            string                       `json:"type" binding:"required"`
	KSignPk         *utils.PublicKey             `json:"ksignpk" binding:"required"`
	ProofOfKSignHex claimsrv.ProofOfClaimUserHex `json:"proofofksignhex" binding:"required"`
	RelayAddr       string                       `json:"relayaddr" binding:"required"`
	Version         uint64                       `json:"version" binding:"required"`
	Nonce           uint                         `json:"nonce" binding:"required"`
}

// IncrementNonce implements the method for the PoWData interface
func (bd BackupDataMsg) IncrementNonce() utils.PoWData {
	bd.Nonce++
	return bd
}
