package claims

/*
// TMP commented until ClaimAuthorizeKSignSecp256k1 is updated to new spec
func TestClaimAuthorizeKSignSecp256k1(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := false
	// Init test
	if err := testgen.InitTest("claimAuthorizeKSignSecp256k1", generateTest); err != nil {
		panic(fmt.Errorf("error initializing test data: %w", err))
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("privateKey", "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f")
	}
	// ClaimAuthorizeKSignSecp256k1
	secKeyHex := testgen.GetTestValue("privateKey").(string)
	secKey, err := crypto.HexToECDSA(secKeyHex)
	if err != nil {
		panic(err)
	}
	pubKey := secKey.Public().(*ecdsa.PublicKey)
	c0 := NewClaimAuthorizeKSignSecp256k1(pubKey)
	c0.Version = 1
	e := c0.Entry()
	// Check claim against test vector
	testgen.CheckTestValue(t, "compressedPubKey", hex.EncodeToString(crypto.CompressPubkey(pubKey)))
	checkClaim(e, t)
	dataTestOutput(&e.Data)
	c1, err := NewClaimAuthorizeKSignSecp256k1FromEntry(e)
	if err != nil {
		panic(err)
	}
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
	// Stop test (write new test vector if needed)
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
}
*/
