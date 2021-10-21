package core

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClaim(t *testing.T) {
	var schemaHash SchemaHash
	claim, err := NewClaim(schemaHash, WithFlagUpdatable(true))
	require.NoError(t, err)
	require.Zero(t, claim.value)
	for i := 1; i < 4; i++ {
		require.Zero(t, claim.index[i])
	}
	for i := 0; i < 32; i++ {
		if i == 16 {
			require.Equal(t, byte(0b10000), claim.index[0][i],
				int253ToString(claim.index[0]))
		} else {
			require.Zero(t, claim.index[0][i], int253ToString(claim.index[0]))
		}
	}

	dt, ok := claim.GetExpirationDate()
	require.True(t, dt.IsZero())
	require.False(t, ok)
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
	claim, err := NewClaim(schemaHash, WithFlagUpdatable(true))
	require.NoError(t, err)
	e := claim.TreeEntry()

	hi, hv, err := e.HiHv()
	require.NoError(t, err)

	hit, err := hi.MarshalText()
	require.NoError(t, err)
	require.Equal(t,
		"19905260441950906049955646784794273651462264973332746773406911374272567544299",
		string(hit))

	hvt, err := hv.MarshalText()
	require.NoError(t, err)
	require.Equal(t,
		"2351654555892372227640888372176282444150254868378439619268573230312091195718",
		string(hvt))
}

func TestClaim_GetSchemaHash(t *testing.T) {
	var sc SchemaHash
	n, err := rand.Read(sc[:])
	require.NoError(t, err)
	require.Equal(t, schemaHashLn, n)
	claim, err := NewClaim(sc)
	require.NoError(t, err)
	require.True(t, bytes.Equal(sc[:], claim.index[0][:schemaHashLn]))
}

func TestClaim_GetFlagUpdatable(t *testing.T) {
	var sc SchemaHash
	claim, err := NewClaim(sc)
	require.NoError(t, err)
	require.False(t, claim.GetFlagUpdatable())

	claim.SetFlagUpdatable(true)
	require.True(t, claim.GetFlagUpdatable())

	claim.SetFlagUpdatable(false)
	require.False(t, claim.GetFlagUpdatable())

	claim, err = NewClaim(sc, WithFlagUpdatable(true))
	require.NoError(t, err)
	require.True(t, claim.GetFlagUpdatable())

	claim, err = NewClaim(sc, WithFlagUpdatable(false))
	require.NoError(t, err)
	require.False(t, claim.GetFlagUpdatable())
}

func TestClaim_GetVersion(t *testing.T) {
	var sc SchemaHash
	ver := uint32(rand.Int63n(math.MaxUint32))
	claim, err := NewClaim(sc, WithVersion(ver))
	require.NoError(t, err)
	require.Equal(t, ver, claim.GetVersion())

	ver2 := uint32(rand.Int63n(math.MaxUint32))
	claim.SetVersion(ver2)
	require.Equal(t, ver2, claim.GetVersion())
}

func TestClaim_GetRevocationNonce(t *testing.T) {
	var sc SchemaHash
	nonce := uint64(rand.Int63())
	claim, err := NewClaim(sc, WithRevocationNonce(nonce))
	require.NoError(t, err)
	require.Equal(t, nonce, claim.GetRevocationNonce())

	nonce2 := uint64(rand.Int63())
	claim.SetRevocationNonce(nonce2)
	require.Equal(t, nonce2, claim.GetRevocationNonce())
}

func TestClaim_ExpirationDate(t *testing.T) {
	var sh SchemaHash
	expDate := time.Now().Truncate(time.Second)
	c1, err := NewClaim(sh, WithExpirationDate(expDate))
	require.NoError(t, err)

	expDate2, ok := c1.GetExpirationDate()
	require.True(t, ok)
	require.True(t, expDate2.Equal(expDate), "%v != %v", expDate, expDate2)

	c1.ResetExpirationDate()
	expDate2, ok = c1.GetExpirationDate()
	require.False(t, ok)
	require.True(t, expDate2.IsZero())

	expDate3 := expDate.Add(10 * time.Second)
	c1.SetExpirationDate(expDate3)
	expDate2, ok = c1.GetExpirationDate()
	require.True(t, ok)
	require.True(t, expDate2.Equal(expDate3), "%v != %v", expDate, expDate3)
}
