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

// IdGenesisFromIdenState calculates the genesis Id from an Identity State.
func IdGenesisFromIdenState(state ElemBytes) *ID { //nolint:revive
	var idGenesisBytes [27]byte
	// we take last 27 bytes, because of swapped endianness
	copy(idGenesisBytes[:], state[len(state)-27:])
	id := NewID(TypeDefault, idGenesisBytes)
	return &id
}

// IdenState calculates the Identity State from the Claims Tree Root,
// Revocation Tree Root and Roots Tree Root.
func IdenState(clr ElemBytes, rer ElemBytes,
	ror ElemBytes) (ElemBytes, error) {

	idenState, err := poseidon.Hash([]*big.Int{
		clr.ToInt(), rer.ToInt(), ror.ToInt()})
	if err != nil {
		return ElemBytes{}, err
	}
	return NewElemBytesFromInt(idenState)
}

// CalculateGenesisID calculate genesis id based on provided claims tree root
func CalculateGenesisID(clr ElemBytes) (*ID, error) {
	idenState, err := poseidon.Hash([]*big.Int{
		clr.ToInt(), big.NewInt(0), big.NewInt(0)})
	if err != nil {
		return nil, err
	}
	idenStateData, err := NewElemBytesFromInt(idenState)
	if err != nil {
		return nil, err
	}
	return IdGenesisFromIdenState(idenStateData), nil
}
