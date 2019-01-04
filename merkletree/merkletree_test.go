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
	"github.com/ethereum/go-ethereum/common"
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
	assert.Equal(t,
		"0x0000000000000000000000000000000000000000000000000000000000000000",
		mt.RootKey().Hex())
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
	assert.Equal(t,
		"114438e8321f62c4a1708f443a5a66f9c8fcb0958e7b7008332b71442610b7a0",
		hex.EncodeToString(e.HIndex()[:]))
}

func TestData(t *testing.T) {
	data := IntsToData(12, 45, 78, 41)
	dataParsed := BytesToData(data.Bytes())
	assert.Equal(t, data, *dataParsed)
}

func TestAddEntry1(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41)
	if err := mt.Add(&e); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0x27c454ae17339dae86b77f0b07a7ff72673201892e281d60394a9b646de29ce3", mt.RootKey().Hex())
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
	assert.Equal(t,
		"0x1475e8f7d486a8d045a04533ad8b27d16ab4850df4e64dc9e39cecb2fcb47cbf",
		mt.RootKey().Hex())
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
	assert.Equal(t,
		"0x1059bfb4f2018d8e15dc5186322b7316d4abda1d534966d0c54e07a4007df51f",
		mt1.RootKey().Hex())
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
	assert.Equal(t,
		"0000000000000000000000000000000000000000000000000000000000000000",
		hex.EncodeToString(proof.Bytes()))
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
	assert.Equal(t, ""+
		"000400000000000000000000000000000000000000000000000000000000000b"+
		"1741ceec35cfc2795e17e4c9ce80992370610dfb25dd01286b33ee5d1a972499"+
		"16ff8f7e5e5ddd7d366eb5758dd44a28823186e85d2d9480d85f45e5e57eba54"+
		"0c719d2afaa5a6769e541b968029fc6ab2ae9c3a4198948f94aebd87a76a3aed",
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
	assert.Equal(t, ""+
		"000400000000000000000000000000000000000000000000000000000000000f"+
		"28df49923aa56a1f3320633c097d56c6f062b5d490698bcca2a84df0c5a7fe87"+
		"0f317f06dfbe10aff5f0a703a7aa09b86011bfc2ad4b465268b52e7dfb0d1ebb"+
		"1e27a375be49f9162136014a9d618b67c6379cf77d77fdc5d6778500e0c93f40"+
		"106e06a9159094113e5d8a94517eb5f468046cd3149b8085391d088b5ab159a4",
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

	verify := VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue())
	assert.True(t, verify)
	assert.Equal(t, ""+
		"000400000000000000000000000000000000000000000000000000000000000f"+
		"28df49923aa56a1f3320633c097d56c6f062b5d490698bcca2a84df0c5a7fe87"+
		"0f317f06dfbe10aff5f0a703a7aa09b86011bfc2ad4b465268b52e7dfb0d1ebb"+
		"1e27a375be49f9162136014a9d618b67c6379cf77d77fdc5d6778500e0c93f40"+
		"106e06a9159094113e5d8a94517eb5f468046cd3149b8085391d088b5ab159a4",
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

	verify := VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue())
	assert.True(t, verify)
	assert.Equal(t, ""+
		"0302000000000000000000000000000000000000000000000000000000000003"+
		"1a97e2325fa70b3ba4958922473b8a55bb24a55d799b583da6f78f89d8d48dea"+
		"3012b3dcbfea0c8d3ecde559c5670ee09a8dfd5b67bbd99dca42d82e4bd06535"+
		"198571b3d34d0989950c7dfd52209ceb5d85400d08137d90cbd96d6223f3a18b"+
		"15331daa10ae035babcaabb76a80198bc449d32240ebb7f456ff2b03cd69bca4",
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
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	assert.Equal(t, ""+
		"000400000000000000000000000000000000000000000000000000000000000f"+
		"1e14ba1e64291bdb0663fd7c4cab8c03115342cdcc6318bcfce8680e5fba816b"+
		"2e3a99b1833362c92ae3c82b8c4b90e7f79cf7335337cb55e3653aa277d72eae"+
		"060a556d978c2b44a12eb35040fc8ce6aa5a4479a34d37f98cb417a59e12f82b"+
		"16ff8f7e5e5ddd7d366eb5758dd44a28823186e85d2d9480d85f45e5e57eba54",
		hex.EncodeToString(proof.Bytes()))

	// Non-existence proof, empty aux
	e = NewEntryFromInts(0, 0, 0, int64(12))
	proof, err = mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.existence, false)
	assert.True(t, proof.nodeAux == nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	assert.Equal(t, ""+
		"010400000000000000000000000000000000000000000000000000000000000b"+
		"1a97e2325fa70b3ba4958922473b8a55bb24a55d799b583da6f78f89d8d48dea"+
		"1741ceec35cfc2795e17e4c9ce80992370610dfb25dd01286b33ee5d1a972499"+
		"25e3a822b71c29996133b4a77b0336e1cb6a07950f6b7765822512640d31e638",
		hex.EncodeToString(proof.Bytes()))

	// Non-existence proof, diff. node aux
	e = NewEntryFromInts(0, 0, 0, int64(10))
	proof, err = mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.existence, false)
	assert.True(t, proof.nodeAux != nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	assert.Equal(t, ""+
		"0302000000000000000000000000000000000000000000000000000000000003"+
		"1a97e2325fa70b3ba4958922473b8a55bb24a55d799b583da6f78f89d8d48dea"+
		"3012b3dcbfea0c8d3ecde559c5670ee09a8dfd5b67bbd99dca42d82e4bd06535"+
		"198571b3d34d0989950c7dfd52209ceb5d85400d08137d90cbd96d6223f3a18b"+
		"15331daa10ae035babcaabb76a80198bc449d32240ebb7f456ff2b03cd69bca4",
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
	assert.True(t, !VerifyProof(mt.RootKey(), proof, e1.HIndex(), e1.HValue()))

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
	proof.nodeAux = &nodeAux{hIndex: e.HIndex(), hValue: e.HValue()}
	assert.True(t, !VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
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

func TestProofFromBytesSmall(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e0 := NewEntryFromInts(0, 0, 0, 0)
	if err := mt.Add(&e0); err != nil {
		t.Fatal(err)
	}

	// Proof of existence, single claim MT
	proof0, err := mt.GenerateProof(e0.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	proof0Parsed, err := NewProofFromBytes(proof0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof0, proof0Parsed)

	// Proof of non-existence with aux node, single claim MT
	e2 := NewEntryFromInts(0, 0, 0, int64(1))
	proof2, err := mt.GenerateProof(e2.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	proof2Parsed, err := NewProofFromBytes(proof2.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof2, proof2Parsed)
}

func TestProofFromBytesBig(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		e := NewEntryFromInts(0, int64(i), 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	// Proof of existence, single claim MT
	e0 := NewEntryFromInts(0, 0, 0, 0)
	proof0, err := mt.GenerateProof(e0.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	proof0Parsed, err := NewProofFromBytes(proof0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof0, proof0Parsed)

	// Proof of non-existence with empty node, single claim MT
	e1 := NewEntryFromInts(0, 0, 0, int64(17))
	proof1, err := mt.GenerateProof(e1.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	proof1Parsed, err := NewProofFromBytes(proof1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof1, proof1Parsed)

	// Proof of non-existence with aux node, single claim MT
	e2 := NewEntryFromInts(0, 0, int64(1), 0)
	proof2, err := mt.GenerateProof(e2.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	proof2Parsed, err := NewProofFromBytes(proof2.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof2, proof2Parsed)
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

func TestDbInsertGet(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	tx, err := mt.storage.NewTx()
	assert.Nil(t, err)
	mt.Lock()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Close()
		}
		mt.Unlock()
	}()

	key := []byte("key")
	mt.dbInsert(tx, key, 9, []byte("value"))
	tx.Commit()

	nodeType, data, err := mt.dbGet(key)
	assert.Nil(t, err)
	assert.Equal(t, byte(9), nodeType)
	assert.Equal(t, []byte("value"), data)

}

func TestMerkleTreeRootStored(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	ethID := common.HexToAddress("0x970E8128AB834E8EAC17Ab8E3812F010678CF791")

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := mt.Storage().WithPrefix(ethID.Bytes())
	// open the MerkleTree of the user
	userMT, err := NewMerkleTree(stoUserID, 140)
	assert.Nil(t, err)

	e := NewEntryFromInts(12, 45, 78, 41)
	err = userMT.Add(&e)
	assert.Nil(t, err)

	// reopen the MerkleTree of the user
	userMTreopened, err := NewMerkleTree(stoUserID, 140)
	assert.Nil(t, err)

	assert.Equal(t, userMT.RootKey(), userMTreopened.RootKey())
}
