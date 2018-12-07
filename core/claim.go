package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/merkletree"
)

var (
	defaultTypeHash        = merkletree.HashBytes([]byte("default"))
	assignNameTypeHash     = merkletree.HashBytes([]byte("assignname"))
	authorizeksignTypeHash = merkletree.HashBytes([]byte("authorizeksign"))
	setRootTypeHash        = merkletree.HashBytes([]byte("setroot"))

	GenericType        = defaultTypeHash[:24]
	AssignNameType     = assignNameTypeHash[:24]
	AuthorizeksignType = authorizeksignTypeHash[:24]
	SetRootType        = setRootTypeHash[:24]

	NamespaceHash = merkletree.HashBytes([]byte("iden3.io"))
)

// BaseIndex is the by default parameters of the index of every Claim
type BaseIndex struct {
	Namespace   merkletree.Hash // keccak("iden3.io")
	Type        [24]byte        // claim type, keccak("<spec>") [32:56]
	IndexLength uint32          // [4]byte
	Version     uint32          // [4] byte
}

// GenericClaim is a default data structure of a claim
type GenericClaim struct {
	BaseIndex
	ExtraIndex struct {
		Data []byte
	}
	Data []byte
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

// ParseGenericClaimBytes returns a GenericClaim struct from an array of bytes
func ParseGenericClaimBytes(b []byte) (GenericClaim, error) {
	if len(b) < 64 {
		return GenericClaim{}, errors.New("[]byte too small")
	}
	var c GenericClaim
	copy(c.BaseIndex.Namespace[:], b[0:32])
	copy(c.BaseIndex.Type[:], b[32:56])
	c.BaseIndex.IndexLength = EthBytesToUint32(b[56:60])
	c.BaseIndex.Version = EthBytesToUint32(b[60:64])
	c.ExtraIndex.Data = b[64:c.BaseIndex.IndexLength]
	c.Data = b[c.BaseIndex.IndexLength:]
	return c, nil
}

// Bytes returns an array of bytes with the GenericClaim data
func (c GenericClaim) Bytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
	b = append(b, indexLengthBytes[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.ExtraIndex.Data[:]...)
	b = append(b, c.Data[:]...)
	return b
}

// IndexLength returns the length of the Index (BaseIndex + ExtraIndex) of the GenericClaim
func (c GenericClaim) IndexLength() uint32 {
	// return uint32(len(c.Bytes()))
	return c.BaseIndex.IndexLength
}

// Hi returns the hash of the index of the claim
func (c GenericClaim) Hi() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes()[:c.BaseIndex.IndexLength])
	// h := merkletree.HashBytes(c.Bytes())
	return h
}

// Ht returns the hash of the full claim
func (c GenericClaim) Ht() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// ParseAssignNameClaimBytes returns an AssignNameClaim struct from an array of bytes
func ParseAssignNameClaimBytes(b []byte) (AssignNameClaim, error) {
	if len(b) < 148 {
		return AssignNameClaim{}, errors.New("[]byte too small")
	}
	var c AssignNameClaim
	copy(c.BaseIndex.Namespace[:], b[0:32])
	copy(c.BaseIndex.Type[:], b[32:56])
	c.BaseIndex.IndexLength = EthBytesToUint32(b[56:60])
	c.BaseIndex.Version = EthBytesToUint32(b[60:64])
	copy(c.ExtraIndex.Name[:], b[64:96])
	copy(c.ExtraIndex.Domain[:], b[96:128])
	copy(c.EthID[:], b[128:148])
	return c, nil
}

// Bytes returns an array of bytes with the AssignNameClaim data
func (c AssignNameClaim) Bytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
	b = append(b, indexLengthBytes[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.ExtraIndex.Name[:]...)
	b = append(b, c.ExtraIndex.Domain[:]...)
	b = append(b, c.EthID[:]...)
	return b
}

