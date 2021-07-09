package writerhttp

import (
	"fmt"
	"github.com/iden3/go-iden3-core/crypto"
	"github.com/iden3/go-merkletree-sql/db/memory"
	"os"
	"strconv"
	"testing"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-merkletree-sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// If generateTest is true, the checked values will be used to generate a test vector
var generateTest = false

func TestHttpPublicGetPublicData(t *testing.T) {
	// create RootsTree & RevocationsTree
	cltMt, err := merkletree.NewMerkleTree(memory.NewMemoryStorage(), 140)
	assert.Nil(t, err)
	rotMt, err := merkletree.NewMerkleTree(memory.NewMemoryStorage(), 140)
	assert.Nil(t, err)
	retMt, err := merkletree.NewMerkleTree(memory.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	// add some leafs to both MerkleTrees
	for i := 0; i < 10; i++ {
		rootBigInt, err := crypto.PoseidonHashBytes([]byte(strconv.Itoa(i)))
		assert.Nil(t, err)
		root := merkletree.NewHashFromBigInt(rootBigInt)
		err = claims.AddLeafRootsTree(rotMt, root)
		assert.Nil(t, err)

		nonce := uint32(i)
		version := uint32(i)
		err = claims.AddLeafRevocationsTree(retMt, nonce, version)
		assert.Nil(t, err)
	}

	testgen.CheckTestValue(t, "rootRootsTree1", rotMt.Root().Hex())
	testgen.CheckTestValue(t, "rootRevocationsTree1", retMt.Root().Hex())

	idenPubOffChainWriteHttp, err := NewIdenPubOffChainWriteHttp(NewConfigDefault("http://foo.bar"), db.NewMemoryStorage())
	require.Nil(t, err)

	idenState := core.IdenState(cltMt.Root(), retMt.Root(), rotMt.Root())

	publicData := idenpuboffchain.PublicData{
		IdenState:           idenState,
		ClaimsTreeRoot:      cltMt.Root(),
		RevocationsTreeRoot: retMt.Root(),
		RevocationsTree:     retMt,
		RootsTreeRoot:       rotMt.Root(),
		RootsTree:           rotMt,
	}

	err = idenPubOffChainWriteHttp.Publish(&core.ID{}, &publicData)
	assert.Nil(t, err)

	pubDataBlobs, err := idenPubOffChainWriteHttp.GetPublicData(nil)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "rootRootsTree1", pubDataBlobs.RootsTreeRoot.Hex())
	assert.Equal(t, rotMt.Root().Hex(), pubDataBlobs.RootsTreeRoot.Hex())
	testgen.CheckTestValue(t, "rootRevocationsTree1", pubDataBlobs.RevocationsTreeRoot.Hex())
	assert.Equal(t, retMt.Root().Hex(), pubDataBlobs.RevocationsTreeRoot.Hex())

	pubDataBlobs2, err := idenPubOffChainWriteHttp.GetPublicData(idenState)
	assert.Nil(t, err)
	assert.Equal(t, pubDataBlobs, pubDataBlobs2)

	_, err = idenpuboffchain.NewPublicDataFromBlobs(pubDataBlobs)
	require.Nil(t, err)
}

// Assert that IdenPubOffChainWrite follows the IdenPubOffChainWriter interface
func TestIdenPubOffChainWriteInterface(t *testing.T) {
	var idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter //nolint:gosimple
	idenPubOffChainWriteHttp, err := NewIdenPubOffChainWriteHttp(NewConfigDefault("http://foo.bar"), db.NewMemoryStorage())
	require.Nil(t, err)
	idenPubOffChainWrite = idenPubOffChainWriteHttp
	require.NotNil(t, idenPubOffChainWrite)
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
		root0, _ := crypto.PoseidonHashBytes([]byte("root0"))
		testgen.SetTestValue("root0", merkletree.NewHashFromBigInt(root0).Hex())
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
