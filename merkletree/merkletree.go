package merkletree

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
	"sync"

	"github.com/iden3/go-iden3-core/common"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/db"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
)

const (
	// ElemBytesLen is the length in bytes of each element used for storing
	// data and hashing.
	ElemBytesLen = 32
	// IndexLen indicates how many elements are used for the index.
	IndexLen = 4
	// DataLen indicates how many elements are used for the data.
	DataLen = 8
)

// ElemBytes is the basic type used to store data in the MT.  ElemBytes
// corresponds to the serialization of an element from mimc7.
type ElemBytes [ElemBytesLen]byte

func NewElemBytesFromBigInt(v *big.Int) (e ElemBytes) {
	bs := common.SwapEndianness(v.Bytes())
	copy(e[:], bs)
	return e
}

func (e *ElemBytes) BigInt() *big.Int {
	return new(big.Int).SetBytes(common3.SwapEndianness(e[:]))
}

// String returns the last 4 bytes of ElemBytes in hex.
func (e *ElemBytes) String() string {
	return fmt.Sprintf("%v...", hex.EncodeToString(e[ElemBytesLen-4:]))
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
// It consists of 8 elements: e0, e1, e2, e3, ...;
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

func (d Data) MarshalText() ([]byte, error) {
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
	// ErrEntryDataNotMatch is used when the entry data doesn't match the expected one.
	ErrEntryDataNotMatch = errors.New("Entry data doesn't match the expected one")

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

type Entrier interface {
	Entry() *Entry
}

func (e *Entry) Index() []ElemBytes {
	return e.Data[:IndexLen]
}

func (e *Entry) Value() []ElemBytes {
	return e.Data[IndexLen:]
}

// HIndex calculates the hash of the Index of the Entry, used to find the path
// from the root to the leaf in the MT.
func (e *Entry) HIndex() (*Hash, error) {
	var err error
	if e.hIndex == nil { // Cache the hIndex.
		//e.hIndex = HashElems(e.Index()[:]...)
		e.hIndex, err = HashElems(e.Index()...)
	}
	return e.hIndex, err
}

// HValue calculates the hash of the Value of the Entry
func (e *Entry) HValue() (*Hash, error) {
	var err error
	if e.hValue == nil { // Cache the hValue.
		e.hValue, err = HashElems(e.Value()...)
	}
	return e.hValue, err
}

// HiHv returns the HIndex and HValue of the Entry
func (e *Entry) HiHv() (*Hash, *Hash, error) {
	hi, err := e.HIndex()
	if err != nil {
		return nil, nil, err
	}
	hv, err := e.HValue()
	if err != nil {
		return nil, nil, err
	}

	return hi, hv, nil
}

func (e *Entry) Bytes() []byte {
	b := e.Data.Bytes()
	return b[:]
}

func (e1 *Entry) Equal(e2 *Entry) bool {
	return e1.Data.Equal(&e2.Data)
}

func (e Entry) MarshalText() ([]byte, error) {
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
		k, err := nodeRoot.Key()
		if err != nil {
			return nil, err
		}
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
			hi, err := n.Entry.HIndex()
			if err != nil {
				return nil, err
			}
			if bytes.Equal(hIndex[:], hi[:]) {
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

// EntryExists checks if a given entry is in the merkle tree starting from the
// rootKey.  If rootKey is nil, the current merkle tree root is used.
func (mt *MerkleTree) EntryExists(entry *Entry, rootKey *Hash) error {
	var err error
	if rootKey != nil {
		mt, err = mt.Snapshot(rootKey)
		if err != nil {
			return err
		}
	}
	hi, err := entry.HIndex()
	if err != nil {
		return err
	}
	data, err := mt.GetDataByIndex(hi)
	if err != nil {
		return err
	}
	foundEntry := &Entry{Data: *data}
	if !foundEntry.Equal(entry) {
		return ErrEntryDataNotMatch
	}
	return nil
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
		oldLeafKey, err := oldLeaf.Key()
		if err != nil {
			return nil, err
		}
		newLeafKey, err := newLeaf.Key()
		if err != nil {
			return nil, err
		}

		if pathNewLeaf[lvl] {
			newNodeMiddle = NewNodeMiddle(oldLeafKey, newLeafKey)
		} else {
			newNodeMiddle = NewNodeMiddle(newLeafKey, oldLeafKey)
		}
		// We can add newLeaf now.  We don't need to add oldLeaf because it's already in the tree.
		_, err = mt.addNode(tx, newLeaf)
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
		hIndex, err := n.Entry.HIndex()
		if err != nil {
			return nil, err
		}
		// Check if leaf node found contains the leaf node we are trying to add
		newLeafHi, err := newLeaf.Entry.HIndex()
		if err != nil {
			return nil, err
		}
		if bytes.Equal(hIndex[:], newLeafHi[:]) {
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

func CheckEntryInField(e Entry) bool {
	bigints := ElemBytesToBigInts(e.Data[:]...)
	ok := cryptoUtils.CheckBigIntArrayInField(bigints)
	return ok
}

// AddClaim adds the Claim that fullfills the Entrier interface to the MerkleTree
func (mt *MerkleTree) AddClaim(e Entrier) error {
	return mt.AddEntry(e.Entry())
}

// AddEntry adds the Entry to the MerkleTree
func (mt *MerkleTree) AddEntry(e *Entry) error {
	// verify that the MerkleTree is writable
	if !mt.writable {
		return ErrNotWritable
	}
	// verfy that the ElemBytes are valid and fit inside the mimc7 field.
	if !CheckEntryInField(*e) {
		return errors.New("Elements not inside the Finite Field over R")
	}
	tx, err := mt.storage.NewTx()
	if err != nil {
		return err
	}
	mt.Lock()
	defer func() {
		if err == nil {
			if err := tx.Commit(); err != nil {
				tx.Close()
			}
		} else {
			tx.Close()
		}
		mt.Unlock()
	}()

	newNodeLeaf := NewNodeLeaf(e)
	hIndex, err := e.HIndex()
	if err != nil {
		return err
	}
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
		if err := mt.walk(n.ChildL, f); err != nil {
			return err
		}
		if err := mt.walk(n.ChildR, f); err != nil {
			return err
		}
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
	var errIn error
	err := mt.Walk(rootKey, func(n *Node) {
		k, err := n.Key()
		if err != nil {
			errIn = err
		}
		switch n.Type {
		case NodeTypeEmpty:
		case NodeTypeLeaf:
			fmt.Fprintf(w, "\"%v\" [style=filled];\n", k)
		case NodeTypeMiddle:
			lr := [2]string{n.ChildL.String(), n.ChildR.String()}
			for i := range lr {
				if lr[i] == "00000000..." {
					lr[i] = fmt.Sprintf("empty%v", cnt)
					fmt.Fprintf(w, "\"%v\" [style=dashed,label=0];\n", lr[i])
					cnt++
				}
			}
			fmt.Fprintf(w, "\"%v\" -> {\"%v\" \"%v\"}\n", k, lr[0], lr[1])
		default:
		}
	})
	fmt.Fprintf(w, "}\n")
	if errIn != nil {
		return errIn
	}
	return err
}

// PrintGraphViz prints directly the GraphViz() output
func (mt *MerkleTree) PrintGraphViz(rootKey *Hash) error {
	if rootKey == nil {
		rootKey = mt.RootKey()
	}
	w := bytes.NewBufferString("")
	fmt.Fprintf(w, "--------\nGraphViz of the MerkleTree with RootKey "+rootKey.Hex()+"\n")
	err := mt.GraphViz(w, nil)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "End of GraphViz of the MerkleTree with RootKey "+rootKey.Hex()+"\n--------\n")

	fmt.Println(w)
	return nil
}

// DumpTree outputs a list of all the key value in hex. Notice that this will
// output the full tree, which is not needed to reconstruct the Tree. To
// reconstruct the tree can be done from the output of DumpClaims funtion.  The
// difference between DumpTree and DumpClaims is that with DumpTree the size of
// the output will be almost the double (in raw bytes, but with the current
// implementation DumpTree output is even smaller than DumpClaims output size, but because DumpTree stores in
// binary while DumpClaims stores in hex) but to recover the tree will not need
// to compute the Tree, while with DumpClaims will require to compute the Tree
// (with the computational cost of each hash)
func (mt *MerkleTree) DumpTree(w io.Writer, rootKey *Hash) error {
	var errS error
	err := mt.Walk(rootKey, func(n *Node) {
		if n.Type != NodeTypeEmpty {
			k, err := n.Key()
			if err != nil {
				errS = err
			}
			err = serializeKV(w, k.Bytes(), n.Value())
			if err != nil {
				errS = err
			}
		}
	})
	if err != nil {
		return err
	}
	if errS != nil {
		return errS
	}

	if rootKey == nil {
		rootKey = mt.RootKey()
	}
	err = serializeKV(w, rootNodeValue, rootKey.Bytes())

	return err
}

func checkKVLen(kLen, vLen int) error {
	if kLen > 0xff {
		return fmt.Errorf("len(k) %d > 0xff", kLen)
	}
	if vLen > 0xffff {
		return fmt.Errorf("len(v) %d > 0xffff", vLen)
	}
	return nil
}

func serializeKV(w io.Writer, k, v []byte) error {
	if err := checkKVLen(len(k), len(v)); err != nil {
		return err
	}
	kH := byte(len(k))
	vH := common3.Uint16ToBytes(uint16(len(v)))
	_, err := w.Write([]byte{kH})
	if err != nil {
		return err
	}
	_, err = w.Write(vH)
	if err != nil {
		return err
	}
	_, err = w.Write(k)
	if err != nil {
		return err
	}
	_, err = w.Write(v)
	if err != nil {
		return err
	}
	return nil
}

func deserializeKV(r io.Reader) ([]byte, []byte, error) {
	header := make([]byte, 3)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, nil, err
	}
	kLen := int(header[0])
	vLen := int(common3.BytesToUint16(header[1:]))
	kv := make([]byte, kLen+vLen)
	_, err = io.ReadFull(r, kv)
	if err == io.EOF {
		return nil, nil, io.ErrUnexpectedEOF
	} else if err != nil {
		return nil, nil, err
	}
	return kv[:kLen], kv[kLen:], nil
}

// ImportTree imports the tree from the output from the DumpTree function
func (mt *MerkleTree) ImportTree(i io.Reader) error {
	tx, err := mt.storage.NewTx()
	if err != nil {
		return err
	}
	mt.Lock()
	defer func() {
		if err == nil {
			if err := tx.Commit(); err != nil {
				tx.Close()
			}
		} else {
			tx.Close()
		}
		mt.Unlock()
	}()

	r := bufio.NewReader(i)
	for {
		k, v, err := deserializeKV(r)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		tx.Put(k, v)
	}

	v, err := tx.Get(rootNodeValue)
	if err != nil {
		return err
	}
	mt.rootKey = &Hash{}
	copy(mt.rootKey[:], v)

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
		c = strings.TrimPrefix(c, "0x")
		if len(c) != 2*ElemBytesLen*DataLen { // 2*ElemBytesLen because is in Hexadecimal string, so each byte is represented by 2 char
			return fmt.Errorf("hex length different than %d", 2*ElemBytesLen*DataLen)
		}
		var err error
		var e Entry
		var d Data
		var dataBytes [ElemBytesLen * DataLen]byte
		err = common3.HexDecodeInto(dataBytes[:], []byte(c))
		if err != nil {
			return err
		}
		d = *NewDataFromBytes(dataBytes)
		e.Data = d

		err = mt.AddEntry(&e)
		if err != nil {
			return err
		}
	}
	return nil
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
		if common.TestBitBigEndian(p.notempties[:], i) {
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

func (p Proof) MarshalJSON() ([]byte, error) {
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
func (p Proof) String() string {
	buf := bytes.NewBufferString("{")
	fmt.Fprintf(buf, "Existence: %v, ", p.Existence)
	fmt.Fprintf(buf, "Depth: %v, ", p.depth)
	fmt.Fprintf(buf, "NotEmpties: ")
	for i := uint(0); i < p.depth; i++ {
		if common.TestBitBigEndian(p.notempties[:], i) {
			fmt.Fprintf(buf, "1")
		} else {
			fmt.Fprintf(buf, "0")
		}
	}
	fmt.Fprintf(buf, ", ")
	fmt.Fprintf(buf, "Siblings: [")
	sibIdx := 0
	for i := uint(0); i < p.depth; i++ {
		if common.TestBitBigEndian(p.notempties[:], i) {
			fmt.Fprintf(buf, "%v", p.Siblings[sibIdx])
			sibIdx++
		} else {
			fmt.Fprintf(buf, "0")
		}
		if i < p.depth-1 {
			fmt.Fprintf(buf, ", ")
		}
	}
	fmt.Fprintf(buf, "]")
	if p.nodeAux != nil {
		fmt.Fprintf(buf, ", NodeAux: {Hi: %v, Hv: %v}}", p.nodeAux.hIndex, p.nodeAux.hValue)
	}
	fmt.Fprintf(buf, "}")
	return buf.String()
}

// SiblingsFromProof returns all the siblings of the proof. This function is used to generate the siblings input for the circom circuits.
func SiblingsFromProof(proof *Proof) []*Hash {
	sibIdx := 0
	var siblings []*Hash
	for lvl := 0; lvl < int(proof.depth); lvl++ {
		if common.TestBitBigEndian(proof.notempties[:], uint(lvl)) {
			siblings = append(siblings, proof.Siblings[sibIdx])
			sibIdx++
		} else {
			siblings = append(siblings, &HashZero)
		}
	}
	return siblings
}

func (p *Proof) AllSiblings() []*Hash {
	return SiblingsFromProof(p)
}

func (p *Proof) AllSiblingsCircom(levels int) []*big.Int {
	siblings := p.AllSiblings()
	// Add the rest of empty levels to the siblings
	for i := len(siblings); i < levels; i++ {
		siblings = append(siblings, &HashZero)
	}
	siblings = append(siblings, &HashZero) // add extra level for circom compatibility
	siblingsBigInt := make([]*big.Int, len(siblings))
	for i, sibling := range siblings {
		siblingsBigInt[i] = sibling.BigInt()
	}
	return siblingsBigInt
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
			nHi, nHv, err := n.Entry.HiHv()
			if err != nil {
				return nil, err
			}
			if bytes.Equal(hIndex[:], nHi[:]) {
				p.Existence = true
				return p, nil
			} else {
				// We found a leaf whose entry didn't match hIndex
				p.nodeAux = &nodeAux{hIndex: nHi, hValue: nHv}
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
			common.SetBitBigEndian(p.notempties[:], uint(p.depth))
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

// RootFromProof calculates the root that would correspond to a tree whose
// siblings are the ones in the proof with the claim hashing to hIndex and
// hValue.
func RootFromProof(proof *Proof, hIndex, hValue *Hash) (*Hash, error) {
	sibIdx := len(proof.Siblings) - 1
	var err error
	var midKey *Hash
	if proof.Existence {
		midKey, err = LeafKey(hIndex, hValue)
		if err != nil {
			return nil, err
		}
	} else {
		if proof.nodeAux == nil {
			midKey = &HashZero
		} else {
			if bytes.Equal(hIndex[:], proof.nodeAux.hIndex[:]) {
				return nil, fmt.Errorf("Non-existence proof being checked against hIndex equal to nodeAux")
			}
			midKey, err = LeafKey(proof.nodeAux.hIndex, proof.nodeAux.hValue)
			if err != nil {
				return nil, err
			}
		}
	}
	path := getPath(int(proof.depth), hIndex)
	var siblingKey *Hash
	for lvl := int(proof.depth) - 1; lvl >= 0; lvl-- {
		if common.TestBitBigEndian(proof.notempties[:], uint(lvl)) {
			siblingKey = proof.Siblings[sibIdx]
			sibIdx--
		} else {
			siblingKey = &HashZero
		}
		if path[lvl] {
			midKey, err = NewNodeMiddle(siblingKey, midKey).Key()
			if err != nil {
				return nil, err
			}
		} else {
			midKey, err = NewNodeMiddle(midKey, siblingKey).Key()
			if err != nil {
				return nil, err
			}
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
		return n.Key()
	}
	k, err := n.Key()
	if err != nil {
		return nil, err
	}
	v := n.Value()
	// Check that the node key doesn't already exist
	if _, err := tx.Get(k[:]); err == nil {
		return nil, ErrNodeKeyAlreadyExists
	}
	tx.Put(k[:], v)
	return k, nil
}

// dbGet is a helper function to get the node of a key from the internal
// storage.
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

// dbInsert is a helper function to insert a node into a key in an open db
// transaction.
func (mt *MerkleTree) dbInsert(tx db.Tx, k []byte, t NodeType, data []byte) {
	v := append([]byte{byte(t)}, data...)
	tx.Put(k, v)
}

// HexStringToHash decodes a hex string into a Hash.
func HexStringToHash(s string) Hash {
	b, err := common3.HexDecode(s)
	if err != nil {
		panic(err)
	}
	var b32 [ElemBytesLen]byte
	copy(b32[:], b[:32])
	return Hash(ElemBytes(b32))
}
