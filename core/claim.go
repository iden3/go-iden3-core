package core

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

// ErrInvalidClaimType indicates a type error when parsing an Entry into a claim.
var ErrInvalidClaimType = errors.New("invalid claim type")

// copyToElemBytes copies the src slice forwards to e, ending at -start of
// e.  This function will panic if src doesn't fit into len(e)-start.
func copyToElemBytes(e *merkletree.ElemBytes, start int, src []byte) {
	copy(e[merkletree.ElemBytesLen-start-len(src):], src)
}

// copyFromElemBytes copies from e to dst, ending at -start of e and going
// forwards.  This function will panic if len(e)-start is smaller than
// len(dst).
func copyFromElemBytes(dst []byte, start int, e *merkletree.ElemBytes) {
	copy(dst, e[merkletree.ElemBytesLen-start-len(dst):])
}

// setClaimTypeVersion is a helper function to set the type and version of a
// claim.
func setClaimTypeVersion(e *merkletree.Entry, claimType ClaimType, version uint32) {
	SetClaimTypeVersionInData(&e.Data, claimType, version)
}

func SetClaimTypeVersionInData(d *merkletree.Data, claimType ClaimType, version uint32) {
	copyToElemBytes(&d[3], 0, claimType[:])
	binary.BigEndian.PutUint32(d[3][merkletree.ElemBytesLen-ClaimTypeVersionLen:], version)
}

// getClaimTypeVersion is a helper function to get the type and version from a
// claim.
func getClaimTypeVersion(e *merkletree.Entry) (c ClaimType, v uint32) {
	return GetClaimTypeVersionFromData(&e.Data)
}

// GetClaimTypeVersionFromData(
func GetClaimTypeVersionFromData(d *merkletree.Data) (c ClaimType, v uint32) {
	copyFromElemBytes(c[:], 0, &d[3])
	v = binary.BigEndian.Uint32(d[3][merkletree.ElemBytesLen-ClaimTypeVersionLen:])
	return c, v
}

// ClaimTypeLen is the length in bytes of the type in a claim.
const ClaimTypeLen = 64 / 8

// ClaimType is the type used to store a claim type.
type ClaimType [ClaimTypeLen]byte

// NewClaimType creates a ClaimType from a type name.
func NewClaimType(name string) *ClaimType {
	t := &ClaimType{}
	h := utils.HashBytes([]byte(name))
	copy(t[:ClaimTypeLen], h[len(h)-ClaimTypeLen:])
	return t
}

var (
	// ClaimTypeBasic is a simple claim type that can be used for anything.
	ClaimTypeBasic = NewClaimType("iden3.claim.basic")
	// ClaimTypeAssignName is a claim type to assign a name to an Eth address.
	ClaimTypeAssignName = NewClaimType("iden3.claim.assign_name")
	// ClaimTypeAuthorizeKSign is a claim type to autorize a public key for signing.
	ClaimTypeAuthorizeKSign = NewClaimType("iden3.claim.authorize_k_sign")
	// ClaimTypeSetRootKey is a claim type of the root key of a merkle tree that goes into the relay.
	ClaimTypeSetRootKey = NewClaimType("iden3.claim.set_root_key")
)

// ClaimVersionLen is the length in bytes of the version in a claim.
const ClaimVersionLen = 32 / 8

// ClaimTypeVersionLen is the length in bytes of the version and length in a claim.
const ClaimTypeVersionLen = ClaimTypeLen + ClaimVersionLen

// ClaimBasic is a simple claim that can be used for anything.
type ClaimBasic struct {
	// Version is the claim version.
	Version uint32
	// IndexSlot is data that goes into the remaining space used for the index.
	IndexSlot [400 / 8]byte
	// DataSlot is the data that goes into the remaining space not used for the index.
	DataSlot [496 / 8]byte
}

// NewClaimBasic returns a ClaimBasic with the provided data.
func NewClaimBasic(indexSlot [400 / 8]byte, dataSlot [496 / 8]byte) *ClaimBasic {
	return &ClaimBasic{
		Version:   0,
		IndexSlot: indexSlot,
		DataSlot:  dataSlot,
	}
}

// NewClaimBasicFromEntry deserializes a ClaimBasic from an Entry.
func NewClaimBasicFromEntry(e *merkletree.Entry) *ClaimBasic {
	c := &ClaimBasic{}
	_, c.Version = getClaimTypeVersion(e)
	copyFromElemBytes(c.IndexSlot[len(c.IndexSlot)-152/8:], ClaimTypeVersionLen, &e.Data[3])
	copyFromElemBytes(c.IndexSlot[:248/8], 0, &e.Data[2])
	copyFromElemBytes(c.DataSlot[248/8:], 0, &e.Data[1])
	copyFromElemBytes(c.DataSlot[:248/8], 0, &e.Data[0])
	return c
}

// ToEntry serializes the claim into an Entry.
func (c *ClaimBasic) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, c.IndexSlot[len(c.IndexSlot)-152/8:])
	copyToElemBytes(&e.Data[2], 0, c.IndexSlot[:248/8])
	copyToElemBytes(&e.Data[1], 0, c.DataSlot[248/8:])
	copyToElemBytes(&e.Data[0], 0, c.DataSlot[:248/8])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimBasic) Type() ClaimType {
	return *ClaimTypeBasic
}

// ClaimAssignName is a claim to assign a name to an Eth address.
type ClaimAssignName struct {
	// Version is the claim version.
	Version uint32
	// NameHash is the hash of the name.
	NameHash [248 / 8]byte
	// EthID is the assigned Ethereum ID.
	EthID common.Address
}

