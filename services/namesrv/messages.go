package namesrv

import (
	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/utils"
)

// RawIdentityTx
type RawIdentityTx struct {
	contractByteCode   []byte
	KRecovery_p        string // ecdsa.PublicKey
	KRevocation_p      string // ecdsa.PublicKey
	KSignOperational_p string // ecdsa.PublicKey
	IDRelayer          *common.Address
}

// VinculateIDMsg is the structure that contains
type VinculateIDMsg struct {
	// This kind of message does not need the caducity
	EthAddr   common.Address     `json:"ethAddr"` // temp, will be calculated directly from RawIdentityTx
	Name      string             `json:"name"`
	Signature string             `json:"signature"` // hex format
	KSignPk   *common3.PublicKey `json:"ksignpk" binding:"required"`
}

// MsgHash returns the Hash(VinculateIDMsg)
func (m *VinculateIDMsg) MsgHash() utils.Hash {
	var b []byte
	b = append(b, m.EthAddr.Bytes()...)
	b = append(b, []byte(m.Name)...)
	return utils.EthHash(b)
}
