package core

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	// ErrUnsupportedID ID with unsupported type.
	ErrUnsupportedID = errors.New("unsupported ID")
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

func (did *DID) UnmarshalJSON(bytes []byte) error {
	var didStr string
	err := json.Unmarshal(bytes, &didStr)
	if err != nil {
		return err
	}

	did3, err := Parse(didStr)
	if err != nil {
		return err
	}
	*did = *did3
	return nil
}

func (did DID) MarshalJSON() ([]byte, error) {
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

func IDFromDID(did DID) ID {
	id, err := idFromDID(did)
	if err != nil {
		return newIDFromDID(did)
	}
	return id
}

func newIDFromDID(did DID) ID {
	checkSum := sha256.Sum256([]byte(did.String()))
	var genesis [27]byte
	copy(genesis[:], checkSum[len(checkSum)-27:])
	return NewID(TypeUnknown, genesis)
}

func idFromDID(did DID) (ID, error) {
	var id ID
	_, ok := DIDMethodNetwork[DIDMethod(did.Method)]
	if !ok {
		return id, ErrUnsupportedID
	}

	if len(did.IDStrings) > 3 || len(did.IDStrings) < 1 {
		return id, ErrUnsupportedID
	}

	var err error
	id, err = IDFromString(did.IDStrings[len(did.IDStrings)-1])
	if err != nil {
		return id, ErrUnsupportedID
	}

	if !CheckChecksum(id) {
		return id, ErrUnsupportedID
	}

	method, blockchain, networkID, err := decodeDIDPartsFromID(id)
	if err != nil {
		return id, ErrUnsupportedID
	}

	if string(method) != did.Method {
		return id, ErrUnsupportedID
	}

	if len(did.IDStrings) > 1 && string(blockchain) != did.IDStrings[0] {
		return id, ErrUnsupportedID
	}

	if len(did.IDStrings) > 2 && string(networkID) != did.IDStrings[1] {
		return id, ErrUnsupportedID
	}

	return id, nil
}

// ParseDIDFromID returns DID from ID
func ParseDIDFromID(id ID) (*DID, error) {

	if id.IsUnknown() {
		return nil, fmt.Errorf("%w: unknown type", ErrUnsupportedID)
	}

	if !CheckChecksum(id) {
		return nil, fmt.Errorf("%w: invalid checksum", ErrUnsupportedID)
	}

	method, blockchain, networkID, err := decodeDIDPartsFromID(id)
	if err != nil {
		return nil, err
	}

	didParts := []string{DIDSchema, string(method), string(blockchain)}
	if string(networkID) != "" {
		didParts = append(didParts, string(networkID))
	}

	didParts = append(didParts, id.String())

	didString := strings.Join(didParts, ":")

	did, err := Parse(didString)
	if err != nil {
		return nil, err
	}
	return did, nil
}

func decodeDIDPartsFromID(id ID) (DIDMethod, Blockchain, NetworkID, error) {
	methodByte := id.MethodByte()
	networkByte := id.BlockchainNetworkByte()

	method, err := FindDIDMethodByValue(methodByte)
	if err != nil {
		return "", UnknownChain, UnknownNetwork, err
	}

	blockchain, err := FindBlockchainForDIDMethodByValue(method, networkByte)
	if err != nil {
		return "", UnknownChain, UnknownNetwork, err
	}

	networkID, err := FindNetworkIDForDIDMethodByValue(method, networkByte)
	if err != nil {
		return "", UnknownChain, UnknownNetwork, err
	}

	return method, blockchain, networkID, nil
}

func MethodFromID(id ID) (DIDMethod, error) {
	if id.IsUnknown() {
		return "", ErrUnsupportedID
	}
	methodByte := id.MethodByte()
	return FindDIDMethodByValue(methodByte)
}

func BlockchainFromID(id ID) (Blockchain, error) {
	if id.IsUnknown() {
		return UnknownChain, ErrUnsupportedID
	}

	method, err := MethodFromID(id)
	if err != nil {
		return UnknownChain, err
	}

	networkByte := id.BlockchainNetworkByte()

	blockchain, err := FindBlockchainForDIDMethodByValue(method, networkByte)
	if err != nil {
		return UnknownChain, err
	}

	return blockchain, nil
}

func NetworkIDFromID(id ID) (NetworkID, error) {
	if id.IsUnknown() {
		return UnknownNetwork, ErrUnsupportedID
	}

	method, err := MethodFromID(id)
	if err != nil {
		return UnknownNetwork, err
	}

	networkByte := id.BlockchainNetworkByte()

	networkID, err := FindNetworkIDForDIDMethodByValue(method, networkByte)
	if err != nil {
		return UnknownNetwork, err
	}

	return networkID, nil
}
