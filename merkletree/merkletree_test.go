package merkletree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	//"strconv"
	"testing"
	//"time"

	//common3 "github.com/iden3/go-iden3-core/common"
	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/db"
	cryptoConstants "github.com/iden3/go-iden3-crypto/constants"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
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

func TestEntry(t *testing.T) {
	e := NewEntryFromInts(12, 45, 78, 41, 35, 80, 54, 42)
	assert.Equal(t,
		"149ce1812eb8bd8bf4b81a81f667faf0c89412705eadcb9153185cbe7ac246f6",
		hex.EncodeToString(e.HIndex()[:]))
}

func TestData(t *testing.T) {
	data := IntsToData(12, 45, 78, 41, 35, 80, 54, 42)
	dataParsed := NewDataFromBytes(data.Bytes())
	assert.Equal(t, data, *dataParsed)
}

func TestAddEntry1(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41, 35, 80, 54, 42)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t,
		"0x1bfaebcf1bc7b14594c3aa47924e8efebe0307eb38c30b0deedf126735eda373",
		mt.RootKey().Hex())
}

func TestAddEntry2(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41, 35, 80, 54, 42)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	e = NewEntryFromInts(33, 44, 55, 66, 51, 62, 73, 84)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t,
		"0x0105c96bc308358c34992a1cfa405f48b545c5f626d83dab632231149e9392ee",
		mt.RootKey().Hex())
}

func TestAddEntry16(t *testing.T) {
	mt1 := newTestingMerkle(t, 140)
	defer mt1.Storage().Close()
	for i := 0; i < 16; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, int64(i), 0, 0, 0)
		if err := mt1.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	mt2 := newTestingMerkle(t, 140)
	defer mt2.Storage().Close()
	for i := 16 - 1; i >= 0; i-- {
		e := NewEntryFromInts(int64(i), 0, 0, 0, int64(i), 0, 0, 0)
		if err := mt2.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	assert.Equal(t, mt1.RootKey().Hex(), mt2.RootKey().Hex())
	assert.Equal(t,
		"0x21f0ea8226adfd8a00b4041a894bdaf8ed2600ff2935d24c7e31776eb988b5f5",
		mt1.RootKey().Hex())
}

func TestAddEntryRepeatIndex(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()
	e0 := NewEntryFromInts(0, 0, 0, 3, 0, 0, 0, 12)
	if err := mt.AddEntry(&e0); err != nil {
		t.Fatal(err)
	}
	e1 := NewEntryFromInts(0, 0, 0, 3, 0, 0, 0, 45)
	err := mt.AddEntry(&e1)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)
}

func TestEntriesIndex(t *testing.T) {
	// Two entries with different Index generate different hash index
	a := NewEntryFromInts(0, 0, 0, 1, 0, 0, 0, 0)
	b := NewEntryFromInts(0, 0, 0, 2, 0, 0, 0, 0)
	assert.NotEqual(t, a.HIndex(), b.HIndex())

	// Two entries with same Index generate the same hash index
	c := NewEntryFromInts(0, 0, 0, 3, 0, 0, 0, 1)
	d := NewEntryFromInts(0, 0, 0, 3, 0, 0, 0, 2)
	assert.Equal(t, c.HIndex(), d.HIndex())
}

func TestGetEntry2(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41, 35, 80, 54, 42)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	e = NewEntryFromInts(33, 44, 55, 66, 51, 62, 73, 84)
	if err := mt.AddEntry(&e); err != nil {
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

	e := NewEntryFromInts(0, 0, 0, 0, int64(42), 0, 0, 0)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}

	proof, err := mt.GenerateProof(e.HIndex(), nil)
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
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	e := NewEntryFromInts(int64(2), 0, 0, 0, 0, 0, 0, 0)

	data, err := mt.GetDataByIndex(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(NewNodeLeaf(&e).Value()[:], NewNodeLeaf(&Entry{Data: *data}).Value()[:]) {
		t.Fatal(err)
	}

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0002000000000000000000000000000000000000000000000000000000000002"+
		"0c7bb8ca050d0d4d57d1e1be5674f217379627bf2e9abfc52d74a2524c07620d",
		hex.EncodeToString(proof.Bytes()))
}

