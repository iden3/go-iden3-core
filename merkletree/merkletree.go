package merkletree

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/iden3/go-iden3/db"
)

// TODO: Remove
// Value is a placeholder interface of a generic claim, a key value object
// stored in the leveldb, for the transition to the new Merklee Tree.
type Value interface {
	IndexLength() uint32 // returns the index length value
	Bytes() []byte       // returns the value in byte array representation
}

// ElemBytes is the basic type used to store data in the MT.  ElemBytes
// corresponds to the serialization of an element from mimc7.
type ElemBytes [ElemBytesLen]byte

// String returns the last 4 bytes of ElemBytes in hex.
func (e *ElemBytes) String() string {
	return hex.EncodeToString(e[ElemBytesLen-4:])
}

// ElemsBytesToBytes serializes an array of ElemBytes to []byte.
func ElemsBytesToBytes(es []ElemBytes) []byte {
	bs := make([]byte, len(es)*ElemBytesLen)
	for i := 0; i < len(es); i++ {
		copy(bs[i*ElemBytesLen:(i+1)*ElemBytesLen], es[i][:])
	}
	return bs
}

// Index is the type used to represent the index of an entry in the MT,
// used to find the path from the root to the leaf that contains such entry.
type Index [IndexLen]ElemBytes

// Data is the type used to represent the data stored in an entry of the MT.
// It consists of 4 elements: e0, e1, e2, e3;
// where v = [e0,e1], index = [e2,e3].
type Data [DataLen]ElemBytes

func (d *Data) String() string {
	return fmt.Sprintf("%s%s%s%s", hex.EncodeToString(d[0][:]), hex.EncodeToString(d[1][:]),
		hex.EncodeToString(d[2][:]), hex.EncodeToString(d[3][:]))
}

const (
	// ElemBytesLen is the length in bytes of each element used for storing
	// data and hashing.
	ElemBytesLen = 32
	// IndexLen indicates how many elements are used for the index.
	IndexLen = 2
	// DataLen indicates how many elements are used for the data.
	DataLen = 4
)

