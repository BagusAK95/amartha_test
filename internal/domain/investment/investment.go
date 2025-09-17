package investment

import "github.com/BagusAK95/amarta_test/internal/domain/common/model"

type Investment struct {
	model.BaseModel
	LoanID     string
	InvestorID string
	Amount     float64
}

func (Investment) TableName() string {
	return "investments"
}
