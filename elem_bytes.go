package core

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/utils"
)

// ElemBytes length is 32 bytes. But not all 32-byte values are valid.
// The value should be not greater than Q constant
// 21888242871839275222246405745257275088548364400416034343698204186575808495617
type ElemBytes [32]byte

// ToInt returns *big.Int representation of ElemBytes.
func (el ElemBytes) ToInt() *big.Int {
	return new(big.Int).SetBytes(utils.SwapEndianness(el[:]))
}

// SetInt sets element's data to serialized value of *big.Int in little-endian.
// And checks that the value is valid (fits in Field Q).
// Returns ErrDataOverflow if the value is too large
func (el *ElemBytes) SetInt(value *big.Int) error {
	if !utils.CheckBigIntInField(value) {
		return ErrDataOverflow
	}

	val := utils.SwapEndianness(value.Bytes())
	copy((*el)[:], val)
	memset((*el)[len(val):], 0)
	return nil
}

// NewElementBytesFromInt creates new ElemBytes from *big.Int.
// Returns error ErrDataOverflow if value is too large to fill the Field Q.
func NewElementBytesFromInt(i *big.Int) (ElemBytes, error) {
	var s ElemBytes
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

// ElemBytesToInts converts slice of ElemBytes to slice of *big.Int
func ElemBytesToInts(elements []ElemBytes) []*big.Int {
	result := make([]*big.Int, len(elements))
	for i := range elements {
		result[i] = elements[i].ToInt()
	}
	return result
}
