package core

import (
	"fmt"
	"math/big"
	"strings"
)

// DIDSchema DID Schema
const DIDSchema = "did"

// DIDMethod represents did methods
type DIDMethod string

const (
	// DIDMethodIden3 DID method-name
	DIDMethodIden3 DIDMethod = "iden3"
	// DIDMethodPolygonID DID method-name
	DIDMethodPolygonID DIDMethod = "polygonID"
)

// Blockchain id of the network "eth", "polygon", etc.
type Blockchain string

const (
	// Ethereum is ethereum blockchain network
	Ethereum Blockchain = "eth"
	// Polygon is polygon blockchain network
	Polygon Blockchain = "polygon"
	// UnknownChain is used when it's not possible to retrieve blockchain type from identifier
	UnknownChain Blockchain = "unknown"
	// NoChain should be used for readonly identity to build readonly flag
	NoChain Blockchain = ""
)

// NetworkID is method specific network identifier
type NetworkID string

const (
	// Main is ethereum main network
	Main NetworkID = "main"
	// Mumbai is polygon mumbai test network
	Mumbai NetworkID = "mumbai"

	// Goerli is ethereum goerli test network
	Goerli NetworkID = "goerli" // goerli
	// UnknownNetwork is used when it's not possible to retrieve network from identifier
	UnknownNetwork NetworkID = "unknown"

	// NoNetwork should be used for readonly identity to build readonly flag
	NoNetwork NetworkID = ""
)

// DIDMethodByte did method flag representation
var DIDMethodByte = map[DIDMethod]byte{
	DIDMethodIden3: 0b00000001,
}

// DIDNetworkFlag is a structure to represent DID blockchain and network id
type DIDNetworkFlag struct {
	Blockchain Blockchain
	NetworkID  NetworkID
}

// DIDMethodNetwork is map for did methods and their blockchain networks
var DIDMethodNetwork = map[DIDMethod]map[DIDNetworkFlag]byte{
	DIDMethodIden3: {
		{Blockchain: NoChain, NetworkID: NoNetwork}: 0b00000000,

		{Blockchain: Polygon, NetworkID: Main}:   0b00010000 | 0b00000001,
		{Blockchain: Polygon, NetworkID: Mumbai}: 0b00010000 | 0b00000010,

		{Blockchain: Ethereum, NetworkID: Main}:   0b00100000 | 0b00000001,
		{Blockchain: Ethereum, NetworkID: Goerli}: 0b00100000 | 0b00000010,
	},
}

// BuildDIDType builds bytes type from chain and network
func BuildDIDType(method DIDMethod, blockchain Blockchain, network NetworkID) ([2]byte, error) {

	fb, ok := DIDMethodByte[method]
	if !ok {
		return [2]byte{}, fmt.Errorf("method %s is not defined in core lib", method)
	}

	sb, ok := DIDMethodNetwork[method][DIDNetworkFlag{Blockchain: blockchain, NetworkID: network}]
	if !ok {
		return [2]byte{}, fmt.Errorf("blockchain network %s %s is not defined in core lib", blockchain, network)
	}
	return [2]byte{fb, sb}, nil
}

// FindNetworkIDForDIDMethodByValue finds network by byte value
func FindNetworkIDForDIDMethodByValue(method DIDMethod, _v byte) (NetworkID, error) {
	_, ok := DIDMethodNetwork[method]
	if !ok {
		return UnknownNetwork, fmt.Errorf("did method %x is not defined in core lib", method)
	}
	for k, v := range DIDMethodNetwork[method] {
		if v == _v {
			return k.NetworkID, nil
		}
	}
	return UnknownNetwork, fmt.Errorf("bytes %x for did method %s is not defined in core lib as a valid network identifer", _v, method)
}

