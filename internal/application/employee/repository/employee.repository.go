package repository

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/BagusAK95/amarta_test/internal/domain/employee"
	"gorm.io/gorm"
)

type employeeRepo struct {
	repository.BaseRepo[employee.Employee]
	writeConn *gorm.DB
	readConn  *gorm.DB
}

func NewEmployeeRepo(dbMaster *gorm.DB, dbSlave *gorm.DB) employee.IEmployeeRepository {
	baseRepo := repository.NewBaseRepo[employee.Employee](dbMaster, dbSlave)

	return &employeeRepo{
		BaseRepo:  *baseRepo,
		writeConn: dbMaster,
		readConn:  dbSlave,
	}
}
