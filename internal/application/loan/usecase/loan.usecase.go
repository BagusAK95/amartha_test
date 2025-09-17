package usecase

import (
	"context"

	"github.com/BagusAK95/amarta_test/internal/domain/loan"
)

type loanUsecase struct {
	loanRepo loan.ILoanRepository
}

func NewLoanUsecase(loanRepo loan.ILoanRepository) loan.ILoanUsecase {
	return &loanUsecase{
		loanRepo: loanRepo,
	}
}

func (u *loanUsecase) CreateLoan(ctx context.Context, body loan.CreateLoanRequest) (*loan.Loan, error) {
	loan, err := u.loanRepo.Create(ctx, loan.Loan{
		BorrowerID:      body.BorrowerID,
		PrincipalAmount: body.Principal,
		Rate:            body.Rate,
		ROI:             body.ROI,
	})
	if err != nil {
		return nil, err
	}

	return &loan, nil
}
