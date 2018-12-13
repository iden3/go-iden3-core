package merkletree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	//"strconv"
	"testing"
	//"time"
	"math/big"

	//common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/db"
	"github.com/stretchr/testify/assert"
)

type Fatalable interface {
	Fatal(args ...interface{})
}

func newTestingMerkle(f Fatalable, numLevels int) *MerkleTree {
	mt, err := NewMerkleTree(db.NewMemoryStorage(), numLevels)
	if err != nil {
		f.Fatal(err)
		return nil
	}
	return mt
}

func TestNewMT(t *testing.T) {
	//create a new MT
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", mt.RootKey().Hex())
}

func NewEntryFromInts(a, b, c, d int64) (e Entry) {
	e.Data = IntsToData(a, b, c, d)
	return e
}

func IntsToData(_a, _b, _c, _d int64) Data {
	a, b, c, d := big.NewInt(_a), big.NewInt(_b), big.NewInt(_c), big.NewInt(_d)
	return BigIntsToData(a, b, c, d)
}

func BigIntsToData(a, b, c, d *big.Int) (data Data) {
	di := []*big.Int{a, b, c, d}
	for i, v := range di {
		copy(data[i][(ElemBytesLen-len(v.Bytes())):], v.Bytes())
	}
	return
}

func TestEntry(t *testing.T) {
	e := NewEntryFromInts(12, 45, 78, 41)
	assert.Equal(t, "114438e8321f62c4a1708f443a5a66f9c8fcb0958e7b7008332b71442610b7a0", hex.EncodeToString(e.HIndex()[:]))
}

func TestAddEntry1(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41)
	if err := mt.Add(&e); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0x2d49fc39bb8f19f26ad47f63b45f77eb4ca50e6548244140a63a105d7c4535d2", mt.RootKey().Hex())
}

func TestAddEntry2(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41)
	if err := mt.Add(&e); err != nil {
		t.Fatal(err)
	}
	e = NewEntryFromInts(33, 44, 55, 66)
	if err := mt.Add(&e); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0x24dcdcb8b10bed49ed2c7795972f2bea478750fc9940eeb64f42440fe0db7cbe", mt.RootKey().Hex())
}

