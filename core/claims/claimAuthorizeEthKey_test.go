package claims

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

func TestClaimAuthEthKey(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := false
	// Init test
	err := testgen.InitTest("claimAuthorizeEthKey", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("addr", "0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	}
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	ethKeyType := EthKeyTypeUpgrade

	c0 := NewClaimAuthEthKey(ethKey, ethKeyType)

	c1 := NewClaimAuthEthKeyFromEntry(c0.Entry())
	c2, err := NewClaimFromEntry(c0.Entry())
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)

	assert.Equal(t, c0.EthKey, ethKey)
	assert.Equal(t, c0.EthKeyType, binary.BigEndian.Uint32(ethKeyType[:]))
	assert.Equal(t, c0.EthKey, c1.EthKey)
	assert.Equal(t, c0.EthKeyType, c1.EthKeyType)
	assert.Equal(t, c0.Type(), c1.Type())
	assert.Equal(t, c0.Type(), *ClaimTypeAuthEthKey)

	assert.Equal(t, c0.Entry().Bytes(), c1.Entry().Bytes())
	assert.Equal(t, c0.Entry().Bytes(), c2.Entry().Bytes())

	e := c0.Entry()
	// Check claim against test vector
	checkClaim(e, t)
	dataTestOutput(&e.Data)
	// Stop test (write new test vector if needed)
	err = testgen.StopTest()
	if err != nil {
		fmt.Println("Error stopping test:", err)
	}
}
