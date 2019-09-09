package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/merkletree"
)

// ClaimEthId is a claim to authorize an ethereum address for the identity. The address can be of a counterfactual smart contract, or a direct address from a private key
type ClaimEthId struct {
	// Version is the claim version
	Version uint32

	// Addr is the EthId that will use this identity in the ethereum blockchain
	Address common.Address

	// IdentityFactory specifies that the ClaimEthId.Address is an smartcontract, and how this identity is created. It can be just an identitied of the method, or an smartcontract that creates the identity.
	// IDEN3 specifies the contract address 0x9827348723984729834234 as the factory for its contrafactual identities
	// If 0x000.0000, means that is not using an identity creator, and the identity is always available.
	IdentityFactory common.Address
}

// NewClaimEthId returns a ClaimEthId
func NewClaimEthId(addr, identityFactory common.Address) *ClaimEthId {
	return &ClaimEthId{
		Version:         0,
		Address:         addr,
		IdentityFactory: identityFactory,
	}
}

// NewClaimEthId deserializes a ClaimEthId from an Entry.
func NewClaimEthIdFromEntry(e *merkletree.Entry) *ClaimEthId {
	c := &ClaimEthId{}
	_, c.Version = GetClaimTypeVersion(e)
	copyFromElemBytes(c.Address[:], 0, &e.Data[2])
	copyFromElemBytes(c.IdentityFactory[:], 0, &e.Data[1])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimEthId) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	SetClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[2], 0, c.Address[:])
	copyToElemBytes(&e.Data[1], 0, c.IdentityFactory[:])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimEthId) Type() ClaimType {
	return *ClaimTypeEthId
}
