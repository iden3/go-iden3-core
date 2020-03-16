package local

import (
	"testing"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/stretchr/testify/require"
)

// Assert that IdenPubOffChainfollows the IdenPubOffChainReader and IdenPubOffChainWriter interface
func TestIdenPubOffChainInterface(t *testing.T) {
	var idenPubOffChainRead idenpuboffchain.IdenPubOffChainReader //nolint:gosimple
	idenPubOffChainRead = NewIdenPubOffChain("http://foo.bar")
	require.NotNil(t, idenPubOffChainRead)

	var idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter //nolint:gosimple
	idenPubOffChainWrite = NewIdenPubOffChain("http://foo.bar")
	require.NotNil(t, idenPubOffChainWrite)
}
