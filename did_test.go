package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DIDString(t *testing.T) {
	tests := []struct {
		description string
		identifier  string
		did         string
		options     DIDOption
	}{
		{"Test readonly did",
			"tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa",
			"did:iden3:tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa",
			nil,
		},
		{"Test eth did",
			"zyaYCrj27j7gJfrBboMW49HFRSkQznyy12ABSVzTy",
			"did:iden3:eth:main:zyaYCrj27j7gJfrBboMW49HFRSkQznyy12ABSVzTy",
			WithNetwork("eth", "main"),
		},
		{"Test polygon did",
			"wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ",
			"did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ",
			WithNetwork("polygon", "mumbai"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got, err := NewDID(test.identifier, test.options)
			require.NoError(t, err)
			require.Equal(t, got.String(), test.did)
		})
	}
}

func TestParseDID(t *testing.T) {

	// did
	didStr := "did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"

	did, err := ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ",
		did.ID.String())
	require.Equal(t, Mumbai, did.NetworkID)
	require.Equal(t, Polygon, did.Blockchain)

	// readonly did
	didStr = "did:iden3:tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa"

	did, err = ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa", did.ID.String())
	require.Equal(t, NetworkID(""), did.NetworkID)
	require.Equal(t, Blockchain(""), did.Blockchain)

	require.Equal(t, [2]byte{DIDMethodByte[DIDMethodIden3], 0b0}, did.ID.Type())
}

func TestDIDGenesisFromState(t *testing.T) {

	typ0, err := BuildDIDType(DIDMethodIden3, NoChain, NoNetwork)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	require.Equal(t, DIDMethodIden3, did.Method)
	require.Equal(t, NoChain, did.Blockchain)
	require.Equal(t, NoNetwork, did.NetworkID)
	require.Equal(t, "did:iden3:tJ93RwaVfE1PEMxd5rpZZuPtLCwbEaDCrNBhAy8HM", did.String())
}
