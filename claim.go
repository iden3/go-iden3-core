package core

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-iden3-crypto/utils"
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
          [3] Merklized: data is merklized root is stored in the:
            000: none
            001: C.i Root Index (root located in i_2)
            010: C.v Root Value (root located in v_2)
          [24] 0
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
//	Q constant: 21888242871839275222246405745257275088548364400416034343698204186575808495617
var ErrDataOverflow = errors.New("data does not fits SNARK size")

// ErrIncorrectIDPosition means that passed position is not one of predefined:
// IDPositionIndex or IDPositionValue
var ErrIncorrectIDPosition = errors.New("incorrect ID position")

// ErrIncorrectMerklizedPosition means that passed position is not one of predefined:
// MerklizedRootPositionIndex or MerklizedRootPositionValue
var ErrIncorrectMerklizedPosition = errors.New("incorrect Merklized position")

// ErrNoID returns when ID not found in the Claim.
var ErrNoID = errors.New("ID is not set")

// ErrNoMerklizedRoot returns when Merklized Root is not found in the Claim.
var ErrNoMerklizedRoot = errors.New("Merklized root is not set")

// ErrInvalidSubjectPosition returns when subject position flags sets in invalid value.
var ErrInvalidSubjectPosition = errors.New("invalid subject position")

// ErrSlotOverflow means some ElemBytes overflows Q Field. And wraps the name
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

// AuthSchemaHash predefined value of auth schema, used for auth claim during identity creation.
// This schema is hardcoded in the identity circuits and used to verify user's auth claim.
// Keccak256(https://schema.iden3.io/core/jsonld/auth.jsonld#AuthBJJCredential) last 16 bytes
// Hex: cca3371a6cb1b715004407e325bd993c
// BigInt: 80551937543569765027552589160822318028
var AuthSchemaHash = SchemaHash{204, 163, 55, 26, 108, 177, 183, 21, 0, 68, 7, 227, 37, 189, 153, 60}

const schemaHashLn = 16

// SchemaHash is a 16-bytes hash of file's content, that describes claim
// structure.
type SchemaHash [schemaHashLn]byte

// MarshalText returns HEX representation of SchemaHash.
//
// Returning error is always nil.
func (sh SchemaHash) MarshalText() ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(sh)))
	hex.Encode(dst, sh[:])
	return dst, nil
}

// NewSchemaHashFromHex creates new SchemaHash from hex string
func NewSchemaHashFromHex(s string) (SchemaHash, error) {
	var sh SchemaHash
	schemaEncodedBytes, err := hex.DecodeString(s)
	if err != nil {
		return SchemaHash{}, err
	}

	if len(schemaEncodedBytes) != len(sh) {
		return SchemaHash{}, fmt.Errorf("invalid schema hash length: %d",
			len(schemaEncodedBytes))
	}
	copy(sh[:], schemaEncodedBytes)

	return sh, nil
}

// NewSchemaHashFromInt creates new SchemaHash from big.Int
func NewSchemaHashFromInt(i *big.Int) SchemaHash {
	var sh SchemaHash
	b := intToBytes(i)
	copy(sh[:], b)

	return sh
}

// BigInt returns a bigInt presentation of SchemaHash
func (sh SchemaHash) BigInt() *big.Int {
	return bytesToInt(sh[:])
}

type Claim struct {
	index [4]ElemBytes
	value [4]ElemBytes
}

