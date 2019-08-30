package keystore

import (
	"crypto/rand"
	"time"

	// "encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"runtime"
	"sync"

	"github.com/gofrs/flock"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/utils"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

// Constants taken from
// https://github.com/ethereum/go-ethereum/blob/master/accounts/keystore/passphrase.go
const (
	// StandardScryptN is the N parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptN = 1 << 18

	// StandardScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptP = 1

	// LightScryptN is the N parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	LightScryptN = 1 << 12

	// LightScryptP is the P parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	LightScryptP = 6

	scryptR     = 8
	scryptDKLen = 32
)

// prefixes for msg to be signed
type PrefixType []byte

var (
	// PrefixMinorUpdate is for signatures related to update the root of an identity as minor update
	PrefixMinorUpdate = []byte("minorupdate")
)

// KeyStoreParams are the Key Store parameters
type KeyStoreParams struct {
	ScryptN int
	ScryptP int
}

// LightKeyStoreParams are parameters for fast key derivation
var LightKeyStoreParams = KeyStoreParams{
	ScryptN: LightScryptN,
	ScryptP: LightScryptP,
}

// StandardKeyStoreParams are parameters for very secure derivation
var StandardKeyStoreParams = KeyStoreParams{
	ScryptN: StandardScryptN,
	ScryptP: StandardScryptP,
}

// EncryptedData contains the key derivation parameters and encryption
// parameters with the encrypted data.
type EncryptedData struct {
	Salt          common3.Hex
	ScryptN       int
	ScryptP       int
	Nonce         common3.Hex
	EncryptedData common3.Hex
}

// EncryptedData encrypts data with a key derived from pass
func EncryptData(data, pass []byte, scryptN, scryptP int) (*EncryptedData, error) {
	var salt [32]byte
	if _, err := io.ReadFull(rand.Reader, salt[:]); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	derivedKey, err := scrypt.Key(pass, salt[:], scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	var key [32]byte
	copy(key[:], derivedKey)
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	var encryptedData []byte
	encryptedData = secretbox.Seal(encryptedData, data, &nonce, &key)

	return &EncryptedData{
		Salt:          common3.Hex(salt[:]),
		ScryptN:       scryptN,
		ScryptP:       scryptP,
		Nonce:         common3.Hex(nonce[:]),
		EncryptedData: common3.Hex(encryptedData),
	}, nil
}

// DecryptData decrypts the encData with the key derived from pass.
func DecryptData(encData *EncryptedData, pass []byte) ([]byte, error) {
	derivedKey, err := scrypt.Key(pass, encData.Salt[:],
		encData.ScryptN, scryptR, encData.ScryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	var key [32]byte
	copy(key[:], derivedKey)
	var nonce [24]byte
	copy(nonce[:], encData.Nonce)
	var data []byte
	data, ok := secretbox.Open(data, encData.EncryptedData, &nonce, &key)
	if !ok {
		return nil, fmt.Errorf("Invalid encrypted data")
	}
	return data, nil
}

// KeysStored is the datastructure of stored keys in the storage.
type KeysStored map[babyjub.PublicKeyComp]EncryptedData

// Storage is an interface for a storage container.
type Storage interface {
	Read() ([]byte, error)
	Write(data []byte) error
	TryLock() (bool, error)
	Unlock() error
}

// FileStorage is a storage backed by a file.
type FileStorage struct {
	path string
	lock *flock.Flock
}

// NewFileStorage returns a new FileStorage backed by a file in path.
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path, lock: flock.New(path + ".lock")}
}

// Read reads the file contents.
func (fs *FileStorage) Read() ([]byte, error) {
	return ioutil.ReadFile(fs.path)
}

// Write writes the data to the file.
func (fs *FileStorage) Write(data []byte) error {
	return ioutil.WriteFile(fs.path, data, 0600)
}

// TryLocks the storage file with a .lock file.
func (fs *FileStorage) TryLock() (bool, error) {
	return fs.lock.TryLock()
}

// Unlocks the storage file and removes the .lock file.
func (fs *FileStorage) Unlock() error {
	if err := fs.lock.Unlock(); err != nil {
		return err
	}
	return os.Remove(fs.path + ".lock")
}

// MemStorage is a storage backed by a slice.
type MemStorage []byte

// Read reads the slice contents.
func (ms *MemStorage) Read() ([]byte, error) {
	return []byte(*ms), nil
}

// Write copies the data to the slice.
func (ms *MemStorage) Write(data []byte) error {
	*ms = data
	return nil
}

// TryLock does nothing.
func (ms *MemStorage) TryLock() (bool, error) { return true, nil }

// Unlock does nothing.
func (ms *MemStorage) Unlock() error { return nil }

// KeyStore is the object used to access create keys and sign with them.
type KeyStore struct {
	storage       Storage
	params        KeyStoreParams
	encryptedKeys KeysStored
	cache         map[babyjub.PublicKeyComp]*babyjub.PrivateKey
	rw            sync.RWMutex
}

// NewKeyStore creates a new key store or opens it if it already exists.
func NewKeyStore(storage Storage, params KeyStoreParams) (*KeyStore, error) {
	if ok, err := storage.TryLock(); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("Unable to acquire storage lock")
	}
	log.Info("BabyJub KeyStore storage locked successfully")
	encryptedKeysJSON, err := storage.Read()
	if os.IsNotExist(err) {
		encryptedKeysJSON = []byte{}
	} else if err != nil {
		storage.Unlock()
		return nil, err
	}
	var encryptedKeys KeysStored
	if len(encryptedKeysJSON) == 0 {
		encryptedKeys = make(map[babyjub.PublicKeyComp]EncryptedData)
	} else {
		if err := json.Unmarshal(encryptedKeysJSON, &encryptedKeys); err != nil {
			storage.Unlock()
			return nil, err
		}
	}
	ks := &KeyStore{
		storage:       storage,
		params:        params,
		encryptedKeys: encryptedKeys,
		cache:         make(map[babyjub.PublicKeyComp]*babyjub.PrivateKey),
	}
	runtime.SetFinalizer(ks, func(ks *KeyStore) {
		// When there are no more references to the key store, clear
		// the secret keys in the cache and unlock the locked storage.
		ks.Close()
	})
	return ks, nil
}

