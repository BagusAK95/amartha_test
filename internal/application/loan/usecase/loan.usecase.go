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

func (u *loanUsecase) CreateLoan(ctx context.Context, req loan.CreateLoanRequest) (*loan.Loan, error) {
	newLoan, err := u.loanRepo.Create(ctx, loan.Loan{
		BorrowerID:         req.BorrowerID,
		PrincipalAmount:    req.PrincipalAmount,
		Rate:               req.Rate,
		ROI:                req.ROI,
		AgreementLetterURL: req.AgreementLetterURL,
		State:              loan.StateProposed,
	})
	if err != nil {
		return nil, err
	}

	return &newLoan, nil
}
