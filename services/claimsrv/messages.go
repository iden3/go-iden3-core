package claimsrv

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
)

// BytesSignedMsg contains the value and its signature in Hex representation
type BytesSignedMsg struct {
	ValueHex     string             `json:"valueHex"` // claim.Bytes() in a hex format
	SignatureHex string             `json:"signatureHex"`
	KSignPk      *common3.PublicKey `json:"ksignpk" binding:"required"`
}

// ClaimBasicMsg contains a core.ClaimBasic with its signature in Hex
type ClaimBasicMsg struct {
	ClaimBasic core.ClaimBasic
	Signature  string
}

// ClaimAssignNameMsg contains a core.ClaimAssignName with its signature in Hex
type ClaimAssignNameMsg struct {
	ClaimAssignName core.ClaimAssignName
	Signature       string
}

// ClaimAuthorizeKSignMsg contains a core.AuthorizeKSignClaim with its signature in Hex
type ClaimAuthorizeKSignMsg struct {
	ClaimAuthorizeKSign core.ClaimAuthorizeKSign
	Signature           string
	KSignPk             *common3.PublicKey
}

// ClaimAuthorizeKSignSecp256k1Msg contains a core.ClaimAuthorizeKSignP256 with its signature in Hex
type ClaimAuthorizeKSignSecp256k1Msg struct {
	ClaimAuthorizeKSignSecp256k1 core.ClaimAuthorizeKSignSecp256k1
	Signature                    string
	KSignP256                    *ecdsa.PublicKey
}

// SetRootMsg contains the data to set the SetRootClaim with its signature in Hex
type SetRootMsg struct {
	Root      string
	IdAddr    string
	KSignPk   *common3.PublicKey
	Timestamp uint64
	Signature string
}

// ClaimValueMsg contains a core.ClaimValue with its signature in Hex
type ClaimValueMsg struct {
	ClaimValue merkletree.Entry
	Signature  string
	KSignPk    *common3.PublicKey
}

// ProofOfTreeLeaf contains all the parameters needed to proof that a Leaf is in a merkletree with a given Root
type ProofOfTreeLeaf struct {
	Leaf  []byte
	Proof []byte
	Root  merkletree.Hash
}

// ProofOfTreeLeafHex is the same data structure than ProofOfTreeLeaf but in Hexadecimal string representation
type ProofOfTreeLeafHex struct {
	Leaf  string
	Proof string
	Root  string
}

func (plh *ProofOfTreeLeafHex) Unhex() ProofOfTreeLeaf {
	var r ProofOfTreeLeaf
	r.Leaf, _ = common3.HexToBytes(plh.Leaf)
	r.Proof, _ = common3.HexToBytes(plh.Proof)
	rootBytes, _ := common3.HexToBytes(plh.Root)
	copy(r.Root[:], rootBytes[:32])
	return r
}

// Hex returns a ProofOfTreeLeafHex data structure
func (pl *ProofOfTreeLeaf) Hex() ProofOfTreeLeafHex {
	r := ProofOfTreeLeafHex{
		common3.BytesToHex(pl.Leaf),
		common3.BytesToHex(pl.Proof),
		pl.Root.Hex(),
	}
	return r
}

// ProofOfClaimUser is the proof of a claim in the Identity MerkleTree, and the SetRootClaim of that MerkleTree inside the Relay's MerkleTree. Also with the proofs of non revocation of both claims
type ProofOfClaimUser struct {
	ClaimProof                     ProofOfTreeLeaf
	SetRootClaimProof              ProofOfTreeLeaf
	ClaimNonRevocationProof        ProofOfTreeLeaf
	SetRootClaimNonRevocationProof ProofOfTreeLeaf
	Date                           uint64
	Signature                      []byte // signature of the Root of the Relay
}

type ProofOfClaimUserHex struct {
	ClaimProof                     ProofOfTreeLeafHex
	SetRootClaimProof              ProofOfTreeLeafHex
	ClaimNonRevocationProof        ProofOfTreeLeafHex
	SetRootClaimNonRevocationProof ProofOfTreeLeafHex
	Date                           uint64
	Signature                      string // signature of the Root of the Relay
}

func (pc *ProofOfClaimUser) Hex() ProofOfClaimUserHex {
	r := ProofOfClaimUserHex{
		pc.ClaimProof.Hex(),
		pc.SetRootClaimProof.Hex(),
		pc.ClaimNonRevocationProof.Hex(),
		pc.SetRootClaimNonRevocationProof.Hex(),
		pc.Date,
		common3.BytesToHex(pc.Signature),
	}
	return r
}
func (pch *ProofOfClaimUserHex) Unhex() (ProofOfClaimUser, error) {
	sigBytes, err := common3.HexToBytes(pch.Signature)
	if err != nil {
		return ProofOfClaimUser{}, err
	}
	r := ProofOfClaimUser{
		pch.ClaimProof.Unhex(),
		pch.SetRootClaimProof.Unhex(),
		pch.ClaimNonRevocationProof.Unhex(),
		pch.SetRootClaimNonRevocationProof.Unhex(),
		pch.Date,
		sigBytes,
	}
	return r, nil
}

type ProofOfClaimPartial struct {
	Mtp0 *merkletree.Proof
	Mtp1 *merkletree.Proof
	Root *merkletree.Hash
	Aux  *SetRootAux
}

func (p *ProofOfClaimPartial) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Mtp0": hex.EncodeToString(p.Mtp0.Bytes()),
		"Mtp1": hex.EncodeToString(p.Mtp1.Bytes()),
		"Root": hex.EncodeToString(p.Root[:]),
		"Aux":  p.Aux,
	})
}

type SetRootAux struct {
	Version uint32
	Era     uint32
	EthAddr common.Address
}

// ProofOfClaim is the proof of a claim in the Relay's MerkleTree, and the proof of non revocation of the claim
type ProofOfClaim struct {
	//ClaimProof              ProofOfTreeLeaf
	//ClaimNonRevocationProof ProofOfTreeLeaf
	Proofs    []ProofOfClaimPartial
	Leaf      *merkletree.Data
	Date      uint64
	Signature []byte // signature of the Root of the Relay
}

func (p *ProofOfClaim) MarshalJSON() ([]byte, error) {
	leafBytes := p.Leaf.Bytes()
	return json.Marshal(map[string]interface{}{
		"Proofs":    p.Proofs,
		"Leaf":      hex.EncodeToString(leafBytes[:]),
		"Date":      p.Date,
		"Signature": hex.EncodeToString(p.Signature),
	})
}

//type ProofOfClaimHex struct {
//	ClaimProof              ProofOfTreeLeafHex
//	ClaimNonRevocationProof ProofOfTreeLeafHex
//	Date                    uint64
//	Signature               string // signature of the Root of the Relay
//}

//func (pc *ProofOfClaim) Hex() ProofOfClaimHex {
//	r := ProofOfClaimHex{
//		pc.ClaimProof.Hex(),
//		pc.ClaimNonRevocationProof.Hex(),
//		pc.Date,
//		common3.BytesToHex(pc.Signature),
//	}
//	return r
//}
//func (pch *ProofOfClaimHex) Unhex() (ProofOfClaim, error) {
//	sigBytes, err := common3.HexToBytes(pch.Signature)
//	if err != nil {
//		return ProofOfClaim{}, err
//	}
//	r := ProofOfClaim{
//		pch.ClaimProof.Unhex(),
//		pch.ClaimNonRevocationProof.Unhex(),
//		pch.Date,
//		sigBytes,
//	}
//	return r, nil
//}