var (
	// ErrNodeAlreadyExists is used when a node already exists.
	ErrNodeAlreadyExists = errors.New("node already exists")
	// ErrEntryIndexNotFound is used when no entry is found for an index.
	ErrEntryIndexNotFound = errors.New("node index not found in the DB")
	// ErrNodeDataBadSize is used when the data of a node has an incorrect
	// size and can't be parsed.
	ErrNodeDataBadSize = errors.New("node data has incorrect size in the DB")
	// ErrReachedMaxLevel is used when a traversal of the MT reaches the
	// maximum level.
	ErrReachedMaxLevel = errors.New("reached maximum level of the merkle tree")
	// ErrInvalidNodeFound is used when an invalid node is found and can't
	// be parsed.
	ErrInvalidNodeFound = errors.New("found an invalid node in the DB")
	// ErrInvalidProofBytes is used when a serialized proof is invalid.
	ErrInvalidProofBytes = errors.New("the serialized proof is invalid")
	// HashZero is a hash value of zeros, and is the key of an empty node.
	HashZero = Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	// ElemBytesOne is a constant element used as a prefix to compute leaf node keys.
	ElemBytesOne = ElemBytes{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
)

// Entry is the generic type that is stored in the MT.
type Entry struct {
	Data   Data
	hIndex *Hash
	hValue *Hash
}

type Entrier interface {
	ToEntry() Entry
}

// HIndex calculates the hash of the Index of the entry, used to find the path
// from the root to the leaf in the MT.
func (e *Entry) HIndex() *Hash {
	if e.hIndex == nil { // Cache the hIndex.
		//e.hIndex = HashElems(e.Index()[:]...)
		e.hIndex = HashElems(e.Data[2:]...)
	}
	return e.hIndex
}

func (e *Entry) HValue() *Hash {
	if e.hValue == nil { // Cache the hIndex.
		e.hValue = HashElems(e.Data[:2]...)
	}
	return e.hValue
}

//MerkleTree is the struct with the main elements of the Merkle Tree
type MerkleTree struct {
	sync.RWMutex
	// storage is the backend database.
	storage db.Storage
	// rootKey is the key of the root node.
	rootKey *Hash
	// maxLevels is the maximum number of levels of the Merkle Tree.
	maxLevels int
}

// NewMerkleTree generates a new Merkle Tree
func NewMerkleTree(storage db.Storage, maxLevels int) (*MerkleTree, error) {
	mt := MerkleTree{storage: storage, maxLevels: maxLevels}
	nodeRoot := NewNodeEmpty()
	k, _ := nodeRoot.Key(), nodeRoot.Value()
	mt.rootKey = k
	return &mt, nil
}

// Storage returns the MT storage
func (mt *MerkleTree) Storage() db.Storage {
	return mt.storage
}

// RootKey returns the MT root node key
func (mt *MerkleTree) RootKey() *Hash {
	return mt.rootKey
}

// MaxLevels returns the MT maximum level
func (mt *MerkleTree) MaxLevels() int {
	return mt.maxLevels
}

// GetDataByIndex returns the data from the MT in the position of the hash of
// the index (hIndex)
func (mt *MerkleTree) GetDataByIndex(hIndex *Hash) (*Data, error) {
	path := getPath(mt.maxLevels, hIndex)
	nextKey := mt.rootKey
	for lvl := 0; lvl < mt.maxLevels; lvl++ {
		n, err := mt.GetNode(nextKey)
		if err != nil {
			return nil, err
		}
		switch n.Type {
		case NodeTypeEmpty:
			return nil, ErrEntryIndexNotFound
		case NodeTypeLeaf:
			if bytes.Equal(hIndex[:], n.Entry.HIndex()[:]) {
				return &n.Entry.Data, nil
			} else {
				return nil, ErrEntryIndexNotFound
			}
		case NodeTypeMiddle:
			if path[lvl] {
				nextKey = n.ChildR
			} else {
				nextKey = n.ChildL
			}
		default:
			return nil, ErrInvalidNodeFound
		}
	}
	return nil, ErrEntryIndexNotFound
}

// pushLeaf recursively pushes an existing oldLeaf down until its path diverges
// from newLeaf, at which point both leafs are stored, all while updating the
// path.
func (mt *MerkleTree) pushLeaf(tx db.Tx, newLeaf *Node, oldLeaf *Node,
	lvl int, pathNewLeaf []bool, pathOldLeaf []bool) (*Hash, error) {
	if lvl > mt.maxLevels-2 {
		return nil, ErrReachedMaxLevel
	}
	var newNodeMiddle *Node
	if pathNewLeaf[lvl] == pathOldLeaf[lvl] { // We need to go deeper!
		nextKey, err := mt.pushLeaf(tx, newLeaf, oldLeaf, lvl+1, pathNewLeaf, pathOldLeaf)
		if err != nil {
			return nil, err
		}
		if pathNewLeaf[lvl] {
			newNodeMiddle = NewNodeMiddle(&HashZero, nextKey) // go right
		} else {
			newNodeMiddle = NewNodeMiddle(nextKey, &HashZero) // go left
		}
		return mt.AddNode(tx, newNodeMiddle)
	} else {
		if pathNewLeaf[lvl] {
			newNodeMiddle = NewNodeMiddle(oldLeaf.Key(), newLeaf.Key())
		} else {
			newNodeMiddle = NewNodeMiddle(newLeaf.Key(), oldLeaf.Key())
		}
		// We can add newLeaf now.  We don't need to add oldLeaf because it's already in the tree.
		_, err := mt.AddNode(tx, newLeaf)
		if err != nil {
			return nil, err
		}
		return mt.AddNode(tx, newNodeMiddle)
	}
}

// addLeaf recursively adds a newLeaf in the MT while updating the path.
func (mt *MerkleTree) addLeaf(tx db.Tx, newLeaf *Node, key *Hash,
	lvl int, path []bool) (*Hash, error) {
	var err error
	var nextKey *Hash
	if lvl > mt.maxLevels-1 {
		return nil, ErrReachedMaxLevel
	}
	n, err := mt.GetNode(key)
	if err != nil {
		return nil, err
	}
	switch n.Type {
	case NodeTypeEmpty:
		// We can add newLeaf now
		return mt.AddNode(tx, newLeaf)
	case NodeTypeLeaf:
		// TODO: delete old node n???  Make this optional???
		hIndex := n.Entry.HIndex()
		pathOldLeaf := getPath(mt.maxLevels, hIndex)
		// We need to push newLeaf down until its path diverges from n's path
		return mt.pushLeaf(tx, newLeaf, n, lvl, path, pathOldLeaf)
	case NodeTypeMiddle:
		// We need to go deeper, continue traversing the tree, left or right depending on path
		var newNodeMiddle *Node
		if path[lvl] {
			nextKey, err = mt.addLeaf(tx, newLeaf, n.ChildR, lvl+1, path) // go right
			newNodeMiddle = NewNodeMiddle(n.ChildL, nextKey)
		} else {
			nextKey, err = mt.addLeaf(tx, newLeaf, n.ChildL, lvl+1, path) // go left
			newNodeMiddle = NewNodeMiddle(nextKey, n.ChildR)
		}
		if err != nil {
			return nil, err
		}
		// TODO: delete old node n???  Make this optional???
		// Update the node to reflect the modified child
		return mt.AddNode(tx, newNodeMiddle)
	default:
		return nil, ErrInvalidNodeFound
	}
}

// Add adds the Entry to the MerkleTree
func (mt *MerkleTree) Add(e *Entry) error {
	// First of all, verfy that the ElemBytes are valid and fit inside the
	// mimc7 field.
	_, err := ElemsBytesToRElems(e.Data[:]...)
	if err != nil {
		return err
	}
	tx, err := mt.storage.NewTx()
	if err != nil {
		return err
	}
	mt.Lock()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Close()
		}
		mt.Unlock()
	}()

	newNodeLeaf := NewNodeLeaf(e)
	hIndex := e.HIndex()
	path := getPath(mt.maxLevels, hIndex)

	newRootKey, err := mt.addLeaf(tx, newNodeLeaf, mt.rootKey, 0, path)
	if err != nil {
		return err
	}
	mt.rootKey = newRootKey
	return nil
}

