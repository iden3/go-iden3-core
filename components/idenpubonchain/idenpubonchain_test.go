package idenpubonchain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Assert that IdenPubOnChain follows the IdenPubOnChainer interface
func TestIdenPubOnChainInterface(t *testing.T) {
	var idenPubOnChain IdenPubOnChainer //nolint:gosimple
	idenPubOnChain = New(nil, ContractAddresses{})
	require.NotNil(t, idenPubOnChain)
}
