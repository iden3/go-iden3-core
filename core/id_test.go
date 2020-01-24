package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-iden3-core/utils"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

// WARNING:	all the functions must be executed when tested.º
// First test function to be executed must call initializeTest
// First test function to be executed must call finalizeTest

// Avoids reinitializing tests
var testInitialized = false

func initializeTest() {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := true
	if !testInitialized {
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
		testInitialized = true
	}
}

func finalizeTest() {
	// Stop test (write new test vector if needed)
	err := testgen.StopTest()
	if err != nil {
		fmt.Println("Error stopping test:", err)
	}
}

func TestIDparsers(t *testing.T) {
	initializeTest()
	// Generate ID0
	var typ0 [2]byte
	typ0Hex, _ := hex.DecodeString(testgen.GetTestValue("typ0").(string))
	copy(typ0[:], typ0Hex[:2])
	var genesis0 [27]byte
	genesis032bytes := utils.HashBytes([]byte(testgen.GetTestValue("genesisUnhashedString0").(string)))
	copy(genesis0[:], genesis032bytes[:])
	id0 := NewID(typ0, genesis0)
	// Check ID0
	testgen.CheckTestValue("idString0", id0.String(), t)
	// Generate ID1
	var typ1 [2]byte
	typ1Hex, _ := hex.DecodeString(testgen.GetTestValue("typ1").(string))
	copy(typ1[:], typ1Hex[:2])
	var genesis1 [27]byte
	genesis132bytes := utils.HashBytes([]byte(testgen.GetTestValue("genesisUnhashedString0").(string)))
	copy(genesis1[:], genesis132bytes[:])
	id1 := NewID(typ1, genesis1)
	// Check ID1
	testgen.CheckTestValue("idString1", id1.String(), t)

	emptyChecksum := []byte{0x00, 0x00}
	assert.True(t, !bytes.Equal(emptyChecksum, id0[29:]))
	assert.True(t, !bytes.Equal(emptyChecksum, id1[29:]))

	id0FromBytes, err := IDFromBytes(id0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromBytes.Bytes())
	assert.Equal(t, id0.String(), id0FromBytes.String())
	testgen.CheckTestValue("idString0", id0FromBytes.String(), t)

	id1FromBytes, err := IDFromBytes(id1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id1.Bytes(), id1FromBytes.Bytes())
	assert.Equal(t, id1.String(), id1FromBytes.String())
	testgen.CheckTestValue("idString1", id1FromBytes.String(), t)

	id0FromString, err := IDFromString(id0.String())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromString.Bytes())
	assert.Equal(t, id0.String(), id0FromString.String())
	testgen.CheckTestValue("idString0", id0FromString.String(), t)
}

func TestIDjsonParser(t *testing.T) {
	id, err := IDFromString(testgen.GetTestValue("idStringInput").(string))
	assert.Nil(t, err)
	idj, err := json.Marshal(&id)
	assert.Nil(t, err)
	testgen.CheckTestValue("idString2", strings.Replace(string(idj), "\"", "", 2), t)
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
	genesis32bytes := utils.HashBytes([]byte(testgen.GetTestValue("genesisUnhashedString0").(string)))
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
	changedGenesis32bytes := utils.HashBytes([]byte(testgen.GetTestValue("genesisUnhashedString1").(string)))
	copy(changedGenesis[:], changedGenesis32bytes[:27])

	copy(id[2:27], changedGenesis[:])
	assert.True(t, !CheckChecksum(id))

	// test with a empty id
	var empty [31]byte
	_, err := IDFromBytes(empty[:])
	assert.Equal(t, errors.New("IDFromBytes error: byte array empty"), err)
}

func TestCalculateIdGenesisFrom4Keys(t *testing.T) {
	var sk babyjub.PrivateKey
	hex.Decode(sk[:], []byte(testgen.GetTestValue("babyJub").(string)))
	kopPub := sk.Public()
	kDis := common.HexToAddress(testgen.GetTestValue("addr").(string))
	kReen := kDis
	kUpdateRoot := kDis

	id, _, err := CalculateIdGenesisFrom4Keys(kopPub, kDis, kReen, kUpdateRoot)
	assert.Nil(t, err)
	if debug {
		fmt.Println("id", id)
		fmt.Println("id (hex)", id.String())
	}
	testgen.CheckTestValue("idString3", id.String(), t)
}

func TestCalculateIdGenesis(t *testing.T) {
	kopStr := testgen.GetTestValue("kOp").(string)
	var kopComp babyjub.PublicKeyComp
	err := kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	claimKOp := NewClaimAuthorizeKSignBabyJub(kopPub)

	id, _, err := CalculateIdGenesis(claimKOp, []*merkletree.Entry{})
	assert.Nil(t, err)
	if debug {
		fmt.Println("id", id)
		fmt.Println("id (hex)", id.String())
	}
	testgen.CheckTestValue("idString4", id.String(), t)
	finalizeTest()
}
