package claims

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

func testClaimAuthorizeKSignBabyJub(t *testing.T, i, testKey string) {
	// Create new claim
	var k babyjub.PrivateKey
	hexK := testgen.GetTestValue(i + testKey).(string)
	if _, err := hex.Decode(k[:], []byte(hexK)); err != nil {
		panic(err)
	}
	pk := k.Public()
	c0 := NewClaimAuthorizeKSignBabyJub(pk)
	assert.True(t, merkletree.CheckEntryInField(*c0.Entry()))
	c0.Version = 1
	e := c0.Entry()
	// Check claim against test vector
	testgen.CheckTestValue(t, i+"_HIndex", e.HIndex().Hex())
	testgen.CheckTestValue(t, i+"_HValue", e.HValue().Hex())
	testgen.CheckTestValue(t, i+"_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimAuthorizeKSignBabyJubFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
	assert.True(t, merkletree.CheckEntryInField(*e))
}

func TestClaimAuthorizeKSignBabyJub(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := false
	// Init test
	if err := testgen.InitTest("claimAuthorizeKSignBabyJub", generateTest); err != nil {
		panic(fmt.Errorf("error initializing test data: %w", err))
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("0_privateKey", "28156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f")
		testgen.SetTestValue("1_privateKey", "9b3260823e7b07dd26ef357ccfed23c10bcef1c85940baa3d02bbf29461bbbbe")
	}
	testClaimAuthorizeKSignBabyJub(t, "0", "_privateKey")
	testClaimAuthorizeKSignBabyJub(t, "1", "_privateKey")

	// Stop test (write new test vector if needed)
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
}

func TestRandomClaimAuthorizeKSignBabyJub(t *testing.T) {
	for i := 0; i < 100; i++ {
		k := babyjub.NewRandPrivKey()
		pk := k.Public()

		c0 := NewClaimAuthorizeKSignBabyJub(pk)
		assert.True(t, merkletree.CheckEntryInField(*c0.Entry()))
		c0.Version = 9999
		e := c0.Entry()
		c1 := NewClaimAuthorizeKSignBabyJubFromEntry(e)
		c2, err := NewClaimFromEntry(e)
		assert.Nil(t, err)
		assert.Equal(t, c0, c1)
		assert.Equal(t, c0, c2)
		assert.True(t, merkletree.CheckEntryInField(*e))
	}
}
