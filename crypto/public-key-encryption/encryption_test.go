package encryption

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	_kp := GenKP()
	kpHex := hex.EncodeToString(_kp.PublicKey.Bytes[:]) + hex.EncodeToString(_kp.SecretKey.Bytes[:])
	pkHex := hex.EncodeToString(_kp.PublicKey.Bytes[:])

	kp, err := ImportBoxKP(kpHex)
	assert.Equal(t, nil, err)
	pk, err := ImportBoxPublicKey(pkHex)
	assert.Equal(t, nil, err)

	msg := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

	encMsg := Encrypt(pk, msg)
	decMsg, err := Decrypt(kp, encMsg)
	assert.Equal(t, nil, err)
	assert.Equal(t, msg, decMsg)

}
