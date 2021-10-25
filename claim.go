package core

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/iden3/go-merkletree-sql"
)

/*
Claim structure

Index:
 i_0: [ 128  bits ] claim schema
      [ 32 bits ] option flags
          [3] Subject:
            000: A.1 Self
            001: invalid
            010: A.2.i OtherIden Index
            011: A.2.v OtherIden Value
            100: B.i Object Index
            101: B.v Object Value
          [1] Expiration: bool
          [1] Updatable: bool
          [27] 0
      [ 32 bits ] version (optional?)
      [ 61 bits ] 0 - reserved for future use
 i_1: [ 248 bits] identity (case b) (optional)
      [  5 bits ] 0
 i_2: [ 253 bits] 0
 i_3: [ 253 bits] 0
Value:
 v_0: [ 64 bits ]  revocation nonce
      [ 64 bits ]  expiration date (optional)
      [ 125 bits] 0 - reserved
 v_1: [ 248 bits] identity (case c) (optional)
      [  5 bits ] 0
 v_2: [ 253 bits] 0
 v_3: [ 253 bits] 0
*/

var ErrDataOverflow = errors.New("data does not fits SNARK size")
var ErrIncorrectIDPosition = errors.New("incorrect ID position")
var ErrNoID = errors.New("ID is not set")

const schemaHashLn = 16

type SchemaHash [schemaHashLn]byte

func (sc SchemaHash) MarshalText() ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(sc)))
	hex.Encode(dst, sc[:])
	return dst, nil
}

// DataSlot length is 253 bits, highest 3 bits should be zeros
type DataSlot [32]byte

func (ds DataSlot) ToInt() *big.Int {
	return new(big.Int).SetBytes(utils.SwapEndianness(ds[:]))
}

func NewDataSlotFromInt(i *big.Int) (DataSlot, error) {
	var s DataSlot
	bs := i.Bytes()
	// may be this check is redundant because of CheckBigIntInField, but just
	// in case.
	if len(bs) > len(s) {
		return s, ErrDataOverflow
	}
	if !utils.CheckBigIntInField(i) {
		return s, ErrDataOverflow
	}
	copy(s[:], utils.SwapEndianness(bs))
	return s, nil
}

type int253 [32]byte

type Claim struct {
	index [4]int253
	value [4]int253
}

type Subject byte

const (
	SubjectSelf           Subject = iota // 000
	SubjectInvalid                       // 001
	SubjectOtherIdenIndex                // 010
	SubjectOtherIdenValue                // 011
	SubjectObjectIndex                   // 100
	SubjectObjectValue                   // 101
)

type IDPosition uint8

const (
	idPositionUndefined IDPosition = iota
	IDPositionIndex
	IDPositionValue
)

const (
	flagsByteIdx         = 16
	flagExpirationBitIdx = 3
	flagUpdatableBitIdx  = 4
)

type Option func(*Claim) error

func WithFlagUpdatable(val bool) Option {
	return func(c *Claim) error {
		c.SetFlagUpdatable(val)
		return nil
	}
}

func WithVersion(ver uint32) Option {
	return func(c *Claim) error {
		c.SetVersion(ver)
		return nil
	}
}

func WithIndexID(id ID) Option {
	return func(c *Claim) error {
		c.SetIndexID(id)
		return nil
	}
}

func WithValueID(id ID) Option {
	return func(c *Claim) error {
		c.SetValueID(id)
		return nil
	}
}

func WithID(id ID, pos IDPosition) Option {
	return func(c *Claim) error {
		switch pos {
		case IDPositionIndex:
			c.SetIndexID(id)
		case IDPositionValue:
			c.SetValueID(id)
		default:
			return ErrIncorrectIDPosition
		}
		return nil
	}
}

func WithRevocationNonce(nonce uint64) Option {
	return func(c *Claim) error {
		c.SetRevocationNonce(nonce)
		return nil
	}
}

func WithExpirationDate(dt time.Time) Option {
	return func(c *Claim) error {
		c.SetExpirationDate(dt)
		return nil
	}
}

func WithIndexData(slotA, slotB DataSlot) Option {
	return func(c *Claim) error {
		return c.SetIndexData(slotA, slotB)
	}
}

func WithValueData(slotA, slotB DataSlot) Option {
	return func(c *Claim) error {
		return c.SetValueData(slotA, slotB)
	}
}

