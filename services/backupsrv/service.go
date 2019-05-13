package backupsrv

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/mongosrv"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type Service interface {
	// BACKUP SERVICE
	Register(user User) error
	BackupUpload(user User, backupPacket BackupPacket) error
	BackupDownload(user User) (BackupPacket, error)

	// SYNCHRONIZATION SERVICE
	GetPoWDifficulty() int
	GetLastVersion(idaddr core.ID) (uint64, error)
	//Save(idaddr common.Address, saveBackupMsg BackupDataMsg) (uint64, error)
	RecoverAll(idaddr core.ID) ([]BackupDataMsg, error)
	RecoverSinceVersion(idaddr core.ID, version uint64) ([]BackupDataMsg, error)
	RecoverByType(idaddr core.ID, dataType string) ([]BackupDataMsg, error)
	RecoverSinceVersionByType(idaddr core.ID, version uint64, dataType string) ([]BackupDataMsg, error)
}

type ServiceImpl struct {
	mongodb       mongosrv.Service
	powDifficulty int
}

func New(mongoservice mongosrv.Service, powDiff int) *ServiceImpl {
	return &ServiceImpl{mongoservice, powDiff}
}

// Register adds a new user into the db if it already not exists
func (bs *ServiceImpl) Register(user User) error {

	existingUser := User{}
	err := bs.mongodb.GetCollections()["users"].Find(bson.M{"username": user.Username}).One(&existingUser)
	if err == nil {
		// already existing user
		return errors.New("already existing user")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashStr := string(hash)
	user.Password = hashStr

	err = bs.mongodb.GetCollections()["users"].Insert(user)
	if err != nil {
		return err
	}

	return nil
}

func (bs *ServiceImpl) BackupUpload(user User, backupPacket BackupPacket) error {
	dbUser := User{}
	err := bs.mongodb.GetCollections()["users"].Find(bson.M{"username": user.Username}).One(&dbUser)
	if err != nil {
		return errors.New("user not exists")
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return errors.New("error with the password")
	}

	fmt.Println(backupPacket)
	existingBackupPacket := BackupPacket{}
	err = bs.mongodb.GetCollections()["backup"].Find(bson.M{"username": user.Username}).One(&existingBackupPacket)
	if err != nil {
		// first backup of the user
		err = bs.mongodb.GetCollections()["backup"].Insert(backupPacket)
	} else {
		// not first backup of the user
		err = bs.mongodb.GetCollections()["backup"].Update(bson.M{"username": user.Username}, backupPacket)
	}

	return err
}

func (bs *ServiceImpl) BackupDownload(user User) (BackupPacket, error) {
	dbUser := User{}
	err := bs.mongodb.GetCollections()["users"].Find(bson.M{"username": user.Username}).One(&dbUser)
	if err != nil {
		return BackupPacket{}, errors.New("user not exists")
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return BackupPacket{}, errors.New("error with the password")
	}

	backupPacket := BackupPacket{}
	err = bs.mongodb.GetCollections()["backup"].Find(bson.M{"username": user.Username}).One(&backupPacket)
	if err != nil {
		return BackupPacket{}, err
	}

	return backupPacket, nil
}

// GetPoWDifficulty returns the configured Proof-of-Work difficulty, setted in the config file of the backupserver
func (bs *ServiceImpl) GetPoWDifficulty() int {
	return bs.powDifficulty
}

func (bs *ServiceImpl) GetLastVersion(idaddr core.ID) (uint64, error) {
	var currVer BackupDataMsg
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.String()), "ksign": "currentversion"}).One(&currVer)
	return currVer.Version, err
}

