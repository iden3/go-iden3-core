package zk

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	witnesscalc "github.com/iden3/go-circom-witnesscalc"
	"github.com/iden3/go-iden3-crypto/babyjub"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
	"github.com/iden3/go-wasm3"

	log "github.com/sirupsen/logrus"
)

func CalculateWitnessBinWASM(wasmBytes []byte, inputs []witnesscalc.Input) ([]*big.Int, error) {
	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   64 * 1024,
	})
	defer runtime.Destroy()

	module, err := runtime.ParseModule(wasmBytes)
	if err != nil {
		return nil, err
	}
	module, err = runtime.LoadModule(module)
	if err != nil {
		return nil, err
	}

	witnessCalculator, err := witnesscalc.NewWitnessCalculator(runtime, module)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	witness, err := witnessCalculator.CalculateWitness(inputs, true)
	if err != nil {
		return nil, err
	}
	log.WithField("elapsed", time.Now().Sub(start)).Debug("Witness calculated")

	sum := new(big.Int)
	m := new(big.Int).SetUint64(0x10000)
	for _, w := range witness {
		sum.Add(sum, w)
		sum.Mod(sum, m)
	}
	fmt.Printf(">>> SUM: %v\n", sum)

	return witness, err
}

func CalculateWitness(wasmFilePath string, inputs []witnesscalc.Input) ([]*big.Int, error) {
	wasmBytes, err := ioutil.ReadFile(wasmFilePath)
	if err != nil {
		return nil, err
	}
	return CalculateWitnessBinWASM(wasmBytes, inputs)
}

func PrivateKeyToBigInt(k *babyjub.PrivateKey) *big.Int {
	sBuf := babyjub.Blake512(k[:])
	// sBuf, _ := hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	sBuf32 := [32]byte{}
	copy(sBuf32[:], sBuf[:32])
	pruneBuffer(&sBuf32)
	fmt.Printf("### %v\n", hex.EncodeToString(sBuf32[:]))
	s := new(big.Int)
	cryptoUtils.SetBigIntFromLEBytes(s, sBuf32[:])
	s.Rsh(s, 3)
	return s
}

func pruneBuffer(buf *[32]byte) *[32]byte {
	buf[0] = buf[0] & 0xF8
	buf[31] = buf[31] & 0x7F
	buf[31] = buf[31] | 0x40
	return buf
}
