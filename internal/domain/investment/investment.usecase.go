package investment

import (
	"context"

	"github.com/google/uuid"
)

type IInvestmentUsecase interface {
	AddInvestment(ctx context.Context, investorID uuid.UUID, req CreateInvestmentRequest) (res *Investment, err error)
}

// TODO: move to DTO folder
type CreateInvestmentRequest struct {
	LoanID uuid.UUID `json:"loan_id"`
	Amount float64   `json:"amount"`
}
