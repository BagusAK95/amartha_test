package loan

import (
	"context"

	"github.com/google/uuid"
)

type ILoanUsecase interface {
	CreateLoan(context.Context, CreateLoanRequest) (*Loan, error)
}

// TODO: Move to DTO folder
type CreateLoanRequest struct {
	BorrowerID         uuid.UUID `json:"borrower_id" binding:"required"`
	PrincipalAmount    float64   `json:"principal_amount" binding:"required"`
	Rate               float32   `json:"rate" binding:"required"`
	ROI                float32   `json:"roi" binding:"required"`
	AgreementLetterURL string    `json:"agreement_letter_url" binding:"required"`
}
