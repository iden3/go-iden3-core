package issuer

// TODO: Move this to a more appropiate place.

import (
	"fmt"
	"sync"

	"github.com/iden3/go-iden3-core/db"
)

// type PersistentValue interface {
// 	Get() (uint32, error)
// 	Set(v uint32) error
// }

// UniqueNonceGen is a generator of unique nonces with persistent state.
type UniqueNonceGen struct {
	index *StorageValue
	mutex sync.Mutex
}

// NewUniqueNonceGen creates a new unique nonce generator, storing the
// persistent state in the index.
func NewUniqueNonceGen(index *StorageValue) *UniqueNonceGen {
	return &UniqueNonceGen{index: index, mutex: sync.Mutex{}}
}

// Init is required to initialize the unique nonce generator.
func (u *UniqueNonceGen) Init(tx db.Tx) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.index.Set(tx, 0)
}

// Next returns a new unique nonce.
func (u *UniqueNonceGen) Next(tx db.Tx) (uint32, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	i, err := u.index.Get(tx)
	if err != nil {
		return 0, err
	}
	if i == 0xffffffff {
		return 0, fmt.Errorf("Reached maximum nonce value")
	}
	u.index.Set(tx, i+1)
	return i, nil
}
