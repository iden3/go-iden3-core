package backupsrv

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/db"
)

var ErrNotFound = errors.New("value not found")

type Service interface {
	SaveBackup(idaddr common.Address) error
	RecoverBackup(idaddr common.Address) error
}
type ServiceImpl struct {
	sto db.Storage
}

func New(sto db.Storage) *ServiceImpl {
	return &ServiceImpl{sto}
}

func (bs *ServiceImpl) SaveBackup(idaddr common.Address) error {
	// check data signature

	// check ksignClaim proof (in user identity tree and in the relay tree)

	// store in database

	return nil
}

func (bs *ServiceImpl) RecoverBackup(idaddr common.Address) error {
	// check data signature

	// check ksignClaim proof (in user identity tree and in the relay tree)

	// store in database

	return nil
}
