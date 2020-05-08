package merkletree

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

func TestGetSetBitmap(t *testing.T) {
	var v [32]byte

	common.SetBitBigEndian(v[:], 7)
	common.SetBitBigEndian(v[:], 8)
	common.SetBitBigEndian(v[:], 255)
	testgen.CheckTestValue(t, "TestGetSetBitmap", hex.EncodeToString(v[:]))

	assert.Equal(t, false, common.TestBitBigEndian(v[:], 6))
	assert.Equal(t, true, common.TestBitBigEndian(v[:], 7))
	assert.Equal(t, true, common.TestBitBigEndian(v[:], 8))
	assert.Equal(t, false, common.TestBitBigEndian(v[:], 9))
	assert.Equal(t, true, common.TestBitBigEndian(v[:], 255))

}

func TestHashElems(t *testing.T) {
	in := interfaceToInt64Array(testgen.GetTestValue("TestHashElems0"))
	d := IntsToData(in[0], in[1], in[2], in[3], in[4], in[5], in[6], in[7])
	h, err := HashElems(d[:4]...)
	assert.Nil(t, err)
	hashTestOutput(h)
	testgen.CheckTestValue(t, "TestHashElems0a", h.Hex())
	h, err = HashElems(d[4:]...)
	assert.Nil(t, err)
	hashTestOutput(h)
	testgen.CheckTestValue(t, "TestHashElems0b", h.Hex())

	in = interfaceToInt64Array(testgen.GetTestValue("TestHashElems1"))
	d = IntsToData(in[0], in[1], in[2], in[3], in[4], in[5], in[6], in[7])
	h, err = HashElems(d[:4]...)
	assert.Nil(t, err)
	hashTestOutput(h)
	testgen.CheckTestValue(t, "TestHashElems1a", h.Hex())
	h, err = HashElems(d[4:]...)
	assert.Nil(t, err)
	hashTestOutput(h)
	testgen.CheckTestValue(t, "TestHashElems1b", h.Hex())

	in = interfaceToInt64Array(testgen.GetTestValue("TestHashElems2"))
	d = IntsToData(in[0], in[1], in[2], in[3], in[4], in[5], in[6], in[7])
	h, err = HashElems(d[:4]...)
	assert.Nil(t, err)
	hashTestOutput(h)
	testgen.CheckTestValue(t, "TestHashElems2a", h.Hex())
	h, err = HashElems(d[4:]...)
	assert.Nil(t, err)
	hashTestOutput(h)
	testgen.CheckTestValue(t, "TestHashElems2b", h.Hex())

	in = interfaceToInt64Array(testgen.GetTestValue("EntryInts0"))
	d = IntsToData(in[0], in[1], in[2], in[3], in[4], in[5], in[6], in[7])
	h, err = HashElems(d[:4]...)
	assert.Nil(t, err)
	hashTestOutput(h)
	testgen.CheckTestValue(t, "TestHashElems3a", h.Hex())
	h, err = HashElems(d[4:]...)
	assert.Nil(t, err)
	hashTestOutput(h)
	testgen.CheckTestValue(t, "TestHashElems3b", h.Hex())
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
		HashElems(ds[i][:]...) //nolint:errcheck
	}
}