func TestGenerateProof64(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 64; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	e := NewEntryFromInts(int64(4), 0, 0, 0, 0, 0, 0, 0)

	data, err := mt.GetDataByIndex(e.HIndex())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(NewNodeLeaf(&e).Value()[:], NewNodeLeaf(&Entry{Data: *data}).Value()[:]) {
		t.Fatal(err)
	}

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"000600000000000000000000000000000000000000000000000000000000002f"+
		"2dcbbc24abf11f9747839eac113fb411b6e352514b58165e625e6cc24ac93469"+
		"11933509e8d91ff31366f6849a890087b78bd05fb7b7ed61aa94c961bf09e12f"+
		"2bc176b781327bfb55fd89e6d51d51750c667ad11c55147ccb28951e29eac7c6"+
		"142f8e2800ac660ee471f90ecb7ae0f91e810c6063f34766d31ee2f5fb519a8c"+
		"047685208787b2167dc23713ff34f18de482b227d9322bfb92b46d643602abe3",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProof1(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 64; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	e := NewEntryFromInts(int64(4), 0, 0, 0, 0, 0, 0, 0)

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}

	verify := VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue())
	assert.True(t, verify)
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"000600000000000000000000000000000000000000000000000000000000002f"+
		"2dcbbc24abf11f9747839eac113fb411b6e352514b58165e625e6cc24ac93469"+
		"11933509e8d91ff31366f6849a890087b78bd05fb7b7ed61aa94c961bf09e12f"+
		"2bc176b781327bfb55fd89e6d51d51750c667ad11c55147ccb28951e29eac7c6"+
		"142f8e2800ac660ee471f90ecb7ae0f91e810c6063f34766d31ee2f5fb519a8c"+
		"047685208787b2167dc23713ff34f18de482b227d9322bfb92b46d643602abe3",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofEmpty(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	e := NewEntryFromInts(int64(42), 0, 0, 0, 0, 0, 0, 0)

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}

	verify := VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue())
	assert.True(t, verify)
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0305000000000000000000000000000000000000000000000000000000000017"+
		"068cde0339e040ec5e4a57146280924be1d14809409f4fe4952472ab7e315fdf"+
		"282368154c3ce63e0b000263a7aabdb7830eb7225ddc33b3ec9e661f90c4f458"+
		"057f8f41e819702fe43598a3ad31745c7d69b0beac4e3bb807d16ab578a9cef4"+
		"03f7ebea6b53386e90e3e1f7e1a0c920b631d48c93940f590e5c606e86d4cb7f"+
		"021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781"+
		"021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofCases(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	// Existence proof
	e := NewEntryFromInts(int64(4), 0, 0, 0, 0, 0, 0, 0)
	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, true)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0003000000000000000000000000000000000000000000000000000000000007"+
		"068cde0339e040ec5e4a57146280924be1d14809409f4fe4952472ab7e315fdf"+
		"28c12f4cab2f5c975e032d44cdcdefb54384aaa35e460ff2ab6c832deffb2335"+
		"15865c65c58165511167835c4a26a53c2984e9ac706f8bade9b48cb794dfab12",
		hex.EncodeToString(proof.Bytes()))

	for i := 8; i < 32; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		proof, err = mt.GenerateProof(e.HIndex(), nil)
		if debug {
			fmt.Println(i, proof)
		}
	}
	// Non-existence proof, empty aux
	e = NewEntryFromInts(int64(12), 0, 0, 0, 0, 0, 0, 0)
	proof, err = mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, false)
	// assert.True(t, proof.nodeAux == nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0302000000000000000000000000000000000000000000000000000000000003"+
		"22c56e9d8bd7705b1caa4162b65f6787f571a407772a2ae7381869b2dbbc0451"+
		"173c150a21af227e15075bbf6a567392a493c92f0bf50545d11eeef29e636d1c"+
		"284a3bbf1769b5b8b76afbc8245b4b0fc6f4750895c8859461bf93d87377bb96"+
		"021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		hex.EncodeToString(proof.Bytes()))

	// Non-existence proof, diff. node aux
	e = NewEntryFromInts(int64(10), 0, 0, 0, 0, 0, 0, 0)
	proof, err = mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, false)
	// assert.True(t, proof.nodeAux != nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0303000000000000000000000000000000000000000000000000000000000007"+
		"068cde0339e040ec5e4a57146280924be1d14809409f4fe4952472ab7e315fdf"+
		"28c12f4cab2f5c975e032d44cdcdefb54384aaa35e460ff2ab6c832deffb2335"+
		"15865c65c58165511167835c4a26a53c2984e9ac706f8bade9b48cb794dfab12"+
		"1ab7751a853c6aa24b06f41eb3d108607208a057c1faa7ba2382d5c0245cdf5f"+
		"021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofFalse(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	// Invalid existence proof (node used for verification doesn't
	// correspond to node in the proof)
	e := NewEntryFromInts(int64(4), 0, 0, 0, 0, 0, 0, 0)
	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, true)
	e1 := NewEntryFromInts(int64(5), 0, 0, 0, int64(5), 0, 0, 0)
	assert.True(t, !VerifyProof(mt.RootKey(), proof, e1.HIndex(), e1.HValue()))

	// Invalid non-existence proof (Non-existence proof, diff. node aux)
	e = NewEntryFromInts(int64(4), 0, 0, 0, 0, 0, 0, 0)
	proof, err = mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, true)
	// Now we change the proof from existence to non-existence, and add e's
	// data as auxiliary node.
	proof.Existence = false
	proof.nodeAux = &nodeAux{hIndex: e.HIndex(), hValue: e.HValue()}
	assert.True(t, !VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
}

