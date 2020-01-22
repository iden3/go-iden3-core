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
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t,
		"0x2ab68ca530d96aa00eafb1b1f40aa2ed0c8ce6e7aec8bed48b7efd52e919f91b",
		mt.RootKey().Hex())
}

func TestAddEntry2(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	e := NewEntryFromInts(12, 45, 78, 41)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	e = NewEntryFromInts(33, 44, 55, 66)
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t,
		"0x0f42da5cbc007dcf89c72ba7212e078c00aba09e9b8d6114515b56daacca95dd",
		mt.RootKey().Hex())
}

func TestAddEntry16(t *testing.T) {
	mt1 := newTestingMerkle(t, 140)
	defer mt1.Storage().Close()
	for i := 0; i < 16; i++ {
		e := NewEntryFromInts(0, int64(i), 0, int64(i))
		if err := mt1.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	mt2 := newTestingMerkle(t, 140)
	defer mt2.Storage().Close()
	for i := 16 - 1; i >= 0; i-- {
		e := NewEntryFromInts(0, int64(i), 0, int64(i))
		if err := mt2.AddEntry(&e); err != nil {
			t.Fatal(err)
		}
	}

	assert.Equal(t, mt1.RootKey().Hex(), mt2.RootKey().Hex())
	assert.Equal(t,
		"0x27de12d35c012988b4fcc275829217c79ae1a2322634ebead16afccd95ffd6fc",
		mt1.RootKey().Hex())
}

func TestAddEntryRepeatIndex(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()
	e0 := NewEntryFromInts(0, 12, 0, 3)
	if err := mt.AddEntry(&e0); err != nil {
		t.Fatal(err)
	}
	e1 := NewEntryFromInts(0, 45, 0, 3)
	err := mt.AddEntry(&e1)
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
	if err := mt.AddEntry(&e); err != nil {
		t.Fatal(err)
	}
	e = NewEntryFromInts(33, 44, 55, 66)
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

	e := NewEntryFromInts(0, int64(42), 0, 0)
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
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.AddEntry(&e); err != nil {
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
		"255b46f58791e3d5220cc9fecfaf2750a02fb8f6135b4561af541dc0ebfe080d",
		hex.EncodeToString(proof.Bytes()))
}

func TestGenerateProof64(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 64; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.AddEntry(&e); err != nil {
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
		"0a5fbf7d26e753c9fff385ff56484d19f526e770f653653a76f44f0f9d282178"+
		"21eeecff4ffeee35051fb5d7aded04885cbd9a1b7de89a384543cb2e076a7968"+
		"255ef74595c3c548c7a31520e90e7978a059fce4914cbee071262731aa6ba9ea"+
		"1e3fe0e032a292052a300ac3012c4053aadb5c6122d953834e3eb23580bd4565"+
		"12c1f7e29dfedf21d5acd1e01f86d9bfaaaa2657b25bacbab3fc3899e43bf372"+
		"08e3ecd1496bd3a4b74babe63a81032898c46b8145ebfe0176de4ede8e868069",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProof1(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 64; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.AddEntry(&e); err != nil {
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
		"0a5fbf7d26e753c9fff385ff56484d19f526e770f653653a76f44f0f9d282178"+
		"21eeecff4ffeee35051fb5d7aded04885cbd9a1b7de89a384543cb2e076a7968"+
		"255ef74595c3c548c7a31520e90e7978a059fce4914cbee071262731aa6ba9ea"+
		"1e3fe0e032a292052a300ac3012c4053aadb5c6122d953834e3eb23580bd4565"+
		"12c1f7e29dfedf21d5acd1e01f86d9bfaaaa2657b25bacbab3fc3899e43bf372"+
		"08e3ecd1496bd3a4b74babe63a81032898c46b8145ebfe0176de4ede8e868069",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofEmpty(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.AddEntry(&e); err != nil {
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
		"00a5ed99c9977778c45f239769be642d5ba9fd333ad4c4ada5924d0eeb9f1a55"+
		"1cae8054c843d56971b4359a7fc563cbb3afe7fc0e312fb5965cedcf875c4716"+
		"1d1e2f713d0a058b705da55c59cc7b1600f577ba6ddfccb8b640312d7cd26230"+
		"021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofCases(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.AddEntry(&e); err != nil {
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
		"12a103b79aa775d294a27c32dbbee228f04f0dde5f8f39a8bdea0ac99623eb9d"+
		"1e38815ed13fa58751bb8d6474e554a59f9f89a8046f085e257379871a65138c"+
		"1f735fbdc197cc1b1a3aee2f96eb6cb11f37cafd253811e6ca8a0464744394a9"+
		"07ee92c7377df5e038dc2092eb643820c50e98a6f2f6f2ea285d4a84d29c1718",
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
		"12a103b79aa775d294a27c32dbbee228f04f0dde5f8f39a8bdea0ac99623eb9d"+
		"1e38815ed13fa58751bb8d6474e554a59f9f89a8046f085e257379871a65138c"+
		"1f735fbdc197cc1b1a3aee2f96eb6cb11f37cafd253811e6ca8a0464744394a9"+
		"1e031226a2e76231b4c29a182a0083c495b7fdf712ee6df9f030b7212abb925d"+
		"2b2e10bd6e4ba7b655949bca7da3b36e751f0e989c495cb9f9d2c8f884b3e73f",
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
		"12a103b79aa775d294a27c32dbbee228f04f0dde5f8f39a8bdea0ac99623eb9d"+
		"21505710a088ff53242477c909cbdf0bb7753eb15c76fb3862ab6b3da80dfb67"+
		"265d97df25853b130e1769fb421277c93a286487a0876b29db6e87d344e5b271"+
		"17ddf6f66c73719745eeca828537ee30394123a28d16eb51cf51f3bcc0bd03a3"+
		"021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		hex.EncodeToString(proof.Bytes()))
}

func TestVerifyProofFalse(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 8; i++ {
		e := NewEntryFromInts(0, 0, 0, int64(i))
		if err := mt.AddEntry(&e); err != nil {
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
		if err := mt.AddEntry(&e); err != nil {
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

	e := NewEntryFromInts(12, 45, 78, 41)
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
		e := NewEntryFromInts(0, 0, 0, int64(i))
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
		e := NewEntryFromInts(0, 0, 0, int64(i))
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
	assert.Equal(t, "0x1caa9d8ad7a8ad3c6e46fa7101ec5239f869533aa41161db4288c665163c4486", mt.RootKey().Hex())
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

	err := mt.AddEntry(&e)
	assert.Nil(t, err)
	err = mt.AddEntry(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)
	err = mt.AddEntry(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)

	assert.Equal(t,
		"0x2ab68ca530d96aa00eafb1b1f40aa2ed0c8ce6e7aec8bed48b7efd52e919f91b",
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

	err = mt.AddEntry(&e)
	assert.Nil(t, err)

	assert.Equal(t,
		"0x2db0bce5cbafdf5ca0161ea969c71bb7e35f413610f25beeaf776fcd984cb40a",
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
