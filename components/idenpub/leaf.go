package idenpub

import (
	"encoding/binary"

	"github.com/iden3/go-iden3-core/merkletree"
)

// LeafRoT contains the root to be inserted in the leaf
type LeafRoT struct {
	Root merkletree.Hash
}

// NewLeafRoT returns a LeafRoT with the provided root.
func NewLeafRoT(root merkletree.Hash) *LeafRoT {
	var r [32]byte
	copy(r[:], root[:31]) // cropped to fit inside the FF
	return &LeafRoT{
		Root: r,
	}
}

// NewLeafRoTFromEntry deserializes a LeafRoT from an Entry.
func NewLeafRoTFromEntry(e *merkletree.Entry) *LeafRoT {
	l := &LeafRoT{}
	l.Root = merkletree.Hash(e.Data[0])
	return l
}

// Entry serializes the leaf into an Entry.
func (l *LeafRoT) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	e.Data[0] = merkletree.ElemBytes(l.Root)
	return e
}

// LeafReT contains the root to be inserted in the leaf
type LeafReT struct {
	Nonce   uint32
	Version uint32
}

// NewLeafReT returns a LeafReT with the provided root.
func NewLeafReT(nonce, version uint32) *LeafReT {
	return &LeafReT{
		Nonce:   nonce,
		Version: version,
	}
}

// NewLeafReTFromEntry deserializes a LeafReT from an Entry.
func NewLeafReTFromEntry(e *merkletree.Entry) *LeafReT {
	l := &LeafReT{}
	l.Nonce = binary.BigEndian.Uint32(e.Data[0][:4])
	l.Version = binary.BigEndian.Uint32(e.Data[4][:4])
	return l
}

// Entry serializes the leaf into an Entry.
func (l *LeafReT) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	binary.BigEndian.PutUint32(e.Data[0][:4], l.Nonce)
	binary.BigEndian.PutUint32(e.Data[4][:4], l.Version)
	return e
}
