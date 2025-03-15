package loan

import (
	"context"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/loan"
)

// LoanService interface
type LoanService interface {
	CreateLoan(ctx context.Context, loan *loan.Loan) (*loan.Loan, error)
	ApproveLoan(ctx context.Context, approval *loan.LoanApproval) error
	GetLoanDetail(ctx context.Context, loanID uuid.UUID) (*loan.LoanDetail, error)
	CreateInvestment(ctx context.Context, investment *loan.LoanInvestment) error
	CreateDisbursement(ctx context.Context, disbursement *loan.LoanDisbursement) error
}
