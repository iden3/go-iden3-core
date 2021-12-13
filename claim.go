package core

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
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

// ErrDataOverflow means that given *big.Int value does not fit in Field Q
// e.g. greater than Q constant:
//
//     Q constant: 21888242871839275222246405745257275088548364400416034343698204186575808495617
var ErrDataOverflow = errors.New("data does not fits SNARK size")

// ErrIncorrectIDPosition means that passed position is not one of predefined:
// IDPositionIndex or IDPositionValue
var ErrIncorrectIDPosition = errors.New("incorrect ID position")

// ErrNoID returns when ID not found in the Claim.
var ErrNoID = errors.New("ID is not set")

// ErrSlotOverflow means some DataSlot overflows Q Field. And wraps the name
// of overflowed slot.
type ErrSlotOverflow struct {
	Field SlotName
}

func (e ErrSlotOverflow) Error() string {
	return fmt.Sprintf("Slot %v not in field (too large)", e.Field)
}

type SlotName string

const (
	SlotNameIndexA = SlotName("IndexA")
	SlotNameIndexB = SlotName("IndexB")
	SlotNameValueA = SlotName("ValueA")
	SlotNameValueB = SlotName("ValueB")
)

const schemaHashLn = 16

// SchemaHash is a 16-bytes hash of file's content, that describes claim
// structure.
type SchemaHash [schemaHashLn]byte

// MarshalText returns HEX representation of SchemaHash.
//
// Returning error is always nil.
func (sc SchemaHash) MarshalText() ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(sc)))
	hex.Encode(dst, sc[:])
	return dst, nil
}

// DataSlot length is 32 bytes. But not all 32-byte values are valid.
// The value should be not greater than Q constant
// 21888242871839275222246405745257275088548364400416034343698204186575808495617
type DataSlot [32]byte

// ToInt returns *big.Int representation of DataSlot.
func (ds DataSlot) ToInt() *big.Int {
	return new(big.Int).SetBytes(utils.SwapEndianness(ds[:]))
}

// SetInt sets data slot to serialized value of *big.Int. And checks that the
// value is valid (fills in Field Q).
// Returns ErrDataOverflow if the value is too large
func (ds *DataSlot) SetInt(value *big.Int) error {
	if !utils.CheckBigIntInField(value) {
		return ErrDataOverflow
	}

	val := utils.SwapEndianness(value.Bytes())
	copy((*ds)[:], val)
	memset((*ds)[len(val):], 0)
	return nil
}

// NewDataSlotFromInt creates new DataSlot from *big.Int.
// Returns error ErrDataOverflow if value is too large to fill the Field Q.
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

type Claim struct {
	index [4]DataSlot
	value [4]DataSlot
}

// Subject for the time being describes the location of ID (in index or value
// slots or nowhere at all).
//
// Values SubjectInvalid, SubjectObjectIndex, SubjectObjectValue
// presents for backward compatibility and for now means nothing.
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
	// IDPositionIndex means ID value is in index slots
	IDPositionIndex IDPosition = iota + 1 // value 0 is position undefined
	// IDPositionValue means ID value is in value slots
	IDPositionValue
)

const (
	flagsByteIdx         = 16
	flagExpirationBitIdx = 3
	flagUpdatableBitIdx  = 4
)

// Option provides the ability to set different Claim's fields on construction
type Option func(*Claim) error

// WithFlagUpdatable sets claim's flag `updatable`
func WithFlagUpdatable(val bool) Option {
	return func(c *Claim) error {
		c.SetFlagUpdatable(val)
		return nil
	}
}

// WithVersion sets claim's version
func WithVersion(ver uint32) Option {
	return func(c *Claim) error {
		c.SetVersion(ver)
		return nil
	}
}

// WithIndexID sets ID to claim's index
func WithIndexID(id ID) Option {
	return func(c *Claim) error {
		c.SetIndexID(id)
		return nil
	}
}

