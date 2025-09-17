package loan

import "context"

type ILoanUsecase interface {
	CreateLoan(context.Context, CreateLoanRequest) (*Loan, error)
}

type CreateLoanRequest struct {
	BorrowerID string  `json:"borrower_id" binding:"required"`
	Principal  float64 `json:"principal" binding:"required"`
	Rate       float64 `json:"rate" binding:"required"`
	ROI        float64 `json:"roi" binding:"required"`
}
