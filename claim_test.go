package core

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClaim(t *testing.T) {
	var schemaHash SchemaHash
	claim, err := NewClaim(schemaHash, WithFlagExpiration(true))
	require.NoError(t, err)
	require.Zero(t, claim.value)
	for i := 1; i < 4; i++ {
		require.Zero(t, claim.index[i])
	}
	for i := 0; i < 32; i++ {
		if i == 16 {
			require.Equal(t, byte(0b1000), claim.index[0][i],
				int253ToString(claim.index[0]))
		} else {
			require.Zero(t, claim.index[0][i], int253ToString(claim.index[0]))
		}
	}
}

func int253ToString(i int253) string {
	var b bytes.Buffer
	for j := len(i) - 1; j >= 0; j-- {
		b.WriteString(fmt.Sprintf("% 08b", i[j]))
	}
	return b.String()
}

func TestMerketreeEntryHash(t *testing.T) {
	var schemaHash SchemaHash
	claim, err := NewClaim(schemaHash, WithFlagExpiration(true))
	require.NoError(t, err)
	e := claim.TreeEntry()

	hi, hv, err := e.HiHv()
	require.NoError(t, err)

	hit, err := hi.MarshalText()
	require.NoError(t, err)
	require.Equal(t,
		"19580667809762269956733050858122189693671367180755664752787058228636503413656",
		string(hit))

	hvt, err := hv.MarshalText()
	require.NoError(t, err)
	require.Equal(t,
		"2351654555892372227640888372176282444150254868378439619268573230312091195718",
		string(hvt))
}
