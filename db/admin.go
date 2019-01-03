package db

import (
	"encoding/hex"
	"fmt"
)

func (l *LevelDbStorage) RawDump() error {
	iter := l.ldb.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Println(hex.EncodeToString(iter.Key()), " ", hex.EncodeToString(iter.Value()))
	}
	iter.Release()
	return nil
}

func IPFSexport() error {
	return nil
}
