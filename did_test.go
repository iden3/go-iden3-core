package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/stretchr/testify/require"
)

func TestParseDID(t *testing.T) {

	// did
	didStr := "did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"

	did, err := ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ",
		did.ID.String())
	require.Equal(t, Polygon, did.Blockchain)
	require.Equal(t, Mumbai, did.NetworkID)

	// readonly did
	didStr = "did:iden3:readonly:tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa"

	did, err = ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa",
		did.ID.String())
	require.Equal(t, ReadOnly, did.Blockchain)
	require.Equal(t, NoNetwork, did.NetworkID)

	require.Equal(t, [2]byte{DIDMethodByte[DIDMethodIden3], 0b0}, did.ID.Type())
}

func TestDID_MarshalJSON(t *testing.T) {
	id, err := IDFromString("wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ")
	require.NoError(t, err)
	did := DID{
		ID:         id,
		Method:     DIDMethodIden3,
		Blockchain: Polygon,
		NetworkID:  Mumbai,
	}

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
	require.Equal(t, id, obj.Obj.ID)
	require.Equal(t, DIDMethodIden3, obj.Obj.Method)
	require.Equal(t, Polygon, obj.Obj.Blockchain)
	require.Equal(t, Mumbai, obj.Obj.NetworkID)
}

func TestDID_UnmarshalJSON_Error(t *testing.T) {
	inBytes := `{"obj": "did:iden3:eth:goerli:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"}`
	var obj struct {
		Obj *DID `json:"obj"`
	}
	err := json.Unmarshal([]byte(inBytes), &obj)
	require.EqualError(t, err,
		"invalid did format: network method of core identity mumbai differs from given did network specific id goerli")
}

func TestDIDGenesisFromState(t *testing.T) {

	typ0, err := BuildDIDType(DIDMethodIden3, ReadOnly, NoNetwork)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	require.Equal(t, DIDMethodIden3, did.Method)
	require.Equal(t, ReadOnly, did.Blockchain)
	require.Equal(t, NoNetwork, did.NetworkID)
	require.Equal(t, "did:iden3:readonly:tJ93RwaVfE1PEMxd5rpZZuPtLCwbEaDCrNBhAy8HM", did.String())
}

func TestDID_PolygonID_Types(t *testing.T) {

	// Polygon no chain, no network
	did := helperBuildDIDFromType(t, DIDMethodPolygonID, ReadOnly, NoNetwork)

	require.Equal(t, DIDMethodPolygonID, did.Method)
	require.Equal(t, ReadOnly, did.Blockchain)
	require.Equal(t, NoNetwork, did.NetworkID)
	require.Equal(t, "did:polygonid:readonly:2mbH5rt9zKT1mTivFAie88onmfQtBU9RQhjNPLwFZh", did.String())

	// Polygon | Polygon chain, Main
	did2 := helperBuildDIDFromType(t, DIDMethodPolygonID, Polygon, Main)

	require.Equal(t, DIDMethodPolygonID, did2.Method)
	require.Equal(t, Polygon, did2.Blockchain)
	require.Equal(t, Main, did2.NetworkID)
	require.Equal(t, "did:polygonid:polygon:main:2pzr1wiBm3Qhtq137NNPPDFvdk5xwRsjDFnMxpnYHm", did2.String())

	// Polygon | Polygon chain, Mumbai
	did3 := helperBuildDIDFromType(t, DIDMethodPolygonID, Polygon, Mumbai)

	require.Equal(t, DIDMethodPolygonID, did3.Method)
	require.Equal(t, Polygon, did3.Blockchain)
	require.Equal(t, Mumbai, did3.NetworkID)
	require.Equal(t, "did:polygonid:polygon:mumbai:2qCU58EJgrELNZCDkSU23dQHZsBgAFWLNpNezo1g6b", did3.String())

}

func TestDID_PolygonID_Types_OnChain(t *testing.T) {
	// Polygon | Polygon chain, Mumbai
	did1 := helperBuildDIDFromTypeOnchain(t, DIDMethodPolygonID, Polygon, Mumbai)

	idInt, _ := big.NewInt(0).SetString("20318741244951419790279970260471061329920784847526059589266894371380597378", 10)
	idBytes := utils.SwapEndianness(idInt.Bytes())
	fmt.Printf("\ndid hex from solidity: %x\n\n", idBytes)

	idExp, err := IDFromInt(idInt)
	require.NoError(t, err)
	require.Equal(t, "2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6", idExp.String())

	var addressBytesExp [20]byte
	_, err = hex.Decode(addressBytesExp[:], []byte("A51c1fc2f0D1a1b8494Ed1FE312d7C3a78Ed91C0"))
	require.NoError(t, err)

	require.Equal(t, DIDMethodPolygonID, did1.Method)
	require.Equal(t, Polygon, did1.Blockchain)
	require.Equal(t, Mumbai, did1.NetworkID)
	require.Equal(t, true, did1.ID.IsOnChain())

	addressBytes, err := did1.ID.EthAddress()
	require.NoError(t, err)
	require.Equal(t, addressBytesExp, addressBytes)

	require.Equal(t, "did:polygonid:polygon:mumbai:2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6", did1.String())
}

func TestDID_PolygonID_ParseDIDFromID_OnChain(t *testing.T) {
	id1, err := IDFromString("2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6")
	require.NoError(t, err)

	did1, err := ParseDIDFromID(id1)
	require.NoError(t, err)

	var addressBytesExp [20]byte
	_, err = hex.Decode(addressBytesExp[:], []byte("A51c1fc2f0D1a1b8494Ed1FE312d7C3a78Ed91C0"))
	require.NoError(t, err)

	require.Equal(t, DIDMethodPolygonID, did1.Method)
	require.Equal(t, Polygon, did1.Blockchain)
	require.Equal(t, Mumbai, did1.NetworkID)
	require.Equal(t, true, did1.ID.IsOnChain())

	addressBytes, err := did1.ID.EthAddress()
	require.NoError(t, err)
	require.Equal(t, addressBytesExp, addressBytes)

	require.Equal(t, "did:polygonid:polygon:mumbai:2z39iB1bPjY2STTFSwbzvK8gqJQMsv5PLpvoSg3opa6", did1.String())
}

func helperBuildDIDFromType(t testing.TB, method DIDMethod, blockchain Blockchain, network NetworkID) *DID {
	t.Helper()

	typ, err := BuildDIDType(method, blockchain, network)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ, genesisState)
	require.NoError(t, err)

	return did
}

func helperBuildDIDFromTypeOnchain(t testing.TB, method DIDMethod, blockchain Blockchain, network NetworkID) *DID {
	t.Helper()

	typ, err := BuildDIDTypeOnChain(method, blockchain, network)
	require.NoError(t, err)

	var addressBytes [20]byte
	_, err = hex.Decode(addressBytes[:], []byte("A51c1fc2f0D1a1b8494Ed1FE312d7C3a78Ed91C0"))
	require.NoError(t, err)
	fmt.Printf("eth address: 0x%x\n", addressBytes)

	genesisState := GenesisFromEthAddress(addressBytes)

	did, err := DIDGenesisFromIdenState(typ, genesisState)
	require.NoError(t, err)

	fmt.Printf("did: %s\n", did.String())
	fmt.Printf("did hex: %x\n", did.ID.Bytes())

	return did
}
