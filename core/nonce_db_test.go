package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonceDb(t *testing.T) {
	ndb := NewNonceDb()

	// Test NonceDb.Add()
	for i := int64(0); i < 256; i++ {
		nObj := ndb.Add(fmt.Sprintf("nonce-a-%v", i), 10, nil)
		assert.NotNil(t, nObj)
	}
	// Can't add a repeated nonce
	nObj := ndb.Add(fmt.Sprintf("nonce-a-%v", 0), 10, nil)
	assert.Nil(t, nObj)

	// Adding an Aux
	ok := ndb.AddAux("nonce-a-0", 42)
	assert.Equal(t, true, ok)

	// Adding an Aux to a nonce that already has one must fail
	ok = ndb.AddAux("nonce-a-0", 64)
	assert.Equal(t, false, ok)

	nObj, ok = ndb.Search("nonce-a-0")
	assert.Equal(t, true, ok)
	assert.Equal(t, 42, nObj.Aux)

	// Test NonceDb.Search()
	for i := int64(0); i < 256; i++ {
		_, ok := ndb.Search(fmt.Sprintf("nonce-a-%v", i))
		assert.Equal(t, true, ok)
	}

	// Test NonceDb.SearchAndDelete()
	for i := int64(0); i < 256; i++ {
		_, ok := ndb.SearchAndDelete(fmt.Sprintf("nonce-a-%v", i))
		assert.Equal(t, true, ok)
	}

	// Must not exists because it was deleted
	_, ok = ndb.Search("nonce-a-0")
	assert.Equal(t, false, ok)

	// Must not exists because it has expired
	nObj = ndb.Add("nonce-b-0", -1, nil)
	assert.NotNil(t, nObj)
	_, ok = ndb.Search("nonce-b-0")
	assert.Equal(t, false, ok)

	ndb = NewNonceDb()

	// DeleteOld should delete half of the nonces
	for i := int64(0); i < 8; i++ {
		nObj := ndb.Add(fmt.Sprintf("nonce-c-%v", i), -60, nil)
		assert.NotNil(t, nObj)
	}
	for i := int64(0); i < 8; i++ {
		nObj := ndb.Add(fmt.Sprintf("nonce-d-%v", i), 60, nil)
		assert.NotNil(t, nObj)
	}
	assert.Equal(t, 16, len(ndb.nonceObjsByNonce))
	ndb.DeleteOld()
	assert.Equal(t, 8, len(ndb.nonceObjsByNonce))

	ndb = NewNonceDb()

	// DeleteOldOportunistic should delete half of the nonces after 128 searches
	for i := int64(0); i < 8; i++ { // Add 8 expired nonces
		nObj := ndb.Add(fmt.Sprintf("nonce-e-%v", i), -60, nil)
		assert.NotNil(t, nObj)
	}
	for i := int64(0); i < 8; i++ { // Add 8 non-expired nonces
		nObj := ndb.Add(fmt.Sprintf("nonce-f-%v", i), 60, nil)
		assert.NotNil(t, nObj)

	}
	assert.Equal(t, 16, len(ndb.nonceObjsByNonce))
	ndb.Search("nope") // counter = 1
	assert.Equal(t, 16, len(ndb.nonceObjsByNonce))
	// counter += 127.  When counter == 128, oportunistic delete clears expired nonces.
	for i := 0; i < 127; i++ {
		ndb.Search("nope")
	}
	assert.Equal(t, 8, len(ndb.nonceObjsByNonce))
	// We add more expired nonces
	for i := int64(0); i < 8; i++ {
		nObj := ndb.Add(fmt.Sprintf("nonce-g-%v", i), -60, nil)
		assert.NotNil(t, nObj)
	}
	// counter += 100
	for i := 0; i < 100; i++ {
		ndb.Search("nope")
	}
	// counter < 128, so oportunistic delete hasn't cleared expired nonces yet.
	assert.Equal(t, 16, len(ndb.nonceObjsByNonce))
}
