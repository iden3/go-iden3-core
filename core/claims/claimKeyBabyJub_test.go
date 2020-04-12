package claims

import (
	"encoding/hex"
	"testing"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

func testClaimKeyBabyJub(t *testing.T, i, testKey string) {
	// Create new claim
	var k babyjub.PrivateKey
	hexK := testgen.GetTestValue(i + testKey).(string)
	if _, err := hex.Decode(k[:], []byte(hexK)); err != nil {
		panic(err)
	}
	pk := k.Public()
	c0 := NewClaimKeyBabyJub(pk)
	c0.Metadata().RevNonce = 5678
	assert.True(t, merkletree.CheckEntryInField(*c0.Entry()))
	e := c0.Entry()
	// Check claim against test vector
	hi, hv, err := e.HiHv()
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "ClaimKeyBabyJub"+i+"_HIndex", hi.Hex())
	testgen.CheckTestValue(t, "ClaimKeyBabyJub"+i+"_HValue", hv.Hex())
	testgen.CheckTestValue(t, "ClaimKeyBabyJub"+i+"_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimKeyBabyJubFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0.Metadata(), c1.Metadata())
	assert.Equal(t, c0, c2)
	assert.True(t, merkletree.CheckEntryInField(*e))
}

func TestClaimKeyBabyJub(t *testing.T) {
	testClaimKeyBabyJub(t, "0", "_privateKey")
	testClaimKeyBabyJub(t, "1", "_privateKey")
}

func TestRandomClaimKeyBabyJub(t *testing.T) {
	for i := 0; i < 100; i++ {
		k := babyjub.NewRandPrivKey()
		pk := k.Public()

		c0 := NewClaimKeyBabyJub(pk)
		assert.True(t, merkletree.CheckEntryInField(*c0.Entry()))
		e := c0.Entry()
		c1 := NewClaimKeyBabyJubFromEntry(e)
		c2, err := NewClaimFromEntry(e)
		assert.Nil(t, err)
		assert.Equal(t, c0, c1)
		assert.Equal(t, c0, c2)
		assert.True(t, merkletree.CheckEntryInField(*e))
	}
}
