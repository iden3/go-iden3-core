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
		"f646c27abe5c185391cbad5e701294c8f0fa67f6811ab8f48bbdb82e81e19c14",
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
		"0x73a3ed356712dfee0d0bc338eb0703befe8e4e9247aac39445b1c71bcfebfa1b",
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
		"0xd88207018be295e3100883146194580161e94cf2e6f062bcdc28354bb61b2c1d",
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
		"0xe0137eca5998549c9ce3624563f474e4249a8b8bb6c580dee8098823c8f45129",
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
		"0004000000000000000000000000000000000000000000000000000000000009"+
		"1633242c5b399fa782e69208fd9d21849a25fbcc6c92f43727943af946f81e06"+
		"1b867fc1143906343471012e1a7f40433f9d6a9ccf27b0fe802924376663fc0d",
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
		"000c00000000000000000000000000000000000000000000000000000000091f"+
		"636f8d10fce0fdc27d85f19db3d2d017d26e5a31600d9d95cdf974df97111c25"+
		"81a6c5a4c8ba07aa9cdb2738a586bdf30c4c52667e7e6c6358e03cdcfc67462b"+
		"315d5d2c24510aac8f1ff22fede9b56a4a6c69fd830475dd5771377c01082b2b"+
		"9f97cc867924fc793fa517f77dbf42d655b23f415e1873a702e6876e561a0323"+
		"e30d32898c52c63cef4d98f4d0b04ec8da486b901e822cbe2e781e5500f30705"+
		"f722ebcab864fe6b752bc2c7b68f310c11a8a1316ee229a539ee46326b4c9912"+
		"7fcbd4866e605c0e590f94938cd431b620c9a0e1f7e1e3906e38536beaebf703",
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
		"000c00000000000000000000000000000000000000000000000000000000091f"+
		"636f8d10fce0fdc27d85f19db3d2d017d26e5a31600d9d95cdf974df97111c25"+
		"81a6c5a4c8ba07aa9cdb2738a586bdf30c4c52667e7e6c6358e03cdcfc67462b"+
		"315d5d2c24510aac8f1ff22fede9b56a4a6c69fd830475dd5771377c01082b2b"+
		"9f97cc867924fc793fa517f77dbf42d655b23f415e1873a702e6876e561a0323"+
		"e30d32898c52c63cef4d98f4d0b04ec8da486b901e822cbe2e781e5500f30705"+
		"f722ebcab864fe6b752bc2c7b68f310c11a8a1316ee229a539ee46326b4c9912"+
		"7fcbd4866e605c0e590f94938cd431b620c9a0e1f7e1e3906e38536beaebf703",
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
		"030400000000000000000000000000000000000000000000000000000000000b"+
		"39878ba2a45789446ad36d92bee0a6dab93fc0db8c6f7de397fa1c017b98cf13"+
		"3c28a309bff2050135fede3a2dd3af33270e3078b77f7fbf04b958a283658b17"+
		"05837ab58511e451f89a941ea9b76dbefe03e480526665df2974374f9e371314"+
		"81b7e168ae02356bb67bef5d540125f040ee4d7bff6eb64a35cfcdf2d5761a02"+
		"81b7e168ae02356bb67bef5d540125f040ee4d7bff6eb64a35cfcdf2d5761a02",
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
		"000c00000000000000000000000000000000000000000000000000000000080b"+
		"39878ba2a45789446ad36d92bee0a6dab93fc0db8c6f7de397fa1c017b98cf13"+
		"3c28a309bff2050135fede3a2dd3af33270e3078b77f7fbf04b958a283658b17"+
		"a994437464048acae6113825fdca371fb16ceb962fee3a1a1bcc97c1bd5f731f"+
		"7fcbd4866e605c0e590f94938cd431b620c9a0e1f7e1e3906e38536beaebf703",
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
		"0103000000000000000000000000000000000000000000000000000000000007"+
		"39878ba2a45789446ad36d92bee0a6dab93fc0db8c6f7de397fa1c017b98cf13"+
		"3c28a309bff2050135fede3a2dd3af33270e3078b77f7fbf04b958a283658b17"+
		"9a5c3e573d133d33f34e3f114200967925b062a62f1a9bd24da2d95645b3ed2e",
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
		"0103000000000000000000000000000000000000000000000000000000000007"+
		"afe9665c235eda8df07635ff38f056b763aabdfb50b05767e9c6b1d2bd36fd11"+
		"1c6d639ef2ee1ed14505f50b2fc993a49273566abf5b07157e22af210a153c17"+
		"27d248811eb867a6d6dea170694d0d62aec70c406616a48ba65a43cf7c4f900c",
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

func TestMTWalkDumpClaims(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		rawIndex := strconv.Itoa(i) + "-testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest"
		rawData := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest-"
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

	dumpedClaims := []string{"0x000000000000000000000000000000000000000000000000362d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400",
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
		"0x000000000000000000000000000000000000000000000000382d7465737474006573747465737474657374746573747465737474657374746573747465737400746573747465737474657374746573747465737474657374746573747465730074746573747465737474657374746573747465737474657374746573747465007465737474657374746573747465737474657374746573747465730000000000747465737474657374746573747465737474657374746573747465737474650073747465737474657374746573747465737474657374746573747465737474006573747465737474657374746573747465737474657374746573747465737400"}

	err := mt.ImportDumpedClaims(dumpedClaims)
	assert.Nil(t, err)
	assert.Equal(t, "0xd8ddc80b67bccb0e1d0c50a170cf7199fddfc92acd156614ab1897b1929ffc1e", mt.RootKey().Hex())
}

func TestMTWalkDumpClaimsAndImportDumpedClaims(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 16; i++ {
		rawIndex := strconv.Itoa(i) + "-testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest"
		rawData := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest-"
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

	e := NewEntryFromInts(12, 45, 78, 41, 35, 80, 54, 42)

	err := mt.AddEntry(&e)
	assert.Nil(t, err)
	err = mt.AddEntry(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)
	err = mt.AddEntry(&e)
	assert.Equal(t, err, ErrEntryIndexAlreadyExists)

	assert.Equal(t,
		"0x73a3ed356712dfee0d0bc338eb0703befe8e4e9247aac39445b1c71bcfebfa1b",
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
		"0x2123d42746602583b0e9dcb41353292238684f18c506b59dbf0edcdbc5e35d15",
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
