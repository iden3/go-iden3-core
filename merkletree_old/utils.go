package merkletree_old

import (
	"fmt"

	// "encoding/json"
	// "fmt"
	"bytes"
	"math/big"

	"github.com/iden3/go-iden3-core/common"
	common3 "github.com/iden3/go-iden3-core/common"
)

// Hash is the type used to represent a hash used in the MT.
type Hash ElemBytes

func NewHashFromBigInt(e *big.Int) *Hash {
	h := Hash(NewElemBytesFromBigInt(e))
	return &h
}

func (h *Hash) BigInt() *big.Int {
	return (*ElemBytes)(h).BigInt()
}

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

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(common3.HexEncode(h.Bytes())), nil
}

func (h *Hash) UnmarshalText(bs []byte) error {
	return common3.HexDecodeInto(h[:], bs)
}

func (h1 *Hash) Equals(h2 *Hash) bool {
	return bytes.Equal(h1[:], h2[:])
}

func ElemBytesToBigInts(elems ...ElemBytes) []*big.Int {
	ints := make([]*big.Int, len(elems))
	for i, elem := range elems {
		ints[i] = elem.BigInt()
	}
	return ints
}

//func ElemBytesToPoseidonInput(elems ...ElemBytes) ([poseidon.T]*big.Int, error) {
//	bigints := ElemBytesToBigInts(elems...)
//
//	z := big.NewInt(0)
//	b := [poseidon.T]*big.Int{z, z, z, z, z, z}
//	copy(b[:poseidon.T], bigints[:])
//
//	return b, nil
//}
//
//// HashElems performs a poseidon hash over the array of ElemBytes.
//// Uses poseidon.PoseidonHash to be compatible with the circom circuits
//// implementations.
//// The maxim slice input size is poseidon.T
//func HashElems(elems ...ElemBytes) (*Hash, error) {
//	if len(elems) > poseidon.T {
//		return nil, fmt.Errorf("HashElems input can not be bigger than %v", poseidon.T)
//	}
//
//	bi, err := ElemBytesToPoseidonInput(elems...)
//	if err != nil {
//		return nil, err
//	}
//
//	poseidonHash, err := poseidon.PoseidonHash(bi)
//	if err != nil {
//		return nil, err
//	}
//	return NewHashFromBigInt(poseidonHash), nil
//}
//
//// HashElemsKey performs a poseidon hash over the array of ElemBytes.
//func HashElemsKey(key *big.Int, elems ...ElemBytes) (*Hash, error) {
//	if len(elems) > poseidon.T-1 {
//		return nil, fmt.Errorf("HashElemsKey input can not be bigger than %v", poseidon.T-1)
//	}
//	if key == nil {
//		key = new(big.Int).SetInt64(0)
//	}
//	bi, err := ElemBytesToPoseidonInput(elems...)
//	if err != nil {
//		return nil, err
//	}
//	copy(bi[len(elems):], []*big.Int{key})
//	poseidonHash, err := poseidon.PoseidonHash(bi)
//	if err != nil {
//		return nil, err
//	}
//	return NewHashFromBigInt(poseidonHash), nil
//}

// getPath returns the binary path, from the root to the leaf.
func getPath(numLevels int, hIndex *Hash) []bool {
	path := make([]bool, numLevels)
	for n := 0; n < numLevels; n++ {
		path[n] = common.TestBit(hIndex[:], uint(n))
	}
	return path
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
		vBytes, err := common.HexDecode(v)
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
