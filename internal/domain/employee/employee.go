package employee

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
)

type Employee struct {
	model.BaseModel
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (Employee) TableName() string {
	return "employees"
}

func (e Employee) TracerName() string {
	return "employeeRepo"
}
