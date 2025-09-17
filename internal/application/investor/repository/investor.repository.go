package repository

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/BagusAK95/amarta_test/internal/domain/investor"
	"gorm.io/gorm"
)

type investorRepo struct {
	repository.BaseRepo[investor.Investor]
	writeConn *gorm.DB
	readConn  *gorm.DB
}

func NewInvestorRepo(dbMaster *gorm.DB, dbSlave *gorm.DB) investor.IInvestorRepository {
	baseRepo := repository.NewBaseRepo[investor.Investor](dbMaster, dbSlave)

	return &investorRepo{
		BaseRepo:  *baseRepo,
		writeConn: dbMaster,
		readConn:  dbSlave,
	}
}