// TODO: Update to new claim ksign secp256k1
// Save verifies the proofs for auth, and stores the data packet in the database
//func (bs *ServiceImpl) Save(idaddr common.Address, m BackupDataMsg) (uint64, error) {
//	// check PoW
//	b, err := json.Marshal(m)
//	if err != nil {
//		return 0, err
//	}
//	hash := utils.HashBytes(b)
//	if !utils.CheckPoW(hash, bs.GetPoWDifficulty()) {
//		return 0, errors.New("PoW not passed")
//	}
//
//	// check ksignClaim proof (in user identity tree and in the relay tree)
//	proofOfKSign, err := m.ProofKSignHex.Unhex()
//	if err != nil {
//		return 0, err
//	}
//	kSignComp := crypto.CompressPubkey(m.KSignPk)
//	relayAddr := common.HexToAddress(m.RelayAddr)
//	verified := claimsrv.CheckProofClaimUser(relayAddr, proofOfKSign, 140)
//	if !verified {
//		return 0, errors.New("ProofKSign can not be verified")
//	}
//
//	// check saveBackupMsg.KSign match with authorizedksign from the ProofKSign, Leaf[64:84] is where is placed the KeyToAuthorize (KSign authorized) in the Claim data
//	if !bytes.Equal(kSignComp, proofOfKSign.ClaimProof.Leaf[64:84]) {
//		return 0, errors.New("KSign not equal to the ProofKSign.ClaimProof.Leaf[KeyToAuthorize]")
//	}
//
//	// check idaddr match with setRootClaim from the proofOfKSign, Leaf[64:84] is where is placed the idaddr in the SetRootClaim
//	if !bytes.Equal(idaddr.Bytes(), proofOfKSign.SetRootClaimProof.Leaf[64:84]) {
//		return 0, errors.New("idaddr don't match with the idaddr from the ProofKSign.SetRootClaimProof.Leaf[EthAddr]")
//	}
//
//	// verify data signature
//	sigBytes, err := common3.HexDecode(m.DataSignature)
//	if err != nil {
//		return 0, err
//	}
//	if !utils.VerifySigEthMsg(kSign, sigBytes, []byte(m.Data)) {
//		return 0, errors.New("signature of the data can not be verified")
//	}
//
//	// check version (check that the current version is == lastversion+1)
//	var aux BackupDataMsg
//	err = bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "version": m.Version}).One(&aux)
//	if err == nil {
//		// if data exists, the given version is not valid
//		return m.Version, errors.New("given version not valid")
//	}
//
//	// store in database
//	err = bs.mongodb.GetCollections()["data"].Insert(m)
//	if err != nil {
//		return 0, err
//	}
//
//	// TODO store in leveldb instead of mongodb. key: idaddr+version, value: type+dataencrypted
//	// the currentVersion will be stored as key: idaddr+"currver", value: currver
//	currVer := BackupDataMsg{
//		IdAddrHex:       idaddr.Hex(),
//		Data:            "",
//		DataSignature:   "",
//		Type:            "",
//		KSign:           "currentversion",
//		ProofKSignHex: claimsrv.ProofClaimUserHex{},
//		RelayAddr:       "",
//		Version:         m.Version,
//		Nonce:           0,
//	}
//	err = bs.mongodb.GetCollections()["data"].Update(bson.M{"idaddrhex": strings.ToLower(idaddr.Hex()), "ksign": "currentversion"}, currVer)
//	return m.Version, nil
//}

// RecoverAll returns all the data packets stored by an idaddr
func (bs *ServiceImpl) RecoverAll(idaddr core.ID) ([]BackupDataMsg, error) {

	// TODO auth verifications
	// check ksignClaim proof (in user identity tree and in the relay tree)

	// get from database
	var dataBackups []BackupDataMsg
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.String())}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}

// RecoverSinceVersion returns all the data packets stored by an idaddr since after the version specified in the parameter
func (bs *ServiceImpl) RecoverSinceVersion(idaddr core.ID, version uint64) ([]BackupDataMsg, error) {
	color.Blue("version")
	fmt.Println(version)

	// TODO auth verifications

	// get from database
	var dataBackups []BackupDataMsg
	// get data with version greather or equal to 'version'
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.String()), "version": bson.M{"$gt": version}}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}

// RecoverByType returns all the data packets stored by an idaddr with the requested type
func (bs *ServiceImpl) RecoverByType(idaddr core.ID, dataType string) ([]BackupDataMsg, error) {

	// TODO auth verifications

	// get from database
	var dataBackups []BackupDataMsg
	// get data by type
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.String()), "type": dataType}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}

// RecoverSinceVersionByType returns all the data packets stored by an idaddr with the requested type since after the version specified in the parameter
func (bs *ServiceImpl) RecoverSinceVersionByType(idaddr core.ID, version uint64, dataType string) ([]BackupDataMsg, error) {

	// TODO auth verifications

	// get from database
	var dataBackups []BackupDataMsg
	// get data by type
	err := bs.mongodb.GetCollections()["data"].Find(bson.M{"idaddrhex": strings.ToLower(idaddr.String()), "version": bson.M{"$gt": version}, "type": dataType}).Limit(100).All(&dataBackups)
	if err != nil {
		return dataBackups, err
	}
	return dataBackups, nil
}
