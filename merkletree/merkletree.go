package merkletree

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/db"
	cryptoConstants "github.com/iden3/go-iden3-crypto/constants"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
)

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

func (d *Data) Bytes() (b [ElemBytesLen * DataLen]byte) {
	for i := 0; i < DataLen; i++ {
		copy(b[i*ElemBytesLen:(i+1)*ElemBytesLen], d[i][:])
	}
	return b
}

func (d1 *Data) Equal(d2 *Data) bool {
	return bytes.Equal(d1[0][:], d2[0][:]) && bytes.Equal(d1[1][:], d2[1][:]) &&
		bytes.Equal(d1[2][:], d2[2][:]) && bytes.Equal(d1[3][:], d2[3][:])
}

func (d *Data) MarshalText() ([]byte, error) {
	dataBytes := d.Bytes()
	return []byte(common3.HexEncode(dataBytes[:])), nil
}

func (d *Data) UnmarshalText(text []byte) error {
	var dataBytes [ElemBytesLen * DataLen]byte
	err := common3.HexDecodeInto(dataBytes[:], text)
	if err != nil {
		return err
	}
	*d = *NewDataFromBytes(dataBytes)
	return nil
}

func NewDataFromBytes(b [ElemBytesLen * DataLen]byte) *Data {
	d := &Data{}
	for i := 0; i < DataLen; i++ {
		copy(d[i][:], b[i*ElemBytesLen : (i+1)*ElemBytesLen][:])
	}
	return d
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
	// ErrNodeKeyAlreadyExists is used when a node key already exists.
	ErrNodeKeyAlreadyExists = errors.New("node already exists")
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
	// ErrInvalidDBValue is used when a value in the key value DB is
	// invalid (for example, it doen't contain a byte header and a []byte
	// body of at least len=1.
	ErrInvalidDBValue = errors.New("the value in the DB is invalid")
	// ErrEntryIndexAlreadyExists is used when the entry index already
	// exists in the tree.
	ErrEntryIndexAlreadyExists = errors.New("the entry index already exists in the tree")
	// ErrNotWritable is used when the MerkleTree is not writable and a write function is called
	ErrNotWritable = errors.New("Merkle Tree not writable")
	// HashZero is a hash value of zeros, and is the key of an empty node.
	HashZero = Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	// ElemBytesOne is a constant element used as a prefix to compute leaf node keys.
	ElemBytesOne = ElemBytes{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	// rootNodeVValue is the Key used to store the current Root in the database
	rootNodeValue = []byte("currentroot")
)

// Entry is the generic type that is stored in the MT.  The entry should not be
// modified after creating because the cached hIndex and hValue won't be
// updated.
type Entry struct {
	Data Data
	// hIndex is a cache used to avoid recalculating hIndex
	hIndex *Hash
	// hValue is a cache used to avoid recalculating hValue
	hValue *Hash
}

type Claim interface {
	Entry() *Entry
}

func NewEntryFromBytes(b []byte) (*Entry, error) {
	if len(b) != ElemBytesLen*DataLen {
		return nil, fmt.Errorf("Invalid length for Entry Data")
	}
	var data [ElemBytesLen * DataLen]byte
	copy(data[:], b)
	return &Entry{Data: *NewDataFromBytes(data)}, nil
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
	if e.hValue == nil { // Cache the hValue.
		e.hValue = HashElems(e.Data[:2]...)
	}
	return e.hValue
}

func (e *Entry) Bytes() []byte {
	b := e.Data.Bytes()
	return b[:]
}

func (e1 *Entry) Equal(e2 *Entry) bool {
	return e1.Data.Equal(&e2.Data)
}

func (e *Entry) MarshalText() ([]byte, error) {
	return []byte(common3.HexEncode(e.Bytes())), nil
}

func (e *Entry) UnmarshalText(text []byte) error {
	return e.Data.UnmarshalText(text)
}

func (e *Entry) Clone() *Entry {
	data := NewDataFromBytes(e.Data.Bytes())
	return &Entry{Data: *data}
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
	// writable indicates if the Merkle Tree allows to write or only to read
	writable bool
}

// NewMerkleTree generates a new Merkle Tree
func NewMerkleTree(storage db.Storage, maxLevels int) (*MerkleTree, error) {
	mt := MerkleTree{storage: storage, maxLevels: maxLevels, writable: true}
	_, gettedRoot, err := mt.dbGet(rootNodeValue)
	if err != nil {
		tx, err := mt.storage.NewTx()
		if err != nil {
			return nil, err
		}
		nodeRoot := NewNodeEmpty()
		k, _ := nodeRoot.Key(), nodeRoot.Value()
		mt.rootKey = k
		mt.dbInsert(tx, rootNodeValue, DBEntryTypeRoot, mt.rootKey[:])
		if err = tx.Commit(); err != nil {
			tx.Close()
			return nil, err
		}
		return &mt, nil
	}
	mt.rootKey = &Hash{}
	copy(mt.rootKey[:], gettedRoot)
	return &mt, nil
}

func (mt *MerkleTree) Snapshot(rootKey *Hash) (*MerkleTree, error) {
	mt.RLock()
	defer mt.RUnlock()
	_, err := mt.GetNode(rootKey)
	if err != nil {
		return nil, err
	}
	return &MerkleTree{storage: mt.storage, maxLevels: mt.maxLevels, rootKey: rootKey, writable: false}, nil
}

// Storage returns the MT storage
func (mt *MerkleTree) Storage() db.Storage {
	return mt.storage
}

// RootKey returns the MT root node key
func (mt *MerkleTree) RootKey() *Hash {
	mt.RLock()
	defer mt.RUnlock()
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
	nextKey := mt.RootKey()
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
		return mt.addNode(tx, newNodeMiddle)
	} else {
		if pathNewLeaf[lvl] {
			newNodeMiddle = NewNodeMiddle(oldLeaf.Key(), newLeaf.Key())
		} else {
			newNodeMiddle = NewNodeMiddle(newLeaf.Key(), oldLeaf.Key())
		}
		// We can add newLeaf now.  We don't need to add oldLeaf because it's already in the tree.
		_, err := mt.addNode(tx, newLeaf)
		if err != nil {
			return nil, err
		}
		return mt.addNode(tx, newNodeMiddle)
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
		return mt.addNode(tx, newLeaf)
	case NodeTypeLeaf:
		// TODO: delete old node n???  Make this optional???
		hIndex := n.Entry.HIndex()
		// Check if leaf node found contains the leaf node we are trying to add
		if bytes.Equal(hIndex[:], newLeaf.Entry.HIndex()[:]) {
			return nil, ErrEntryIndexAlreadyExists
		}
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
		return mt.addNode(tx, newNodeMiddle)
	default:
		return nil, ErrInvalidNodeFound
	}
}