func TestProofFromBytesSmall(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e0 := NewEntryFromInts(0, 0, 0, 0, 0, 0, 0, 0)
	if err := mt.AddEntry(&e0); err != nil {
		t.Fatal(err)
	}

	// Proof of existence, single claim MT
	proof0, err := mt.GenerateProof(e0.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proof0Parsed, err := NewProofFromBytes(proof0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof0, proof0Parsed)

	// Proof of non-existence with aux node, single claim MT
	e2 := NewEntryFromInts(int64(1), 0, 0, 0, 0, 0, 0, 0)
	proof2, err := mt.GenerateProof(e2.HIndex(), nil)
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
		e := NewEntryFromInts(int64(i), 0, 0, 0, int64(i), 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	// Proof of existence, single claim MT
	e0 := NewEntryFromInts(0, 0, 0, 0, 0, 0, 0, 0)
	proof0, err := mt.GenerateProof(e0.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proof0Parsed, err := NewProofFromBytes(proof0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof0, proof0Parsed)

	// Proof of non-existence with empty node, single claim MT
	e1 := NewEntryFromInts(int64(17), 0, 0, 0, 0, 0, 0, 0)
	proof1, err := mt.GenerateProof(e1.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proof1Parsed, err := NewProofFromBytes(proof1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof1, proof1Parsed)

	// Proof of non-existence with aux node, single claim MT
	e2 := NewEntryFromInts(0, int64(17), 0, 0, 0, 0, 0, 0)
	proof2, err := mt.GenerateProof(e2.HIndex(), nil)
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
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
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

	ethAddr := common.HexToAddress("0x970E8128AB834E8EAC17Ab8E3812F010678CF791")

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := mt.Storage().WithPrefix(ethAddr.Bytes())
	// open the MerkleTree of the user
	userMT, err := NewMerkleTree(stoUserID, 140)
	assert.Nil(t, err)

	e := NewEntryFromInts(12, 45, 78, 41, 35, 80, 54, 42)
	err = userMT.AddEntry(&e)
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
	if debug {
		fmt.Println(s.String())
	}
}

func TestMTWalk(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	w := bytes.NewBufferString("")
	err := mt.Walk(mt.RootKey(), func(n *Node) {
		if n.Type != NodeTypeEmpty {
			fmt.Fprintf(w, "node \"%v\"\n", common3.HexEncode(n.Value()))
		}
	})
	assert.Nil(t, err)
	if debug {
		fmt.Println(w)
	}
}

func TestMTWalkGraphViz(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		if err := mt.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	w := bytes.NewBufferString("")
	fmt.Fprintf(w, "--------\nGraphViz of the MerkleTree with RootKey "+mt.RootKey().Hex()+"\n")
	err := mt.GraphViz(w, nil)
	fmt.Fprintf(w, "End of GraphViz of the MerkleTree with RootKey "+mt.RootKey().Hex()+"\n--------\n")
	assert.Nil(t, err)
	if debug {
		fmt.Println(w)
	}
}

func copyToElemBytes(e *ElemBytes, start int, src []byte) {
	copy(e[ElemBytesLen-start-len(src):], src)
}
func newClaimBasicEntry(indexSlot [400 / 8]byte, dataSlot [496 / 8]byte) *Entry {
	e := &Entry{}
	claimTypeVersionLen := (64 / 8) + (32 / 8)
	copyToElemBytes(&e.Data[3], claimTypeVersionLen, indexSlot[len(indexSlot)-152/8:])
	copyToElemBytes(&e.Data[2], 0, indexSlot[:248/8])
	copyToElemBytes(&e.Data[1], 0, dataSlot[248/8:])
	copyToElemBytes(&e.Data[0], 0, dataSlot[:248/8])
	return e
}

func TestMTWalkDumpClaims(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		rawIndex := strconv.Itoa(i) + "-testtesttesttesttesttesttesttesttesttesttesttest"
		rawData := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttest-"
		var indexSlot [400 / 8]byte
		var dataSlot [496 / 8]byte
		copy(indexSlot[:], rawIndex[:400/8])
		copy(dataSlot[:], rawData[:496/8])
		e := newClaimBasicEntry(indexSlot, dataSlot)

		if err := mt.AddEntry(e); err != nil {
			t.Fatal(err)
		}
	}

	w := bytes.NewBufferString("")
	fmt.Fprintf(w, "--------\nDumpClaims of the MerkleTree with RootKey "+mt.RootKey().Hex()+"\n")
	err := mt.DumpClaimsIoWriter(w, nil)
	fmt.Fprintf(w, "End of DumpClaims of the MerkleTree with RootKey "+mt.RootKey().Hex()+"\n--------\n")
	assert.Nil(t, err)
	if debug {
		fmt.Println(w)
	}
}

func TestImportClaims(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	dumpedClaims := []string{"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031322d7465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465730000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500372d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500382d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031342d7465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465730000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031332d7465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465730000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500302d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031312d7465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465730000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031352d7465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465730000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031302d7465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465730000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500332d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500392d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500352d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500312d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500322d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500342d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500362d746573747465737474657374746573747465737474657374746573747400657374746573747465737474657374746573740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"}

	err := mt.ImportDumpedClaims(dumpedClaims)
	assert.Nil(t, err)
	assert.Equal(t, "0x112090fc177384d3d3fda930d2a2b31a047129b4592fdae3a723fbada86f2041", mt.RootKey().Hex())
}

func TestMTWalkDumpClaimsAndImportDumpedClaims(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		rawIndex := strconv.Itoa(i) + "-testtesttesttesttesttesttesttesttesttesttesttest"
		rawData := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttest-"
		var indexSlot [400 / 8]byte
		var dataSlot [496 / 8]byte
		copy(indexSlot[:], rawIndex[:400/8])
		copy(dataSlot[:], rawData[:496/8])
		e := newClaimBasicEntry(indexSlot, dataSlot)

		if err := mt.AddEntry(e); err != nil {
			t.Fatal(err)
		}
	}

	// export claims
	dumpedClaims, err := mt.DumpClaims(nil)
	assert.Nil(t, err)
	assert.Equal(t, 64*DataLen+2, len(dumpedClaims[0]))

	// import claims
	mt2 := newTestingMerkle(t, 140)
	defer mt2.Storage().Close()
	err = mt2.ImportDumpedClaims(dumpedClaims)
	assert.Nil(t, err)

	assert.Equal(t, mt.RootKey().Hex(), mt2.RootKey().Hex())
}

func TestAddRepeatedClaim(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41, 35, 80, 54, 42)

	err := mt.AddEntry(&e)
	assert.Nil(t, err)
	err = mt.AddEntry(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)
	err = mt.AddEntry(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)

	assert.Equal(t,
		"0x1bfaebcf1bc7b14594c3aa47924e8efebe0307eb38c30b0deedf126735eda373",
		mt.RootKey().Hex())
}

func TestAddBigIntEntries(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e, err := NewEntryFromHexs("0x0000000000000000000000000000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000000000000000000000000000000",
		"0x00036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd",
		"0x0000000000000000000000000000000000004d59000000000000000000000004",
		"0x0000000000000000000000000000000000004d59000000000000000000000009",
		"0x0000000000000000000000000000000000004d59000000000000000000000008",
		"0x0000000000000000000000000000000000004d59000000000000000000000007",
		"0x0000000000000000000000000000000000004d59000000000000000000000006")
	assert.Nil(t, err)

	err = mt.AddEntry(&e)
	assert.Nil(t, err)

	assert.Equal(t,
		"0x2de6e91744026bda8a865c71c7a35fd0cad1f58ed0ac03c89659110d988e61ca",
		mt.RootKey().Hex())
}

type testClaim struct {
	E *Entry
}

func (tc *testClaim) Entry() *Entry {
	return tc.E
}

func TestEntryToBytesToEntry(t *testing.T) {

	e, err := NewEntryFromHexs("0x0000000000000000000000000000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000000000000000000000000000000",
		"0x00036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd",
		"0x0000000000000000000000000000000000004d59000000000000000000000004",
		"0x0000000000000000000000000000000000004d59000000000000000000000009",
		"0x0000000000000000000000000000000000004d59000000000000000000000008",
		"0x0000000000000000000000000000000000004d59000000000000000000000007",
		"0x0000000000000000000000000000000000004d59000000000000000000000006")
	assert.Nil(t, err)

	claim := testClaim{
		E: &e,
	}
	cBytes := claim.Entry().Bytes()

	var leafBytes [ElemBytesLen * DataLen]byte
	copy(leafBytes[:], cBytes[:ElemBytesLen*DataLen])
	leafData := NewDataFromBytes(leafBytes)
	leafDataBytes := leafData.Bytes()

	assert.Equal(t, cBytes, leafDataBytes[:])
	assert.Equal(t, cBytes, leafBytes[:])

	entry := Entry{
		Data: *leafData,
	}
	for _, elemBytes := range entry.Data {
		bigints := ElemBytesToBigInt(elemBytes)
		ok := cryptoUtils.CheckBigIntInField(bigints, cryptoConstants.Q)
		assert.True(t, ok)
	}
}
