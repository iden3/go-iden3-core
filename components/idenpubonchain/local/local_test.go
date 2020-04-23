package local

import (
	"testing"
	"time"

	zktypes "github.com/iden3/go-circom-prover-verifier/types"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/stretchr/testify/require"
)

// Assert that IdenPubOnChain follows the IdenPubOnChainer interface
func TestLocalIdenPubOnChainInterface(t *testing.T) {
	var idenPubOnChain idenpubonchain.IdenPubOnChainer //nolint:gosimple
	idenPubOnChain = New(
		func() time.Time {
			return time.Now()
		},
		func() uint64 {
			return 0
		},
		&zktypes.Vk{},
	)
	require.NotNil(t, idenPubOnChain)
}