// Add adds the Entry to the MerkleTree
func (mt *MerkleTree) Add(e *Entry) error {
	// verify that the MerkleTree is writable
	if !mt.writable {
		return ErrNotWritable
	}
	// verfy that the ElemBytes are valid and fit inside the mimc7 field.
	bigints := ElemBytesToBigInts(e.Data[:]...)
	ok := cryptoUtils.CheckBigIntArrayInField(bigints, cryptoConstants.Q)
	if !ok {
		return errors.New("Elements not inside the Finite Field over R")
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
	mt.dbInsert(tx, rootNodeValue, DBEntryTypeRoot, mt.rootKey[:])
	return nil
}

// walk is a helper recursive function to iterate over all tree branches
func (mt *MerkleTree) walk(key *Hash, f func(*Node)) error {
	n, err := mt.GetNode(key)
	if err != nil {
		return err
	}
	switch n.Type {
	case NodeTypeEmpty:
		f(n)
	case NodeTypeLeaf:
		f(n)
	case NodeTypeMiddle:
		f(n)
		mt.walk(n.ChildL, f)
		mt.walk(n.ChildR, f)
	default:
		return ErrInvalidNodeFound
	}
	return nil
}

// Walk iterates over all the branches of a MerkleTree with the given rootKey
// if rootKey is nil, it will get the current RootKey of the current state of the MerkleTree.
// For each node, it calls the f function given in the parameters.
// See some examples of the Walk function usage in the merkletree_test.go
// test functions: TestMTWalk, TestMTWalkGraphViz, TestMTWalkDumpClaims
func (mt *MerkleTree) Walk(rootKey *Hash, f func(*Node)) error {
	if rootKey == nil {
		rootKey = mt.RootKey()
	}
	err := mt.walk(rootKey, f)
	return err
}

// GraphViz uses Walk function to generate a string GraphViz representation of the
// tree and writes it to w
func (mt *MerkleTree) GraphViz(w io.Writer, rootKey *Hash) error {
	fmt.Fprintf(w, `digraph hierarchy {
node [fontname=Monospace,fontsize=10,shape=box]
`)
	cnt := 0
	err := mt.Walk(rootKey, func(n *Node) {
		switch n.Type {
		case NodeTypeEmpty:
		case NodeTypeLeaf:
			fmt.Fprintf(w, "\"%v\" [style=filled];\n", n.Key())
		case NodeTypeMiddle:
			lr := [2]string{n.ChildL.String(), n.ChildR.String()}
			for i, _ := range lr {
				if lr[i] == "00000000" {
					lr[i] = fmt.Sprintf("empty%v", cnt)
					fmt.Fprintf(w, "\"%v\" [style=dashed,label=0];\n", lr[i])
					cnt++
				}
			}
			fmt.Fprintf(w, "\"%v\" -> {\"%v\" \"%v\"}\n", n.Key(), lr[0], lr[1])
		default:
		}
	})
	fmt.Fprintf(w, "}\n")
	return err
}

// DumpClaims outputs a list of all the claims in hex.
func (mt *MerkleTree) DumpClaims(rootKey *Hash) ([]string, error) {
	var dumpedClaims []string
	err := mt.Walk(rootKey, func(n *Node) {
		if n.Type == NodeTypeLeaf {
			dumpedClaims = append(dumpedClaims, common3.HexEncode(n.Entry.Bytes()))
		}
	})
	return dumpedClaims, err
}

// ImportClaims parses and adds the dumped list of claims in hex from the
// DumpClaims function.
func (mt *MerkleTree) ImportDumpedClaims(dumpedClaims []string) error {
	for _, c := range dumpedClaims {
		if strings.HasPrefix(c, "0x") {
			c = c[2:]
		}
		if len(c) != 256 {
			return errors.New("hex length different than 256")
		}
		var err error
		var e Entry
		e, err = NewEntryFromHexs(c[:64], c[64:128], c[128:192], c[192:])
		if err != nil {
			return err
		}

		err = mt.Add(&e)
		if err != nil {
			return err
		}
	}
	return nil
}

// DumpClaimsIoWriter uses Walk function to get all the Claims of the tree and write
// them to w.  The output is JSON encoded with claims in hex.
func (mt *MerkleTree) DumpClaimsIoWriter(w io.Writer, rootKey *Hash) error {
	fmt.Fprintf(w, "[\n")
	err := mt.Walk(rootKey, func(n *Node) {
		if n.Type == NodeTypeLeaf {
			fmt.Fprintf(w, "	\"%v\",\n", common3.HexEncode(n.Entry.Bytes()))
		}
	})
	fmt.Fprintf(w, "]\n")
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
	Existence bool
	// depth indicates how deep in the tree the proof goes.
	depth uint
	// notempties is a bitmap of non-empty Siblings found in Siblings.
	notempties [ElemBytesLen - proofFlagsLen]byte
	// Siblings is a list of non-empty sibling keys.
	Siblings []*Hash
	nodeAux  *nodeAux
}

// NewProofFromBytes parses a byte array into a Proof.
func NewProofFromBytes(bs []byte) (*Proof, error) {
	if len(bs) < ElemBytesLen {
		return nil, ErrInvalidProofBytes
	}
	p := &Proof{}
	if (bs[0] & 0x01) == 0 {
		p.Existence = true
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
			p.Siblings = append(p.Siblings, &sib)
			sibIdx++
		}
	}

	if !p.Existence && ((bs[0] & 0x02) != 0) {
		p.nodeAux = &nodeAux{hIndex: &Hash{}, hValue: &Hash{}}
		nodeAuxBytes := siblingBytes[len(p.Siblings)*ElemBytesLen:]
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
	bsLen := proofFlagsLen + len(p.notempties) + ElemBytesLen*len(p.Siblings)
	if p.nodeAux != nil {
		bsLen += 2 * ElemBytesLen
	}
	bs := make([]byte, bsLen)

	if !p.Existence {
		bs[0] |= 0x01
	}
	bs[1] = byte(p.depth)
	copy(bs[proofFlagsLen:len(p.notempties)+proofFlagsLen], p.notempties[:])
	siblingsBytes := bs[len(p.notempties)+proofFlagsLen:]
	for i, k := range p.Siblings {
		copy(siblingsBytes[i*ElemBytesLen:(i+1)*ElemBytesLen], k[:])
	}
	if p.nodeAux != nil {
		bs[0] |= 0x02
		copy(bs[len(bs)-2*ElemBytesLen:], p.nodeAux.hIndex[:])
		copy(bs[len(bs)-1*ElemBytesLen:], p.nodeAux.hValue[:])
	}
	return bs
}

func (p *Proof) MarshalJSON() ([]byte, error) {
	return json.Marshal(common3.HexEncode(p.Bytes()))
}

func (p *Proof) UnmarshalJSON(bs []byte) error {
	proofBytes, err := common3.UnmarshalJSONHexDecode(bs)
	if err != nil {
		return err
	}
	proof, err := NewProofFromBytes(proofBytes)
	if err != nil {
		return err
	}
	*p = *proof
	return nil
}

// String outputs a multiline string representation of the Proof.
func (p *Proof) String() string {
	buf := bytes.NewBufferString("Proof:\n")
	fmt.Fprintf(buf, "\texistence: %v\n", p.Existence)
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
			fmt.Fprintf(buf, "%v ", p.Siblings[sibIdx])
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
// If the rootKey is nil, the current merkletree root is used
func (mt *MerkleTree) GenerateProof(hIndex *Hash, rootKey *Hash) (*Proof, error) {
	p := &Proof{}
	var siblingKey *Hash

	path := getPath(mt.maxLevels, hIndex)
	if rootKey == nil {
		rootKey = mt.RootKey()
	}
	nextKey := rootKey
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
				p.Existence = true
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
			p.Siblings = append(p.Siblings, siblingKey)
		}
	}
	return nil, ErrEntryIndexNotFound
}

