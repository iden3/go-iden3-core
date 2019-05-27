package common

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Base64ToBytes converts a base64 encoded string into array of bytes
func Base64ToBytes(base64String string) ([]byte, error) {
	hashBytes, err := base64.StdEncoding.DecodeString(base64String)
	return hashBytes, err
}

// Base64Decode decodes a string in base64 with optional padding into an array
// of bytes.
func Base64Decode(b64 string) ([]byte, error) {
	l := len(b64)
	if b64[l-3] == '=' {
		b64 = b64[:l-3]
	} else if b64[l-2] == '=' {
		b64 = b64[:l-2]
	} else if b64[l-1] == '=' {
		b64 = b64[:l-1]
	}
	return base64.RawStdEncoding.DecodeString(b64)
}

// BytesToBase64 converts an array of bytes to a base64 encoded string
func BytesToBase64(bytesArray []byte) string {
	h := base64.StdEncoding.EncodeToString(bytesArray)
	return h
}

// Hex is a byte slice type that can be marshalled and unmarshaled in hex
type Hex []byte

// MarshalText encodes buf as hex
func (buf Hex) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(buf)), nil
}

// String encodes buf as hex
func (buf Hex) String() string {
	return hex.EncodeToString(buf)
}

// UnmarshalText decodes a hex into buf
func (buf *Hex) UnmarshalText(h []byte) error {
	*buf = make([]byte, hex.DecodedLen(len(h)))
	if _, err := hex.Decode(*buf, h); err != nil {
		return err
	}
	return nil
}

// HexEncode encodes an array of bytes into a string in hex.
func HexEncode(bs []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(bs))
}

// HexDecode decodes a hex string into an array of bytes.
func HexDecode(h string) ([]byte, error) {
	if strings.HasPrefix(h, "0x") {
		h = h[2:]
	}
	return hex.DecodeString(h)
}

// HexDecodeInto decodes a hex string into an array of bytes (dst), verifying
// that the decoded array has the same length as dst.
func HexDecodeInto(dst []byte, h []byte) error {
	if bytes.HasPrefix(h, []byte("0x")) {
		h = h[2:]
	}
	n, err := hex.Decode(dst, h)
	if err != nil {
		return err
	} else if n != len(dst) {
		return fmt.Errorf("expected %v bytes when decoding hex string, got %v", len(dst), n)
	}
	return nil
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

// UnmarshalJSONHexDecodeInto decodes the JSON Hex string into bs
func UnmarshalJSONHexDecodeInto(dst []byte, bs []byte) error {
	hexStr := ""
	err := json.Unmarshal(bs, &hexStr)
	if err != nil {
		return err
	}
	return HexDecodeInto(dst[:], []byte(hexStr))
}

// UnmarshalJSONHexDecode decodes the JSON Hex string and returns the []byte
func UnmarshalJSONHexDecode(bs []byte) ([]byte, error) {
	hexStr := ""
	err := json.Unmarshal(bs, &hexStr)
	if err != nil {
		return nil, err
	}
	return HexDecode(hexStr)
}
