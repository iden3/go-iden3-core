package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/stretchr/testify/assert"
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

func (el ElemBytes) String() string {
	var b bytes.Buffer
	for j := len(el) - 1; j >= 0; j-- {
		b.WriteString(fmt.Sprintf("% 08b", el[j]))
	}
	return b.String()
}

func TestRawSlots(t *testing.T) {
	var schemaHash SchemaHash
	claim, err := NewClaim(schemaHash, WithFlagUpdatable(true))
	require.NoError(t, err)
	index, value := claim.RawSlots()
	indexHash, err := poseidon.Hash([]*big.Int{
		index[0].ToInt(), index[1].ToInt(), index[2].ToInt(), index[3].ToInt()})
	require.NoError(t, err)
	valueHash, err := poseidon.Hash([]*big.Int{
		value[0].ToInt(), value[1].ToInt(), value[2].ToInt(), value[3].ToInt()})
	require.NoError(t, err)

	require.Equal(t,
		"19905260441950906049955646784794273651462264973332746773406911374272567544299",
		indexHash.Text(10))

	require.Equal(t,
		"2351654555892372227640888372176282444150254868378439619268573230312091195718",
		valueHash.Text(10))
}

func TestClaim_GetSchemaHash(t *testing.T) {
	var sc SchemaHash
	n, err := rand.Read(sc[:])
	require.NoError(t, err)
	require.Equal(t, schemaHashLn, n)
	claim, err := NewClaim(sc)
	require.NoError(t, err)
	require.True(t,
		bytes.Equal(sc[:], claim.index[0][:schemaHashLn]))

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
	iX := toInt(t,
		"16243864111864693853212588481963275789994876191154110553066821559749894481761")
	iY := toInt(t,
		"7078462697308959301666117070269719819629678436794910510259518359026273676830")
	vX := toInt(t,
		"12448278679517811784508557734102986855579744384337338465055621486538311281772")
	vY := toInt(t,
		"9260608685281348956030279125705000716237952776955782848598673606545494194823")

	ixSlot, err := NewElemBytesFromInt(iX)
	require.NoError(t, err)
	iySlot, err := NewElemBytesFromInt(iY)
	require.NoError(t, err)
	vxSlot, err := NewElemBytesFromInt(vX)
	require.NoError(t, err)
	vySlot, err := NewElemBytesFromInt(vY)
	require.NoError(t, err)
	_, err = NewClaim(SchemaHash{},
		WithIndexData(ixSlot, iySlot),
		WithValueData(vxSlot, vySlot))
	require.NoError(t, err)
}

func TestNewDataSlotFromInt(t *testing.T) {
	ds, err := NewElemBytesFromInt(toInt(t,
		"16243864111864693853212588481963275789994876191154110553066821559749894481761"))
	require.NoError(t, err)
	expected := ElemBytes{
		0x61, 0x27, 0xa0, 0xeb, 0x58, 0x7a, 0x6c, 0x2b,
		0x4a, 0xa8, 0xc1, 0x2e, 0xf5, 0x01, 0xb2, 0xdb,
		0xd0, 0x9c, 0xb1, 0xa5, 0x9c, 0x83, 0x42, 0x57,
		0x91, 0xa5, 0x20, 0xbf, 0x86, 0xb3, 0xe9, 0x23,
	}
	require.Equal(t, expected, ds)

	_, err = NewElemBytesFromInt(toInt(t,
		"9916243864111864693853212588481963275789994876191154110553066821559749894481761"))
	require.EqualError(t, err, ErrDataOverflow.Error())
}

func TestClaim_WithIndexDataInts(t *testing.T) {

	expSlot := ElemBytes{}
	err := expSlot.SetInt(big.NewInt(0))
	require.NoError(t, err)

	value := new(big.Int).SetInt64(64)

	claim, err := NewClaim(SchemaHash{},
		WithIndexDataInts(value, nil))
	require.NoError(t, err)
	require.Equal(t, expSlot, claim.index[3])

	claim2, err := NewClaim(SchemaHash{},
		WithIndexDataInts(nil, value))
	require.NoError(t, err)
	require.Equal(t, expSlot, claim2.index[2])
}

