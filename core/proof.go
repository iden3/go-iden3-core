package core

import "github.com/iden3/go-iden3/merkletree"

// MerkleProof is the data structure of the proof of a leaf in the Merkle Tree
type MerkleProof struct {
	Root     merkletree.Hash
	Proof    []byte
	Leaf     []byte          // claim
	LeafHash merkletree.Hash // claim hash
}

// ClaimWithProof is the data structure of the needed MerkleProofs to proof a Claim
type ClaimWithProof struct {
	// claim and merkleproof of the Claim in the user identity's merkletree
	Claim          []byte      // the claim in bytes format
	ExistenceProof MerkleProof // exists
	SoundnessProof MerkleProof // non-revocated

	// merkleproof of the IDRootClaim of that user identity in the Relay merkletree
	RootSignature   []byte // signature of the root
	RelayRootsProof MerkleProof
}
