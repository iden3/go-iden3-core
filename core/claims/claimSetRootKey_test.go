package claims

/*
// TMP commented until ClaimSetRootKey is updated to new spec
func TestClaimSetRootKey(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := false
	// Init test
	if err := testgen.InitTest("claimSetRootKey", generateTest); err != nil {
		panic(fmt.Errorf("error initializing test data: %w", err))
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("idString", "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
		testgen.SetTestValue("rootKey", hex.EncodeToString([]byte{
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c}))
	}
	// ClaimSetRootKey
	id, err := core.IDFromString(testgen.GetTestValue("idString").(string))
	assert.Nil(t, err)
	rootKey := hexStringToKey(testgen.GetTestValue("rootKey").(string))
	c0, err := NewClaimSetRootKey(&id, &rootKey)
	assert.Nil(t, err)
	c0.Version = 1
	c0.Era = 1
	e := c0.Entry()
	// Check claim against test vector
	checkClaim(e, t)
	dataTestOutput(&e.Data)
	c1 := NewClaimSetRootKeyFromEntry(e)
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
