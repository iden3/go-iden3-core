package identityagentsrv

import (
	// "encoding/hex"
	// "fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"

	"github.com/stretchr/testify/require"
)

//// BEGIN Helper functions

func createIdentityLoadAgent(t *testing.T) (*core.ID, *babyjub.PublicKey, *Agent) {
	sto, err := NewTestingStorage()
	require.Nil(t, err)

	ia := New(sto, &RootUpdaterMock{})

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kopComp babyjub.PublicKeyComp
	err = kopComp.UnmarshalText([]byte(kopStr))
	require.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	require.Nil(t, err)
	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub).Entry()

	id, _, err := ia.CreateIdentity(claimKOp, nil)
	require.Nil(t, err)
	require.Equal(t, "1GURWwRa5YQA8KA2AdmGANhXSpAupfpy2VsHse2QU", id.String())

	agent, err := ia.NewAgent(id)
	require.Nil(t, err)
	return id, kopPub, agent
}

//// END

var service *Service

var rmDirs []string

func NewTestingStorage() (db.Storage, error) {
	dir, err := ioutil.TempDir("", "db")
	rmDirs = append(rmDirs, dir)
	if err != nil {
		return nil, err
	}
	sto, err := db.NewLevelDbStorage(dir, false)
	return sto, err
}

func testServiceInterficeFunction(ia *Service) {

}

func TestServiceInterface(t *testing.T) {
	sto, err := NewTestingStorage()
	require.Nil(t, err)

	ia := New(sto, &RootUpdaterMock{})

	testServiceInterficeFunction(ia)
}

func TestNewIdentity(t *testing.T) {
	sto, err := NewTestingStorage()
	require.Nil(t, err)

	ia := New(sto, &RootUpdaterMock{})

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kopComp babyjub.PublicKeyComp
	err = kopComp.UnmarshalText([]byte(kopStr))
	require.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	require.Nil(t, err)
	kDis := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kReen := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kUpdateRoot := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub).Entry()
	claimKDis := core.NewClaimAuthEthKey(kDis, core.EthKeyTypeDisable).Entry()
	claimKReen := core.NewClaimAuthEthKey(kReen, core.EthKeyTypeReenable).Entry()
	claimKUpdateRoot := core.NewClaimAuthEthKey(kUpdateRoot, core.EthKeyTypeUpdateRoot).Entry()

	id, proofKOp, err := ia.CreateIdentity(claimKOp, []*merkletree.Entry{claimKDis, claimKReen, claimKUpdateRoot})
	require.Nil(t, err)

	require.Equal(t, "1FJS9Bb6LE5GFpnNAyuS657jmpdjSc1MVHts54FUP", id.String())
	proofKOpVerified, err := proofKOp.Verify(proofKOp.Proof.Root)
	require.Nil(t, err)
	require.True(t, proofKOpVerified)
}

func TestAddClaims(t *testing.T) {
	_, _, agent := createIdentityLoadAgent(t)

	// create claim to be added
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	ethKeyType := core.EthKeyTypeUpgrade
	c0 := core.NewClaimAuthEthKey(ethKey, ethKeyType).Entry()
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	ethKeyType = core.EthKeyTypeUpgrade
	c1 := core.NewClaimAuthEthKey(ethKey, ethKeyType).Entry()

	err := agent.AddClaims([]*merkletree.Entry{c0, c1})
	require.Nil(t, err)

	// should give collision error when adding the claim already added
	err = agent.AddClaims([]*merkletree.Entry{c0, c1})
	require.Equal(t, merkletree.ErrEntryIndexAlreadyExists, err)
}

func TestGetClaims(t *testing.T) {
	_, kopPub, agent := createIdentityLoadAgent(t)
	emittedClaimsAfterGenesis, err := agent.ClaimsEmitted()
	require.Nil(t, err)

	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub).Entry()

	// create claims to be added
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	c0 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade).Entry()
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	c1 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade).Entry()

	err = agent.AddClaims([]*merkletree.Entry{c0, c1})
	require.Nil(t, err)

	emittedClaims, err := agent.ClaimsEmitted()
	require.Nil(t, err)
	receivedClaims, err := agent.ClaimsReceived()
	require.Nil(t, err)
	require.Equal(t, c0.Bytes(), emittedClaims[1].Bytes())
	require.Equal(t, claimKOp.Bytes(), emittedClaims[2].Bytes())
	require.Equal(t, c1.Bytes(), emittedClaims[0].Bytes())
	require.Equal(t, 3, len(emittedClaims)) // 3 emitted claims, 1 on genesistree, and 2 after genesistree
	require.Equal(t, 0, len(receivedClaims))

	genesisClaims, err := agent.ClaimsGenesis()
	require.Nil(t, err)
	require.Equal(t, 1, len(genesisClaims))

	for i, claim := range emittedClaimsAfterGenesis {
		require.Equal(t, claim.Bytes(), genesisClaims[i].Bytes())
	}
}

func TestGetClaimByHi(t *testing.T) {
	_, _, agent := createIdentityLoadAgent(t)

	// create claims to be added
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	c0 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade).Entry()
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	c1 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade).Entry()

	err := agent.AddClaims([]*merkletree.Entry{c0, c1})
	require.Nil(t, err)

	claim, _, err := agent.GetClaimByHi(c0.HIndex())
	require.Nil(t, err)
	require.Equal(t, c0.Bytes(), claim.Bytes())

	claim, _, err = agent.GetClaimByHi(c1.HIndex())
	require.Nil(t, err)
	require.Equal(t, c1.Bytes(), claim.Bytes())
}

func TestExportMT(t *testing.T) {
	_, _, agent := createIdentityLoadAgent(t)

	// create claims to be added
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	c0 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade).Entry()
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	c1 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade).Entry()

	err := agent.AddClaims([]*merkletree.Entry{c0, c1})
	require.Nil(t, err)

	mt, err := agent.ExportMT()
	require.Nil(t, err)
	require.Equal(t, agent.mt.RootKey().Hex(), mt[0][0]) // crop first 4 from mt map, as the first 2 are for '03' indicating the node type of the MerkleTree, the other 2 are for the '0x'

	require.Equal(t, 5, len(mt))
}

func TestGetCurrentRoot(t *testing.T) {
	_, _, agent := createIdentityLoadAgent(t)

	// create claims to be added
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	c0 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade).Entry()
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	c1 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade).Entry()

	err := agent.AddClaims([]*merkletree.Entry{c0, c1})
	require.Nil(t, err)

	mt, err := agent.ExportMT()
	require.Nil(t, err)

	require.Equal(t, agent.mt.RootKey().Hex(), mt[0][0])

	// this will be for integration tests
	// root, err := agent.GetCurrentRoot()
	// require.Nil(t, err)
	// require.Equal(t, agent.mt.RootKey().Hex(), root.Local.Hex())
}

func TestMain(m *testing.M) {
	result := m.Run()
	for _, dir := range rmDirs {
		os.RemoveAll(dir)
	}
	os.Exit(result)
}
