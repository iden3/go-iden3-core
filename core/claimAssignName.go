package core

import (
	"github.com/iden3/go-iden3-core/merkletree"
)

// ClaimAssignName is a claim to assign a name to an id.
type ClaimAssignName struct {
	// Version is the claim version.
	Version uint32
	// NameHash is the hash of the name.
	NameHash [248 / 8]byte
	// Id is the assigned ID
	Id ID
}

// NewClaimAssignName returns a ClaimAssignName with the name and Id.
func NewClaimAssignName(name string, id ID) *ClaimAssignName {
	c := &ClaimAssignName{}
	c.Version = 0
	c.NameHash = HashString(name)
	c.Id = id
	return c
}

// NewClaimAssignNameFromEntry deserializes a ClaimAssignName from an Entry.
func NewClaimAssignNameFromEntry(e *merkletree.Entry) *ClaimAssignName {
	c := &ClaimAssignName{}
	_, c.Version = getClaimTypeVersion(e)
	copyFromElemBytes(c.NameHash[:], 0, &e.Data[2])
	copyFromElemBytes(c.Id[:], 0, &e.Data[1])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimAssignName) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[2], 0, c.NameHash[:])
	copyToElemBytes(&e.Data[1], 0, c.Id[:31])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAssignName) Type() ClaimType {
	return *ClaimTypeAssignName
}