func TestAddEntry16(t *testing.T) {
	mt1 := newTestingMerkle(t, 140)
	defer mt1.Storage().Close()
	for i := 0; i < 16; i++ {
		e := NewEntryFromInts(0, int64(i), 0, int64(i))
		if err := mt1.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	mt2 := newTestingMerkle(t, 140)
	defer mt2.Storage().Close()
	for i := 16 - 1; i >= 0; i-- {
		e := NewEntryFromInts(0, int64(i), 0, int64(i))
		if err := mt2.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	assert.Equal(t, mt1.RootKey().Hex(), mt2.RootKey().Hex())
	assert.Equal(t, "0x0bd26ed069568d6db1032f2761b56167d8b618204c5c1b0dd54bb4a4010fe36e", mt1.RootKey().Hex())
}

func TestEntriesIndex(t *testing.T) {
	// Two entries with different Index generate different hash index
	a := NewEntryFromInts(0, 0, 0, 1)
	b := NewEntryFromInts(0, 0, 0, 2)
	assert.NotEqual(t, a.HIndex(), b.HIndex())

	// Two entries with same Index generate the same hash index
	c := NewEntryFromInts(0, 1, 0, 3)
	d := NewEntryFromInts(0, 2, 0, 3)
	assert.Equal(t, c.HIndex(), d.HIndex())
}

func TestGetEntry2(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41)
	if err := mt.Add(&e); err != nil {
		t.Fatal(err)
	}
	e = NewEntryFromInts(33, 44, 55, 66)
	if err := mt.Add(&e); err != nil {
		t.Fatal(err)
	}

	data, err := mt.GetDataByIndex(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(NewNodeLeaf(&e).Value()[:], NewNodeLeaf(&Entry{Data: *data}).Value()[:]) {
		t.Fatal(err)
	}
}

func TestGenerateProof1(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(0, int64(42), 0, 0)
	if err := mt.Add(&e); err != nil {
		t.Fatal(err)
	}

	proof, err := mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", hex.EncodeToString(proof.Bytes()))
}

func TestGenerateProof4(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 4; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	e := NewEntryFromInts(0, 0, 0, int64(2))

	data, err := mt.GetDataByIndex(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(NewNodeLeaf(&e).Value()[:], NewNodeLeaf(&Entry{Data: *data}).Value()[:]) {
		t.Fatal(err)
	}

	proof, err := mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(proof)
	assert.Equal(t,
		"000400000000000000000000000000000000000000000000000000000000000b293a5e97fccfafe457fa22796168cfce0ff8928fbbd7da9d2cf983287ec52f3317f0c4fe7ebb238a42891bce3d7d2cdf288f1a0237f97530611a591c6deae08e17f267633bb0021e42eac5a0662921709310747225d4fcae6b6c63187b0e7a62",
		hex.EncodeToString(proof.Bytes()))
}

func TestGenerateProof64(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 64; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	e := NewEntryFromInts(0, 0, 0, int64(4))

	data, err := mt.GetDataByIndex(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(NewNodeLeaf(&e).Value()[:], NewNodeLeaf(&Entry{Data: *data}).Value()[:]) {
		t.Fatal(err)
	}

	proof, err := mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(proof)
	assert.Equal(t,
		"000400000000000000000000000000000000000000000000000000000000000f29b583ac7aa28489977b8383d310b47282010a30d3ef76d31a462845cec334a304574e4a467d11f2c53c3653548ea4c6b4884194c1efb32d96d0f6ef5d70c420301af026598b737db5fad61d12769c1c350e9b395dffc1c42b46ebf888c31bd61b6514b48b7da109066e8a0952ae47f4d7b031a2f690bdcde9e2e84746bc036c",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProof1(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 64; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	e := NewEntryFromInts(0, 0, 0, int64(4))

	proof, err := mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println()
	fmt.Println(proof)
	fmt.Println()

	verify := VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HTotal())
	assert.True(t, verify)
	assert.Equal(t,
		"000400000000000000000000000000000000000000000000000000000000000f29b583ac7aa28489977b8383d310b47282010a30d3ef76d31a462845cec334a304574e4a467d11f2c53c3653548ea4c6b4884194c1efb32d96d0f6ef5d70c420301af026598b737db5fad61d12769c1c350e9b395dffc1c42b46ebf888c31bd61b6514b48b7da109066e8a0952ae47f4d7b031a2f690bdcde9e2e84746bc036c",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofEmpty(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	e := NewEntryFromInts(0, 0, 0, int64(42))

	proof, err := mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println()
	fmt.Println(proof)
	fmt.Println()

	verify := VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HTotal())
	assert.True(t, verify)
	assert.Equal(t,
		"03020000000000000000000000000000000000000000000000000000000000032457c8e7eabebeeef71726e920f7c8b63da2f6b3cd97743ea8fb49eae76e46641ba6d011509e611076162c1f94e6e099a0e9fc0f992282f881324213dd4e3e40198571b3d34d0989950c7dfd52209ceb5d85400d08137d90cbd96d6223f3a18b2d957252161c7a359052be895ef1bf56228e3a2977ad906d1d4b25dfb8aadb1c",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofCases(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	// Existence proof
	e := NewEntryFromInts(0, 0, 0, int64(4))
	proof, err := mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.existence, true)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HTotal()))
	assert.Equal(t,
		"000400000000000000000000000000000000000000000000000000000000000f2b724aa8a314c8da446586a0636329c4815794b913e4dfa4a15bdf58ef34b507209978f585bd0ac41e12c9089441a52a68baf5b08959f9c68a89d72eb630c48b2d36441b75d605e210812607b31c35be8210e011fb7faf830cf74cd13cb3686f17f0c4fe7ebb238a42891bce3d7d2cdf288f1a0237f97530611a591c6deae08e",
		hex.EncodeToString(proof.Bytes()))

	// Non-existence proof, empty aux
	e = NewEntryFromInts(0, 0, 0, int64(12))
	proof, err = mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.existence, false)
	assert.True(t, proof.nodeAux == nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HTotal()))
	assert.Equal(t,
		"010400000000000000000000000000000000000000000000000000000000000b2457c8e7eabebeeef71726e920f7c8b63da2f6b3cd97743ea8fb49eae76e4664293a5e97fccfafe457fa22796168cfce0ff8928fbbd7da9d2cf983287ec52f332fe2cda15e178196dc20854bbe646533657523f56a92930aaed3eb2dc88369ff",
		hex.EncodeToString(proof.Bytes()))

	// Non-existence proof, diff. node aux
	e = NewEntryFromInts(0, 0, 0, int64(10))
	proof, err = mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.existence, false)
	assert.True(t, proof.nodeAux != nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HTotal()))
	assert.Equal(t,
		"03020000000000000000000000000000000000000000000000000000000000032457c8e7eabebeeef71726e920f7c8b63da2f6b3cd97743ea8fb49eae76e46641ba6d011509e611076162c1f94e6e099a0e9fc0f992282f881324213dd4e3e40198571b3d34d0989950c7dfd52209ceb5d85400d08137d90cbd96d6223f3a18b2d957252161c7a359052be895ef1bf56228e3a2977ad906d1d4b25dfb8aadb1c",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofFalse(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	// Invalid existence proof (node used for verification doesn't
	// correspond to node in the proof)
	e := NewEntryFromInts(0, 0, 0, int64(4))
	proof, err := mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.existence, true)
	e1 := NewEntryFromInts(0, int64(5), 0, int64(5))
	assert.True(t, !VerifyProof(mt.RootKey(), proof, e1.HIndex(), e1.HTotal()))

	// Invalid non-existence proof (Non-existence proof, diff. node aux)
	e = NewEntryFromInts(0, 0, 0, int64(4))
	proof, err = mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.existence, true)
	// Now we change the proof from existence to non-existence, and add e's
	// data as auxiliary node.
	proof.existence = false
	proof.nodeAux = &nodeAux{hIndex: e.HIndex(), hTotal: e.HTotal()}
	assert.True(t, !VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HTotal()))
}

func TestMTGraphViz(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	s := bytes.NewBufferString("")
	mt.GraphViz(s)
	fmt.Println(s)
}

func BenchmarkAddEntry(b *testing.B) {
	mt := newTestingMerkle(b, 140)
	defer mt.Storage().Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
			b.Fatal(err)
		}
	}
}
