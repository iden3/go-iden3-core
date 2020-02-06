package merkletree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	//"strconv"
	"testing"
	//"time"

	//common3 "github.com/iden3/go-iden3-core/common"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/testgen"
	cryptoConstants "github.com/iden3/go-iden3-crypto/constants"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
	"github.com/stretchr/testify/assert"
)

var debug = false

// If generateTest is true, the checked values will be used to generate a test vector
var generateTest = false

type Fatalable interface {
	Fatal(args ...interface{})
}

func interfaceToInt64Array(in interface{}) []int64 {
	var o []int64
	switch t := in.(type) {
	case []interface{}:
		o = make([]int64, len(t))
		for i, v := range t {
			o[i] = int64(v.(float64))
		}
	case []int64:
		o = t
	default:
		panic("Error parsing interface to []int64")
	}
	return o
}
func interfaceToStringArray(in interface{}) []string {
	var o []string
	switch t := in.(type) {
	case []interface{}:
		o = make([]string, len(t))
		for i, v := range t {
			o[i] = v.(string)
		}
	case []string:
		o = t
	default:
		panic("Error parsing interface to []int64")
	}
	return o
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
	testgen.CheckTestValue(t, "TestNewMT0", mt.RootKey().Hex())
}

func TestEntry(t *testing.T) {
	// in := testgen.GetTestValue("EntryInts0").([]int64)
	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts0"))
	e := NewEntryFromIntArray(in)
	testgen.CheckTestValue(t, "TestEntry0", hex.EncodeToString(e.HIndex()[:]))
}

func TestData(t *testing.T) {
	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts0"))
	data := IntArrayToData(in)
	dataParsed := NewDataFromBytes(data.Bytes())
	assert.Equal(t, data, *dataParsed)
}

func TestAddEntry1(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts0"))
	e := NewEntryFromIntArray(in)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	testgen.CheckTestValue(t, "TestAddEntry1", mt.RootKey().Hex())
}

func TestAddEntry2(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts0"))
	e := NewEntryFromIntArray(in)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	in = interfaceToInt64Array(testgen.GetTestValue("EntryInts1"))
	e = NewEntryFromIntArray(in)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	testgen.CheckTestValue(t, "TestAddEntry2", mt.RootKey().Hex())
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
	testgen.CheckTestValue(t, "TestAddEntry16", mt1.RootKey().Hex())
}

func TestAddEntryRepeatIndex(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()
	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts2"))
	e0 := NewEntryFromIntArray(in)
	if err := mt.AddEntry(&e0); err != nil {
		t.Fatal(err)
	}
	in = interfaceToInt64Array(testgen.GetTestValue("EntryInts2"))
	e1 := NewEntryFromIntArray(in)
	err := mt.AddEntry(&e1)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)
}

func TestEntriesIndex(t *testing.T) {
	// Two entries with different Index generate different hash index
	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts4"))
	a := NewEntryFromIntArray(in)
	in = interfaceToInt64Array(testgen.GetTestValue("EntryInts5"))
	b := NewEntryFromIntArray(in)
	assert.NotEqual(t, a.HIndex(), b.HIndex())

	// Two entries with same Index generate the same hash index
	in = interfaceToInt64Array(testgen.GetTestValue("EntryInts6"))
	c := NewEntryFromIntArray(in)
	in = interfaceToInt64Array(testgen.GetTestValue("EntryInts7"))
	d := NewEntryFromIntArray(in)
	assert.Equal(t, c.HIndex(), d.HIndex())
}

