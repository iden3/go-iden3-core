package core

import (
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3/crypto/babyjub"
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

// SetClaimTypeVersionInData is a helper function to set the type and version of a
// claim.
func SetClaimTypeVersionInData(d *merkletree.Data, claimType ClaimType, version uint32) {
	copyToElemBytes(&d[3], 0, claimType[:])
	binary.BigEndian.PutUint32(d[3][merkletree.ElemBytesLen-ClaimTypeVersionLen:], version)
}

// getClaimTypeVersion is a helper function to get the type and version from a
// claim.
func getClaimTypeVersion(e *merkletree.Entry) (c ClaimType, v uint32) {
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

// ClaimAuthorizeKSignBabyJub is a claim to authorize a baby jub public key for
// signing.
type ClaimAuthorizeKSignBabyJub struct {
	// Version is the claim version.
	Version uint32
	// Sign means positive if false, negative if true.
	Sign bool
	// Ay is the y coordinate of the baby jub curve point which corresponds
	// to the public key.
	Ay *big.Int
}

// NewClaimAuthorizeKSignBabyJub returns a ClaimAuthorizeKSignBabyJub with the
// given elliptic public key parameters.
func NewClaimAuthorizeKSignBabyJub(pk *babyjub.PublicKey) *ClaimAuthorizeKSignBabyJub {
	return &ClaimAuthorizeKSignBabyJub{
		Version: 0,
		Sign:    babyjub.PointCoordSign(pk.X),
		Ay:      pk.Y,
	}
}

// NewClaimAuthorizeKSignBabyJubFromEntry deserializes a
// ClaimAuthorizeKSignBabyJubFrom from an Entry.
func NewClaimAuthorizeKSignBabyJubFromEntry(e *merkletree.Entry) *ClaimAuthorizeKSignBabyJub {
	c := &ClaimAuthorizeKSignBabyJub{}
	_, c.Version = getClaimTypeVersion(e)
	sign := []byte{0}
	copyFromElemBytes(sign, ClaimTypeVersionLen, &e.Data[3])
	if sign[0] == 1 {
		c.Sign = true
	}
	c.Ay = new(big.Int).SetBytes(e.Data[2][:])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimAuthorizeKSignBabyJub) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	sign := []byte{0}
	if c.Sign {
		sign = []byte{1}
	}
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, sign)
	copy(e.Data[2][:], c.Ay.Bytes())
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAuthorizeKSignBabyJub) Type() ClaimType {
	return *ClaimTypeAuthorizeKSignBabyJub
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
	// Id is the ID related to the root key.
	Id ID
	// RootKey is the root of the mekrlee tree.
	RootKey merkletree.Hash
}

// NewClaimSetRootKey returns a ClaimSetRootKey with the given Eth ID and
// merklee tree root key.
func NewClaimSetRootKey(id ID, rootKey merkletree.Hash) *ClaimSetRootKey {
	return &ClaimSetRootKey{
		Version: 0,
		Era:     0,
		Id:      id,
		RootKey: rootKey,
	}
}

