package loan

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/loan"

	"github.com/go-pg/pg/v10"
)

type loanRepo struct {
	db *pg.DB
}

func NewLoanRepository(db *pg.DB) LoanRepository {
	return &loanRepo{db: db}
}

// Create Loan
func (r *loanRepo) CreateLoan(ctx context.Context, loan *loan.Loan) error {
	data := loanModelPG{Loan: loan}
	_, err := r.db.Model(&data).Context(ctx).Insert()
	return err
}

// Approval is Approve & Update Loan Status
func (r *loanRepo) Approve(ctx context.Context, approval *loan.LoanApproval) error {
	_, err := r.db.ExecContext(ctx, `
		WITH inserted_approval AS (
			INSERT INTO loan.approvals (loan_id, visited_file, approved_by)
			VALUES ( ?, ?, ?)
			RETURNING loan_id
		)
		UPDATE loan.loans
		SET state = 'approved', approval_date = ?, rate = ?
		WHERE loan_id = (SELECT loan_id FROM inserted_approval);
	`, approval.LoanId, approval.VisitedFile, approval.ApprovedBy,
		approval.ApprovalDate, approval.Rate)
	return err
}

func (r *loanRepo) GetLoanDetail(ctx context.Context, loanID uuid.UUID) (*loan.LoanDetail, error) {
	loanDetail := &loan.LoanDetail{}

	_, err := r.db.QueryOneContext(ctx, loanDetail, `
		SELECT 
			l.loan_id, 
			l.description, 
			l.proposed_by, 
			l.amount, 
			l.duration_month, 
			l.rate, 
			l.state, 
			l.approval_date, 
			l.aggrement_file, 
			l.created_at, 
			l.updated_at,
			COALESCE(SUM(i.amount), 0) AS total_investment
		FROM loan.loans l
		LEFT JOIN loan.investments i ON l.loan_id = i.loan_id
		WHERE l.loan_id = ?
		GROUP BY l.loan_id
	`, loanID)

	if err != nil {
		return nil, err
	}
	return loanDetail, nil
}

// Create Loan Investment
func (r *loanRepo) CreateInvestment(ctx context.Context, investment *loan.LoanInvestment) error {
	_, err := r.db.ExecContext(ctx, `
			INSERT INTO loan.investments (loan_id, amount, invested_by)
			VALUES (?, ?, ?)
	`, investment.LoanId, investment.Amount, investment.InvestedBy)
	// on create/update row of table loan.investments will check constraint function check_total_investment

	if err != nil {
		return fmt.Errorf("failed to create investment: %w", err)
	}
	return nil
}

// Create Loan Disbursement
func (r *loanRepo) CreateDisbursement(ctx context.Context, disbursement *loan.LoanDisbursement) error {
	_, err := r.db.ExecContext(ctx, `
		WITH inserted_disbursement AS (
			INSERT INTO loan.disbursments (disbursment_id, loan_id, disbursment_date, disbursment_by)
			VALUES (gen_random_uuid(), ?, ?, ?)
			RETURNING loan_id
		)
		UPDATE loan.loans
		SET state = 'disbursed'
		WHERE loan_id = (SELECT loan_id FROM inserted_disbursement);
	`, disbursement.LoanId, disbursement.DisbursmentDate, disbursement.DisbursementBy)
	return err
}
