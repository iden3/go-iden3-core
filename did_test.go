package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DIDString(t *testing.T) {

	tests := []struct {
		description string
		identifier  string
		did         string
		options     DIDOption
	}{
		{"Test readonly did",
			"114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE",
			"did:iden3:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE",
			nil,
		},
		{"Test eth did",
			"114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE",
			"did:iden3:eth:main:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE",
			WithNetwork("eth", "main"),
		},
		{"Test polygon did",
			"114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE",
			"did:iden3:polygon:test:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE",
			WithNetwork("polygon", "test"),
		},
		{"Test identifier 41 char",
			"11FjRaFUGZA5yBXREaH6P11yezYsxwJLMsEUerTVj",
			"did:iden3:polygon:test:11FjRaFUGZA5yBXREaH6P11yezYsxwJLMsEUerTVj",
			WithNetwork("polygon", "test"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got, err := New(test.identifier, test.options)
			assert.NoError(t, err)

			assert.Equal(t, got.String(), test.did)
		})
	}

}

func TestParseDID(t *testing.T) {

	didStr := "did:iden3:eth:test:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE"

	did, err := ParseDID(didStr)
	assert.NoError(t, err)

	assert.Equal(t, "114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE", did.ID.String())
	assert.Equal(t, Network("test"), did.Network)
	assert.Equal(t, NetworkID("eth"), did.NetworkID)
}

func TestDID_ParseDID_DoesntMatchRegexp(t *testing.T) {
	didStr := "dididen3:eth:test:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE"

	_, err := ParseDID(didStr)
	assert.ErrorIs(t, err, ErrDoesnotMatchRegexp)
}
