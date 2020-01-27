package idenstatewriter

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"

	// "github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/components/idenstatereader"
	"github.com/iden3/go-iden3-core/eth/contracts"
	"github.com/iden3/go-iden3-core/merkletree"
	log "github.com/sirupsen/logrus"
)

type IdenStateWriter interface {
	Start()
	StopAndJoin()
	// GetRoot(addr common.Address) (merkletree.Hash, error)
	GetRoot(id *core.ID) (*proof.RootData, error)
	SetRoot(hash merkletree.Hash)
}

type IdenStateWrite struct {
	lastRoot      merkletree.Hash
	lastRootMutex sync.RWMutex
	stopch        chan (interface{})
	stoppedch     chan (interface{})
	// rootcommits    *eth.Contract
	idenstatereader.IdenStateReader
	id             *core.ID
	kUpdateRootMtp []byte
	contractAddr   common.Address
}

func New(idenStateReader idenstatereader.IdenStateReader, id *core.ID, kUpdateRootMtp []byte, contractAddr common.Address) *IdenStateWrite {
	return &IdenStateWrite{
		stopch:          make(chan (interface{})),
		stoppedch:       make(chan (interface{})),
		lastRoot:        merkletree.Hash{},
		IdenStateReader: idenStateReader,
		id:              id,
		kUpdateRootMtp:  kUpdateRootMtp,
		contractAddr:    contractAddr,
	}
}

func (s *IdenStateWrite) Start() {
	go func() {
		lastRoot := s.getLastRoot()
		log.Info("Starting root publisher")
		for {
			select {
			case <-s.stopch:
				log.Info("Root publisher finalized")
				s.stoppedch <- nil
				return
			case <-time.After(time.Second):
				sLastRoot := s.getLastRoot()
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

func (s *IdenStateWrite) SetRoot(hash merkletree.Hash) {
	s.lastRootMutex.Lock()
	s.lastRoot = hash
	s.lastRootMutex.Unlock()
}

func (s *IdenStateWrite) getLastRoot() (hash merkletree.Hash) {
	s.lastRootMutex.RLock()
	defer s.lastRootMutex.RUnlock()
	return s.lastRoot
}

func (s *IdenStateWrite) updateRoot(hash merkletree.Hash) error {
	if tx, err := s.Client().CallAuth(
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
		_, err = s.Client().WaitReceipt(tx)
		if err != nil {
			return fmt.Errorf("Error waiting for receipt: %v", err)
		}
	}
	return nil
}

func (s *IdenStateWrite) StopAndJoin() {
	go func() {
		s.stopch <- nil
	}()
	<-s.stoppedch
}
