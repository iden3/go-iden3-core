package merkletree

import (
	"encoding/hex"
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
	assert.Equal(t, "0x01bcec76425a7b51ecbbfa7cbf37f4ec55df516ae9f8b28ff012a2e5e61f44b5", h.Hex())

	d = IntsToData(1, 0, 0, 0)
	h = HashElems(d[:]...)
	assert.Equal(t, "0x13c75402700abf74fadf64e67d61576045a43113f3b36fd5557fe1f5ee5402fa", h.Hex())

	d = IntsToData(0, 0, 0, 1)
	h = HashElems(d[:]...)
	assert.Equal(t, "0x1354710c5d868228ca7b854a2ad9a47074761d78581996315067203272a2b5c0", h.Hex())

	d = IntsToData(12, 45, 78, 41)
	h = HashElems(d[:]...)
	assert.Equal(t, "0x1e027004fed670669c5ac756f7cf39cd607299252c241a14d49f478dbd52c3a5", h.Hex())
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
