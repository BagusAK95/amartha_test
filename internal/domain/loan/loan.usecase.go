package loan

import (
	"context"

	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/google/uuid"
)

type ILoanUsecase interface {
	CreateLoan(ctx context.Context, req CreateLoanRequest) (*Loan, error)
	RejectLoan(ctx context.Context, loanID uuid.UUID, rejectReason string) (*Loan, error)
	ApproveLoan(ctx context.Context, loanID uuid.UUID, req ApproveLoanRequest) (*Loan, error)
	DisburseLoan(ctx context.Context, loanID uuid.UUID, req DisburseLoanRequest) (*Loan, error)
	ListLoan(ctx context.Context, state *string, page int, limit int) (repository.Pagination[Loan], error)
	DetailLoan(ctx context.Context, loanID uuid.UUID) (*Loan, error)
	GetLoanAgreementDetail(ctx context.Context, loanID uuid.UUID) (*LoanAgreementResponse, error)
}
