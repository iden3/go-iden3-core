package zk

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/iden3/go-circom-prover-verifier/parsers"
	zktypes "github.com/iden3/go-circom-prover-verifier/types"
	"github.com/iden3/go-iden3-core/common"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// G1ToBigInts transforms a `*bn256.G1` into `*big.Int` format, to be used for
// example in snarkjs solidity verifiers.
func G1ToBigInts(g1 *bn256.G1) [2]*big.Int {
	numBytes := 256 / 8
	bs := g1.Marshal()
	x := new(big.Int).SetBytes(bs[:numBytes])
	y := new(big.Int).SetBytes(bs[numBytes:])
	return [2]*big.Int{x, y}
}

// G2ToBigInts transforms a `*bn256.G2` into `*big.Int` format, to be used for
// example in snarkjs solidity verifiers.
func G2ToBigInts(g2 *bn256.G2) [2][2]*big.Int {
	numBytes := 256 / 8
	bs := g2.Marshal()
	xx := new(big.Int).SetBytes(bs[0*numBytes : 1*numBytes])
	xy := new(big.Int).SetBytes(bs[1*numBytes : 2*numBytes])
	yx := new(big.Int).SetBytes(bs[2*numBytes : 3*numBytes])
	yy := new(big.Int).SetBytes(bs[3*numBytes : 4*numBytes])
	// return [2][2]*big.Int{[2]*big.Int{xy, xx}, [2]*big.Int{yy, yx}}
	return [2][2]*big.Int{[2]*big.Int{xx, xy}, [2]*big.Int{yx, yy}}
}

// ProofToBigInts transforms a zkp (that uses `*bn256.G1` and `*bn256.G2`) into
// `*big.Int` format, to be used for example in snarkjs solidity verifiers.
func ProofToBigInts(proof *zktypes.Proof) (a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) {
	a = G1ToBigInts(proof.A)
	b = G2ToBigInts(proof.B)
	c = G1ToBigInts(proof.C)
	return a, b, c
}

// PrintProof prints the zkp in JSON pretty format.
func PrintProof(proof *zktypes.Proof) {
	proofA, proofB, proofC := ProofToBigInts(proof)
	fmt.Printf(
		`    "a": ["%v",
	    "%v"],
`,
		proofA[0], proofA[1])
	fmt.Printf(
		`    "b": [
           ["%v",
            "%v"],
           ["%v",
            "%v"]],
`,
		proofB[0][0], proofB[0][1], proofB[1][0], proofB[1][1])
	fmt.Printf(
		`    "c": ["%v",
	    "%v"]
`,
		proofC[0], proofC[1])
}

func download(url, filename string) error {
	dialTimeout := func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, time.Duration(2*time.Second))
	}
	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

