package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/iden3/go-iden3-core/crypto"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

var generateTest = false

func TestIDparsers(t *testing.T) {
	// Generate ID0
	var typ0 [2]byte
	typ0Hex, _ := hex.DecodeString(testgen.GetTestValue("typ0").(string))
	copy(typ0[:], typ0Hex[:2])
	var genesis0 [27]byte
	genesis032bytes := crypto.HashBytes([]byte(testgen.GetTestValue("genesisUnhashedString0").(string)))
	copy(genesis0[:], genesis032bytes[:])
	id0 := NewID(typ0, genesis0)
	// Check ID0
	testgen.CheckTestValue(t, "idString0", id0.String())
	// Generate ID1
	var typ1 [2]byte
	typ1Hex, _ := hex.DecodeString(testgen.GetTestValue("typ1").(string))
	copy(typ1[:], typ1Hex[:2])
	var genesis1 [27]byte
	genesis132bytes := crypto.HashBytes([]byte(testgen.GetTestValue("genesisUnhashedString0").(string)))
	copy(genesis1[:], genesis132bytes[:])
	id1 := NewID(typ1, genesis1)
	// Check ID1
	testgen.CheckTestValue(t, "idString1", id1.String())

	emptyChecksum := []byte{0x00, 0x00}
	assert.True(t, !bytes.Equal(emptyChecksum, id0[29:]))
	assert.True(t, !bytes.Equal(emptyChecksum, id1[29:]))

	id0FromBytes, err := IDFromBytes(id0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromBytes.Bytes())
	assert.Equal(t, id0.String(), id0FromBytes.String())
	testgen.CheckTestValue(t, "idString0", id0FromBytes.String())

	id1FromBytes, err := IDFromBytes(id1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id1.Bytes(), id1FromBytes.Bytes())
	assert.Equal(t, id1.String(), id1FromBytes.String())
	testgen.CheckTestValue(t, "idString1", id1FromBytes.String())

	id0FromString, err := IDFromString(id0.String())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromString.Bytes())
	assert.Equal(t, id0.String(), id0FromString.String())
	testgen.CheckTestValue(t, "idString0", id0FromString.String())
}

func TestIDjsonParser(t *testing.T) {
	id, err := IDFromString(testgen.GetTestValue("idStringInput").(string))
	assert.Nil(t, err)
	idj, err := json.Marshal(&id)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "idString2", strings.Replace(string(idj), "\"", "", 2))
	var idp ID
	err = json.Unmarshal(idj, &idp)
	assert.Nil(t, err)

	assert.Equal(t, id, idp)

	idsMap := make(map[ID]string)
	idsMap[id] = "first"
	idsMapJSON, err := json.Marshal(idsMap)
	assert.Nil(t, err)

	var idsMapUnmarshaled map[ID]string
	err = json.Unmarshal(idsMapJSON, &idsMapUnmarshaled)
	assert.Nil(t, err)
}

func TestCheckChecksum(t *testing.T) {
	typ := TypeBJP0
	var genesis [27]byte
	genesis32bytes := crypto.HashBytes([]byte(testgen.GetTestValue("genesisUnhashedString0").(string)))
	copy(genesis[:], genesis32bytes[:])

	id := NewID(typ, genesis)

	var checksum [2]byte
	copy(checksum[:], id[len(id)-2:])
	assert.Equal(t, CalculateChecksum(typ, genesis), checksum)

	assert.True(t, CheckChecksum(id))

	// check that if we change the checksum, returns false on CheckChecksum
	id = NewID(typ, genesis)
	copy(id[29:], []byte{0x00, 0x01})
	assert.True(t, !CheckChecksum(id))

	// check that if we change the type, returns false on CheckChecksum
	id = NewID(typ, genesis)
	copy(id[:2], []byte{0x00, 0x01})
	assert.True(t, !CheckChecksum(id))

	// check that if we change the genesis, returns false on CheckChecksum
	id = NewID(typ, genesis)
	// changedGenesis := utils.HashBytes([]byte("changedgenesis"))
	var changedGenesis [27]byte
	changedGenesis32bytes := crypto.HashBytes([]byte(testgen.GetTestValue("genesisUnhashedString1").(string)))
	copy(changedGenesis[:], changedGenesis32bytes[:27])

	copy(id[2:27], changedGenesis[:])
	assert.True(t, !CheckChecksum(id))

	// test with a empty id
	var empty [31]byte
	_, err := IDFromBytes(empty[:])
	assert.Equal(t, errors.New("IDFromBytes error: byte array empty"), err)
}

func initTest() {
	// If generateTest is true, the checked values will be used to generate a test vector
	// Init test
	err := testgen.InitTest("id", generateTest)
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
}

func TestMain(m *testing.M) {
	initTest()
	result := m.Run()
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
	os.Exit(result)
}
