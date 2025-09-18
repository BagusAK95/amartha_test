package borrower

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
)

type Borrower struct {
	model.BaseModel
	FullName     string `json:"full_name"`
	IDCardNumber string `json:"id_card_number"`
	Address      string `json:"address"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Status       string `json:"status"`
}

func (Borrower) TableName() string {
	return "borrowers"
}

func (e Borrower) TracerName() string {
	return "borrowerRepo"
}
