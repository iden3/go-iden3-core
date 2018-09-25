package eth

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

type TestContractStore struct {
	t *testing.T
}

func (sto *TestContractStore) Get(id string) (*abi.ABI, []byte, error) {
	path := "../../contracts/build/contracts"
	switch id {
	case deployerContract:
		path = path + "/Deployer.json"
	case proxyContract:
		path = path + "/IDen3DelegateProxy.json"
	case implContract:
		path = path + "/IDen3Impl.json"
	default:
		sto.t.Error("bad contract id")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	abi, bytecode, err := UnmarshallSolcAbiJson(f)

	if id == deployerContract {
		bytecode, err = hex.DecodeString("608060405234801561001057600080fd5b506101a3806100206000396000f3fe6080604052600436106100405763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663cf5ba53f8114610045575b600080fd5b34801561005157600080fd5b506100f86004803603602081101561006857600080fd5b81019060208101813564010000000081111561008357600080fd5b82018360208201111561009557600080fd5b803590602001918460018302840111640100000000831117156100b757600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506100fa945050505050565b005b6000808251602084016000f5905073ffffffffffffffffffffffffffffffffffffffff8116151561012a57600080fd5b6040805173ffffffffffffffffffffffffffffffffffffffff8316815290517f1449abf21e49fd025f33495e77f7b1461caefdd3d4bb646424a3f445c4576a5b9181900360200190a1505056fea165627a7a7230582029a23322cadc2c71f486df1679a9fe8055788bd43d52d633681181b436f31b790029")
		if err != nil {
			return nil, nil, err
		}
	}

	return abi, bytecode, err
}

func TestA(t *testing.T) {

	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	ks := keystore.NewKeyStore(dir, 2, 2)
	accRelayer, err := ks.NewAccount("")
	assert.Nil(t, err)
	assert.Nil(t, ks.Unlock(accRelayer, ""))

	rpc, err := rpc.DialContext(context.TODO(), "http://localhost:8545")
	assert.Nil(t, err)
	defer func() {
		rpc.Close()
	}()

	var raw json.RawMessage
	err = rpc.CallContext(context.TODO(), &raw, "eth_coinbase")
	assert.Nil(t, err)
	assert.True(t, len(raw) > 0)

	var coinbase common.Address
	assert.Nil(t, json.Unmarshal(raw, &coinbase))
	t.Log("coinbase ", coinbase.Hex())

	hundredEth := new(big.Int).Mul(big.NewInt(100), big.NewInt(params.Ether))
	t.Log("100eth ", hundredEth.Text(10))

	err = rpc.CallContext(context.TODO(), &raw, "eth_sendTransaction", map[string]string{
		"from":     coinbase.Hex(),
		"to":       accRelayer.Address.Hex(),
		"gas":      "0x76c0",                   // 30400
		"gasPrice": "0x9184e72a000",            // 10000000000000
		"value":    "0x" + hundredEth.Text(16), // 2441406250
		"data":     "",
	})
	assert.Nil(t, err)

	rpc.Close()
	client, err := NewWeb3Client("http://localhost:8545", ks, &accRelayer)
	assert.Nil(t, err)
	client.MaxGasPrice = 0

	im, err := NewIdentityManager(client, &TestContractStore{t}, nil, nil)
	assert.Nil(t, err)

	assert.Nil(t, im.Initialize())

	id := Identity{
		Operational: common.Address{},
		Relayer:     accRelayer.Address,
		Recovery:    common.Address{},
		Revoke:      common.Address{},
		Impl:        *im.ImplAddr(),
	}
	idaddr, err := im.Deploy(&id)
	assert.Nil(t, err)

	t.Log("IDADDR ", idaddr.Hex())

	assert.Nil(t, im.Ping(&id))

	assert.False(t, true)

}
