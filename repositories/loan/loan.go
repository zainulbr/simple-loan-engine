package loan

import (
	"context"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/loan"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *loan.Loan) error
	Approve(ctx context.Context, approval *loan.LoanApproval) error
	GetLoanDetail(ctx context.Context, loanID uuid.UUID) (*loan.LoanDetail, error)
	CreateInvestment(ctx context.Context, investment *loan.LoanInvestment) error
	CreateDisbursement(ctx context.Context, disbursement *loan.LoanDisbursement) error
}

type loanModelPG struct {
	tableName struct{} `pg:"loan.loans"` // Schema "loan", Table "loans"
	*loan.Loan
}
