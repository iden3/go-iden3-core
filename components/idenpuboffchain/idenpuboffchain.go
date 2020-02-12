package idenpuboffchain

import (
	"bytes"
	"math"
	"sync"

	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
)

var (
	dbKeyCacheI          = []byte("cacheI")
	dbKeyIdenState       = []byte("idenState")
	dbKeyClaimsRoot      = []byte("claimsroot")
	dbKeyRootsRoot       = []byte("rootsroot")
	dbKeyRevocationsRoot = []byte("revocationsroot")
	dbKeyRootsTree       = []byte("rootstree")
	dbKeyRevocationsTree = []byte("revocationstree")
)

// IdenPubOffChainer is a interface, that for the moment will be satisfied at least by IdenPubOffChainHttp & IdenPubIPFS.
type IdenPubOffChainer interface {
	Publish()
}

// IdenPubOffChainHttp satisfies the IdenPubOffChainer interface, and stores in a leveldb the published RootsTree & RevocationsTree to be returned when requested.
type IdenPubOffChainHttp struct {
	rw              sync.RWMutex
	db              db.Storage
	rootsTree       *merkletree.MerkleTree
	revocationsTree *merkletree.MerkleTree
}

// NewIdenPubOffChainHttp returns a new IdenPubOffChainHttp
func NewIdenPubOffChainHttp(db db.Storage, rootsTree *merkletree.MerkleTree, revocationsTree *merkletree.MerkleTree) *IdenPubOffChainHttp {
	return &IdenPubOffChainHttp{
		db:              db,
		rootsTree:       rootsTree,
		revocationsTree: revocationsTree,
	}
}

// Publish publishes the RootsTree and RevocationsTree to the configured way of publishing
func (i *IdenPubOffChainHttp) Publish(idenState, claimsRoot, rootsRoot, revocationsRoot *merkletree.Hash) error {
	// RootsTree
	w := bytes.NewBufferString("")
	err := i.rootsTree.DumpTree(w, rootsRoot)
	if err != nil {
		return err
	}
	rotBlob := w.Bytes()

	// RevocationsTree
	w = bytes.NewBufferString("")
	err = i.revocationsTree.DumpTree(w, revocationsRoot)
	if err != nil {
		return err
	}
	retBlob := w.Bytes()

	tx, err := i.db.NewTx()
	if err != nil {
		return err
	}
	i.rw.Lock()
	defer func() {
		if err == nil {
			if err := tx.Commit(); err != nil {
				tx.Close()
			}
		} else {
			tx.Close()
		}
		i.rw.Unlock()
	}()

	cacheI, err := i.getCacheI(tx)
	if err != nil {
		return err
	}
	cacheI = nextCacheI(cacheI)

	tx.Put(append(dbKeyIdenState, byte(cacheI)), idenState[:])
	tx.Put(append(dbKeyClaimsRoot, byte(cacheI)), claimsRoot[:])
	tx.Put(append(dbKeyRootsRoot, byte(cacheI)), rootsRoot[:])
	tx.Put(append(dbKeyRootsTree, byte(cacheI)), rotBlob)
	tx.Put(append(dbKeyRevocationsRoot, byte(cacheI)), revocationsRoot[:])
	tx.Put(append(dbKeyRevocationsTree, byte(cacheI)), retBlob)

	tx.Put(dbKeyCacheI, []byte{byte(cacheI)})

	return nil
}

func nextCacheI(i int) int {
	return int(math.Mod(float64(i), 1))
}

func (i *IdenPubOffChainHttp) getCacheI(tx db.Tx) (int, error) {
	cacheI, err := tx.Get(dbKeyCacheI)
	if err == db.ErrNotFound {
		cacheI = []byte{1}
	} else if err != nil {
		return 0, err
	}
	return int(cacheI[0]), nil
}

// PublicData contains the RootsTree + Root, and the RevocationTree + Root
type PublicData struct {
	IdenState           merkletree.Hash
	ClaimsTreeRoot      merkletree.Hash
	RootsTreeRoot       merkletree.Hash
	RootsTree           []byte
	RevocationsTreeRoot merkletree.Hash
	RevocationsTree     []byte
}

// GetPublicData returns the public data of the IdenPubOffChainHttp.
func (i *IdenPubOffChainHttp) GetPublicData() (*PublicData, error) {
	tx, err := i.db.NewTx()
	if err != nil {
		return nil, err
	}
	i.rw.RLock()
	defer func() {
		tx.Close()
		i.rw.RUnlock()
	}()

	cacheI, err := i.getCacheI(tx)
	if err != nil {
		return nil, err
	}

	// idenState
	idenState, err := tx.Get(append(dbKeyIdenState, byte(cacheI)))
	if err != nil {
		return nil, err
	}

	// clt
	cltRoot, err := tx.Get(append(dbKeyClaimsRoot, byte(cacheI)))
	if err != nil {
		return nil, err
	}

	// rot
	rotRoot, err := tx.Get(append(dbKeyRootsRoot, byte(cacheI)))
	if err != nil {
		return nil, err
	}
	rot, err := tx.Get(append(dbKeyRootsTree, byte(cacheI)))
	if err != nil {
		return nil, err
	}

	// ret
	retRoot, err := tx.Get(append(dbKeyRevocationsRoot, byte(cacheI)))
	if err != nil {
		return nil, err
	}
	ret, err := tx.Get(append(dbKeyRevocationsTree, byte(cacheI)))
	if err != nil {
		return nil, err
	}

	var idenState32 [merkletree.ElemBytesLen]byte
	var cltRoot32 [merkletree.ElemBytesLen]byte
	var rotRoot32 [merkletree.ElemBytesLen]byte
	var retRoot32 [merkletree.ElemBytesLen]byte
	copy(idenState32[:], idenState[:32])
	copy(cltRoot32[:], cltRoot[:32])
	copy(rotRoot32[:], rotRoot[:32])
	copy(retRoot32[:], retRoot[:32])

	p := &PublicData{
		IdenState:           merkletree.Hash(merkletree.ElemBytes(idenState32)),
		ClaimsTreeRoot:      merkletree.Hash(merkletree.ElemBytes(cltRoot32)),
		RootsTreeRoot:       merkletree.Hash(merkletree.ElemBytes(rotRoot32)),
		RootsTree:           rot,
		RevocationsTreeRoot: merkletree.Hash(merkletree.ElemBytes(retRoot32)),
		RevocationsTree:     ret,
	}
	return p, nil
}
