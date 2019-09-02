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
		"0x021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		h.Hex())

	d = IntsToData(1, 0, 0, 0)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x050a05b5d53f6f01b1629db59138e94b0827e70cbf91b1f66255b90ca700450d",
		h.Hex())

	d = IntsToData(0, 0, 0, 1)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x0f4fcef5783b61f7bb7424de30dd83476ef2cb9b60bb631cdad8822445ec00d9",
		h.Hex())

	d = IntsToData(12, 45, 78, 41)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x149ce1812eb8bd8bf4b81a81f667faf0c89412705eadcb9153185cbe7ac246f6",
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