func (ks *KeyStore) Close() {
	zero := [32]byte{}
	for _, sk := range ks.cache {
		copy(sk[:], zero[:])
	}
	err := ks.storage.Unlock()
	if err != nil {
		log.Error("Failed unlocking BabyJub KeyStore storage ", err)
	} else {
		log.Info("BabyJub KeyStore storage unlocked")
	}
}

// Keys returns the compressed public keys of the key storage.
func (ks *KeyStore) Keys() []babyjub.PublicKeyComp {
	ks.rw.RLock()
	defer ks.rw.RUnlock()
	keys := make([]babyjub.PublicKeyComp, 0, len(ks.encryptedKeys))
	for pk, _ := range ks.encryptedKeys {
		keys = append(keys, pk)
	}
	return keys
}

// NewKey creates a new key in the key store encrypted with pass.
func (ks *KeyStore) NewKey(pass []byte) (*babyjub.PublicKeyComp, error) {
	sk := babyjub.NewRandPrivKey()
	return ks.ImportKey(sk, pass)
}

// ImportKey imports a secret key into the storage and encrypts it with pass.
func (ks *KeyStore) ImportKey(sk babyjub.PrivateKey, pass []byte) (*babyjub.PublicKeyComp, error) {
	ks.rw.Lock()
	defer ks.rw.Unlock()
	encryptedKey, err := EncryptData(sk[:], pass, ks.params.ScryptN, ks.params.ScryptP)
	if err != nil {
		return nil, err
	}
	pk := sk.Public()
	pubComp := pk.Compress()
	ks.encryptedKeys[pubComp] = *encryptedKey
	encryptedKeysJSON, err := json.Marshal(ks.encryptedKeys)
	if err != nil {
		return nil, err
	}
	if err := ks.storage.Write(encryptedKeysJSON); err != nil {
		return nil, err
	}
	return &pubComp, nil
}

