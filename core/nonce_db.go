package core

import (
	"container/heap"
	"sync"
	"time"
)

// NonceObj represents a nonce with an expiration date and auxiliary data
type NonceObj struct {
	Nonce      string
	Expiration int64
	Aux        interface{}
}

type nonceObjs []*NonceObj

func (ns nonceObjs) Len() int           { return len(ns) }
func (ns nonceObjs) Less(i, j int) bool { return ns[i].Expiration < ns[j].Expiration }
func (ns nonceObjs) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }

func (ns *nonceObjs) Push(x interface{}) {
	*ns = append(*ns, x.(*NonceObj))
}

func (ns *nonceObjs) Pop() interface{} {
	old := *ns
	n := len(old)
	x := old[n-1]
	*ns = old[0 : n-1]
	return x
}

type nonceObjHeap struct {
	elems nonceObjs
}

func newNonceObjHeap() *nonceObjHeap {
	h := &nonceObjs{}
	heap.Init(h)
	return &nonceObjHeap{*h}
}

func (h *nonceObjHeap) Push(nObj *NonceObj) {
	heap.Push(&h.elems, nObj)
}

func (h *nonceObjHeap) Pop() *NonceObj {
	return heap.Pop(&h.elems).(*NonceObj)
}

func (h *nonceObjHeap) Peek() *NonceObj {
	if len(h.elems) == 0 {
		return nil
	} else {
		return h.elems[0]
	}
}

func (h *nonceObjHeap) Len() int {
	return h.elems.Len()
}

// NonceDb is a collection of nonces with expiration dates.
type NonceDb struct {
	mutex              sync.RWMutex
	deleteCounterMutex sync.Mutex
	noncesHeap         *nonceObjHeap
	nonceObjsByNonce   map[string]*NonceObj
	deleteCounter      uint64
}

// NewNonceDb creates an empty NonceDb.
func NewNonceDb() *NonceDb {
	return &NonceDb{
		noncesHeap:       newNonceObjHeap(),
		nonceObjsByNonce: make(map[string]*NonceObj),
	}
}

func (ndb *NonceDb) add(nonce string, expiration int64, aux interface{}) bool {
	if _, ok := ndb.nonceObjsByNonce[nonce]; ok {
		return false
	}
	nObj := &NonceObj{nonce, expiration, aux}
	ndb.nonceObjsByNonce[nonce] = nObj
	ndb.noncesHeap.Push(nObj)
	return true
}

// Add adds a nonce with aux data that expires after delta seconds>
func (ndb *NonceDb) Add(nonce string, delta int64, aux interface{}) bool {
	expiration := time.Now().Unix() + delta
	ndb.mutex.Lock()
	defer ndb.mutex.Unlock()
	return ndb.add(nonce, expiration, aux)
}

// AddAux adds aux data to a nonceObj only if it doesn't have an Aux already.
// Returns true on success.
func (ndb *NonceDb) AddAux(nonce string, aux interface{}) bool {
	ndb.mutex.Lock()
	defer ndb.mutex.Unlock()
	nObj, ok := ndb.nonceObjsByNonce[nonce]
	if !ok {
		return false
	} else if nObj.Aux != nil {
		return false
	}
	nObj.Aux = aux
	return true
}

// Search searches a nonce object by nonce.  Returns false if the nonce is not
// found or expired.
func (ndb *NonceDb) Search(nonce string) (*NonceObj, bool) {
	ndb.DeleteOldOportunistic()
	ndb.mutex.RLock()
	nObj, ok := ndb.nonceObjsByNonce[nonce]
	ndb.mutex.RUnlock()
	if !ok {
		return nil, false
	}
	if nObj.Expiration < time.Now().Unix() {
		ok = false
	}
	return nObj, ok
}

// SearchAndDelete searches a nonce object by nonce, and if found, deletes it
// from the NonceDb.  Returns false if the nonce is not found or expired.
func (ndb *NonceDb) SearchAndDelete(nonce string) (*NonceObj, bool) {
	ndb.DeleteOldOportunistic()
	ndb.mutex.Lock()
	nObj, ok := ndb.nonceObjsByNonce[nonce]
	if !ok {
		ndb.mutex.Unlock()
		return nil, false
	}
	delete(ndb.nonceObjsByNonce, nonce)
	ndb.mutex.Unlock()
	if nObj.Expiration < time.Now().Unix() {
		ok = false
	}
	return nObj, ok
}

// DeleteOldOportunistic deletes expired nonces once every N calls (where N is
// 128 for now).
func (ndb *NonceDb) DeleteOldOportunistic() {
	mustDelete := false
	ndb.deleteCounterMutex.Lock()
	ndb.deleteCounter++
	if ndb.deleteCounter >= 128 {
		mustDelete = true
		ndb.deleteCounter = 0
	}
	ndb.deleteCounterMutex.Unlock()
	if mustDelete {
		ndb.DeleteOld()
	}
}

// DeleteOld deletes all the expired nonces.
func (ndb *NonceDb) DeleteOld() {
	now := time.Now().Unix()
	ndb.mutex.Lock()
	for {
		nObj := ndb.noncesHeap.Peek()
		if nObj == nil {
			break
		}
		if nObj.Expiration < now {
			ndb.noncesHeap.Pop()
			delete(ndb.nonceObjsByNonce, nObj.Nonce)
		} else {
			break
		}
	}
	ndb.mutex.Unlock()
}
