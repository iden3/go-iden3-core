package idenpub

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	cryptoConstants "github.com/iden3/go-iden3-crypto/constants"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/assert"
)

// If generateTest is true, the checked values will be used to generate a test vector
var generateTest = false

func TestAddLeafRoT(t *testing.T) {
	root := hexStringToHash(testgen.GetTestValue("root0").(string))

	mt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	err = AddLeafRoT(mt, root)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "rootRoT0", mt.RootKey().Hex())

	proof, err := mt.GenerateProof(NewLeafRoT(root).Entry().HIndex(), nil)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "proofLeafRoT", hex.EncodeToString(proof.Bytes()))
}

func TestAddLeafReT(t *testing.T) {
	nonce := uint32(testgen.GetTestValue("nonce0").(float64))
	version := uint32(testgen.GetTestValue("version0").(float64))

	mt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	err = AddLeafReT(mt, nonce, version)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "rootReT0", mt.RootKey().Hex())

	proof, err := mt.GenerateProof(NewLeafReT(nonce, version).Entry().HIndex(), nil)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "proofReT", hex.EncodeToString(proof.Bytes()))
}

func TestHttpPublicGetPublicData(t *testing.T) {
	// create RoT & ReT
	cltMt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)
	rotMt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)
	retMt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	// add some leafs to both MerkleTrees
	for i := 0; i < 10; i++ {
		root, err := poseidon.HashBytes([]byte(strconv.Itoa(i)))
		assert.Nil(t, err)
		err = AddLeafRoT(rotMt, merkletree.BigIntToHash(new(big.Int).Mod(root, cryptoConstants.Q)))
		assert.Nil(t, err)

		nonce := uint32(i)
		version := uint32(i)
		err = AddLeafReT(retMt, nonce, version)
		assert.Nil(t, err)
	}

	testgen.CheckTestValue(t, "rootRoT1", rotMt.RootKey().Hex())
	testgen.CheckTestValue(t, "rootReT1", retMt.RootKey().Hex())

	idenPubHTTP := NewIdenPubHTTP(db.NewMemoryStorage(), rotMt, retMt)

	idenState := hexStringToHash(testgen.GetTestValue("idenState0").(string))

	err = idenPubHTTP.Publish(&idenState, cltMt.RootKey(), rotMt.RootKey(), retMt.RootKey())
	assert.Nil(t, err)

	pubData, err := idenPubHTTP.GetPublicData()
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "rootRoT1", pubData.RoTRoot.Hex())
	assert.Equal(t, rotMt.RootKey().Hex(), pubData.RoTRoot.Hex())
	testgen.CheckTestValue(t, "rootReT1", pubData.ReTRoot.Hex())
	assert.Equal(t, retMt.RootKey().Hex(), pubData.ReTRoot.Hex())

	testgen.CheckTestValue(t, "RoT1", hex.EncodeToString(pubData.RoT))
	testgen.CheckTestValue(t, "ReT1", hex.EncodeToString(pubData.ReT))
}

func initTest() {
	// Init test
	err := testgen.InitTest("idenpub", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
		idenState0, err := poseidon.HashBytes([]byte("idenState0"))
		if err != nil {
			panic(err)
		}
		testgen.SetTestValue("idenState0", merkletree.BigIntToHash(idenState0).Hex())
		root0, err := poseidon.HashBytes([]byte("root0"))
		if err != nil {
			panic(err)
		}
		testgen.SetTestValue("root0", merkletree.BigIntToHash(root0).Hex())
		testgen.SetTestValue("nonce0", float64(5))
		testgen.SetTestValue("version0", float64(5))
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

func hexStringToHash(s string) merkletree.Hash {
	b, err := common3.HexDecode(s)
	if err != nil {
		panic(err)
	}
	var b32 [merkletree.ElemBytesLen]byte
	copy(b32[:], b[:32])
	return merkletree.Hash(merkletree.ElemBytes(b32))
}