// NewClaimSetRootKeyFromEntry deserializes a ClaimSetRootKey from an Entry.
func NewClaimSetRootKeyFromEntry(e *merkletree.Entry) *ClaimSetRootKey {
	c := &ClaimSetRootKey{}
	_, c.Version = getClaimTypeVersion(e)
	var era [32 / 8]byte
	copyFromElemBytes(era[:], ClaimTypeVersionLen, &e.Data[3])
	c.Era = binary.BigEndian.Uint32(era[:])
	copyFromElemBytes(c.Id[:], 0, &e.Data[2])
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
	copyToElemBytes(&e.Data[2], 0, c.Id[:])
	e.Data[1] = merkletree.ElemBytes(c.RootKey)
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimSetRootKey) Type() ClaimType {
	return *ClaimTypeSetRootKey
}

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
	ObjectHash [248 / 8]byte
	// Auxiliary data to complement claim information.
	AuxData [248 / 8]byte
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
	objectHash []byte, auxData []byte) *ClaimLinkObjectIdentity {
	var objectHashSlice [31]byte
	minLen := minInt(len(objectHash), 32)
	copy(objectHashSlice[:], objectHash[1:minLen])
	var auxDataSlice [31]byte
	minLen = minInt(len(auxData), 32)
	copy(auxDataSlice[:], auxData[1:minLen])
	return &ClaimLinkObjectIdentity{
		Version:     0,
		ObjectType:  objectType,
		ObjectIndex: objectIndex,
		Id:          id,
		ObjectHash:  objectHashSlice,
		AuxData:     auxDataSlice,
	}
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

// ServiceType
var (
	// ServiceTypeRelay is the type for authorize Relays
	ServiceTypeRelay = NewServiceType(0)
	// ServiceTypeNotificationsServer is the type for authorize Notification Server
	ServiceTypeNotificationsServer = NewServiceType(1)
	// ServiceTypeDiscoveryNode is the type for authorize DiscoveryNode
	ServiceTypeDiscoveryNode = NewServiceType(2)
)

// ServiceTypeLen is the length in bytes of the type of the Services
const ServiceTypeLen = 64 / 8

// ServiceType is the type used to store a claim type.
type ServiceType [ServiceTypeLen]byte

// NewServiceType to set type of authorized services
func NewServiceType(num uint64) *ServiceType {
	st := ServiceType{}
	binary.BigEndian.PutUint64(st[:], num)
	return &st
}

// ClaimAuthorizeService is a claim to authorize a Service for the identity that performs the claim
type ClaimAuthorizeService struct {
	// Version is the claim version.
	Version uint32
	// ServiceType is the type of the authorized service
	ServiceType *ServiceType
	// ServiceAddr is the hash of the addr
	ServiceAddr [248 / 8]byte
	// ServicePubK is the hash of the pubK
	ServicePubK [248 / 8]byte
	// ServiceUrl is the hash of the domain
	ServiceUrl [248 / 8]byte
}

// NewClaimAuthorizeService returns a ClaimAuthorizeService with the provided data.
func NewClaimAuthorizeService(serviceType *ServiceType, serviceAddr, servicePubK, serviceUrl string) *ClaimAuthorizeService {
	return &ClaimAuthorizeService{
		Version:     0,
		ServiceType: serviceType,
		ServiceAddr: HashString(serviceAddr),
		ServicePubK: HashString(servicePubK),
		ServiceUrl:  HashString(serviceUrl),
	}
}

// NewClaimAuthorizeServiceFromEntry deserializes a ClaimAuthorizeService from an Entry.
func NewClaimAuthorizeServiceFromEntry(e *merkletree.Entry) *ClaimAuthorizeService {
	c := &ClaimAuthorizeService{}
	_, c.Version = getClaimTypeVersion(e)
	var serviceType [64 / 8]byte
	copyFromElemBytes(serviceType[:], ClaimTypeVersionLen, &e.Data[3])
	c.ServiceType = NewServiceType(binary.BigEndian.Uint64(serviceType[:]))
	copyFromElemBytes(c.ServiceAddr[:], 0, &e.Data[2])
	copyFromElemBytes(c.ServicePubK[:], 0, &e.Data[1])
	copyFromElemBytes(c.ServiceUrl[:], 0, &e.Data[0])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimAuthorizeService) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, c.ServiceType[:])
	copyToElemBytes(&e.Data[2], 0, c.ServiceAddr[:])
	copyToElemBytes(&e.Data[1], 0, c.ServicePubK[:])
	copyToElemBytes(&e.Data[0], 0, c.ServiceUrl[:])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAuthorizeService) Type() ClaimType {
	return *ClaimTypeAuthorizeService
}

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
	_, c.Version = getClaimTypeVersion(e)
	copyFromElemBytes(c.Address[:], 0, &e.Data[2])
	copyFromElemBytes(c.IdentityFactory[:], 0, &e.Data[1])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimEthId) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[2], 0, c.Address[:])
	copyToElemBytes(&e.Data[1], 0, c.IdentityFactory[:])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimEthId) Type() ClaimType {
	return *ClaimTypeEthId
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
	default:
		return nil, ErrInvalidClaimType
	}
}