// WithValueID sets ID to claim's value
func WithValueID(id ID) Option {
	return func(c *Claim) error {
		c.SetValueID(id)
		return nil
	}
}

// WithID sets ID to claim's index or value depending on `pos`.
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

// WithRevocationNonce sets claim's revocation nonce.
func WithRevocationNonce(nonce uint64) Option {
	return func(c *Claim) error {
		c.SetRevocationNonce(nonce)
		return nil
	}
}

// WithExpirationDate sets claim's expiration date to `dt`.
func WithExpirationDate(dt time.Time) Option {
	return func(c *Claim) error {
		c.SetExpirationDate(dt)
		return nil
	}
}

// WithIndexData sets data to index slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func WithIndexData(slotA, slotB DataSlot) Option {
	return func(c *Claim) error {
		return c.SetIndexData(slotA, slotB)
	}
}

// WithIndexDataBytes sets data to index slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func WithIndexDataBytes(slotA, slotB []byte) Option {
	return func(c *Claim) error {
		return c.SetIndexDataBytes(slotA, slotB)
	}
}

// WithIndexDataInts sets data to index slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func WithIndexDataInts(slotA, slotB *big.Int) Option {
	return func(c *Claim) error {
		return c.SetIndexDataInts(slotA, slotB)
	}
}

// WithValueData sets data to value slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func WithValueData(slotA, slotB DataSlot) Option {
	return func(c *Claim) error {
		return c.SetValueData(slotA, slotB)
	}
}

// WithValueDataBytes sets data to value slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func WithValueDataBytes(slotA, slotB []byte) Option {
	return func(c *Claim) error {
		return c.SetValueDataBytes(slotA, slotB)
	}
}

// WithValueDataInts sets data to value slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func WithValueDataInts(slotA, slotB *big.Int) Option {
	return func(c *Claim) error {
		return c.SetValueDataInts(slotA, slotB)
	}
}

// NewClaim creates new Claim with specified SchemaHash and any number of
// options. Using options you can specify any field in claim.
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

// SetSchemaHash updates claim's schema hash.
func (c *Claim) SetSchemaHash(schema SchemaHash) {
	copy(c.index[0][:schemaHashLn], schema[:])
}

// GetSchemaHash return copy of claim's schema hash.
func (c *Claim) GetSchemaHash() SchemaHash {
	var schemaHash SchemaHash
	copy(schemaHash[:], c.index[0][:schemaHashLn])
	return schemaHash
}

func (c *Claim) setSubject(s Subject) {
	// clean first 3 bits
	c.index[0][flagsByteIdx] &= 0b11111000
	c.index[0][flagsByteIdx] |= byte(s)
}

func (c *Claim) getSubject() Subject {
	sbj := c.index[0][flagsByteIdx]
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

// SetFlagUpdatable sets claim's flag `updatable`
func (c *Claim) SetFlagUpdatable(val bool) {
	if val {
		c.index[0][flagsByteIdx] |= byte(1) << flagUpdatableBitIdx
	} else {
		c.index[0][flagsByteIdx] &= ^(byte(1) << flagUpdatableBitIdx)
	}
}

// GetFlagUpdatable returns claim's flag `updatable`
func (c *Claim) GetFlagUpdatable() bool {
	mask := byte(1) << flagUpdatableBitIdx
	return c.index[0][flagsByteIdx]&mask > 0
}

// SetVersion sets claim's version
func (c *Claim) SetVersion(ver uint32) {
	binary.LittleEndian.PutUint32(c.index[0][20:24], ver)
}

// GetVersion returns claim's version
func (c *Claim) GetVersion() uint32 {
	return binary.LittleEndian.Uint32(c.index[0][20:24])
}

// SetIndexID sets id to index. Removes id from value if any.
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

// SetValueID sets id to value. Removes id from index if any.
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

// ResetID deletes ID from index and from value.
func (c *Claim) ResetID() {
	c.resetIndexID()
	c.resetValueID()
	c.setSubject(SubjectSelf)
}

// GetID returns ID from claim's index of value.
// Returns error ErrNoID if ID is not set.
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

// SetRevocationNonce sets claim's revocation nonce
func (c *Claim) SetRevocationNonce(nonce uint64) {
	binary.LittleEndian.PutUint64(c.value[0][:8], nonce)
}

// GetRevocationNonce returns revocation nonce
func (c *Claim) GetRevocationNonce() uint64 {
	return binary.LittleEndian.Uint64(c.value[0][:8])
}

// SetExpirationDate sets expiration date to dt
func (c *Claim) SetExpirationDate(dt time.Time) {
	c.setFlagExpiration(true)
	binary.LittleEndian.PutUint64(c.value[0][8:16], uint64(dt.Unix()))
}

// ResetExpirationDate removes expiration date from claim
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
	}
	return time.Time{}, false
}

