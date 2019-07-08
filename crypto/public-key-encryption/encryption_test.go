package encryption

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	_kp := GenKP()
	kpHex := ExportBoxKP(&_kp)
	pkHex := ExportBoxPublicKey(&_kp.PublicKey)

	kp, err := ImportBoxKP(kpHex)
	assert.Nil(t, err)
	assert.Equal(t, *kp, _kp)
	pk, err := ImportBoxPublicKey(pkHex)
	assert.Nil(t, err)
	assert.Equal(t, *pk, _kp.PublicKey)

	msg := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

	encMsg := Encrypt(pk, msg)
	decMsg, err := Decrypt(kp, encMsg)
	assert.Nil(t, err)
	assert.Equal(t, msg, decMsg)

}
