package zk

import (
	"io/ioutil"
	"math/big"
	"time"

	witnesscalc "github.com/iden3/go-circom-witnesscalc"
	"github.com/iden3/go-wasm3"

	log "github.com/sirupsen/logrus"
)

func CalculateWitnessBinWASM(wasmBytes []byte, inputs map[string]interface{}) ([]*big.Int, error) {
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

	return witness, err
}

func CalculateWitness(wasmFilePath string, inputs map[string]interface{}) ([]*big.Int, error) {
	wasmBytes, err := ioutil.ReadFile(wasmFilePath)
	if err != nil {
		return nil, err
	}
	return CalculateWitnessBinWASM(wasmBytes, inputs)
}
