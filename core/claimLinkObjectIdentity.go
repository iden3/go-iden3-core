package core

import (
	"encoding/binary"

	"github.com/iden3/go-iden3/merkletree"
)

// HashType defines the type of hash used in objectHash.
type HashType uint32

const (
	// HashTypeKeccak256 indicates that hash keccak256 is used.
	HashTypeKeccak256 HashType = 0
	// HashTypeSha256 indicates that hash sha256 is used.
	HashTypeSha256 HashType = 1
)

// ObjectType defines the type of object that objectHash is representing.
type ObjectType uint32

const (
	// ObjectTypePassport indicates that hash represents a passport.
	ObjectTypePassport ObjectType = 0
	// ObjectTypeAddress indicates that hash represents an address.
	ObjectTypeAddress ObjectType = 1
	// ObjectTypePhone indicates that hash represents a phone number.
	ObjectTypePhone ObjectType = 2
	// ObjectTypeDob indicates that hash represents date of birth.
	ObjectTypeDob ObjectType = 3
	// ObjectTypeGivenName indicates that hash represents a given name.
	ObjectTypeGivenName ObjectType = 4
	// ObjectTypeFamilyName indicates that hash represents a family name.
	ObjectTypeFamilyName ObjectType = 5
	// ObjectTypeCertificate indicates that hash represents a certificate.
	ObjectTypeCertificate ObjectType = 6
	// ObjectTypeStorage indicates that hash represents a stored file.
	ObjectTypeStorage ObjectType = 7
)

// ClaimLinkObjectIdentity aims to link a hash of an object to an identity.
type ClaimLinkObjectIdentity struct {
	// Version is the claim version.
	Version uint32
	// ObjectType is the representation of the objectHash.
	ObjectType ObjectType
	// ObjectIndex is the index of this object which the identity has.
	ObjectIndex uint16
	// Id is the ID.
	Id ID
	// ObjectHash is the hash of the object.
	ObjectHash [256 / 8]byte
	// Auxiliary data to complement claim information.
	AuxData [256 / 8]byte
}

// minInt returns the minimum between two inputs
func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// NewClaimLinkObjectIdentity returns a ClaimLinkObjectIdentity.
func NewClaimLinkObjectIdentity(objectType ObjectType, objectIndex uint16, id ID,
	objectHash [256 / 8]byte, auxData [256 / 8]byte) (*ClaimLinkObjectIdentity, error) {
	if _, err := merkletree.ElemBytesToRElem(merkletree.ElemBytes(objectHash)); err != nil {
		return nil, err
	}
	if _, err := merkletree.ElemBytesToRElem(merkletree.ElemBytes(auxData)); err != nil {
		return nil, err
	}
	return &ClaimLinkObjectIdentity{
		Version:     0,
		ObjectType:  objectType,
		ObjectIndex: objectIndex,
		Id:          id,
		ObjectHash:  objectHash,
		AuxData:     auxData,
	}, nil
}

// NewClaimLinkObjectIdentityFromEntry deserializes a ClaimLinkObjectIdentity from an Entry.
func NewClaimLinkObjectIdentityFromEntry(entry *merkletree.Entry) *ClaimLinkObjectIdentity {
	claim := &ClaimLinkObjectIdentity{}
	_, claim.Version = getClaimTypeVersion(entry)
	var objectType [32 / 8]byte
	var objectIndex [16 / 8]byte
	var indexLen = ClaimTypeVersionLen
	// object type
	copyFromElemBytes(objectType[:], indexLen, &entry.Data[3])
	claim.ObjectType = ObjectType(binary.BigEndian.Uint32(objectType[:]))
	// object index
	indexLen += len(objectType)
	copyFromElemBytes(objectIndex[:], indexLen, &entry.Data[3])
	claim.ObjectIndex = binary.BigEndian.Uint16(objectIndex[:])
	// identity address
	copyFromElemBytes(claim.Id[:], 0, &entry.Data[2])
	// hash object
	copyFromElemBytes(claim.ObjectHash[:], 0, &entry.Data[1])
	// hash type
	copyFromElemBytes(claim.AuxData[:], 0, &entry.Data[0])
	return claim
}

// Entry serializes the claim into an Entry.
func (claim *ClaimLinkObjectIdentity) Entry() *merkletree.Entry {
	entry := &merkletree.Entry{}
	var indexLen = ClaimTypeVersionLen
	// type and version
	setClaimTypeVersion(entry, claim.Type(), claim.Version)
	// object type
	var objectType [32 / 8]byte
	binary.BigEndian.PutUint32(objectType[:], uint32(claim.ObjectType))
	copyToElemBytes(&entry.Data[3], indexLen, objectType[:])
	// object index
	indexLen += len(objectType)
	var objectIndex [16 / 8]byte
	binary.BigEndian.PutUint16(objectIndex[:], claim.ObjectIndex)
	copyToElemBytes(&entry.Data[3], indexLen, objectIndex[:])
	// identity address
	copyToElemBytes(&entry.Data[2], 0, claim.Id[:])
	// object hash
	copyToElemBytes(&entry.Data[1], 0, claim.ObjectHash[:])
	// aux data
	copyToElemBytes(&entry.Data[0], 0, claim.AuxData[:])
	return entry
}

// Type returns the ClaimType of the claim.
func (c *ClaimLinkObjectIdentity) Type() ClaimType {
	return *ClaimTypeLinkObjectIdentity
}
