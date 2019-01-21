package core

import (
	"crypto/ecdsa"
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func NewClaimTypeNum(num uint64) *ClaimType {
	ct := ClaimType{}
	binary.BigEndian.PutUint64(ct[:], num)
	return &ct
}

var (
	// ClaimTypeBasic is a simple claim type that can be used for anything.
	ClaimTypeBasic = NewClaimTypeNum(0)
	// ClaimTypeAuthorizeKSign is a claim type to autorize a public key for signing.
	ClaimTypeAuthorizeKSign = NewClaimTypeNum(1)
	// ClaimTypeSetRootKey is a claim type of the root key of a merkle tree that goes into the relay.
	ClaimTypeSetRootKey = NewClaimTypeNum(2)
	// ClaimTypeAssignName is a claim type to assign a name to an Eth address.
	ClaimTypeAssignName = NewClaimTypeNum(3)
	// ClaimTypeAuthorizeKSignSecp256k1 is a claim type to autorize a secp256k1 public key for signing.
	ClaimTypeAuthorizeKSignSecp256k1 = NewClaimTypeNum(4)
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

// Entry serializes the claim into an Entry.
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
	// EthAddr is the assigned Ethereum Address.
	EthAddr common.Address
}

// NewClaimAssignName returns a ClaimAssignName with the name and Eth address.
func NewClaimAssignName(name string, ethAddr common.Address) *ClaimAssignName {
	c := &ClaimAssignName{}
	c.Version = 0
	hash := utils.HashBytes([]byte(name))
	copy(c.NameHash[:], hash[len(hash)-248/8:])
	c.EthAddr = ethAddr
	return c
}

// NewClaimAssignNameFromEntry deserializes a ClaimAssignName from an Entry.
func NewClaimAssignNameFromEntry(e *merkletree.Entry) *ClaimAssignName {
	c := &ClaimAssignName{}
	_, c.Version = getClaimTypeVersion(e)
	copyFromElemBytes(c.NameHash[:], 0, &e.Data[2])
	copyFromElemBytes(c.EthAddr[:], 0, &e.Data[1])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimAssignName) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[2], 0, c.NameHash[:])
	copyToElemBytes(&e.Data[1], 0, c.EthAddr[:])
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
	// Ay is the y coordinate of the elliptic curve public key.
	Ay merkletree.ElemBytes
}

// NewClaimAuthorizeKSign returns a ClaimAuthorizeKSign with the given elliptic
// public key parameters.
func NewClaimAuthorizeKSign(sign bool, ay merkletree.ElemBytes) *ClaimAuthorizeKSign {
	return &ClaimAuthorizeKSign{
		Version: 0,
		Sign:    sign,
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
	c.Ay = e.Data[2]
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimAuthorizeKSign) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	sign := []byte{0}
	if c.Sign {
		sign = []byte{1}
	}
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, sign)
	e.Data[2] = c.Ay
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAuthorizeKSign) Type() ClaimType {
	return *ClaimTypeAuthorizeKSign
}

// ClaimAuthorizeKSignSecp256k1 is a claim to autorize a public key for signing.
type ClaimAuthorizeKSignSecp256k1 struct {
	// Version is the claim version.
	Version uint32
	// PubKey is the ECDSA public key.
	PubKey *ecdsa.PublicKey
}

// NewClaimAuthorizeKSignSecp256k1 returns a ClaimAuthorizeKSignSecp256k1 with the given elliptic
// public key parameters.
func NewClaimAuthorizeKSignSecp256k1(pk *ecdsa.PublicKey) *ClaimAuthorizeKSignSecp256k1 {
	return &ClaimAuthorizeKSignSecp256k1{
		Version: 0,
		PubKey:  pk,
	}
}

// NewClaimAuthorizeKSignSecp256k1FromEntry deserializes a ClaimAuthorizeKSignSecp256k1 from an Entry.
func NewClaimAuthorizeKSignSecp256k1FromEntry(e *merkletree.Entry) (*ClaimAuthorizeKSignSecp256k1, error) {
	c := &ClaimAuthorizeKSignSecp256k1{}
	_, c.Version = getClaimTypeVersion(e)
	var cpk [33]byte
	copyFromElemBytes(cpk[len(cpk)-2:], ClaimTypeVersionLen, &e.Data[3])
	copyFromElemBytes(cpk[:len(cpk)-2], 0, &e.Data[2])
	var err error
	c.PubKey, err = crypto.DecompressPubkey(cpk[:])
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Entry serializes the claim into an Entry.
func (c *ClaimAuthorizeKSignSecp256k1) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	cpk := crypto.CompressPubkey(c.PubKey)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, cpk[len(cpk)-2:])
	copyToElemBytes(&e.Data[2], 0, cpk[:len(cpk)-2])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAuthorizeKSignSecp256k1) Type() ClaimType {
	return *ClaimTypeAuthorizeKSignSecp256k1
}

// ClaimSetRootKey is a claim of the root key of a merkle tree that goes into the relay.
type ClaimSetRootKey struct {
	// Version is the claim version.
	Version uint32
	// Era is used for labeling epochs.
	Era uint32
	// EthAddr is the Ethereum Address related to the root key.
	EthAddr common.Address
	// RootKey is the root of the mekrlee tree.
	RootKey merkletree.Hash
}

// NewClaimSetRootKey returns a ClaimSetRootKey with the given Eth ID and
// merklee tree root key.
func NewClaimSetRootKey(ethAddr common.Address, rootKey merkletree.Hash) *ClaimSetRootKey {
	return &ClaimSetRootKey{
		Version: 0,
		Era:     0,
		EthAddr: ethAddr,
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
	copyFromElemBytes(c.EthAddr[:], 0, &e.Data[2])
	c.RootKey = merkletree.Hash(e.Data[1])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimSetRootKey) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	var era [32 / 8]byte
	binary.BigEndian.PutUint32(era[:], c.Era)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, era[:])
	copyToElemBytes(&e.Data[2], 0, c.EthAddr[:])
	e.Data[1] = merkletree.ElemBytes(c.RootKey)
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
	case *ClaimTypeAuthorizeKSignSecp256k1:
		return NewClaimAuthorizeKSignSecp256k1FromEntry(e)
	default:
		return nil, ErrInvalidClaimType
	}
}
