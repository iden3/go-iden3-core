package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// DIDMethod DID method-name
const DIDMethod = "iden3"

// DIDSchema DID Schema
const DIDSchema = "did"

// NetworkID id of the network "eth", "poligon", etc.
type NetworkID string

const (
	ETHEREUM NetworkID = "eth"     // ETHEREUM ethereum network
	POLIGON  NetworkID = "poligon" // POLIGON poligon network
)

type Network string

const (
	MAIN    Network = "main"    // main net
	TEST    Network = "test"    // test net
	ROPSTEN Network = "ropsten" // ropsten net
	RINKEBY Network = "rinkeby" // rinkeby net
	KOVAN   Network = "kovan"   // kovan net
)

var (
	// valid id for regexp
	// did:iden3:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE -readonly id. For readonly identifier networkID and
	// network can be empty as this identifier is newer published on chain
	// did:iden3:eth:main:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE - eth network eth networkID, main - network

	didRegex = regexp.MustCompile(`^\b(did):\b(iden3):\b(eth|poligon):\b(main|test|ropsten|rinkeby|kovan):(
[1-9a-km-zA-HJ-NP-Z]{41,42})$`)
	DoesntMatchRegexp = errors.New("did doesnt matching regex")
)

// DID Decentralized Identifiers (DIDs)
// https://w3c.github.io/did-core/#did-syntax
type DID struct {
	ID        ID        // ID did specific id
	NetworkID NetworkID // NetworkID network identifier eth / poligon,...
	Network   Network   // Network specific network identifier eth {main, ropsten, rinkeby, kovan}
}

type DIDOption func(*DID) error

func New(didStr string, options ...DIDOption) (*DID, error) {

	did := &DID{}
	var err error

	did.ID, err = IDFromString(didStr)
	if err != nil {
		return &DID{}, err
	}

	for _, o := range options {
		if o == nil {
			continue
		}
		err := o(did)
		if err != nil {
			return nil, err
		}
	}
	return did, nil

}

// WithNetwork sets NetworkID and Network (eth:main)
func WithNetwork(networkID, network string) DIDOption {
	return func(d *DID) error {
		d.NetworkID = NetworkID(networkID)
		d.Network = Network(network)
		return nil
	}
}

// String did as a string
func (did *DID) String() string {
	if did.NetworkID == "" {
		return fmt.Sprintf("%s:%s:%s", DIDSchema, DIDMethod, did.ID.String())
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s", DIDSchema, DIDMethod, did.NetworkID, did.Network, did.ID.String())
}

// ParseDID method parse string and extract DID if string is valid Iden3 identifier
func ParseDID(didStr string) (*DID, error) {
	did := DID{}
	emptyDID := &DID{}
	var err error

	matched := didRegex.MatchString(didStr)
	if !matched {
		return emptyDID, DoesntMatchRegexp
	}

	arg := strings.Split(didStr, ":")

	// validate id
	did.ID, err = IDFromString(arg[4])
	if err != nil {
		return emptyDID, err
	}

	did.Network = Network(arg[3])
	did.NetworkID = NetworkID(arg[2])

	return &did, nil
}
