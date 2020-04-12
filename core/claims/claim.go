package claims

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/crypto"
	"github.com/iden3/go-iden3-core/merkletree"

	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
)

// ErrInvalidClaimType indicates a type error when parsing an Entry into a claim.
var ErrInvalidClaimType = errors.New("invalid claim type")

// ClearMostSigByte sets the most significant byte of the element to 0 to make sure it fits
// inside the FiniteField over R.
func ClearMostSigByte(e [merkletree.ElemBytesLen]byte) merkletree.ElemBytes {
	e[0] = 0
	return merkletree.ElemBytes(e)
}

func GetRevocationNonce(e *merkletree.Entry) uint32 {
	return binary.BigEndian.Uint32(e.Data[4][:4])
}

const (
	// ClaimTypeLen is the length in bytes of the type in a claim.
	ClaimTypeLen       = 64 / 8
	ClaimFlagsLen      = 32 / 8
	ClaimHeaderLen     = ClaimTypeLen + ClaimFlagsLen
	ClaimVersionLen    = 32 / 8
	ClaimRevNonceLen   = 32 / 8
	ClaimExpirationLen = 64 / 8
	EntryFullBytesLen  = 248 / 8
)

// HashString takes the first 31 bytes of a hash applied to string
func HashString(s string) (stringHashed [EntryFullBytesLen]byte) {
	hash := crypto.HashBytes([]byte(s))
	copy(stringHashed[:], hash[len(hash)-EntryFullBytesLen:])
	return stringHashed
}

// ClaimType is the type used to store a claim type.
type ClaimType [ClaimTypeLen]byte

func (ct ClaimType) MarshalText() ([]byte, error) {
	var str string
	switch ct {
	case ClaimTypeBasic:
		str = fmt.Sprintf("str:%v", ClaimTypeStringBasic)
	case ClaimTypeKeyBabyJub:
		str = fmt.Sprintf("str:%v", ClaimTypeStringKeyBabyJub)
	default:
		str = fmt.Sprintf("hex:%v", common.Hex(ct[:]))
	}
	return []byte(str), nil
}

