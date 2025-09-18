package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/BagusAK95/amarta_test/internal/application/investment/usecase"
	"github.com/BagusAK95/amarta_test/internal/domain/borrower"
	borrowerMock "github.com/BagusAK95/amarta_test/internal/domain/borrower/mock"
	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
	"github.com/BagusAK95/amarta_test/internal/domain/investment"
	investmentMock "github.com/BagusAK95/amarta_test/internal/domain/investment/mock"
	"github.com/BagusAK95/amarta_test/internal/domain/investor"
	investorMock "github.com/BagusAK95/amarta_test/internal/domain/investor/mock"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	loanMock "github.com/BagusAK95/amarta_test/internal/domain/loan/mock"
	"github.com/BagusAK95/amarta_test/internal/domain/mail"
	busMock "github.com/BagusAK95/amarta_test/internal/infrastructure/bus/mock"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAddInvestment(t *testing.T) {
	ctx := context.Background()
	investorID := uuid.New()
	loanID := uuid.New()
	req := investment.CreateInvestmentRequest{
		LoanID: loanID,
		Amount: 1000,
	}
	loanData := loan.Loan{
		BaseModel:       model.BaseModel{ID: loanID},
		State:           loan.StateApproved,
		PrincipalAmount: 2000,
	}
	investorData := investor.Investor{
		BaseModel: model.BaseModel{ID: investorID},
		Balance:   5000,
	}
	investmentData := investment.Investment{
		BaseModel:  model.BaseModel{ID: uuid.New()},
		LoanID:     loanID,
		InvestorID: investorID,
		Amount:     1000,
	}

	t.Run("success", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investmentRepo.On("BeginTransaction", mock.Anything).Return(&gorm.DB{})
		loanRepo.On("GetByIDLockTx", mock.Anything, loanID, mock.Anything).Return(loanData, nil)
		investorRepo.On("GetByIDLockTx", mock.Anything, investorID, mock.Anything).Return(investorData, nil)
		investmentRepo.On("GetTotalInvestmentByLoanID", mock.Anything, loanID).Return(float64(0), nil)
		investmentRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("investment.Investment"), mock.Anything).Return(investmentData, nil)
		investorRepo.On("UpdateWithMapTx", mock.Anything, investorID, mock.Anything, mock.Anything).Return(investor.Investor{}, nil)
		investmentRepo.On("Commit", mock.Anything).Return(&gorm.DB{})
		mailBus.On("Publish", "mail.send", mock.Anything)

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.AddInvestment(ctx, investorID, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		investmentRepo.AssertExpectations(t)
		investorRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		mailBus.AssertExpectations(t)
	})

	t.Run("loan not found", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investmentRepo.On("BeginTransaction", mock.Anything).Return(&gorm.DB{})
		loanRepo.On("GetByIDLockTx", mock.Anything, loanID, mock.Anything).Return(loan.Loan{}, nil)
		investmentRepo.On("Rollback", mock.Anything).Return(&gorm.DB{})

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.AddInvestment(ctx, investorID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("loan not found"), err)
		loanRepo.AssertExpectations(t)
	})

	t.Run("loan not in approved state", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		loanData.State = loan.StateProposed
		investmentRepo.On("BeginTransaction", mock.Anything).Return(&gorm.DB{})
		loanRepo.On("GetByIDLockTx", mock.Anything, loanID, mock.Anything).Return(loanData, nil)
		investmentRepo.On("Rollback", mock.Anything).Return(&gorm.DB{})

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.AddInvestment(ctx, investorID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewBadRequestError("loan is not in approved"), err)
		loanRepo.AssertExpectations(t)
	})

	t.Run("investor not found", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		loanData.State = loan.StateApproved
		investmentRepo.On("BeginTransaction", mock.Anything).Return(&gorm.DB{})
		loanRepo.On("GetByIDLockTx", mock.Anything, loanID, mock.Anything).Return(loanData, nil)
		investorRepo.On("GetByIDLockTx", mock.Anything, investorID, mock.Anything).Return(investor.Investor{}, nil)
		investmentRepo.On("Rollback", mock.Anything).Return(&gorm.DB{})

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.AddInvestment(ctx, investorID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("investor not found"), err)
		loanRepo.AssertExpectations(t)
		investorRepo.AssertExpectations(t)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investorData.Balance = 500
		investmentRepo.On("BeginTransaction", mock.Anything).Return(&gorm.DB{})
		loanRepo.On("GetByIDLockTx", mock.Anything, loanID, mock.Anything).Return(loanData, nil)
		investorRepo.On("GetByIDLockTx", mock.Anything, investorID, mock.Anything).Return(investorData, nil)
		investmentRepo.On("Rollback", mock.Anything).Return(&gorm.DB{})

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.AddInvestment(ctx, investorID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewBadRequestError("insufficient balance"), err)
		loanRepo.AssertExpectations(t)
		investorRepo.AssertExpectations(t)
	})

	t.Run("total investment exceeds principal", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investorData.Balance = 5000
		investmentRepo.On("BeginTransaction", mock.Anything).Return(&gorm.DB{})
		loanRepo.On("GetByIDLockTx", mock.Anything, loanID, mock.Anything).Return(loanData, nil)
		investorRepo.On("GetByIDLockTx", mock.Anything, investorID, mock.Anything).Return(investorData, nil)
		investmentRepo.On("GetTotalInvestmentByLoanID", mock.Anything, loanID).Return(float64(1500), nil)
		investmentRepo.On("Rollback", mock.Anything).Return(&gorm.DB{})

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.AddInvestment(ctx, investorID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewBadRequestError("total investment would exceed loan principal amount"), err)
		loanRepo.AssertExpectations(t)
		investorRepo.AssertExpectations(t)
		investmentRepo.AssertExpectations(t)
	})

	t.Run("check loan invested", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		loanData.PrincipalAmount = 1000
		investmentRepo.On("BeginTransaction", mock.Anything).Return(&gorm.DB{})
		loanRepo.On("GetByIDLockTx", mock.Anything, loanID, mock.Anything).Return(loanData, nil)
		investorRepo.On("GetByIDLockTx", mock.Anything, investorID, mock.Anything).Return(investorData, nil)
		investmentRepo.On("GetTotalInvestmentByLoanID", mock.Anything, loanID).Return(float64(0), nil)
		investmentRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("investment.Investment"), mock.Anything).Return(investmentData, nil)
		investorRepo.On("UpdateWithMapTx", mock.Anything, investorID, mock.Anything, mock.Anything).Return(investor.Investor{}, nil)
		loanRepo.On("UpdateWithMapTx", mock.Anything, loanID, mock.Anything, mock.Anything).Return(loan.Loan{}, nil)
		borrowerRepo.On("GetByID", mock.Anything, mock.Anything).Return(borrower.Borrower{}, nil)
		investmentRepo.On("Commit", mock.Anything).Return(&gorm.DB{})
		mailBus.On("Publish", "mail.send", mock.Anything).Times(2)

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.AddInvestment(ctx, investorID, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		investmentRepo.AssertExpectations(t)
		investorRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		borrowerRepo.AssertExpectations(t)
		mailBus.AssertExpectations(t)
	})
}

