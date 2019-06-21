package core

import (
	"bytes"
	"errors"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/crypto/babyjub"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
)

var (
	// TypeBJM7 specifies the BJ-M7
	// - first 2 bytes: `00000000 00000000`
	// - curve of k_op: babyjub
	// - hash function: `MIMC7`
	TypeBJM7 = [2]byte{0x00, 0x00}

	// TypeS2M7 specifies the S2-M7
	// - first 2 bytes: `00000000 00000100`
	// - curve of k_op: secp256k1
	// - hash function: `MIMC7`
	TypeS2M7 = [2]byte{0x00, 0x04}
)

// ID is a byte array with
// [  type  | root_genesis | checksum ]
// [2 bytes |   27 bytes   | 2 bytes  ]
// where the root_genesis are the first 28 bytes from the hash root_genesis
type ID [31]byte

// NewID creates a new ID from a type and genesis
func NewID(typ [2]byte, genesis [27]byte) ID {
	checksum := CalculateChecksum(typ, genesis)
	var b [31]byte
	copy(b[:2], typ[:])
	copy(b[2:], genesis[:])
	copy(b[29:], checksum[:])
	return ID(b)
}

// String returns a base58 from the ID.Bytes()
func (id *ID) String() string {
	return base58.Encode(id[:])
}

// Bytes returns the bytes from the ID
func (id *ID) Bytes() []byte {
	return id[:]
}

// func (id ID) MarshalJSON() ([]byte, error) {
//         fmt.Println(id.String())
//         return json.Marshal(id.String())
// }

func (id ID) MarshalText() ([]byte, error) {
	// return json.Marshal(id.String())
	return []byte(id.String()), nil
}

func (id *ID) UnmarshalText(b []byte) error {
	var err error
	var idFromString ID
	idFromString, err = IDFromString(string(b))
	copy(id[:], idFromString[:])
	return err
}

// IDFromString returns the ID from a given string
func IDFromString(s string) (ID, error) {
	b := base58.Decode(s)
	return IDFromBytes(b)
}

var emptyID [31]byte

// IDFromBytes returns the ID from a given byte array
func IDFromBytes(b []byte) (ID, error) {
	if len(b) != 31 {
		return ID{}, errors.New("IDFromBytes error: byte array incorrect length")
	}
	if bytes.Equal(b, emptyID[:]) {
		return ID{}, errors.New("IDFromBytes error: byte array empty")
	}
	var bId [31]byte
	copy(bId[:], b[:])
	id := ID(bId)
	if !CheckChecksum(id) {
		return ID{}, errors.New("IDFromBytes error: checksum error")
	}
	return id, nil
}

// DecomposeID returns type, genesis and checksum from an ID
func DecomposeID(id ID) ([2]byte, [27]byte, [2]byte, error) {
	var typ [2]byte
	var genesis [27]byte
	var checksum [2]byte
	copy(typ[:], id[:2])
	copy(genesis[:], id[2:len(id)-2])
	copy(checksum[:], id[len(id)-2:])
	return typ, genesis, checksum, nil
}

// CalculateChecksum, returns the checksum for a given type and genesis_root,
// where checksum: hash( [type | root_genesis ] )
func CalculateChecksum(typ [2]byte, genesis [27]byte) [2]byte {
	var toChecksum [29]byte
	copy(toChecksum[:], typ[:])
	copy(toChecksum[2:], genesis[:])

	s := uint16(0)
	for _, b := range toChecksum {
		s += uint16(b)
	}
	var checksum [2]byte
	checksum[0] = byte(s >> 8)
	checksum[1] = byte(s & 0xff)
	return checksum
}

// CheckChecksum returns a bool indicating if the ID.Checksum is consistent with the rest of the ID data
func CheckChecksum(id ID) bool {
	typ, genesis, checksum, err := DecomposeID(id)
	if err != nil {
		return false
	}
	if bytes.Equal(checksum[:], []byte{0, 0}) {
		return false
	}
	c := CalculateChecksum(typ, genesis)
	return bytes.Equal(c[:], checksum[:])
}

// CalculateIdGenesis calculates the ID given the input parameters.
// Adds the given parameters into an ephemeral MerkleTree to calculate the MerkleRoot.
// ID: base58 ( [ type | root_genesis | checksum ] )
// where checksum: hash( [type | root_genesis ] )
// where the hash function is MIMC7
func CalculateIdGenesis(kop *babyjub.PublicKey, kdis, kreen, kupdateRoot common.Address) (*ID, []merkletree.Entrier, error) {
	// use an efimer merkletree in memory to calculate the genesis root to get that identity
	mt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	if err != nil {
		return nil, nil, err
	}

	claimKOp := NewClaimAuthorizeKSignBabyJub(kop)
	err = mt.Add(claimKOp.Entry())
	if err != nil {
		return nil, nil, err
	}

	claimKDis := NewClaimAuthEthKey(kdis, EthKeyTypeDisable)
	err = mt.Add(claimKDis.Entry())
	if err != nil {
		return nil, nil, err
	}
	claimKReen := NewClaimAuthEthKey(kreen, EthKeyTypeReenable)
	err = mt.Add(claimKReen.Entry())
	if err != nil {
		return nil, nil, err
	}
	claimKUpdateRoot := NewClaimAuthEthKey(kupdateRoot, EthKeyTypeUpdateRoot)
	err = mt.Add(claimKUpdateRoot.Entry())
	if err != nil {
		return nil, nil, err
	}

	idGenesis := mt.RootKey()

	var idGenesisBytes [27]byte
	copy(idGenesisBytes[:], idGenesis.Bytes()[len(idGenesis.Bytes())-27:])
	id := NewID(TypeBJM7, idGenesisBytes)
	return &id, []merkletree.Entrier{claimKOp}, nil
}
