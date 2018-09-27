package rootsrv

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	lastRoot    merkletree.Hash
	stopch      chan (interface{})
	stoppedch   chan (interface{})
	rootcommits *eth.Contract
}

func New(rootcommits *eth.Contract) *ServiceImpl {
	return &ServiceImpl{
		stopch:      make(chan (interface{})),
		rootcommits: rootcommits,
	}
}

func (s *ServiceImpl) Start() {
	s.lastRoot = merkletree.Hash{}
	lastRoot := s.lastRoot
	go func() {
		log.Info("Starting root publisher")
		for {
			switch {
			case <-s.stopch:
				s.stoppedch <- nil
				break
			case lastRoot != s.lastRoot:
				lastRoot = s.lastRoot
				_, _, err := s.rootcommits.SendTransactionSync(nil, 0, "AddRoot", lastRoot)
				if err != nil {
					log.Error("Failed to add root", err)
					lastRoot = merkletree.Hash{}
					time.Sleep(1)
				}
			default:
				time.Sleep(1)
			}
		}
		log.Info("Root publisher finialized")
	}()
}

func (s *ServiceImpl) GetRoot(addr common.Address) (merkletree.Hash, error) {
	var res merkletree.Hash
	err := s.rootcommits.Call(&res, "GetRoot", addr)
	return res, err
}

func (s *ServiceImpl) SetRoot(hash merkletree.Hash) {
	s.lastRoot = hash
}

func (s *ServiceImpl) StopAndJoin() {
	s.stopch <- nil
	<-s.stoppedch
}
