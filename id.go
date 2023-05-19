package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/mr-tron/base58"
)

var (
	// TypeDefault specifies the regular identity
	// - first 2 bytes: `00000000 00000000`
	TypeDefault = [2]byte{0x00, 0x00}

	// TypeDID specifies the identity with iden3 method in specific networks
	// - first byte: did method e.g. 00000001 - iden3 did method
	// - second byte - blockchain network
	// - 0-3 bits of 2nd byte: blockchain network e.g. 0001 - polygon
	// - 4-7 bits of 2nd byte: network id e.g. 0010 - mumbai
	//  example of 2nd byte: 00010010 - polygon mumbai, 00000000 - readonly identities.
	// valid iden3 method {0b00000001,0b00010010}, readonly {0b00000001, 0b00000000}
)

const idLength = 31
const genesisLn = 27

// ID is a byte array with
// [  type  | root_genesis | checksum ]
// [2 bytes |   27 bytes   | 2 bytes  ]
// where the root_genesis are the first 28 bytes from the hash root_genesis
type ID [idLength]byte

// NewID creates a new ID from a type and genesis
func NewID(typ [2]byte, genesis [genesisLn]byte) ID {
	checksum := CalculateChecksum(typ, genesis)
	var b ID
	copy(b[:2], typ[:])
	copy(b[2:], genesis[:])
	copy(b[29:], checksum[:])
	return b
}

// ProfileID calculates the Profile ID from the Identity and profile nonce. If nonce is empty or zero ID is returned
func ProfileID(id ID, nonce *big.Int) (ID, error) {

	if nonce == nil || big.NewInt(0).Cmp(nonce) == 0 {
		return id, nil
	}

	hash, err := poseidon.Hash([]*big.Int{id.BigInt(), nonce})
	if err != nil {
		return ID{}, err
	}

	typ, _, _, err := DecomposeID(id)
	if err != nil {
		return ID{}, err
	}

	var genesis [genesisLn]byte
	copy(genesis[:], firstNBytes(hash, genesisLn))
	return NewID(typ, genesis), nil
}

// firstNBytes encodes big int in little endian representation and return
// lowers n bytes
func firstNBytes(i *big.Int, n uint) []byte {
	b := intToBytes(i)
	if len(b) > int(n) {
		return b[:n]
	}
	b2 := make([]byte, n)
	copy(b2, b)
	return b2
}

// String returns a base58 from the ID
func (id *ID) String() string {
	return base58.Encode(id[:])
}

// Bytes returns the bytes from the ID
func (id *ID) Bytes() []byte {
	return id[:]
}

func (id *ID) BigInt() *big.Int {
	var s ElemBytes
	copy(s[:], id[:])
	return s.ToInt()
}

func (id *ID) Equal(id2 *ID) bool {
	return bytes.Equal(id[:], id2[:])
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *ID) UnmarshalText(b []byte) error {
	var err error
	var idFromString ID
	idFromString, err = IDFromString(string(b))
	copy(id[:], idFromString[:])
	return err
}

func (id *ID) Equals(id2 *ID) bool {
	return bytes.Equal(id[:], id2[:])
}

func (id *ID) Type() [2]byte {
	var typ [2]byte
	copy(typ[:], id[:2])
	return typ
}

// IDFromString returns the ID from a given string
func IDFromString(s string) (ID, error) {
	b, err := base58.Decode(s)
	if err != nil {
		return ID{}, err
	}
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
	var bID [31]byte
	copy(bID[:], b[:])
	id := ID(bID)
	if !CheckChecksum(id) {
		return ID{}, errors.New("IDFromBytes error: checksum error")
	}
	return id, nil
}

// IDFromInt returns the ID from a given big.Int
func IDFromInt(i *big.Int) (ID, error) {
	b := intToBytes(i)
	if len(b) > idLength {
		return ID{}, errors.New("IDFromInt error: big.Int too large")
	}
	for len(b) < idLength {
		b = append(b, make([]byte, idLength-len(b))...)
	}
	return IDFromBytes(b)
}

// DecomposeID returns type, genesis and checksum from an ID
func DecomposeID(id ID) (typ [2]byte, genesis [genesisLn]byte, checksum [2]byte,
	err error) {

	copy(typ[:], id[:2])
	copy(genesis[:], id[2:len(id)-2])
	copy(checksum[:], id[len(id)-2:])
	return typ, genesis, checksum, nil
}

// CalculateChecksum returns the checksum for a given type and genesis_root,
// where checksum:
//
//	hash( [type | root_genesis ] )
func CalculateChecksum(typ [2]byte, genesis [genesisLn]byte) [2]byte {
	var toChecksum [29]byte
	copy(toChecksum[:], typ[:])
	copy(toChecksum[2:], genesis[:])

	s := uint16(0)
	for _, b := range toChecksum {
		s += uint16(b)
	}
	var checksum [2]byte
	binary.LittleEndian.PutUint16(checksum[:], s)
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

// NewIDFromIdenState calculates the genesis ID from an Identity State.
func NewIDFromIdenState(typ [2]byte, state *big.Int) (*ID, error) {
	var idGenesisBytes [genesisLn]byte

	idenStateData, err := NewElemBytesFromInt(state)
	if err != nil {
		return nil, err
	}

	// we take last 27 bytes, because of swapped endianness
	copy(idGenesisBytes[:], idenStateData[len(idenStateData)-genesisLn:])
	id := NewID(typ, idGenesisBytes)
	return &id, nil
}

// IdenState calculates the Identity State from the Claims Tree Root,
// Revocation Tree Root and Roots Tree Root.
func IdenState(clr, rer, ror *big.Int) (*big.Int, error) {
	return poseidon.Hash([]*big.Int{clr, rer, ror})
}

// CheckGenesisStateID check if the state is genesis for the id.
func CheckGenesisStateID(id, state *big.Int) (bool, error) {
	userID, err := IDFromInt(id)
	if err != nil {
		return false, err
	}
	identifier, err := NewIDFromIdenState(userID.Type(), state)
	if err != nil {
		return false, err
	}

	return id.Cmp(identifier.BigInt()) == 0, nil
}
