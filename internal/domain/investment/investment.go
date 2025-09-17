package investment

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
	"github.com/google/uuid"
)

type Investment struct {
	model.BaseModel
	LoanID     uuid.UUID
	InvestorID uuid.UUID
	Amount     float64
}

func (Investment) TableName() string {
	return "investments"
}
