package usecase

import (
	"context"
	"time"

	"github.com/BagusAK95/amarta_test/internal/config"
	"github.com/BagusAK95/amarta_test/internal/domain/borrower"
	"github.com/BagusAK95/amarta_test/internal/domain/investment"
	"github.com/BagusAK95/amarta_test/internal/domain/investor"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	"github.com/BagusAK95/amarta_test/internal/domain/mail"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/bus"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

var tracerName = "InvestmentUsecase"
var tracer = otel.Tracer(tracerName)

type investmentUsecase struct {
	investmentRepo investment.IInvestmentRepository
	investorRepo   investor.IInvestorRepository
	loanRepo       loan.ILoanRepository
	borrowerRepo   borrower.IBorrowerRepository
	mailBus        bus.Bus[mail.MailSendRequest]
}

func NewInvestmentUsecase(investmentRepo investment.IInvestmentRepository, investorRepo investor.IInvestorRepository, loanRepo loan.ILoanRepository, borrowerRepo borrower.IBorrowerRepository, mailBus bus.Bus[mail.MailSendRequest]) investment.IInvestmentUsecase {
	return &investmentUsecase{
		investmentRepo: investmentRepo,
		investorRepo:   investorRepo,
		loanRepo:       loanRepo,
		borrowerRepo:   borrowerRepo,
		mailBus:        mailBus,
	}
}

func (u *investmentUsecase) AddInvestment(ctx context.Context, investorID uuid.UUID, req investment.CreateInvestmentRequest) (res *investment.Investment, err error) {
	ctx, span := tracer.Start(ctx, tracerName+".AddInvestment")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, config.CONTEXT_TIMEOUT)
	defer cancel()

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

	mailRequest := mail.MailSendRequest{
		To:       validInvestor.Email,
		Subject:  "Your Investment is Confirmed",
		Template: "investment_confirmed.html",
		Data: map[string]any{
			"InvestmentID":     newInvestment.ID.String(),
			"LoanID":           validLoan.ID.String(),
			"InvestorName":     validInvestor.FullName,
			"InvestmentAmount": req.Amount,
			"ROI":              validLoan.ROI,
			"AgreementDate":    newInvestment.CreatedAt,
			"AppUrl":           config.APP_URL,
			"Year":             time.Now().Year(),
		},
	}

	u.mailBus.Publish("mail.send", mailRequest)

	return &newInvestment, nil
}

func (u *investmentUsecase) checkLoanInvested(ctx context.Context, validLoan loan.Loan, lastInvestment float64, trx *gorm.DB) (err error) {
	ctx, span := tracer.Start(ctx, tracerName+".CheckLoanInvested")
	defer span.End()

	if lastInvestment != validLoan.PrincipalAmount {
		return
	}

	_, err = u.loanRepo.UpdateWithMapTx(ctx, validLoan.ID, map[string]any{
		"state": loan.StateInvested,
	}, trx)
	if err != nil {
		return err
	}

	validBorrower, err := u.borrowerRepo.GetByID(ctx, validLoan.BorrowerID)
	if err != nil {
		return err
	}

	mailRequest := mail.MailSendRequest{
		To:       validBorrower.Email,
		Subject:  "Your Loan Has Been Funded",
		Template: "loan_invested.html",
		Data: map[string]any{
			"BorrowerName": validBorrower.FullName,
			"LoanID":       validLoan.ID.String(),
			"LoanAmount":   validLoan.PrincipalAmount,
			"InterestRate": validLoan.Rate,
			"AppUrl":       config.APP_URL,
			"Year":         time.Now().Year(),
		},
	}

	u.mailBus.Publish("mail.send", mailRequest)

	return nil
}

func (u *investmentUsecase) GetInvestmentAgreementDetail(ctx context.Context, investmentID uuid.UUID) (*investment.InvestmentAgreementResponse, error) {
	ctx, span := tracer.Start(ctx, tracerName+".GetInvestmentAgreementDetail")
	defer span.End()

	inv, err := u.investmentRepo.GetByID(ctx, investmentID)
	if err != nil {
		return nil, err
	} else if inv.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("investment not found")
	}

	loanData, err := u.loanRepo.GetByID(ctx, inv.LoanID)
	if err != nil {
		return nil, err
	} else if loanData.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("loan not found")
	}

	investorData, err := u.investorRepo.GetByID(ctx, inv.InvestorID)
	if err != nil {
		return nil, err
	} else if investorData.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("investor not found")
	}

	borrowerData, err := u.borrowerRepo.GetByID(ctx, loanData.BorrowerID)
	if err != nil {
		return nil, err
	} else if borrowerData.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("borrower not found")
	}

	return &investment.InvestmentAgreementResponse{
		AgreementID:      inv.ID,
		AgreementDate:    *inv.CreatedAt,
		InvestmentAmount: inv.Amount,
		ROI:              loanData.ROI,
		LoanID:           loanData.ID,
		InvestorName:     investorData.FullName,
		BorrowerName:     borrowerData.FullName,
	}, nil
}
