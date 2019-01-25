package centrauthsrv

import (
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/claimsrv"
)

type AuthMsg struct {
	Address    string                       `json:"address" binding:"required"`
	Challenge  string                       `json:"challenge" binding:"required"`
	Signature  string                       `json:"signature" binding:"required"`
	KSignPk    *common3.PublicKey           `json:"ksignpk" binding:"required"`
	KSignProof claimsrv.ProofOfClaimUserHex `json:"ksignProof" binding:"required"`
}

type AuthTokenMsg struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}
