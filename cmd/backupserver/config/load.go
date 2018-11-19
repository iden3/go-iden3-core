package config

import (
	"os"

	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/services/backupsrv"
	"github.com/iden3/go-iden3/services/mongosrv"

	log "github.com/sirupsen/logrus"
)

func assert(msg string, err error) {
	if err != nil {
		log.Error(msg, " ", err.Error())
		os.Exit(1)
	}
}

func LoadStorage() db.Storage {
	// Open database
	storage, err := db.NewLevelDbStorage(C.Storage.Path, false)
	assert("Cannot open storage", err)
	log.WithField("path", C.Storage.Path).Info("Storage opened")
	return storage
}
func LoadMongoService() mongosrv.Service {
	collectionsArray := []string{"data"}
	mongoservice, err := mongosrv.New(C.Mongodb.Url, C.Mongodb.Database, collectionsArray)
	assert("Cannot open mongodb storage", err)
	log.WithField("path", C.Mongodb.Url).Info("Mongodb storage opened")
	return mongoservice
}

func LoadBackupService(sto mongosrv.Service) backupsrv.Service {
	return backupsrv.New(sto)
}
