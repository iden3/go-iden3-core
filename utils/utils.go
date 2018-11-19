package utils

import (
	"encoding/json"

	"github.com/iden3/go-iden3/merkletree"
)

type PoWData interface {
	IncrementNonce() PoWData
}

func CheckPoW(hash merkletree.Hash, difficulty int) bool {
	for i := 0; i < difficulty; i++ {
		if hash[i] != byte(0) {
			return false
		}
	}
	return true
}

func PoW(data PoWData, difficulty int) (PoWData, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	hash := merkletree.HashBytes(b)
	for !CheckPoW(hash, difficulty) {
		data = data.IncrementNonce()

		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		hash = merkletree.HashBytes(b)
	}
	return data, nil
}
