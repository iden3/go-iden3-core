package core

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/utils"
	"github.com/stretchr/testify/assert"
)

func TestIDparsers(t *testing.T) {
	typ0 := [2]byte{0x00, 0x00}
	var genesis0 [28]byte
	genesis032bytes := utils.HashBytes([]byte("genesistest"))
	copy(genesis0[:], genesis032bytes[:])

	id0 := NewID(typ0, genesis0)
	assert.Equal(t, id0.String(), "11iu5zg52EFoKcKcE6ChPiChftxBdidpxzh4vsmzf6F")

	typ1 := [2]byte{0x00, 0x01}
	var genesis1 [28]byte
	genesis132bytes := utils.HashBytes([]byte("genesistest"))
	copy(genesis1[:], genesis132bytes[:])

	id1 := NewID(typ1, genesis1)
	assert.Equal(t, id1.String(), "12cCE371biv2KzYyBiHZDGmc4nHPaJsQBCYS7Zwxkyk")

	emptyChecksum := []byte{0x00, 0x00}
	assert.True(t, !bytes.Equal(emptyChecksum, id0[30:]))
	assert.True(t, !bytes.Equal(emptyChecksum, id1[30:]))

	id0FromBytes, err := IDFromBytes(id0.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromBytes.Bytes())
	assert.Equal(t, id0.String(), id0FromBytes.String())
	assert.Equal(t, id0FromBytes.String(), "11iu5zg52EFoKcKcE6ChPiChftxBdidpxzh4vsmzf6F")

	id1FromBytes, err := IDFromBytes(id1.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, id1.Bytes(), id1FromBytes.Bytes())
	assert.Equal(t, id1.String(), id1FromBytes.String())
	assert.Equal(t, id1FromBytes.String(), "12cCE371biv2KzYyBiHZDGmc4nHPaJsQBCYS7Zwxkyk")

	id0FromString, err := IDFromString(id0.String())
	assert.Nil(t, err)
	assert.Equal(t, id0.Bytes(), id0FromString.Bytes())
	assert.Equal(t, id0.String(), id0FromString.String())
	assert.Equal(t, id0FromString.String(), "11iu5zg52EFoKcKcE6ChPiChftxBdidpxzh4vsmzf6F")
}

func TestCheckChecksum(t *testing.T) {
	typ := TypeS2M7
	var genesis [28]byte
	genesis32bytes := utils.HashBytes([]byte("genesistest"))
	copy(genesis[:], genesis32bytes[:])

	id := NewID(typ, genesis)

	var checksum [2]byte
	copy(checksum[:], id[len(id)-2:])
	assert.Equal(t, CalculateChecksum(typ, genesis), checksum)

	assert.True(t, CheckChecksum(id))

	// check that if we change the checksum, returns false on CheckChecksum
	id = NewID(typ, genesis)
	copy(id[30:], []byte{0x00, 0x00})
	assert.True(t, !CheckChecksum(id))

	// check that if we change the type, returns false on CheckChecksum
	id = NewID(typ, genesis)
	copy(id[:2], []byte{0x00, 0x00})
	assert.True(t, !CheckChecksum(id))

	// check that if we change the genesis, returns false on CheckChecksum
	id = NewID(typ, genesis)
	// changedGenesis := utils.HashBytes([]byte("changedgenesis"))
	var changedGenesis [28]byte
	changedGenesis32bytes := utils.HashBytes([]byte("changedgenesis"))
	copy(changedGenesis[:], changedGenesis32bytes[:28])

	copy(id[2:28], changedGenesis[:])
	assert.True(t, !CheckChecksum(id))
}

func TestCalculateIdGenesis(t *testing.T) {
	kopStr := "0x037e211781efef4687e78be4fb008768acca8101b6f1f7ea099599f02a8813f386"
	krecStr := "0x03f9737be33b5829e3da80160464b2891277dae7d7c23609f9bb34bd4ede397bbf"
	krevStr := "0x02d2da59d3022b4c1589e4910baa6cbaddd01f95ed198fdc3068d9dc1fb784a9a4"

	kopBytes, err := common3.HexDecode(kopStr)
	assert.Nil(t, err)
	kopPub, err := crypto.DecompressPubkey(kopBytes[:])
	assert.Nil(t, err)

	krecBytes, err := common3.HexDecode(krecStr)
	assert.Nil(t, err)
	krecPub, err := crypto.DecompressPubkey(krecBytes[:])
	assert.Nil(t, err)

	krevBytes, err := common3.HexDecode(krevStr)
	assert.Nil(t, err)
	krevPub, err := crypto.DecompressPubkey(krevBytes[:])
	assert.Nil(t, err)

	idAddr, err := CalculateIdGenesis(kopPub, krecPub, krevPub)
	assert.Nil(t, err)
	if debug {
		fmt.Println("idAddr", idAddr)
		fmt.Println("idAddr (hex)", idAddr.String())
	}
	assert.Equal(t, idAddr.String(), "14ZjYvqmB4sNQG5o1orKUWinDt1Zk13iVvs4zPCPmvE")
}
