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

func GetIdenStateZKFiles(urlBase string) error {
	downloadPath := "/tmp/iden3/idenstatezk"
	filenamesHash := map[string]string{
		"circuit.wasm":          "8eafd9314c4d2664a23bf98a4f42cd0c29984960ae3544747ba5fbd60905c41f",
		"proving_key.json":      "972373336851ef6c51366db4795c8821140ee95e1340bb4d01d35fb8e38d0116",
		"verification_key.json": "b6a43b38be6b855c0be06b5b5b1b871fa4a774eed6a3bfd81df99c51147449b1",
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
