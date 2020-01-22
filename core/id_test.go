package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/utils"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

func TestIDparsers(t *testing.T) {
	typ0 := [2]byte{0x00, 0x00}
	var genesis0 [27]byte
	genesis032bytes := utils.HashBytes([]byte("genesistest"))
	copy(genesis0[:], genesis032bytes[:])

	id0 := NewID(typ0, genesis0)
	assert.Equal(t, "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf", id0.String())

	typ1 := [2]byte{0x00, 0x01}
	var genesis1 [27]byte
	genesis132bytes := utils.HashBytes([]byte("genesistest"))
	copy(genesis1[:], genesis132bytes[:])

	id1 := NewID(typ1, genesis1)
	assert.Equal(t, "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2Xh6", id1.String())

	emptyChecksum := []byte{0x00, 0x00}
	assert.True(t, !bytes.Equal(emptyChecksum, id0[29:]))
	assert.True(t, !bytes.Equal(emptyChecksum, id1[29:]))

	id0FromBytes, err := IDFromBytes(id0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromBytes.Bytes())
	assert.Equal(t, id0.String(), id0FromBytes.String())
	assert.Equal(t, "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf", id0FromBytes.String())

	id1FromBytes, err := IDFromBytes(id1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id1.Bytes(), id1FromBytes.Bytes())
	assert.Equal(t, id1.String(), id1FromBytes.String())
	assert.Equal(t, "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2Xh6", id1FromBytes.String())

	id0FromString, err := IDFromString(id0.String())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromString.Bytes())
	assert.Equal(t, id0.String(), id0FromString.String())
	assert.Equal(t, "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf", id0FromString.String())
}

func TestIDjsonParser(t *testing.T) {
	id, err := IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf")
	assert.Nil(t, err)
	idj, err := json.Marshal(&id)
	assert.Nil(t, err)
	assert.Equal(t, `"11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf"`, string(idj))

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
	typ := TypeS2M7
	var genesis [27]byte
	genesis32bytes := utils.HashBytes([]byte("genesistest"))
	copy(genesis[:], genesis32bytes[:])

	id := NewID(typ, genesis)

	var checksum [2]byte
	copy(checksum[:], id[len(id)-2:])
	assert.Equal(t, CalculateChecksum(typ, genesis), checksum)

	assert.True(t, CheckChecksum(id))

	// check that if we change the checksum, returns false on CheckChecksum
	id = NewID(typ, genesis)
	copy(id[29:], []byte{0x00, 0x00})
	assert.True(t, !CheckChecksum(id))

	// check that if we change the type, returns false on CheckChecksum
	id = NewID(typ, genesis)
	copy(id[:2], []byte{0x00, 0x00})
	assert.True(t, !CheckChecksum(id))

	// check that if we change the genesis, returns false on CheckChecksum
	id = NewID(typ, genesis)
	// changedGenesis := utils.HashBytes([]byte("changedgenesis"))
	var changedGenesis [27]byte
	changedGenesis32bytes := utils.HashBytes([]byte("changedgenesis"))
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
	hex.Decode(sk[:], []byte("28156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"))
	kopPub := sk.Public()
	kDis := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kReen := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kUpdateRoot := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	id, _, err := CalculateIdGenesisFrom4Keys(kopPub, kDis, kReen, kUpdateRoot)
	assert.Nil(t, err)
	if debug {
		fmt.Println("id", id)
		fmt.Println("id (hex)", id.String())
	}
	assert.Equal(t, "1LzwQet8DMLnYKBz2WgUvL3WDfjbbPrkAmcekMSUP", id.String())
}

func TestCalculateIdGenesis(t *testing.T) {
	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
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
	assert.Equal(t, "1GURWwRa5YQA8KA2AdmGANhXSpAupfpy2VsHse2QU", id.String())
}
