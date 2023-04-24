package core

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDID(t *testing.T) {

	// did
	didStr := "did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"

	did3, err := Parse(didStr)
	require.NoError(t, err)

	blockchain, networkID, id, err := Decompose(*did3)
	require.NoError(t, err)
	require.Equal(t, "wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ", id.String())
	require.Equal(t, Polygon, blockchain)
	require.Equal(t, Mumbai, networkID)

	// readonly did
	didStr = "did:iden3:readonly:tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa"

	did3, err = Parse(didStr)
	require.NoError(t, err)

	blockchain, networkID, id, err = Decompose(*did3)
	require.NoError(t, err)

	require.Equal(t, "tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa", id.String())
	require.Equal(t, ReadOnly, blockchain)
	require.Equal(t, NoNetwork, networkID)

	require.Equal(t, [2]byte{DIDMethodByte[DIDMethodIden3], 0b0}, id.Type())
}

func TestDID_MarshalJSON(t *testing.T) {
	id, err := IDFromString("wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ")
	require.NoError(t, err)
	did, err := ParseDIDFromID(id)
	require.NoError(t, err)

	b, err := did.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t,
		`"did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"`,
		string(b))
}

func TestDID_UnmarshalJSON(t *testing.T) {
	inBytes := `{"obj": "did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"}`
	id, err := IDFromString("wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ")
	require.NoError(t, err)
	var obj struct {
		Obj *DID `json:"obj"`
	}
	err = json.Unmarshal([]byte(inBytes), &obj)
	require.NoError(t, err)
	require.NotNil(t, obj.Obj)
	require.Equal(t, string(DIDMethodIden3), obj.Obj.Method)
	blockchain, networkID, id2, err := Decompose(*obj.Obj)
	require.NoError(t, err)
	require.Equal(t, id, id2)
	require.Equal(t, Polygon, blockchain)
	require.Equal(t, Mumbai, networkID)
}

func TestDID_UnmarshalJSON_Error(t *testing.T) {
	inBytes := `{"obj": "did:iden3:eth:goerli:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"}`
	var obj struct {
		Obj *DID `json:"obj"`
	}
	err := json.Unmarshal([]byte(inBytes), &obj)
	require.NoError(t, err)

	_, err = IDFromDID(*obj.Obj)
	require.EqualError(t, err, "invalid did format: blockchain mismatch: "+
		"found polygon in ID but eth in DID")
}

func TestDIDGenesisFromState(t *testing.T) {

	typ0, err := BuildDIDType(DIDMethodIden3, ReadOnly, NoNetwork)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	require.Equal(t, string(DIDMethodIden3), did.Method)
	blockchain, networkID, _, err := Decompose(*did)
	require.NoError(t, err)
	require.Equal(t, ReadOnly, blockchain)
	require.Equal(t, NoNetwork, networkID)
	require.Equal(t,
		"did:iden3:readonly:tJ93RwaVfE1PEMxd5rpZZuPtLCwbEaDCrNBhAy8HM",
		did.String())
}

func TestDID_PolygonID_Types(t *testing.T) {

	// Polygon no chain, no network
	did := helperBuildDIDFromType(t, DIDMethodPolygonID, ReadOnly, NoNetwork)

	require.Equal(t, string(DIDMethodPolygonID), did.Method)
	blockchain, networkID, _, err := Decompose(*did)
	require.NoError(t, err)
	require.Equal(t, ReadOnly, blockchain)
	require.Equal(t, NoNetwork, networkID)
	require.Equal(t,
		"did:polygonid:readonly:2mbH5rt9zKT1mTivFAie88onmfQtBU9RQhjNPLwFZh",
		did.String())

	// Polygon | Polygon chain, Main
	did2 := helperBuildDIDFromType(t, DIDMethodPolygonID, Polygon, Main)

	require.Equal(t, string(DIDMethodPolygonID), did2.Method)
	blockchain2, networkID2, _, err := Decompose(*did2)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain2)
	require.Equal(t, Main, networkID2)
	require.Equal(t,
		"did:polygonid:polygon:main:2pzr1wiBm3Qhtq137NNPPDFvdk5xwRsjDFnMxpnYHm",
		did2.String())

	// Polygon | Polygon chain, Mumbai
	did3 := helperBuildDIDFromType(t, DIDMethodPolygonID, Polygon, Mumbai)

	require.Equal(t, string(DIDMethodPolygonID), did3.Method)
	blockchain3, networkID3, _, err := Decompose(*did3)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain3)
	require.Equal(t, Mumbai, networkID3)
	require.Equal(t,
		"did:polygonid:polygon:mumbai:2qCU58EJgrELNZCDkSU23dQHZsBgAFWLNpNezo1g6b",
		did3.String())

}

func TestDID_PolygonID_ParseDIDFromID_OnChain(t *testing.T) {
	id1, err := IDFromString("2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6")
	require.NoError(t, err)

	did1, err := ParseDIDFromID(id1)
	require.NoError(t, err)

	var addressBytesExp [20]byte
	_, err = hex.Decode(addressBytesExp[:],
		[]byte("A51c1fc2f0D1a1b8494Ed1FE312d7C3a78Ed91C0"))
	require.NoError(t, err)

	require.Equal(t, string(DIDMethodPolygonID), did1.Method)
	wantIDs := []string{"polygon", "mumbai",
		"2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6"}
	require.Equal(t, wantIDs, did1.IDStrings)
	bc, nID, id, err := Decompose(*did1)
	require.NoError(t, err)
	require.Equal(t, Polygon, bc)
	require.Equal(t, Mumbai, nID)
	require.Equal(t, true, id.IsOnChain())

	addressBytes, err := id.EthAddress()
	require.NoError(t, err)
	require.Equal(t, addressBytesExp, addressBytes)

	require.Equal(t,
		"did:polygonid:polygon:mumbai:2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6",
		did1.String())
}

func TestDecompose(t *testing.T) {
	s := "did:polygonid:polygon:mumbai:2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6"

	did3, err := Parse(s)
	require.NoError(t, err)

	wantID, err := IDFromString("2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6")
	require.NoError(t, err)

	bch, nt, id, err := Decompose(*did3)
	require.NoError(t, err)
	require.Equal(t, Polygon, bch)
	require.Equal(t, Mumbai, nt)
	require.Equal(t, wantID, id)
}

func helperBuildDIDFromType(t testing.TB, method DIDMethod,
	blockchain Blockchain, network NetworkID) *DID {
	t.Helper()

	typ, err := BuildDIDType(method, blockchain, network)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ, genesisState)
	require.NoError(t, err)

	return did
}