// calcHash uses sha256
func calcHash(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

// checkHash uses sha256
func checkHash(filename, hashStr string) error {
	hash, err := hex.DecodeString(hashStr)
	if err != nil {
		return err
	}
	h, err := calcHash(filename)
	if err != nil {
		return err
	}
	if !bytes.Equal(h, hash) {
		fmt.Printf("\"%s\": \"%s\",\n", path.Base(filename), hex.EncodeToString(h))
		return fmt.Errorf("hash mismatch: expected %v but got %v", hashStr, hex.EncodeToString(h))
	}
	return nil
}

// ZkFilesHashes are the sha256 hash in hex of the zk files
type ZkFilesHashes struct {
	ProvingKey      string
	VerificationKey string
	WitnessCalcWASM string
}

// ZkFiles allows convenient access to the files required for zk proving and verifying.
type ZkFiles struct {
	Url                 string
	Path                string
	hashes              ZkFilesHashes
	cacheProvingKey     bool
	pathProvingKey      string
	provingKey          *zktypes.Pk
	pathVerificationKey string
	verificationKey     *zktypes.Vk
	pathWitnessCalcWASM string
	witnessCalcWASM     []byte
}

// NewZkFiles creates a new ZkFiles that will try to use the zk files from
// `path` checking that the `hashes` match with the files.  If the files don't
// exist, they are downloaded into `path` from `url`.  The proving key can be
// quite big: setting `cacheProvingKey` to false will make the ZkFiles not
// keep it in memory after requesting it, parsing it from disk every time it is
// required.  The rest of the files are always cached.
func NewZkFiles(url, path string, hashes ZkFilesHashes, cacheProvingKey bool) *ZkFiles {
	return &ZkFiles{
		Url:             url,
		Path:            path,
		hashes:          hashes,
		cacheProvingKey: cacheProvingKey,
	}
}

func (z *ZkFiles) insecureDownload(basename string) error {
	if err := os.MkdirAll(z.Path, 0700); err != nil {
		return err
	}
	filename := path.Join(z.Path, basename)
	_, err := os.Stat(filename)
	if err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	url := fmt.Sprintf("%s/%s", z.Url, basename)
	log.WithField("filename", filename).WithField("url", url).Debug("Downloading zk file")
	if err := download(url, filename); err != nil {
		return err
	}
	return nil
}

// InsecureDownloadAll downloads all the zk files but doesn't check the hashes.
func (z *ZkFiles) InsecureDownloadAll() error {
	for _, basename := range []string{"proving_key.json", "verification_key.json", "circuit.wasm"} {
		if err := z.insecureDownload(basename); err != nil {
			return err
		}
	}
	return nil
}

// InsecureCalcHashes calculates the hashes of the zkfiles without checking them.
func (z *ZkFiles) InsecureCalcHashes() (*ZkFilesHashes, error) {
	var hashes [3][]byte
	for i, basename := range []string{"proving_key.json", "verification_key.json", "circuit.wasm"} {
		filename := path.Join(z.Path, basename)
		h, err := calcHash(filename)
		if err != nil {
			return nil, err
		}
		hashes[i] = h
	}
	return &ZkFilesHashes{
		ProvingKey:      hex.EncodeToString(hashes[0]),
		VerificationKey: hex.EncodeToString(hashes[1]),
		WitnessCalcWASM: hex.EncodeToString(hashes[2]),
	}, nil
}

// DebugDownloadPrintHashes is a helper function that downloads all the zk
// files in a temporary directory, calculates their hashes, and prints the code
// of the `ZkFilesHashes` with the calculated hashes, ready to be pasted in
// real code.
func (z *ZkFiles) DebugDownloadPrintHashes() error {
	dir, err := ioutil.TempDir("", "zkfiles")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up
	z0 := NewZkFiles(z.Url, dir, ZkFilesHashes{}, false)
	if err := z0.InsecureDownloadAll(); err != nil {
		return nil
	}
	hashes, err := z0.InsecureCalcHashes()
	if err != nil {
		return err
	}
	s := fmt.Sprintf("%#v", hashes)
	s = strings.ReplaceAll(s, "{", "{\n\t")
	s = strings.ReplaceAll(s, ", ", ",\n\t")
	s = strings.ReplaceAll(s, "}", ",\n}")
	fmt.Println(s)
	return nil
}

func (z *ZkFiles) downloadCheckFile(basename, hash string) error {
	filename := path.Join(z.Path, basename)
	_, err := os.Stat(filename)
	if err == nil {
		if err := checkHash(filename, hash); err != nil {
			return err
		}
		// log.WithField("filename", filename).Debug("Skipping downloading zk file")
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	url := fmt.Sprintf("%s/%s", z.Url, basename)
	log.WithField("filename", filename).WithField("url", url).Debug("Downloading zk file")
	if err := download(url, filename); err != nil {
		return err
	}
	if err := checkHash(filename, hash); err != nil {
		return err
	}
	return nil
}

// DownloadProvingKey downloads the ProvingKey and checks its hash.
func (z *ZkFiles) DownloadProvingKey() error {
	if err := os.MkdirAll(z.Path, 0700); err != nil {
		return err
	}
	basename := "proving_key.json"
	if err := z.downloadCheckFile(basename, z.hashes.ProvingKey); err != nil {
		return err
	}
	z.pathProvingKey = path.Join(z.Path, basename)
	return nil
}

// DownloadVerificationKey downloads the VerificationKey and checks its hash.
func (z *ZkFiles) DownloadVerificationKey() error {
	if err := os.MkdirAll(z.Path, 0700); err != nil {
		return err
	}
	basename := "verification_key.json"
	if err := z.downloadCheckFile(basename, z.hashes.VerificationKey); err != nil {
		return err
	}
	z.pathVerificationKey = path.Join(z.Path, basename)
	return nil
}

// DownloadWitnessCalcWASM downloads the WitnessCalcWASM and checks its hash.
func (z *ZkFiles) DownloadWitnessCalcWASM() error {
	if err := os.MkdirAll(z.Path, 0700); err != nil {
		return err
	}
	basename := "circuit.wasm"
	if err := z.downloadCheckFile(basename, z.hashes.WitnessCalcWASM); err != nil {
		return err
	}
	z.pathWitnessCalcWASM = path.Join(z.Path, basename)
	return nil
}

// DownloadAll downloads all the zk files and checks their hashes.
func (z *ZkFiles) DownloadAll() error {
	if err := z.DownloadProvingKey(); err != nil {
		return err
	}
	if err := z.DownloadVerificationKey(); err != nil {
		return err
	}
	if err := z.DownloadWitnessCalcWASM(); err != nil {
		return err
	}
	return nil
}

func (z *ZkFiles) loadProvingKey() (*zktypes.Pk, error) {
	provingKeyJson, err := ioutil.ReadFile(z.pathProvingKey)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	pk, err := parsers.ParsePk(provingKeyJson)
	if err != nil {
		return nil, err
	}
	log.WithField("elapsed", time.Since(start)).Debug("Parsed proving key")
	return pk, nil
}

// LoadProvingKey loads the ProvingKey, downloading it if necessary.
func (z *ZkFiles) LoadProvingKey() error {
	if z.provingKey != nil {
		// log.Debug("zkfiles: proving key already loaded")
		return nil
	}
	if z.pathProvingKey == "" {
		if err := z.DownloadProvingKey(); err != nil {
			return err
		}
	}
	if z.cacheProvingKey {
		if pk, err := z.loadProvingKey(); err != nil {
			return err
		} else {
			z.provingKey = pk
		}
	}
	return nil
}

// LoadVerificationKey loads the VerificationKey, downloading it if necessary.
func (z *ZkFiles) LoadVerificationKey() error {
	if z.verificationKey != nil {
		// log.Debug("zkfiles: verification key already loaded")
		return nil
	}
	if z.pathVerificationKey == "" {
		if err := z.DownloadVerificationKey(); err != nil {
			return err
		}
	}
	vkJSON, err := ioutil.ReadFile(z.pathVerificationKey)
	if err != nil {
		return err
	}
	vk, err := parsers.ParseVk(vkJSON)
	if err != nil {
		return err
	}
	z.verificationKey = vk
	return nil
}

// LoadWitnessCalcWASM loads the WitnessCalcWASM, downloading it if necessary.
func (z *ZkFiles) LoadWitnessCalcWASM() error {
	if z.witnessCalcWASM != nil {
		// log.Debug("zkfiles: witnessCalc WASM already loaded")
		return nil
	}
	if z.pathWitnessCalcWASM == "" {
		if err := z.DownloadWitnessCalcWASM(); err != nil {
			return err
		}
	}
	wasmBytes, err := ioutil.ReadFile(z.pathWitnessCalcWASM)
	if err != nil {
		return err
	}
	z.witnessCalcWASM = wasmBytes
	return nil
}

// LoadAll loads all the zk files, downloading them if necessary.
func (z *ZkFiles) LoadAll() error {
	if err := z.LoadProvingKey(); err != nil {
		return err
	}
	if err := z.LoadVerificationKey(); err != nil {
		return err
	}
	if err := z.LoadWitnessCalcWASM(); err != nil {
		return err
	}
	return nil
}

// ProvingKey returns the ProvingKey, downloading and loading it if necessary.
func (z *ZkFiles) ProvingKey() (*zktypes.Pk, error) {
	if err := z.LoadProvingKey(); err != nil {
		return nil, err
	}
	var pk *zktypes.Pk
	if !z.cacheProvingKey {
		var err error
		pk, err = z.loadProvingKey()
		if err != nil {
			return nil, err
		}
	} else {
		pk = z.provingKey
	}
	return pk, nil
}

// VerificationKey returns the VerificationKey, downloading and loading it if necessary.
func (z *ZkFiles) VerificationKey() (*zktypes.Vk, error) {
	if err := z.LoadVerificationKey(); err != nil {
		return nil, err
	}
	return z.verificationKey, nil
}

// WitnessCalcWASM returns the WitnessCalcWASM byte slice, downloading and loading it if necessary.
func (z *ZkFiles) WitnessCalcWASM() ([]byte, error) {
	if err := z.LoadWitnessCalcWASM(); err != nil {
		return nil, err
	}
	return z.witnessCalcWASM, nil
}

// InputsToMapStrings transforms the input signals map from *big.Int type (as
// used in witnesscalc) to quoted strings (as used in JSON encoding).
func InputsToMapStrings(inputs interface{}) (map[string]interface{}, error) {
	var inputsMap map[string]interface{}
	if err := mapstructure.Decode(inputs, &inputsMap); err != nil {
		return nil, err
	}
	inputsStrings := make(map[string]interface{})
	for key, value := range inputsMap {
		switch v := value.(type) {
		case *big.Int:
			inputsStrings[key] = v.String()
		case []*big.Int:
			vs := make([]string, len(v))
			for i, v := range v {
				vs[i] = v.String()
			}
			inputsStrings[key] = vs
		default:
			panic(fmt.Sprintf("Type: %T", value))
		}
	}
	return inputsStrings, nil
}

// ZkProofOut is the output of calculating a zkp.
type ZkProofOut struct {
	Proof      zktypes.Proof
	PubSignals []*big.Int
}

// PubSignals is a helper wrapper type over []*big.Int that is JSON friendly.
type PubSignals []*big.Int

func (p PubSignals) MarshalJSON() ([]byte, error) {
	aux := make([]string, len(p))
	for i, v := range p {
		// We use LittleEndian here!
		aux[i] = hex.EncodeToString(common.SwapEndianness(v.Bytes()))
	}
	return json.Marshal(aux)
}

func (p *PubSignals) UnmarshalJSON(data []byte) error {
	var aux []string
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	vs := make([]*big.Int, len(aux))
	for i, v := range aux {
		bs, err := hex.DecodeString(v)
		if err != nil {
			return err
		}
		vs[i] = new(big.Int).SetBytes(common.SwapEndianness(bs))
	}
	*p = vs
	return nil
}
