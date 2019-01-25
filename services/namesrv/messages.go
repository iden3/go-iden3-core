package namesrv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/utils"
)

// RawIdentityTx is TODO
type RawIdentityTx struct {
	contractByteCode   []byte
	KRecovery_p        string // ecdsa.PublicKey
	KRevocation_p      string // ecdsa.PublicKey
	KSignOperational_p string // ecdsa.PublicKey
	IDRelayer          *common.Address
}

// VinculateIDMsg is the structure that contains a request to assign an
// ethereum address to a name.
type VinculateIDMsg struct {
	// This kind of message does not need the caducity
	EthAddr   common.Address         `json:"ethAddr"` // temp, will be calculated directly from RawIdentityTx
	Name      string                 `json:"name"`
	Signature *utils.SignatureEthMsg `json:"signature"` // hex format
	KSignPk   *utils.PublicKey       `json:"ksignpk" binding:"required"`
}

// Bytes returns the byte array serialization of VinculateIDMsg
func (m *VinculateIDMsg) Bytes() []byte {
	var b []byte
	b = append(b, m.EthAddr.Bytes()...)
	b = append(b, []byte(m.Name)...)
	return b
}
