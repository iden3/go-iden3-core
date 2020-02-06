package idenpub

import (
	"bytes"
	"math"
	"sync"

	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
)

var (
	cacheIKey    = []byte("cacheI")
	idenStateKey = []byte("idenState")
	cltRootKey   = []byte("cltRoot")
	rotRootKey   = []byte("rotRoot")
	retRootKey   = []byte("retRoot")
	rotKey       = []byte("rot")
	retKey       = []byte("ret")
)

// IdenPub is a interface, that for the moment will be satisfied at least by IdenPubHTTP & IdenPubIPFS.
type IdenPub interface {
	Publish()
}

// IdenPubHTTP satisfies the IdenPub interface, and stores in a leveldb the published RoT & ReT to be returned when requested.
type IdenPubHTTP struct {
	rw  sync.RWMutex
	db  db.Storage
	rot *merkletree.MerkleTree
	ret *merkletree.MerkleTree
}

// AddLeafRoT adds a new leaf to the given MerkleTree, which contains the Root
func AddLeafRoT(mt *merkletree.MerkleTree, root merkletree.Hash) error {
	l := NewLeafRoT(root)
	return mt.AddEntry(l.Entry())
}

// AddLeafReT adds a new leaf to the given MerkleTree, which contains the Nonce & Version
func AddLeafReT(mt *merkletree.MerkleTree, nonce, version uint32) error {
	l := NewLeafReT(nonce, version)
	return mt.AddEntry(l.Entry())
}

// NewIdenPubHTTP returns a new IdenPubHTTP
func NewIdenPubHTTP(db db.Storage, rot *merkletree.MerkleTree, ret *merkletree.MerkleTree) *IdenPubHTTP {
	return &IdenPubHTTP{
		db:  db,
		rot: rot,
		ret: ret,
	}
}

// Publish publishes the RoT and ReT to the configured way of publishing
func (i *IdenPubHTTP) Publish(idenState, claimsRoot, rootsRoot, revocationsRoot *merkletree.Hash) error {
	// RoT
	w := bytes.NewBufferString("")
	err := i.rot.DumpTree(w, rootsRoot)
	if err != nil {
		return err
	}
	rotBlob := w.Bytes()

	// ReT
	w = bytes.NewBufferString("")
	err = i.ret.DumpTree(w, revocationsRoot)
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

	tx.Put(append(idenStateKey, byte(cacheI)), idenState[:])
	tx.Put(append(cltRootKey, byte(cacheI)), claimsRoot[:])
	tx.Put(append(rotRootKey, byte(cacheI)), rootsRoot[:])
	tx.Put(append(rotKey, byte(cacheI)), rotBlob)
	tx.Put(append(retRootKey, byte(cacheI)), revocationsRoot[:])
	tx.Put(append(retKey, byte(cacheI)), retBlob)

	tx.Put(cacheIKey, []byte{byte(cacheI)})

	return nil
}

func nextCacheI(i int) int {
	return int(math.Mod(float64(i), 1))
}

func (i *IdenPubHTTP) getCacheI(tx db.Tx) (int, error) {
	cacheI, err := tx.Get(cacheIKey)
	if err == db.ErrNotFound {
		cacheI = []byte{1}
	} else if err != nil {
		return 0, err
	}
	return int(cacheI[0]), nil
}

// PublicData contains the RootsTree + Root, and the RevocationTree + Root
type PublicData struct {
	IdenState merkletree.Hash
	ClTRoot   merkletree.Hash
	RoTRoot   merkletree.Hash
	RoT       []byte
	ReTRoot   merkletree.Hash
	ReT       []byte
}

// GetPublicData returns the public data of the IdenPubHTTP.
func (i *IdenPubHTTP) GetPublicData() (*PublicData, error) {
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
	idenState, err := tx.Get(append(idenStateKey, byte(cacheI)))
	if err != nil {
		return nil, err
	}

	// clt
	cltRoot, err := tx.Get(append(cltRootKey, byte(cacheI)))
	if err != nil {
		return nil, err
	}

	// rot
	rotRoot, err := tx.Get(append(rotRootKey, byte(cacheI)))
	if err != nil {
		return nil, err
	}
	rot, err := tx.Get(append(rotKey, byte(cacheI)))
	if err != nil {
		return nil, err
	}

	// ret
	retRoot, err := tx.Get(append(retRootKey, byte(cacheI)))
	if err != nil {
		return nil, err
	}
	ret, err := tx.Get(append(retKey, byte(cacheI)))
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
		IdenState: merkletree.Hash(merkletree.ElemBytes(idenState32)),
		ClTRoot:   merkletree.Hash(merkletree.ElemBytes(cltRoot32)),
		RoTRoot:   merkletree.Hash(merkletree.ElemBytes(rotRoot32)),
		RoT:       rot,
		ReTRoot:   merkletree.Hash(merkletree.ElemBytes(retRoot32)),
		ReT:       ret,
	}
	return p, nil
}
