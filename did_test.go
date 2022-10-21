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

	didStr := "did:iden3:eth:test:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E"

	did, err := ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E",
		did.ID.String())
	require.Equal(t, NetworkID("test"), did.NetworkID)
	require.Equal(t, Blockchain("eth"), did.Blockchain)
}

func TestDID_ParseDID_DoesntMatchRegexp(t *testing.T) {
	didStr := "dididen3:eth:test:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4m2E"

	_, err := ParseDID(didStr)
	require.ErrorIs(t, err, ErrDoesNotMatchRegexp)
}