// IndexLength returns the length of the Index (BaseIndex + ExtraIndex) of the AssignNameClaim
func (c AssignNameClaim) IndexLength() uint32 {
	// var bytesIndex []byte
	// bytesIndex = append(bytesIndex, c.BaseIndex.Namespace[:]...)
	// bytesIndex = append(bytesIndex, c.BaseIndex.Type[:]...)
	// indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
	// bytesIndex = append(bytesIndex, indexLengthBytes[:]...)
	// versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	// bytesIndex = append(bytesIndex, versionBytes[:]...)
	// bytesIndex = append(bytesIndex, c.ExtraIndex.Name[:]...)
	// bytesIndex = append(bytesIndex, c.ExtraIndex.Domain[:]...)
	// return uint32(len(bytesIndex))
	return c.BaseIndex.IndexLength
}

// Hi returns the hash of the index of the claim
func (c AssignNameClaim) Hi() merkletree.Hash {
	// var bytesIndex []byte
	// bytesIndex = append(bytesIndex, c.BaseIndex.Namespace[:]...)
	// bytesIndex = append(bytesIndex, c.BaseIndex.Type[:]...)
	// indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
	// bytesIndex = append(bytesIndex, indexLengthBytes[:]...)
	// versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	// bytesIndex = append(bytesIndex, versionBytes[:]...)
	// bytesIndex = append(bytesIndex, c.ExtraIndex.Name[:]...)
	// bytesIndex = append(bytesIndex, c.ExtraIndex.Domain[:]...)
	// h := merkletree.HashBytes(bytesIndex)
	h := merkletree.HashBytes(c.Bytes()[:c.BaseIndex.IndexLength])
	return h
}

// Ht returns the hash of the full claim
func (c AssignNameClaim) Ht() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// ParseAuthorizeKSignClaimBytes returns an KSignClaim struct from an array of bytes
func ParseAuthorizeKSignClaimBytes(b []byte) (AuthorizeKSignClaim, error) {
	if len(b) < 164 {
		return AuthorizeKSignClaim{}, errors.New("[]byte too small")
	}
	var c AuthorizeKSignClaim
	copy(c.BaseIndex.Namespace[:], b[0:32])
	copy(c.BaseIndex.Type[:], b[32:56])
	c.BaseIndex.IndexLength = EthBytesToUint32(b[56:60])
	c.BaseIndex.Version = EthBytesToUint32(b[60:64])
	copy(c.ExtraIndex.KeyToAuthorize[:], b[64:84])
	copy(c.Application[:], b[84:116])
	copy(c.ApplicationAuthz[:], b[116:148])
	c.ValidFrom = EthBytesToUint64(b[148:156])
	c.ValidUntil = EthBytesToUint64(b[156:164])
	return c, nil
}

// Bytes returns an array of bytes with the KSignClaim data
func (c AuthorizeKSignClaim) Bytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
	b = append(b, indexLengthBytes[:]...)
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

// func (c AuthorizeKSignClaim) indexBytes() (b []byte) {
// 	b = append(b, c.BaseIndex.Namespace[:]...)
// 	b = append(b, c.BaseIndex.Type[:]...)
// 	indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
// 	b = append(b, indexLengthBytes[:]...)
// 	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
// 	b = append(b, versionBytes[:]...)
// 	b = append(b, c.ExtraIndex.KeyToAuthorize[:]...)
// 	return b
// }

// IndexLength returns the length of the Index (BaseIndex + ExtraIndex) of the AuthorizeKSignClaim
func (c AuthorizeKSignClaim) IndexLength() uint32 {
	// return uint32(len(c.indexBytes()))
	return c.BaseIndex.IndexLength
}

// Hi returns the hash of the index of the claim
func (c AuthorizeKSignClaim) Hi() merkletree.Hash {
	// return merkletree.HashBytes(c.indexBytes())
	return merkletree.HashBytes(c.Bytes()[:c.BaseIndex.IndexLength])
}

