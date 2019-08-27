package core

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

var rmDirs []string

func TestProof(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	rmDirs = append(rmDirs, dir)
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	id0, err := IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0, err := NewClaimSetRootKey(id0, rootKey0)
	assert.Nil(t, err)
	err = mt.Add(claim0.Entry())
	assert.Nil(t, err)

	id1, err := IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	assert.Nil(t, err)
	rootKey1 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b})
	claim1, err := NewClaimSetRootKey(id1, rootKey1)
	assert.Nil(t, err)
	err = mt.Add(claim1.Entry())
	assert.Nil(t, err)

	mtp, err := GetClaimProofByHi(mt, claim0.Entry().HIndex())
	assert.Nil(t, err)

	// j, err := json.Marshal(mtp)
	// assert.Nil(t, err)

	// id := ID{}
	verified, err := VerifyProofClaim(nil, mtp)
	assert.Nil(t, err)
	assert.True(t, verified)
}

func NewEntryFromInts(a, b, c, d int64) (e merkletree.Entry) {
	e.Data = IntsToData(a, b, c, d)
	return e
}

func IntsToData(_a, _b, _c, _d int64) merkletree.Data {
	a, b, c, d := big.NewInt(_a), big.NewInt(_b), big.NewInt(_c), big.NewInt(_d)
	return BigIntsToData(a, b, c, d)
}

func BigIntsToData(a, b, c, d *big.Int) (data merkletree.Data) {
	di := []*big.Int{a, b, c, d}
	for i, v := range di {
		copy(data[i][(merkletree.ElemBytesLen-len(v.Bytes())):], v.Bytes())
	}
	return
}

func TestClaimProof(t *testing.T) {
	mt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	claim1 := NewEntryFromInts(33, 44, 55, 66)
	claim2 := NewEntryFromInts(1111, 2222, 3333, 4444)
	claim3 := NewEntryFromInts(5555, 6666, 7777, 8888)

	mt.Add(&claim1)
	mt.Add(&claim2)
	mt.Add(&claim3)

	mtp, err := GetClaimProofByHi(mt, claim1.HIndex())
	assert.Nil(t, err)

	fmt.Println("mtp", mtp.Leaf,
		hex.EncodeToString(mtp.Proofs[0].Mtp0.Bytes()),
		hex.EncodeToString(mtp.Proofs[0].Mtp1.Bytes()))
}

func TestProofClaimGenesis(t *testing.T) {
	kOpStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	var kOp babyjub.PublicKey
	err := kOp.UnmarshalText([]byte(kOpStr))
	assert.Nil(t, err)

	claimKOp := NewClaimAuthorizeKSignBabyJub(&kOp).Entry()

	id, proofClaimKOp, err := CalculateIdGenesis(claimKOp, []*merkletree.Entry{})
	assert.Nil(t, err)

	proofClaimGenesis := ProofClaimGenesis{
		Claim: claimKOp,
		Mtp:   proofClaimKOp.Proofs[0].Mtp0,
		Root:  proofClaimKOp.Proofs[0].Root,
		Id:    id,
	}
	assert.Nil(t, proofClaimGenesis.Verify())

	// Invalid Id
	proofClaimGenesis = ProofClaimGenesis{
		Claim: claimKOp,
		Mtp:   proofClaimKOp.Proofs[0].Mtp0,
		Root:  proofClaimKOp.Proofs[0].Root,
		Id:    &ID{},
	}
	assert.NotNil(t, proofClaimGenesis.Verify())

	// Invalid Mtp of non-existence
	claimKOp2 := NewClaimAuthorizeKSignBabyJub(&kOp)
	claimKOp2.Version = 1
	proofClaimGenesis = ProofClaimGenesis{
		Claim: claimKOp2.Entry(),
		Mtp:   proofClaimKOp.Proofs[0].Mtp1,
		Root:  proofClaimKOp.Proofs[0].Root,
		Id:    &ID{},
	}
	assert.NotNil(t, proofClaimGenesis.Verify())

	// Invalid Claim
	proofClaimGenesis = ProofClaimGenesis{
		Claim: NewClaimBasic([50]byte{}, [62]byte{}).Entry(),
		Mtp:   proofClaimKOp.Proofs[0].Mtp0,
		Root:  proofClaimKOp.Proofs[0].Root,
		Id:    &ID{},
	}
	assert.NotNil(t, proofClaimGenesis.Verify())
}

func TestGetPredicateProof(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	rmDirs = append(rmDirs, dir)
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	id0, err := IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0, err := NewClaimSetRootKey(id0, rootKey0)
	assert.Nil(t, err)
	err = mt.Add(claim0.Entry())
	assert.Nil(t, err)
	oldRoot := mt.RootKey()

	mtp, err := GetClaimProofByHi(mt, claim0.Entry().HIndex())
	assert.Nil(t, err)

	// id := ID{}
	verified, err := VerifyProofClaim(nil, mtp)
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
	rmDirs = append(rmDirs, dir)
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	id0, err := IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0, err := NewClaimSetRootKey(id0, rootKey0)
	assert.Nil(t, err)
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
	rmDirs = append(rmDirs, dir)
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	id0, err := IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	assert.Nil(t, err)
	rootKey0 := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a})
	claim0, err := NewClaimSetRootKey(id0, rootKey0)
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
	assert.Equal(t, "0x1ca5dadd000f7e5874fe434888d916fbe81c9a3376675184dde0a213ab286d63", predicateProof.OldRoot.Hex())
	assert.NotEqual(t, predicateProof.OldRoot.Hex(), predicateProof.Root.Hex())

	assert.Equal(t, predicateProof.MtpNonExistInOldRoot.Siblings[0], predicateProof.MtpExist.Siblings[0])

	assert.True(t, VerifyPredicateProof(predicateProof))
}

func TestMain(m *testing.M) {
	result := m.Run()
	for _, dir := range rmDirs {
		os.RemoveAll(dir)
	}
	os.Exit(result)
}
