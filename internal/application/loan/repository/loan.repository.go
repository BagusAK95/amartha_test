package repository

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	"gorm.io/gorm"
)

type loanRepo struct {
	repository.BaseRepo[loan.Loan]
	writeConn *gorm.DB
	readConn  *gorm.DB
}

func NewLoanRepo(dbMaster *gorm.DB, dbSlave *gorm.DB) loan.ILoanRepository {
	baseRepo := repository.NewBaseRepo[loan.Loan](dbMaster, dbSlave)

	return &loanRepo{
		BaseRepo:  *baseRepo,
		writeConn: dbMaster,
		readConn:  dbSlave,
	}
}