// SetIndexData sets data to index slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func (c *Claim) SetIndexData(slotA, slotB DataSlot) error {
	slotsAsInts := []*big.Int{slotA.ToInt(), slotB.ToInt()}
	if !utils.CheckBigIntArrayInField(slotsAsInts) {
		return ErrDataOverflow
	}

	copy(c.index[2][:], slotA[:])
	copy(c.index[3][:], slotB[:])
	return nil
}

// SetIndexDataBytes sets data to index slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func (c *Claim) SetIndexDataBytes(slotA, slotB []byte) error {
	err := setSlotBytes(&(c.index[2]), slotA, SlotNameIndexA)
	if err != nil {
		return err
	}
	return setSlotBytes(&(c.index[3]), slotB, SlotNameIndexB)
}

// SetIndexDataInts sets data to index slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func (c *Claim) SetIndexDataInts(slotA, slotB *big.Int) error {
	err := setSlotInt(&c.index[2], slotA, SlotNameIndexA)
	if err != nil {
		return err
	}
	return setSlotInt(&c.index[3], slotB, SlotNameIndexB)
}

// SetValueData sets data to value slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func (c *Claim) SetValueData(slotA, slotB DataSlot) error {
	slotsAsInts := []*big.Int{slotA.ToInt(), slotB.ToInt()}
	if !utils.CheckBigIntArrayInField(slotsAsInts) {
		return ErrDataOverflow
	}

	copy(c.value[2][:], slotA[:])
	copy(c.value[3][:], slotB[:])
	return nil
}

// SetValueDataBytes sets data to value slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func (c *Claim) SetValueDataBytes(slotA, slotB []byte) error {
	err := setSlotBytes(&(c.value[2]), slotA, SlotNameValueA)
	if err != nil {
		return err
	}
	return setSlotBytes(&(c.value[3]), slotB, SlotNameValueB)
}

// SetValueDataInts sets data to value slots A & B.
// Returns ErrSlotOverflow if slotA or slotB value are too big.
func (c *Claim) SetValueDataInts(slotA, slotB *big.Int) error {
	err := setSlotInt(&c.value[2], slotA, SlotNameValueA)
	if err != nil {
		return err
	}
	return setSlotInt(&c.value[3], slotB, SlotNameValueB)
}

func setSlotBytes(slot *DataSlot, value []byte, slotName SlotName) error {
	if len(value) > len(*slot) {
		return ErrSlotOverflow{slotName}
	}
	copy((*slot)[:], value)
	if !utils.CheckBigIntInField(slot.ToInt()) {
		return ErrSlotOverflow{slotName}
	}
	memset((*slot)[len(value):], 0)
	return nil
}

func setSlotInt(slot *DataSlot, value *big.Int, slotName SlotName) error {
	err := slot.SetInt(value)
	if err == ErrDataOverflow {
		return ErrSlotOverflow{slotName}
	}
	return err
}

// TreeEntry creates new merkletree.Entry from the claim. Following changes to
// claim does not change returned merkletree.Entry.
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

// Clone returns full deep copy of claim
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
