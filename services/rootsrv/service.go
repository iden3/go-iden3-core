package rootsrv

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/merkletree"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	Start()
	StopAndJoin()
	GetRoot(addr common.Address) (merkletree.Hash, error)
	SetRoot(hash merkletree.Hash)
}

type ServiceImpl struct {
	lastRoot      merkletree.Hash
	lastRootMutex sync.RWMutex
	stopch        chan (interface{})
	stoppedch     chan (interface{})
	rootcommits   *eth.Contract
}

func New(rootcommits *eth.Contract) *ServiceImpl {
	return &ServiceImpl{
		stopch:      make(chan (interface{})),
		stoppedch:   make(chan (interface{})),
		rootcommits: rootcommits,
		lastRoot:    merkletree.Hash{},
	}
}

func (s *ServiceImpl) Start() {
	go func() {
		s.lastRootMutex.RLock()
		lastRoot := s.lastRoot
		s.lastRootMutex.RUnlock()
		log.Info("Starting root publisher")
		for {
			select {
			case <-s.stopch:
				log.Info("Root publisher finalized")
				s.stoppedch <- nil
				return
			case <-time.After(time.Second):
				s.lastRootMutex.RLock()
				sLastRoot := s.lastRoot
				s.lastRootMutex.RUnlock()
				if lastRoot != sLastRoot {
					lastRoot = sLastRoot
					log.Debugf("Upading root in smart contract to %v\n",
						common3.HexEncode(lastRoot[:]))
					if tx, err := s.rootcommits.SendTransaction(nil, 0,
						"setRoot", lastRoot); err != nil {
						log.Error("Failed to add root: ", err)
						lastRoot = merkletree.Hash{}
					} else {
						_, err = s.rootcommits.Client().WaitReceipt(tx.Hash())
						if err != nil {
							log.Error("Error waiting for receipt: ", err)
						}
					}
				}
			}
		}
	}()
}

func (s *ServiceImpl) GetRoot(addr common.Address) (merkletree.Hash, error) {
	var res merkletree.Hash
	err := s.rootcommits.Call(&res, "getRoot", addr)
	return res, err
}

func (s *ServiceImpl) SetRoot(hash merkletree.Hash) {
	s.lastRootMutex.Lock()
	s.lastRoot = hash
	s.lastRootMutex.Unlock()
}

func (s *ServiceImpl) StopAndJoin() {
	go func() {
		s.stopch <- nil
	}()
	<-s.stoppedch
}