// graphViz is a helper recursive function to output the tree in GraphViz syntax.
func (mt *MerkleTree) graphViz(w io.Writer, key *Hash, cnt *int) error {
	n, err := mt.GetNode(key)
	if err != nil {
		return err
	}
	switch n.Type {
	case NodeTypeEmpty:
	case NodeTypeLeaf:
		fmt.Fprintf(w, "\"%v\" [style=filled];\n", n.Key())
	case NodeTypeMiddle:
		lr := [2]string{n.ChildL.String(), n.ChildR.String()}
		for i, _ := range lr {
			if lr[i] == "00000000" {
				lr[i] = fmt.Sprintf("empty%v", *cnt)
				fmt.Fprintf(w, "\"%v\" [style=dashed,label=0];\n", lr[i])
				(*cnt)++
			}
		}
		fmt.Fprintf(w, "\"%v\" -> {\"%v\" \"%v\"}\n", n.Key(), lr[0], lr[1])
		mt.graphViz(w, n.ChildL, cnt)
		mt.graphViz(w, n.ChildR, cnt)
	default:
		return ErrInvalidNodeFound
	}
	return nil
}

// GraphViz generates a string GraphViz representation of the tree and writes
// it to w.
func (mt *MerkleTree) GraphViz(w io.Writer) error {
	fmt.Fprintf(w, `digraph hierarchy {
node [fontname=Monospace,fontsize=10,shape=box]
`)
	cnt := 0
	err := mt.graphViz(w, mt.RootKey(), &cnt)
	fmt.Fprintf(w, "}")
	return err
}

// nodeAux contains the auxiliary node used in a non-existence proof.
type nodeAux struct {
	//key    *Hash
	hIndex *Hash
	hValue *Hash
}

// proofFlagsLen is the byte length of the flags in the proof header (first 32
// bytes).
const proofFlagsLen = 2

