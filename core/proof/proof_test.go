package proof

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

// If generateTest is true, the checked values will be used to generate a test vector
var generateTest = false

var rmDirs []string

/*
// TMP commented due ClaimSetRootKey is not updated to new spec
func TestProof(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	rmDirs = append(rmDirs, dir)
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)
	id0, err := core.IDFromString(testgen.GetTestValue("idString0").(string))
	assert.Nil(t, err)

	rootKey0 := hexStringToKey(testgen.GetTestValue("rootKey0").(string))

	claim0, err := claims.NewClaimSetRootKey(&id0, &rootKey0)
	assert.Nil(t, err)
	err = mt.AddClaim(claim0)
	assert.Nil(t, err)

	//idString
	id1, err := core.IDFromString(testgen.GetTestValue("idString1").(string))
	assert.Nil(t, err)
	rootKey1 := hexStringToKey(testgen.GetTestValue("rootKey1").(string))
	claim1, err := claims.NewClaimSetRootKey(&id1, &rootKey1)
	assert.Nil(t, err)
	err = mt.AddClaim(claim1)
	assert.Nil(t, err)

	cp, err := GetClaimProofByHi(mt, claim0.Entry().HIndex())
	assert.Nil(t, err)

	verified, err := cp.Verify(cp.Proof.Root)
	assert.Nil(t, err)
	assert.True(t, verified)
}
*/

func TestClaimProof(t *testing.T) {
	mt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	entry1 := merkletree.NewEntryFromInts(33, 44, 55, 66, 11, 22, 33, 44)
	entry2 := merkletree.NewEntryFromInts(5, 2222, 3333, 4444, 1, 2, 3, 4)
	entry3 := merkletree.NewEntryFromInts(5555, 6666, 7777, 8888, 1, 2, 3, 4)

	if err := mt.AddEntry(&entry1); err != nil {
		panic(err)
	}
	if err := mt.AddEntry(&entry2); err != nil {
		panic(err)
	}
	if err := mt.AddEntry(&entry3); err != nil {
		panic(err)
	}

	mtp, err := GetClaimProofByHi(mt, entry1.HIndex())
	assert.Nil(t, err)

	fmt.Println("mtp", mtp.Claim,
		hex.EncodeToString(mtp.Proof.Mtp0.Bytes()),
		hex.EncodeToString(mtp.Proof.Mtp1.Bytes()))
}

func TestGetPredicateProof(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	rmDirs = append(rmDirs, dir)
	assert.Nil(t, err)
	sto, err := db.NewLevelDbStorage(dir, false)
	assert.Nil(t, err)

	mt, err := merkletree.NewMerkleTree(sto, 140)
	assert.Nil(t, err)

	id0, err := core.IDFromString(testgen.GetTestValue("idString1").(string))
	assert.Nil(t, err)
	rootKey0 := merkletree.HexStringToHash(testgen.GetTestValue("rootKey2").(string))
	claim0, err := claims.NewClaimSetRootKey(&id0, &rootKey0)
	assert.Nil(t, err)
	err = mt.AddClaim(claim0)
	assert.Nil(t, err)
	oldRoot := mt.RootKey()

	cp, err := GetClaimProofByHi(mt, claim0.Entry().HIndex())
	assert.Nil(t, err)

	// id := core.ID{}
	verified, err := cp.Verify(cp.Proof.Root)
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

	id0, err := core.IDFromString(testgen.GetTestValue("idString1").(string))
	assert.Nil(t, err)
	rootKey0 := merkletree.HexStringToHash(testgen.GetTestValue("rootKey2").(string))
	claim0, err := claims.NewClaimSetRootKey(&id0, &rootKey0)
	assert.Nil(t, err)
	// oldRoot is the root before adding the claim that we want to prove that we added correctly
	oldRoot := mt.RootKey()
	err = mt.AddClaim(claim0)
	assert.Nil(t, err)

	predicateProof, err := GetPredicateProof(mt, oldRoot, claim0.Entry().HIndex())
	assert.Nil(t, err)

	_, v := claims.GetClaimTypeVersion(predicateProof.LeafEntry)
	assert.Equal(t, uint32(0), v)
	testgen.CheckTestValue(t, "predicateProof0", predicateProof.OldRoot.Hex())
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

	id0, err := core.IDFromString(testgen.GetTestValue("idString1").(string))
	assert.Nil(t, err)
	rootKey0 := merkletree.HexStringToHash(testgen.GetTestValue("rootKey2").(string))
	claim0, err := claims.NewClaimSetRootKey(&id0, &rootKey0)
	if err != nil {
		panic(err)
	}
	err = mt.AddClaim(claim0)
	assert.Nil(t, err)

	// add the same claim, but with version 1
	claim1 := &claims.ClaimSetRootKey{
		Version: claim0.Version + 1,
		Era:     0,
		Id:      claim0.Id,
		RootKey: claim0.RootKey,
	}
	err = mt.AddClaim(claim1)
	assert.Nil(t, err)
	// add the same claim, but with version 2
	claim2 := &claims.ClaimSetRootKey{
		Version: claim1.Version + 1,
		Era:     0,
		Id:      claim0.Id,
		RootKey: claim0.RootKey,
	}
	err = mt.AddClaim(claim2)
	assert.Nil(t, err)

	// oldRoot is the root before adding the claim that we want to prove that we added correctly
	oldRoot := mt.RootKey()

	// add the same claim, but with version 3
	claim3 := &claims.ClaimSetRootKey{
		Version: claim2.Version + 1,
		Era:     0,
		Id:      claim0.Id,
		RootKey: claim0.RootKey,
	}
	err = mt.AddClaim(claim3)
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

	_, v := claims.GetClaimTypeVersion(predicateProof.LeafEntry)
	assert.Equal(t, uint32(3), v)
	testgen.CheckTestValue(t, "predicateProof1", predicateProof.OldRoot.Hex())
	assert.NotEqual(t, predicateProof.OldRoot.Hex(), predicateProof.Root.Hex())

	assert.Equal(t, predicateProof.MtpNonExistInOldRoot.Siblings[0], predicateProof.MtpExist.Siblings[0])

	assert.True(t, VerifyPredicateProof(predicateProof))
}

func initTest() {
	// Init test
	err := testgen.InitTest("proof", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("idString0", "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf")
		testgen.SetTestValue("idString1", "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
		testgen.SetTestValue("rootKey0", hex.EncodeToString([]byte{
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c}))
		testgen.SetTestValue("rootKey1", hex.EncodeToString([]byte{
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b}))
		testgen.SetTestValue("rootKey2", hex.EncodeToString([]byte{
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0a}))
		testgen.SetTestValue("kOp", "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796")
	}
}

func TestMain(m *testing.M) {
	initTest()
	result := m.Run()
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
	for _, dir := range rmDirs {
		os.RemoveAll(dir)
	}
	os.Exit(result)
}
