package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/iden3/go-iden3-core/v2/w3c"
	"github.com/stretchr/testify/require"
)

func TestParseDID(t *testing.T) {

	// did
	didStr := "did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"

	did, err := w3c.ParseDID(didStr)
	require.NoError(t, err)

	id, err := IDFromDID(*did)
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

	did, err = w3c.ParseDID(didStr)
	require.NoError(t, err)

	id, err = IDFromDID(*did)
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
		Obj *w3c.DID `json:"obj"`
	}
	err = json.Unmarshal([]byte(inBytes), &obj)
	require.NoError(t, err)
	require.NotNil(t, obj.Obj)
	require.Equal(t, string(DIDMethodIden3), obj.Obj.Method)

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
		Obj *w3c.DID `json:"obj"`
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
	did, err := NewDIDFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	require.Equal(t, string(DIDMethodIden3), did.Method)

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

func TestDIDFromID(t *testing.T) {
	typ0, err := BuildDIDType(DIDMethodIden3, ReadOnly, NoNetwork)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	id, err := NewIDFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	did, err := ParseDIDFromID(*id)
	require.NoError(t, err)

	require.Equal(t,
		"did:iden3:readonly:tJ93RwaVfE1PEMxd5rpZZuPtLCwbEaDCrNBhAy8HM",
		did.String())
}

func TestDID_PolygonID_Types(t *testing.T) {
	testCases := []struct {
		title   string
		method  DIDMethod
		chain   Blockchain
		net     NetworkID
		wantDID string
	}{
		{
			title:   "Polygon no chain, no network",
			method:  DIDMethodPolygonID,
			chain:   ReadOnly,
			net:     NoNetwork,
			wantDID: "did:polygonid:readonly:2mbH5rt9zKT1mTivFAie88onmfQtBU9RQhjNPLwFZh",
		},
		{
			title:   "Polygon | Polygon chain, Main",
			method:  DIDMethodPolygonID,
			chain:   Polygon,
			net:     Main,
			wantDID: "did:polygonid:polygon:main:2pzr1wiBm3Qhtq137NNPPDFvdk5xwRsjDFnMxpnYHm",
		},
		{
			title:   "Polygon | Polygon chain, Mumbai",
			method:  DIDMethodPolygonID,
			chain:   Polygon,
			net:     Mumbai,
			wantDID: "did:polygonid:polygon:mumbai:2qCU58EJgrELNZCDkSU23dQHZsBgAFWLNpNezo1g6b",
		},
		{
			title:   "Polygon | Polygon chain, Amoy",
			method:  DIDMethodPolygonID,
			chain:   Polygon,
			net:     Amoy,
			wantDID: "did:polygonid:polygon:amoy:2qQ68JkRcf3xrHPQPWZei3YeVzHPP58wYNxx2mEouR",
		},
		{
			title:   "Polygon | zkEVM chain, main",
			method:  DIDMethodPolygonID,
			chain:   ZkEVM,
			net:     Main,
			wantDID: "did:polygonid:zkevm:main:2wQjmkL1SsgqC7AuZdUcaXsUVfEi1i58VEhm3r2r8F",
		},
		{
			title:   "Polygon | zkEVM chain, test",
			method:  DIDMethodPolygonID,
			chain:   ZkEVM,
			net:     Test,
			wantDID: "did:polygonid:zkevm:test:2wcMpvr8NgWTfqN6ChaFEx1qRnLREXhjeoJ45pFyw5",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			did := helperBuildDIDFromType(t, tc.method, tc.chain, tc.net)
			require.Equal(t, string(tc.method), did.Method)
			id, err := IDFromDID(*did)
			require.NoError(t, err)
			method, err := MethodFromID(id)
			require.NoError(t, err)
			require.Equal(t, tc.method, method)
			blockchain, err := BlockchainFromID(id)
			require.NoError(t, err)
			require.Equal(t, tc.chain, blockchain)
			networkID, err := NetworkIDFromID(id)
			require.NoError(t, err)
			require.Equal(t, tc.net, networkID)
			require.Equal(t, tc.wantDID, did.String())
		})
	}
}