// VerifyProof verifies the Merkle Proof for the entry and root.
func VerifyProof(rootKey *Hash, proof *Proof, hIndex, hValue *Hash) bool {
	rootFromProof, err := RootFromProof(proof, hIndex, hValue)
	if err != nil {
		return false
	}
	return bytes.Equal(rootKey[:], rootFromProof[:])
}

func RootFromProof(proof *Proof, hIndex, hValue *Hash) (*Hash, error) {
	sibIdx := len(proof.Siblings) - 1
	var midKey *Hash
	if proof.Existence {
		midKey = LeafKey(hIndex, hValue)
	} else {
		if proof.nodeAux == nil {
			midKey = &HashZero
		} else {
			if bytes.Equal(hIndex[:], proof.nodeAux.hIndex[:]) {
				return nil, fmt.Errorf("Non-existence proof being checked against hIndex equal to nodeAux")
			}
			midKey = LeafKey(proof.nodeAux.hIndex, proof.nodeAux.hValue)
		}
	}
	path := getPath(int(proof.depth), hIndex)
	var siblingKey *Hash
	for lvl := int(proof.depth) - 1; lvl >= 0; lvl-- {
		if testBitBigEndian(proof.notempties[:], uint(lvl)) {
			siblingKey = proof.Siblings[sibIdx]
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
	return midKey, nil
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

// addNode adds a node into the MT.  Empty nodes are not stored in the tree;
// they are all the same and assumed to always exist.
func (mt *MerkleTree) addNode(tx db.Tx, n *Node) (*Hash, error) {
	// verify that the MerkleTree is writable
	if !mt.writable {
		return nil, ErrNotWritable
	}
	if n.Type == NodeTypeEmpty {
		return n.Key(), nil
	}
	k, v := n.Key(), n.Value()
	// Check that the node key doesn't already exist
	if _, err := tx.Get(k[:]); err == nil {
		return nil, ErrNodeKeyAlreadyExists
	}
	tx.Put(k[:], v)
	return k, nil
}

func (mt *MerkleTree) dbGet(k []byte) (NodeType, []byte, error) {
	if bytes.Equal(k, HashZero[:]) {
		return 0, nil, nil
	}

	value, err := mt.storage.Get(k)
	if err != nil {
		return 0, nil, err
	}

	if len(value) < 2 {
		return 0, nil, ErrInvalidDBValue
	}
	nodeType := value[0]
	nodeBytes := value[1:]

	return NodeType(nodeType), nodeBytes, nil
}

func (mt *MerkleTree) dbInsert(tx db.Tx, k []byte, t NodeType, data []byte) {
	v := append([]byte{byte(t)}, data...)
	tx.Put(k, v)
}
