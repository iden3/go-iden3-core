package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"
)

// PoWData is the interface for the data that have the Nonce parameter to calculate the Proof-of-Work
type PoWData interface {
	IncrementNonce() PoWData
}

// CheckPoW verifies the PoW difficulty of a Hash
func CheckPoW(hash Hash, difficulty int) bool {
	var empty [32]byte
	if !bytes.Equal(hash[:][0:difficulty], empty[0:difficulty]) {
		return false
	}
	return true
}

// PoW calculates the nonce for the given data in order to fit in the current Proof of Work difficulty
func PoW(data PoWData, difficulty int) (PoWData, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	hash := HashBytes(b)
	for !CheckPoW(hash, difficulty) {
		data = data.IncrementNonce()

		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		hash = HashBytes(b)
	}
	return data, nil
}

func VerifyTimestamp(timestamp int64, timelimit int) bool {
	t := time.Unix(timestamp, 10)
	elapsed := time.Since(t)
	if int(elapsed.Seconds()) > timelimit {
		return false
	}
	return true
}

// Uint32ToEthBytes converts a uint32 to bytes in big endian.
func Uint32ToEthBytes(u uint32) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, u)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

// Uint64ToEthBytes convets a uint64 to bytes in big endian.
func Uint64ToEthBytes(u uint64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, u)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

// EthBytesToUint32 converts bytes as big endian to uint32.
func EthBytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// EthBytesToUint64 converts bytes as big endian to uint64.
func EthBytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
