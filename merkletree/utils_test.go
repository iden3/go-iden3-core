package merkletree

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetBitmap(t *testing.T) {
	var v [32]byte

	setBitBigEndian(v[:], 7)
	setBitBigEndian(v[:], 8)
	setBitBigEndian(v[:], 255)
	expected := "8000000000000000000000000000000000000000000000000000000000000180"
	assert.Equal(t, expected, hex.EncodeToString(v[:]))

	assert.Equal(t, false, testBitBigEndian(v[:], 6))
	assert.Equal(t, true, testBitBigEndian(v[:], 7))
	assert.Equal(t, true, testBitBigEndian(v[:], 8))
	assert.Equal(t, false, testBitBigEndian(v[:], 9))
	assert.Equal(t, true, testBitBigEndian(v[:], 255))

}

func TestHashElems(t *testing.T) {
	d := IntsToData(0, 0, 0, 0)
	h := HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x242b5c83eda327009115428013c28a3eefde6e0283dedee194d80ab70bb71904",
		h.Hex())

	d = IntsToData(1, 0, 0, 0)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x06c81898c60327061330f4bd12d75087a303e181fc9ad3a62e10e4d290f1fec0",
		h.Hex())

	d = IntsToData(0, 0, 0, 1)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x148cd22d6573de244de1bccf18a283c0ce91ddd6edec11b91f934c51cec80bfa",
		h.Hex())

	d = IntsToData(12, 45, 78, 41)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x09a31174e6239fadbe9143f2c1e1b1bbf6aab3ad4df395449bbf74154bb329f9",
		h.Hex())
}

func hashTestOutput(h *Hash) {
	if !debug {
		return
	}
	fmt.Printf("\t\t\"%v\",\n", h.Hex())
}

func BenchmarkHashElems(b *testing.B) {
	ds := make([]Data, b.N)
	for i := 0; i < b.N; i++ {
		ds[i] = IntsToData(0, int64(i), 0, int64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashElems(ds[i][:]...)
	}
}
