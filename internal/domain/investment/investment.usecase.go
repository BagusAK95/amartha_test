package investment

import (
	"context"

	"github.com/google/uuid"
)

type IInvestmentUsecase interface {
	AddInvestment(ctx context.Context, investorID uuid.UUID, req CreateInvestmentRequest) (res *Investment, err error)
	GetInvestmentAgreementDetail(ctx context.Context, investmentID uuid.UUID) (*InvestmentAgreementResponse, error)
}
