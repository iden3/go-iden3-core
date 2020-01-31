package merkletree

import (
	"encoding/hex"
	"fmt"

	// "encoding/json"
	// "fmt"
	"bytes"
	"math/big"
	"strings"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-crypto/poseidon"
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

func SwapEndianness(b []byte) []byte {
	o := make([]byte, len(b))
	for i := range b {
		o[len(b)-1-i] = b[i]
	}
	return o
}

func ElemBytesToBigInt(elem ElemBytes) *big.Int {
	return big.NewInt(0).SetBytes(SwapEndianness(elem[:]))
}

func (h1 *Hash) Equals(h2 *Hash) bool {
	return bytes.Equal(h1[:], h2[:])
}

func ElemBytesToBigInts(elems ...ElemBytes) []*big.Int {
	ints := make([]*big.Int, len(elems))
	for i, elem := range elems {
		ints[i] = ElemBytesToBigInt(elem)
	}
	return ints
}

// BigIntToHash converts a *big.Int to a Hash.
func BigIntToHash(e *big.Int) (h Hash) {
	bs := SwapEndianness(e.Bytes())
	copy(h[:], bs)
	return h
}

// HashElems performs a mimc7 hash over the array of ElemBytes.
func HashElems(elems ...ElemBytes) *Hash {
	bigints := ElemBytesToBigInts(elems...)
	// mimcHash, err := mimc7.Hash(bigints, nil)
	poseidonHash, err := poseidon.Hash(bigints)
	if err != nil {
		panic(err)
	}
	h := BigIntToHash(poseidonHash)
	return &h
}

// HashElemsKey performs a mimc7 hash over the array of ElemBytes.
func HashElemsKey(key *big.Int, elems ...ElemBytes) *Hash {
	bigints := ElemBytesToBigInts(elems...)
	// mimcHash, err := mimc7.Hash(bigints, key)
	if key != nil {
		bigints = append(bigints, []*big.Int{key}...)
	}
	poseidonHash, err := poseidon.Hash(bigints)
	if err != nil {
		panic(err)
	}
	h := BigIntToHash(poseidonHash)
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
	h = strings.TrimPrefix(h, "0x")
	return hex.DecodeString(h)
}

func NewEntryFromHexs(a, b, c, d, e, f, g, h string) (entry Entry, err error) {
	entry.Data, err = HexsToData(a, b, c, d, e, f, g, h)
	if err != nil {
		return entry, err
	}
	return entry, nil
}

func HexsToData(_a, _b, _c, _d, _e, _f, _g, _h string) (Data, error) {
	var bints []*big.Int
	aux := []string{_a, _b, _c, _d, _e, _f, _g, _h}
	for _, v := range aux {
		vBytes, err := HexDecode(v)
		if err != nil {
			return Data{}, err
		}
		el := new(big.Int).SetBytes(vBytes)
		bints = append(bints, el)
	}

	return BigIntsToData(bints[0], bints[1], bints[2], bints[3], bints[4], bints[5], bints[6], bints[7]), nil
}

func NewDataFromBytes(b [ElemBytesLen * DataLen]byte) *Data {
	d := &Data{}
	for i := 0; i < DataLen; i++ {
		copy(d[i][:], b[i*ElemBytesLen : (i+1)*ElemBytesLen][:])
	}
	return d
}

func NewEntryFromBytes(b []byte) (*Entry, error) {
	if len(b) != ElemBytesLen*DataLen {
		return nil, fmt.Errorf("Invalid length for Entry Data")
	}
	var data [ElemBytesLen * DataLen]byte
	copy(data[:], b)
	return &Entry{Data: *NewDataFromBytes(data)}, nil
}

func NewEntryFromIntArray(a []int64) Entry {
	return NewEntryFromInts(a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7])
}

func NewEntryFromInts(a, b, c, d, e, f, g, h int64) (entry Entry) {
	entry.Data = IntsToData(a, b, c, d, e, f, g, h)
	return entry
}

func IntArrayToData(a []int64) Data {
	return IntsToData(a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7])
}

func IntsToData(_a, _b, _c, _d, _e, _f, _g, _h int64) Data {
	a, b, c, d, e, f, g, h := big.NewInt(_a), big.NewInt(_b), big.NewInt(_c), big.NewInt(_d), big.NewInt(_e), big.NewInt(_f), big.NewInt(_g), big.NewInt(_h)
	return BigIntsToData(a, b, c, d, e, f, g, h)
}

func BigIntsToData(a, b, c, d, e, f, g, h *big.Int) (data Data) {
	di := []*big.Int{a, b, c, d, e, f, g, h}
	for i := 0; i < len(di); i++ {
		copy(data[i][:], di[i].Bytes())
	}
	return
}
