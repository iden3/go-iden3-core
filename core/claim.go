package core

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/merkletree"
)

var (
	defaultTypeHash        = merkletree.HashBytes([]byte("default"))
	assignnameTypeHash     = merkletree.HashBytes([]byte("assignname"))
	authorizeksignTypeHash = merkletree.HashBytes([]byte("authorizeksign"))
	setRootTypeHash        = merkletree.HashBytes([]byte("setroot"))
)

// BaseIndex is the by default parameters of the index of every Claim
type BaseIndex struct {
	Namespace merkletree.Hash // keccak("iden3.io")
	Type      merkletree.Hash // claim type, keccak("<spec>")
	Version   uint32          // [4] byte
}

// ClaimDefault is a default data structure of a claim
type ClaimDefault struct {
	BaseIndex
	ExtraIndex struct {
		Data []byte
	}
}

// AssignNameClaim is the claim to assign a name to an identity
type AssignNameClaim struct {
	BaseIndex
	ExtraIndex struct {
		Name   merkletree.Hash // keccak("bob")
		Domain merkletree.Hash // ens_namehash("barcelona.eth")
	}
	EthID common.Address // EthID address of identity
}

// AuthorizeKSignClaim is the claim to authorize a KSign key
type AuthorizeKSignClaim struct {
	BaseIndex
	ExtraIndex struct {
		KeyToAuthorize common.Address
	}
	Application      merkletree.Hash
	ApplicationAuthz merkletree.Hash
	ValidFrom        uint64
	ValidUntil       uint64
}

// SetRootClaim is the Claim that goes inside the Relay's merkletree, that sets the ID merkle root
type SetRootClaim struct {
	BaseIndex
	ExtraIndex struct {
		EthID common.Address
	}
	Root merkletree.Hash
}

// ParseClaimDefaultBytes returns a ClaimDefault struct from an array of bytes
func ParseClaimDefaultBytes(b []byte) (ClaimDefault, error) {
	if len(b) < 68 {
		return ClaimDefault{}, errors.New("[]byte too small")
	}
	var c ClaimDefault
	copy(c.BaseIndex.Namespace[:], b[0:32])
	copy(c.BaseIndex.Type[:], b[32:64])
	versionBytes := b[64:68]
	c.BaseIndex.Version = EthBytesToUint32(versionBytes)
	c.ExtraIndex.Data = b[68:]
	return c, nil
}

// Bytes returns an array of bytes with the ClaimDefault data
func (c ClaimDefault) Bytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.ExtraIndex.Data[:]...)
	return b
}

// IndexLength returns the length of the Index (BaseIndex + ExtraIndex) of the ClaimDefault
func (c ClaimDefault) IndexLength() uint32 {
	return uint32(len(c.Bytes()))
}

// Hi returns the hash of the index of the claim
func (c ClaimDefault) Hi() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// Ht returns the hash of the full claim
func (c ClaimDefault) Ht() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// ParseAssignNameClaimBytes returns an AssignNameClaim struct from an array of bytes
func ParseAssignNameClaimBytes(b []byte) (AssignNameClaim, error) {
	if len(b) < 152 {
		return AssignNameClaim{}, errors.New("[]byte too small")
	}
	var c AssignNameClaim
	copy(c.BaseIndex.Namespace[:], b[0:32])
	copy(c.BaseIndex.Type[:], b[32:64])
	versionBytes := b[64:68]
	c.BaseIndex.Version = EthBytesToUint32(versionBytes)
	copy(c.ExtraIndex.Name[:], b[68:100])
	copy(c.ExtraIndex.Domain[:], b[100:132])
	copy(c.EthID[:], b[132:152])
	return c, nil
}

// Bytes returns an array of bytes with the AssignNameClaim data
func (c AssignNameClaim) Bytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.ExtraIndex.Name[:]...)
	b = append(b, c.ExtraIndex.Domain[:]...)
	b = append(b, c.EthID[:]...)
	return b
}

