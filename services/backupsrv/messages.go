package backupsrv

import (
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/utils"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type BackupPacket struct {
	Username string `json:"username" binding:"required"`
	Backup   string `json:"backup" binding:"required"`
}

type BackupDataMsg struct {
	IdAddrHex       string                       `json:"idAddrHex" binding:"required"`
	Data            string                       `json:"data" binding:"required"`
	DataSignature   string                       `json:"dataSignature" binding:"required"`
	Type            string                       `json:"type" binding:"required"`
	KSignPk         *utils.PublicKey             `json:"ksignpk" binding:"required"`
	ProofOfKSignHex claimsrv.ProofOfClaimUserHex `json:"proofKSignHex" binding:"required"`
	RelayAddr       string                       `json:"relayAddr" binding:"required"`
	Version         uint64                       `json:"version" binding:"required"`
	Nonce           uint                         `json:"nonce" binding:"required"`
}

// IncrementNonce implements the method for the PoWData interface
func (bd BackupDataMsg) IncrementNonce() utils.PoWData {
	bd.Nonce++
	return bd
}