// Ht returns the hash of the full claim
func (c AuthorizeKSignClaim) Ht() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// ParseSetRootClaimBytes returns a SetRootClaim struct from an array of bytes
func ParseSetRootClaimBytes(b []byte) (SetRootClaim, error) {
	if len(b) < 116 {
		return SetRootClaim{}, errors.New("[]byte too small")
	}
	var c SetRootClaim
	copy(c.BaseIndex.Namespace[:], b[0:32])
	copy(c.BaseIndex.Type[:], b[32:56])
	c.BaseIndex.IndexLength = EthBytesToUint32(b[56:60])
	c.BaseIndex.Version = EthBytesToUint32(b[60:64])
	copy(c.ExtraIndex.EthID[:], b[64:84])
	copy(c.Root[:], b[84:116])
	return c, nil
}

// Bytes returns an array of bytes with the SetRootClaim data
func (c SetRootClaim) Bytes() (b []byte) {
	b = append(b, c.BaseIndex.Namespace[:]...)
	b = append(b, c.BaseIndex.Type[:]...)
	indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
	b = append(b, indexLengthBytes[:]...)
	versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.ExtraIndex.EthID[:]...)
	b = append(b, c.Root[:]...)
	return b
}

// IndexLength returns the length of the Index (BaseIndex + ExtraIndex) of the SetRootClaim
func (c SetRootClaim) IndexLength() uint32 {
	// var bytesIndex []byte
	// bytesIndex = append(bytesIndex, c.BaseIndex.Namespace[:]...)
	// bytesIndex = append(bytesIndex, c.BaseIndex.Type[:]...)
	// indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
	// bytesIndex = append(bytesIndex, indexLengthBytes[:]...)
	// versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	// bytesIndex = append(bytesIndex, versionBytes[:]...)
	// bytesIndex = append(bytesIndex, c.ExtraIndex.EthID[:]...)
	// return uint32(len(bytesIndex))
	return c.BaseIndex.IndexLength
}

// Hi returns the hash of the index of the claim
func (c SetRootClaim) Hi() merkletree.Hash {
	// var bytesIndex []byte
	// bytesIndex = append(bytesIndex, c.BaseIndex.Namespace[:]...)
	// bytesIndex = append(bytesIndex, c.BaseIndex.Type[:]...)
	// indexLengthBytes, _ := Uint32ToEthBytes(c.BaseIndex.IndexLength)
	// bytesIndex = append(bytesIndex, indexLengthBytes[:]...)
	// versionBytes, _ := Uint32ToEthBytes(c.BaseIndex.Version)
	// bytesIndex = append(bytesIndex, versionBytes[:]...)
	// bytesIndex = append(bytesIndex, c.ExtraIndex.EthID[:]...)
	// h := merkletree.HashBytes(bytesIndex)
	h := merkletree.HashBytes(c.Bytes()[:c.BaseIndex.IndexLength])
	return h
}

// Ht returns the hash of the full claim
func (c SetRootClaim) Ht() merkletree.Hash {
	h := merkletree.HashBytes(c.Bytes())
	return h
}

// ParseTypeClaimBytes returns the type of the claim from an array of bytes
func ParseTypeClaimBytes(b []byte) (string, error) {
	if len(b) < 64 { // 64, as is the minimum length of the BaseIndex
		return "", errors.New("[]byte too small")
	}
	if int(EthBytesToUint32(b[56:60])) > len(b) {
		return "", errors.New("claim.BaseIndex.IndexLength can not be bigger than claim bytes length")
	}
	typeBytes := b[32:56]
	if bytes.Equal(GenericType, typeBytes) {
		return "default", nil
	} else if bytes.Equal(AssignNameType, typeBytes) {
		return "assignname", nil
	} else if bytes.Equal(AuthorizeksignType, typeBytes) {
		return "authorizeksign", nil
	} else if bytes.Equal(SetRootType, typeBytes) {
		return "setroot", nil
	}
	return "", errors.New("type unrecognized")
}

