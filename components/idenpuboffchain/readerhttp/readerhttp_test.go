package readerhttp

import (
	"testing"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/stretchr/testify/require"
)

// Assert that IdenPubOffChainReadHttp follows the IdenPubOffChainReader interface
func TestIdenPubOffChainReadHttpInterface(t *testing.T) {
	var idenPubOffChainRead idenpuboffchain.IdenPubOffChainReader //nolint:gosimple
	idenPubOffChainRead = NewIdenPubOffChainHttp()
	require.NotNil(t, idenPubOffChainRead)
}
