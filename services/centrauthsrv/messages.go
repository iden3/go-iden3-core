package centrauthsrv

import "github.com/iden3/go-iden3/services/claimsrv"

type AuthMsg struct {
	Address    string                   `json:"address" binding:"required"`
	Challenge  string                   `json:"challenge" binding:"required"`
	Signature  string                   `json:"signature" binding:"required"`
	KSign      string                   `json:"ksign" binding:"required"`
	KSignProof claimsrv.ProofOfClaimHex `json:"ksignProof" binding:"required"`
}

type AuthTokenMsg struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}