func TestGetInvestmentAgreementDetail(t *testing.T) {
	ctx := context.Background()
	investmentID := uuid.New()
	loanID := uuid.New()
	investorID := uuid.New()
	borrowerID := uuid.New()
	now := time.Now()

	investmentData := investment.Investment{
		BaseModel:  model.BaseModel{ID: investmentID, CreatedAt: &now},
		LoanID:     loanID,
		InvestorID: investorID,
	}
	loanData := loan.Loan{
		BaseModel:  model.BaseModel{ID: loanID},
		BorrowerID: borrowerID,
	}
	investorData := investor.Investor{
		BaseModel: model.BaseModel{ID: investorID},
	}
	borrowerData := borrower.Borrower{
		BaseModel: model.BaseModel{ID: borrowerID},
	}

	t.Run("success", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investmentRepo.On("GetByID", mock.Anything, investmentID).Return(investmentData, nil)
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		investorRepo.On("GetByID", mock.Anything, investorID).Return(investorData, nil)
		borrowerRepo.On("GetByID", mock.Anything, borrowerID).Return(borrowerData, nil)

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.GetInvestmentAgreementDetail(ctx, investmentID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		investmentRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		investorRepo.AssertExpectations(t)
		borrowerRepo.AssertExpectations(t)
	})

	t.Run("investment not found", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investmentRepo.On("GetByID", mock.Anything, investmentID).Return(investment.Investment{}, nil)

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.GetInvestmentAgreementDetail(ctx, investmentID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("investment not found"), err)
		investmentRepo.AssertExpectations(t)
	})

	t.Run("loan not found", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investmentRepo.On("GetByID", mock.Anything, investmentID).Return(investmentData, nil)
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loan.Loan{}, nil)

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.GetInvestmentAgreementDetail(ctx, investmentID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("loan not found"), err)
		investmentRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
	})

	t.Run("investor not found", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investmentRepo.On("GetByID", mock.Anything, investmentID).Return(investmentData, nil)
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		investorRepo.On("GetByID", mock.Anything, investorID).Return(investor.Investor{}, nil)

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.GetInvestmentAgreementDetail(ctx, investmentID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("investor not found"), err)
		investmentRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		investorRepo.AssertExpectations(t)
	})

	t.Run("borrower not found", func(t *testing.T) {
		investmentRepo := new(investmentMock.MockIInvestmentRepository)
		investorRepo := new(investorMock.MockIInvestorRepository)
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		mailBus := new(busMock.MockBus[mail.MailSendRequest])

		investmentRepo.On("GetByID", mock.Anything, investmentID).Return(investmentData, nil)
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		investorRepo.On("GetByID", mock.Anything, investorID).Return(investorData, nil)
		borrowerRepo.On("GetByID", mock.Anything, borrowerID).Return(borrower.Borrower{}, nil)

		uc := usecase.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
		res, err := uc.GetInvestmentAgreementDetail(ctx, investmentID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("borrower not found"), err)
		investmentRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		investorRepo.AssertExpectations(t)
		borrowerRepo.AssertExpectations(t)
	})
}
