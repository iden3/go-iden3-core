package core

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/iden3/go-iden3-core/v2/w3c"
)

var (
	// ErrUnsupportedID ID with unsupported type.
	ErrUnsupportedID = errors.New("unsupported ID")
	// ErrIncorrectDID return if DID method is known, but format of DID is incorrect.
	ErrIncorrectDID = errors.New("incorrect DID")
	// ErrMethodUnknown return if DID method is unknown.
	ErrMethodUnknown = errors.New("unknown DID method")
	// ErrDIDMethodNotSupported unsupported did method.
	ErrDIDMethodNotSupported = errors.New("not supported did method")
	// ErrBlockchainNotSupportedForDID unsupported network for did.
	ErrBlockchainNotSupportedForDID = errors.New("not supported blockchain")
	// ErrNetworkNotSupportedForDID unsupported network for did.
	ErrNetworkNotSupportedForDID = errors.New("not supported network")
)

// DIDMethod represents did methods
type DIDMethod string

const (
	// DIDMethodIden3
	DIDMethodIden3 DIDMethod = "iden3"
	// DIDMethodPolygonID
	DIDMethodPolygonID DIDMethod = "polygonid"
	// DIDMethodOther any other method not listed before
	DIDMethodOther DIDMethod = ""
)

var didMethods = map[DIDMethod]DIDMethod{
	DIDMethodIden3:     DIDMethodIden3,
	DIDMethodPolygonID: DIDMethodPolygonID,
	DIDMethodOther:     DIDMethodOther,
}

// GetDIDMethod returns DID method by name
func GetDIDMethod(name string) (DIDMethod, error) {

	method, ok := didMethods[DIDMethod(name)]
	if !ok {
		return DIDMethodOther, fmt.Errorf("DID method '%s' not found", name)
	}
	return method, nil
}

// Blockchain id of the network "eth", "polygon", etc.
type Blockchain string

const (
	// Ethereum is ethereum blockchain network
	Ethereum Blockchain = "eth"
	// Polygon is polygon blockchain network
	Polygon Blockchain = "polygon"
	// Privado is Privado blockchain network
	Privado Blockchain = "privado"
	// Linea is Linea blockchain network
	Linea Blockchain = "linea"
	// UnknownChain is used when it's not possible to retrieve blockchain type from identifier
	UnknownChain Blockchain = "unknown"
	// ReadOnly should be used for readonly identity to build readonly flag
	ReadOnly Blockchain = "readonly"
	// NoChain can be used for identity to build readonly flag
	NoChain Blockchain = ""
)

var blockchains = map[Blockchain]Blockchain{
	Ethereum:     Ethereum,
	Polygon:      Polygon,
	Privado:      Privado,
	Linea:        Linea,
	UnknownChain: UnknownChain,
	ReadOnly:     ReadOnly,
	NoChain:      NoChain,
}

// GetBlockchain returns blockchain by name
func GetBlockchain(name string) (Blockchain, error) {
	blockchain, ok := blockchains[Blockchain(name)]
	if !ok {
		return UnknownChain, fmt.Errorf("blockchain '%s' not found", name)
	}
	return blockchain, nil
}

// RegisterBlockchain registers new blockchain
func RegisterBlockchain(b Blockchain) error {
	blockchains[b] = b
	return nil
}

// NetworkID is method specific network identifier
type NetworkID string

// Generic NetworkIDs
const (
	// Main is main network
	Main NetworkID = "main"
	// Test is test network
	Test NetworkID = "test"
	// UnknownNetwork is used when it's not possible to retrieve network from identifier
	UnknownNetwork NetworkID = "unknown"
	// NoNetwork should be used for readonly identity to build readonly flag
	NoNetwork NetworkID = ""
)

// Ethereum-specific NetworkIDs
const (
	// Goerli is Ethereum goerli test network
	Goerli NetworkID = "goerli"
	// Sepolia is Ethereum Sepolia test network
	Sepolia NetworkID = "sepolia"
)

// Polygon-specific NetworkIDs
const (
	// Mumbai is Polygon mumbai test network
	Mumbai NetworkID = "mumbai"
	// Amoy is Polygon amoy test network
	Amoy NetworkID = "amoy"
	// Zkevm is zkEVM network in Polygon and potentially other blockchains
	Zkevm NetworkID = "zkevm"
	// Cardona is Polygon zkEVM Cardona test network
	Cardona NetworkID = "cardona"
)

