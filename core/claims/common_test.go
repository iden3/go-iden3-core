package claims

import (
	"encoding/hex"

	"github.com/iden3/go-iden3-core/merkletree"
)

func hexStringToKey(s string) merkletree.Hash {
	var keyBytes [merkletree.ElemBytesLen]byte
	keyBytesHex, _ := hex.DecodeString(s)
	copy(keyBytes[:], keyBytesHex[:merkletree.ElemBytesLen])
	return merkletree.Hash(merkletree.ElemBytes(keyBytes))
}
