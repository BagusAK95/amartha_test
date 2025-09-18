package loan

import (
	"time"

	"github.com/google/uuid"
)

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

type LoanAgreementResponse struct {
	LoanID          uuid.UUID
	PrincipalAmount float64
	InterestRate    float32
	BorrowerName    string
}

type DisburseLoanRequest struct {
	SignedAgreementURL string    `json:"signed_agreement_url" binding:"required"`
	OfficerEmployeeID  uuid.UUID `json:"officer_employee_id" binding:"required"`
	DisbursementDate   time.Time `json:"disbursement_date" binding:"required"`
}
