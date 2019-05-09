package core

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/stretchr/testify/assert"
)

func TestProof(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	idAddr0, err := IDFromString("1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0 := NewClaimSetRootKey(idAddr0, rootKey0)
	err = mt.Add(claim0.Entry())
	assert.Nil(t, err)

	idAddr1, err := IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZohfu")
	assert.Nil(t, err)
	rootKey1 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b})
	claim1 := NewClaimSetRootKey(idAddr1, rootKey1)
	err = mt.Add(claim1.Entry())
	assert.Nil(t, err)

	mtp, err := GetClaimProofByHi(mt, claim0.Entry().HIndex())
	assert.Nil(t, err)

	j, err := json.Marshal(mtp)
	assert.Nil(t, err)

	relayAddr := common.Address{}
	verified, err := VerifyProofClaim(relayAddr, mtp)
	assert.Nil(t, err)
	assert.True(t, verified)
}