func TestGetEntry2(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts0"))
	e := NewEntryFromIntArray(in)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	in = interfaceToInt64Array(testgen.GetTestValue("EntryInts1"))
	e = NewEntryFromIntArray(in)
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

	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts8"))
	e := NewEntryFromIntArray(in)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	testgen.CheckTestValue(t, "TestGenerateProof1", hex.EncodeToString(proof.Bytes()))
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
	testgen.CheckTestValue(t, "TestGenerateProof4", hex.EncodeToString(proof.Bytes()))
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
	testgen.CheckTestValue(t, "TestGenerateProof64", hex.EncodeToString(proof.Bytes()))
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
	testgen.CheckTestValue(t, "TestVerifyProof1", hex.EncodeToString(proof.Bytes()))
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
	testgen.CheckTestValue(t, "TestVerifyProofEmpty", hex.EncodeToString(proof.Bytes()))
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
	testgen.CheckTestValue(t, "TestVerifyProofCases0", hex.EncodeToString(proof.Bytes()))

	for i := 8; i < 32; i++ {
		e := NewEntryFromInts(int64(i), 0, 0, 0, 0, 0, 0, 0)
		proof, err = mt.GenerateProof(e.HIndex(), nil)
		assert.Nil(t, err)
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
	testgen.CheckTestValue(t, "TestVerifyProofCases1", hex.EncodeToString(proof.Bytes()))

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
	testgen.CheckTestValue(t, "TestVerifyProofCases2", hex.EncodeToString(proof.Bytes()))
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
			if err = tx.Commit(); err != nil {
				panic(err)
			}
		} else {
			tx.Close()
		}
		mt.Unlock()
	}()

	key := []byte("key")
	mt.dbInsert(tx, key, 9, []byte("value"))
	if err = tx.Commit(); err != nil {
		panic(err)
	}

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

	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts0"))
	e := NewEntryFromIntArray(in)
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

func newClaimBasicEntry(indexSlot [800 / 8]byte, dataSlot [960 / 8]byte) *Entry {
	e := &Entry{}

	copy(e.Data[0][ElemBytesLen-(64/8):], indexSlot[0:56/8])
	copy(e.Data[1][0:], indexSlot[56/8:304/8])
	copy(e.Data[2][0:], indexSlot[304/8:552/8])
	copy(e.Data[3][0:], indexSlot[552/8:800/8])

	copy(e.Data[4][0:], dataSlot[:216/8])
	copy(e.Data[5][0:], dataSlot[216/8:464/8])
	copy(e.Data[6][0:], dataSlot[464/8:712/8])
	copy(e.Data[7][0:], dataSlot[712/8:960/8])

	return e
}

func TestDumpTreeImportTree(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		rawIndex := strconv.Itoa(i) + testgen.GetTestValue("RawIndex0").(string)
		rawData := testgen.GetTestValue("RawData0").(string)
		var indexSlot [800 / 8]byte
		var dataSlot [960 / 8]byte
		copy(indexSlot[:], rawIndex[:800/8])
		copy(dataSlot[:], rawData[:960/8])
		e := newClaimBasicEntry(indexSlot, dataSlot)

		if err := mt.AddEntry(e); err != nil {
			t.Fatal(err)
		}
	}

	w := bytes.NewBufferString("")
	err := mt.DumpTree(w, nil)
	assert.Nil(t, err)
	if debug {
		fmt.Println(w)
	}

	dumpedTree := w.Bytes()

	imt := newTestingMerkle(t, 140)
	defer imt.Storage().Close()

	err = imt.ImportTree(bytes.NewReader(dumpedTree))
	assert.Nil(t, err)
	assert.Equal(t, mt.RootKey(), imt.RootKey())

	w = bytes.NewBufferString("")
	err = imt.DumpTree(w, nil)
	assert.Nil(t, err)
	dumpedTree2 := w.Bytes()
	assert.Equal(t, dumpedTree, dumpedTree2)
}

func TestMTWalkDumpClaims(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		rawIndex := strconv.Itoa(i) + testgen.GetTestValue("RawIndex0").(string)
		rawData := testgen.GetTestValue("RawData0").(string)
		var indexSlot [800 / 8]byte
		var dataSlot [960 / 8]byte
		copy(indexSlot[:], rawIndex[:800/8])
		copy(dataSlot[:], rawData[:960/8])
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

	dumpedClaims := interfaceToStringArray(testgen.GetTestValue("DumpedClaims"))

	err := mt.ImportDumpedClaims(dumpedClaims)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "TestImportClaims", mt.RootKey().Hex())
}

func TestMTWalkDumpClaimsAndImportDumpedClaims(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		rawIndex := strconv.Itoa(i) + testgen.GetTestValue("RawIndex0").(string)
		rawData := testgen.GetTestValue("RawData0").(string)
		var indexSlot [800 / 8]byte
		var dataSlot [960 / 8]byte
		copy(indexSlot[:], rawIndex[:800/8])
		copy(dataSlot[:], rawData[:960/8])
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

	in := interfaceToInt64Array(testgen.GetTestValue("EntryInts0"))
	e := NewEntryFromInts(in[0], in[1], in[2], in[3], in[4], in[5], in[6], in[7])

	err := mt.AddEntry(&e)
	assert.Nil(t, err)
	err = mt.AddEntry(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)
	err = mt.AddEntry(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)

	testgen.CheckTestValue(t, "TestAddRepeatedClaim", mt.RootKey().Hex())
}