func NewClaim(schemaHash SchemaHash, options ...Option) (*Claim, error) {
	c := &Claim{}
	c.SetSchemaHash(schemaHash)
	for _, o := range options {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Claim) SetSchemaHash(schema SchemaHash) {
	copy(c.index[0][:schemaHashLn], schema[:])
}

func (c *Claim) GetSchemaHash() SchemaHash {
	var schemaHash SchemaHash
	copy(schemaHash[:], c.index[0][:schemaHashLn])
	return schemaHash
}

func (c *Claim) setSubject(s Subject) {
	// clean first 3 bits
	c.index[0][9] &= 0b11111000
	c.index[0][9] |= byte(s)
}

func (c *Claim) getSubject() Subject {
	sbj := c.index[0][9]
	// clean all except first 3 bits
	sbj &= 0b00000111
	return Subject(sbj)
}

func (c *Claim) setFlagExpiration(val bool) {
	if val {
		c.index[0][flagsByteIdx] |= byte(1) << flagExpirationBitIdx
	} else {
		c.index[0][flagsByteIdx] &= ^(byte(1) << flagExpirationBitIdx)
	}
}

func (c *Claim) getFlagExpiration() bool {
	mask := byte(1) << flagExpirationBitIdx
	return c.index[0][flagsByteIdx]&mask > 0
}

func (c *Claim) SetFlagUpdatable(val bool) {
	if val {
		c.index[0][flagsByteIdx] |= byte(1) << flagUpdatableBitIdx
	} else {
		c.index[0][flagsByteIdx] &= ^(byte(1) << flagUpdatableBitIdx)
	}
}

func (c *Claim) GetFlagUpdatable() bool {
	mask := byte(1) << flagUpdatableBitIdx
	return c.index[0][flagsByteIdx]&mask > 0
}

func (c *Claim) SetVersion(ver uint32) {
	binary.LittleEndian.PutUint32(c.index[0][20:24], ver)
}

func (c *Claim) GetVersion() uint32 {
	return binary.LittleEndian.Uint32(c.index[0][20:24])
}

func (c *Claim) SetIndexID(id ID) {
	c.resetValueID()
	c.setSubject(SubjectOtherIdenIndex)
	copy(c.index[1][:], id[:])
}

func (c *Claim) resetIndexID() {
	var zeroID ID
	copy(c.index[1][:], zeroID[:])
}

func (c *Claim) getIndexID() ID {
	var id ID
	copy(id[:], c.index[1][:])
	return id
}

func (c *Claim) SetValueID(id ID) {
	c.resetIndexID()
	c.setSubject(SubjectOtherIdenValue)
	copy(c.value[1][:], id[:])
}

func (c *Claim) resetValueID() {
	var zeroID ID
	copy(c.value[1][:], zeroID[:])
}

func (c *Claim) getValueID() ID {
	var id ID
	copy(id[:], c.value[1][:])
	return id
}

func (c *Claim) ResetID() {
	c.resetIndexID()
	c.resetValueID()
	c.setSubject(SubjectSelf)
}

func (c *Claim) GetID() (ID, error) {
	var id ID
	switch c.getSubject() {
	case SubjectOtherIdenIndex:
		return c.getIndexID(), nil
	case SubjectOtherIdenValue:
		return c.getValueID(), nil
	default:
		return id, ErrNoID
	}
}

func (c *Claim) SetRevocationNonce(nonce uint64) {
	binary.LittleEndian.PutUint64(c.value[0][:8], nonce)
}

func (c *Claim) GetRevocationNonce() uint64 {
	return binary.LittleEndian.Uint64(c.value[0][:8])
}

func (c *Claim) SetExpirationDate(dt time.Time) {
	c.setFlagExpiration(true)
	binary.LittleEndian.PutUint64(c.value[0][8:16], uint64(dt.Unix()))
}

func (c *Claim) ResetExpirationDate() {
	c.setFlagExpiration(false)
	memset(c.value[0][8:16], 0)
}

// GetExpirationDate returns expiration date and flag. Flag is true if
// expiration date is present, false if null.
func (c *Claim) GetExpirationDate() (time.Time, bool) {
	if c.getFlagExpiration() {
		expirationDate :=
			time.Unix(int64(binary.LittleEndian.Uint64(c.value[0][8:16])), 0)
		return expirationDate, true
	} else {
		return time.Time{}, false
	}
}

func (c *Claim) SetIndexData(slotA, slotB DataSlot) error {
	slotsAsInts := []*big.Int{slotA.ToInt(), slotB.ToInt()}
	if !utils.CheckBigIntArrayInField(slotsAsInts) {
		return ErrDataOverflow
	}

	copy(c.index[2][:], slotA[:])
	copy(c.index[3][:], slotB[:])
	return nil
}

func (c *Claim) SetValueData(slotA, slotB DataSlot) error {
	slotsAsInts := []*big.Int{slotA.ToInt(), slotB.ToInt()}
	if !utils.CheckBigIntArrayInField(slotsAsInts) {
		return ErrDataOverflow
	}

	copy(c.value[2][:], slotA[:])
	copy(c.value[3][:], slotB[:])
	return nil
}

func (c *Claim) TreeEntry() merkletree.Entry {
	var e merkletree.Entry
	for i := range c.index {
		copy(e.Data[i][:], c.index[i][:])
	}
	for i := range c.value {
		copy(e.Data[i+len(c.index)][:], c.value[i][:])
	}
	return e
}

func (c *Claim) Clone() *Claim {
	var newClaim Claim
	for i := range c.index {
		copy(newClaim.index[i][:], c.index[i][:])
	}
	for i := range c.value {
		copy(newClaim.value[i][:], c.value[i][:])
	}
	return &newClaim
}

func memset(arr []byte, v byte) {
	if len(arr) == 0 {
		return
	}
	arr[0] = v
	for ptr := 1; ptr < len(arr); ptr *= 2 {
		copy(arr[ptr:], arr[:ptr])
	}
}
