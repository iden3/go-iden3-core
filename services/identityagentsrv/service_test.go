package identityagentsrv

import (
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
