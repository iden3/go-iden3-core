package rootsrv

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/contracts"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/merkletree"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	Start()
	StopAndJoin()
	// GetRoot(addr common.Address) (merkletree.Hash, error)
	GetRoot(id *core.ID) (merkletree.Hash, error)
	SetRoot(hash merkletree.Hash)
}

type ServiceImpl struct {
	lastRoot      merkletree.Hash
	lastRootMutex sync.RWMutex
	stopch        chan (interface{})
	stoppedch     chan (interface{})
	// rootcommits    *eth.Contract
	client         *eth.Client2
	id             *core.ID
	kUpdateRootMtp []byte
	contractAddr   common.Address
}

func New(client *eth.Client2, id *core.ID, kUpdateRootMtp []byte, contractAddr common.Address) *ServiceImpl {
	return &ServiceImpl{
		stopch:         make(chan (interface{})),
		stoppedch:      make(chan (interface{})),
		lastRoot:       merkletree.Hash{},
		client:         client,
		id:             id,
		kUpdateRootMtp: kUpdateRootMtp,
		contractAddr:   contractAddr,
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
					if err := s.updateRoot(lastRoot); err != nil {
						log.Error(err)
						lastRoot = merkletree.Hash{}
					}

				}
			}
		}
	}()
}

// func (s *ServiceImpl) GetRootOld(addr common.Address) (merkletree.Hash, error) {
// 	var res merkletree.Hash
// 	err := s.rootcommits.Call(&res, "getRoot", addr)
// 	return res, err
// }

func (s *ServiceImpl) GetRoot(id *core.ID) (merkletree.Hash, error) {
	var res [32]byte
	err := s.client.Call(func(c *ethclient.Client) error {
		rootcommits, err := contracts.NewRootCommits(s.contractAddr, c)
		if err != nil {
			return err
		}
		res, err = rootcommits.GetRoot(nil, *id)
		return err
	})
	return res, err
}

func (s *ServiceImpl) SetRoot(hash merkletree.Hash) {
	s.lastRootMutex.Lock()
	s.lastRoot = hash
	s.lastRootMutex.Unlock()
}

func (s *ServiceImpl) updateRoot(hash merkletree.Hash) error {
	if tx, err := s.client.CallAuth(
		func(c *ethclient.Client, auth *bind.TransactOpts) (*types.Transaction, error) {
			rootcommits, err := contracts.NewRootCommits(s.contractAddr, c)
			if err != nil {
				return nil, err
			}
			return rootcommits.SetRoot(auth, hash, *s.id, s.kUpdateRootMtp)
		},
	); err != nil {
		return fmt.Errorf("Failed to add root: %v", err)
	} else {
		_, err = s.client.WaitReceipt(tx)
		if err != nil {
			return fmt.Errorf("Error waiting for receipt: %v", err)
		}
	}
	return nil
}

func (s *ServiceImpl) StopAndJoin() {
	go func() {
		s.stopch <- nil
	}()
	<-s.stoppedch
}
