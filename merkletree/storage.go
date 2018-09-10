package merkletree

import (
	"bytes"

	common3 "github.com/iden3/go-iden3/common"
	"github.com/syndtr/goleveldb/leveldb"
)

type Storage interface {
	NewTx() (StorageTx, error)
	WithPrefix(prefix []byte) Storage
	Get(key Hash) (byte, uint32, []byte, error)
	Close()
}

type StorageTx interface {
	Get(key Hash) (byte, uint32, []byte, error)
	Insert(stKey Hash, nodeType byte, indexLength uint32, nodeBytes []byte)
	Commit() error
	Close()
}

type LevelDbStorage struct {
	ldb    *leveldb.DB
	prefix []byte
}

type LevelDbStorageTx struct {
	*LevelDbStorage
	cache map[Hash][]byte
}

func NewLevelDbStorage(path string) (*LevelDbStorage, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDbStorage{ldb, []byte{}}, nil
}

func (l *LevelDbStorage) WithPrefix(prefix []byte) Storage {
	return &LevelDbStorage{l.ldb, append(l.prefix, prefix...)}
}

func (l *LevelDbStorage) NewTx() (StorageTx, error) {
	return &LevelDbStorageTx{l, make(map[Hash][]byte)}, nil
}

// Get retreives a value from a key in the mt.Lvl
func (l *LevelDbStorage) Get(key Hash) (byte, uint32, []byte, error) {

	var value []byte
	var err error

	// if key is EMPTY node
	if bytes.Equal(key[:], EmptyNodeValue[:]) {
		return 0, 0, EmptyNodeValue[:], nil
	}

	value, err = l.ldb.Get(append(l.prefix, key[:]...), nil)

	if err != nil { // not found
		return 0, 0, EmptyNodeValue[:], err
	}

	// get nodetype of the first byte of the value
	nodeType := value[0]
	indexLength := common3.BytesToUint32(value[1:5])
	nodeBytes := value[5:]
	return nodeType, indexLength, nodeBytes, err
}

// Get retreives a value from a key in the mt.Lvl
func (l *LevelDbStorageTx) Get(key Hash) (byte, uint32, []byte, error) {

	var ok bool
	var value []byte
	var err error

	if value, ok = l.cache[key]; !ok {

		// if key is EMPTY node
		if bytes.Equal(key[:], EmptyNodeValue[:]) {
			return 0, 0, EmptyNodeValue[:], nil
		}

		value, err = l.ldb.Get(append(l.prefix, key[:]...), nil)

		if err != nil { // not found
			return 0, 0, EmptyNodeValue[:], err
		}
	}

	// get nodetype of the first byte of the value
	nodeType := value[0]
	indexLength := common3.BytesToUint32(value[1:5])
	nodeBytes := value[5:]
	return nodeType, indexLength, nodeBytes, err
}

// Insert saves a key:value into the mt.Lvl
func (l *LevelDbStorageTx) Insert(stKey Hash, nodeType byte, indexLength uint32, nodeBytes []byte) {

	// add nodetype at the first byte of the value
	var stValue []byte
	stValue = append(stValue, nodeType)
	indexLengthBytes, err := common3.Uint32ToBytes(indexLength)
	if err != nil {
		panic(err)
	}
	stValue = append(stValue, indexLengthBytes[:]...)
	stValue = append(stValue, nodeBytes[:]...)

	l.cache[stKey] = stValue
}

func (l *LevelDbStorageTx) Commit() error {

	var batch leveldb.Batch
	for k, v := range l.cache {
		batch.Put(append(l.prefix, k[:]...), v)
	}

	l.cache = nil
	return l.ldb.Write(&batch, nil)
}

func (l *LevelDbStorageTx) Close() {
	l.cache = nil
}

func (l *LevelDbStorage) Close() {
	if err := l.ldb.Close(); err != nil {
		panic(err)
	}
}
