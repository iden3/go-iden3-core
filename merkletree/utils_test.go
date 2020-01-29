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
		"0x026fc3d15c056ad66cf7debba8024ae081dc9bf6fedd6c956a9e9be5abed3404",
		h.Hex())

	d = IntsToData(1, 0, 0, 0, 0, 0, 0, 0)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x8efce10fbbbb8acdac2d811d61e84cf88bd7860cb50c19fc363e0dc88b7c2407",
		h.Hex())

	d = IntsToData(0, 0, 0, 0, 0, 0, 0, 1)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0x8079479cb512aab7567814a36fde824fc38eadbdf0b13508274847b38ec20012",
		h.Hex())

	d = IntsToData(12, 45, 78, 41, 35, 80, 54, 42)
	h = HashElems(d[:]...)
	hashTestOutput(h)
	assert.Equal(t,
		"0xeb775c5100a85d670705f77c240692c8e7de19a0951a38cc1ade6098590f8708",
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
