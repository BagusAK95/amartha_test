package employee

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
)

type IEmployeeRepository interface {
	repository.IBaseRepo[Employee]
}