// Proof defines the required elements for a MT proof of existence or non-existence.
type Proof struct {
	// existence indicates wether this is a proof of existence or non-existence.
	existence bool
	// depth indicates how deep in the tree the proof goes.
	depth uint
	// notempties is a bitmap of non-empty siblings found in siblings.
	notempties [ElemBytesLen - proofFlagsLen]byte
	// siblings is a list of non-empty sibling keys.
	siblings []*Hash
	nodeAux  *nodeAux
}

// NewProofFromBytes parses a byte array into a Proof.
func NewProofFromBytes(bs []byte) (*Proof, error) {
	if len(bs) < ElemBytesLen {
		return nil, ErrInvalidProofBytes
	}
	p := &Proof{}
	if (bs[0] & 0x01) == 0 {
		p.existence = true
	}
	p.depth = uint(bs[1])
	copy(p.notempties[:], bs[proofFlagsLen:ElemBytesLen])
	siblingBytes := bs[ElemBytesLen:]
	sibIdx := 0
	for i := uint(0); i < p.depth; i++ {
		if testBitBigEndian(p.notempties[:], i) {
			if len(siblingBytes) < (sibIdx+1)*ElemBytesLen {
				return nil, ErrInvalidProofBytes
			}
			var sib Hash
			copy(sib[:], siblingBytes[sibIdx*ElemBytesLen:(sibIdx+1)*ElemBytesLen])
			p.siblings = append(p.siblings, &sib)
			sibIdx++
		}
	}

	if p.existence && ((bs[0] & 0x02) != 0) {
		p.nodeAux = &nodeAux{hIndex: &Hash{}, hValue: &Hash{}}
		nodeAuxBytes := siblingBytes[len(p.siblings)*ElemBytesLen:]
		if len(nodeAuxBytes) != 2*ElemBytesLen {
			return nil, ErrInvalidProofBytes
		}
		copy(p.nodeAux.hIndex[:], nodeAuxBytes[:ElemBytesLen])
		copy(p.nodeAux.hValue[:], nodeAuxBytes[ElemBytesLen:2*ElemBytesLen])
	}
	return p, nil
}

// Bytes serializes a Proof into a byte array.
func (p *Proof) Bytes() []byte {
	bsLen := proofFlagsLen + len(p.notempties) + ElemBytesLen*len(p.siblings)
	if p.nodeAux != nil {
		bsLen += 2 * ElemBytesLen
	}
	bs := make([]byte, bsLen)

	if !p.existence {
		bs[0] |= 0x01
	}
	bs[1] = byte(p.depth)
	copy(bs[proofFlagsLen:len(p.notempties)+proofFlagsLen], p.notempties[:])
	siblingsBytes := bs[len(p.notempties)+proofFlagsLen:]
	for i, k := range p.siblings {
		copy(siblingsBytes[i*ElemBytesLen:(i+1)*ElemBytesLen], k[:])
	}
	if p.nodeAux != nil {
		bs[0] |= 0x02
		copy(bs[len(bs)-2*ElemBytesLen:], p.nodeAux.hIndex[:])
		copy(bs[len(bs)-1*ElemBytesLen:], p.nodeAux.hValue[:])
	}
	return bs
}

// String outputs a multiline string representation of the Proof.
func (p *Proof) String() string {
	buf := bytes.NewBufferString("Proof:")
	fmt.Fprintf(buf, "\texistence: %v\n", p.existence)
	fmt.Fprintf(buf, "\tdepth: %v\n", p.depth)
	fmt.Fprintf(buf, "\tnotempties: ")
	for i := uint(0); i < p.depth; i++ {
		if testBitBigEndian(p.notempties[:], i) {
			fmt.Fprintf(buf, "1")
		} else {
			fmt.Fprintf(buf, "0")
		}
	}
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "\tsiblings: ")
	sibIdx := 0
	for i := uint(0); i < p.depth; i++ {
		if testBitBigEndian(p.notempties[:], i) {
			fmt.Fprintf(buf, "%v ", p.siblings[sibIdx])
			sibIdx++
		} else {
			fmt.Fprintf(buf, "0 ")
		}
	}
	fmt.Fprintf(buf, "\n")
	if p.nodeAux != nil {
		fmt.Fprintf(buf, "\tnode aux: hi: %v, ht: %v\n", p.nodeAux.hIndex, p.nodeAux.hValue)
	}
	return buf.String()
}

