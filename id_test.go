package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	assert.Equal(t, "114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE", id0.String())
	// Generate ID1
	var typ1 [2]byte
	typ1Hex, _ := hex.DecodeString("0001")
	copy(typ1[:], typ1Hex[:2])
	var genesis1 [27]byte
	genesis132bytes := hashBytes([]byte("genesistest"))
	copy(genesis1[:], genesis132bytes[:])
	id1 := NewID(typ1, genesis1)
	// Check ID1
	assert.Equal(t, "1GYjyJKqdDyzo927FqJkAdLWB64kV2NVAjaQFHbAf", id1.String())

	emptyChecksum := []byte{0x00, 0x00}
	assert.True(t, !bytes.Equal(emptyChecksum, id0[29:]))
	assert.True(t, !bytes.Equal(emptyChecksum, id1[29:]))

	id0FromBytes, err := IDFromBytes(id0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromBytes.Bytes())
	assert.Equal(t, id0.String(), id0FromBytes.String())
	assert.Equal(t, "114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE",
		id0FromBytes.String())

	id1FromBytes, err := IDFromBytes(id1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id1.Bytes(), id1FromBytes.Bytes())
	assert.Equal(t, id1.String(), id1FromBytes.String())
	assert.Equal(t, "1GYjyJKqdDyzo927FqJkAdLWB64kV2NVAjaQFHbAf",
		id1FromBytes.String())

	id0FromString, err := IDFromString(id0.String())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromString.Bytes())
	assert.Equal(t, id0.String(), id0FromString.String())
	assert.Equal(t, "114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE",
		id0FromString.String())
}

func TestIDAsDID(t *testing.T) {
	typ := [2]byte{DIDIden3BlockchainType[Polygon], DIDNetworkType[Mumbai]}
	var genesis1 [27]byte
	genesisbytes := hashBytes([]byte("genesistes1t2"))
	copy(genesis1[:], genesisbytes[:])

	id := NewID(typ, genesis1)
	fmt.Println(id.String())
}

func TestIDjsonParser(t *testing.T) {
	id, err := IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf")
	assert.Nil(t, err)
	idj, err := json.Marshal(&id)
	assert.Nil(t, err)
	assert.Equal(t, "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf",
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
	id, err := IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf")
	assert.Nil(t, err)

	intID := id.BigInt()

	got, err := IDFromInt(intID)
	assert.Nil(t, err)

	assert.Equal(t, id, got)
}

func TestIDFromIntStr(t *testing.T) {
	idStr := "11BBCPZ6Zq9HX1JhHrHT3QKUFD9kFDEyJFoAVMpuZR"

	idFromStr, err := IDFromString(idStr)
	require.NoError(t, err)

	intFromIDFromStr := idFromStr.BigInt()

	id, err := IDFromInt(intFromIDFromStr)
	require.NoError(t, err)

	require.Equal(t, idStr, id.String())
}

func TestIDinDIDFormat(t *testing.T) {

	typ := [2]byte{DIDIden3BlockchainType[Polygon], DIDNetworkType[Mumbai]}
	var genesis [27]byte
	genesis32bytes := hashBytes([]byte("genesistest"))
	copy(genesis[:], genesis32bytes[:])

	id := NewID(typ, genesis)

	var checksum [2]byte
	copy(checksum[:], id[len(id)-2:])
	assert.Equal(t, CalculateChecksum(typ, genesis), checksum)

	fmt.Println(id.String())
	did := DID{
		ID:         id,
		Blockchain: Polygon,
		NetworkID:  Mumbai,
	}
	fmt.Println(did.String())
}
func TestIDFromDIDString(t *testing.T) {

	didFromStr, err := ParseDID("did:iden3:polygon:mumbai:4RzkkAj2G1ugUEdSo676p5ot7dgQqZ8riTfv4Ev1YX2")
	require.NoError(t, err)
	typ := [2]byte{DIDIden3BlockchainType[didFromStr.Blockchain], DIDNetworkType[didFromStr.NetworkID]}
	var genesis [27]byte
	genesis32bytes := hashBytes([]byte("genesistest"))
	copy(genesis[:], genesis32bytes[:])

	id := NewID(typ, genesis)

	var checksum [2]byte
	copy(checksum[:], id[len(id)-2:])
	assert.Equal(t, CalculateChecksum(typ, genesis), checksum)
	assert.Equal(t, didFromStr.ID.String(), id.String())

}
