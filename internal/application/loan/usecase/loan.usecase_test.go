package usecase_test

import (
	"context"
	"testing"

	"github.com/BagusAK95/amarta_test/internal/application/loan/usecase"
	"github.com/BagusAK95/amarta_test/internal/domain/borrower"
	borrowerMock "github.com/BagusAK95/amarta_test/internal/domain/borrower/mock"
	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/BagusAK95/amarta_test/internal/domain/employee"
	employeeMock "github.com/BagusAK95/amarta_test/internal/domain/employee/mock"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	loanMock "github.com/BagusAK95/amarta_test/internal/domain/loan/mock"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoan(t *testing.T) {
	ctx := context.Background()
	borrowerID := uuid.New()
	req := loan.CreateLoanRequest{
		BorrowerID:      borrowerID,
		PrincipalAmount: 1000,
		Rate:            0.1,
		ROI:             100,
	}
	borrowerData := borrower.Borrower{
		BaseModel: model.BaseModel{ID: borrowerID},
		FullName:  "test borrower",
	}
	loanData := loan.Loan{
		BaseModel:       model.BaseModel{ID: uuid.New()},
		BorrowerID:      borrowerID,
		PrincipalAmount: 1000,
		Rate:            0.1,
		ROI:             100,
		State:           loan.StateProposed,
	}

	t.Run("success", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		borrowerRepo.On("GetByID", mock.Anything, borrowerID).Return(borrowerData, nil)
		loanRepo.On("Create", mock.Anything, mock.AnythingOfType("loan.Loan")).Return(loanData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.CreateLoan(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, loanData.ID, res.ID)
		borrowerRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
	})

	t.Run("borrower not found", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		borrowerRepo.On("GetByID", mock.Anything, borrowerID).Return(borrower.Borrower{}, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.CreateLoan(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("borrower not found"), err)
		borrowerRepo.AssertExpectations(t)
	})

	t.Run("create loan error", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		borrowerRepo.On("GetByID", mock.Anything, borrowerID).Return(borrowerData, nil)
		loanRepo.On("Create", mock.Anything, mock.AnythingOfType("loan.Loan")).Return(loan.Loan{}, assert.AnError)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.CreateLoan(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		borrowerRepo.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
	})
}

func TestRejectLoan(t *testing.T) {
	ctx := context.Background()
	loanID := uuid.New()
	reason := "some reason"
	loanData := loan.Loan{
		BaseModel: model.BaseModel{ID: loanID},
		State:     loan.StateProposed,
	}

	t.Run("success", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		loanRepo.On("UpdateWithMap", mock.Anything, loanID, mock.Anything).Return(loanData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.RejectLoan(ctx, loanID, reason)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		loanRepo.AssertExpectations(t)
	})

	t.Run("loan not found", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loan.Loan{}, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.RejectLoan(ctx, loanID, reason)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("loan not found"), err)
		loanRepo.AssertExpectations(t)
	})

	t.Run("loan not in proposed state", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanData.State = loan.StateApproved
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.RejectLoan(ctx, loanID, reason)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewBadRequestError("loan is not in proposed state"), err)
		loanRepo.AssertExpectations(t)
	})
}

func TestApproveLoan(t *testing.T) {
	ctx := context.Background()
	loanID := uuid.New()
	employeeID := uuid.New()
	req := loan.ApproveLoanRequest{
		ValidatorEmployeeID: employeeID,
	}
	loanData := loan.Loan{
		BaseModel: model.BaseModel{ID: loanID},
		State:     loan.StateProposed,
	}
	employeeData := employee.Employee{
		BaseModel: model.BaseModel{ID: employeeID},
	}

	t.Run("success", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		employeeRepo.On("GetByID", mock.Anything, employeeID).Return(employeeData, nil)
		loanRepo.On("UpdateWithMap", mock.Anything, loanID, mock.Anything).Return(loanData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.ApproveLoan(ctx, loanID, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		loanRepo.AssertExpectations(t)
		employeeRepo.AssertExpectations(t)
	})

	t.Run("loan not found", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loan.Loan{}, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.ApproveLoan(ctx, loanID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("loan not found"), err)
		loanRepo.AssertExpectations(t)
	})

	t.Run("loan not in proposed state", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanData.State = loan.StateApproved
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.ApproveLoan(ctx, loanID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewBadRequestError("loan is not in proposed state"), err)
		loanRepo.AssertExpectations(t)
	})

	t.Run("validator employee not found", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanData.State = loan.StateProposed
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		employeeRepo.On("GetByID", mock.Anything, employeeID).Return(employee.Employee{}, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.ApproveLoan(ctx, loanID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("validator employee not found"), err)
		loanRepo.AssertExpectations(t)
		employeeRepo.AssertExpectations(t)
	})
}

