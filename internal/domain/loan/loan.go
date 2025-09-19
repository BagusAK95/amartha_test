package loan

import (
	"time"

	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
	"github.com/google/uuid"
)

type Loan struct {
	model.BaseModel
	BorrowerID          uuid.UUID           `json:"borrower_id"`
	PrincipalAmount     float64             `json:"principal_amount"`
	Rate                float32             `json:"rate"`
	ROI                 float32             `json:"roi"`
	State               State               `json:"state"`
	AgreementLetterURL  string              `json:"agreement_letter_url"`
	ApprovalDetails     ApprovalDetails     `json:"approval_details" gorm:"embedded"`
	DisbursementDetails DisbursementDetails `json:"disbursement_details" gorm:"embedded"`
	RejectReason        *string             `json:"reject_reason"`
}

func (Loan) TableName() string {
	return "loans"
}

type State string

const (
	StateProposed  State = "proposed"
	StateApproved  State = "approved"
	StateRejected  State = "rejected"
	StateInvested  State = "invested"
	StateDisbursed State = "disbursed"
)

type ApprovalDetails struct {
	ValidatorEmployeeID  *string    `json:"validator_employee_id"`
	VisitProofPictureURL *string    `json:"visit_proof_picture_url"`
	ApprovalDate         *time.Time `json:"approval_date"`
}

type DisbursementDetails struct {
	OfficerEmployeeID  *string    `json:"officer_employee_id"`
	SignedAgreementURL *string    `json:"signed_agreement_url"`
	DisbursementDate   *time.Time `json:"disbursement_date"`
}
