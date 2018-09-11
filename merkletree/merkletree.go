package merkletree

import (
	"bytes"
	"errors"
	"sync"
)

const (
	// emptyNodeType indicates the type of an EmptyNodeValue Node
	EmptyNodeType = 00
	// normalNodeType indicates the type of a middle Node
	normalNodeType = 01
	// finalNodeType indicates the type of middle Node that is in an optimized branch, then in the value contains the value of the final leaf node of that branch
	finalNodeType = 02
	// valueNodeType indicates the type of a value Node
	valueNodeType = 03
	// rootNodeType indicates the type of a root Node
	rootNodeType = 04
)

// EmptyNodeValue is a [32]byte EmptyNodeValue array, all to zero
var (
	ErrNodeAlreadyExists = errors.New("node already exists")
	rootNodeValue        = HashBytes([]byte("root"))
	EmptyNodeValue       = Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

// Hash used in this tree, is the [32]byte keccak()
type Hash [32]byte

// Value is the interface of a generic claim, a key value object stored in the leveldb
type Value interface {
	IndexLength() uint32
	Bytes() []byte
}

//MerkleTree struct with the main elements of the Merkle Tree
type MerkleTree struct {
	sync.RWMutex
	storage   Storage
	root      Hash
	numLevels int // Height of the Merkle Tree, number of levels
}

// New generates a new Merkle Tree
func New(storage Storage, numLevels int) (*MerkleTree, error) {
	var mt MerkleTree
	mt.storage = storage
	mt.numLevels = numLevels
	var err error
	_, _, rootHash, err := mt.storage.Get(rootNodeValue)
	if err != nil {
		mt.root = EmptyNodeValue
		tx, err := mt.storage.NewTx()
		if err != nil {
			return nil, err
		}
		tx.Insert(rootNodeValue, rootNodeType, 0, mt.root[:])
		if err = tx.Commit(); err != nil {
			return nil, err
		}
		// return &mt
	}
	copy(mt.root[:], rootHash)
	return &mt, nil
}

// Storage returns the merkletree.Storage
func (mt *MerkleTree) Storage() Storage {
	return mt.storage
}

// Root returns the merkletree.Root
func (mt *MerkleTree) Root() Hash {
	return mt.root
}

// NumLevels returns the merkletree.NumLevels
func (mt *MerkleTree) NumLevels() int {
	return mt.numLevels
}

// Add adds the claim to the MT
func (mt *MerkleTree) Add(v Value) error {
	var err error
	var tx StorageTx

	tx, err = mt.storage.NewTx()
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

	hi := HashBytes(v.Bytes()[:v.IndexLength()])
	path := getPath(mt.numLevels, hi)

	nodeHash := mt.root
	var siblings []Hash
	for i := mt.numLevels - 2; i >= 0; i-- {
		nodeType, indexLength, nodeBytes, err := tx.Get(nodeHash)
		if err != nil {
			return err
		}
		if nodeType == byte(finalNodeType) {
			/*
				if node is Final Tree, let's:
					- search until where the path of each leafHash is shared --> posDiff
					- in that layer, make two childs, one for each Final Tree
						- one for the already existing leaf (from the old Final Tree)
						- another for the new claim added
			*/
			hiChild := HashBytes(nodeBytes[:indexLength])
			pathChild := getPath(mt.numLevels, hiChild)
			posDiff := comparePaths(pathChild, path)
			if posDiff == -1 {
				return ErrNodeAlreadyExists
			}
			finalNode1Hash := calcHashFromLeafAndLevel(posDiff, pathChild, HashBytes(nodeBytes))
			tx.Insert(finalNode1Hash, finalNodeType, indexLength, nodeBytes)
			finalNode2Hash := calcHashFromLeafAndLevel(posDiff, path, HashBytes(v.Bytes()))
			tx.Insert(finalNode2Hash, finalNodeType, v.IndexLength(), v.Bytes())
			// now the parent
			var parentNode treeNode
			if path[posDiff] {
				parentNode = treeNode{
					ChildL: finalNode1Hash,
					ChildR: finalNode2Hash,
				}
			} else {
				parentNode = treeNode{
					ChildL: finalNode2Hash,
					ChildR: finalNode1Hash,
				}
			}
			siblings = append(siblings, getEmptiesBetweenIAndPosHash(mt, i, posDiff+1)...)
			if mt.root, err = mt.replaceLeaf(tx, siblings, path[posDiff+1:], parentNode.Ht(), normalNodeType, 0, parentNode.Bytes()); err != nil {
				return err
			}
			tx.Insert(rootNodeValue, rootNodeType, 0, mt.root[:])
			return nil
		}
		node := parseNodeBytes(nodeBytes)
		var sibling Hash
		if !path[i] {
			nodeHash = node.ChildL
			sibling = node.ChildR
		} else {
			nodeHash = node.ChildR
			sibling = node.ChildL
		}
		siblings = append(siblings, sibling)

		if bytes.Equal(nodeHash[:], EmptyNodeValue[:]) {
			// if the node is EmptyNodeValue, the claim will go directly at that height, as a Final Node
			if i == mt.numLevels-2 && bytes.Equal(siblings[len(siblings)-1][:], EmptyNodeValue[:]) {
				// if the pt node is the unique in the tree, just put it into the root node
				// this means to be in i==mt.NumLevels-2 && nodeHash==EmptyNodeValue
				finalNodeHash := calcHashFromLeafAndLevel(i+1, path, HashBytes(v.Bytes()))
				tx.Insert(finalNodeHash, finalNodeType, v.IndexLength(), v.Bytes())
				mt.root = finalNodeHash
				tx.Insert(rootNodeValue, rootNodeType, 0, mt.root[:])
				return nil
			}
			finalNodeHash := calcHashFromLeafAndLevel(i, path, HashBytes(v.Bytes()))
			if mt.root, err = mt.replaceLeaf(tx, siblings, path[i:], finalNodeHash, finalNodeType, v.IndexLength(), v.Bytes()); err != nil {
				return err
			}
			tx.Insert(rootNodeValue, rootNodeType, 0, mt.root[:])
			return nil
		}
	}

	mt.root, err = mt.replaceLeaf(tx, siblings, path, HashBytes(v.Bytes()), valueNodeType, v.IndexLength(), v.Bytes())
	tx.Insert(rootNodeValue, rootNodeType, 0, mt.root[:])
	return nil
}

// GenerateProof generates the Merkle Proof from a given claimHash for the current root
func (mt *MerkleTree) GenerateProof(v Value) ([]byte, error) {
	mt.RLock()
	defer mt.RUnlock()

	var empties [32]byte

	hi := HashBytes(v.Bytes()[:v.IndexLength()])
	path := getPath(mt.numLevels, hi)
	var siblings []Hash
	nodeHash := mt.root

	for level := 0; level < mt.numLevels-1; level++ {
		nodeType, _, nodeBytes, err := mt.storage.Get(nodeHash)
		if err != nil {
			return nil, err
		}
		if nodeType == byte(finalNodeType) {
			break
		}
		node := parseNodeBytes(nodeBytes)

		var sibling Hash
		if !path[mt.numLevels-level-2] {
			nodeHash = node.ChildL
			sibling = node.ChildR
		} else {
			nodeHash = node.ChildR
			sibling = node.ChildL
		}
		if !bytes.Equal(sibling[:], EmptyNodeValue[:]) {
			setbitmap(empties[:], uint(level))
			siblings = append([]Hash{sibling}, siblings...)
		}
	}
	// merge empties and siblings
	var mp []byte
	mp = append(mp, empties[:]...)
	for k := range siblings {
		mp = append(mp, siblings[k][:]...)
	}
	return mp, nil
}

// GetValueInPos returns the merkletree value in the position of the Hash of the Index (Hi)
func (mt *MerkleTree) GetValueInPos(hi Hash) ([]byte, error) {

	mt.RLock()
	defer mt.RUnlock()

	path := getPath(mt.numLevels, hi)
	nodeHash := mt.root
	for i := mt.numLevels - 2; i >= 0; i-- {
		nodeType, indexLength, nodeBytes, err := mt.storage.Get(nodeHash)
		if err != nil {
			return nodeBytes, err
		}
		if nodeType == byte(finalNodeType) {
			// check if nodeBytes path is different of hi
			index := nodeBytes[:indexLength]
			hi := HashBytes(index)
			nodePath := getPath(mt.numLevels, hi)
			posDiff := comparePaths(path, nodePath)
			// if is different, return an EmptyNodeValue, else return the nodeBytes
			if posDiff != -1 {
				return EmptyNodeValue[:], nil
			}
			return nodeBytes, nil
		}
		node := parseNodeBytes(nodeBytes)
		if !path[i] {
			nodeHash = node.ChildL
		} else {
			nodeHash = node.ChildR
		}
	}
	_, _, valueBytes, err := mt.storage.Get(nodeHash)
	if err != nil {
		return valueBytes, err
	}
	return valueBytes, nil
}

func calcHashFromLeafAndLevel(untilLevel int, path []bool, leafHash Hash) Hash {
	nodeCurrLevel := leafHash
	for i := 0; i < untilLevel; i++ {
		if path[i] {
			node := treeNode{
				ChildL: EmptyNodeValue,
				ChildR: nodeCurrLevel,
			}
			nodeCurrLevel = node.Ht()
		} else {
			node := treeNode{
				ChildL: nodeCurrLevel,
				ChildR: EmptyNodeValue,
			}
			nodeCurrLevel = node.Ht()
		}
	}
	return nodeCurrLevel
}

func (mt *MerkleTree) replaceLeaf(tx StorageTx, siblings []Hash, path []bool, newLeafHash Hash, nodetype byte, indexLength uint32, newLeafValue []byte) (Hash, error) {
	// add the new claim
	tx.Insert(newLeafHash, nodetype, indexLength, newLeafValue)
	currNode := newLeafHash
	// here the path is only the path[posDiff+1]
	for i := 0; i < len(siblings); i++ {
		if !path[i] {
			node := treeNode{
				ChildL: currNode,
				ChildR: siblings[len(siblings)-1-i],
			}
			tx.Insert(node.Ht(), normalNodeType, 0, node.Bytes())
			currNode = node.Ht()
		} else {

			node := treeNode{
				ChildL: siblings[len(siblings)-1-i],
				ChildR: currNode,
			}
			tx.Insert(node.Ht(), normalNodeType, 0, node.Bytes())
			currNode = node.Ht()
		}
	}

	return currNode, nil // currNode = root
}

// CheckProof validates the Merkle Proof for the claimHash and root
func CheckProof(root Hash, proof []byte, v Value, numLevels int) bool {

	var empties [32]byte
	copy(empties[:], proof[:len(empties)])
	hashLen := len(EmptyNodeValue)

	var siblings []Hash
	for i := len(empties); i < len(proof); i += hashLen {
		var siblingHash Hash
		copy(siblingHash[:], proof[i:i+hashLen])
		siblings = append(siblings, siblingHash)
	}

	hi := HashBytes(v.Bytes()[:v.IndexLength()])
	path := getPath(numLevels, hi)

	nodeHash := HashBytes(v.Bytes())
	siblingUsedPos := 0

	for level := numLevels - 2; level >= 0; level-- {

		var sibling Hash

		if testbitmap(empties[:], uint(level)) {
			sibling = siblings[siblingUsedPos]
			siblingUsedPos++
		} else {
			sibling = EmptyNodeValue
		}

		// calculate the nodeHash with the current nodeHash and the sibling
		if path[numLevels-level-2] {
			node := treeNode{
				ChildL: sibling,
				ChildR: nodeHash,
			}
			nodeHash = node.Ht()
		} else {
			node := treeNode{
				ChildL: nodeHash,
				ChildR: sibling,
			}
			nodeHash = node.Ht()
		}
	}
	return bytes.Equal(nodeHash[:], root[:])
}
