package core

import (
	"encoding/hex"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/merkletree"
)

type AuxClaim struct {
	entry *merkletree.Entry
}

func (a AuxClaim) Entry() *merkletree.Entry {
	return a.entry
}

func HexToClaim(h string) (merkletree.Claim, error) {
	bytesValue, err := common3.HexDecode(h)
	if err != nil {
		return AuxClaim{}, err
	}
	var dataBytes [128]byte
	copy(dataBytes[:], bytesValue)
	data := merkletree.NewDataFromBytes(dataBytes)
	entry := merkletree.Entry{
		Data: *data,
	}
	claim := AuxClaim{
		entry: &entry,
	}
	return claim, nil
}

func HexArrayToClaimArray(arrH []string) ([]merkletree.Claim, error) {
	var claims []merkletree.Claim
	for _, h := range arrH {
		claim, err := HexToClaim(h)
		if err != nil {
			return claims, err
		}
		claims = append(claims, claim)
	}
	return claims, nil
}

func ClaimToHex(c merkletree.Claim) string {
	h := hex.EncodeToString(c.Entry().Bytes())
	return h
}

func ClaimArrayToHexArray(claims []merkletree.Claim) []string {
	var hexs []string
	for _, c := range claims {
		h := ClaimToHex(c)
		hexs = append(hexs, h)
	}
	return hexs
}

type ClaimObj struct {
	Claim merkletree.Claim
	Proof ProofClaimPartial // TODO once the RootUpdater is ready we can use here the proof part of the Relay (or direct from blockchain)
}
type ClaimObjHex struct {
	Claim string
	Proof ProofClaimPartial
}

func (co *ClaimObj) Hex() ClaimObjHex {
	return ClaimObjHex{
		Claim: ClaimToHex(co.Claim),
		Proof: co.Proof,
	}
}
func ClaimObjArrayToHexArray(cos []ClaimObj) []ClaimObjHex {
	var cosHex []ClaimObjHex
	for _, co := range cos {
		cosHex = append(cosHex, co.Hex())
	}
	return cosHex
}
