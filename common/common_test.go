package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToBase64(t *testing.T) {
	base64 := BytesToBase64([]byte("test"))
	assert.Equal(t, "dGVzdA==", base64)
}
func TestBase64ToBytes(t *testing.T) {
	s, err := Base64ToBytes("dGVzdA==")
	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), s)
}

func TestBytesToHex(t *testing.T) {
	h := BytesToHex([]byte("test"))
	assert.Equal(t, "0x74657374", h)
}
func TestHexToBytes(t *testing.T) {
	s, err := HexToBytes("0x74657374")
	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), s)
}

func TestUint32ToBytes(t *testing.T) {
	b, err := Uint32ToBytes(999)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0xe7, 0x3, 0x0, 0x0}, b)
}
func TestBytesToUint32(t *testing.T) {
	u := BytesToUint32([]byte{0xe7, 0x3, 0x0, 0x0})
	assert.Equal(t, uint32(999), u)
}
