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
	h = strings.TrimPrefix(h, "0x")
	return hex.DecodeString(h)
}

// HexDecodeInto decodes a hex string into an array of bytes (dst), verifying
// that the decoded array has the same length as dst.
func HexDecodeInto(dst []byte, h []byte) error {
	if bytes.HasPrefix(h, []byte("0x")) {
		h = h[2:]
	}
	if len(h)/2 != len(dst) {
		return fmt.Errorf("expected %v bytes in hex string, got %v", len(dst), len(h)/2)
	}
	n, err := hex.Decode(dst, h)
	if err != nil {
		return err
	} else if n != len(dst) {
		return fmt.Errorf("expected %v bytes when decoding hex string, got %v", len(dst), n)
	}
	return nil
}

// Uint16ToBytes returns a byte array from a uint16
func Uint16ToBytes(u uint16) []byte {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], u)
	return b[:]
}

// BytesToUint16 returns a uint16 from a byte array
func BytesToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b[:2])
}

// Uint32ToBytes returns a byte array from a uint32
func Uint32ToBytes(u uint32) []byte {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], u)
	return b[:]
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

// Uint32ToEthBytes converts a uint32 to bytes in big endian.
func Uint32ToEthBytes(u uint32) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, u)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

// Uint64ToEthBytes convets a uint64 to bytes in big endian.
func Uint64ToEthBytes(u uint64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, u)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

// EthBytesToUint32 converts bytes as big endian to uint32.
func EthBytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// EthBytesToUint64 converts bytes as big endian to uint64.
func EthBytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// SetBit sets the bit n in the bitmap to 1.
func SetBit(bitmap []byte, n uint) {
	bitmap[n/8] |= 1 << (n % 8)
}

// SetBitBigEndian sets the bit n in the bitmap to 1, in Big Endian.
func SetBitBigEndian(bitmap []byte, n uint) {
	bitmap[uint(len(bitmap))-n/8-1] |= 1 << (n % 8)
}

// TestBit tests whether the bit n in bitmap is 1.
func TestBit(bitmap []byte, n uint) bool {
	return bitmap[n/8]&(1<<(n%8)) != 0
}

// TestBitBigEndian tests whether the bit n in bitmap is 1, in Big Endian.
func TestBitBigEndian(bitmap []byte, n uint) bool {
	return bitmap[uint(len(bitmap))-n/8-1]&(1<<(n%8)) != 0
}

// SwapEndianness swaps the order of the bytes in the slice.
func SwapEndianness(b []byte) []byte {
	o := make([]byte, len(b))
	for i := range b {
		o[len(b)-1-i] = b[i]
	}
	return o
}
