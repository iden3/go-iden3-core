package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
)

// Signature is a secp256k1 ecdsa signature.
type Signature [65]byte

// UnmarshalJSON deserializes a signature from a hex string.
func (s *Signature) UnmarshalJSON(bs []byte) error {
	return common3.UnmarshalJSONHexDecodeInto(s[:], bs)
}

// MarshalJSON serializes a signature as a hex string.
func (s *Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(common3.HexEncode(s[:]))
}

// SignatureEthMsg is a secp256k1 ecdsa signature of an ethereum message:
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sig://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
type SignatureEthMsg [65]byte

// UnmarshalJSON deserializes a signature from a hex string.
func (s *SignatureEthMsg) UnmarshalJSON(bs []byte) error {
	return common3.UnmarshalJSONHexDecodeInto(s[:], bs)
}

// MarshalJSON serializes a signature as a hex string.
func (s *SignatureEthMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(common3.HexEncode(s[:]))
}

const web3SignaturePrefix = "\x19Ethereum Signed Message:\n"

// Sign performs the signature over a Hash
func Sign(h Hash, ks *keystore.KeyStore, acc accounts.Account) ([]byte, error) {
	return ks.SignHash(acc, h[:])
}

// SignEthMsg performs an ethereum message signature over a Hash.
func SignEthMsg(ks *keystore.KeyStore, acc accounts.Account, msg []byte) (*SignatureEthMsg, error) {
	hash := EthHash(msg)
	sig, err := ks.SignHash(acc, hash[:])
	if err != nil {
		return nil, err
	}
	sig[64] += 27
	sigEthMsg := &SignatureEthMsg{}
	copy(sigEthMsg[:], sig)
	return sigEthMsg, nil
}

// EthHash is the hashing function used before signing ethereum messages.
func EthHash(b []byte) Hash {
	header := fmt.Sprintf("%s%d", web3SignaturePrefix, len(b))
	return HashBytes([]byte(header), b)
}

func VerifySigEthMsg(addr common.Address, sig *SignatureEthMsg, msg []byte) bool {
	hash := EthHash(msg)
	var _sig SignatureEthMsg
	copy(_sig[:], sig[:])
	_sig[64] -= 27
	return VerifySig(addr, (*Signature)(&_sig), hash[:])
}

// VerifySigEthMsgDate verifies the signature of a byte array with a date
// appended given an ethereum address.
func VerifySigEthMsgDate(addr common.Address, sig *SignatureEthMsg, msg []byte, date int64) bool {
	dateBytes := Uint64ToEthBytes(uint64(date))
	return VerifySigEthMsg(addr, sig, append(msg[:], dateBytes...))
}

// VerifySig verifies a given signature and the msgHash with the expected address
func VerifySig(addr common.Address, sig *Signature, msgHash []byte) bool {
	recoveredPub, err := crypto.Ecrecover(msgHash, sig[:])
	if err != nil {
		fmt.Printf("ECRecover error: %s\n", err)
		return false
	}
	pubKey, _ := crypto.UnmarshalPubkey(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return bytes.Equal(addr.Bytes(), recoveredAddr.Bytes())
}

// GetPkFromKeyStore is a hack to obtain the public key of an addres who's
// private key is stored in a key store.  It does this by signing an empty hash
// and recovering the public key from the signature.
func GetPkFromKeyStore(ks *keystore.KeyStore, addr common.Address) (*ecdsa.PublicKey, error) {
	var h [256 / 8]byte
	sig, err := ks.SignHash(accounts.Account{Address: addr}, h[:])
	if err != nil {
		return nil, err
	}
	pk, err := crypto.Ecrecover(h[:], sig)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPubkey(pk)
}

// VerifySigBytes verifies the signature of a byte array given an ethereum address.
//func VerifySigBytes(addr common.Address, sig *Signature, msg []byte) bool {
//	msgHash := EthHash(msg)
//	return VerifySig(addr, sig, msgHash[:])
//}

// VerifySigBytesDate verifies the signature of a byte array with a date
// appended given an ethereum address.
//func VerifySigBytesDate(addr common.Address, sig *Signature, msg []byte, date uint64) bool {
//	dateBytes := Uint64ToEthBytes(date)
//	return VerifySigBytes(addr, sig, append(msg[:], dateBytes...))
//}
