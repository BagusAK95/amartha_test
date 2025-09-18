package usecase

import (
	"context"
	"time"

	"github.com/BagusAK95/amarta_test/internal/domain/borrower"
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/BagusAK95/amarta_test/internal/domain/employee"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type loanUsecase struct {
	loanRepo     loan.ILoanRepository
	borrowerRepo borrower.IBorrowerRepository
	employeeRepo employee.IEmployeeRepository
}

func NewLoanUsecase(loanRepo loan.ILoanRepository, borrowerRepo borrower.IBorrowerRepository, employeeRepo employee.IEmployeeRepository) loan.ILoanUsecase {
	return &loanUsecase{
		loanRepo:     loanRepo,
		borrowerRepo: borrowerRepo,
		employeeRepo: employeeRepo,
	}
}

func (u *loanUsecase) CreateLoan(ctx context.Context, req loan.CreateLoanRequest) (*loan.Loan, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loanUsecase.CreateLoan")
	defer span.Finish()

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

func (u *loanUsecase) RejectLoan(ctx context.Context, loanID uuid.UUID, rejectReason string) (*loan.Loan, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loanUsecase.RejectLoan")
	defer span.Finish()

	validLoan, err := u.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	} else if validLoan.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("loan not found")
	} else if validLoan.State != loan.StateProposed {
		return nil, httpError.NewBadRequestError("loan is not in proposed state")
	}

	updatedLoan, err := u.loanRepo.UpdateWithMap(ctx, loanID, map[string]any{
		"state":         loan.StateRejected,
		"reject_reason": rejectReason,
	})
	if err != nil {
		return nil, err
	}

	return &updatedLoan, nil
}

func (u *loanUsecase) ApproveLoan(ctx context.Context, loanID uuid.UUID, req loan.ApproveLoanRequest) (*loan.Loan, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loanUsecase.ApproveLoan")
	defer span.Finish()

	validLoan, err := u.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	} else if validLoan.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("loan not found")
	} else if validLoan.State != loan.StateProposed {
		return nil, httpError.NewBadRequestError("loan is not in proposed state")
	}

	validEmployee, err := u.employeeRepo.GetByID(ctx, req.ValidatorEmployeeID)
	if err != nil {
		return nil, err
	} else if validEmployee.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("validator employee not found")
	}

	updatedLoan, err := u.loanRepo.UpdateWithMap(ctx, loanID, map[string]any{
		"state":                   loan.StateApproved,
		"approval_date":           time.Now(),
		"validator_employee_id":   req.ValidatorEmployeeID,
		"visit_proof_picture_url": req.VisitProofPictureURL,
	})
	if err != nil {
		return nil, err
	}

	return &updatedLoan, nil
}

func (u *loanUsecase) DisburseLoan(ctx context.Context, loanID uuid.UUID, req loan.DisburseLoanRequest) (*loan.Loan, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loanUsecase.DisburseLoan")
	defer span.Finish()

	validLoan, err := u.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	} else if validLoan.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("loan not found")
	} else if validLoan.State != loan.StateInvested {
		return nil, httpError.NewBadRequestError("loan is not in invested state")
	}

	validEmployee, err := u.employeeRepo.GetByID(ctx, req.OfficerEmployeeID)
	if err != nil {
		return nil, err
	} else if validEmployee.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("officer employee not found")
	}

	updatedLoan, err := u.loanRepo.UpdateWithMap(ctx, loanID, map[string]any{
		"state":                loan.StateDisbursed,
		"disbursement_date":    req.DisbursementDate,
		"officer_employee_id":  req.OfficerEmployeeID,
		"signed_agreement_url": req.SignedAgreementURL,
	})
	if err != nil {
		return nil, err
	}

	return &updatedLoan, nil
}

func (u *loanUsecase) DetailLoan(ctx context.Context, loanID uuid.UUID) (*loan.Loan, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loanUsecase.DetailLoan")
	defer span.Finish()

	validLoan, err := u.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	} else if validLoan.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("loan not found")
	}

	return &validLoan, nil
}

func (u *loanUsecase) ListLoan(ctx context.Context, state *string, page int, limit int) (repository.Pagination[loan.Loan], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loanUsecase.ListLoan")
	defer span.Finish()

	filter := map[string]any{}
	if state != nil {
		filter["state"] = *state
	}

	loans, err := u.loanRepo.Pagination(ctx, filter, page, limit)
	if err != nil {
		return repository.Pagination[loan.Loan]{}, err
	}

	return loans, nil
}

func (u *loanUsecase) GetLoanAgreementDetail(ctx context.Context, loanID uuid.UUID) (*loan.LoanAgreementResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loanUsecase.GetLoanAgreementDetail")
	defer span.Finish()

	loanData, err := u.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	} else if loanData.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("loan not found")
	} else if loanData.State == loan.StateProposed || loanData.State == loan.StateApproved {
		return nil, httpError.NewNotFoundError("loan not found")
	}

	borrowerData, err := u.borrowerRepo.GetByID(ctx, loanData.BorrowerID)
	if err != nil {
		return nil, err
	} else if borrowerData.ID == uuid.Nil {
		return nil, httpError.NewNotFoundError("borrower not found")
	}

	return &loan.LoanAgreementResponse{
		LoanID:          loanData.ID,
		PrincipalAmount: loanData.PrincipalAmount,
		InterestRate:    loanData.Rate,
		BorrowerName:    borrowerData.FullName,
	}, nil
}
