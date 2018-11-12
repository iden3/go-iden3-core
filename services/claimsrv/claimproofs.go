package claimsrv

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

// CheckProofOfClaim checks the Merkle Proof of the Claim, the SetRootClaim, and the non revocation proof of both claims
func CheckProofOfClaim(relayAddr common.Address, pc ProofOfClaim, numLevels int) bool {
	hiClaim := core.HiFromClaimBytes(pc.ClaimProof.Leaf)
	vClaimProof := merkletree.CheckProof(pc.ClaimProof.Root, pc.ClaimProof.Proof,
		hiClaim, merkletree.HashBytes(pc.ClaimProof.Leaf), numLevels)

	hiSetRootClaim := core.HiFromClaimBytes(pc.SetRootClaimProof.Leaf)
	vSetRootClaimProof := merkletree.CheckProof(pc.SetRootClaimProof.Root, pc.SetRootClaimProof.Proof,
		hiSetRootClaim, merkletree.HashBytes(pc.SetRootClaimProof.Leaf), numLevels)

	hiNonRevocationClaim := core.HiFromClaimBytes(pc.ClaimNonRevocationProof.Leaf)
	vClaimNonRevocationProof := merkletree.CheckProof(pc.ClaimNonRevocationProof.Root, pc.ClaimNonRevocationProof.Proof,
		hiNonRevocationClaim, merkletree.EmptyNodeValue, numLevels)

	hiNonRevocationSetRootClaim := core.HiFromClaimBytes(pc.SetRootClaimNonRevocationProof.Leaf)
	vSetRootClaimNonRevocationProof := merkletree.CheckProof(pc.SetRootClaimNonRevocationProof.Root, pc.SetRootClaimNonRevocationProof.Proof,
		hiNonRevocationSetRootClaim, merkletree.EmptyNodeValue, numLevels)

	// additional, check caducity of the pc.Date

	// check signature of the ProofOfClaim.SetRootClaimProof.Root with the identity of the Relay
	// checking this Root and the four Merkle Proofs, we check the full ProofOfClaim
	dateBytes, err := core.Uint64ToEthBytes(pc.Date)
	if err != nil {
		return false
	}
	rootdate := pc.SetRootClaimProof.Root[:]
	rootdate = append(rootdate, dateBytes...)
	rootdateHash := merkletree.HashBytes(rootdate)
	if !utils.VerifySig(relayAddr, pc.Signature, rootdateHash[:]) {
		return false
	}

	if vClaimProof && vSetRootClaimProof && vClaimNonRevocationProof && vSetRootClaimNonRevocationProof {
		return true
	}
	return false
}

// CheckKSignInIDdb checks that a given KSign is in an AuthorizeKSignClaim in the Identity Merkle Tree (in this version, as the Merkle Tree don't allows to delete data, the verification only needs to check if the AuthorizeKSignClaim is in the key-value)
func CheckKSignInIDdb(mt *merkletree.MerkleTree, ksign common.Address) bool {
	// generate the AuthorizeKSignClaim
	authorizeKSignClaim := core.NewOperationalKSignClaim(ksign)
	ht := authorizeKSignClaim.Ht()
	value, err := mt.Storage().Get(ht[:])
	if err != nil {
		return false
	}
	if !bytes.Equal(authorizeKSignClaim.Bytes(), value[5:]) { // value[5:] to skip the db prefix
		return false
	}

	// non revocation
	authorizeKSignClaim.BaseIndex.Version++
	ht = authorizeKSignClaim.Ht()
	value, err = mt.Storage().Get(ht[:])
	if err.Error() != "key not found" {
		return false
	}

	return true
}
