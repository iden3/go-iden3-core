package core

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/merkletree"
)

// EthKeyTypeLen is the length in bytes of the type in a claim
const EthKeyTypeLen = 32 / 8

// ClaimType is the type used to store a claim type
type EthKeyType [EthKeyTypeLen]byte

// NewClaimTypeNum to set type to a claim.
func NewEthKeyType(num uint32) EthKeyType {
	ct := EthKeyType{}
	binary.BigEndian.PutUint32(ct[:], num)
	return ct
}

var (
	// EthKeyTypeDisable specifies a Ethereum Key (Addr) that is allowed to Disable the ID
	EthKeyTypeDisable = NewEthKeyType(0)
	// EthKeyTypeReenable specifies a Ethereum Key (Addr) that is allowed to Reenable the ID
	EthKeyTypeReenable = NewEthKeyType(1)
	// EthKeyTypeUpgrade specifies a Ethereum Key (Addr) that is allowed to Upgrade the ID
	EthKeyTypeUpgrade = NewEthKeyType(2)
	// EthKeyTypeUpdateRoot specifies a Ethereum Key (Addr) that is allowed to Update the Root in roots smart contract in name of the ID
	EthKeyTypeUpdateRoot = NewEthKeyType(3)
)

// ClaimAuthEthKey is a claim type to authorize an Eth Address directly from a private key, allowing to specify if is used as KDisable (revoke), KReenable (recover), etc
type ClaimAuthEthKey struct {
	// Version is the claim version
	Version uint32
	// EthKey is the ethereum address of the Key that is being authorized
	EthKey common.Address
	// EthKeyType specifies the type of the EthKey, for which use is authorized
	EthKeyType uint32
}

// NewClaimAuthEthKey returns a ClaimAuthEthKey
func NewClaimAuthEthKey(ethKey common.Address, typ EthKeyType) *ClaimAuthEthKey {
	return &ClaimAuthEthKey{
		Version:    0,
		EthKey:     ethKey,
		EthKeyType: binary.BigEndian.Uint32(typ[:]),
	}
}

// NewClaimAuthEthKey deserializes a ClaimAuthEthKey from an Entry
func NewClaimAuthEthKeyFromEntry(e *merkletree.Entry) *ClaimAuthEthKey {
	c := &ClaimAuthEthKey{}
	_, c.Version = GetClaimTypeVersion(e)
	copyFromElemBytes(c.EthKey[:], 0, &e.Data[2])
	var typ [EthKeyTypeLen]byte
	copyFromElemBytes(typ[:], 20, &e.Data[2])
	c.EthKeyType = binary.BigEndian.Uint32(typ[:])
	return c
}

// Entry serializes the claim into an Entry
func (c *ClaimAuthEthKey) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	SetClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[2], 0, c.EthKey[:])
	var typ [EthKeyTypeLen]byte
	binary.BigEndian.PutUint32(typ[:], c.EthKeyType)
	copyToElemBytes(&e.Data[2], 20, typ[:])
	return e
}

// Type returns the ClaimType of the claim
func (c *ClaimAuthEthKey) Type() ClaimType {
	return *ClaimTypeAuthEthKey
}