func TestDID_PolygonID_ParseDIDFromID(t *testing.T) {
	id1, err := IDFromString("2qCU58EJgrEM9NKvHkvg5NFWUiJPgN3M3LnCr98j3x")
	require.NoError(t, err)

	did, err := ParseDIDFromID(id1)
	require.NoError(t, err)

	var addressBytesExp [20]byte
	_, err = hex.Decode(addressBytesExp[:],
		[]byte("A51c1fc2f0D1a1b8494Ed1FE312d7C3a78Ed91C0"))
	require.NoError(t, err)

	require.Equal(t, string(DIDMethodPolygonID), did.Method)
	wantIDs := []string{"polygon", "mumbai",
		"2qCU58EJgrEM9NKvHkvg5NFWUiJPgN3M3LnCr98j3x"}
	require.Equal(t, wantIDs, did.IDStrings)
	id, err := IDFromDID(*did)
	require.NoError(t, err)
	method, err := MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethodPolygonID, method)
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
		"did:polygonid:polygon:mumbai:2qCU58EJgrEM9NKvHkvg5NFWUiJPgN3M3LnCr98j3x",
		did.String())
}

func TestDecompose(t *testing.T) {
	wantIDHex := "2qCU58EJgrEM9NKvHkvg5NFWUiJPgN3M3LnCr98j3x"
	ethAddrHex := "a51c1fc2f0d1a1b8494ed1fe312d7c3a78ed91c0"
	genesis := genFromHex("00000000000000" + ethAddrHex)
	tp, err := BuildDIDType(DIDMethodPolygonID, Polygon, Mumbai)
	require.NoError(t, err)
	id0 := NewID(tp, genesis)

	s := fmt.Sprintf("did:polygonid:polygon:mumbai:%v", id0.String())

	did, err := w3c.ParseDID(s)
	require.NoError(t, err)

	wantID, err := IDFromString(wantIDHex)
	require.NoError(t, err)

	id, err := IDFromDID(*did)
	require.NoError(t, err)
	require.Equal(t, wantID, id)

	method, err := MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethodPolygonID, method)

	blockchain, err := BlockchainFromID(id)
	require.NoError(t, err)
	require.Equal(t, Polygon, blockchain)

	networkID, err := NetworkIDFromID(id)
	require.NoError(t, err)
	require.Equal(t, Mumbai, networkID)

	ethAddr, err := EthAddressFromID(id)
	require.NoError(t, err)
	require.Equal(t, ethAddrFromHex(ethAddrHex), ethAddr)
}

func helperBuildDIDFromType(t testing.TB, method DIDMethod,
	blockchain Blockchain, network NetworkID) *w3c.DID {
	t.Helper()

	typ, err := BuildDIDType(method, blockchain, network)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := NewDIDFromIdenState(typ, genesisState)
	require.NoError(t, err)

	return did
}

func TestNewIDFromDID(t *testing.T) {
	did, err := w3c.ParseDID("did:something:x")
	require.NoError(t, err)
	id := newIDFromUnsupportedDID(*did)
	require.Equal(t, []byte{0xff, 0xff}, id[:2])
	wantID, err := hex.DecodeString(
		"ffff84b1e6d0d9ecbe951348ea578dbacc022cdbbff4b11218671dca871c11")
	require.NoError(t, err)
	require.Equal(t, wantID, id[:])

	id2, err := IDFromDID(*did)
	require.NoError(t, err)
	require.Equal(t, id, id2)
}

func TestGenesisFromEthAddress(t *testing.T) {

	ethAddrHex := "accb91a7d1d9ad0d33b83f2546ed30285c836c6e"
	wantGenesisHex := "00000000000000accb91a7d1d9ad0d33b83f2546ed30285c836c6e"
	require.Len(t, ethAddrHex, 20*2)
	require.Len(t, wantGenesisHex, 27*2)

	ethAddrBytes, err := hex.DecodeString(ethAddrHex)
	require.NoError(t, err)
	var ethAddr [20]byte
	copy(ethAddr[:], ethAddrBytes)

	genesis := GenesisFromEthAddress(ethAddr)
	wantGenesis, err := hex.DecodeString(wantGenesisHex)
	require.NoError(t, err)
	require.Equal(t, wantGenesis, genesis[:])

	tp2, err := BuildDIDType(DIDMethodPolygonID, Polygon, Mumbai)
	require.NoError(t, err)

	id := NewID(tp2, genesis)
	ethAddr2, err := EthAddressFromID(id)
	require.NoError(t, err)
	require.Equal(t, ethAddr, ethAddr2)

	var wantID ID
	copy(wantID[:], tp2[:])
	copy(wantID[len(tp2):], genesis[:])
	ch := CalculateChecksum(tp2, genesis)
	copy(wantID[len(tp2)+len(genesis):], ch[:])
	require.Equal(t, wantID, id)

	// make genesis not look like an address
	genesis[0] = 1
	id = NewID(tp2, genesis)
	_, err = EthAddressFromID(id)
	require.EqualError(t, err,
		"can't get Ethereum address: high bytes of genesis are not zero")
}

