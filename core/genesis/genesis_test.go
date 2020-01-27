package genesis

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

var debug = false

// WARNING:	all the functions must be executed when tested.º
// First test function to be executed must call initializeTest
// First test function to be executed must call finalizeTest

// Avoids reinitializing tests
var testInitialized = false

func initializeTest() {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := false
	if !testInitialized {
		// Init test
		err := testgen.InitTest("genesis", generateTest)
		if err != nil {
			fmt.Println("error initializing test data:", err)
			return
		}
		// Add input data to the test vector
		if generateTest {
			testgen.SetTestValue("genesisUnhashedString0", "genesistest")
			testgen.SetTestValue("genesisUnhashedString1", "changedgenesis")
			testgen.SetTestValue("typ0", hex.EncodeToString([]byte{0x00, 0x00}))
			testgen.SetTestValue("typ1", hex.EncodeToString([]byte{0x00, 0x01}))
			testgen.SetTestValue("babyJub", "28156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f")
			testgen.SetTestValue("addr", "0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
			testgen.SetTestValue("idStringInput", "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf")
			testgen.SetTestValue("kOp", "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796")
		}
		testInitialized = true
	}
}

func finalizeTest() {
	// Stop test (write new test vector if needed)
	err := testgen.StopTest()
	if err != nil {
		fmt.Println("Error stopping test:", err)
	}
}

func TestCalculateIdGenesisFrom4Keys(t *testing.T) {
	initializeTest()
	var sk babyjub.PrivateKey
	hex.Decode(sk[:], []byte(testgen.GetTestValue("babyJub").(string)))
	kopPub := sk.Public()
	kDis := common.HexToAddress(testgen.GetTestValue("addr").(string))
	kReen := kDis
	kUpdateRoot := kDis

	id, _, err := CalculateIdGenesisFrom4Keys(kopPub, kDis, kReen, kUpdateRoot)
	assert.Nil(t, err)
	if debug {
		fmt.Println("id", id)
		fmt.Println("id (hex)", id.String())
	}
	testgen.CheckTestValue("idString3", id.String(), t)
}

func TestCalculateIdGenesis(t *testing.T) {
	kopStr := testgen.GetTestValue("kOp").(string)
	var kopComp babyjub.PublicKeyComp
	err := kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	claimKOp := claims.NewClaimAuthorizeKSignBabyJub(kopPub)

	id, _, err := CalculateIdGenesis(claimKOp, []*merkletree.Entry{})
	assert.Nil(t, err)
	if debug {
		fmt.Println("id", id)
		fmt.Println("id (hex)", id.String())
	}
	testgen.CheckTestValue("idString4", id.String(), t)
	finalizeTest()
}

// TODO: Review if this goes here or in proof
func TestProofClaimGenesis(t *testing.T) {
	kOpStr := testgen.GetTestValue("kOp").(string)
	var kOp babyjub.PublicKey
	err := kOp.UnmarshalText([]byte(kOpStr))
	assert.Nil(t, err)

	claimKOp := claims.NewClaimAuthorizeKSignBabyJub(&kOp)

	id, proofClaimKOp, err := CalculateIdGenesis(claimKOp, []*merkletree.Entry{})
	assert.Nil(t, err)

	proofClaimGenesis := proof.ProofClaimGenesis{
		Mtp: proofClaimKOp.Proof.Mtp0,
		Id:  id,
	}
	_, err = proofClaimGenesis.Verify(claimKOp.Entry())
	assert.Nil(t, err)

	// Invalid Id
	proofClaimGenesis = proof.ProofClaimGenesis{
		Mtp: proofClaimKOp.Proof.Mtp0,
		Id:  &core.ID{},
	}
	_, err = proofClaimGenesis.Verify(claimKOp.Entry())
	assert.NotNil(t, err)

	// Invalid Mtp of non-existence
	claimKOp2 := claims.NewClaimAuthorizeKSignBabyJub(&kOp)
	claimKOp2.Version = 1
	proofClaimGenesis = proof.ProofClaimGenesis{
		Mtp: proofClaimKOp.Proof.Mtp1,
		Id:  &core.ID{},
	}
	_, err = proofClaimGenesis.Verify(claimKOp2.Entry())
	assert.NotNil(t, err)

	// Invalid Claim
	proofClaimGenesis = proof.ProofClaimGenesis{
		Mtp: proofClaimKOp.Proof.Mtp0,
		Id:  &core.ID{},
	}
	_, err = proofClaimGenesis.Verify(claims.NewClaimBasic([50]byte{}, [62]byte{}).Entry())
	assert.NotNil(t, err)
}