func TestClaim_WithValueDataInts(t *testing.T) {

	expSlot := ElemBytes{}
	err := expSlot.SetInt(big.NewInt(0))
	require.NoError(t, err)

	value := new(big.Int).SetInt64(64)

	claim, err := NewClaim(SchemaHash{},
		WithValueDataInts(value, nil))
	require.NoError(t, err)

	require.Equal(t, expSlot, claim.value[3])

	claim2, err := NewClaim(SchemaHash{},
		WithValueDataInts(nil, value))
	require.NoError(t, err)
	require.Equal(t, expSlot, claim2.value[2])
}

func TestClaim_WithIndexDataBytes(t *testing.T) {

	iX := toInt(t,
		"124482786795178117845085577341029868555797443843373384650556214865383112817")
	expSlot := ElemBytes{}
	err := expSlot.SetInt(big.NewInt(0))
	require.NoError(t, err)

	claim, err := NewClaim(SchemaHash{},
		WithIndexDataBytes(iX.Bytes(), nil))
	require.NoError(t, err)
	require.Equal(t, expSlot, claim.index[3])
}

func TestClaimJSONSerialization(t *testing.T) {
	in := `[
"15163995036539824738096525342132337704181738148399168403057770094395141110111",
"3206594817839378626027676511482956481343861686313501795018892230311002175077",
"7420031054231607091230846181053275837604749850669737756447914128096832575029",
"6843256246667081694694856844555135410358903741435158507252727716055448769466",
"18335061644187980192028500482619331449203987338928612566250871337402164885236",
"4747739418092571675618239353368909204965774632269590366651599441049750269324",
"10060277146294090095035892104009266064127776406104429246320070556972379481946",
"5835715034681704899254417398745238273415614452113785384300119694985241103333"
]`
	var want = Claim{
		index: [4]ElemBytes{
			slotFromHex("5fb90badb37c5821b6d95526a41a9504680b4e7c8b763a1b1d49d4955c848621"),
			slotFromHex("65f606f6a63b7f3dfd2567c18979e4d60f26686d9bf2fb26c901ff354cde1607"),
			slotFromHex("35d6042c4160f38ee9e2a9f3fb4ffb0019b454d522b5ffa17604193fb8966710"),
			slotFromHex("ba53af19779cb2948b6570ffa0b773963c130ad797ddeafe4e3ad29b5125210f"),
		},
		value: [4]ElemBytes{
			slotFromHex("f4b6f44090a32711f3208e4e4b89cb5165ce64002cbd9c2887aa113df2468928"),
			slotFromHex("8ced323cb76f0d3fac476c9fb03fc9228fbae88fd580663a0454b68312207f0a"),
			slotFromHex("5a27db029de37ae37a42318813487685929359ca8c5eb94e152dc1af42ea3d16"),
			slotFromHex("e50be1a6dc1d5768e8537988fddce562e9b948c918bba3e933e5c400cde5e60c"),
		},
	}

	t.Run("unmarshal", func(t *testing.T) {
		var result Claim
		err := json.Unmarshal([]byte(in), &result)
		require.NoError(t, err)
		require.Equal(t, want, result)
	})

	t.Run("marshal", func(t *testing.T) {
		result, err := json.Marshal(want)
		require.NoError(t, err)
		require.JSONEq(t, in, string(result))
	})
}

func slotFromHex(in string) ElemBytes {
	var eb ElemBytes
	data, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}
	if len(data) != 32 {
		panic("data len not 32")
	}
	copy(eb[:], data)
	return eb
}

