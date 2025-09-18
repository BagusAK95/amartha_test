package investment

import (
	"time"

	"github.com/google/uuid"
)

type CreateInvestmentRequest struct {
	LoanID uuid.UUID `json:"loan_id"`
	Amount float64   `json:"amount"`
}

type InvestmentAgreementResponse struct {
	AgreementID      uuid.UUID
	AgreementDate    time.Time
	InvestmentAmount float64
	ROI              float32
	LoanID           uuid.UUID
	LoanTerm         int
	InvestorName     string
	BorrowerName     string
}
