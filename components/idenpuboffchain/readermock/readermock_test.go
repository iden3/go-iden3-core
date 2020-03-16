package readermock

import (
	"testing"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/stretchr/testify/require"
)

// Assert that IdenPubOffChainReadMock follows the IdenPubOffChainReader interface
func TestIdenPubOffChainReadMockInterface(t *testing.T) {
	var idenPubOffChainRead idenpuboffchain.IdenPubOffChainReader //nolint:gosimple
	idenPubOffChainRead = New()
	require.NotNil(t, idenPubOffChainRead)
}