// GenerateProof generates the proof of existence (or non-existence) of an
// Entry's hash Index for a Merkle Tree given the root.
func (mt *MerkleTree) GenerateProof(hIndex *Hash) (*Proof, error) {
	p := &Proof{}
	var siblingKey *Hash

	path := getPath(mt.maxLevels, hIndex)
	nextKey := mt.rootKey
	for p.depth = 0; p.depth < uint(mt.maxLevels); p.depth++ {
		n, err := mt.GetNode(nextKey)
		if err != nil {
			return nil, err
		}
		switch n.Type {
		case NodeTypeEmpty:
			return p, nil
		case NodeTypeLeaf:
			if bytes.Equal(hIndex[:], n.Entry.HIndex()[:]) {
				p.existence = true
				return p, nil
			} else {
				// We found a leaf whose entry didn't match hIndex
				p.nodeAux = &nodeAux{hIndex: n.Entry.HIndex(), hValue: n.Entry.HValue()}
				return p, nil
			}
		case NodeTypeMiddle:
			if path[p.depth] {
				nextKey = n.ChildR
				siblingKey = n.ChildL
			} else {
				nextKey = n.ChildL
				siblingKey = n.ChildR
			}
		default:
			return nil, ErrInvalidNodeFound
		}
		if !bytes.Equal(siblingKey[:], HashZero[:]) {
			setBitBigEndian(p.notempties[:], uint(p.depth))
			p.siblings = append(p.siblings, siblingKey)
		}
	}
	return nil, ErrEntryIndexNotFound
}

// VerifyProof verifies the Merkle Proof for the entry and root.
func VerifyProof(rootKey *Hash, proof *Proof, hIndex, hValue *Hash) bool {
	sibIdx := len(proof.siblings) - 1
	var midKey *Hash
	if proof.existence {
		midKey = LeafKey(hIndex, hValue)
	} else {
		if proof.nodeAux == nil {
			midKey = &HashZero
		} else {
			if bytes.Equal(hIndex[:], proof.nodeAux.hIndex[:]) {
				return false
			}
			midKey = LeafKey(proof.nodeAux.hIndex, proof.nodeAux.hValue)
		}
	}
	path := getPath(int(proof.depth), hIndex)
	var siblingKey *Hash
	for lvl := int(proof.depth) - 1; lvl >= 0; lvl-- {
		if testBitBigEndian(proof.notempties[:], uint(lvl)) {
			siblingKey = proof.siblings[sibIdx]
			sibIdx--
		} else {
			siblingKey = &HashZero
		}
		if path[lvl] {
			midKey = NewNodeMiddle(siblingKey, midKey).Key()
		} else {
			midKey = NewNodeMiddle(midKey, siblingKey).Key()
		}
	}
	return bytes.Equal(rootKey[:], midKey[:])
}

// GetNode gets a node by key from the MT.  Empty nodes are not stored in the
// tree; they are all the same and assumed to always exist.
func (mt *MerkleTree) GetNode(key *Hash) (*Node, error) {
	if bytes.Equal(key[:], HashZero[:]) {
		return NewNodeEmpty(), nil
	}
	nBytes, err := mt.storage.Get(key[:])
	if err != nil {
		return nil, err
	}
	return NewNodeFromBytes(nBytes)
}

// AddNode adds a node into the MT.  Empty nodes are not stored in the tree;
// they are all the same and assumed to always exist.
func (mt *MerkleTree) AddNode(tx db.Tx, n *Node) (*Hash, error) {
	if n.Type == NodeTypeEmpty {
		return n.Key(), nil
	}
	k, v := n.Key(), n.Value()
	// Check that the node key doesn't already exist
	if _, err := tx.Get(k[:]); err == nil {
		return nil, ErrNodeAlreadyExists
	}
	tx.Put(k[:], v)
	return k, nil
}