func TestDisburseLoan(t *testing.T) {
	ctx := context.Background()
	loanID := uuid.New()
	employeeID := uuid.New()
	req := loan.DisburseLoanRequest{
		OfficerEmployeeID: employeeID,
	}
	loanData := loan.Loan{
		BaseModel: model.BaseModel{ID: loanID},
		State:     loan.StateInvested,
	}
	employeeData := employee.Employee{
		BaseModel: model.BaseModel{ID: employeeID},
	}

	t.Run("success", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		employeeRepo.On("GetByID", mock.Anything, employeeID).Return(employeeData, nil)
		loanRepo.On("UpdateWithMap", mock.Anything, loanID, mock.Anything).Return(loanData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.DisburseLoan(ctx, loanID, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		loanRepo.AssertExpectations(t)
		employeeRepo.AssertExpectations(t)
	})

	t.Run("loan not found", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loan.Loan{}, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.DisburseLoan(ctx, loanID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("loan not found"), err)
		loanRepo.AssertExpectations(t)
	})

	t.Run("loan not in invested state", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanData.State = loan.StateProposed
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.DisburseLoan(ctx, loanID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewBadRequestError("loan is not in invested state"), err)
		loanRepo.AssertExpectations(t)
	})

	t.Run("officer employee not found", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanData.State = loan.StateInvested
		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		employeeRepo.On("GetByID", mock.Anything, employeeID).Return(employee.Employee{}, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.DisburseLoan(ctx, loanID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("officer employee not found"), err)
		loanRepo.AssertExpectations(t)
		employeeRepo.AssertExpectations(t)
	})
}

func TestListLoan(t *testing.T) {
	ctx := context.Background()
	state := "proposed"
	page := 1
	limit := 10
	loanData := []loan.Loan{
		{BaseModel: model.BaseModel{ID: uuid.New()}, State: loan.StateProposed},
	}
	pagination := repository.Pagination[loan.Loan]{
		Data: loanData,
	}

	t.Run("success", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("Pagination", mock.Anything, mock.Anything, page, limit).Return(pagination, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.ListLoan(ctx, &state, page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		loanRepo.AssertExpectations(t)
	})
}

func TestDetailLoan(t *testing.T) {
	ctx := context.Background()
	loanID := uuid.New()
	loanData := loan.Loan{
		BaseModel: model.BaseModel{ID: loanID},
	}

	t.Run("success", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.DetailLoan(ctx, loanID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		loanRepo.AssertExpectations(t)
	})
}

func TestGetLoanAgreementDetail(t *testing.T) {
	ctx := context.Background()
	loanID := uuid.New()
	borrowerID := uuid.New()
	loanData := loan.Loan{
		BaseModel:  model.BaseModel{ID: loanID},
		BorrowerID: borrowerID,
		State:      loan.StateInvested,
	}
	borrowerData := borrower.Borrower{
		BaseModel: model.BaseModel{ID: borrowerID},
	}

	t.Run("success", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		borrowerRepo.On("GetByID", mock.Anything, borrowerID).Return(borrowerData, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.GetLoanAgreementDetail(ctx, loanID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		loanRepo.AssertExpectations(t)
		borrowerRepo.AssertExpectations(t)
	})

	t.Run("loan not found", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loan.Loan{}, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.GetLoanAgreementDetail(ctx, loanID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("loan not found"), err)
		loanRepo.AssertExpectations(t)
	})

	t.Run("borrower not found", func(t *testing.T) {
		loanRepo := new(loanMock.MockILoanRepository)
		borrowerRepo := new(borrowerMock.MockIBorrowerRepository)
		employeeRepo := new(employeeMock.MockIEmployeeRepository)

		loanRepo.On("GetByID", mock.Anything, loanID).Return(loanData, nil)
		borrowerRepo.On("GetByID", mock.Anything, borrowerID).Return(borrower.Borrower{}, nil)

		uc := usecase.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
		res, err := uc.GetLoanAgreementDetail(ctx, loanID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, httpError.NewNotFoundError("borrower not found"), err)
		loanRepo.AssertExpectations(t)
		borrowerRepo.AssertExpectations(t)
	})
}
