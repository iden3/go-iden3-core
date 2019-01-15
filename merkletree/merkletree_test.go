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

var debug = false

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
		"0a546e62170e636895398e3831a7834affd35adad67016ff2d5215828e643aaf",
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
	assert.Equal(t,
		"0x2e4bd3b61579f9ed4d814dfa4228c743b853fb6d0b5d6a80735bd5aab579231f",
		mt.RootKey().Hex())
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
		"0x2c0efd3fd64110611ee03714790060b506fc52229146df05474f89b50b8018fa",
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
		"0x171726a6de9a2e6112f467543ce72f67871e7804c070355ff98f6ba63631c25f",
		mt1.RootKey().Hex())
}

func TestAddEntryRepeatIndex(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()
	e0 := NewEntryFromInts(0, 12, 0, 3)
	if err := mt.Add(&e0); err != nil {
		t.Fatal(err)
	}
	e1 := NewEntryFromInts(0, 45, 0, 3)
	err := mt.Add(&e1)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)
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
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"000400000000000000000000000000000000000000000000000000000000000d"+
		"2f5e46500bd35bd62bcd14a6587f33003270e26a74a5a0f606e6380d209c338c"+
		"0a108d874a3f027a9f3e8e2b6f3dad0faa5304fd3c5c7a70c6289b7fbe44dd43"+
		"04e61a72b69b24e24458f33ed6c07fcb8126782b47fba27b9123182a9b4e3ff2",
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
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"000700000000000000000000000000000000000000000000000000000000005f"+
		"292f31c3effbd4ef90051e62e453ce5ba9253a790c33e99ee7989d3778b841c9"+
		"176f8d38c4bda0ced90edd02a1131d9f328ed83659dff2a6a5b7072a53484944"+
		"01eae6da67bc0a8033b7f7ea961a9b88a14477f8f7a90fa0ac7fcb10a7f06626"+
		"0fe17cbefb6ad59ca52b64788a80f77636bb0a03b5ee8ce41bdd327323ded8f5"+
		"2f26ab8d6d6a3eae1d7e9b3522b55ede77619cbc4e157c23f63522a6298b88b0"+
		"260716f66b66072ba0c6b418995af4ba604953f91d376c76d5ca2a655f4a073b",
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
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"000700000000000000000000000000000000000000000000000000000000005f"+
		"292f31c3effbd4ef90051e62e453ce5ba9253a790c33e99ee7989d3778b841c9"+
		"176f8d38c4bda0ced90edd02a1131d9f328ed83659dff2a6a5b7072a53484944"+
		"01eae6da67bc0a8033b7f7ea961a9b88a14477f8f7a90fa0ac7fcb10a7f06626"+
		"0fe17cbefb6ad59ca52b64788a80f77636bb0a03b5ee8ce41bdd327323ded8f5"+
		"2f26ab8d6d6a3eae1d7e9b3522b55ede77619cbc4e157c23f63522a6298b88b0"+
		"260716f66b66072ba0c6b418995af4ba604953f91d376c76d5ca2a655f4a073b",
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
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0302000000000000000000000000000000000000000000000000000000000003"+
		"0397d8b2557e1fd4d6ca48508c0a976bfea8a5b81e99e0f08c1dfddd93ddd16b"+
		"2b0905533da6edc91cd889feca537ffb36ec12546bebd31554d090d9a751a966"+
		"21299cf503028e491b03603fbff7714e25aad3195eb76fe44f7d03638f5a3d0a"+
		"13580fd5d3ca0f7604a3a50f663cb4fd23c214f1955fa5b3ee9ed5ed06bb70a3",
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
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0002000000000000000000000000000000000000000000000000000000000003"+
		"06704e71b5ba8daa15428ae04764b14c403211e7a4dd57dca9e1ff3cd7b7f0e7"+
		"128c1ba622cda99a3144e7178a6bdc59ad8d180b642aee16622613ab65c97419",
		hex.EncodeToString(proof.Bytes()))

	//for i := 8; i < 32; i++ {
	//	e = NewEntryFromInts(0, 0, 0, int64(i))
	//	proof, err = mt.GenerateProof(e.HIndex())
	//	fmt.Println(i, proof)
	//}
	// Non-existence proof, empty aux
	e = NewEntryFromInts(0, 0, 0, int64(11))
	proof, err = mt.GenerateProof(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.existence, false)
	assert.True(t, proof.nodeAux == nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0105000000000000000000000000000000000000000000000000000000000017"+
		"06704e71b5ba8daa15428ae04764b14c403211e7a4dd57dca9e1ff3cd7b7f0e7"+
		"0ca6fc5e3a7efbcd8253564af04933f8a5a4fcdc0416b2a7aa951a9a92a6b15a"+
		"0200592d33c08050ca689c511bfba5366727f9a440c444f248837cf5f85b0a74"+
		"22d492ff12960c3b3214afd2a18f762b04651c2d35ce2a4edb5ab5b26b23ce8c",
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
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0302000000000000000000000000000000000000000000000000000000000003"+
		"06704e71b5ba8daa15428ae04764b14c403211e7a4dd57dca9e1ff3cd7b7f0e7"+
		"128c1ba622cda99a3144e7178a6bdc59ad8d180b642aee16622613ab65c97419"+
		"01b64831f1504ebb67612647c67a30640023df5b51236ab409cd85923b0206d5"+
		"13580fd5d3ca0f7604a3a50f663cb4fd23c214f1955fa5b3ee9ed5ed06bb70a3",
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
	if debug {
		fmt.Println(s)
	}
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
	assert.Equal(t, NodeType(9), nodeType)
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

func proofTestOutput(p *Proof) {
	if !debug {
		return
	}
	s := bytes.NewBufferString("")
	pHex := hex.EncodeToString(p.Bytes())
	chunks := len(pHex) / 64
	for i := 0; i < chunks; i++ {
		fmt.Fprintf(s, "\t\t\"%v\"", pHex[i*(64):(i+1)*(64)])
		if i == chunks-1 {
			fmt.Fprintf(s, ",\n")
		} else {
			fmt.Fprintf(s, "+\n")

		}
	}
	fmt.Println(s.String())
}
