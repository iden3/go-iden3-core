package idenpuboffchainreader

import (
	"fmt"

	"github.com/iden3/go-iden3-core/components/idenpuboffchainwriter"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
)

// IdenPubOffChainReader is a interface to read the off chain public state of an identity.
type IdenPubOffChainReader interface {
	Publish()
}

// IdenPubOffChainReadHttp satisfies the IdenPubOffChainRead interface, and reads the off chain public state of an identity from a IdenPubOffChainWriteHttp.
type IdenPubOffChainReadHttp struct {
}

func (i *IdenPubOffChainReadHttp) GetPublicData(idPubUrl string, id *core.ID, idenState *merkletree.Hash) (*idenpuboffchainwriter.PublicData, error) {

	return nil, fmt.Errorf("TODO")
}
