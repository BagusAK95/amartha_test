package investment

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
	"github.com/google/uuid"
)

type Investment struct {
	model.BaseModel
	LoanID     uuid.UUID `json:"loan_id"`
	InvestorID uuid.UUID `json:"investor_id"`
	Amount     float64   `json:"amount"`
}

func (Investment) TableName() string {
	return "investments"
}

func (e Investment) TracerName() string {
	return "investmentRepo"
}
