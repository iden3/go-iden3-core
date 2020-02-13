package idenpuboffchain

import (
	"bytes"
	"fmt"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
)

var (
	ErrCalculatedIdenStateDoesntMatch = fmt.Errorf("Calculated IdenState doesn't match the one in PublicDataBlobs")
)

// IdenPubOffChainReader is a interface to read the off chain public state of an identity.
type IdenPubOffChainReader interface {
	GetPublicData(idenPubUrl string, id *core.ID, idenState *merkletree.Hash) (*PublicData, error)
}

// IdenPubOffChainWriter is a interface to write the off chain public state of an identity.
type IdenPubOffChainWriter interface {
	Publish(publicData *PublicData) error
}

// PublicDataBlobs contains the RootsTree (blob) + Root, and the RevocationTree (blob) + Root
type PublicDataBlobs struct {
	IdenState           merkletree.Hash
	ClaimsTreeRoot      merkletree.Hash
	RevocationsTreeRoot merkletree.Hash
	RevocationsTree     []byte
	RootsTreeRoot       merkletree.Hash
	RootsTree           []byte
}

// PublicData contains the IdenState, ClaimsRoot, RootsTree and RevocationsTree
type PublicData struct {
	IdenState           *merkletree.Hash
	ClaimsTreeRoot      *merkletree.Hash
	RevocationsTreeRoot *merkletree.Hash
	RevocationsTree     *merkletree.MerkleTree
	RootsTreeRoot       *merkletree.Hash
	RootsTree           *merkletree.MerkleTree
}

// NewPublicDataFromBlobs builds the revocation tree and the roots tree in
// memory storage.  It also checks the validity of the tree roots against the
// identity state.
func NewPublicDataFromBlobs(publicDataBlobs *PublicDataBlobs) (*PublicData, error) {
	revocationsTree, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	if err != nil {
		return nil, err
	}
	rootsTree, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	if err != nil {
		return nil, err
	}

	if err = revocationsTree.ImportTree(bytes.NewReader(publicDataBlobs.RevocationsTree)); err != nil {
		return nil, err
	}
	if !revocationsTree.RootKey().Equals(&publicDataBlobs.RevocationsTreeRoot) {
		return nil, fmt.Errorf("Imported revocations tree root (%v) doesn't match the expected root (%v)",
			revocationsTree.RootKey(), publicDataBlobs.RevocationsTreeRoot)
	}

	if err = rootsTree.ImportTree(bytes.NewReader(publicDataBlobs.RootsTree)); err != nil {
		return nil, err
	}
	if !rootsTree.RootKey().Equals(&publicDataBlobs.RootsTreeRoot) {
		return nil, fmt.Errorf("Imported roots tree root (%v) doesn't match the expected root (%v)",
			rootsTree.RootKey(), publicDataBlobs.RootsTreeRoot)
	}
	idenState := core.IdenState(&publicDataBlobs.ClaimsTreeRoot,
		&publicDataBlobs.RevocationsTreeRoot, &publicDataBlobs.RootsTreeRoot)
	if !idenState.Equals(&publicDataBlobs.IdenState) {
		return nil, ErrCalculatedIdenStateDoesntMatch
	}
	return &PublicData{
		IdenState:           &publicDataBlobs.IdenState,
		ClaimsTreeRoot:      &publicDataBlobs.ClaimsTreeRoot,
		RevocationsTreeRoot: &publicDataBlobs.RevocationsTreeRoot,
		RevocationsTree:     revocationsTree,
		RootsTreeRoot:       &publicDataBlobs.RootsTreeRoot,
		RootsTree:           rootsTree,
	}, nil
}
