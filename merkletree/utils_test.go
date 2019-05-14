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
		"0x0b077ec0500876486f9b8860e222cee2a5fc339da0c9b953ce54cb6b7a21c431",
		h.Hex())

	d = IntsToData(1, 0, 0, 0)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x069203c1040df8f57cdbf2670f8c55c7be32fb8553ceb63c02d9a73f4e2cb1a4",
		h.Hex())

	d = IntsToData(0, 0, 0, 1)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x24031c96d5476a281a19e7aa9f0f6efa6ce5662c0a720245a7b75198e6e4e8b1",
		h.Hex())

	d = IntsToData(12, 45, 78, 41)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x284bc1f34f335933a23a433b6ff3ee179d682cd5e5e2fcdd2d964afa85104beb",
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
