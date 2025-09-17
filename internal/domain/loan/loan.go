package loan

import (
	"time"

	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
)

type Loan struct {
	model.BaseModel
	BorrowerID         string
	PrincipalAmount    float64
	Rate               float64
	ROI                float64
	State              State
	AgreementLetterURL string

	ApprovalDetails     ApprovalDetails     `gorm:"embedded"`
	DisbursementDetails DisbursementDetails `gorm:"embedded"`
}

func (Loan) TableName() string {
	return "loans"
}

type State string

const (
	StateProposed  State = "proposed"
	StateApproved  State = "approved"
	StateInvested  State = "invested"
	StateDisbursed State = "disbursed"
)

type ApprovalDetails struct {
	FieldValidatorEmployeeID string    `json:"field_validator_employee_id"`
	VisitProofPictureURL     string    `json:"visit_proof_picture_url"`
	ApprovalDate             time.Time `json:"approval_date"`
}

type DisbursementDetails struct {
	FieldOfficerEmployeeID string    `json:"field_officer_employee_id"`
	SignedAgreementURL     string    `json:"signed_agreement_url"`
	DisbursementDate       time.Time `json:"disbursement_date"`
}
