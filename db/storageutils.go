package db

import (
	"encoding/binary"
	"encoding/json"
)

// StorageValue allows storing a uint32 persistently
type StorageValue struct {
	dbKey []byte
}

// NewStorageValue creates a new StorageValue that uses the dbKey in a
// Storage to store the value.
func NewStorageValue(dbKey []byte) *StorageValue {
	return &StorageValue{dbKey: dbKey}
}

// Set sets the value in an open db transaction.
func (sv *StorageValue) Set(tx Tx, v uint32) {
	var vBytes [4]byte
	binary.LittleEndian.PutUint32(vBytes[:], v)
	tx.Put(sv.dbKey, vBytes[:])
}

// Get returns the current value in an open db transaction.
func (sv *StorageValue) Get(tx Tx) (uint32, error) {
	vBytes, err := tx.Get(sv.dbKey)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(vBytes), nil
}

// StorageList allows storing a list of key values that are also stored by
// index number.
type StorageList struct {
	length            *StorageValue
	dbPrefixList      []byte
	dbPrefixListByIdx []byte
}

// NewStorageList creates a new StorageList that will store the contents under
// the dbPrefix in a Storage.
func NewStorageList(dbPrefix []byte) *StorageList {
	return &StorageList{
		length:            NewStorageValue(append(dbPrefix, []byte("len")...)),
		dbPrefixList:      append(dbPrefix, []byte("list:")...),
		dbPrefixListByIdx: append(dbPrefix, []byte("byidx:")...),
	}
}

// Init initializes the Storage list in an open db transaction.
func (sl *StorageList) Init(tx Tx) {
	sl.length.Set(tx, 0)
}

// Append adds a new key value entry to the StorageList in an open db transaction.
func (sl *StorageList) Append(tx Tx, key []byte, value interface{}) error {
	idx, err := sl.length.Get(tx)
	if err != nil {
		return err
	}
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var idxBytes [4]byte
	binary.LittleEndian.PutUint32(idxBytes[:], idx)
	tx.Put(append(sl.dbPrefixList, key...), valueJSON)
	tx.Put(append(sl.dbPrefixListByIdx, idxBytes[:]...), key)
	sl.length.Set(tx, idx+1)
	return nil
}

// GetByIdx returns the key value given the index of the StorageList in an open db transaction.
func (sl *StorageList) GetByIdx(tx Tx, idx uint32, value interface{}) ([]byte, error) {
	var idxBytes [4]byte
	binary.LittleEndian.PutUint32(idxBytes[:], idx)
	key, err := tx.Get(append(sl.dbPrefixListByIdx, idxBytes[:]...))
	if err != nil {
		return nil, err
	}
	err = sl.Get(tx, key, value)
	return key, err
}

// GetByIdx returns the value given the key of the StorageList in an open db transaction.
func (sl *StorageList) Get(tx Tx, key []byte, value interface{}) error {
	valueJSON, err := tx.Get(append(sl.dbPrefixList, key...))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(valueJSON, value); err != nil {
		return err
	}
	return err
}

// Length returns the number of elements in the StorageList in an open db transaction.
func (sl *StorageList) Length(tx Tx) (uint32, error) {
	return sl.length.Get(tx)
}

func StoreJSON(tx Tx, key []byte, v interface{}) error {
	vJSON, err := json.Marshal(v)
	if err != nil {
		return err
	}
	tx.Put(key, vJSON)
	return nil
}

func LoadJSON(storage Storage, key []byte, v interface{}) error {
	vJSON, err := storage.Get(key)
	if err == ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}
	return json.Unmarshal(vJSON, v)
}