func TestGenerateRandomSlots(t *testing.T) {
	t.Skip("use this to generate random slots for other tests")

	var dataRnd [32]byte
	watchdog := 0
	cnt := 0
	for {
		watchdog++
		if watchdog > 1000 {
			t.Fatal(watchdog)
		}

		n, err := rand.Read(dataRnd[:])
		require.NoError(t, err)
		require.Equal(t, 32, n)
		intVal := bytesToInt(dataRnd[:])
		if !utils.CheckBigIntInField(intVal) {
			continue
		}

		t.Logf("%v: %v", hex.EncodeToString(dataRnd[:]), intVal.Text(10))

		cnt++
		if cnt >= 8 {
			break
		}
	}
}

func TestClaimBinarySerialization(t *testing.T) {
	binDataStr := strings.Join([]string{
		"5fb90badb37c5821b6d95526a41a9504680b4e7c8b763a1b1d49d4955c848621",
		"65f606f6a63b7f3dfd2567c18979e4d60f26686d9bf2fb26c901ff354cde1607",
		"35d6042c4160f38ee9e2a9f3fb4ffb0019b454d522b5ffa17604193fb8966710",
		"ba53af19779cb2948b6570ffa0b773963c130ad797ddeafe4e3ad29b5125210f",
		"f4b6f44090a32711f3208e4e4b89cb5165ce64002cbd9c2887aa113df2468928",
		"8ced323cb76f0d3fac476c9fb03fc9228fbae88fd580663a0454b68312207f0a",
		"5a27db029de37ae37a42318813487685929359ca8c5eb94e152dc1af42ea3d16",
		"e50be1a6dc1d5768e8537988fddce562e9b948c918bba3e933e5c400cde5e60c",
	}, "")
	binData, err := hex.DecodeString(binDataStr)
	require.NoError(t, err)

	var want = Claim{
		index: [4]ElemBytes{
			slotFromHex("5fb90badb37c5821b6d95526a41a9504680b4e7c8b763a1b1d49d4955c848621"),
			slotFromHex("65f606f6a63b7f3dfd2567c18979e4d60f26686d9bf2fb26c901ff354cde1607"),
			slotFromHex("35d6042c4160f38ee9e2a9f3fb4ffb0019b454d522b5ffa17604193fb8966710"),
			slotFromHex("ba53af19779cb2948b6570ffa0b773963c130ad797ddeafe4e3ad29b5125210f"),
		},
		value: [4]ElemBytes{
			slotFromHex("f4b6f44090a32711f3208e4e4b89cb5165ce64002cbd9c2887aa113df2468928"),
			slotFromHex("8ced323cb76f0d3fac476c9fb03fc9228fbae88fd580663a0454b68312207f0a"),
			slotFromHex("5a27db029de37ae37a42318813487685929359ca8c5eb94e152dc1af42ea3d16"),
			slotFromHex("e50be1a6dc1d5768e8537988fddce562e9b948c918bba3e933e5c400cde5e60c"),
		},
	}

	t.Run("unmarshal", func(t *testing.T) {
		var result Claim
		err := result.UnmarshalBinary(binData)
		require.NoError(t, err)
		require.Equal(t, want, result)
	})

	t.Run("marshal", func(t *testing.T) {
		result, err := want.MarshalBinary()
		require.NoError(t, err)
		require.Equal(t, binData, result)
	})
}

func TestNewSchemaHashFromHex(t *testing.T) {

	hash := "ca938857241db9451ea329256b9c06e5"
	got, err := NewSchemaHashFromHex(hash)
	require.NoError(t, err)

	exp, err := hex.DecodeString(hash)
	require.NoError(t, err)

	assert.Equal(t, exp[:], got[:])

}

func TestSchemaHash_BigInt(t *testing.T) {
	schema, err := NewSchemaHashFromHex("ca938857241db9451ea329256b9c06e5")
	require.NoError(t, err)

	exp, b := new(big.Int).SetString("304427537360709784173770334266246861770", 10)
	require.True(t, b)

	got := schema.BigInt()

	assert.Equal(t, exp, got)

}

