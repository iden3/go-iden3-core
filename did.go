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
	// ErrIncorrectDID return if DID method is known, but format of DID is incorrect.
	ErrIncorrectDID = errors.New("incorrect DID")
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
type DIDMethod byte

const (
	// DIDMethodIden3
	DIDMethodIden3 DIDMethod = 0b00000001
	// DIDMethodPolygonID
	DIDMethodPolygonID DIDMethod = 0b00000010
	// DIDMethodPolygonIDOnChain
	DIDMethodPolygonIDOnChain DIDMethod = 0b10000010
	// DIDMethodOther any other method not listed before
	DIDMethodOther DIDMethod = 0b11111111
)

var knownMethods = map[DIDMethod]struct{}{
	DIDMethodIden3:            {},
	DIDMethodPolygonID:        {},
	DIDMethodPolygonIDOnChain: {},
}

var onChainMethods = map[DIDMethod]struct{}{
	DIDMethodPolygonIDOnChain: {},
}

func (m DIDMethod) String() string {
	switch m {
	case DIDMethodIden3:
		return "iden3"
	case DIDMethodPolygonID:
		return "polygonid"
	case DIDMethodPolygonIDOnChain:
		return "polygonid"
	case DIDMethodOther:
		return ""
	default:
		return fmt.Sprintf("unknown<%v>", uint8(m))
	}
}

func (m DIDMethod) Byte() byte {
	return byte(m)
}

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
	DIDMethodPolygonIDOnChain: {
		{Blockchain: Polygon, NetworkID: Main}:   0b00010000 | 0b00000001,
		{Blockchain: Polygon, NetworkID: Mumbai}: 0b00010000 | 0b00000010,
	},
	DIDMethodOther: {},
}

// BuildDIDType builds bytes type from chain and network
func BuildDIDType(method DIDMethod, blockchain Blockchain,
	network NetworkID) ([2]byte, error) {

	if _, ok := knownMethods[method]; !ok {
		return [2]byte{}, ErrDIDMethodNotSupported
	}

	if blockchain == NoChain {
		blockchain = ReadOnly
	}

	netFlag := DIDNetworkFlag{Blockchain: blockchain, NetworkID: network}
	sb, ok := DIDMethodNetwork[method][netFlag]
	if !ok {
		return [2]byte{}, ErrNetworkNotSupportedForDID
	}

	return [2]byte{method.Byte(), sb}, nil
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
	return UnknownChain, ErrBlockchainNotSupportedForDID
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
	hash := sha256.Sum256([]byte(did.String()))
	var genesis [27]byte
	copy(genesis[:], hash[len(hash)-27:])
	return NewID(TypeUnknown, genesis)
}

func idFromDID(did DID) (ID, error) {
	found := false
	for method := range knownMethods {
		if method.String() == did.Method {
			found = true
			break
		}
	}
	if !found {
		return ID{}, ErrUnsupportedID
	}

	var id ID

	if len(did.IDStrings) > 3 || len(did.IDStrings) < 1 {
		return id, fmt.Errorf("%w: unexpected number of ID strings",
			ErrIncorrectDID)
	}

	var err error
	id, err = IDFromString(did.IDStrings[len(did.IDStrings)-1])
	if err != nil {
		return id, fmt.Errorf("%w: can't parse ID string", ErrIncorrectDID)
	}

	if !CheckChecksum(id) {
		return id, fmt.Errorf("%w: incorrect ID checksum", ErrIncorrectDID)
	}

	method, blockchain, networkID, err := decodeDIDPartsFromID(id)
	if err != nil {
		return id, ErrUnsupportedID
	}

	if method.String() != did.Method {
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

	didParts := []string{DIDSchema, method.String(), string(blockchain)}
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
	method := id.Method()
	networkByte := id.BlockchainNetworkByte()

	blockchain, err := FindBlockchainForDIDMethodByValue(method, networkByte)
	if err != nil {
		return DIDMethodOther, UnknownChain, UnknownNetwork, err
	}

	networkID, err := FindNetworkIDForDIDMethodByValue(method, networkByte)
	if err != nil {
		return DIDMethodOther, UnknownChain, UnknownNetwork, err
	}

	return method, blockchain, networkID, nil
}

func MethodFromID(id ID) (DIDMethod, error) {
	if id.IsUnknown() {
		return DIDMethodOther, fmt.Errorf("%w: unknown type", ErrUnsupportedID)
	}
	method := id.Method()
	if _, ok := knownMethods[method]; !ok && method != DIDMethodOther {
		return DIDMethodOther, fmt.Errorf("%w: unknown method",
			ErrUnsupportedID)
	}
	return method, nil
}

func BlockchainFromID(id ID) (Blockchain, error) {
	if id.IsUnknown() {
		return UnknownChain, fmt.Errorf("%w: unknown type", ErrUnsupportedID)
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
		return UnknownNetwork, fmt.Errorf("%w: unknown type", ErrUnsupportedID)
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

func EthAddressFromID(id ID) ([20]byte, error) {
	if _, ok := onChainMethods[id.Method()]; !ok {
		return [20]byte{}, errors.New(
			"can't get Ethereum address of not on-chain identity")
	}
	var address [20]byte
	copy(address[:], id[2:22])
	return address, nil
}
