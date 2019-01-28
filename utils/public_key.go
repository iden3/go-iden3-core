package utils

import (
	"crypto/ecdsa"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
)

// PublicKey is a secp256k1 public key used to verify ecdsa signatures.
type PublicKey struct {
	ecdsa.PublicKey
}

// MarshalJSON serializes the public key as a hex string.
func (pk *PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(common3.HexEncode(crypto.CompressPubkey(&pk.PublicKey)))
}

// UnmarshalJSON deserializes the public key from a hex string.
func (pk *PublicKey) UnmarshalJSON(bs []byte) error {
	pkBytes := [33]byte{}
	err := common3.UnmarshalJSONHexDecodeInto(pkBytes[:], bs)
	if err != nil {
		return err
	}
	pk1, err := crypto.DecompressPubkey(pkBytes[:])
	if err != nil {
		return err
	}
	pk.PublicKey = *pk1
	return nil
}