func TestGetIDPosition(t *testing.T) {
	tests := []struct {
		name             string
		claim            func(t *testing.T) *Claim
		expectedPosition IDPosition
	}{
		{
			name: "self claim",
			claim: func(t *testing.T) *Claim {
				c, err := NewClaim(SchemaHash{})
				require.NoError(t, err)
				return c
			},
			expectedPosition: IDPositionNone,
		},
		{
			name: "subject stored in index",
			claim: func(t *testing.T) *Claim {
				c, err := NewClaim(SchemaHash{})
				require.NoError(t, err)

				var genesis [27]byte
				genesis32bytes := hashBytes([]byte("genesistest"))
				copy(genesis[:], genesis32bytes[:])

				c.SetIndexID(NewID(TypeDefault, genesis))
				return c
			},
			expectedPosition: IDPositionIndex,
		},
		{
			name: "subject stored in value",
			claim: func(t *testing.T) *Claim {
				c, err := NewClaim(SchemaHash{})
				require.NoError(t, err)

				var genesis [27]byte
				genesis32bytes := hashBytes([]byte("genesistest"))
				copy(genesis[:], genesis32bytes[:])

				c.SetValueID(NewID(TypeDefault, genesis))
				return c
			},
			expectedPosition: IDPositionValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.claim(t)
			position, err := c.GetIDPosition()
			require.NoError(t, err)
			require.Equal(t, tt.expectedPosition, position)
		})
	}
}

func TestGetIDPosition_ErrorCase(t *testing.T) {
	tests := []struct {
		name             string
		claim            func(t *testing.T) *Claim
		expectedPosition IDPosition
		expectedError    error
	}{
		{
			name: "invalid position",
			claim: func(t *testing.T) *Claim {
				c, err := NewClaim(SchemaHash{})
				require.NoError(t, err)
				c.setSubject(_subjectFlagInvalid)
				return c
			},
			expectedPosition: IDPositionNone,
			expectedError:    ErrInvalidSubjectPosition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.claim(t)
			position, err := c.GetIDPosition()
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedPosition, position)
		})
	}
}

func TestGetMerklizePosition(t *testing.T) {
	tests := []struct {
		name             string
		claim            func(t *testing.T) *Claim
		expectedPosition MerklizedRootPosition
	}{
		{
			name: "not merklized",
			claim: func(t *testing.T) *Claim {
				c, err := NewClaim(SchemaHash{})
				require.NoError(t, err)
				return c
			},
			expectedPosition: MerklizedRootPositionNone,
		},
		{
			name: "mt root stored in index",
			claim: func(t *testing.T) *Claim {
				c, err := NewClaim(SchemaHash{})
				require.NoError(t, err)

				c.setFlagMerklized(MerklizedRootPositionIndex)
				return c
			},
			expectedPosition: MerklizedRootPositionIndex,
		},
		{
			name: "mt root stored in value",
			claim: func(t *testing.T) *Claim {
				c, err := NewClaim(SchemaHash{})
				require.NoError(t, err)

				c.setFlagMerklized(MerklizedRootPositionValue)
				return c
			},
			expectedPosition: MerklizedRootPositionValue,
		},
		{
			name: "mt root random bits",
			claim: func(t *testing.T) *Claim {
				c, err := NewClaim(SchemaHash{})
				require.NoError(t, err)

				c.setFlagMerklized(MerklizedRootPositionValue)
				return c
			},
			expectedPosition: MerklizedRootPositionValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.claim(t)
			position, err := c.GetMerklizedPosition()
			require.NoError(t, err)
			require.Equal(t, tt.expectedPosition, position)
		})
	}
}

func TestGetMerklizePosition_ErrorCase(t *testing.T) {
	c, err := NewClaim(SchemaHash{})
	require.NoError(t, err)
	c.index[0][flagsByteIdx] &= 0b11111000
	c.index[0][flagsByteIdx] |= byte(_merklizedFlagInvalid)

	position, err := c.GetMerklizedPosition()
	require.ErrorIs(t, err, ErrIncorrectMerklizedPosition)
	require.Equal(t, 0, int(position))
}

