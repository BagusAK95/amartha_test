package investment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type IInvestmentUsecase interface {
	AddInvestment(ctx context.Context, investorID uuid.UUID, req CreateInvestmentRequest) (res *Investment, err error)
	GetInvestmentAgreementDetail(ctx context.Context, investmentID uuid.UUID) (*InvestmentAgreementResponse, error)
}

// TODO: move to DTO folder
type CreateInvestmentRequest struct {
	LoanID uuid.UUID `json:"loan_id"`
	Amount float64   `json:"amount"`
}

type InvestmentAgreementResponse struct {
	AgreementID      uuid.UUID
	AgreementDate    time.Time
	InvestmentAmount float64
	InterestRate     float32
	LoanID           uuid.UUID
	LoanTerm         int
	InvestorName     string
	BorrowerName     string
}
