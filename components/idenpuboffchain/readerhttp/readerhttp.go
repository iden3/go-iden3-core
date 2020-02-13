package readerhttp

import (
	"fmt"

	"github.com/iden3/go-iden3-core/components/httpclient"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
)

// IdenPubOffChainReader is a interface to read the off chain public state of an identity.
type IdenPubOffChainReader interface {
	GetPublicData(idenPubUrl string, id *core.ID, idenState *merkletree.Hash) (*idenpuboffchain.PublicData, error)
}

// IdenPubOffChainReadHttp satisfies the IdenPubOffChainRead interface, and reads the off chain public state of an identity from a IdenPubOffChainWriteHttp.
type IdenPubOffChainReadHttp struct {
}

func NewIdenPubOffChainHttp() *IdenPubOffChainReadHttp {
	return &IdenPubOffChainReadHttp{}
}

func (i *IdenPubOffChainReadHttp) GetPublicData(idenPubUrl string, id *core.ID, idenState *merkletree.Hash) (*idenpuboffchain.PublicData, error) {
	httpClient := httpclient.NewHttpClient(idenPubUrl)

	var publicDataBlobs idenpuboffchain.PublicDataBlobs
	var err error
	if idenState != nil {
		err = httpClient.DoRequest(httpClient.NewRequest().Path(
			fmt.Sprintf("%s/state/%s", id.String(), idenState.Hex())).Get(""), &publicDataBlobs)
	} else {
		err = httpClient.DoRequest(httpClient.NewRequest().Path(
			fmt.Sprintf("%s/laststate", id.String())).Get(""), &publicDataBlobs)
	}
	publicData, err := idenpuboffchain.NewPublicDataFromBlobs(&publicDataBlobs)
	if err != nil {
		return nil, err
	}
	return publicData, nil
}
