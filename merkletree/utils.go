package merkletree

import (
	"encoding/hex"
	// "encoding/json"
	// "fmt"
	"math/big"
	"strings"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-crypto/mimc7"
)

// Hash is the type used to represent a hash used in the MT.
type Hash ElemBytes

// String returns the last 4 bytes of Hash in hex.
func (h *Hash) String() string {
	//return hex.EncodeToString(h[ElemBytesLen-4:])
	return (*ElemBytes)(h).String()
}

// Hex returns a hex string from the Hash type.
func (h Hash) Hex() string {
	return common3.HexEncode(h[:])
}

// Bytes returns a byte array from a Hash.
func (h Hash) Bytes() []byte {
	return h[:]
}

func (h *Hash) MarshalText() ([]byte, error) {
	return []byte(common3.HexEncode(h.Bytes())), nil
}

func (h *Hash) UnmarshalText(bs []byte) error {
	return common3.HexDecodeInto(h[:], bs)
}

func ElemBytesToBigInt(elem ElemBytes) *big.Int {
	return big.NewInt(0).SetBytes(elem[:])
}

func ElemBytesToBigInts(elems ...ElemBytes) []*big.Int {
	ints := make([]*big.Int, len(elems))
	for i, elem := range elems {
		ints[i] = big.NewInt(0).SetBytes(elem[:])
	}
	return ints
}

// BigIntToHash converts a *big.Int to a Hash.
func BigIntToHash(e *big.Int) (h Hash) {
	bs := e.Bytes()
	copy(h[ElemBytesLen-len(bs):], bs)
	return h
}

// HashElems performs a mimc7 hash over the array of ElemBytes.
func HashElems(elems ...ElemBytes) *Hash {
	bigints := ElemBytesToBigInts(elems...)
	mimcHash, err := mimc7.Hash(bigints, nil)
	if err != nil {
		panic(err)
	}
	h := BigIntToHash(mimcHash)
	return &h
}

// HashElemsKey performs a mimc7 hash over the array of ElemBytes.
func HashElemsKey(key *big.Int, elems ...ElemBytes) *Hash {
	bigints := ElemBytesToBigInts(elems...)
	mimcHash, err := mimc7.Hash(bigints, key)
	if err != nil {
		panic(err)
	}
	h := BigIntToHash(mimcHash)
	return &h
}

// getPath returns the binary path, from the root to the leaf.
func getPath(numLevels int, hIndex *Hash) []bool {
	path := make([]bool, numLevels)
	for n := 0; n < numLevels; n++ {
		path[n] = testBitBigEndian(hIndex[:], uint(n))
	}
	return path
}

// setBit sets the bit n in the bitmap to 1.
func setBit(bitmap []byte, n uint) {
	bitmap[n/8] |= 1 << (n % 8)
}

// setBitBigEndian sets the bit n in the bitmap to 1, in Big Endian.
func setBitBigEndian(bitmap []byte, n uint) {
	bitmap[uint(len(bitmap))-n/8-1] |= 1 << (n % 8)
}

// testBit tests whether the bit n in bitmap is 1.
func testBit(bitmap []byte, n uint) bool {
	return bitmap[n/8]&(1<<(n%8)) != 0
}

// testBitBigEndian tests whether the bit n in bitmap is 1, in Big Endian.
func testBitBigEndian(bitmap []byte, n uint) bool {
	return bitmap[uint(len(bitmap))-n/8-1]&(1<<(n%8)) != 0
}

func HexDecode(h string) ([]byte, error) {
	if strings.HasPrefix(h, "0x") {
		h = h[2:]
	}
	return hex.DecodeString(h)
}
func NewEntryFromHexs(a, b, c, d string) (e Entry, err error) {
	e.Data, err = HexsToData(a, b, c, d)
	if err != nil {
		return e, err
	}
	return e, nil
}
func HexsToData(_a, _b, _c, _d string) (Data, error) {
	aBytes, err := HexDecode(_a)
	if err != nil {
		return Data{}, err
	}
	a := new(big.Int).SetBytes(aBytes)

	bBytes, err := HexDecode(_b)
	if err != nil {
		return Data{}, err
	}
	b := new(big.Int).SetBytes(bBytes)

	cBytes, err := HexDecode(_c)
	if err != nil {
		return Data{}, err
	}
	c := new(big.Int).SetBytes(cBytes)

	dBytes, err := HexDecode(_d)
	if err != nil {
		return Data{}, err
	}
	d := new(big.Int).SetBytes(dBytes)

	return BigIntsToData(a, b, c, d), nil
}

func NewEntryFromInts(a, b, c, d int64) (e Entry) {
	e.Data = IntsToData(a, b, c, d)
	return e
}

func IntsToData(_a, _b, _c, _d int64) Data {
	a, b, c, d := big.NewInt(_a), big.NewInt(_b), big.NewInt(_c), big.NewInt(_d)
	return BigIntsToData(a, b, c, d)
}

func BigIntsToData(a, b, c, d *big.Int) (data Data) {
	di := []*big.Int{a, b, c, d}
	for i, v := range di {
		copy(data[i][(ElemBytesLen-len(v.Bytes())):], v.Bytes())
	}
	return
}
