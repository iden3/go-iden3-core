package merkletree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	//"strconv"
	"testing"
	//"time"
	"math/big"

	//common3 "github.com/iden3/go-iden3/common"
	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
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
		"1fd4bc970a697084ec1f83ecf81936d4a047e27c654752ddbc89f9ed1728e0ab",
		hex.EncodeToString(e.HIndex()[:]))
}

func TestData(t *testing.T) {
	data := IntsToData(12, 45, 78, 41)
	dataParsed := NewDataFromBytes(data.Bytes())
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
		"0x112bae1c89a7a51a9a09e88c2f095bfe8a7d94d7c0cf5ba017a491c3e0b95c8f",
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
		"0x1fb755a3677f8fd6c47b5462b69778ef6383c31d2d498b765e953f8cacaa6744",
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
		"0x173fd27f6622526dfb21c4d8d83e3c95adba5d8f46a397113e4e80e629c6de76",
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

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0001000000000000000000000000000000000000000000000000000000000001"+
		"05086d2d031b3aeb91b850c7a0280499ded7ba4b8b25caffff5dc754ed207eb8",
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

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"000700000000000000000000000000000000000000000000000000000000005f"+
		"1a412100b0a427796bf11048f3759c9c2163dc39f44d2b243dd8d8fee531800d"+
		"21284e40a8d3e9658d429ebedbef7d509de684b9143aac0b1be3c3979a4e5ed0"+
		"005844e6e7d83169766472da339496a1663ebc14a1016dd39105ca13848f68d4"+
		"102e955ce94001fa069fdd144a3fda637020dd21f1a4a9f23ae0665a4ed27457"+
		"28b82cce3a8d858c295e42c8f058b4a73ca39b0222754165db9f8c7ecf5f431a"+
		"1e220bd9bea43d106d41b7bc2a5d2e88d0fc4f0429c3f8de7059a0bf93e40212",
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

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}

	verify := VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue())
	assert.True(t, verify)
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"000700000000000000000000000000000000000000000000000000000000005f"+
		"1a412100b0a427796bf11048f3759c9c2163dc39f44d2b243dd8d8fee531800d"+
		"21284e40a8d3e9658d429ebedbef7d509de684b9143aac0b1be3c3979a4e5ed0"+
		"005844e6e7d83169766472da339496a1663ebc14a1016dd39105ca13848f68d4"+
		"102e955ce94001fa069fdd144a3fda637020dd21f1a4a9f23ae0665a4ed27457"+
		"28b82cce3a8d858c295e42c8f058b4a73ca39b0222754165db9f8c7ecf5f431a"+
		"1e220bd9bea43d106d41b7bc2a5d2e88d0fc4f0429c3f8de7059a0bf93e40212",
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

	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}

	verify := VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue())
	assert.True(t, verify)
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0105000000000000000000000000000000000000000000000000000000000017"+
		"2e2e61a54ec48cb031effbf00420cd06d707535616965f1ffda2edd1006b807c"+
		"0b9a1a9cc13e5fe12e380fb702c10fde1a9201b7f89e25051f57e00862f20522"+
		"1e1b8f66c3bd26be093e358ed6c54f9d1986411ea404d01482fcf26b04912a0c"+
		"17974735b062a464506127e92858392e853b6abbcb3a1d93a5924af42198c3d1",
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
	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, true)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0003000000000000000000000000000000000000000000000000000000000007"+
		"2e2e61a54ec48cb031effbf00420cd06d707535616965f1ffda2edd1006b807c"+
		"0b9a1a9cc13e5fe12e380fb702c10fde1a9201b7f89e25051f57e00862f20522"+
		"00b1574ea5a96e97ff7b0c964c14d0ad7b9da5b56d068b3aabe10fd3051b0d2a",
		hex.EncodeToString(proof.Bytes()))

	for i := 8; i < 32; i++ {
		e = NewEntryFromInts(0, 0, 0, int64(i))
		proof, err = mt.GenerateProof(e.HIndex(), nil)
		if debug {
			fmt.Println(i, proof)
		}
	}
	// Non-existence proof, empty aux
	e = NewEntryFromInts(0, 0, 0, int64(12))
	proof, err = mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, false)
	assert.True(t, proof.nodeAux == nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"0105000000000000000000000000000000000000000000000000000000000017"+
		"2e2e61a54ec48cb031effbf00420cd06d707535616965f1ffda2edd1006b807c"+
		"0b9a1a9cc13e5fe12e380fb702c10fde1a9201b7f89e25051f57e00862f20522"+
		"1e1b8f66c3bd26be093e358ed6c54f9d1986411ea404d01482fcf26b04912a0c"+
		"17974735b062a464506127e92858392e853b6abbcb3a1d93a5924af42198c3d1",
		hex.EncodeToString(proof.Bytes()))

	// Non-existence proof, diff. node aux
	e = NewEntryFromInts(0, 0, 0, int64(10))
	proof, err = mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, false)
	assert.True(t, proof.nodeAux != nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"030400000000000000000000000000000000000000000000000000000000000b"+
		"0a439bd423b069c01717cd8641d610f286ec5062c9ecee6f2412af76ff551cb5"+
		"1c0fd5c25407d0220a0bcbc6734908153fd18ec43ee62ee157030462c43f537d"+
		"2d5899d3a66630d86c3b9ff896a99b0f2b9e7afc49b69175ed07773ae39263f5"+
		"149648851923be5e707629f0619a1f391452d3c252291d5492d5f9280542380f"+
		"1541a6b5aa9bf7d9be3d5cb0bcc7cacbca26242016a0feebfc19c90f2224baed",
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
	proof, err := mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, true)
	e1 := NewEntryFromInts(0, int64(5), 0, int64(5))
	assert.True(t, !VerifyProof(mt.RootKey(), proof, e1.HIndex(), e1.HValue()))

	// Invalid non-existence proof (Non-existence proof, diff. node aux)
	e = NewEntryFromInts(0, 0, 0, int64(4))
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

	e0 := NewEntryFromInts(0, 0, 0, 0)
	if err := mt.Add(&e0); err != nil {
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
	e2 := NewEntryFromInts(0, 0, 0, int64(1))
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
		e := NewEntryFromInts(0, int64(i), 0, int64(i))
		if err := mt.Add(&e); err != nil {
			t.Fatal(err)
		}
	}

	// Proof of existence, single claim MT
	e0 := NewEntryFromInts(0, 0, 0, 0)
	proof0, err := mt.GenerateProof(e0.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proof0Parsed, err := NewProofFromBytes(proof0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof0, proof0Parsed)

	// Proof of non-existence with empty node, single claim MT
	e1 := NewEntryFromInts(0, 0, 0, int64(17))
	proof1, err := mt.GenerateProof(e1.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	proof1Parsed, err := NewProofFromBytes(proof1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, proof1, proof1Parsed)

	// Proof of non-existence with aux node, single claim MT
	e2 := NewEntryFromInts(0, 0, int64(1), 0)
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

	ethAddr := common.HexToAddress("0x970E8128AB834E8EAC17Ab8E3812F010678CF791")

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := mt.Storage().WithPrefix(ethAddr.Bytes())
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
	if debug {
		fmt.Println(s.String())
	}
}

func TestMTWalk(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
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
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.Add(&e); err != nil {
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
func TestMTWalkDumpLeafs(t *testing.T) {
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

		if err := mt.Add(e); err != nil {
			t.Fatal(err)
		}
	}

	w := bytes.NewBufferString("")
	fmt.Fprintf(w, "--------\nDumpClaims of the MerkleTree with RootKey "+mt.RootKey().Hex()+"\n")
	err := mt.DumpClaims(w, nil)
	fmt.Fprintf(w, "End of DumpClaims of the MerkleTree with RootKey "+mt.RootKey().Hex()+"\n--------\n")
	assert.Nil(t, err)
	if debug {
		fmt.Println(w)
	}
}