var networks = map[NetworkID]NetworkID{
	Main:           Main,
	Mumbai:         Mumbai,
	Amoy:           Amoy,
	Zkevm:          Zkevm,
	Cardona:        Cardona,
	Goerli:         Goerli,
	Sepolia:        Sepolia,
	Test:           Test,
	UnknownNetwork: UnknownNetwork,
	NoNetwork:      NoNetwork,
}

// GetNetwork returns network by name
func GetNetwork(name string) (NetworkID, error) {
	network, ok := networks[NetworkID(name)]
	if !ok {
		return UnknownNetwork, fmt.Errorf("network '%s' not found", name)
	}
	return network, nil
}

// RegisterNetwork registers new network
func RegisterNetwork(n NetworkID) error {
	networks[n] = n
	return nil
}

// DIDMethodByte did method flag representation
var DIDMethodByte = map[DIDMethod]byte{
	DIDMethodIden3:     0b00000001,
	DIDMethodPolygonID: 0b00000010,
	DIDMethodOther:     0b11111111,
}

// RegisterDIDMethod registers new DID method with byte flag
func RegisterDIDMethod(m DIDMethod, b byte) error {

	max := DIDMethodByte[DIDMethodOther]
	if b >= max {
		return fmt.Errorf("Can't register DID method byte: current %b, maximum byte allowed: %b", b, max-1)
	}

	existingByte, ok := DIDMethodByte[m]
	if ok && existingByte == b {
		return nil
	}

	for _, v := range DIDMethodByte {
		if v == b {
			return fmt.Errorf(`can't register method '%s' because DID method byte '%b' already registered for another method`, m, b)
		}
	}

	didMethods[m] = m
	DIDMethodByte[m] = b

	return nil
}

// DIDNetworkFlag is a structure to represent DID blockchain and network id
type DIDNetworkFlag struct {
	Blockchain Blockchain
	NetworkID  NetworkID
}

var blockchainNetworkMap = map[DIDNetworkFlag]byte{
	{Blockchain: ReadOnly, NetworkID: NoNetwork}: 0b0000_0000,

	{Blockchain: Polygon, NetworkID: Main}:    0b0001_0000 | 0b0000_0001,
	{Blockchain: Polygon, NetworkID: Mumbai}:  0b0001_0000 | 0b0000_0010,
	{Blockchain: Polygon, NetworkID: Amoy}:    0b0001_0000 | 0b0000_0011,
	{Blockchain: Polygon, NetworkID: Zkevm}:   0b0001_0000 | 0b0000_0100,
	{Blockchain: Polygon, NetworkID: Cardona}: 0b0001_0000 | 0b0000_0101,

	{Blockchain: Ethereum, NetworkID: Main}:    0b0010_0000 | 0b0000_0001,
	{Blockchain: Ethereum, NetworkID: Goerli}:  0b0010_0000 | 0b0000_0010,
	{Blockchain: Ethereum, NetworkID: Sepolia}: 0b0010_0000 | 0b0000_0011,

	{Blockchain: Privado, NetworkID: Main}: 0b1010_0000 | 0b0000_0001,
	{Blockchain: Privado, NetworkID: Test}: 0b1010_0000 | 0b0000_0010,

	{Blockchain: Linea, NetworkID: Main}:    0b0100_0000 | 0b0000_1001,
	{Blockchain: Linea, NetworkID: Sepolia}: 0b0100_0000 | 0b0000_1000,
}

// DIDMethodNetwork is map for did methods and their blockchain networks
var DIDMethodNetwork = map[DIDMethod]map[DIDNetworkFlag]byte{
	DIDMethodIden3:     blockchainNetworkMap,
	DIDMethodPolygonID: blockchainNetworkMap,
	DIDMethodOther: {
		{Blockchain: UnknownChain, NetworkID: UnknownNetwork}: 0b1111_1111,
	},
}

// DIDMethodNetworkParams is a structure to represent DID method network options
type DIDMethodNetworkParams struct {
	Method      DIDMethod
	Blockchain  Blockchain
	Network     NetworkID
	NetworkFlag byte
}

type registrationOptions struct {
	chainID    *int
	methodByte *byte
}

// RegistrationOptions is a type for DID method network options
type RegistrationOptions func(params *registrationOptions)

