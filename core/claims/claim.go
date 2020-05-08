package claims

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/crypto"
	"github.com/iden3/go-iden3-core/merkletree"
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
	return binary.LittleEndian.Uint32(e.Data[4][:4])
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
	// TODO: Update to LittleEndian (needs update in circuits/buildClaimKeyBBJJ.circom
	binary.BigEndian.PutUint64(ct[:], num)
	return ct
}

// ClaimTypeVersionLen is the length in bytes of the version and length in a claim.
const ClaimTypeVersionLen = ClaimTypeLen + ClaimFlagsLen + ClaimVersionLen

// ClaimSubject is the flag option to specify a recipient of a claim
type ClaimSubject byte

const (
	// ClaimSubjectSelf is a claim that refers to a property of the issuing
	// identity.
	ClaimSubjectSelf       ClaimSubject = 0b00
	ClaimSubjectStringSelf string       = "Self"
	// ClaimSubjectOtherIden is a claim that refers to a property of
	// another identity.
	ClaimSubjectOtherIden       ClaimSubject = 0b10
	ClaimSubjectStringOtherIden string       = "OtherIden"
	// ClaimSubjectObject is a claim that refers to a property of an
	// object.
	ClaimSubjectObject       ClaimSubject = 0b01
	ClaimSubjectStringObject string       = "Object"
)

func (cs ClaimSubject) MarshalText() ([]byte, error) {
	var str string
	switch cs {
	case ClaimSubjectSelf:
		str = ClaimSubjectStringSelf
	case ClaimSubjectOtherIden:
		str = ClaimSubjectStringOtherIden
	case ClaimSubjectObject:
		str = ClaimSubjectStringObject
	default:
		return nil, fmt.Errorf("invalid ClaimSubject")
	}
	return []byte(str), nil
}

func (cs *ClaimSubject) UnmarshalText(b []byte) error {
	switch string(b) {
	case ClaimSubjectStringSelf:
		*cs = ClaimSubjectSelf
	case ClaimSubjectStringOtherIden:
		*cs = ClaimSubjectOtherIden
	case ClaimSubjectStringObject:
		*cs = ClaimSubjectObject
	default:
		return fmt.Errorf("invalid ClaimSubject")
	}
	return nil
}

// ClaimSubjectPos is the flag option to specify the position of the subject in the claim
type ClaimSubjectPos byte

const (
	// ClaimSubjectPosIndex means that the subject is found in the Index of the claim
	ClaimSubjectPosIndex       ClaimSubjectPos = 0b0
	ClaimSubjectPosStringIndex string          = "Index"
	// ClaimSubjectPosIndex means that the subject is found in the Value of the claim
	ClaimSubjectPosValue       ClaimSubjectPos = 0b1
	ClaimSubjectPosStringValue string          = "Value"
)

func (csp ClaimSubjectPos) MarshalText() ([]byte, error) {
	var str string
	switch csp {
	case ClaimSubjectPosIndex:
		str = ClaimSubjectPosStringIndex
	case ClaimSubjectPosValue:
		str = ClaimSubjectPosStringValue
	default:
		return nil, fmt.Errorf("invalid ClaimSubjectPos")
	}
	return []byte(str), nil
}

func (csp *ClaimSubjectPos) UnmarshalText(b []byte) error {
	switch string(b) {
	case ClaimSubjectPosStringIndex:
		*csp = ClaimSubjectPosIndex
	case ClaimSubjectPosStringValue:
		*csp = ClaimSubjectPosValue
	default:
		return fmt.Errorf("invalid ClaimSubjectPos")
	}
	return nil
}

// ClaimHeader represents the first bytes of the claim index and contains its
// type and flags.
type ClaimHeader struct {
	Type       ClaimType
	Subject    ClaimSubject
	SubjectPos ClaimSubjectPos
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
	*flags0 |= byte(c.Subject)
	*flags0 |= byte(c.SubjectPos) << 2
	*flags0 |= bool2byte(c.Expiration) << 3
	*flags0 |= bool2byte(c.Version) << 4
}

