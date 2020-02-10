package testgen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Input  map[string]interface{}
	Output map[string]interface{}
}

var generate bool
var fileName string
var testData TestData

// InitTest initializes the testgen framework.
func InitTest(name string, gen bool) error {
	filePath := "testVectors"
	generate = gen
	err := os.MkdirAll(filePath, 0744)
	if err != nil {
		return err
	}
	fileName = path.Join(filePath, name+".json")
	// file doesnt exist yet!
	if generate {
		testData.Input = make(map[string]interface{})
		testData.Output = make(map[string]interface{})
		return nil
	} else if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return errors.New("No test vector has been created yet and won't be create since gen is set to false")
	} else {
		// load test data from file
		testData, err = getTestData()
		return err
	}
}

// CheckTestValue will check that the passed value is equal to the testVectors
// value under the specified output value key.  If generate is set to true,
// instead of checking the equality, the value will be stored in the testVector
// under the specified output value key.
func CheckTestValue(t *testing.T, key string, value interface{}) {
	if generate {
		if val, ok := testData.Output[key]; ok && val != value {
			panic(fmt.Sprint("Key already used with different value: \nKey:", key,
				"\nPrevious value:", val,
				"\nCurrent value:", value))
		}
		testData.Output[key] = value
	} else {
		assert.Equal(t, testData.Output[key], value)
	}
}

// GetTestValue takes the value from the testVectors under the specified input
// value key.
func GetTestValue(key string) interface{} {
	return testData.Input[key]
}

// SetTestValue sets the value to the testVectors under the specified input
// value key.
func SetTestValue(key string, value interface{}) {
	if generate {
		testData.Input[key] = value
	}
}

// StopTest will write the testVectors to the corresponding file if generate is
// true.
func StopTest() error {
	if generate {
		return writeGeneratedTest()
	}
	return nil
}

func getTestData() (TestData, error) {
	// Read file
	var td TestData
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return td, err
	}
	defer jsonFile.Close()

	// Decode file
	dec := json.NewDecoder(jsonFile)
	err = dec.Decode(&td)
	if err != nil {
		return td, err
	}

	return td, nil
}

func writeGeneratedTest() error {
	// write genrated test data into a json file
	j, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, j, 0644)
	return err
}