// WithChainID registers new chain ID method with byte flag
func WithChainID(chainID int) RegistrationOptions {
	return func(opts *registrationOptions) {
		opts.chainID = &chainID
	}
}

// WithDIDMethodByte registers new DID method with byte flag
func WithDIDMethodByte(methodByte byte) RegistrationOptions {
	return func(opts *registrationOptions) {
		opts.methodByte = &methodByte
	}
}

// RegisterDIDMethodNetwork registers new DID method network
func RegisterDIDMethodNetwork(params DIDMethodNetworkParams, opts ...RegistrationOptions) error {
	var err error
	o := registrationOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	b := params.Blockchain
	n := params.Network
	m := params.Method

	err = RegisterBlockchain(b)
	if err != nil {
		return err
	}

	err = RegisterNetwork(n)
	if err != nil {
		return err
	}

	if o.methodByte != nil {
		err = RegisterDIDMethod(m, *o.methodByte)
		if err != nil {
			return err
		}
	}

	flg := DIDNetworkFlag{Blockchain: b, NetworkID: n}

	if _, ok := DIDMethodNetwork[m]; !ok {
		DIDMethodNetwork[m] = map[DIDNetworkFlag]byte{}
	}

	if o.chainID != nil {
		err = RegisterChainID(b, n, *o.chainID)
		if err != nil {
			return err
		}
	}
	existedFlag, ok := DIDMethodNetwork[m][flg]
	if ok && existedFlag == params.NetworkFlag {
		return nil
	}

	for _, v := range DIDMethodNetwork[m] {
		if v == params.NetworkFlag {
			return fmt.Errorf(`DID network flag %b is already registered for the another network id for '%s' method`, v, m)
		}
	}

	DIDMethodNetwork[m][flg] = params.NetworkFlag
	return nil

}