func (ct *ClaimType) UnmarshalText(b []byte) error {
	str := string(b)
	if strings.HasPrefix(str, "str:") {
		str := strings.TrimPrefix(str, "str:")
		switch str {
		case ClaimTypeStringBasic:
			*ct = ClaimTypeBasic
		case ClaimTypeStringKeyBabyJub:
			*ct = ClaimTypeKeyBabyJub
		default:
			return fmt.Errorf("Unknown ClaimType str:%v", str)
		}
	} else if strings.HasPrefix(str, "hex:") {
		str := strings.TrimPrefix(str, "hex:")
		if err := common.HexDecodeInto(ct[:], []byte(str)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid ClaimType prefix")
	}
	return nil
}

// NewClaimType creates a ClaimType from a type name.
func NewClaimType(name string) ClaimType {
	t := ClaimType{}
	h := crypto.HashBytes([]byte(name))
	copy(t[:ClaimTypeLen], h[len(h)-ClaimTypeLen:])
	return t
}

// NewClaimTypeNum to set type to a claim.
func NewClaimTypeNum(num uint64) ClaimType {
	ct := ClaimType{}
	binary.BigEndian.PutUint64(ct[:], num)
	return ct
}

var (
	// ClaimTypeBasic is a simple claim type that can be used for anything.
	ClaimTypeBasic       = NewClaimTypeNum(0)
	ClaimTypeStringBasic = "Basic"

	// ClaimTypeKeyBabyJub is a claim type to autorize a babyjub public key for signing.
	ClaimTypeKeyBabyJub       = NewClaimTypeNum(1)
	ClaimTypeStringKeyBabyJub = "KeyBabyJub"

// 	// ClaimTypeSetRootKey is a claim type of the root key of a merkle tree that goes into the relay.
// 	ClaimTypeSetRootKey = NewClaimTypeNum(2)
// 	// ClaimTypeAssignName is a claim type to assign a name to an ID
// 	ClaimTypeAssignName = NewClaimTypeNum(3)
// 	// ClaimTypeAuthorizeKSignSecp256k1 is a claim type to autorize a secp256k1 public key for signing.
// 	ClaimTypeAuthorizeKSignSecp256k1 = NewClaimTypeNum(4)
// 	// ClaimTypeLinkObjectIdentity is a claim type to link an object (represented by a hash) to an identity.
// 	ClaimTypeLinkObjectIdentity = NewClaimTypeNum(5)
// 	// ClaimTypeAuthorizeService is a claim type to authorize a Service for the identity that performs the claim
// 	ClaimTypeAuthorizeService = NewClaimTypeNum(6)
// 	// ClaimTypeNonce is a claim used to increment the tree nonce to modify the root hash
// 	ClaimTypeNonce = NewClaimTypeNum(7)
// 	// ClaimTypeEthId is a claim type to autorize an Eth Address to be used as Id inside Ethereum
// 	ClaimTypeEthId = NewClaimTypeNum(8)
// 	// ClaimTypeAuthEthKey is a claim type to authorize an Eth Address directly from a private key, allowing to specify if is used as KDisable (revoke), KReenable (recover), etc
// 	ClaimTypeAuthEthKey = NewClaimTypeNum(9)
)

// ClaimTypeVersionLen is the length in bytes of the version and length in a claim.
const ClaimTypeVersionLen = ClaimTypeLen + ClaimFlagsLen + ClaimVersionLen

// NewClaimFromEntry deserializes a valid claim type into a Claim.
func NewClaimFromEntry(e *merkletree.Entry) (merkletree.Entrier, error) {
	for _, elemBytes := range e.Data {
		bigints := merkletree.ElemBytesToBigInt(elemBytes)
		ok := cryptoUtils.CheckBigIntInField(bigints)
		if !ok {
			return nil, errors.New("Elements not in the Finite Field over R")
		}
	}
	var metadata Metadata
	metadata.Unmarshal(e)
	switch metadata.Type() {
	case ClaimTypeBasic:
		c := NewClaimBasicFromEntry(e)
		return c, nil
	// case *ClaimTypeAssignName:
	// 	c := NewClaimAssignNameFromEntry(e)
	// 	return c, nil
	case ClaimTypeKeyBabyJub:
		c := NewClaimKeyBabyJubFromEntry(e)
		return c, nil
	// case *ClaimTypeSetRootKey:
	// 	c := NewClaimSetRootKeyFromEntry(e)
	// 	return c, nil
	// case *ClaimTypeAuthorizeKSignSecp256k1:
	// 	return NewClaimAuthorizeKSignSecp256k1FromEntry(e)
	// case *ClaimTypeLinkObjectIdentity:
	// 	c := NewClaimLinkObjectIdentityFromEntry(e)
	// 	return c, nil
	// case *ClaimTypeAuthorizeService:
	// 	c := NewClaimAuthorizeServiceFromEntry(e)
	// 	return c, nil
	// case *ClaimTypeEthId:
	// 	c := NewClaimEthIdFromEntry(e)
	// 	return c, nil
	// case *ClaimTypeAuthEthKey:
	// 	c := NewClaimAuthEthKeyFromEntry(e)
	// 	return c, nil
	default:
		return nil, ErrInvalidClaimType
	}
}

// ClaimRecip is the flag option to specify a recipient of a claim
type ClaimRecip byte

const (
	// ClaimRecipSelf is a claim that refers to a property of the issuing
	// identity.
	ClaimRecipSelf       ClaimRecip = 0b00
	ClaimRecipStringSelf string     = "Self"
	// ClaimRecipIdenIndex is a claim that refers to a property of an
	// identity found in the index part of the claim.
	ClaimRecipIdenIndex       ClaimRecip = 0b01
	ClaimRecipStringIdenIndex string     = "IdenIndex"
	// ClaimRecipIdenIndex is a claim that refers to a property of an
	// identity found in the value part of the claim.
	ClaimRecipIdenValue       ClaimRecip = 0b10
	ClaimRecipStringIdenValue string     = "IdenValue"
)

func (cr ClaimRecip) MarshalText() ([]byte, error) {
	var str string
	switch cr {
	case ClaimRecipSelf:
		str = ClaimRecipStringSelf
	case ClaimRecipIdenIndex:
		str = ClaimRecipStringIdenIndex
	case ClaimRecipIdenValue:
		str = ClaimRecipStringIdenValue
	default:
		return nil, fmt.Errorf("invalid ClaimRecip")
	}
	return []byte(str), nil
}

func (cr *ClaimRecip) UnmarshalText(b []byte) error {
	switch string(b) {
	case ClaimRecipStringSelf:
		*cr = ClaimRecipSelf
	case ClaimRecipStringIdenIndex:
		*cr = ClaimRecipIdenIndex
	case ClaimRecipStringIdenValue:
		*cr = ClaimRecipIdenValue
	default:
		return fmt.Errorf("invalid ClaimRecip")
	}
	return nil
}

// ClaimHeader represents the first bytes of the claim index and contains its
// type and flags.
type ClaimHeader struct {
	Type       ClaimType
	Dest       ClaimRecip
	Expiration bool
	Version    bool
}

func bool2byte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func byte2bool(b byte) bool {
	return b != 0
}

// Marshal the ClaimHeader into an entry
func (c ClaimHeader) Marshal(e *merkletree.Entry) {
	index := e.Index()
	copy(index[0][:ClaimTypeLen], c.Type[:])
	flags0 := &index[0][ClaimTypeLen]
	*flags0 = 0
	*flags0 |= byte(c.Dest)
	*flags0 |= bool2byte(c.Expiration) << 2
	*flags0 |= bool2byte(c.Version) << 3
}

// Unmarshal the ClaimHeader from an entry
func (c *ClaimHeader) Unmarshal(e *merkletree.Entry) {
	index := e.Index()
	copy(c.Type[:], index[0][:ClaimTypeLen])
	flags0 := index[0][ClaimTypeLen]
	c.Dest = ClaimRecip(flags0 & 0b00000011)
	c.Expiration = byte2bool(flags0 & (1 << 2))
	c.Version = byte2bool(flags0 & (1 << 3))
}

var (
	// ClaimHeaderBasic is a simple claim type that can be used for anything.
	ClaimHeaderBasic = ClaimHeader{
		Type:       ClaimTypeBasic,
		Dest:       ClaimRecipSelf,
		Expiration: false,
		Version:    false}
	// ClaimTypeKeyBabyJub is a claim type issued about a babyjub public key.
	ClaimHeaderKeyBabyJub = ClaimHeader{
		Type:       ClaimTypeKeyBabyJub,
		Dest:       ClaimRecipSelf,
		Expiration: false,
		Version:    false}
)

// Claimer is an intefrace that extends Entrier with a function that
// returns the claim metadata.
type Claimer interface {
	merkletree.Entrier
	Metadata() *Metadata
}

// Metadata is a header and generic (some optional) values of a claim.
type Metadata struct {
	header     ClaimHeader
	Dest       *core.ID
	Expiration int64
	Version    uint32
	RevNonce   uint32
}

// NewMetadata creates a new Metadata with a specific header.
func NewMetadata(header ClaimHeader) Metadata {
	return Metadata{header: header}
}

// Header returns the header from the metadata.
func (m *Metadata) Header() ClaimHeader {
	return m.header
}

// Type returns the claim type from the header in the metadata.
func (m *Metadata) Type() ClaimType {
	return m.header.Type
}

// Marshal the Metadata into an entry
func (m Metadata) Marshal(e *merkletree.Entry) {
	m.header.Marshal(e)
	index := e.Index()
	value := e.Value()
	switch m.header.Dest {
	case ClaimRecipSelf:
	case ClaimRecipIdenIndex:
		copy(index[1][:], m.Dest[:])
	case ClaimRecipIdenValue:
		copy(value[1][:], m.Dest[:])
	default:
		panic(fmt.Sprintf("Unexpected header.Dest %v", m.header.Dest))
	}
	if m.header.Version {
		binary.BigEndian.PutUint32(index[0][ClaimTypeLen+ClaimFlagsLen:], m.Version)
	}
	if m.header.Expiration {
		binary.BigEndian.PutUint64(value[0][ClaimRevNonceLen:], uint64(m.Expiration))
	}
	binary.BigEndian.PutUint32(value[0][:], m.RevNonce)
}

// Unmarshal the Metadata from an entry
func (m *Metadata) Unmarshal(e *merkletree.Entry) {
	m.header.Unmarshal(e)
	index := e.Index()
	value := e.Value()
	if m.header.Dest != ClaimRecipSelf {
		m.Dest = &core.ID{}
	}
	switch m.header.Dest {
	case ClaimRecipSelf:
	case ClaimRecipIdenIndex:
		copy(m.Dest[:], index[1][:])
	case ClaimRecipIdenValue:
		copy(m.Dest[:], value[1][:])
	default:
		panic(fmt.Sprintf("Unexpected header.Dest %v", m.header.Dest))
	}
	if m.header.Version {
		m.Version = binary.BigEndian.Uint32(index[0][ClaimTypeLen+ClaimFlagsLen:])
	}
	if m.header.Expiration {
		m.Expiration = int64(binary.BigEndian.Uint64(value[0][ClaimRevNonceLen:]))
	}
	m.RevNonce = binary.BigEndian.Uint32(value[0][:])
}

type metadataJSON struct {
	Type       ClaimType
	Recip      ClaimRecip
	ID         *core.ID
	Expiration *int64
	Version    *uint32
	RevNonce   uint32
}

func (m Metadata) MarshalJSON() ([]byte, error) {
	var metadata metadataJSON
	h := m.Header()
	metadata.Type = h.Type
	metadata.Recip = h.Dest
	if h.Dest != ClaimRecipSelf {
		metadata.ID = m.Dest
	}
	if h.Expiration {
		metadata.Expiration = &m.Expiration
	}
	if h.Version {
		metadata.Version = &m.Version
	}
	metadata.RevNonce = m.RevNonce
	return json.Marshal(metadata)
}

func (m *Metadata) UnmarshalJSON(b []byte) error {
	var metadata metadataJSON
	if err := json.Unmarshal(b, &metadata); err != nil {
		return err
	}
	m.header = ClaimHeader{
		Type:       metadata.Type,
		Dest:       metadata.Recip,
		Expiration: metadata.Expiration != nil,
		Version:    metadata.Version != nil,
	}
	switch m.header.Type {
	case ClaimTypeBasic:
		if m.header != ClaimHeaderBasic {
			return fmt.Errorf("claim header for ClaimType %v is different than expected",
				ClaimTypeStringBasic)
		}
	case ClaimTypeKeyBabyJub:
		if m.header != ClaimHeaderKeyBabyJub {
			return fmt.Errorf("claim header for ClaimType %v is different than expected",
				ClaimTypeStringKeyBabyJub)
		}
	default:
	}
	if m.header.Dest != ClaimRecipSelf {
		m.Dest = metadata.ID
	}
	if m.header.Expiration {
		m.Expiration = *metadata.Expiration
	}
	if m.header.Version {
		m.Version = *metadata.Version
	}
	m.RevNonce = metadata.RevNonce
	return nil
}
