package core

import (
	"encoding/hex"
	"math/big"

	"github.com/iden3/go-iden3-crypto/utils"
)

// ElemBytes length is 32 bytes. But not all 32-byte values are valid.
// The value should be not greater than Q constant
// 21888242871839275222246405745257275088548364400416034343698204186575808495617
type ElemBytes [32]byte

// ToInt returns *big.Int representation of ElemBytes.
func (el ElemBytes) ToInt() *big.Int {
	return bytesToInt(el[:])
}

// SetInt sets element's data to serialized value of *big.Int in little-endian.
// And checks that the value is valid (fits in Field Q).
// Returns ErrDataOverflow if the value is too large
func (el *ElemBytes) SetInt(value *big.Int) error {
	val, err := fieldIntToBytes(value)
	if err != nil {
		return err
	}
	copy((*el)[:], val)
	memset((*el)[len(val):], 0)
	return nil
}

// Hex returns HEX representation of ElemBytes
func (el ElemBytes) Hex() string {
	return hex.EncodeToString(el[:])
}

// NewElemBytesFromInt creates new ElemBytes from *big.Int.
// Returns error ErrDataOverflow if value is too large to fill the Field Q.
func NewElemBytesFromInt(i *big.Int) (ElemBytes, error) {
	val, err := fieldIntToBytes(i)
	if err != nil {
		return ElemBytes{}, err
	}
	var s ElemBytes
	copy(s[:], val)
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

func bytesToInt(in []byte) *big.Int {
	return new(big.Int).SetBytes(utils.SwapEndianness(in))
}

func fieldBytesToInt(in []byte) (*big.Int, error) {
	i := bytesToInt(in)
	if !utils.CheckBigIntInField(i) {
		return nil, ErrDataOverflow
	}

	return i, nil
}

func intToBytes(in *big.Int) []byte {
	return utils.SwapEndianness(in.Bytes())
}

func fieldIntToBytes(in *big.Int) ([]byte, error) {
	if !utils.CheckBigIntInField(in) {
		return nil, ErrDataOverflow
	}

	return intToBytes(in), nil
}
