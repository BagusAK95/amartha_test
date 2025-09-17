package investment

import (
	"context"

	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/google/uuid"
)

type IInvestmentRepository interface {
	repository.IBaseRepo[Investment]
	GetTotalInvestmentByLoanID(ctx context.Context, loanID uuid.UUID) (float64, error)
}
