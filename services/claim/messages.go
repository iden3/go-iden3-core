package claimsrv

import (
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
)

// BytesSignedMsg contains the value and its signature in Hex representation
type BytesSignedMsg struct {
	ValueHex     string `json:"valueHex"` // claim.Bytes() in a hex format
	SignatureHex string `json:"signatureHex"`
}

// ClaimDefaultMsg contains a core.ClaimDefault with its signature in Hex
type ClaimDefaultMsg struct {
	ClaimDefault core.ClaimDefault
	Signature    string
}

// AssignNameClaimMsg contains a core.AssignNameClaim with its signature in Hex
type AssignNameClaimMsg struct {
	AssignNameClaim core.AssignNameClaim
	Signature       string
}

// AuthorizeKSignClaimMsg contains a core.AuthorizeKSignClaim with its signature in Hex
type AuthorizeKSignClaimMsg struct {
	AuthorizeKSignClaim core.AuthorizeKSignClaim
	Signature           string
}

// SetRootClaimMsg contains a core.SetRootClaim with its signature in Hex
type SetRootClaimMsg struct {
	SetRootClaim core.SetRootClaim
	Signature    string
}

// ClaimValueMsg contains a core.ClaimValue with its signature in Hex
type ClaimValueMsg struct {
	ClaimValue merkletree.Value
	Signature  string
}
