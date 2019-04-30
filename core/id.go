package core

import (
	"bytes"
	"crypto/ecdsa"
	"errors"

	"github.com/btcsuite/btcutil/base58"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

var (
	// TypeBJM7140 specifies the BJ-M7-140
	// - first 2 bytes: `00000000 00000000`
	// - max depth tree: 140 levels
	// - curve of k_op$: babyjub
	// - hash function: `MIMC7`
	TypeBJM7140 = [2]byte{0x00, 0x00}

	// TypeS2M7140 specifies the S2-M7-140
	// - first 2 bytes: `00000000 00000100`
	// - max depth tree: 140 levels
	// - curve of k_op: secp256k1
	// - hash function: `MIMC7`
	TypeS2M7140 = [2]byte{0x00, 0x04}
)

// ID is a byte array with [ type | root_genesis | checksum ]
type ID [36]byte

// NewID creates a new ID from a type and genesis
func NewID(typ [2]byte, genesis [32]byte) ID {
	checksum := CalculateChecksum(typ, genesis)
	var b [36]byte
	copy(b[:2], typ[:])
	copy(b[2:], genesis[:])
	copy(b[34:], checksum[:])
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

// IDFromString returns the ID from a given string
func IDFromString(s string) (*ID, error) {
	b := base58.Decode(s)
	return IDFromBytes(b)
}

// IDFromBytes returns the ID from a given byte array
func IDFromBytes(b []byte) (*ID, error) {
	if len(b) != 36 {
		return nil, errors.New("byte array incorrect")
	}
	var bId [36]byte
	copy(bId[:], b[:])
	id := ID(bId)
	if !CheckChecksum(id) {
		return nil, errors.New("checksum error")
	}
	return &id, nil
}

// DecomposeID returns
func DecomposeID(id ID) ([2]byte, [32]byte, [2]byte, error) {
	var typ [2]byte
	var genesis [32]byte
	var checksum [2]byte
	copy(typ[:], id[:2])
	copy(genesis[:], id[2:len(id)-2])
	copy(checksum[:], id[len(id)-2:])
	return typ, genesis, checksum, nil
}

// CalculateChecksum, returns the checksum for a given type and genesis_root,
// where checksum: hash( [type | root_genesis ] )
func CalculateChecksum(typ [2]byte, genesis [32]byte) [2]byte {
	var toHash [34]byte
	copy(toHash[:], typ[:])
	copy(toHash[2:], genesis[:])
	h := utils.HashBytes(toHash[:])
	var checksum [2]byte
	copy(checksum[:], h[len(h)-2:]) // last two bytes
	return checksum
}

// CheckChecksum returns a bool indicating if the ID.Checksum is consistent with the rest of the ID data
func CheckChecksum(id ID) bool {
	typ, genesis, checksum, err := DecomposeID(id)
	if err != nil {
		return false
	}
	c := CalculateChecksum(typ, genesis)
	return bytes.Equal(c[:], checksum[:])
}

// GenerateArrayClaimAuthorizeKSignFromPublicKeys returns an array of ClaimAuthorizeKSignSecp256k1 from the given public keys
func GenerateArrayClaimAuthorizeKSignFromPublicKeys(keys ...*ecdsa.PublicKey) []*ClaimAuthorizeKSignSecp256k1 {
	var claims []*ClaimAuthorizeKSignSecp256k1
	for _, key := range keys {
		claims = append(claims, NewClaimAuthorizeKSignSecp256k1(key))
	}
	return claims
}

// CalculateIdGenesis calculates the ID given the input parameters.
// Adds the given parameters into an efimeral MerkleTree to calculate the MerkleRoot.
// ID: base58 ( [ type | root_genesis | checksum ] )
// where checksum: hash( [type | root_genesis ] )
// where the hash function is MIMC7
func CalculateIdGenesis(kop, krec, krev *ecdsa.PublicKey) (*ID, error) {
	// add the claims into an efimer merkletree to calculate the genesis root to get that identity
	mt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	if err != nil {
		return nil, err
	}
	// generate the Authorize KSign Claims for the given public Keys
	claims := GenerateArrayClaimAuthorizeKSignFromPublicKeys(kop, krec, krev)
	for _, claim := range claims {
		err = mt.Add(claim.Entry())
		if err != nil {
			return nil, err
		}
	}

	idGenesis := mt.RootKey()

	var idGenesisBytes [32]byte
	copy(idGenesisBytes[:], idGenesis.Bytes())
	id := NewID(TypeS2M7140, idGenesisBytes)
	return &id, nil
}
