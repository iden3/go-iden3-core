package claims

/*
// TMP commented until ClaimAssignName is updated to new spec

func TestClaimAssignName(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := false

	// Init test
	if err := testgen.InitTest("claimAssignName", generateTest); err != nil {
		panic(fmt.Errorf("error initializing test data: %w", err))
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("name", "example.iden3.eth")
		testgen.SetTestValue("IDString", "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	}
	// Get input data from test vector
	name := testgen.GetTestValue("name").(string)
	id, err := core.IDFromString(testgen.GetTestValue("IDString").(string))
	assert.Nil(t, err)
	// Create new claim
	c0 := NewClaimAssignName(name, id)
	c0.Version = 1
	e := c0.Entry()
	// Check claim against test vector
	checkClaim(e, t)
	dataTestOutput(&e.Data)
	c1 := NewClaimAssignNameFromEntry(e)
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
