package claims

import (
	"encoding/hex"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-merkletree-sql"
)

type ClaimGeneric struct {
	metadata Metadata
	entry    *merkletree.Entry
}

func NewClaimGeneric(entry *merkletree.Entry) *ClaimGeneric {
	var metadata Metadata
	metadata.Unmarshal(entry)
	return &ClaimGeneric{metadata: metadata, entry: entry}
}

func (c *ClaimGeneric) Entry() *merkletree.Entry {
	c.metadata.Marshal(c.entry)
	return c.entry
}

func (c *ClaimGeneric) Metadata() *Metadata {
	return &c.metadata
}

func HexToClaimGeneric(h string) (ClaimGeneric, error) {
	bytesValue, err := common3.HexDecode(h)
	if err != nil {
		return ClaimGeneric{}, err
	}
	var dataBytes [merkletree.ElemBytesLen * merkletree.DataLen]byte
	copy(dataBytes[:], bytesValue)
	data := merkletree.NewDataFromBytes(dataBytes)
	entry := merkletree.Entry{
		Data: *data,
	}
	entrier := ClaimGeneric{
		entry: &entry,
	}
	return entrier, nil
}

func HexArrayToClaimGenericArray(arrH []string) ([]ClaimGeneric, error) {
	var claims []ClaimGeneric
	for _, h := range arrH {
		claim, err := HexToClaimGeneric(h)
		if err != nil {
			return claims, err
		}
		claims = append(claims, claim)
	}
	return claims, nil
}

func ClaimToHex(c merkletree.Entrier) string {
	h := hex.EncodeToString(c.Entry().Bytes())
	return h
}

func ClaimArrayToHexArray(claims []merkletree.Entrier) []string {
	var hexs []string
	for _, c := range claims {
		h := ClaimToHex(c)
		hexs = append(hexs, h)
	}
	return hexs
}

// HexStringToHash decodes a hex string into a Hash.
func HexStringToHash(s string) merkletree.Hash {
	b, err := common3.HexDecode(s)
	if err != nil {
		panic(err)
	}
	var hash merkletree.Hash
	copy(hash[:], b[:32])
	return hash
}
