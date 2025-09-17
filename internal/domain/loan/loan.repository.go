package loan

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
)

type ILoanRepository interface {
	repository.IBaseRepo[Loan]
}
