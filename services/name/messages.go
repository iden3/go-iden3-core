package namesrv

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/merkletree"
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
	Msg struct {
		// This kind of message does not need the caducity
		RawIdentityTx
		EthID string // temp, will be calculated directly from RawIdentityTx
		Name  string
	}
	MsgSignature string // hex format
}

// MsgHash returns the Hash(VinculateIDMsg)
func (m *VinculateIDMsg) MsgHash() merkletree.Hash {
	// var b []byte
	// b = append(b, m.Msg.KRecovery_p...)
	// b = append(b, m.Msg.KRevocation_p...)
	// b = append(b, m.Msg.KSign_p...)
	// b = append(b, m.Msg.EthId...)
	// b = append(b, []byte(m.Msg.UsernameRequested)...)
	b := fmt.Sprintf("%v", m.Msg)
	return merkletree.HashBytes([]byte(b))
}
