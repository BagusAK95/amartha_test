package loan

import (
	"context"

	"github.com/google/uuid"
)

type ILoanUsecase interface {
	CreateLoan(ctx context.Context, req CreateLoanRequest) (*Loan, error)
	RejectLoan(ctx context.Context, loanID uuid.UUID, rejectReason string) (*Loan, error)
	ApproveLoan(ctx context.Context, loanID uuid.UUID, req ApproveLoanRequest) (*Loan, error)
}

// TODO: Move to DTO folder
type CreateLoanRequest struct {
	BorrowerID         uuid.UUID `json:"borrower_id" binding:"required"`
	PrincipalAmount    float64   `json:"principal_amount" binding:"required"`
	Rate               float32   `json:"rate" binding:"required"`
	ROI                float32   `json:"roi" binding:"required"`
	AgreementLetterURL string    `json:"agreement_letter_url" binding:"required"`
}

type RejectLoanRequest struct {
	RejectReason string `json:"reject_reason" binding:"required"`
}

type ApproveLoanRequest struct {
	ValidatorEmployeeID  uuid.UUID `json:"validator_employee_id" binding:"required"`
	VisitProofPictureURL string    `json:"visit_proof_picture_url" binding:"required"`
}
