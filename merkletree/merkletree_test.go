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
	e := NewEntryFromInts(12, 45, 78, 41)
	assert.Equal(t,
		"20555fb8cece3be661b574f3788f7bc44a404c133ade1af360fb2027f7823330",
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
		"0x2788998445b8a4890c87b1ddbfda680910c15d9b303a4027e55af03e27bfc6be",
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
		"0x08aed47b11012036cdac701d95468e4a187181c674f030dc1add043948003f95",
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
		"0x0d44c5a7064fdf0c9417abb5c06bfd56b50ac411c380b6124017a8c7a49c2d84",
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
		"0ce372c6ef3b10b3a0eea195a065c73762012cbbe5fb91767a4e8727566a4373",
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
		"000a00000000000000000000000000000000000000000000000000000000022f"+
		"001f57a4709353008d237db4fec5b6b4d074c44ed2e4d48e22f16f4b44bd9786"+
		"20866629c0f8eb30350d65b21f0379de7ee8a5a9f7c25a902f703a44c52b8582"+
		"0847a9e05f3a5d2825681139ba84518e641a71fd8156b6567d56b30cbfe98cfd"+
		"0d23a4acf74ff511200d53d1f35defb007a8a73bbd433eb3b7d0005869828156"+
		"2daa9c84699538f4aa098dd083348bef1b353a855467306c1ecbffbd19d78511"+
		"011aede300b24621f54e34a0f10bf3cd383e9cd5933060033840c19454eafde2",
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
		"000a00000000000000000000000000000000000000000000000000000000022f"+
		"001f57a4709353008d237db4fec5b6b4d074c44ed2e4d48e22f16f4b44bd9786"+
		"20866629c0f8eb30350d65b21f0379de7ee8a5a9f7c25a902f703a44c52b8582"+
		"0847a9e05f3a5d2825681139ba84518e641a71fd8156b6567d56b30cbfe98cfd"+
		"0d23a4acf74ff511200d53d1f35defb007a8a73bbd433eb3b7d0005869828156"+
		"2daa9c84699538f4aa098dd083348bef1b353a855467306c1ecbffbd19d78511"+
		"011aede300b24621f54e34a0f10bf3cd383e9cd5933060033840c19454eafde2",
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
		"0303000000000000000000000000000000000000000000000000000000000005"+
		"1db2e965cfa2a371abae3fd268618047dd0f839c672e49d3454204a371302682"+
		"2522f2630cf4128a40b79325331c0958852829d761bf0fda11f6c2e84f9cabee"+
		"1d1e2f713d0a058b705da55c59cc7b1600f577ba6ddfccb8b640312d7cd26230"+
		"021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
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
		"000400000000000000000000000000000000000000000000000000000000000f"+
		"2df8ffefef78f5b1971405b93893c12f062aaf9923993138eb9dba2294b5ec5a"+
		"001be301244f384800a0a771a2eaa8353beac02132f3c8cf975b0215bf4c04ff"+
		"1d7241136dad6e2e52d2499a2c92b35716e82b3264732ab54566c6820327a635"+
		"18c5938adf95f76b55448f59c5d6517ea5f56ccb5aeee7ebec1061f32625d147",
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
	// assert.True(t, proof.nodeAux == nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"010500000000000000000000000000000000000000000000000000000000001f"+
		"2df8ffefef78f5b1971405b93893c12f062aaf9923993138eb9dba2294b5ec5a"+
		"001be301244f384800a0a771a2eaa8353beac02132f3c8cf975b0215bf4c04ff"+
		"1d7241136dad6e2e52d2499a2c92b35716e82b3264732ab54566c6820327a635"+
		"2e00575d995e4955576ff334696a442d75b28b70c408d150d083f36f5bb06aad"+
		"216fea80c63e2cd4eb74fb907ff17cfe0439710ed2f3e7b594a443c9bcf6ee00",
		hex.EncodeToString(proof.Bytes()))

	// Non-existence proof, diff. node aux
	e = NewEntryFromInts(0, 0, 0, int64(10))
	proof, err = mt.GenerateProof(e.HIndex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, proof.Existence, false)
	// assert.True(t, proof.nodeAux != nil)
	assert.True(t, VerifyProof(mt.RootKey(), proof, e.HIndex(), e.HValue()))
	proofTestOutput(proof)
	assert.Equal(t, ""+
		"030400000000000000000000000000000000000000000000000000000000000b"+
		"2df8ffefef78f5b1971405b93893c12f062aaf9923993138eb9dba2294b5ec5a"+
		"256fa18d0b93a9e85316b602317b651136d6d545e113d1bf23702948ed8f2742"+
		"27c1d9ed3fd46fb3e4044f680af4f42a749a72964251424fcc3eb3753767ffe7"+
		"17ddf6f66c73719745eeca828537ee30394123a28d16eb51cf51f3bcc0bd03a3"+
		"021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
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

		if err := mt.Add(e); err != nil {
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

	dumpedClaims := []string{"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031342d746573747465737474657374746573747465737474657374746573740074657374746573747465737474657374746573000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500382d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500302d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500332d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500342d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500362d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031322d746573747465737474657374746573747465737474657374746573740074657374746573747465737474657374746573000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031312d746573747465737474657374746573747465737474657374746573740074657374746573747465737474657374746573000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031302d746573747465737474657374746573747465737474657374746573740074657374746573747465737474657374746573000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031352d746573747465737474657374746573747465737474657374746573740074657374746573747465737474657374746573000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500372d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500392d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x007465737474657374746573747465737474657374746573747465737474657300747465737474657374746573747465737474657374746573747465737474650031332d746573747465737474657374746573747465737474657374746573740074657374746573747465737474657374746573000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500312d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500352d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000",
		"0x0074657374746573747465737474657374746573747465737474657374746573007474657374746573747465737474657374746573747465737474657374746500322d74657374746573747465737474657374746573747465737474657374740065737474657374746573747465737474657374000000000000000000000000"}

	err := mt.ImportDumpedClaims(dumpedClaims)
	assert.Nil(t, err)
	assert.Equal(t, "0x2798a8423908359062141846c0acd6d61ef4ac1050195e5cfe2f4e17a544960b", mt.RootKey().Hex())
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

		if err := mt.Add(e); err != nil {
			t.Fatal(err)
		}
	}

	// export claims
	dumpedClaims, err := mt.DumpClaims(nil)
	assert.Nil(t, err)
	assert.Equal(t, 256+2, len(dumpedClaims[0]))

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

	e := NewEntryFromInts(12, 45, 78, 41)

	err := mt.Add(&e)
	assert.Nil(t, err)
	err = mt.Add(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)
	err = mt.Add(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)

	assert.Equal(t,
		"0x2788998445b8a4890c87b1ddbfda680910c15d9b303a4027e55af03e27bfc6be",
		mt.RootKey().Hex())
}

func TestAddBigIntEntries(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e, err := NewEntryFromHexs("0x0000000000000000000000000000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000000000000000000000000000000",
		"0x00036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd",
		"0x0000000000000000000000000000000000004d59000000000000000000000004")
	assert.Nil(t, err)

	err = mt.Add(&e)
	assert.Nil(t, err)

	assert.Equal(t,
		"0x304b899b2bf2cceba0147b2459ccdd56a08d50f5b8c0c700b39042758b941df5",
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
		"0x0000000000000000000000000000000000004d59000000000000000000000004")
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
