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
		"0x29a6a240e2d8f8bf39b5338b9664d414c5d793f4dead5414b55645f2087e2acb",
		h.Hex())

	d = IntsToData(1, 0, 0, 0)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x118163a0cff49aa500bbf3c6d1bc9ab0de41f3d95d35de7c10d1ad3dbf6af2e9",
		h.Hex())

	d = IntsToData(0, 0, 0, 1)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x0ca28f247163f514b1c4d1db84e9b06d159b54026cf8e4f22a383b44f474b070",
		h.Hex())

	d = IntsToData(12, 45, 78, 41)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x10e02cc6c8fc40cda121602903df911f6398d65f84ff1f27c680d0b7d85b7418",
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
