package loan

import (
	"time"

	"github.com/google/uuid"
)

type CreateLoanRequest struct {
	BorrowerID         uuid.UUID `json:"borrower_id" validate:"required"`
	PrincipalAmount    float64   `json:"principal_amount" validate:"required,min=1"`
	Rate               float32   `json:"rate" validate:"required,min=0"`
	ROI                float32   `json:"roi" validate:"required,min=0"`
	AgreementLetterURL string    `json:"agreement_letter_url" validate:"required,url"`
}

type RejectLoanRequest struct {
	RejectReason string `json:"reject_reason" validate:"required"`
}

type ApproveLoanRequest struct {
	ValidatorEmployeeID  uuid.UUID `json:"validator_employee_id" validate:"required"`
	VisitProofPictureURL string    `json:"visit_proof_picture_url" validate:"required,url"`
}

type LoanAgreementResponse struct {
	LoanID          uuid.UUID
	PrincipalAmount float64
	InterestRate    float32
	BorrowerName    string
}

type DisburseLoanRequest struct {
	SignedAgreementURL string    `json:"signed_agreement_url" validate:"required,url"`
	OfficerEmployeeID  uuid.UUID `json:"officer_employee_id" validate:"required"`
	DisbursementDate   time.Time `json:"disbursement_date" validate:"required"`
}