// Unmarshal the ClaimHeader from an entry
func (c *ClaimHeader) Unmarshal(e *merkletree.Entry) {
	index := e.Index()
	copy(c.Type[:], index[0][:ClaimTypeLen])
	flags0 := index[0][ClaimTypeLen]
	c.Subject = ClaimSubject(flags0 & 0b00000011)
	c.SubjectPos = ClaimSubjectPos(flags0 & (1 << 2))
	c.Expiration = byte2bool(flags0 & (1 << 3))
	c.Version = byte2bool(flags0 & (1 << 4))
}

// Claimer is an intefrace that extends Entrier with a function that
// returns the claim metadata.
type Claimer interface {
	merkletree.Entrier
	Metadata() *Metadata
}

// Metadata is a header and generic (some optional) values of a claim.
type Metadata struct {
	header     ClaimHeader
	Subject    *core.ID
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
	if m.header.Subject != ClaimSubjectSelf {
		switch m.header.SubjectPos {
		case ClaimSubjectPosIndex:
			copy(index[1][:], m.Subject[:])
		case ClaimSubjectPosValue:
			copy(value[1][:], m.Subject[:])
		default:
			panic(fmt.Sprintf("Unexpected header.SubjectPos %v", m.header.SubjectPos))
		}
	}
	if m.header.Version {
		binary.LittleEndian.PutUint32(index[0][ClaimTypeLen+ClaimFlagsLen:], m.Version)
	}
	if m.header.Expiration {
		binary.LittleEndian.PutUint64(value[0][ClaimRevNonceLen:], uint64(m.Expiration))
	}
	binary.LittleEndian.PutUint32(value[0][:], m.RevNonce)
}

// Unmarshal the Metadata from an entry
func (m *Metadata) Unmarshal(e *merkletree.Entry) {
	m.header.Unmarshal(e)
	index := e.Index()
	value := e.Value()
	if m.header.Subject != ClaimSubjectSelf {
		m.Subject = &core.ID{}
	}
	if m.header.Subject != ClaimSubjectSelf {
		switch m.header.SubjectPos {
		case ClaimSubjectPosIndex:
			copy(m.Subject[:], index[1][:])
		case ClaimSubjectPosValue:
			copy(m.Subject[:], value[1][:])
		default:
			panic(fmt.Sprintf("Unexpected header.SubjectPos %v", m.header.SubjectPos))
		}
	}
	if m.header.Version {
		m.Version = binary.LittleEndian.Uint32(index[0][ClaimTypeLen+ClaimFlagsLen:])
	}
	if m.header.Expiration {
		m.Expiration = int64(binary.LittleEndian.Uint64(value[0][ClaimRevNonceLen:]))
	}
	m.RevNonce = binary.LittleEndian.Uint32(value[0][:])
}

type metadataJSON struct {
	Type       ClaimType
	Subject    ClaimSubject
	SubjectPos ClaimSubjectPos
	ID         *core.ID
	Expiration *int64
	Version    *uint32
	RevNonce   uint32
}

func (m Metadata) MarshalJSON() ([]byte, error) {
	var metadata metadataJSON
	h := m.Header()
	metadata.Type = h.Type
	metadata.SubjectPos = h.SubjectPos
	metadata.Subject = h.Subject
	if h.Subject != ClaimSubjectSelf {
		metadata.ID = m.Subject
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
		Subject:    metadata.Subject,
		SubjectPos: metadata.SubjectPos,
		Expiration: metadata.Expiration != nil,
		Version:    metadata.Version != nil,
	}
	if err := checkHeader(&m.header); err != nil {
		return err
	}
	if m.header.Subject != ClaimSubjectSelf {
		m.Subject = metadata.ID
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