// IndexLength returns the length of the Index (BaseIndex + ExtraIndex) of the AssignNameClaim
func (c AssignNameClaim) IndexLength() uint32 {
	var bytesIndex []byte
	bytesIndex = append(bytesIndex, c.BaseIndex.Namespace[:]...)
	bytesIndex = append(bytesIndex, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	bytesIndex = append(bytesIndex, versionBytes[:]...)
	bytesIndex = append(bytesIndex, c.ExtraIndex.Name[:]...)
	bytesIndex = append(bytesIndex, c.ExtraIndex.Domain[:]...)
	return uint32(len(bytesIndex))
}

// Hi returns the hash of the index of the claim
func (c AssignNameClaim) Hi() merkletree.Hash {
	var bytesIndex []byte
	bytesIndex = append(bytesIndex, c.BaseIndex.Namespace[:]...)
	bytesIndex = append(bytesIndex, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	bytesIndex = append(bytesIndex, versionBytes[:]...)
	bytesIndex = append(bytesIndex, c.ExtraIndex.Name[:]...)
	bytesIndex = append(bytesIndex, c.ExtraIndex.Domain[:]...)
	h := merkletree.HashBytes(bytesIndex)
	return h
}

// Ht returns the hash of the full claim
func (c AssignNameClaim) Ht() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// ParseAuthorizeKSignClaimBytes returns an KSignClaim struct from an array of bytes
func ParseAuthorizeKSignClaimBytes(b []byte) (AuthorizeKSignClaim, error) {
	if len(b) < 168 {
		return AuthorizeKSignClaim{}, errors.New("[]byte too small")
	}
	var c AuthorizeKSignClaim
	copy(c.BaseIndex.Namespace[:], b[0:32])
	copy(c.BaseIndex.Type[:], b[32:64])
	versionBytes := b[64:68]
	c.BaseIndex.Version = EthBytesToUint32(versionBytes)
	copy(c.ExtraIndex.KeyToAuthorize[:], b[68:88])
	copy(c.Application[:], b[88:120])
	copy(c.ApplicationAuthz[:], b[120:152])
	c.ValidFrom = EthBytesToUint64(b[152:160])
	c.ValidUntil = EthBytesToUint64(b[160:168])
	return c, nil
}

// Bytes returns an array of bytes with the KSignClaim data
func (c AuthorizeKSignClaim) Bytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.ExtraIndex.KeyToAuthorize[:]...)
	b = append(b, c.Application[:]...)
	b = append(b, c.ApplicationAuthz[:]...)
	validFromBytes, _ := Uint64ToEthBytes(c.ValidFrom)
	validUntilBytes, _ := Uint64ToEthBytes(c.ValidUntil)
	b = append(b, validFromBytes...)
	b = append(b, validUntilBytes...)
	return b
}

func (c AuthorizeKSignClaim) indexBytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.ExtraIndex.KeyToAuthorize[:]...)
	return b
}

// IndexLength returns the length of the Index (BaseIndex + ExtraIndex) of the AuthorizeKSignClaim
func (c AuthorizeKSignClaim) IndexLength() uint32 {
	return uint32(len(c.indexBytes()))
}

// Hi returns the hash of the index of the claim
func (c AuthorizeKSignClaim) Hi() merkletree.Hash {
	return merkletree.HashBytes(c.indexBytes())
}

// Ht returns the hash of the full claim
func (c AuthorizeKSignClaim) Ht() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// ParseSetRootClaimBytes returns a SetRootClaim struct from an array of bytes
func ParseSetRootClaimBytes(b []byte) (SetRootClaim, error) {
	if len(b) < 120 {
		return SetRootClaim{}, errors.New("[]byte too small")
	}
	var c SetRootClaim
	copy(c.BaseIndex.Namespace[:], b[0:32])
	copy(c.BaseIndex.Type[:], b[32:64])
	c.BaseIndex.Version = EthBytesToUint32(b[64:68])
	copy(c.ExtraIndex.EthID[:], b[68:88])
	copy(c.Root[:], b[88:120])
	return c, nil
}

// Bytes returns an array of bytes with the SetRootClaim data
func (c SetRootClaim) Bytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.ExtraIndex.EthID[:]...)
	b = append(b, c.Root[:]...)
	return b
}

