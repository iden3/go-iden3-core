package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetworkByChainID(t *testing.T) {
	chainID, err := GetChainID(Ethereum, Main)
	require.NoError(t, err)
	blockchain, networkID, err := NetworkByChainID(chainID)
	require.NoError(t, err)
	require.Equal(t, Ethereum, blockchain)
	require.Equal(t, Main, networkID)
}