func TestAddBigIntEntries(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	in := interfaceToStringArray(testgen.GetTestValue("EntryFromHexs"))
	e, err := NewEntryFromHexs(in[0], in[1], in[2], in[3], in[4], in[5], in[6], in[7])
	assert.Nil(t, err)

	err = mt.AddEntry(&e)
	assert.Nil(t, err)

	testgen.CheckTestValue(t, "TestAddBigIntEntries", mt.RootKey().Hex())
}

type testClaim struct {
	E *Entry
}

func (tc *testClaim) Entry() *Entry {
	return tc.E
}

func TestEntryToBytesToEntry(t *testing.T) {

	in := interfaceToStringArray(testgen.GetTestValue("EntryFromHexs"))
	e, err := NewEntryFromHexs(in[0], in[1], in[2], in[3], in[4], in[5], in[6], in[7])
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

func initTest() {
	// Init test
	err := testgen.InitTest("merkletree", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("EntryInts0", []int64{12, 45, 78, 41, 35, 80, 54, 42})
		testgen.SetTestValue("EntryInts1", []int64{33, 44, 55, 66, 51, 62, 73, 84})
		testgen.SetTestValue("EntryInts2", []int64{0, 0, 0, 3, 0, 0, 0, 12})
		testgen.SetTestValue("EntryInts3", []int64{0, 0, 0, 3, 0, 0, 0, 45})
		testgen.SetTestValue("EntryInts4", []int64{0, 0, 0, 1, 0, 0, 0, 0})
		testgen.SetTestValue("EntryInts5", []int64{0, 0, 0, 2, 0, 0, 0, 0})
		testgen.SetTestValue("EntryInts6", []int64{0, 0, 0, 3, 0, 0, 0, 1})
		testgen.SetTestValue("EntryInts7", []int64{0, 0, 0, 3, 0, 0, 0, 2})
		testgen.SetTestValue("EntryInts8", []int64{0, 0, 0, 0, 42, 0, 0, 0})
		testgen.SetTestValue("RawIndex0", "-testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest")
		testgen.SetTestValue("RawData0", "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest-")
		testgen.SetTestValue("DumpedClaims", []string{"0x000000000000000000000000000000000000000000000000362d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000322d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x00000000000000000000000000000000000000000000000031322d74657374007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000342d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000352d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x00000000000000000000000000000000000000000000000031332d74657374007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000332d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x00000000000000000000000000000000000000000000000031302d74657374007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x00000000000000000000000000000000000000000000000031312d74657374007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x00000000000000000000000000000000000000000000000031342d74657374007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000372d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000392d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000312d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x00000000000000000000000000000000000000000000000031352d74657374007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000302d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
			"0x000000000000000000000000000000000000000000000000382d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400"})
		testgen.SetTestValue("EntryFromHexs", []string{"0x0000000000000000000000000000000000000000000000000000000000000000",
			"0x0000000000000000000000000000000000000000000000000000000000000000",
			"0x00036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd",
			"0x0000000000000000000000000000000000004d59000000000000000000000004",
			"0x0000000000000000000000000000000000004d59000000000000000000000009",
			"0x0000000000000000000000000000000000004d59000000000000000000000008",
			"0x0000000000000000000000000000000000004d59000000000000000000000007",
			"0x0000000000000000000000000000000000004d59000000000000000000000006"})

		// utils
		testgen.SetTestValue("TestHashElems0", []int64{0, 0, 0, 0, 0, 0, 0, 0})
		testgen.SetTestValue("TestHashElems1", []int64{1, 0, 0, 0, 0, 0, 0, 0})
		testgen.SetTestValue("TestHashElems2", []int64{0, 0, 0, 0, 0, 0, 0, 1})
	}
}

func TestMain(m *testing.M) {
	initTest()
	result := m.Run()
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
	os.Exit(result)
}
