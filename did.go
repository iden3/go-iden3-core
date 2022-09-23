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
	// ETHEREUM is ethereum blockchain network
	ETHEREUM Blockchain = "eth"
	// POLYGON is polygon blockchain network
	POLYGON Blockchain = "polygon"
	// UNKNOWN_CHAIN is used when it's not possible to retrieve blockchain type from identifier
	UNKNOWN_CHAIN Blockchain = "unknown"
)

// NetworkID is method specific network identifier
type NetworkID string

const (
	// MAIN is ethereum main network
	MAIN NetworkID = "main"
	// MUMBAI is polygon mumbai test network
	MUMBAI NetworkID = "mumbai"
	// ROPSTEN is ethereum ropsten test network
	ROPSTEN NetworkID = "ropsten"
	// RINKEBY is ethereum rinkeby test network
	RINKEBY NetworkID = "rinkeby"
	// KOVAN is ethereum kovan test network
	KOVAN NetworkID = "kovan"
	// GOERLI is ethereum goerli test network
	GOERLI NetworkID = "goerli" // goerli
	// UNKNOWN_NETWORK is used when it's not possible to retrieve network from identifier
	UNKNOWN_NETWORK NetworkID = "unknown"
)

// DIDTypeIDEN3Flag is binary represantation of IDEN3 method flag.
var DIDTypeIDEN3Flag byte = 0b11100000 // 3 bytes for did method

// DIDIden3BlockchainType is mapping between blockchain network and its binary representation
var DIDIden3BlockchainType = map[Blockchain]byte{
	ETHEREUM: DIDTypeIDEN3Flag | 0b00000000, // nolint - reason: explicit declaration of 0 byte value
	POLYGON:  DIDTypeIDEN3Flag | 0b00000001,
}

// DIDNetworkType is mapping between network id and its binary representation
var DIDNetworkType = map[NetworkID]byte{
	MAIN:    0b00000000, // nolint - reason: explicit declaration of 0 byte value
	MUMBAI:  0b00000001,
	ROPSTEN: 0b00000010,
	RINKEBY: 0b00000011,
	KOVAN:   0b00000100,
	GOERLI:  0b00000101,
}

// BuildDIDType builds bytes type from chain and network
func BuildDIDType(blockchain Blockchain, network NetworkID) ([2]byte, error) {
	fb, ok := DIDIden3BlockchainType[blockchain]
	if !ok {
		return [2]byte{}, errors.New(fmt.Sprintf("blockchain %s is not defined in core lib", blockchain))
	}
	sb, ok := DIDNetworkType[network]
	if !ok {
		return [2]byte{}, errors.New(fmt.Sprintf("network %s is not defined in core lib", network))
	}
	return [2]byte{fb, sb}, nil
}

// FindNetworkIDByValue finds network by byte value
func FindNetworkIDByValue(_v byte) (NetworkID, error) {
	for k, v := range DIDNetworkType {
		if v == _v {
			return k, nil
		}
	}
	return UNKNOWN_NETWORK, errors.New(fmt.Sprintf("network %x is not defined in core lib", _v))

}

//FindBlockchainByValue finds blockchain type by byte value
func FindBlockchainByValue(_v byte) (Blockchain, error) {
	for k, v := range DIDIden3BlockchainType {
		if v == _v {
			return k, nil
		}
	}
	return UNKNOWN_CHAIN, errors.New(fmt.Sprintf("blockchain %x is not defined in core lib", _v))

}

var (
	// valid id for regexp
	// did:iden3:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE -readonly id. For readonly identifier networkID and
	// network can be empty as this identifier is newer published on chain
	// did:iden3:eth:main:114vgnnCupQMX4wqUBjg5kUya3zMXfPmKc9HNH4TSE - eth network eth networkID, main - network

	didRegex = regexp.MustCompile(`^\b(did):\b(iden3):(\b(eth|polygon):\b(main|mumbai|ropsten|rinkeby|kovan):)?([1-9a-km-zA-HJ-NP-Z]{41,43})$`)

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

func NewDID(didStr string, options ...DIDOption) (*DID, error) {

	did := &DID{}
	var err error

	did.ID, err = IDFromString(didStr)
	if err != nil {
		return nil, err
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
func WithNetwork(blockchain Blockchain, network NetworkID) DIDOption {
	return func(d *DID) error {
		d.Blockchain = blockchain
		d.NetworkID = network
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
	var err error

	matched := didRegex.MatchString(didStr)
	if !matched {
		return nil, ErrDoesNotMatchRegexp
	}

	arg := strings.Split(didStr, ":")

	// validate id
	did.ID, err = IDFromString(arg[4])
	if err != nil {
		return nil, err
	}

	did.NetworkID = NetworkID(arg[3])
	did.Blockchain = Blockchain(arg[2])

	return &did, nil
}

// ParseDIDFromID returns did from ID
func ParseDIDFromID(id ID) (*DID, error) {
	var err error
	did := DID{}
	did.ID = id
	idBytes := id.Bytes()
	did.Blockchain, err = FindBlockchainByValue(idBytes[0])
	if err != nil {
		return nil, err
	}
	did.NetworkID, err = FindNetworkIDByValue(idBytes[1])
	if err != nil {
		return nil, err
	}
	return &did, nil
}
