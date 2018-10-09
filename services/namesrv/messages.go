package namesrv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/merkletree"
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
	EthID     common.Address `json:"ethID"` // temp, will be calculated directly from RawIdentityTx
	Name      string         `json:"name"`
	Signature string         `json:"signature"` // hex format
}

// MsgHash returns the Hash(VinculateIDMsg)
func (m *VinculateIDMsg) MsgHash() merkletree.Hash {
	var b []byte
	b = append(b, m.EthID.Bytes()...)
	b = append(b, []byte(m.Name)...)
	return utils.EthHash(b)
}