func genFromHex(gh string) [genesisLn]byte {
	genBytes, err := hex.DecodeString(gh)
	if err != nil {
		panic(err)
	}
	var gen [genesisLn]byte
	copy(gen[:], genBytes)
	return gen
}

func ethAddrFromHex(ea string) [20]byte {
	eaBytes, err := hex.DecodeString(ea)
	if err != nil {
		panic(err)
	}
	var ethAddr [20]byte
	copy(ethAddr[:], eaBytes)
	return ethAddr
}

func TestCustomDIDRegistration(t *testing.T) {
	testCases := []struct {
		Description string
		Data        DIDMethodNetworkParams
		opts        []RegistrationOptions
	}{
		{
			Description: "register new did method network",
			Data: DIDMethodNetworkParams{
				Method:      "test_method",
				Blockchain:  "test_chain",
				Network:     "test_net",
				NetworkFlag: 0b0001_0001,
			},
			opts: []RegistrationOptions{WithChainID(101), WithDIDMethodByte(0b00000011)},
		},
		{
			Description: "register one more new did method network",
			Data: DIDMethodNetworkParams{
				Method:      "method",
				Blockchain:  "chain",
				Network:     "network",
				NetworkFlag: 0b0001_0001,
			},
			opts: []RegistrationOptions{WithChainID(102), WithDIDMethodByte(0b00000100)},
		},
		{
			Description: "register the same new did method network",
			Data: DIDMethodNetworkParams{
				Method:      "method",
				Blockchain:  "chain",
				Network:     "network",
				NetworkFlag: 0b0001_0001,
			},
			opts: []RegistrationOptions{WithChainID(102), WithDIDMethodByte(0b00000100)},
		},
		{
			Description: "register network to existing did method",
			Data: DIDMethodNetworkParams{
				Method:      DIDMethodIden3,
				Blockchain:  "chain",
				Network:     Test,
				NetworkFlag: 0b01000000 | 0b00000011,
			},
			opts: []RegistrationOptions{WithChainID(103)},
		},
		{
			Description: "register network to existing did method and chainId",
			Data: DIDMethodNetworkParams{
				Method:      DIDMethodIden3,
				Blockchain:  ReadOnly,
				Network:     NoNetwork,
				NetworkFlag: 0b00000000,
			},
			opts: []RegistrationOptions{WithChainID(104)},
		},
		{
			Description: "register one more network to existing did method",
			Data: DIDMethodNetworkParams{
				Method:      DIDMethodIden3,
				Blockchain:  ReadOnly,
				Network:     "network",
				NetworkFlag: 0b11000000 | 0b00000011,
			},
			opts: []RegistrationOptions{WithChainID(105)},
		},
		{
			Description: "register known chain id to new did method",
			Data: DIDMethodNetworkParams{
				Method:      "method2",
				Blockchain:  Polygon,
				Network:     Mumbai,
				NetworkFlag: 0b0001_0001,
			},
			opts: []RegistrationOptions{WithDIDMethodByte(0b0000111)},
		},
		{
			Description: "register known chain id to new did method",
			Data: DIDMethodNetworkParams{
				Method:      "iden3",
				Blockchain:  ReadOnly,
				Network:     NoNetwork,
				NetworkFlag: 0b0000_0000,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			err := RegisterDIDMethodNetwork(tc.Data, tc.opts...)
			require.NoError(t, err)
		})
	}

	d := helperBuildDIDFromType(t, "method", "chain", "network")
	require.Equal(t, "4bb86obLkMrifHixMY62WM4iQQVr7u29cxWjMAinrT", d.IDStrings[2])

	did3, err := w3c.ParseDID("did:method:chain:network:4bb86obLkMrifHixMY62WM4iQQVr7u29cxWjMAinrT")
	require.NoError(t, err)

	id, err := idFromDID(*did3)
	require.NoError(t, err)

	require.Equal(t, "4bb86obLkMrifHixMY62WM4iQQVr7u29cxWjMAinrT", id.String())
	method, err := MethodFromID(id)
	require.NoError(t, err)
	require.Equal(t, DIDMethod("method"), method)

	blockchain, err := BlockchainFromID(id)
	require.NoError(t, err)
	require.Equal(t, Blockchain("chain"), blockchain)

	networkID, err := NetworkIDFromID(id)
	require.NoError(t, err)
	require.Equal(t, NetworkID("network"), networkID)

}