// NewClaimAssignName returns a ClaimAssignName with the name and Eth address.
func NewClaimAssignName(name string, ethID common.Address) *ClaimAssignName {
	c := &ClaimAssignName{}
	c.Version = 0
	hash := utils.HashBytes([]byte(name))
	copy(c.NameHash[:], hash[len(hash)-248/8:])
	c.EthID = ethID
	return c
}

// NewClaimAssignNameFromEntry deserializes a ClaimAssignName from an Entry.
func NewClaimAssignNameFromEntry(e *merkletree.Entry) *ClaimAssignName {
	c := &ClaimAssignName{}
	_, c.Version = getClaimTypeVersion(e)
	copyFromElemBytes(c.NameHash[:], 0, &e.Data[2])
	copyFromElemBytes(c.EthID[:], 0, &e.Data[1])
	return c
}

// ToEntry serializes the claim into an Entry.
func (c *ClaimAssignName) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[2], 0, c.NameHash[:])
	copyToElemBytes(&e.Data[1], 0, c.EthID[:])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAssignName) Type() ClaimType {
	return *ClaimTypeAssignName
}

// ClaimAuthorizeKSign is a claim to autorize a public key for signing.
type ClaimAuthorizeKSign struct {
	// Version is the claim version.
	Version uint32
	// Sign means positive if false, negative if true.
	Sign bool
	// Ax is the x coordinate of the elliptic curve public key.
	Ax [128 / 8]byte
	// Ay is the x coordinate of the elliptic curve public key.
	Ay [128 / 8]byte
}

// NewClaimAuthorizeKSign returns a ClaimAuthorizeKSign with the given elliptic
// public key parameters.
func NewClaimAuthorizeKSign(sign bool, ax, ay [128 / 8]byte) *ClaimAuthorizeKSign {
	return &ClaimAuthorizeKSign{
		Version: 0,
		Sign:    sign,
		Ax:      ax,
		Ay:      ay,
	}
}

// NewClaimAuthorizeKSign deserializes a ClaimAuthorizeKSign from an Entry.
func NewClaimAuthorizeKSignFromEntry(e *merkletree.Entry) *ClaimAuthorizeKSign {
	c := &ClaimAuthorizeKSign{}
	_, c.Version = getClaimTypeVersion(e)
	sign := []byte{0}
	copyFromElemBytes(sign, ClaimTypeVersionLen, &e.Data[3])
	if sign[0] == 1 {
		c.Sign = true
	}
	copyFromElemBytes(c.Ax[:], ClaimTypeVersionLen+1, &e.Data[3])
	copyFromElemBytes(c.Ay[:], 0, &e.Data[2])
	return c
}

// ToEntry serializes the claim into an Entry.
func (c *ClaimAuthorizeKSign) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	sign := []byte{0}
	if c.Sign {
		sign = []byte{1}
	}
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, sign)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen+1, c.Ax[:])
	copyToElemBytes(&e.Data[2], 0, c.Ay[:])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAuthorizeKSign) Type() ClaimType {
	return *ClaimTypeAuthorizeKSign
}

// ClaimSetRootKey is a claim of the root key of a merkle tree that goes into the relay.
type ClaimSetRootKey struct {
	// Version is the claim version.
	Version uint32
	// Era is used for labeling epochs.
	Era uint32
	// EthID is the Ethereum ID related to the root key.
	EthID common.Address
	// RootKey is the root of the mekrlee tree.
	RootKey merkletree.Hash
}

// NewClaimSetRootKey returns a ClaimSetRootKey with the given Eth ID and
// merklee tree root key.
func NewClaimSetRootKey(ethID common.Address, rootKey merkletree.Hash) *ClaimSetRootKey {
	return &ClaimSetRootKey{
		Version: 0,
		Era:     0,
		EthID:   ethID,
		RootKey: rootKey,
	}
}

// NewClaimSetRootKey deserializes a ClaimSetRootKey from an Entry.
func NewClaimSetRootKeyFromEntry(e *merkletree.Entry) *ClaimSetRootKey {
	c := &ClaimSetRootKey{}
	_, c.Version = getClaimTypeVersion(e)
	var era [32 / 8]byte
	copyFromElemBytes(era[:], ClaimTypeVersionLen, &e.Data[3])
	c.Era = binary.BigEndian.Uint32(era[:])
	copyFromElemBytes(c.EthID[:], 0, &e.Data[2])
	copyFromElemBytes(c.RootKey[:], 0, &e.Data[1])
	return c
}

// ToEntry serializes the claim into an Entry.
func (c *ClaimSetRootKey) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	var era [32 / 8]byte
	binary.BigEndian.PutUint32(era[:], c.Era)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, era[:])
	copyToElemBytes(&e.Data[2], 0, c.EthID[:])
	copyToElemBytes(&e.Data[1], 0, c.RootKey[:])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimSetRootKey) Type() ClaimType {
	return *ClaimTypeSetRootKey
}

// NewClaimFromEntry deserializes a valid claim type into a Claim.
func NewClaimFromEntry(e *merkletree.Entry) (merkletree.Entrier, error) {
	claimType, _ := getClaimTypeVersion(e)
	switch claimType {
	case *ClaimTypeBasic:
		c := NewClaimBasicFromEntry(e)
		return c, nil
	case *ClaimTypeAssignName:
		c := NewClaimAssignNameFromEntry(e)
		return c, nil
	case *ClaimTypeAuthorizeKSign:
		c := NewClaimAuthorizeKSignFromEntry(e)
		return c, nil
	case *ClaimTypeSetRootKey:
		c := NewClaimSetRootKeyFromEntry(e)
		return c, nil
	default:
		return nil, ErrInvalidClaimType
	}
}
