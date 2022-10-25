package core

import (
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
			"114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
			"did:iden3:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
			nil,
		},
		{"Test eth did",
			"114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
			"did:iden3:eth:main:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
			WithNetwork("eth", "main"),
		},
		{"Test polygon did",
			"114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
			"did:iden3:polygon:test:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
			WithNetwork("polygon", "test"),
		},
		{"Test identifier 41 char",
			"11FjRaFUGZA5yBXREaH6P11yezYsxwJLMsEUerXWk",
			"did:iden3:polygon:test:11FjRaFUGZA5yBXREaH6P11yezYsxwJLMsEUerXWk",
			WithNetwork("polygon", "test"),
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
	didStr := "did:iden3:polygon:mumbai:11FjRaFUGZA5yBXREaH6P11yezYsxwJLMsEUerXWk"

	did, err := ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "11FjRaFUGZA5yBXREaH6P11yezYsxwJLMsEUerXWk",
		did.ID.String())
	require.Equal(t, NetworkID("mumbai"), did.NetworkID)
	require.Equal(t, Blockchain("polygon"), did.Blockchain)

	// readonly did
	didStr = "did:iden3:1MWtoAdZESeiphxp3bXupZcfS9DhMTdWNSjRwVYc2"

	did, err = ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "1MWtoAdZESeiphxp3bXupZcfS9DhMTdWNSjRwVYc2", did.ID.String())
	require.Equal(t, NetworkID(""), did.NetworkID)
	require.Equal(t, Blockchain(""), did.Blockchain)
	require.Equal(t, TypeReadOnly, did.ID.Type())
}

func TestDID_ParseDID_DoesntMatchRegexp(t *testing.T) {
	didStr := "did:iden3:ethereum:ropsten:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E"

	_, err := ParseDID(didStr)
	require.ErrorIs(t, err, ErrDoesNotMatchRegexp)
}
