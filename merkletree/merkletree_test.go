package merkletree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
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

func HexDecode(h string) ([]byte, error) {
	if strings.HasPrefix(h, "0x") {
		h = h[2:]
	}
	return hex.DecodeString(h)
}
func NewEntryFromHexs(a, b, c, d string) (e Entry, err error) {
	e.Data, err = HexsToData(a, b, c, d)
	if err != nil {
		return e, err
	}
	return e, nil
}
func HexsToData(_a, _b, _c, _d string) (Data, error) {
	aBytes, err := HexDecode(_a)
	if err != nil {
		return Data{}, err
	}
	a := new(big.Int).SetBytes(aBytes)

	bBytes, err := HexDecode(_b)
	if err != nil {
		return Data{}, err
	}
	b := new(big.Int).SetBytes(bBytes)

	cBytes, err := HexDecode(_c)
	if err != nil {
		return Data{}, err
	}
	c := new(big.Int).SetBytes(cBytes)

	dBytes, err := HexDecode(_d)
	if err != nil {
		return Data{}, err
	}
	d := new(big.Int).SetBytes(dBytes)

	return BigIntsToData(a, b, c, d), nil
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
		"067f3202335ea256ae6e6aadcd2d5f7f4b06a00b2d1e0de903980d5ab552dc70",
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
		"0x2bf39430aa2482fc1e2f170179c8cab126b0f55f71edc8d333f4c80cb4e798f5",
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
		"0x18cbbb4a507b66fb20c7f7548629e38ceba0d4feb61c709d3655128fd2cb86c1",
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
		"0x2f398ded371e7ecd5b1728ce86aba7398272728c2a6d685ae37df0f7645bf254",
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
		"0005000000000000000000000000000000000000000000000000000000000011"+
		"19b861dd1137c0355e2387a4485c5be6fb6e96b58bbbd6ec3dc0685aac6ab21b"+
		"1db94892d26681f45bd59a26d9b63b4a3b72e49426cc99413de511ce936d0ea2",
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
		"00080000000000000000000000000000000000000000000000000000000000ff"+
		"27d4dd8a257bfdf256f2228d5d0efe1e76f122724d6b758fcb955e29eca181b5"+
		"19ef7310fb6f536f3bae38426596bafe09565f6a5e10a6a00b15b6ea2a358dcb"+
		"0dbc0dbdf300d56b3f856c857143078885be23c9f4af5858e01f744ce786a8af"+
		"255374c08e2fb29e5b718d47ba53cd59cc8865c5e83601c2e12e097416b75a81"+
		"13f5e61b352fe2e8214e11c8c8cb955a40908b4821806be4f8cc44c4a7da5d93"+
		"044d24b423c6812821268c877b68c9a6a5e8e027a2ec5622091544a9a1efd861"+
		"2c29b6a10a02185110721050ff704afac110e0c1d843bac8b6649252eb8af9d6"+
		"07138a66cae812df472805aadd26784e0da092bbdba296e7700268d6f6812bd8",
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
		"00080000000000000000000000000000000000000000000000000000000000ff"+
		"27d4dd8a257bfdf256f2228d5d0efe1e76f122724d6b758fcb955e29eca181b5"+
		"19ef7310fb6f536f3bae38426596bafe09565f6a5e10a6a00b15b6ea2a358dcb"+
		"0dbc0dbdf300d56b3f856c857143078885be23c9f4af5858e01f744ce786a8af"+
		"255374c08e2fb29e5b718d47ba53cd59cc8865c5e83601c2e12e097416b75a81"+
		"13f5e61b352fe2e8214e11c8c8cb955a40908b4821806be4f8cc44c4a7da5d93"+
		"044d24b423c6812821268c877b68c9a6a5e8e027a2ec5622091544a9a1efd861"+
		"2c29b6a10a02185110721050ff704afac110e0c1d843bac8b6649252eb8af9d6"+
		"07138a66cae812df472805aadd26784e0da092bbdba296e7700268d6f6812bd8",
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
		"0103000000000000000000000000000000000000000000000000000000000007"+
		"21e0fd75e6b2a35c03e77f0926d3527f53c7756306e0a0fa79f4e6f74f572d37"+
		"02094a7bacee0d07157147d4e7ef169943f6b8b404776b381d4f6935b0316a25"+
		"1c536405a78590470e8ab304afcca954afcbbfdc5818e2d00fbc223e9e4a0cf3",
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
		"0007000000000000000000000000000000000000000000000000000000000043"+
		"0b17095aa925dc78a07fdb53f1589890d018b61d9ec2dedc800ffd5a06c6d37a"+
		"137d66cc357c7ce5d3b8dd239492f04a84c189c7784a081e705def9846ac7d02"+
		"1f777df2736b8f339bc822ce0e82681c3d5712d6c8c1681b9f8a8da33f027bd3",
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
		"0307000000000000000000000000000000000000000000000000000000000043"+
		"0b17095aa925dc78a07fdb53f1589890d018b61d9ec2dedc800ffd5a06c6d37a"+
		"137d66cc357c7ce5d3b8dd239492f04a84c189c7784a081e705def9846ac7d02"+
		"1f777df2736b8f339bc822ce0e82681c3d5712d6c8c1681b9f8a8da33f027bd3"+
		"2bd5a51b497073b952be60de9a83802c64c49aececbd504cbc8269bef021e9e6"+
		"06d4571fb9634e4bed32e265f91a373a852c476656c5c13b09bc133ac61bc5a6",
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
		"0103000000000000000000000000000000000000000000000000000000000007"+
		"0b17095aa925dc78a07fdb53f1589890d018b61d9ec2dedc800ffd5a06c6d37a"+
		"137d66cc357c7ce5d3b8dd239492f04a84c189c7784a081e705def9846ac7d02"+
		"12ad30ae8e4ba3c5bed549923763c6cdf1d32466fba659448ce987a6b6cf3d08",
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
		"0x2bf39430aa2482fc1e2f170179c8cab126b0f55f71edc8d333f4c80cb4e798f5",
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
		"0x2b5433214c653859dde8ba916fb6b31cb93c26bcca70630219e3bb4f587886df",
		mt.RootKey().Hex())
}