func TestWithFlagMerklized(t *testing.T) {
	claim, err := NewClaim(SchemaHash{},
		WithFlagMerklized(MerklizedRootPositionIndex))
	require.NoError(t, err)

	require.Equal(t, byte(merklizedFlagIndex), claim.index[0][flagsByteIdx]&0b11100000)
}

func TestWithIndexMerklizedRoot(t *testing.T) {
	expVal := big.NewInt(9999)
	expSlot := ElemBytes{}
	err := expSlot.SetInt(expVal)
	require.NoError(t, err)

	claim, err := NewClaim(SchemaHash{},
		WithIndexMerklizedRoot(expVal))
	require.NoError(t, err)
	require.Equal(t, expSlot, claim.index[2])

	position, err := claim.GetMerklizedPosition()
	require.NoError(t, err)
	require.Equal(t, MerklizedRootPositionIndex, position)
}

func TestWithValueMerklizedRoot(t *testing.T) {
	expVal := big.NewInt(9999)
	expSlot := ElemBytes{}
	err := expSlot.SetInt(expVal)
	require.NoError(t, err)

	claim, err := NewClaim(SchemaHash{},
		WithValueMerklizedRoot(expVal))
	require.NoError(t, err)
	require.Equal(t, expSlot, claim.value[2])

	position, err := claim.GetMerklizedPosition()
	require.NoError(t, err)
	require.Equal(t, MerklizedRootPositionValue, position)
}

func TestWithMerklizedRoot(t *testing.T) {
	expVal := big.NewInt(9999)
	expSlot := ElemBytes{}
	err := expSlot.SetInt(expVal)
	require.NoError(t, err)

	claim, err := NewClaim(SchemaHash{},
		WithMerklizedRoot(expVal, MerklizedRootPositionIndex))
	require.NoError(t, err)
	require.Equal(t, expSlot, claim.index[2])

	position, err := claim.GetMerklizedPosition()
	require.NoError(t, err)
	require.Equal(t, MerklizedRootPositionIndex, position)

	claim2, err := NewClaim(SchemaHash{},
		WithMerklizedRoot(expVal, MerklizedRootPositionValue))
	require.NoError(t, err)
	require.Equal(t, expSlot, claim2.value[2])

	position2, err := claim2.GetMerklizedPosition()
	require.NoError(t, err)
	require.Equal(t, MerklizedRootPositionValue, position2)
}

func TestClaim_SetMerklizedRoot(t *testing.T) {
	expVal := big.NewInt(9999)
	expSlot := ElemBytes{}
	err := expSlot.SetInt(expVal)
	require.NoError(t, err)

	claim, err := NewClaim(SchemaHash{})
	require.NoError(t, err)

	err = claim.SetIndexMerklizedRoot(expVal)
	require.NoError(t, err)
	require.Equal(t, expSlot, claim.index[2])

	position, err := claim.GetMerklizedPosition()
	require.NoError(t, err)
	require.Equal(t, MerklizedRootPositionIndex, position)

	r, err := claim.GetMerklizedRoot()
	require.NoError(t, err)
	require.Equal(t, expVal, r)

	claim2, err := NewClaim(SchemaHash{})
	require.NoError(t, err)

	err = claim2.SetValueMerklizedRoot(expVal)
	require.NoError(t, err)
	require.Equal(t, expSlot, claim2.value[2])

	position2, err := claim2.GetMerklizedPosition()
	require.NoError(t, err)
	require.Equal(t, MerklizedRootPositionValue, position2)

	r2, err := claim2.GetMerklizedRoot()
	require.NoError(t, err)
	require.Equal(t, expVal, r2)

	claim3, err := NewClaim(SchemaHash{})
	require.NoError(t, err)

	position3, err := claim3.GetMerklizedPosition()
	require.NoError(t, err)
	require.Equal(t, MerklizedRootPositionNone, position3)

	_, err = claim3.GetMerklizedRoot()
	require.Error(t, ErrNoMerklizedRoot, err)
}
