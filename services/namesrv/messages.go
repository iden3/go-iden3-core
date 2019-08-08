package namesrv

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/utils"
)

// RawIdentityTx is TODO
type RawIdentityTx struct {
	contractByteCode   []byte
	KRecovery_p        string // ecdsa.PublicKey
	KRevocation_p      string // ecdsa.PublicKey
	KSignOperational_p string // ecdsa.PublicKey
	IdRelayer          *core.ID
}

// VinculateIdMsg is the structure that contains a request to assign an
// ethereum address to a name.
type VinculateIdMsg struct {
	// This kind of message does not need the caducity
	Id         core.ID                `json:"id" binding:"required"` // temp, will be calculated directly from RawIdentityTx
	Name       string                 `json:"name" binding:"required"`
	Signature  *utils.SignatureEthMsg `json:"signature" binding:"required"` // hex format
	KSignPk    *utils.PublicKey       `json:"kSignPk" binding:"required"`
	ProofKSign core.ProofClaim        `json:"proofKSign" binding:"required"`
}

// Bytes returns the byte array serialization of VinculateIdMsg
func (m *VinculateIdMsg) Bytes() []byte {
	var b []byte
	b = append(b, m.Id.Bytes()...)
	b = append(b, []byte(m.Name)...)
	return b
}
