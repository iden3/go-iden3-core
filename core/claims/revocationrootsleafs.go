package claims

import (
	"encoding/binary"

	"github.com/iden3/go-iden3-core/merkletree"
)

// LeafRootsTree contains the root to be inserted in the leaf
type LeafRootsTree struct {
	Root merkletree.Hash
}

// NewLeafRootsTree returns a LeafRootsTree with the provided root.
func NewLeafRootsTree(root merkletree.Hash) *LeafRootsTree {
	return &LeafRootsTree{
		Root: root,
	}
}

// NewLeafRootsTreeFromEntry deserializes a LeafRootsTree from an Entry.
func NewLeafRootsTreeFromEntry(e *merkletree.Entry) *LeafRootsTree {
	l := &LeafRootsTree{}
	l.Root = merkletree.Hash(e.Data[0])
	return l
}

// Entry serializes the leaf into an Entry.
func (l *LeafRootsTree) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	e.Data[0] = merkletree.ElemBytes(l.Root)
	return e
}

// LeafRevocationsTree contains the root to be inserted in the leaf
type LeafRevocationsTree struct {
	Nonce   uint32
	Version uint32
}

// NewLeafRevocationsTree returns a LeafRevocationsTree with the provided root.
func NewLeafRevocationsTree(nonce, version uint32) *LeafRevocationsTree {
	return &LeafRevocationsTree{
		Nonce:   nonce,
		Version: version,
	}
}

// NewLeafRevocationsTreeFromEntry deserializes a LeafRevocationsTree from an Entry.
func NewLeafRevocationsTreeFromEntry(e *merkletree.Entry) *LeafRevocationsTree {
	l := &LeafRevocationsTree{}
	l.Nonce = binary.BigEndian.Uint32(e.Data[0].Bytes()[:4])
	l.Version = binary.BigEndian.Uint32(e.Data[4].Bytes()[:4])
	return l
}

// Entry serializes the leaf into an Entry.
func (l *LeafRevocationsTree) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	binary.BigEndian.PutUint32(e.Data[0].Bytes()[:4], l.Nonce)
	binary.BigEndian.PutUint32(e.Data[4].Bytes()[:4], l.Version)
	return e
}

// AddLeafRootsTree adds a new leaf to the given MerkleTree, which contains the Root
func AddLeafRootsTree(mt *merkletree.MerkleTree, root *merkletree.Hash) error {
	l := NewLeafRootsTree(*root)
	return mt.AddEntry(l.Entry())
}

// AddLeafRevocationsTree adds a new leaf to the given MerkleTree, which contains the Nonce & Version
func AddLeafRevocationsTree(mt *merkletree.MerkleTree, nonce, version uint32) error {
	l := NewLeafRevocationsTree(nonce, version)
	return mt.AddEntry(l.Entry())
}