func TestCustomDIDRegistration_Negative(t *testing.T) {
	testCases := []struct {
		Description string
		Data        DIDMethodNetworkParams
		opts        []RegistrationOptions
		err         string
	}{

		{
			Description: "try to overwrite existing chain id",
			Data: DIDMethodNetworkParams{
				Method:      DIDMethodIden3,
				Blockchain:  Polygon,
				Network:     Mumbai,
				NetworkFlag: 0b0001_0001,
			},
			opts: []RegistrationOptions{WithChainID(1)},
			err:  "can't register chain id 1 for 'polygon:mumbai' because it's already registered for another chain id",
		},
		{
			Description: "try to overwrite existing DID method byte",
			Data: DIDMethodNetworkParams{
				Method:      DIDMethodIden3,
				Blockchain:  Ethereum,
				Network:     Main,
				NetworkFlag: 0b00100000 | 0b00000001,
			},
			opts: []RegistrationOptions{WithChainID(1), WithDIDMethodByte(0b00000010)},
			err:  "can't register method 'iden3' because DID method byte '10' already registered for another method",
		},
		{
			Description: "try to write max did method byte",
			Data: DIDMethodNetworkParams{
				Method:      "method33",
				Blockchain:  Ethereum,
				Network:     Main,
				NetworkFlag: 0b00100000 | 0b00000001,
			},
			opts: []RegistrationOptions{WithChainID(1), WithDIDMethodByte(0b11111111)},
			err:  "Can't register DID method byte: current 11111111, maximum byte allowed: 11111110",
		},
		{
			Description: "try to rewrite existing DID Method Network Flag",
			Data: DIDMethodNetworkParams{
				Method:      DIDMethodIden3,
				Blockchain:  Ethereum,
				Network:     Main,
				NetworkFlag: 0b00100000 | 0b00000011,
			},
			opts: nil,
			err:  "DID network flag 100011 is already registered for the another network id for 'iden3' method",
		},
		{
			Description: "register new did method with existing method byte",
			Data: DIDMethodNetworkParams{
				Method:      "new_method",
				Blockchain:  "new_chain",
				Network:     "new_net",
				NetworkFlag: 0b0001_0001,
			},
			opts: []RegistrationOptions{WithChainID(101), WithDIDMethodByte(0b00000001)},
			err:  "can't register method 'new_method' because DID method byte '1' already registered for another method",
		},
		{
			Description: "register new did method with existing chain id",
			Data: DIDMethodNetworkParams{
				Method:      "new_method",
				Blockchain:  Ethereum,
				Network:     Main,
				NetworkFlag: 0b0001_0001,
			},
			opts: []RegistrationOptions{WithChainID(137), WithDIDMethodByte(0b10000000)},
			err:  "can't register chain id 137 for 'eth:main' because it's already registered for another chain id",
		},
		{
			Description: "register new network and chain with existing networkFlag for existing existing did method",
			Data: DIDMethodNetworkParams{
				Method:      DIDMethodIden3,
				Blockchain:  "supa_chain",
				Network:     "supa_net",
				NetworkFlag: 0b00010000 | 0b00000001,
			},
			err: "DID network flag 10001 is already registered for the another network id for 'iden3' method",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			err := RegisterDIDMethodNetwork(tc.Data, tc.opts...)
			require.EqualError(t, err, tc.err)
		})
	}

}
