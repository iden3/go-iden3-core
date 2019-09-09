package core

import (
	"encoding/binary"
	"errors"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/utils"
	cryptoConstants "github.com/iden3/go-iden3-crypto/constants"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
)

// ErrInvalidClaimType indicates a type error when parsing an Entry into a claim.
var ErrInvalidClaimType = errors.New("invalid claim type")

// copyToElemBytes copies the src slice forwards to e, ending at -start of
// e.  This function will panic if src doesn't fit into len(e)-start.
func copyToElemBytes(e *merkletree.ElemBytes, start int, src []byte) {
	copy(e[merkletree.ElemBytesLen-start-len(src):], src)
}

// ClearMostSigByte sets the most significant byte of the element to 0 to make sure it fits
// inside the FiniteField over R.
func ClearMostSigByte(e [256 / 8]byte) merkletree.ElemBytes {
	e[0] = 0
	return merkletree.ElemBytes(e)
}

// copyFromElemBytes copies from e to dst, ending at -start of e and going
// forwards.  This function will panic if len(e)-start is smaller than
// len(dst).
func copyFromElemBytes(dst []byte, start int, e *merkletree.ElemBytes) {
	copy(dst, e[merkletree.ElemBytesLen-start-len(dst):])
}

// SetClaimTypeVersion is a helper function to set the type and version of a
// claim.
func SetClaimTypeVersion(e *merkletree.Entry, claimType ClaimType, version uint32) {
	SetClaimTypeVersionInData(&e.Data, claimType, version)
}

// SetClaimTypeVersionInData is a helper function to set the type and version of a
// claim.
func SetClaimTypeVersionInData(d *merkletree.Data, claimType ClaimType, version uint32) {
	copyToElemBytes(&d[3], 0, claimType[:])
	binary.BigEndian.PutUint32(d[3][merkletree.ElemBytesLen-ClaimTypeVersionLen:], version)
}

// GetClaimTypeVersion is a helper function to get the type and version from a
// claim.
func GetClaimTypeVersion(e *merkletree.Entry) (c ClaimType, v uint32) {
	return GetClaimTypeVersionFromData(&e.Data)
}

// GetClaimTypeVersionFromData gets claims fields data and version from a given claim.
func GetClaimTypeVersionFromData(d *merkletree.Data) (c ClaimType, v uint32) {
	copyFromElemBytes(c[:], 0, &d[3])
	v = binary.BigEndian.Uint32(d[3][merkletree.ElemBytesLen-ClaimTypeVersionLen:])
	return c, v
}

// HashString takes the first 31 bytes of a hash applied to string
func HashString(s string) (stringHashed [248 / 8]byte) {
	hash := utils.HashBytes([]byte(s))
	copy(stringHashed[:], hash[len(hash)-248/8:])
	return stringHashed
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

// NewClaimTypeNum to set type to a claim.
func NewClaimTypeNum(num uint64) *ClaimType {
	ct := ClaimType{}
	binary.BigEndian.PutUint64(ct[:], num)
	return &ct
}

var (
	// ClaimTypeBasic is a simple claim type that can be used for anything.
	ClaimTypeBasic = NewClaimTypeNum(0)
	// ClaimTypeAuthorizeKSignBabyJub is a claim type to autorize a babyjub public key for signing.
	ClaimTypeAuthorizeKSignBabyJub = NewClaimTypeNum(1)
	// ClaimTypeSetRootKey is a claim type of the root key of a merkle tree that goes into the relay.
	ClaimTypeSetRootKey = NewClaimTypeNum(2)
	// ClaimTypeAssignName is a claim type to assign a name to an ID
	ClaimTypeAssignName = NewClaimTypeNum(3)
	// ClaimTypeAuthorizeKSignSecp256k1 is a claim type to autorize a secp256k1 public key for signing.
	ClaimTypeAuthorizeKSignSecp256k1 = NewClaimTypeNum(4)
	// ClaimTypeLinkObjectIdentity is a claim type to link an object (represented by a hash) to an identity.
	ClaimTypeLinkObjectIdentity = NewClaimTypeNum(5)
	// ClaimTypeAuthorizeService is a claim type to authorize a Service for the identity that performs the claim
	ClaimTypeAuthorizeService = NewClaimTypeNum(6)
	// ClaimTypeNonce is a claim used to increment the tree nonce to modify the root hash
	ClaimTypeNonce = NewClaimTypeNum(7)
	// ClaimTypeEthId is a claim type to autorize an Eth Address to be used as Id inside Ethereum
	ClaimTypeEthId = NewClaimTypeNum(8)
	// ClaimTypeAuthEthKey is a claim type to authorize an Eth Address directly from a private key, allowing to specify if is used as KDisable (revoke), KReenable (recover), etc
	ClaimTypeAuthEthKey = NewClaimTypeNum(9)
)

// ClaimVersionLen is the length in bytes of the version in a claim.
const ClaimVersionLen = 32 / 8

// ClaimTypeVersionLen is the length in bytes of the version and length in a claim.
const ClaimTypeVersionLen = ClaimTypeLen + ClaimVersionLen

// NewClaimFromEntry deserializes a valid claim type into a Claim.
func NewClaimFromEntry(e *merkletree.Entry) (merkletree.Claim, error) {
	for _, elemBytes := range e.Data {
		bigints := merkletree.ElemBytesToBigInt(elemBytes)
		ok := cryptoUtils.CheckBigIntInField(bigints, cryptoConstants.Q)
		if !ok {
			return nil, errors.New("Elements not in the Finite Field over R")
		}
	}
	claimType, _ := GetClaimTypeVersion(e)
	switch claimType {
	case *ClaimTypeBasic:
		c := NewClaimBasicFromEntry(e)
		return c, nil
	case *ClaimTypeAssignName:
		c := NewClaimAssignNameFromEntry(e)
		return c, nil
	case *ClaimTypeAuthorizeKSignBabyJub:
		c := NewClaimAuthorizeKSignBabyJubFromEntry(e)
		return c, nil
	case *ClaimTypeSetRootKey:
		c := NewClaimSetRootKeyFromEntry(e)
		return c, nil
	case *ClaimTypeAuthorizeKSignSecp256k1:
		return NewClaimAuthorizeKSignSecp256k1FromEntry(e)
	case *ClaimTypeLinkObjectIdentity:
		c := NewClaimLinkObjectIdentityFromEntry(e)
		return c, nil
	case *ClaimTypeAuthorizeService:
		c := NewClaimAuthorizeServiceFromEntry(e)
		return c, nil
	case *ClaimTypeEthId:
		c := NewClaimEthIdFromEntry(e)
		return c, nil
	case *ClaimTypeAuthEthKey:
		c := NewClaimAuthEthKeyFromEntry(e)
		return c, nil
	default:
		return nil, ErrInvalidClaimType
	}
}