// IndexLength returns the length of the Index (BaseIndex + ExtraIndex) of the SetRootClaim
func (c SetRootClaim) IndexLength() uint32 {
	var bytesIndex []byte
	bytesIndex = append(bytesIndex, c.BaseIndex.Namespace[:]...)
	bytesIndex = append(bytesIndex, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	bytesIndex = append(bytesIndex, versionBytes[:]...)
	bytesIndex = append(bytesIndex, c.ExtraIndex.EthID[:]...)
	return uint32(len(bytesIndex))
}

// Hi returns the hash of the index of the claim
func (c SetRootClaim) Hi() merkletree.Hash {
	var bytesIndex []byte
	bytesIndex = append(bytesIndex, c.BaseIndex.Namespace[:]...)
	bytesIndex = append(bytesIndex, c.BaseIndex.Type[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	bytesIndex = append(bytesIndex, versionBytes[:]...)
	bytesIndex = append(bytesIndex, c.ExtraIndex.EthID[:]...)
	h := merkletree.HashBytes(bytesIndex)
	return h
}

// Ht returns the hash of the full claim
func (c SetRootClaim) Ht() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// ParseTypeClaimBytes returns the type of the claim from an array of bytes
func ParseTypeClaimBytes(b []byte) (string, error) {
	if len(b) < 68 { // 68, as is the minimum length of the BaseIndex
		return "", errors.New("[]byte too small")
	}
	typeBytes := b[32:64]
	if bytes.Equal(defaultTypeHash[:], typeBytes) {
		return "default", nil
	} else if bytes.Equal(assignnameTypeHash[:], typeBytes) {
		return "assignname", nil
	} else if bytes.Equal(authorizeksignTypeHash[:], typeBytes) {
		return "authorizeksign", nil
	} else if bytes.Equal(setRootTypeHash[:], typeBytes) {
		return "setroot", nil
	}
	return "", errors.New("type unrecognized")
}

// ParseValueFromBytes returns a merkletree.Value from a given byte array
func ParseValueFromBytes(b []byte) (merkletree.Value, error) {
	if len(b) < 68 { // 68, as is the minimum length of the BaseIndex
		return ClaimDefault{}, errors.New("[]byte too small")
	}
	typeBytes := b[32:64]
	var value merkletree.Value
	var err error
	switch common3.BytesToHex(typeBytes) {
	case defaultTypeHash.Hex():
		value, err = ParseClaimDefaultBytes(b)
		break
	case assignnameTypeHash.Hex():
		value, err = ParseAssignNameClaimBytes(b)
		break
	case authorizeksignTypeHash.Hex():
		value, err = ParseAuthorizeKSignClaimBytes(b)
		break
	case setRootTypeHash.Hex():
		value, err = ParseSetRootClaimBytes(b)
		break
	default:
		value = ClaimDefault{}
		err = errors.New("claim type unrecognized")
		break
	}
	return value, err
}

// NewClaimDefault returns a ClaimDefault object with the given parameters
func NewClaimDefault(namespaceStr, typeStr string, data []byte) ClaimDefault {
	var c ClaimDefault
	c.BaseIndex.Namespace = merkletree.HashBytes([]byte(namespaceStr))
	c.BaseIndex.Type = merkletree.HashBytes([]byte(typeStr))
	c.BaseIndex.Version = 0
	c.ExtraIndex.Data = data
	return c
}

// NewAssignNameClaim returns a AssignNameClaim object with the given parameters
func NewAssignNameClaim(namespaceStr string, name, domain merkletree.Hash, ethID common.Address) AssignNameClaim {
	var c AssignNameClaim
	c.BaseIndex.Namespace = merkletree.HashBytes([]byte(namespaceStr))
	c.BaseIndex.Type = assignnameTypeHash
	c.BaseIndex.Version = 0
	c.ExtraIndex.Name = name
	c.ExtraIndex.Domain = domain
	c.EthID = ethID
	return c
}

// NewKSignClaim returns a KSignClaim object with the given parameters
func NewAuthorizeKSignClaim(namespaceStr string, keyToAuthorize common.Address, applicationName, applicationAuthz string, validFrom, validUntil uint64) AuthorizeKSignClaim {
	var c AuthorizeKSignClaim
	c.BaseIndex.Namespace = merkletree.HashBytes([]byte(namespaceStr))
	c.BaseIndex.Type = authorizeksignTypeHash
	c.BaseIndex.Version = 0
	c.ExtraIndex.KeyToAuthorize = keyToAuthorize
	c.Application = merkletree.HashBytes([]byte(applicationName))
	c.ApplicationAuthz = merkletree.HashBytes([]byte(applicationAuthz))
	c.ValidFrom = validFrom
	c.ValidUntil = validUntil
	return c
}

// NewSetRootClaim returns a SetRootClaim object with the given parameters
func NewSetRootClaim(namespaceStr string, ethID common.Address, root merkletree.Hash) SetRootClaim {
	var c SetRootClaim
	c.BaseIndex.Namespace = merkletree.HashBytes([]byte(namespaceStr))
	c.BaseIndex.Type = setRootTypeHash
	c.BaseIndex.Version = 0
	c.ExtraIndex.EthID = ethID
	c.Root = root
	return c
}

func Uint32ToEthBytes(u uint32) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, u)
	return buff.Bytes(), err
}
func Uint64ToEthBytes(u uint64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, u)
	return buff.Bytes(), err
}

func EthBytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}
func EthBytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