// ParseValueFromBytes returns a merkletree.Value from a given byte array
func ParseValueFromBytes(b []byte) (merkletree.Value, error) {
	if len(b) < 64 { // 64, as is the minimum length of the BaseIndex
		return GenericClaim{}, errors.New("[]byte too small")
	}
	typeBytes := common3.BytesToHex(b[32:56])
	var value merkletree.Value
	var err error
	switch typeBytes {
	case common3.BytesToHex(GenericType):
		value, err = ParseGenericClaimBytes(b)
		break
	case common3.BytesToHex(AssignNameType):
		value, err = ParseAssignNameClaimBytes(b)
		break
	case common3.BytesToHex(AuthorizeksignType):
		value, err = ParseAuthorizeKSignClaimBytes(b)
		break
	case common3.BytesToHex(SetRootType):
		value, err = ParseSetRootClaimBytes(b)
		break
	default:
		value = GenericClaim{}
		err = errors.New("claim type unrecognized")
		break
	}
	return value, err
}

// NewGenericClaim returns a GenericClaim object with the given parameters
func NewGenericClaim(namespaceStr, typeStr string, extraIndexData []byte, data []byte) GenericClaim {
	var c GenericClaim
	c.BaseIndex.Namespace = merkletree.HashBytes([]byte(namespaceStr))
	typeHash := merkletree.HashBytes([]byte(typeStr))
	copy(c.BaseIndex.Type[:], typeHash[:24])
	c.BaseIndex.IndexLength = 64 + uint32(len(extraIndexData))
	c.BaseIndex.Version = 0
	c.ExtraIndex.Data = extraIndexData
	c.Data = data
	return c
}

// NewAssignNameClaim returns a AssignNameClaim object with the given parameters
func NewAssignNameClaim(name, domain merkletree.Hash, ethID common.Address) AssignNameClaim {
	var c AssignNameClaim
	c.BaseIndex.Namespace = NamespaceHash
	copy(c.BaseIndex.Type[:], AssignNameType)
	c.BaseIndex.IndexLength = 128
	c.BaseIndex.Version = 0
	c.ExtraIndex.Name = name
	c.ExtraIndex.Domain = domain
	c.EthID = ethID
	return c
}

// NewKSignClaim returns a KSignClaim object with the given parameters
func NewAuthorizeKSignClaim(keyToAuthorize common.Address, applicationName, applicationAuthz string, validFrom, validUntil uint64) AuthorizeKSignClaim {
	var c AuthorizeKSignClaim
	c.BaseIndex.Namespace = NamespaceHash
	copy(c.BaseIndex.Type[:], AuthorizeksignType)
	c.BaseIndex.IndexLength = 84
	c.BaseIndex.Version = 0
	c.ExtraIndex.KeyToAuthorize = keyToAuthorize
	c.Application = merkletree.HashBytes([]byte(applicationName))
	c.ApplicationAuthz = merkletree.HashBytes([]byte(applicationAuthz))
	c.ValidFrom = validFrom
	c.ValidUntil = validUntil
	return c
}

func NewOperationalKSignClaim(keyToAuthorize common.Address) AuthorizeKSignClaim {
	var c AuthorizeKSignClaim
	c.BaseIndex.Namespace = NamespaceHash
	copy(c.BaseIndex.Type[:], AuthorizeksignType)
	c.BaseIndex.IndexLength = 84
	c.BaseIndex.Version = 0
	c.ExtraIndex.KeyToAuthorize = keyToAuthorize
	c.Application = merkletree.Hash{}
	c.ApplicationAuthz = merkletree.Hash{}
	c.ValidFrom = 0
	c.ValidUntil = math.MaxUint64
	return c
}

// NewSetRootClaim returns a SetRootClaim object with the given parameters
func NewSetRootClaim(ethID common.Address, root merkletree.Hash) SetRootClaim {
	var c SetRootClaim
	c.BaseIndex.Namespace = NamespaceHash
	copy(c.BaseIndex.Type[:], SetRootType)
	c.BaseIndex.IndexLength = 84
	c.BaseIndex.Version = 0
	c.ExtraIndex.EthID = ethID
	c.Root = root
	return c
}

// HiFromClaimBytes returns the Hi Hash of the given Claim bytes
func HiFromClaimBytes(b []byte) merkletree.Hash {
	indexLength := EthBytesToUint32(b[56:60])
	hi := merkletree.HashBytes(b[:indexLength])
	return hi
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