// FindBlockchainForDIDMethodByValue finds blockchain type by byte value
func FindBlockchainForDIDMethodByValue(method DIDMethod, _v byte) (Blockchain, error) {
	_, ok := DIDMethodNetwork[method]
	if !ok {
		return UnknownChain, fmt.Errorf("did method %x is not defined in core lib", method)
	}
	for k, v := range DIDMethodNetwork[method] {
		if v == _v {
			return k.Blockchain, nil
		}
	}
	return UnknownChain, fmt.Errorf("bytes %x for did method %s is not defined in core lib as a valid blockchain network", _v, method)
}

// FindDIDMethodByValue finds did method by its byte value
func FindDIDMethodByValue(_v byte) (DIDMethod, error) {
	for k, v := range DIDMethodByte {
		if v == _v {
			return k, nil
		}
	}
	return "", fmt.Errorf("bytes %x are not defined in core lib as valid did method", _v)
}

// DID Decentralized Identifiers (DIDs)
// https://w3c.github.io/did-core/#did-syntax
type DID struct {
	ID         ID         // ID did specific id
	Method     DIDMethod  // DIDMethod did method
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

// DIDGenesisFromIdenState calculates the genesis ID from an Identity State and returns it as DID
func DIDGenesisFromIdenState(typ [2]byte, state *big.Int) (*DID, error) {
	id, err := IdGenesisFromIdenState(typ, state)
	if err != nil {
		return nil, err
	}
	return ParseDIDFromID(*id)
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
		return fmt.Sprintf("%s:%s:%s", DIDSchema, DIDMethodIden3, did.ID.String())
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s", DIDSchema, DIDMethodIden3, did.Blockchain,
		did.NetworkID, did.ID.String())
}

// ParseDID method parse string and extract DID if string is valid Iden3 identifier
func ParseDID(didStr string) (*DID, error) {
	did := DID{}
	var err error

	arg := strings.Split(didStr, ":")

	did.Method = DIDMethod(arg[1])

	switch len(arg) {
	case 5:
		// validate id
		did.ID, err = IDFromString(arg[4])
		if err != nil {
			return nil, err
		}

		did.Blockchain = Blockchain(arg[2])
		did.NetworkID = NetworkID(arg[3])

	case 3:
		// validate readonly id
		did.ID, err = IDFromString(arg[2])
		if err != nil {
			return nil, err
		}
	}

	// check did method defined in core lib
	_, ok := DIDMethodByte[did.Method]
	if !ok {
		return nil, fmt.Errorf(`did method %s is not defined in core lib`, did.Method)
	}

	// check did network defined in core lib for did method
	_, ok = DIDMethodNetwork[did.Method][DIDNetworkFlag{Blockchain: did.Blockchain, NetworkID: did.NetworkID}]
	if !ok {
		return nil, fmt.Errorf(`blockchain network "%s %s" is not defined for %s did method`, did.Blockchain, did.NetworkID, did.Method)
	}

	// check id contains did network and method

	d, err := ParseDIDFromID(did.ID)
	if err != nil {
		return nil, err
	}
	if d.Method != did.Method {
		return nil, fmt.Errorf(`did method of core identity %s differs from given did method %s`, d.Method, did.Method)
	}
	if d.NetworkID != did.NetworkID {
		return nil, fmt.Errorf(`network method of core identity %s differs from given did network specific id %s`, d.NetworkID, did.NetworkID)
	}
	if d.Blockchain != did.Blockchain {
		return nil, fmt.Errorf(`blockchain network of core identity %s differs from given did blockhain network %s`, d.Blockchain, did.Blockchain)
	}
	return &did, nil
}

// ParseDIDFromID returns did from ID
func ParseDIDFromID(id ID) (*DID, error) {
	var err error
	did := DID{}
	did.ID = id
	typ := id.Type()

	did.Method, err = FindDIDMethodByValue(typ[0])
	if err != nil {
		return nil, err
	}
	did.Blockchain, err = FindBlockchainForDIDMethodByValue(did.Method, typ[1])
	if err != nil {
		return nil, err
	}
	did.NetworkID, err = FindNetworkIDForDIDMethodByValue(did.Method, typ[1])
	if err != nil {
		return nil, err
	}
	return &did, nil
}
