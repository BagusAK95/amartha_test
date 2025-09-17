package loan

import "context"

type ILoanUsecase interface {
	CreateLoan(context.Context, CreateLoanRequest) (*Loan, error)
}

// TODO: Move to DTO folder
type CreateLoanRequest struct {
	BorrowerID         string  `json:"borrower_id" binding:"required"`
	PrincipalAmount    float64 `json:"principal_amount" binding:"required"`
	Rate               float64 `json:"rate" binding:"required"`
	ROI                float64 `json:"roi" binding:"required"`
	AgreementLetterURL string  `json:"agreement_letter_url" binding:"required"`
}
