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

// Blockchain id of the network "eth", "polygon", etc.
type Blockchain string

const (
	ETHEREUM Blockchain = "eth"     // ETHEREUM ethereum network
	POLYGON  Blockchain = "polygon" // POLYGON polygon network
)

type NetworkID string

const (
	MAIN    NetworkID = "main"    // main net
	TEST    NetworkID = "test"    // test net
	ROPSTEN NetworkID = "ropsten" // ropsten net
	RINKEBY NetworkID = "rinkeby" // rinkeby net
	KOVAN   NetworkID = "kovan"   // kovan net
)

var (
	// valid id for regexp
	// did:iden3:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE -readonly id. For readonly identifier networkID and
	// network can be empty as this identifier is newer published on chain
	// did:iden3:eth:main:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE - eth network eth networkID, main - network

	didRegex = regexp.MustCompile(`^\b(did):\b(iden3):(\b(eth|polygon):\b(main|test|ropsten|rinkeby|kovan):)?([1-9a-km-zA-HJ-NP-Z]{42}|[1-9a-km-zA-HJ-NP-Z]{41})$`)

	// ErrDoesNotMatchRegexp is returned when did string parsed
	ErrDoesNotMatchRegexp = errors.New("did does not match regex")
)

// DID Decentralized Identifiers (DIDs)
// https://w3c.github.io/did-core/#did-syntax
type DID struct {
	ID         ID         // ID did specific id
	Blockchain Blockchain // Blockchain network identifier eth / polygon,...
	NetworkID  NetworkID  // NetworkID specific network identifier eth {main, ropsten, rinkeby, kovan}
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

// WithNetwork sets Blockchain and NetworkID (eth:main)
func WithNetwork(blockchain, network string) DIDOption {
	return func(d *DID) error {
		d.Blockchain = Blockchain(blockchain)
		d.NetworkID = NetworkID(network)
		return nil
	}
}

// String did as a string
func (did *DID) String() string {
	if did.Blockchain == "" {
		return fmt.Sprintf("%s:%s:%s", DIDSchema, DIDMethod, did.ID.String())
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s", DIDSchema, DIDMethod, did.Blockchain, did.NetworkID, did.ID.String())
}

// ParseDID method parse string and extract DID if string is valid Iden3 identifier
func ParseDID(didStr string) (*DID, error) {
	did := DID{}
	emptyDID := &DID{}
	var err error

	matched := didRegex.MatchString(didStr)
	if !matched {
		return emptyDID, ErrDoesNotMatchRegexp
	}

	arg := strings.Split(didStr, ":")

	// validate id
	did.ID, err = IDFromString(arg[4])
	if err != nil {
		return emptyDID, err
	}

	did.NetworkID = NetworkID(arg[3])
	did.Blockchain = Blockchain(arg[2])

	return &did, nil
}
