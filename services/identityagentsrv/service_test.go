package identityagentsrv

import (
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"

	"github.com/stretchr/testify/assert"
)

var service *ServiceImpl

func NewTestingStorage() (db.Storage, error) {
	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		return nil, err
	}
	sto, err := db.NewLevelDbStorage(dir, false)
	return sto, err
}

func TestNewIdentity(t *testing.T) {
	sto, err := NewTestingStorage()
	assert.Nil(t, err)

	ia := New(sto)

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kopComp babyjub.PublicKeyComp
	err = kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	kDis := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kReen := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kUpdateRoot := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub)
	claimKDis := core.NewClaimAuthEthKey(kDis, core.EthKeyTypeDisable)
	claimKReen := core.NewClaimAuthEthKey(kReen, core.EthKeyTypeReenable)
	claimKUpdateRoot := core.NewClaimAuthEthKey(kUpdateRoot, core.EthKeyTypeUpdateRoot)

	id, proofKOp, err := ia.NewIdentity(claimKOp, []merkletree.Claim{claimKDis, claimKReen, claimKUpdateRoot})
	assert.Nil(t, err)

	assert.Equal(t, "117aFcVWPyypFbjCuHRKaAaTV7nN3yT9q6PthJpm96", id.String())
	var relayPk *babyjub.PublicKey
	proofKOpVerified, err := core.VerifyProofClaim(relayPk, proofKOp)
	assert.Nil(t, err)
	assert.True(t, proofKOpVerified)
}

func TestAddClaim(t *testing.T) {
	sto, err := NewTestingStorage()
	assert.Nil(t, err)

	ia := New(sto)

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kopComp babyjub.PublicKeyComp
	err = kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub)

	id, _, err := ia.NewIdentity(claimKOp, []merkletree.Claim{})
	assert.Nil(t, err)

	// create claim to be added
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	ethKeyType := core.EthKeyTypeUpgrade
	c0 := core.NewClaimAuthEthKey(ethKey, ethKeyType)

	err = ia.AddClaim(id, c0)
	assert.Nil(t, err)

	// should give collision error when adding the claim already added
	err = ia.AddClaim(id, c0)
	assert.Equal(t, merkletree.ErrEntryIndexAlreadyExists, err)
}

func TestAddClaims(t *testing.T) {
	sto, err := NewTestingStorage()
	assert.Nil(t, err)

	ia := New(sto)

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kopComp babyjub.PublicKeyComp
	err = kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub)

	id, _, err := ia.NewIdentity(claimKOp, []merkletree.Claim{})
	assert.Nil(t, err)

	// create claim to be added
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	ethKeyType := core.EthKeyTypeUpgrade
	c0 := core.NewClaimAuthEthKey(ethKey, ethKeyType)
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	ethKeyType = core.EthKeyTypeUpgrade
	c1 := core.NewClaimAuthEthKey(ethKey, ethKeyType)

	err = ia.AddClaims(id, []merkletree.Claim{c0, c1})
	assert.Nil(t, err)

	// should give collision error when adding the claim already added
	err = ia.AddClaims(id, []merkletree.Claim{c0, c1})
	assert.Equal(t, merkletree.ErrEntryIndexAlreadyExists, err)
}

func TestGetClaims(t *testing.T) {
	sto, err := NewTestingStorage()
	assert.Nil(t, err)

	ia := New(sto)

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kopComp babyjub.PublicKeyComp
	err = kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub)
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	id, _, err := ia.NewIdentity(claimKOp, []merkletree.Claim{})
	assert.Nil(t, err)

	// create claims to be added
	ethKey = common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	c0 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade)
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	c1 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade)

	err = ia.AddClaims(id, []merkletree.Claim{c0, c1})
	assert.Nil(t, err)

	emittedClaims, receivedClaims, err := ia.GetAllClaims(id)
	assert.Nil(t, err)
	assert.Equal(t, c0.Entry().Bytes(), emittedClaims[0].Claim.Entry().Bytes())
	assert.Equal(t, claimKOp.Entry().Bytes(), emittedClaims[1].Claim.Entry().Bytes())
	assert.Equal(t, c1.Entry().Bytes(), emittedClaims[2].Claim.Entry().Bytes())
	assert.Equal(t, 3, len(emittedClaims)) // 3 emitted claims, 1 on genesistree, and 2 after genesistree
	assert.Equal(t, 0, len(receivedClaims))
}

func TestGetClaimByHi(t *testing.T) {
	sto, err := NewTestingStorage()
	assert.Nil(t, err)

	ia := New(sto)

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kopComp babyjub.PublicKeyComp
	err = kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub)
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	id, _, err := ia.NewIdentity(claimKOp, []merkletree.Claim{})
	assert.Nil(t, err)

	// create claims to be added
	ethKey = common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	c0 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade)
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	c1 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade)

	err = ia.AddClaims(id, []merkletree.Claim{c0, c1})
	assert.Nil(t, err)

	claim, _, err := ia.GetClaimByHi(id, c0.Entry().HIndex())
	assert.Nil(t, err)
	assert.Equal(t, c0.Entry().Bytes(), claim.Entry().Bytes())

	claim, _, err = ia.GetClaimByHi(id, c1.Entry().HIndex())
	assert.Nil(t, err)
	assert.Equal(t, c1.Entry().Bytes(), claim.Entry().Bytes())
}

func TestGetFullMT(t *testing.T) {
	sto, err := NewTestingStorage()
	assert.Nil(t, err)

	ia := New(sto)

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kopComp babyjub.PublicKeyComp
	err = kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	claimKOp := core.NewClaimAuthorizeKSignBabyJub(kopPub)
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	id, _, err := ia.NewIdentity(claimKOp, []merkletree.Claim{})
	assert.Nil(t, err)

	// create claims to be added
	ethKey = common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	c0 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade)
	ethKey = common.HexToAddress("0x3d380182Cd261CdcD413e4B8D17c89c943c39b1A")
	c1 := core.NewClaimAuthEthKey(ethKey, core.EthKeyTypeUpgrade)

	err = ia.AddClaims(id, []merkletree.Claim{c0, c1})
	assert.Nil(t, err)

	mt, err := ia.GetFullMT(id)
	assert.Nil(t, err)
	idStorages, err := ia.LoadIdStorages(id)
	assert.Nil(t, err)
	assert.Equal(t, idStorages.mt.RootKey().Hex()[2:], mt["0x"+hex.EncodeToString([]byte("currentroot"))][4:]) // crop first 4 from mt map, as the first 2 are for '03' indicating the node type of the MerkleTree, the other 2 are for the '0x'

	count := 0
	for _, _ = range mt {
		count++
	}
	assert.Equal(t, 7, count)
}
