package merkletree

import "bytes"

// Node is the data structure of an intermediate node of the Merkle Tree
type treeNode struct {
	ChildL Hash // hash of the left child
	ChildR Hash // hash of the right child
}

// Bytes returns an array of bytes with the Node data
func (n *treeNode) Bytes() (b []byte) {
	b = append(b, n.ChildL[:]...)
	b = append(b, n.ChildR[:]...)
	return b
}

// Ht returns the hash of the full node
func (n *treeNode) Ht() Hash {
	h := HashBytes(n.Bytes())
	return h
}

// ParseNodeBytes returns a Node struct from an array of bytes
func parseNodeBytes(b []byte) treeNode {
	if bytes.Equal(b, EmptyNodeValue[:]) {
		var node treeNode
		// TODO version
		node.ChildL = EmptyNodeValue
		node.ChildR = EmptyNodeValue
		return node
	}
	var node treeNode
	copy(node.ChildL[:], b[:32])
	copy(node.ChildR[:], b[32:])
	return node
}
