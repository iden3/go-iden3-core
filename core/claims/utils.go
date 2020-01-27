package claims

import (
	"encoding/hex"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/merkletree"
)

type ClaimGeneric struct {
	entry *merkletree.Entry
}

func (a ClaimGeneric) Entry() *merkletree.Entry {
	return a.entry
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