func (ks *KeyStore) ExportKey(pk *babyjub.PublicKeyComp, pass []byte) (*babyjub.PrivateKey, error) {
	if err := ks.UnlockKey(pk, pass); err != nil {
		return nil, err
	}
	return ks.cache[*pk], nil
}

// UnlockKey decrypts the key corresponding to the public key pk and loads it
// into the cache.
func (ks *KeyStore) UnlockKey(pk *babyjub.PublicKeyComp, pass []byte) error {
	ks.rw.Lock()
	defer ks.rw.Unlock()
	encryptedKey, ok := ks.encryptedKeys[*pk]
	if !ok {
		return fmt.Errorf("Public key not found in the key store")
	}
	skBuf, err := DecryptData(&encryptedKey, pass)
	if err != nil {
		return err
	}
	var sk babyjub.PrivateKey
	copy(sk[:], skBuf)
	ks.cache[*pk] = &sk
	return nil
}

// SignElem uses the key corresponding to the public key pk to sign the field
// element msg.
func (ks *KeyStore) SignElem(pk *babyjub.PublicKeyComp, msg *big.Int) (*babyjub.SignatureComp, error) {
	ks.rw.RLock()
	defer ks.rw.RUnlock()
	sk, ok := ks.cache[*pk]
	if !ok {
		return nil, fmt.Errorf("Public key not found in the cache.  Is it unlocked?")
	}
	sig := sk.SignMimc7(msg)
	sigComp := sig.Compress()
	return &sigComp, nil
}

// Sign uses the key corresponding to the public key pk to sign the mimc7 hash
// of the [prefix | date | msg] byte slice.
func (ks *KeyStore) Sign(pk *babyjub.PublicKeyComp, prefix PrefixType, rawMsg []byte) (*babyjub.SignatureComp, int64, error) {
	date := time.Now()
	msg := append(prefix, utils.Uint64ToEthBytes(uint64(date.Unix()))...)
	msg = append(msg, rawMsg...)
	sig, err := ks.SignRaw(pk, msg)
	return sig, date.Unix(), err
}

// SignRaw uses the key corresponding to the public key pk to sign the mimc7/poseidon hash
// of the msg byte slice.
func (ks *KeyStore) SignRaw(pk *babyjub.PublicKeyComp, msg []byte) (*babyjub.SignatureComp, error) {
	// h, err := mimc7.HashBytes(msg)
	h, err := poseidon.HashBytes(msg)
	if err != nil {
		return nil, err
	}
	return ks.SignElem(pk, h)
}

// VerifySignatureElem verifies that the signature sigComp of the field element
// msg was signed with the public key pkComp.
func VerifySignatureElem(pkComp *babyjub.PublicKeyComp, msg *big.Int, sigComp *babyjub.SignatureComp) (bool, error) {
	pkPoint, err := babyjub.NewPoint().Decompress(*pkComp)
	if err != nil {
		return false, err
	}
	sig, err := new(babyjub.Signature).Decompress(*sigComp)
	if err != nil {
		return false, err
	}
	pk := babyjub.PublicKey(*pkPoint)
	return pk.VerifyMimc7(msg, sig), nil
}

// VerifySignature verifies that the signature sigComp of the poseidon hash of
// the [prefix | date | msg] byte slice was signed with the public key pkComp.
func VerifySignature(pkComp *babyjub.PublicKeyComp, sigComp *babyjub.SignatureComp, prefix PrefixType, date int64, rawMsg []byte) (bool, error) {
	msg := append(prefix, utils.Uint64ToEthBytes(uint64(date))...)
	msg = append(msg, rawMsg...)
	return VerifySignatureRaw(pkComp, sigComp, msg)
}

// VerifySignatureRaw verifies that the signature sigComp of the poseidon hash of
// the msg byte slice was signed with the public key pkComp.
func VerifySignatureRaw(pkComp *babyjub.PublicKeyComp, sigComp *babyjub.SignatureComp, msg []byte) (bool, error) {
	h, err := poseidon.HashBytes(msg)
	if err != nil {
		return false, err
	}
	return VerifySignatureElem(pkComp, h, sigComp)
}
