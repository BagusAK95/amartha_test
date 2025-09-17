package usecase

import (
	"context"

	"github.com/BagusAK95/amarta_test/internal/domain/borrower"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/google/uuid"
)

type loanUsecase struct {
	loanRepo     loan.ILoanRepository
	borrowerRepo borrower.IBorrowerRepository
}

func NewLoanUsecase(loanRepo loan.ILoanRepository, borrowerRepo borrower.IBorrowerRepository) loan.ILoanUsecase {
	return &loanUsecase{
		loanRepo:     loanRepo,
		borrowerRepo: borrowerRepo,
	}
}

func (u *loanUsecase) CreateLoan(ctx context.Context, req loan.CreateLoanRequest) (*loan.Loan, error) {
	borrower, err := u.borrowerRepo.GetByID(ctx, req.BorrowerID)
	if err != nil {
		return nil, err
	} else if borrower.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("borrower not found")
	}

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
