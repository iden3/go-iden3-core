package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/iden3/go-iden3-crypto/utils"

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
				claim.index[0].String())
		} else {
			require.Zero(t, claim.index[0][i], claim.index[0].String())
		}
	}

	dt, ok := claim.GetExpirationDate()
	require.True(t, dt.IsZero())
	require.False(t, ok)
}

func (ds DataSlot) String() string {
	var b bytes.Buffer
	for j := len(ds) - 1; j >= 0; j-- {
		b.WriteString(fmt.Sprintf("% 08b", ds[j]))
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
	require.True(t, bytes.Equal(sc[:], utils.SwapEndianness(claim.index[0][:schemaHashLn])))

	shFromClaim := claim.GetSchemaHash()
	shFromClaimHexBytes, err := shFromClaim.MarshalText()
	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(sc[:]), string(shFromClaimHexBytes))

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

func toInt(t testing.TB, s string) *big.Int {
	t.Helper()
	i, ok := new(big.Int).SetString(s, 10)
	require.True(t, ok, s)
	return i
}

func TestIntSize(t *testing.T) {
	iX := toInt(t, "16243864111864693853212588481963275789994876191154110553066821559749894481761")
	iY := toInt(t, "7078462697308959301666117070269719819629678436794910510259518359026273676830")
	vX := toInt(t, "12448278679517811784508557734102986855579744384337338465055621486538311281772")
	vY := toInt(t, "9260608685281348956030279125705000716237952776955782848598673606545494194823")

	ixSlot, err := NewDataSlotFromInt(iX)
	require.NoError(t, err)
	iySlot, err := NewDataSlotFromInt(iY)
	require.NoError(t, err)
	vxSlot, err := NewDataSlotFromInt(vX)
	require.NoError(t, err)
	vySlot, err := NewDataSlotFromInt(vY)
	require.NoError(t, err)
	_, err = NewClaim(SchemaHash{},
		WithIndexData(ixSlot, iySlot),
		WithValueData(vxSlot, vySlot))
	require.NoError(t, err)
}

func TestNewDataSlotFromInt(t *testing.T) {
	ds, err := NewDataSlotFromInt(toInt(t,
		"16243864111864693853212588481963275789994876191154110553066821559749894481761"))
	require.NoError(t, err)
	expected := DataSlot{
		0x61, 0x27, 0xa0, 0xeb, 0x58, 0x7a, 0x6c, 0x2b,
		0x4a, 0xa8, 0xc1, 0x2e, 0xf5, 0x01, 0xb2, 0xdb,
		0xd0, 0x9c, 0xb1, 0xa5, 0x9c, 0x83, 0x42, 0x57,
		0x91, 0xa5, 0x20, 0xbf, 0x86, 0xb3, 0xe9, 0x23,
	}
	require.Equal(t, expected, ds)

	_, err = NewDataSlotFromInt(toInt(t,
		"9916243864111864693853212588481963275789994876191154110553066821559749894481761"))
	require.EqualError(t, err, ErrDataOverflow.Error())
}
