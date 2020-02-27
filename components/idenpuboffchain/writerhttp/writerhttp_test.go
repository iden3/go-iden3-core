package writerhttp

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// If generateTest is true, the checked values will be used to generate a test vector
var generateTest = false

func TestHttpPublicGetPublicData(t *testing.T) {
	// create RootsTree & RevocationsTree
	cltMt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)
	rotMt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)
	retMt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	// add some leafs to both MerkleTrees
	for i := 0; i < 10; i++ {
		rootBigInt, err := poseidon.HashBytes([]byte(strconv.Itoa(i)))
		assert.Nil(t, err)
		root := merkletree.BigIntToHash(rootBigInt)
		err = claims.AddLeafRootsTree(rotMt, &root)
		assert.Nil(t, err)

		nonce := uint32(i)
		version := uint32(i)
		err = claims.AddLeafRevocationsTree(retMt, nonce, version)
		assert.Nil(t, err)
	}

	testgen.CheckTestValue(t, "rootRootsTree1", rotMt.RootKey().Hex())
	testgen.CheckTestValue(t, "rootRevocationsTree1", retMt.RootKey().Hex())

	idenPubOffChainWriteHttp, err := NewIdenPubOffChainWriteHttp(NewConfigDefault("http://foo.bar"), db.NewMemoryStorage())
	require.Nil(t, err)

	idenState := core.IdenState(cltMt.RootKey(), retMt.RootKey(), rotMt.RootKey())

	publicData := idenpuboffchain.PublicData{
		IdenState:           idenState,
		ClaimsTreeRoot:      cltMt.RootKey(),
		RevocationsTreeRoot: retMt.RootKey(),
		RevocationsTree:     retMt,
		RootsTreeRoot:       rotMt.RootKey(),
		RootsTree:           rotMt,
	}

	err = idenPubOffChainWriteHttp.Publish(&core.ID{}, &publicData)
	assert.Nil(t, err)

	pubDataBlobs, err := idenPubOffChainWriteHttp.GetPublicData(nil)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "rootRootsTree1", pubDataBlobs.RootsTreeRoot.Hex())
	assert.Equal(t, rotMt.RootKey().Hex(), pubDataBlobs.RootsTreeRoot.Hex())
	testgen.CheckTestValue(t, "rootRevocationsTree1", pubDataBlobs.RevocationsTreeRoot.Hex())
	assert.Equal(t, retMt.RootKey().Hex(), pubDataBlobs.RevocationsTreeRoot.Hex())

	pubDataBlobs2, err := idenPubOffChainWriteHttp.GetPublicData(idenState)
	assert.Nil(t, err)
	assert.Equal(t, pubDataBlobs, pubDataBlobs2)

	_, err = idenpuboffchain.NewPublicDataFromBlobs(pubDataBlobs)
	require.Nil(t, err)
}

func initTest() {
	// Init test
	err := testgen.InitTest("idenpuboffchainwriterhttp", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
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
