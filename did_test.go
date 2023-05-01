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

	id, err := IDFromDID(*did3)
	require.NoError(t, err)
	require.Equal(t, "wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ", id.String())
	method, err := MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethodIden3, method)
	blockchain, err := BlockchainFromID(id)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain)
	networkID, err := NetworkIDFromID(id)
	require.NoError(t, err)
	require.Equal(t, Mumbai, networkID)

	// readonly did
	didStr = "did:iden3:readonly:tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa"

	did3, err = Parse(didStr)
	require.NoError(t, err)

	id, err = IDFromDID(*did3)
	require.NoError(t, err)
	require.Equal(t, "tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa", id.String())
	method, err = MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethodIden3, method)
	blockchain, err = BlockchainFromID(id)
	require.NoError(t, err)
	require.Equal(t, ReadOnly, blockchain)
	networkID, err = NetworkIDFromID(id)
	require.NoError(t, err)
	require.Equal(t, NoNetwork, networkID)

	require.Equal(t, [2]byte{DIDMethodIden3.Byte(), 0b0}, id.Type())
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
	require.Equal(t, DIDMethodIden3.String(), obj.Obj.Method)

	id2, err := IDFromDID(*obj.Obj)
	require.NoError(t, err)
	method, err := MethodFromID(id2)
	require.NoError(t, err)
	require.Equal(t, DIDMethodIden3, method)
	blockchain, err := BlockchainFromID(id2)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain)
	networkID, err := NetworkIDFromID(id2)
	require.NoError(t, err)
	require.Equal(t, Mumbai, networkID)

	require.Equal(t, id, id2)
}

func TestDID_UnmarshalJSON_Error(t *testing.T) {
	inBytes := `{"obj": "did:iden3:eth:goerli:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"}`
	var obj struct {
		Obj *DID `json:"obj"`
	}
	err := json.Unmarshal([]byte(inBytes), &obj)
	require.NoError(t, err)

	//_, err = IDFromDID(*obj.Obj)
	//require.EqualError(t, err, "invalid did format: blockchain mismatch: "+
	//	"found polygon in ID but eth in DID")
}

func TestDIDGenesisFromState(t *testing.T) {

	typ0, err := BuildDIDType(DIDMethodIden3, ReadOnly, NoNetwork)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	require.Equal(t, DIDMethodIden3.String(), did.Method)

	id, err := IDFromDID(*did)
	require.NoError(t, err)
	method, err := MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethodIden3, method)
	blockchain, err := BlockchainFromID(id)
	require.NoError(t, err)
	require.Equal(t, ReadOnly, blockchain)
	networkID, err := NetworkIDFromID(id)
	require.NoError(t, err)
	require.Equal(t, NoNetwork, networkID)

	require.Equal(t,
		"did:iden3:readonly:tJ93RwaVfE1PEMxd5rpZZuPtLCwbEaDCrNBhAy8HM",
		did.String())
}

func TestDID_PolygonID_Types(t *testing.T) {

	// Polygon no chain, no network
	did := helperBuildDIDFromType(t, DIDMethodPolygonID, ReadOnly, NoNetwork)

	require.Equal(t, DIDMethodPolygonID.String(), did.Method)
	id, err := IDFromDID(*did)
	require.NoError(t, err)
	method, err := MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethodPolygonID, method)
	blockchain, err := BlockchainFromID(id)
	require.NoError(t, err)
	require.Equal(t, ReadOnly, blockchain)
	networkID, err := NetworkIDFromID(id)
	require.NoError(t, err)
	require.Equal(t, NoNetwork, networkID)
	require.Equal(t,
		"did:polygonid:readonly:2mbH5rt9zKT1mTivFAie88onmfQtBU9RQhjNPLwFZh",
		did.String())

	// Polygon | Polygon chain, Main
	did2 := helperBuildDIDFromType(t, DIDMethodPolygonID, Polygon, Main)

	require.Equal(t, DIDMethodPolygonID.String(), did2.Method)
	id2, err := IDFromDID(*did2)
	require.NoError(t, err)
	method2, err := MethodFromID(id2)
	require.NoError(t, err)
	require.Equal(t, DIDMethodPolygonID, method2)
	blockchain2, err := BlockchainFromID(id2)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain2)
	networkID2, err := NetworkIDFromID(id2)
	require.NoError(t, err)
	require.Equal(t, Main, networkID2)
	require.Equal(t,
		"did:polygonid:polygon:main:2pzr1wiBm3Qhtq137NNPPDFvdk5xwRsjDFnMxpnYHm",
		did2.String())

	// Polygon | Polygon chain, Mumbai
	did3 := helperBuildDIDFromType(t, DIDMethodPolygonID, Polygon, Mumbai)

	require.Equal(t, DIDMethodPolygonID.String(), did3.Method)
	id3, err := IDFromDID(*did3)
	require.NoError(t, err)
	method3, err := MethodFromID(id3)
	require.NoError(t, err)
	require.Equal(t, DIDMethodPolygonID, method3)
	blockchain3, err := BlockchainFromID(id3)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain3)
	networkID3, err := NetworkIDFromID(id3)
	require.NoError(t, err)
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

	require.Equal(t, DIDMethodPolygonIDOnChain.String(), did1.Method)
	wantIDs := []string{"polygon", "mumbai",
		"2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6"}
	require.Equal(t, wantIDs, did1.IDStrings)
	id, err := IDFromDID(*did1)
	require.NoError(t, err)
	method, err := MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethodPolygonIDOnChain, method)
	blockchain, err := BlockchainFromID(id)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain)
	networkID, err := NetworkIDFromID(id)
	require.NoError(t, err)
	require.Equal(t, Mumbai, networkID)

	ethAddr, err := EthAddressFromID(id)
	require.NoError(t, err)
	require.Equal(t, addressBytesExp, ethAddr)

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

	id, err := IDFromDID(*did3)
	require.NoError(t, err)
	require.Equal(t, wantID, id)

	method, err := MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethodPolygonIDOnChain, method)

	blockchain, err := BlockchainFromID(id)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain)

	networkID, err := NetworkIDFromID(id)
	require.NoError(t, err)
	require.Equal(t, Mumbai, networkID)
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

func TestNewIDFromDID(t *testing.T) {
	did, err := Parse("did:something:x")
	require.NoError(t, err)
	id := newIDFromDID(*did)
	wantID, err := hex.DecodeString(
		"ffff84b1e6d0d9ecbe951348ea578dbacc022cdbbff4b11218671dca871c11")
	require.NoError(t, err)
	require.Equal(t, wantID, id[:])

	id2, err := IDFromDID(*did)
	require.NoError(t, err)
	require.Equal(t, id, id2)
}
