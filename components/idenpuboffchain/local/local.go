package local

import (
	"fmt"
	"sync"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-merkletree-sql"
)

type IdenPubOffChainIden struct {
	datas map[merkletree.Hash]*idenpuboffchain.PublicData
	last  *merkletree.Hash
}

type IdenPubOffChain struct {
	rw    sync.RWMutex
	idens map[core.ID]*IdenPubOffChainIden
	url   string
}

func NewIdenPubOffChain(url string) *IdenPubOffChain {
	return &IdenPubOffChain{
		idens: make(map[core.ID]*IdenPubOffChainIden),
		url:   url,
	}
}

func (ip *IdenPubOffChain) GetPublicData(idenPubUrl string, id *core.ID, idenState *merkletree.Hash) (*idenpuboffchain.PublicData, error) {
	if ip.url != idenPubUrl {
		return nil, fmt.Errorf("Bad URL")
	}
	ip.rw.RLock()
	defer ip.rw.RUnlock()
	iden, ok := ip.idens[*id]
	if !ok {
		return nil, fmt.Errorf("ID not found")
	}
	if idenState == nil {
		idenState = iden.last
	}
	data, ok := iden.datas[*idenState]
	if !ok {
		return nil, fmt.Errorf("idenState not found")
	}
	return data, nil
}

func (ip *IdenPubOffChain) Url() string {
	return ip.url
}

func (ip *IdenPubOffChain) Publish(id *core.ID, publicData *idenpuboffchain.PublicData) error {
	ip.rw.Lock()
	defer ip.rw.Unlock()
	iden, ok := ip.idens[*id]
	if !ok {
		iden = &IdenPubOffChainIden{
			datas: make(map[merkletree.Hash]*idenpuboffchain.PublicData),
			last:  nil,
		}
		ip.idens[*id] = iden
	}
	iden.datas[*publicData.IdenState] = publicData
	iden.last = publicData.IdenState
	return nil
}
