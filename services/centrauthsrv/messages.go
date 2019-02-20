package centrauthsrv

import (
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/utils"
)

type AuthMsg struct {
	Address    string                     `json:"address" binding:"required"`
	Challenge  string                     `json:"challenge" binding:"required"`
	Signature  string                     `json:"signature" binding:"required"`
	KSignPk    *utils.PublicKey           `json:"ksignpk" binding:"required"`
	KSignProof claimsrv.ProofClaimUserHex `json:"ksignProof" binding:"required"`
}

type AuthTokenMsg struct {
	Success bool   `json:"success" binding:"required"`
	Token   string `json:"token" binding:"required"`
}
