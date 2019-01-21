package common

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
)

// Base64ToBytes converts a base64 encoded string into array of bytes
func Base64ToBytes(base64String string) ([]byte, error) {
	hashBytes, err := base64.StdEncoding.DecodeString(base64String)
	return hashBytes, err
}

// BytesToBase64 converts an array of bytes to a base64 encoded string
func BytesToBase64(bytesArray []byte) string {
	h := base64.StdEncoding.EncodeToString(bytesArray)
	return h
}

// BytesToHex converts from an array of bytes to a hex encoded string
func BytesToHex(bytesArray []byte) string {
	r := "0x"
	h := hex.EncodeToString(bytesArray)
	r = r + h
	return r
}

// HexToBytes converts from a hex string into an array of bytes
func HexToBytes(h string) ([]byte, error) {
	b, err := hex.DecodeString(h[2:])
	return b, err
}

// Uint32ToBytes returns a byte array from a uint32
func Uint32ToBytes(u uint32) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, u)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

// BytesToUint32 returns a uint32 from a byte array
func BytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

// TODO: Find a better place for this
type PublicKey struct {
	ecdsa.PublicKey
}

func (pk *PublicKey) MarshalJSON() ([]byte, error) {
	s := hex.EncodeToString(crypto.CompressPubkey(&pk.PublicKey))
	return json.Marshal(s)
}

func (pk *PublicKey) UnmarshalJSON(bs []byte) error {
	hexStr := ""
	if err := json.Unmarshal(bs, &hexStr); err != nil {
		return err
	}
	pkBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	pk1, err := crypto.DecompressPubkey(pkBytes)
	pk.PublicKey = *pk1
	return err
}
