package repository

import (
	"github.com/BagusAK95/amarta_test/internal/domain/borrower"
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"gorm.io/gorm"
)

type borrowerRepo struct {
	repository.BaseRepo[borrower.Borrower]
	writeConn *gorm.DB
	readConn  *gorm.DB
}

func NewBorrowerRepo(dbMaster *gorm.DB, dbSlave *gorm.DB) borrower.IBorrowerRepository {
	baseRepo := repository.NewBaseRepo[borrower.Borrower](dbMaster, dbSlave)

	return &borrowerRepo{
		BaseRepo:  *baseRepo,
		writeConn: dbMaster,
		readConn:  dbSlave,
	}
}
