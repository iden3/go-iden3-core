package issuer

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

func GetIdenStateZKFiles(urlBase, downloadPath string) error {
	// downloadPath := "/tmp/iden3/idenstatezk"
	filenamesHash := map[string]string{
		"circuit.wasm":          "8eafd9314c4d2664a23bf98a4f42cd0c29984960ae3544747ba5fbd60905c41f",
		"proving_key.json":      "2c72fceb10323d8b274dbd7649a63c1b6a11fff3a1e4cd7f5ec12516f32ec452",
		"verification_key.json": "473952ff80aef85403005eb12d1e78a3f66b1cc11e7bd55d6bfe94e0b5577640",
	}
	if err := os.MkdirAll(downloadPath, 0700); err != nil {
		return err
	}
	for basename, hash := range filenamesHash {
		filename := path.Join(downloadPath, basename)
		_, err := os.Stat(filename)
		if err == nil {
			if err := checkHash(filename, hash); err != nil {
				return err
			}
			log.WithField("filename", filename).Debug("Skipping downloading zk file")
			continue
		} else if !os.IsNotExist(err) {
			return err
		}
		url := fmt.Sprintf("%s/%s", urlBase, basename)
		log.WithField("filename", filename).WithField("url", url).Debug("Downloading zk file")
		if err := download(url, filename); err != nil {
			return err
		}
		if err := checkHash(filename, hash); err != nil {
			return err
		}
	}
	return nil
}

func checkHash(filename, hashStr string) error {
	hash, err := hex.DecodeString(hashStr)
	if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return err
	}
	h := hasher.Sum(nil)
	if bytes.Compare(h, hash) != 0 {
		fmt.Printf("\"%s\": \"%s\",\n", path.Base(filename), hex.EncodeToString(h))
		return fmt.Errorf("hash mismatch: expected %v but got %v", hashStr, hex.EncodeToString(h))
	}
	return nil
}

func download(url, filename string) error {
	resp, err := http.Get(url)
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
