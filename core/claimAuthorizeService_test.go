package core

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

func TestClaimAuthorizeService(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := true
	// Init test
	err := testgen.InitTest("claimAuthorizeService", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("ethAddr", hex.EncodeToString([]byte{
			0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
			0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
			0x39, 0x39, 0x39, 0x3a}))
		testgen.SetTestValue("publicKey", "af048ddcc131d526699d928e8b8548c5c85fb7d407fc408bb543e4e58f305347f67942a7e56d7dc90bbcecca865f2fbde3118c91516594262f62857136f71dbc")
		testgen.SetTestValue("serviceName", "relay.iden3.io")
	}
	// ClaimAuthorizeService
	ethAddrHex, _ := hex.DecodeString(testgen.GetTestValue("ethAddr").(string))
	ethAddr := common.BytesToAddress(ethAddrHex)
	pubKstr := testgen.GetTestValue("publicKey").(string)
	c0 := NewClaimAuthorizeService(ServiceTypeRelay, ethAddr.Hex(), pubKstr, testgen.GetTestValue("serviceName").(string))
	e := c0.Entry()
	// Check claim against test vector
	checkClaim(e, t)
	dataTestOutput(&e.Data)
	c1 := NewClaimAuthorizeServiceFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
	assert.Equal(t, c0.ServiceType, ServiceTypeRelay)
	// Stop test (write new test vector if needed)
	err = testgen.StopTest()
	if err != nil {
		fmt.Println("Error stopping test:", err)
	}
}