// BuildDIDType builds bytes type from chain and network
func BuildDIDType(method DIDMethod, blockchain Blockchain,
	network NetworkID) ([2]byte, error) {

	fb, ok := DIDMethodByte[method]
	if !ok {
		return [2]byte{}, ErrDIDMethodNotSupported
	}

	netFlag := DIDNetworkFlag{Blockchain: blockchain, NetworkID: network}
	sb, ok := DIDMethodNetwork[method][netFlag]
	if !ok {
		return [2]byte{}, ErrNetworkNotSupportedForDID
	}

	return [2]byte{fb, sb}, nil
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

// FindDIDMethodByValue finds did method by its byte value
func FindDIDMethodByValue(b byte) (DIDMethod, error) {
	for k, v := range DIDMethodByte {
		if v == b {
			return k, nil
		}
	}
	return DIDMethodOther, ErrDIDMethodNotSupported
}

// NewDIDFromIdenState calculates the genesis ID from an Identity State and
// returns it as a DID
func NewDIDFromIdenState(typ [2]byte, state *big.Int) (*w3c.DID, error) {
	id, err := NewIDFromIdenState(typ, state)
	if err != nil {
		return nil, err
	}
	return ParseDIDFromID(*id)
}

// NewDID creates a new *w3c.DID from the type and the genesis
func NewDID(typ [2]byte, genesis [genesisLn]byte) (*w3c.DID, error) {
	return ParseDIDFromID(NewID(typ, genesis))
}

func IDFromDID(did w3c.DID) (ID, error) {
	id, err := idFromDID(did)
	if errors.Is(err, ErrMethodUnknown) {
		return newIDFromUnsupportedDID(did), nil
	}
	return id, err
}

func newIDFromUnsupportedDID(did w3c.DID) ID {
	hash := sha256.Sum256([]byte(did.String()))
	var genesis [genesisLn]byte
	copy(genesis[:], hash[len(hash)-genesisLn:])
	flg := DIDNetworkFlag{Blockchain: UnknownChain, NetworkID: UnknownNetwork}
	var tp = [2]byte{
		DIDMethodByte[DIDMethodOther],
		DIDMethodNetwork[DIDMethodOther][flg],
	}
	return NewID(tp, genesis)
}

func idFromDID(did w3c.DID) (ID, error) {
	method, ok := didMethods[DIDMethod(did.Method)]
	if !ok {
		method = DIDMethodOther
	}
	_, ok = DIDMethodByte[method]
	if !ok || method == DIDMethodOther {
		return ID{}, ErrMethodUnknown
	}

	var id ID

	if len(did.IDStrings) > 3 || len(did.IDStrings) < 2 {
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

	method2, blockchain, networkID, err := decodeDIDPartsFromID(id)
	if err != nil {
		return id, err
	}

	if method2 != method {
		return id, fmt.Errorf("%w: methods in ID and DID are different",
			ErrIncorrectDID)
	}

	if string(blockchain) != did.IDStrings[0] {
		return id, fmt.Errorf("%w: blockchains in ID and DID are different",
			ErrIncorrectDID)
	}

	if len(did.IDStrings) > 2 && string(networkID) != did.IDStrings[1] {
		return id, fmt.Errorf("%w: networkIDs in ID and DID are different",
			ErrIncorrectDID)
	}

	return id, nil
}

// ParseDIDFromID returns DID from ID
func ParseDIDFromID(id ID) (*w3c.DID, error) {

	if !CheckChecksum(id) {
		return nil, fmt.Errorf("%w: invalid checksum", ErrUnsupportedID)
	}

	method, blockchain, networkID, err := decodeDIDPartsFromID(id)
	if err != nil {
		return nil, err
	}

	if isUnsupportedDID(method, blockchain, networkID) {
		return nil, fmt.Errorf("%w: unsupported DID",
			ErrMethodUnknown)
	}

	didParts := []string{"did", string(method), string(blockchain)}
	if string(networkID) != "" {
		didParts = append(didParts, string(networkID))
	}

	didParts = append(didParts, id.String())

	didString := strings.Join(didParts, ":")

	did, err := w3c.ParseDID(didString)
	if err != nil {
		return nil, err
	}
	return did, nil
}

func decodeDIDPartsFromID(id ID) (DIDMethod, Blockchain, NetworkID, error) {
	method, err := FindDIDMethodByValue(id[0])
	if err != nil {
		return DIDMethodOther, UnknownChain, UnknownNetwork, err
	}

	blockchain, err := FindBlockchainForDIDMethodByValue(method, id[1])
	if err != nil {
		return DIDMethodOther, UnknownChain, UnknownNetwork, err
	}

	networkID, err := FindNetworkIDForDIDMethodByValue(method, id[1])
	if err != nil {
		return DIDMethodOther, UnknownChain, UnknownNetwork, err
	}

	return method, blockchain, networkID, nil
}

func MethodFromID(id ID) (DIDMethod, error) {
	method, blockchain, netID, err := decodeDIDPartsFromID(id)
	if err != nil {
		return DIDMethodOther, err
	}

	if isUnsupportedDID(method, blockchain, netID) {
		return DIDMethodOther, fmt.Errorf("%w: unsupported DID",
			ErrMethodUnknown)
	}

	return method, nil
}

func BlockchainFromID(id ID) (Blockchain, error) {
	method, blockchain, netID, err := decodeDIDPartsFromID(id)
	if err != nil {
		return UnknownChain, err
	}

	if isUnsupportedDID(method, blockchain, netID) {
		return UnknownChain, fmt.Errorf("%w: unsupported DID",
			ErrMethodUnknown)
	}

	return blockchain, nil
}

func NetworkIDFromID(id ID) (NetworkID, error) {
	method, blockchain, netID, err := decodeDIDPartsFromID(id)
	if err != nil {
		return UnknownNetwork, err
	}

	if isUnsupportedDID(method, blockchain, netID) {
		return UnknownNetwork, fmt.Errorf("%w: unsupported DID",
			ErrMethodUnknown)
	}

	return netID, nil
}

func EthAddressFromID(id ID) ([20]byte, error) {
	var z [7]byte
	if !bytes.Equal(z[:], id[2:2+len(z)]) {
		return [20]byte{}, errors.New(
			"can't get Ethereum address: high bytes of genesis are not zero")
	}

	var address [20]byte
	copy(address[:], id[2+7:])
	return address, nil
}

func GenesisFromEthAddress(addr [20]byte) [genesisLn]byte {
	var genesis [genesisLn]byte
	copy(genesis[7:], addr[:])
	return genesis
}

func isUnsupportedDID(method DIDMethod, blockchain Blockchain,
	networkID NetworkID) bool {

	return method == DIDMethodOther && blockchain == UnknownChain &&
		networkID == UnknownNetwork
}
