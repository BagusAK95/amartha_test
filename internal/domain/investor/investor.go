package investor

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
)

type Investor struct {
	model.BaseModel
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
}

func (Investor) TableName() string {
	return "investors"
}

func (e Investor) TracerName() string {
	return "investorRepo"
}
