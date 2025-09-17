package repository

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/BagusAK95/amarta_test/internal/domain/investment"
	"gorm.io/gorm"
)

type investmentRepo struct {
	repository.BaseRepo[investment.Investment]
	writeConn *gorm.DB
	readConn  *gorm.DB
}

func NewInvestmentRepo(dbMaster *gorm.DB, dbSlave *gorm.DB) investment.IInvestmentRepository {
	baseRepo := repository.NewBaseRepo[investment.Investment](dbMaster, dbSlave)

	return &investmentRepo{
		BaseRepo:  *baseRepo,
		writeConn: dbMaster,
		readConn:  dbSlave,
	}
}
