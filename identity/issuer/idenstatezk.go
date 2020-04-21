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
		"circuit.wasm":          "7c3958904d30f949187e070a606f87a2483bfad8321818a16b3d5d35739ba2d0",
		"proving_key.json":      "11a04fe1e6566e3e8a2f00032851b2e0999011324d216a0517b006fd0c8695fb",
		"verification_key.json": "ef24eabbb0c172ede61f58c737e07876f70e2fd17e95daa06e63abb80c620883",
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
