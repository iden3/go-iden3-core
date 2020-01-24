package core

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

func TestClaimAuthorizeKSignBabyJub(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := true
	// Init test
	err := testgen.InitTest("claimAuthorizeKSignBabyJub", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("privateKey", "28156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f")
	}
	// Create new claim
	var k babyjub.PrivateKey
	hexK := testgen.GetTestValue("privateKey").(string)
	hex.Decode(k[:], []byte(hexK))
	pk := k.Public()
	c0 := NewClaimAuthorizeKSignBabyJub(pk)
	c0.Version = 1
	e := c0.Entry()
	// Check claim against test vector
	checkClaim(e, t)
	dataTestOutput(&e.Data)
	c1 := NewClaimAuthorizeKSignBabyJubFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
	// Stop test (write new test vector if needed)
	err = testgen.StopTest()
	if err != nil {
		fmt.Println("Error stopping test:", err)
	}
}
