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
	d := IntsToData(0, 0, 0, 0, 0, 0, 0, 0)
	h := HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x0434edabe59b9e6a956cddfef69bdc81e04a02a8bbdef76cd66a055cd1c36f02",
		h.Hex())

	d = IntsToData(1, 0, 0, 0, 0, 0, 0, 0)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x07247c8bc80d3e36fc190cb50c86d78bf84ce8611d812daccd8abbbb0fe1fc8e",
		h.Hex())

	d = IntsToData(0, 0, 0, 0, 0, 0, 0, 1)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x1200c28eb34748270835b1f0bdad8ec34f82de6fa3147856b7aa12b59c477980",
		h.Hex())

	d = IntsToData(12, 45, 78, 41, 35, 80, 54, 42)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x08870f599860de1acc381a95a019dee7c89206247cf70507675da800515c77eb",
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
		ds[i] = IntsToData(int64(i), 0, 0, 0, int64(i), 0, 0, 0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashElems(ds[i][:]...)
	}
}
