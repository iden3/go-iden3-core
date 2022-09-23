package core

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/mr-tron/base58"
)

var (
	// TypeDefault specifies the regular identity
	// - first 2 bytes: `00000000 00000000`
	TypeDefault = [2]byte{0x00, 0x00}

	// TypeReadOnly specifies the readonly identity, this type of identity MUST not be published on chain
	// - first 2 bytes: `00000000 00000001`
	TypeReadOnly = [2]byte{0b00000000, 0b00000001}

	// TypeDID specifies the identity with iden3 method in specific networks
	// - first 1 bytes: `0b00000000 polygon main`
	// 					`0b00000001 polygon mumbai`
	// 					`0b00000010 ethereum main`
	// 					`0b00000011 ethereum ropsten`
	// 					`0b00000012 ethereum kovan`
	// 					`0b00000013 ethereum rinkeby`
)

const idLength = 31

// ID is a byte array with
// [  type  | root_genesis | checksum ]
// [2 bytes |   27 bytes   | 2 bytes  ]
// where the root_genesis are the first 28 bytes from the hash root_genesis
type ID [idLength]byte

// NewID creates a new ID from a type and genesis
func NewID(typ [2]byte, genesis [27]byte) ID {
	checksum := CalculateChecksum(typ, genesis)
	var b [31]byte
	copy(b[:2], typ[:])
	copy(b[2:], genesis[:])
	copy(b[29:], checksum[:])
	return ID(b)
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
func DecomposeID(id ID) ([2]byte, [27]byte, [2]byte, error) {
	var typ [2]byte
	var genesis [27]byte
	var checksum [2]byte
	copy(typ[:], id[:2])
	copy(genesis[:], id[2:len(id)-2])
	copy(checksum[:], id[len(id)-2:])
	return typ, genesis, checksum, nil
}

// CalculateChecksum returns the checksum for a given type and genesis_root,
// where checksum:
//   hash( [type | root_genesis ] )
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

// IdGenesisFromIdenState calculates the genesis ID from an Identity State.
func IdGenesisFromIdenState(typ [2]byte, //nolint:revive
	state *big.Int) (*ID, error) {

	var idGenesisBytes [27]byte

	idenStateData, err := NewElemBytesFromInt(state)
	if err != nil {
		return nil, err
	}

	// we take last 27 bytes, because of swapped endianness
	copy(idGenesisBytes[:], idenStateData[len(idenStateData)-27:])
	id := NewID(typ, idGenesisBytes)
	return &id, nil
}

// IdenState calculates the Identity State from the Claims Tree Root,
// Revocation Tree Root and Roots Tree Root.
func IdenState(clr, rer, ror *big.Int) (*big.Int, error) {
	return poseidon.Hash([]*big.Int{clr, rer, ror})
}