// NewClaimFromBigInts creates new Claim from bigInts.
func NewClaimFromBigInts(raw [8]*big.Int) (*Claim, error) {
	var c Claim
	for i := 0; i < 4; i++ {
		err := c.index[i].SetInt(raw[i])
		if err != nil {
			return nil, err
		}
		err = c.value[i].SetInt(raw[i+4])
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

// subjectFlag for the time being describes the location of ID (in index or value
// slots or nowhere at all).
//
// Values subjectFlagInvalid presents for backward compatibility and for now means nothing.
type subjectFlag byte

const (
	subjectFlagSelf           subjectFlag = iota // 000
	_subjectFlagInvalid                          // nolint // 001
	subjectFlagOtherIdenIndex                    // 010
	subjectFlagOtherIdenValue                    // 011
)

type IDPosition uint8

const (
	// IDPositionNone means ID value not located in claim.
	IDPositionNone IDPosition = iota
	// IDPositionIndex means ID value is in index slots.
	IDPositionIndex
	// IDPositionValue means ID value is in value slots.
	IDPositionValue
)

// merklizedFlag for the time being describes the location of root (in index or value
// slots or nowhere at all).
//
// Values merklizedFlagIndex indicates that root is located in index[2] slots.
// Values merklizedFlagValue indicates that root is located in value[2] slots.
type merklizedFlag byte

const (
	merklizedFlagNone     merklizedFlag = 0b00000000 // 000 00000
	merklizedFlagIndex    merklizedFlag = 0b00100000 // 001 00000
	merklizedFlagValue    merklizedFlag = 0b01000000 // 010 00000
	_merklizedFlagInvalid merklizedFlag = 0b10000000 // 010 00000
)

type MerklizedRootPosition uint8

const (
	// MerklizedRootPositionNone means root data value not located in claim.
	MerklizedRootPositionNone MerklizedRootPosition = iota
	// MerklizedRootPositionIndex means root data value is in index slots.
	MerklizedRootPositionIndex
	// MerklizedRootPositionValue means root data value is in value slots.
	MerklizedRootPositionValue
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

// WithFlagMerklized sets claim's flag `merklize`
func WithFlagMerklized(p MerklizedRootPosition) Option {
	return func(c *Claim) error {
		c.setFlagMerklized(p)
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
func WithIndexData(slotA, slotB ElemBytes) Option {
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
func WithValueData(slotA, slotB ElemBytes) Option {
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

// WithIndexMerklizedRoot sets root to index i_2
// Returns ErrSlotOverflow if root value are too big.
func WithIndexMerklizedRoot(r *big.Int) Option {
	return func(c *Claim) error {
		c.setFlagMerklized(MerklizedRootPositionIndex)
		return setSlotInt(&c.index[2], r, SlotNameIndexA)
	}
}

// WithValueMerklizedRoot sets root to value v_2
// Returns ErrSlotOverflow if root value are too big.
func WithValueMerklizedRoot(r *big.Int) Option {
	return func(c *Claim) error {
		c.setFlagMerklized(MerklizedRootPositionValue)
		return setSlotInt(&c.value[2], r, SlotNameValueA)
	}
}

// NewClaim creates new Claim with specified SchemaHash and any number of
// options. Using options you can specify any field in claim.
func NewClaim(sh SchemaHash, options ...Option) (*Claim, error) {
	c := &Claim{}
	c.SetSchemaHash(sh)
	for _, o := range options {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithMerklizedRoot sets root to value v_2 or index i_2
// Returns ErrSlotOverflow if root value are too big.
func WithMerklizedRoot(r *big.Int, pos MerklizedRootPosition) Option {
	return func(c *Claim) error {
		switch pos {
		case MerklizedRootPositionIndex:
			c.setFlagMerklized(MerklizedRootPositionIndex)
			return setSlotInt(&c.index[2], r, SlotNameIndexA)
		case MerklizedRootPositionValue:
			c.setFlagMerklized(MerklizedRootPositionValue)
			return setSlotInt(&c.value[2], r, SlotNameValueA)
		default:
			return ErrIncorrectMerklizedPosition
		}
	}
}

// HIndex calculates the hash of the Index of the Claim
func (c *Claim) HIndex() (*big.Int, error) {
	return poseidon.Hash(ElemBytesToInts(c.index[:]))
}

// HValue calculates the hash of the Value of the Claim
func (c *Claim) HValue() (*big.Int, error) {
	return poseidon.Hash(ElemBytesToInts(c.value[:]))
}

// HiHv returns the HIndex and HValue of the Claim
func (c *Claim) HiHv() (*big.Int, *big.Int, error) {
	hi, err := c.HIndex()
	if err != nil {
		return nil, nil, err
	}
	hv, err := c.HValue()
	if err != nil {
		return nil, nil, err
	}

	return hi, hv, nil
}

// SetSchemaHash updates claim's schema hash.
func (c *Claim) SetSchemaHash(sh SchemaHash) {
	copy(c.index[0][:schemaHashLn], sh[:])
}

// GetSchemaHash return copy of claim's schema hash.
func (c *Claim) GetSchemaHash() SchemaHash {
	var sh SchemaHash
	copy(sh[:], c.index[0][:schemaHashLn])
	return sh
}

// GetIDPosition returns the position at which the ID is stored.
func (c *Claim) GetIDPosition() (IDPosition, error) {
	switch c.getSubject() {
	case subjectFlagSelf:
		return IDPositionNone, nil
	case subjectFlagOtherIdenIndex:
		return IDPositionIndex, nil
	case subjectFlagOtherIdenValue:
		return IDPositionValue, nil
	default:
		return 0, ErrInvalidSubjectPosition
	}
}

func (c *Claim) setSubject(s subjectFlag) {
	// clean first 3 bits
	c.index[0][flagsByteIdx] &= 0b11111000
	c.index[0][flagsByteIdx] |= byte(s)
}

// setFlagMerklized sets the merklized flag in the claim
func (c *Claim) setFlagMerklized(s MerklizedRootPosition) {
	var f merklizedFlag
	switch s {
	case MerklizedRootPositionIndex:
		f = merklizedFlagIndex
	case MerklizedRootPositionValue:
		f = merklizedFlagValue
	default:
		f = merklizedFlagNone
	}
	// clean last 3 bits
	c.index[0][flagsByteIdx] &= 0b00011111
	c.index[0][flagsByteIdx] |= byte(f)
}

func (c *Claim) getSubject() subjectFlag {
	sbj := c.index[0][flagsByteIdx]
	// clean all except first 3 bits
	sbj &= 0b00000111
	return subjectFlag(sbj)
}

func (c *Claim) getMerklized() merklizedFlag {
	mt := c.index[0][flagsByteIdx]
	// clean all except last 3 bits
	mt &= 0b11100000
	return merklizedFlag(mt)
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

// SetIndexMerklizedRoot sets merklized root to index. Removes root from value[2] if any.
func (c *Claim) SetIndexMerklizedRoot(r *big.Int) error {
	c.resetValueMerklizedRoot()
	c.setFlagMerklized(MerklizedRootPositionIndex)
	return setSlotInt(&c.index[2], r, SlotNameIndexA)
}
func (c *Claim) resetIndexMerklizedRoot() {
	var zeroBytes ElemBytes
	copy(c.index[2][:], zeroBytes[:])
}

// SetValueMerklizedRoot sets merklized root to value. Removes root from index[2] if any.
func (c *Claim) SetValueMerklizedRoot(r *big.Int) error {
	c.resetIndexMerklizedRoot()
	c.setFlagMerklized(MerklizedRootPositionValue)
	return setSlotInt(&c.value[2], r, SlotNameValueA)
}

func (c *Claim) resetValueMerklizedRoot() {
	var zeroBytes ElemBytes
	copy(c.value[2][:], zeroBytes[:])
}

// SetIndexID sets id to index. Removes id from value if any.
func (c *Claim) SetIndexID(id ID) {
	c.resetValueID()
	c.setSubject(subjectFlagOtherIdenIndex)
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
	c.setSubject(subjectFlagOtherIdenValue)
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
	c.setSubject(subjectFlagSelf)
}

// GetID returns ID from claim's index of value.
// Returns error ErrNoID if ID is not set.
func (c *Claim) GetID() (ID, error) {
	var id ID
	switch c.getSubject() {
	case subjectFlagOtherIdenIndex:
		return c.getIndexID(), nil
	case subjectFlagOtherIdenValue:
		return c.getValueID(), nil
	default:
		return id, ErrNoID
	}
}

// GetMerklizedRoot returns merklized root from claim's index of value.
// Returns error ErrNoMerklizedRoot if MerklizedRoot is not set.
func (c *Claim) GetMerklizedRoot() (*big.Int, error) {
	switch c.getMerklized() {
	case merklizedFlagIndex:
		return c.index[2].ToInt(), nil
	case merklizedFlagValue:
		return c.value[2].ToInt(), nil
	default:
		return nil, ErrNoMerklizedRoot
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
func (c *Claim) SetIndexData(slotA, slotB ElemBytes) error {
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
func (c *Claim) SetValueData(slotA, slotB ElemBytes) error {
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

func setSlotBytes(slot *ElemBytes, value []byte, slotName SlotName) error {
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

func setSlotInt(slot *ElemBytes, value *big.Int, slotName SlotName) error {
	if value == nil {
		value = big.NewInt(0)
	}

	err := slot.SetInt(value)
	if err == ErrDataOverflow {
		return ErrSlotOverflow{slotName}
	}
	return err
}

// RawSlots returns raw bytes of claim's index and value
func (c *Claim) RawSlots() (index [4]ElemBytes, value [4]ElemBytes) {
	return c.index, c.value
}

// RawSlotsAsInts returns slots as []*big.Int
func (c *Claim) RawSlotsAsInts() []*big.Int {
	return append(ElemBytesToInts(c.index[:]), ElemBytesToInts(c.value[:])...)
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

func (c Claim) MarshalJSON() ([]byte, error) {
	intVals := c.RawSlotsAsInts()
	var obj = make([]string, len(intVals))
	for i, v := range intVals {
		obj[i] = v.Text(10)
	}
	return json.Marshal(obj)
}

func (c *Claim) UnmarshalJSON(in []byte) error {
	var sVals []string
	err := json.Unmarshal(in, &sVals)
	if err != nil {
		return err
	}

	if len(sVals) != len(c.index)+len(c.value) {
		return errors.New("invalid number of claim's slots")
	}

	var (
		intVal *big.Int
		ok     bool
	)

	for i := 0; i < len(c.index); i++ {
		intVal, ok = new(big.Int).SetString(sVals[i], 10)
		if !ok {
			return fmt.Errorf("can't parse int for index field #%v", i)
		}
		err = c.index[i].SetInt(intVal)
		if err != nil {
			return fmt.Errorf("can't set index field #%v: %w", i, err)
		}
	}

	for i := 0; i < len(c.value); i++ {
		intVal, ok = new(big.Int).SetString(sVals[i+len(c.index)], 10)
		if !ok {
			return fmt.Errorf("can't parse int for value field #%v", i)
		}
		err = c.value[i].SetInt(intVal)
		if err != nil {
			return fmt.Errorf("can't set value field #%v: %w", i, err)
		}
	}

	return nil
}

// Hex returns hex representation of binary claim
func (c Claim) Hex() (string, error) {
	b, err := c.MarshalBinary()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), err
}

func (c Claim) MarshalBinary() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	buf.Grow(len(c.index)*len(c.index[0]) + len(c.value)*len(c.value[0]))
	for i := range c.index {
		buf.Write(c.index[i][:])
	}
	for i := range c.value {
		buf.Write(c.value[i][:])
	}
	return buf.Bytes(), nil
}

func (c *Claim) UnmarshalBinary(data []byte) error {
	wantLen := len(c.index)*len(c.index[0]) + len(c.value)*len(c.value[0])
	if len(data) != wantLen {
		return errors.New("unexpected length of input data")
	}

	offset := 0
	for i := range c.index {
		copy(c.index[i][:], data[offset:])
		offset += len(c.index[i])
		_, err := fieldBytesToInt(c.index[i][:])
		if err != nil {
			return fmt.Errorf("can't set index slot #%v: %w", i, err)
		}
	}
	for i := range c.value {
		copy(c.value[i][:], data[offset:])
		offset += len(c.value[i])
		_, err := fieldBytesToInt(c.value[i][:])
		if err != nil {
			return fmt.Errorf("can't set value slot #%v: %w", i, err)
		}
	}

	return nil
}

func (c *Claim) FromHex(hexStr string) error {

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	return c.UnmarshalBinary(data)
}

// GetMerklizedPosition returns the position at which the Merklize flag is stored.
func (c *Claim) GetMerklizedPosition() (MerklizedRootPosition, error) {
	switch c.getMerklized() {
	case merklizedFlagNone:
		return MerklizedRootPositionNone, nil
	case merklizedFlagIndex:
		return MerklizedRootPositionIndex, nil
	case merklizedFlagValue:
		return MerklizedRootPositionValue, nil
	default:
		return 0, ErrIncorrectMerklizedPosition
	}
}
