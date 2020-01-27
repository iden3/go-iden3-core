package claims

import (
	"encoding/hex"
	"testing"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
)

func hexStringToKey(s string) merkletree.Hash {
	var keyBytes [merkletree.ElemBytesLen]byte
	keyBytesHex, _ := hex.DecodeString(s)
	copy(keyBytes[:], keyBytesHex[:merkletree.ElemBytesLen])
	return merkletree.Hash(merkletree.ElemBytes(keyBytes))
}

func checkClaim(e *merkletree.Entry, t *testing.T) {
	testgen.CheckTestValue("HIndex", e.HIndex().Hex(), t)
	testgen.CheckTestValue("HValue", e.HValue().Hex(), t)
	testgen.CheckTestValue("dataString", e.Data.String(), t)
}
