package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashBytes(t *testing.T) {
	h := HashBytes([]byte("test")).Hex()
	assert.Equal(t, "0x9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658", h)

	h = HashBytes([]byte("authorizeksign")).Hex()
	assert.Equal(t, "0x353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a05ed7726d7932a1f", h)
}
