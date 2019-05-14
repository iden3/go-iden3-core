package merkletree

import (
	"fmt"
	"math/big"
)

// NodeType defines the type of node in the MT.
type NodeType byte

const (
	// NodeTypeMiddle indicates the type of middle Node that has children.
	NodeTypeMiddle NodeType = 0
	// NodeTypeLeaf indicates the type of a leaf Node that contains a claim.
	NodeTypeLeaf NodeType = 1
	// NodeTypeEmpty indicates the type of an empty Node.
	NodeTypeEmpty NodeType = 2

	// DBEntryTypeRoot indicates the type of a DB entry that indicates the current Root of a MerkleTree
	DBEntryTypeRoot NodeType = 3
)

// Node is the struct that represents a node in the MT. The node should not be
// modified after creation because the cached key won't be updated.
type Node struct {
	// Type is the type of node in the tree.
	Type NodeType
	// ChildL is the left child of a middle node.
	ChildL *Hash
	// ChildR is the right child of a middle node.
	ChildR *Hash
	// Entry is the data stored in a leaf node.
	Entry *Entry
	// key is a cache used to avoid recalculating key
	key *Hash
}

// NewNodeLeaf creates a new leaf node.
func NewNodeLeaf(entry *Entry) *Node {
	return &Node{Type: NodeTypeLeaf, Entry: entry}
}

// NewNodeMiddle creates a new middle node.
func NewNodeMiddle(childL *Hash, childR *Hash) *Node {
	return &Node{Type: NodeTypeMiddle, ChildL: childL, ChildR: childR}
}

// NewNodeEmpty creates a new empty node.
func NewNodeEmpty() *Node {
	return &Node{Type: NodeTypeEmpty}
}

// NewNodeFromBytes creates a new node by parsing the input []byte.
func NewNodeFromBytes(b []byte) (*Node, error) {
	if len(b) < 1 {
		return nil, ErrNodeDataBadSize
	}
	n := Node{Type: NodeType(b[0])}
	b = b[1:]
	switch n.Type {
	case NodeTypeMiddle:
		if len(b) != 2*ElemBytesLen {
			return nil, ErrNodeDataBadSize
		}
		n.ChildL, n.ChildR = &Hash{}, &Hash{}
		copy(n.ChildL[:], b[:ElemBytesLen])
		copy(n.ChildR[:], b[ElemBytesLen:ElemBytesLen*2])
	case NodeTypeLeaf:
		if len(b) != 4*ElemBytesLen {
			return nil, ErrNodeDataBadSize
		}
		n.Entry = &Entry{}
		for i := 0; i < DataLen; i++ {
			copy(n.Entry.Data[i][:], b[i*ElemBytesLen:(i+1)*ElemBytesLen])
		}
	case NodeTypeEmpty:
		break
	default:
		return nil, ErrInvalidNodeFound
	}
	return &n, nil
}

// LeafKey computes the key of a leaf node given the hIndex and hValue of the
// entry of the leaf.
func LeafKey(hIndex, hValue *Hash) *Hash {
	// return HashElems(ElemBytesOne, ElemBytes(*hIndex), ElemBytes(*hValue))
	return HashElemsKey(big.NewInt(1), ElemBytes(*hIndex), ElemBytes(*hValue))
}

// Key computes the key of the node by hashing the content in a specific way
// for each type of node.  This key is used as the hash of the merklee tree for
// each node.
func (n *Node) Key() *Hash {
	if n.key == nil { // Cache the key to avoid repeated hash computations.
		// NOTE: We are not using the type to calculate the hash!
		switch n.Type {
		case NodeTypeMiddle: // H(ChildL || ChildR)
			n.key = HashElems(ElemBytes(*n.ChildL), ElemBytes(*n.ChildR))
		case NodeTypeLeaf: // H(Data...)
			n.key = LeafKey(n.Entry.HIndex(), n.Entry.HValue())
		case NodeTypeEmpty: // Zero
			n.key = &HashZero
		default:
			n.key = &HashZero
		}
	}
	return n.key
}

// Value returns the value of the node.  This is the content that is stored in the backend database.
func (n *Node) Value() []byte {
	switch n.Type {
	case NodeTypeMiddle: // {Type || ChildL || ChildR}
		return append([]byte{byte(n.Type)}, append(n.ChildL[:], n.ChildR[:]...)...)
	case NodeTypeLeaf: // {Type || Data...}
		return append([]byte{byte(n.Type)}, ElemsBytesToBytes(n.Entry.Data[:])...)
	case NodeTypeEmpty: // {}
		return []byte{}
	default:
		return []byte{}
	}
}

// String outputs a string representation of a node (different for each type).
func (n *Node) String() string {
	switch n.Type {
	case NodeTypeMiddle: // {Type || ChildL || ChildR}
		return fmt.Sprintf("Middle L:%s R:%s", n.ChildL, n.ChildR)
	case NodeTypeLeaf: // {Type || Data...}
		return fmt.Sprintf("Leaf I:%s D:%s", n.Entry.Data[1], n.Entry.Data[3])
	case NodeTypeEmpty: // {}
		return "Empty"
	default:
		return "Invalid Node"
	}
}
