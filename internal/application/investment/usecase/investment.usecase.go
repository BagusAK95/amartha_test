package usecase

import (
	"context"

	"github.com/BagusAK95/amarta_test/internal/domain/investment"
	"github.com/BagusAK95/amarta_test/internal/domain/investor"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type investmentUsecase struct {
	investmentRepo investment.IInvestmentRepository
	investorRepo   investor.IInvestorRepository
	loanRepo       loan.ILoanRepository
}

func NewInvestmentUsecase(investmentRepo investment.IInvestmentRepository, investorRepo investor.IInvestorRepository, loanRepo loan.ILoanRepository) investment.IInvestmentUsecase {
	return &investmentUsecase{
		investmentRepo: investmentRepo,
		investorRepo:   investorRepo,
		loanRepo:       loanRepo,
	}
}

func (u *investmentUsecase) AddInvestment(ctx context.Context, investorID uuid.UUID, req investment.CreateInvestmentRequest) (res *investment.Investment, err error) {
	trx := u.investmentRepo.BeginTransaction(ctx)
	defer func() {
		if err != nil {
			u.investmentRepo.Rollback(trx)
			return
		}

		u.investmentRepo.Commit(trx)
	}()

	validLoan, err := u.loanRepo.GetByIDLockTx(ctx, req.LoanID, trx)
	if err != nil {
		return nil, err
	} else if validLoan.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("loan not found")
	} else if validLoan.State != loan.StateApproved {
		return nil, httpError.NewBadRequestError("loan is not in approved")
	}

	validInvestor, err := u.investorRepo.GetByIDLockTx(ctx, investorID, trx)
	if err != nil {
		return nil, err
	} else if validInvestor.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("investor not found")
	} else if validInvestor.Balance < req.Amount {
		return nil, httpError.NewBadRequestError("insufficient balance")
	}

	totalInvestment, err := u.investmentRepo.GetTotalInvestmentByLoanID(ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	lastInvestment := totalInvestment + req.Amount
	if lastInvestment > validLoan.PrincipalAmount {
		return nil, httpError.NewBadRequestError("total investment would exceed loan principal amount")
	}

	newInvestment, err := u.investmentRepo.CreateWithTx(ctx, investment.Investment{
		LoanID:     req.LoanID,
		InvestorID: investorID,
		Amount:     req.Amount,
	}, trx)
	if err != nil {
		return nil, err
	}

	_, err = u.investorRepo.UpdateWithMapTx(ctx, investorID, map[string]any{
		"balance": validInvestor.Balance - req.Amount,
	}, trx)
	if err != nil {
		return nil, err
	}

	err = u.checkLoanInvested(ctx, validLoan, lastInvestment, trx)
	if err != nil {
		return nil, err
	}

	// TODO: Sent email to investor

	return &newInvestment, nil
}

func (u *investmentUsecase) checkLoanInvested(ctx context.Context, validLoan loan.Loan, lastInvestment float64, trx *gorm.DB) (err error) {
	if lastInvestment != validLoan.PrincipalAmount {
		return
	}

	_, err = u.loanRepo.UpdateWithMapTx(ctx, validLoan.ID, map[string]any{
		"state": loan.StateInvested,
	}, trx)
	if err != nil {
		return err
	}

	// TODO: Sent email to borrower

	return nil
}
