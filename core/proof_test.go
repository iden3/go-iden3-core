package core

import (
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

	id0, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0 := NewClaimSetRootKey(id0, rootKey0)
	err = mt.Add(claim0.Entry())
	assert.Nil(t, err)

	id1, err := IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoWij")
	assert.Nil(t, err)
	rootKey1 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b})
	claim1 := NewClaimSetRootKey(id1, rootKey1)
	err = mt.Add(claim1.Entry())
	assert.Nil(t, err)

	mtp, err := GetClaimProofByHi(mt, claim0.Entry().HIndex())
	assert.Nil(t, err)

	// j, err := json.Marshal(mtp)
	// assert.Nil(t, err)

	relayAddr := common.Address{}
	verified, err := VerifyProofClaim(relayAddr, mtp)
	assert.Nil(t, err)
	assert.True(t, verified)
}

func TestGetPredicateProof(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	id0, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0 := NewClaimSetRootKey(id0, rootKey0)
	err = mt.Add(claim0.Entry())
	assert.Nil(t, err)
	oldRoot := mt.RootKey()

	mtp, err := GetClaimProofByHi(mt, claim0.Entry().HIndex())
	assert.Nil(t, err)

	relayAddr := common.Address{}
	verified, err := VerifyProofClaim(relayAddr, mtp)
	assert.Nil(t, err)
	assert.True(t, verified)

	p, err := GetPredicateProof(mt, oldRoot, claim0.Entry().HIndex())
	assert.Nil(t, err)

	assert.True(t, merkletree.VerifyProof(mt.RootKey(), p.MtpExist, claim0.Entry().HIndex(), claim0.Entry().HValue()))

	claim0Entry := GetNextVersionEntry(claim0.Entry())
	assert.True(t, merkletree.VerifyProof(mt.RootKey(), p.MtpNonExistNextVersion, claim0Entry.HIndex(), claim0Entry.HValue()))

	assert.True(t, merkletree.VerifyProof(mt.RootKey(), p.MtpNonExistInOldRoot, claim0.Entry().HIndex(), claim0.Entry().HValue()))
}

func TestGenerateAndVerifyPredicateProofOfClaimVersion0(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	id0, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0 := NewClaimSetRootKey(id0, rootKey0)
	// oldRoot is the root before adding the claim that we want to prove that we added correctly
	oldRoot := mt.RootKey()
	err = mt.Add(claim0.Entry())
	assert.Nil(t, err)

	predicateProof, err := GetPredicateProof(mt, oldRoot, claim0.Entry().HIndex())
	assert.Nil(t, err)

	_, v := getClaimTypeVersion(predicateProof.LeafEntry)
	assert.Equal(t, uint32(0), v)
	assert.Equal(t, predicateProof.OldRoot.Hex(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	assert.NotEqual(t, predicateProof.OldRoot.Hex(), predicateProof.Root.Hex())

	assert.True(t, VerifyPredicateProof(predicateProof))
}

func TestGenerateAndVerifyPredicateProofOfClaimVersion1(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	id0, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0 := NewClaimSetRootKey(id0, rootKey0)
	err = mt.Add(claim0.Entry())
	assert.Nil(t, err)

	// add the same claim, but with version 1
	claim1 := &ClaimSetRootKey{
		Version: claim0.Version + 1,
		Era:     0,
		Id:      claim0.Id,
		RootKey: claim0.RootKey,
	}
	err = mt.Add(claim1.Entry())
	assert.Nil(t, err)
	// add the same claim, but with version 2
	claim2 := &ClaimSetRootKey{
		Version: claim1.Version + 1,
		Era:     0,
		Id:      claim0.Id,
		RootKey: claim0.RootKey,
	}
	err = mt.Add(claim2.Entry())
	assert.Nil(t, err)

	// oldRoot is the root before adding the claim that we want to prove that we added correctly
	oldRoot := mt.RootKey()

	// add the same claim, but with version 3
	claim3 := &ClaimSetRootKey{
		Version: claim2.Version + 1,
		Era:     0,
		Id:      claim0.Id,
		RootKey: claim0.RootKey,
	}
	err = mt.Add(claim3.Entry())
	assert.Nil(t, err)

	// expect error, as we are trying to generate a proof of a claim which one the next version
	_, err = GetPredicateProof(mt, oldRoot, claim0.Entry().HIndex())
	assert.Equal(t, err, ErrRevokedClaim)
	_, err = GetPredicateProof(mt, oldRoot, claim1.Entry().HIndex())
	assert.Equal(t, err, ErrRevokedClaim)
	_, err = GetPredicateProof(mt, oldRoot, claim2.Entry().HIndex())
	assert.Equal(t, err, ErrRevokedClaim)

	predicateProof, err := GetPredicateProof(mt, oldRoot, claim3.Entry().HIndex())
	assert.Nil(t, err)

	_, v := getClaimTypeVersion(predicateProof.LeafEntry)
	assert.Equal(t, uint32(3), v)
	assert.Equal(t, predicateProof.OldRoot.Hex(), "0x22e52a06ab4c4824377859745863e4d306b74b451ec9d09637a4d0c0ad14667e")
	assert.NotEqual(t, predicateProof.OldRoot.Hex(), predicateProof.Root.Hex())

	assert.Equal(t, predicateProof.MtpNonExistInOldRoot.Siblings[0], predicateProof.MtpExist.Siblings[0])

	assert.True(t, VerifyPredicateProof(predicateProof))
}
