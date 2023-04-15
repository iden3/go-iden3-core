package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	// ErrInvalidDID invalid did format.
	ErrInvalidDID = errors.New("invalid did format")
	// ErrDIDMethodNotSupported unsupported did method.
	ErrDIDMethodNotSupported = errors.New("not supported did method")
	// ErrBlockchainNotSupportedForDID unsupported network for did.
	ErrBlockchainNotSupportedForDID = errors.New("not supported blockchain")
	// ErrNetworkNotSupportedForDID unsupported network for did.
	ErrNetworkNotSupportedForDID = errors.New("not supported network")
)

// DIDSchema DID Schema
const DIDSchema = "did"

// DIDMethod represents did methods
type DIDMethod string

const (
	// DIDMethodIden3 DID method-name
	DIDMethodIden3 DIDMethod = "iden3"
	// DIDMethodPolygonID DID method-name
	DIDMethodPolygonID DIDMethod = "polygonid"
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
	// ReadOnly should be used for readonly identity to build readonly flag
	ReadOnly Blockchain = "readonly"
	// NoChain can be used for identity to build readonly flag
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
	DIDMethodIden3:     0b00000001,
	DIDMethodPolygonID: 0b00000010,
}

// DIDNetworkFlag is a structure to represent DID blockchain and network id
type DIDNetworkFlag struct {
	Blockchain Blockchain
	NetworkID  NetworkID
}

// DIDMethodNetwork is map for did methods and their blockchain networks
var DIDMethodNetwork = map[DIDMethod]map[DIDNetworkFlag]byte{
	DIDMethodIden3: {
		{Blockchain: ReadOnly, NetworkID: NoNetwork}: 0b00000000,

		{Blockchain: Polygon, NetworkID: Main}:   0b00010000 | 0b00000001,
		{Blockchain: Polygon, NetworkID: Mumbai}: 0b00010000 | 0b00000010,

		{Blockchain: Ethereum, NetworkID: Main}:   0b00100000 | 0b00000001,
		{Blockchain: Ethereum, NetworkID: Goerli}: 0b00100000 | 0b00000010,
	},
	DIDMethodPolygonID: {
		{Blockchain: ReadOnly, NetworkID: NoNetwork}: 0b00000000,

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
		return [2]byte{}, ErrDIDMethodNotSupported
	}

	if blockchain == NoChain {
		blockchain = ReadOnly
	}

	sb, ok := DIDMethodNetwork[method][DIDNetworkFlag{Blockchain: blockchain, NetworkID: network}]
	if !ok {
		return [2]byte{}, ErrNetworkNotSupportedForDID
	}
	return [2]byte{fb, sb}, nil
}

// BuildDIDTypeOnChain builds bytes type from chain and network
func BuildDIDTypeOnChain(method DIDMethod, blockchain Blockchain, network NetworkID) ([2]byte, error) {
	typ, err := BuildDIDType(method, blockchain, network)
	if err != nil {
		return [2]byte{}, err
	}

	// set on-chain flag (first bit of first byte) to 1
	typ[0] |= MethodOnChainFlag

	return typ, nil
}

// FindNetworkIDForDIDMethodByValue finds network by byte value
func FindNetworkIDForDIDMethodByValue(method DIDMethod, _v byte) (NetworkID, error) {
	_, ok := DIDMethodNetwork[method]
	if !ok {
		return UnknownNetwork, ErrDIDMethodNotSupported
	}
	for k, v := range DIDMethodNetwork[method] {
		if v == _v {
			return k.NetworkID, nil
		}
	}
	return UnknownNetwork, ErrNetworkNotSupportedForDID
}

// FindBlockchainForDIDMethodByValue finds blockchain type by byte value
func FindBlockchainForDIDMethodByValue(method DIDMethod, _v byte) (Blockchain, error) {
	_, ok := DIDMethodNetwork[method]
	if !ok {
		return UnknownChain, ErrDIDMethodNotSupported
	}
	for k, v := range DIDMethodNetwork[method] {
		if v == _v {
			return k.Blockchain, nil
		}
	}
	return UnknownChain, ErrNetworkNotSupportedForDID
}

// FindDIDMethodByValue finds did method by its byte value
func FindDIDMethodByValue(_v byte) (DIDMethod, error) {
	for k, v := range DIDMethodByte {
		if v == _v {
			return k, nil
		}
	}
	return "", ErrDIDMethodNotSupported
}

// DID Decentralized Identifiers (DIDs)
// https://w3c.github.io/did-core/#did-syntax
type DID struct {
	ID              ID         // ID did specific id
	Method          DIDMethod  // DIDMethod did method
	Blockchain      Blockchain // Blockchain network identifier eth / polygon,...
	NetworkID       NetworkID  // NetworkID specific network identifier eth {main, mumbai, goerli}
	isUnknownMethod bool       // isUnknownMethod specifies if DID is of unsupported method
	didString       string     // didString full DID string identifier
}

