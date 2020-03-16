package writermock

import (
	"testing"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/stretchr/testify/require"
)

// Assert that IdenPubOffChainWriteMock follows the IdenPubOffChainWriter interface
func TestIdenPubOffChainWriteMockInterface(t *testing.T) {
	var idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter //nolint:gosimple
	idenPubOffChainWrite = New()
	require.NotNil(t, idenPubOffChainWrite)
}
