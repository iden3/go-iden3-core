package claims

import (
	"encoding/binary"
	"errors"
	"time"

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

var ErrDataOverflow = errors.New("data should not take more then 253 bits")
var ErrIncorrectIDPosition = errors.New("incorrect ID position")

const schemaHashLn = 16

type SchemaHash [schemaHashLn]byte
type ID [31]byte

// DataSlot length is 253 bits, highest 3 bits should be zeros
type DataSlot [32]byte

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
	int253mask           = byte(0b11100000)
)

type Option func(*Claim) error

func WithFlagExpiration(val bool) Option {
	return func(c *Claim) error {
		c.SetFlagExpiration(val)
		return nil
	}
}

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

func WithIndexSlot3(data DataSlot) Option {
	return func(c *Claim) error {
		return c.SetIndexSlot3(data)
	}
}

func WithIndexSlot4(data DataSlot) Option {
	return func(c *Claim) error {
		return c.SetIndexSlot4(data)
	}
}

func WithValueSlot3(data DataSlot) Option {
	return func(c *Claim) error {
		return c.SetValueSlot3(data)
	}
}

func WithValueSlot4(data DataSlot) Option {
	return func(c *Claim) error {
		return c.SetValueSlot4(data)
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

func (c *Claim) setSubject(s Subject) {
	// clean first 3 bits
	c.index[0][9] &= 0b11111000
	c.index[0][9] |= byte(s)
}

func (c *Claim) SetFlagExpiration(val bool) {
	if val {
		c.index[0][flagsByteIdx] |= byte(1) << flagExpirationBitIdx
	} else {
		c.index[0][flagsByteIdx] &= ^(byte(1) << flagExpirationBitIdx)
	}
}

func (c *Claim) SetFlagUpdatable(val bool) {
	if val {
		c.index[0][flagsByteIdx] |= byte(1) << flagUpdatableBitIdx
	} else {
		c.index[0][flagsByteIdx] &= ^(byte(1) << flagUpdatableBitIdx)
	}
}

func (c *Claim) SetVersion(ver uint32) {
	binary.LittleEndian.PutUint32(c.index[0][20:24], ver)
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

func (c *Claim) SetValueID(id ID) {
	c.resetIndexID()
	c.setSubject(SubjectOtherIdenValue)
	copy(c.value[1][:], id[:])
}

func (c *Claim) resetValueID() {
	var zeroID ID
	copy(c.value[1][:], zeroID[:])
}

func (c *Claim) ResetID() {
	c.resetIndexID()
	c.resetValueID()
	c.setSubject(SubjectSelf)
}

func (c *Claim) SetRevocationNonce(nonce uint64) {
	binary.LittleEndian.PutUint64(c.value[0][:8], nonce)
}

func (c *Claim) SetExpirationDate(dt time.Time) {
	binary.LittleEndian.PutUint64(c.value[0][8:16], uint64(dt.Unix()))
}

func (c *Claim) SetIndexSlot3(data DataSlot) error {
	if !isInt253compatible(data) {
		return ErrDataOverflow
	}
	copy(c.index[2][:], data[:])
	return nil
}

func (c *Claim) SetIndexSlot4(data DataSlot) error {
	if !isInt253compatible(data) {
		return ErrDataOverflow
	}
	copy(c.index[3][:], data[:])
	return nil
}

func (c *Claim) SetValueSlot3(data DataSlot) error {
	if !isInt253compatible(data) {
		return ErrDataOverflow
	}
	copy(c.value[2][:], data[:])
	return nil
}

func (c *Claim) SetValueSlot4(data DataSlot) error {
	if !isInt253compatible(data) {
		return ErrDataOverflow
	}
	copy(c.value[3][:], data[:])
	return nil
}

func isInt253compatible(data DataSlot) bool {
	return data[len(data)-1]&int253mask == 0
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