func (did *DID) SetString(didStr string) error {
	// Parse method, networks and id
	arg := strings.Split(didStr, ":")
	if len(arg) < 3 {
		return ErrInvalidDID
	}

	if arg[0] != DIDSchema {
		return ErrInvalidDID
	}

	did.didString = didStr
	did.Method = DIDMethod(arg[1])

	// check did method defined in core lib
	_, ok := DIDMethodByte[did.Method]
	if !ok {
		did.isUnknownMethod = true

		// TODO: algo how to encode unknown did to ID
		// keccac256 and take 27 bytes from it

		return nil
	}

	switch len(arg) {
	case 5:
		var err error

		// we have both blockchain and network id specified
		did.Blockchain = Blockchain(arg[2])
		did.NetworkID = NetworkID(arg[3])

		// parse id
		did.ID, err = IDFromString(arg[4])
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidDID, err)
		}

	case 4:
		var err error

		// we have only blockchain specified
		did.Blockchain = Blockchain(arg[2])

		// set default network id
		switch did.Blockchain {
		case ReadOnly:
			did.NetworkID = NoNetwork
		case Polygon:
			did.NetworkID = Main
		case Ethereum:
			did.NetworkID = Main
		default:
			return ErrNetworkNotSupportedForDID
		}

		did.ID, err = IDFromString(arg[3])
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidDID, err)
		}

	case 3:
		var err error

		// we do not have blockchain & network id, set default ones
		did.Blockchain = Polygon
		did.NetworkID = Main

		did.ID, err = IDFromString(arg[2])
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidDID, err)
		}
	default:
		return ErrInvalidDID
	}

	// check did network defined in core lib for did method
	_, ok = DIDMethodNetwork[did.Method][DIDNetworkFlag{
		Blockchain: did.Blockchain,
		NetworkID:  did.NetworkID}]
	if !ok {
		return ErrNetworkNotSupportedForDID
	}

	if !CheckChecksum(did.ID) {
		return fmt.Errorf("%w: %s", ErrInvalidDID, "invalid checksum")
	}

	// check id contains did network and method
	return did.validate()
}

// Return nil on success or error if fields are inconsistent.
func (did *DID) validate() error {
	d, err := ParseDIDFromID(did.ID)
	if err != nil {
		return err
	}

	if d.Method != did.Method {
		return fmt.Errorf(
			"%w: did method of core identity %s differs from given did method %s",
			ErrInvalidDID, d.Method, did.Method)
	}

	if d.NetworkID != did.NetworkID {
		return fmt.Errorf(
			"%w: network method of core identity %s differs from given did network specific id %s",
			ErrInvalidDID, d.NetworkID, did.NetworkID)
	}

	if d.Blockchain != did.Blockchain {
		return fmt.Errorf(
			"%w: blockchain network of core identity %s differs from given did blockhain network %s",
			ErrInvalidDID, d.Blockchain, did.Blockchain)
	}

	if !bytes.Equal(d.ID[:], did.ID[:]) {
		return fmt.Errorf(
			"%w: ID of core identity %s differs from given did ID %s",
			ErrInvalidDID, d.ID.String(), did.ID.String())
	}

	return nil
}

func (did *DID) UnmarshalJSON(bytes []byte) error {
	var didStr string
	err := json.Unmarshal(bytes, &didStr)
	if err != nil {
		return err
	}

	return did.SetString(didStr)
}

func (did *DID) MarshalJSON() ([]byte, error) {
	return json.Marshal(did.String())
}

// DIDGenesisFromIdenState calculates the genesis ID from an Identity State and returns it as DID
func DIDGenesisFromIdenState(typ [2]byte, state *big.Int) (*DID, error) {
	id, err := IdGenesisFromIdenState(typ, state)
	if err != nil {
		return nil, err
	}
	return ParseDIDFromID(*id)
}

// String did as a string
func (did *DID) String() string {
	if did.isUnknownMethod {
		return did.didString
	}

	if did.Blockchain == NoChain || did.Blockchain == ReadOnly {
		return fmt.Sprintf("%s:%s:%s:%s", DIDSchema, did.Method, ReadOnly, did.ID.String())
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s", DIDSchema, did.Method, did.Blockchain,
		did.NetworkID, did.ID.String())
}

func (did *DID) CoreID() ID {
	return did.ID
}

// ParseDID method parse string and extract DID if string is valid Iden3 identifier
func ParseDID(didStr string) (*DID, error) {
	var did DID
	err := did.SetString(didStr)
	return &did, err
}

// ParseDIDFromID returns did from ID
func ParseDIDFromID(id ID) (*DID, error) {
	var err error
	did := DID{}
	did.ID = id
	method := id.MethodByte()
	net := id.BlockchainNetworkByte()

	did.Method, err = FindDIDMethodByValue(method)
	if err != nil {
		return nil, err
	}

	did.Blockchain, err = FindBlockchainForDIDMethodByValue(did.Method, net)
	if err != nil {
		return nil, err
	}

	did.NetworkID, err = FindNetworkIDForDIDMethodByValue(did.Method, net)
	if err != nil {
		return nil, err
	}

	return &did, nil
}
