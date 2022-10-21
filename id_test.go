package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Hash [32]byte

// Hex returns a hex string from the Hash type
func (hash Hash) Hex() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(hash[:]))
}

// hashBytes performs a sha256 hash over the bytes
func hashBytes(b []byte) (hash Hash) {
	h := sha256.Sum256(b)
	copy(hash[:], h[:])
	return hash
}

func TestIDparsers(t *testing.T) {
	// Generate ID0
	var typ0 [2]byte
	typ0Hex, _ := hex.DecodeString("0000")
	copy(typ0[:], typ0Hex[:2])
	var genesis0 [27]byte
	genesis032bytes := hashBytes([]byte("genesistest"))
	copy(genesis0[:], genesis032bytes[:])
	id0 := NewID(typ0, genesis0)
	// Check ID0
	assert.Equal(t, "114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E", id0.String())
	// Generate ID1
	var typ1 [2]byte
	typ1Hex, _ := hex.DecodeString("0001")
	copy(typ1[:], typ1Hex[:2])
	var genesis1 [27]byte
	genesis132bytes := hashBytes([]byte("genesistest"))
	copy(genesis1[:], genesis132bytes[:])
	id1 := NewID(typ1, genesis1)
	// Check ID1
	assert.Equal(t, "1GYjyJKqdDyzo927FqJkAdLWB64kV2NVAjaQFHtq4", id1.String())

	emptyChecksum := []byte{0x00, 0x00}
	assert.True(t, !bytes.Equal(emptyChecksum, id0[29:]))
	assert.True(t, !bytes.Equal(emptyChecksum, id1[29:]))

	id0FromBytes, err := IDFromBytes(id0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromBytes.Bytes())
	assert.Equal(t, id0.String(), id0FromBytes.String())
	assert.Equal(t, "114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
		id0FromBytes.String())

	id1FromBytes, err := IDFromBytes(id1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id1.Bytes(), id1FromBytes.Bytes())
	assert.Equal(t, id1.String(), id1FromBytes.String())
	assert.Equal(t, "1GYjyJKqdDyzo927FqJkAdLWB64kV2NVAjaQFHtq4",
		id1FromBytes.String())

	id0FromString, err := IDFromString(id0.String())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromString.Bytes())
	assert.Equal(t, id0.String(), id0FromString.String())
	assert.Equal(t, "114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
		id0FromString.String())
}

func TestIDjsonParser(t *testing.T) {
	id, err := IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZogFv")
	assert.Nil(t, err)
	idj, err := json.Marshal(&id)
	assert.Nil(t, err)
	assert.Equal(t, "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZogFv",
		strings.Replace(string(idj), "\"", "", 2))
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
	typ := TypeDefault
	var genesis [27]byte
	genesis32bytes := hashBytes([]byte("genesistest"))
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
	changedGenesis32bytes := hashBytes([]byte("changedgenesis"))
	copy(changedGenesis[:], changedGenesis32bytes[:27])

	copy(id[2:27], changedGenesis[:])
	assert.True(t, !CheckChecksum(id))

	// test with a empty id
	var empty [31]byte
	_, err := IDFromBytes(empty[:])
	assert.Equal(t, errors.New("IDFromBytes error: byte array empty"), err)
}

func TestIDFromInt(t *testing.T) {
	id, err := IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZogFv")
	assert.Nil(t, err)

	intID := id.BigInt()

	got, err := IDFromInt(intID)
	assert.Nil(t, err)

	assert.Equal(t, id, got)
}

func TestIDFromIntStr(t *testing.T) {
	idStr := "11BBCPZ6Zq9HX1JhHrHT3QKUFD9kFDEyJFoAVMptVs"

	idFromStr, err := IDFromString(idStr)
	require.NoError(t, err)

	intFromIDFromStr := idFromStr.BigInt()

	id, err := IDFromInt(intFromIDFromStr)
	require.NoError(t, err)

	require.Equal(t, idStr, id.String())
}

func TestProfileID(t *testing.T) {
	idInt, ok := new(big.Int).SetString(
		"23630567111950550539435915649280822148510307443797111728722609533581131776",
		10)
	require.True(t, ok)
	id, err := IDFromInt(idInt)
	require.NoError(t, err)
	nonce := big.NewInt(10)
	id2, err := ProfileID(id, nonce)
	require.NoError(t, err)
	require.Equal(t,
		"25425363284463910957419549722021124450832239517990785975889689633068548096",
		id2.BigInt().String())
}
