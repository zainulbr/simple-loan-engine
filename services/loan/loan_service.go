package loan

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/loan"
	loanRepository "github.com/zainulbr/simple-loan-engine/repositories/loan"
)

// Implementation
type loanService struct {
	loanRepo loanRepository.LoanRepository
}

// Constructor
func NewLoanService(loanRepo loanRepository.LoanRepository) LoanService {
	return &loanService{
		loanRepo: loanRepo,
	}
}

// Create Loan
func (s *loanService) CreateLoan(ctx context.Context, loan *loan.Loan) (*loan.Loan, error) {
	err := s.loanRepo.CreateLoan(ctx, loan)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

// Approve Loan
func (s *loanService) ApproveLoan(ctx context.Context, approval *loan.LoanApproval) error {
	loanDetail, err := s.loanRepo.GetLoanDetail(ctx, approval.LoanId)
	if err != nil {
		return err
	}

	// Check if loan is in valid state for approval
	if loanDetail.State != loan.StateProposed {
		return errors.New("loan is not in a valid state for approval")
	}

	// Create approval & Update loan state to approved
	err = s.loanRepo.Approve(ctx, approval)
	if err != nil {
		return err
	}

	return nil
}

// Get Loan Detail
func (s *loanService) GetLoanDetail(ctx context.Context, loanID uuid.UUID) (*loan.LoanDetail, error) {
	return s.loanRepo.GetLoanDetail(ctx, loanID)
}

// Create Investment
func (s *loanService) CreateInvestment(ctx context.Context, investment *loan.LoanInvestment) error {
	loanDetail, err := s.loanRepo.GetLoanDetail(ctx, investment.LoanId)
	if err != nil {
		return err
	}

	// Check if loan is in valid state for insvestment
	if loanDetail.State != loan.StateApproved {
		return errors.New("loan is not in a valid state for investment")
	}

	// Check if total investment matches loan amount
	if (loanDetail.TotalInvestment + investment.Amount) > investment.Amount {
		return errors.New("investment amount exceeds loan amount")
	}

	err = s.loanRepo.CreateInvestment(ctx, investment)
	if err != nil {
		return err
	}

	return nil
}

// Create Disbursement
func (s *loanService) CreateDisbursement(ctx context.Context, disbursement *loan.LoanDisbursement) error {

	loanDetail, err := s.loanRepo.GetLoanDetail(ctx, disbursement.LoanId)
	if err != nil {
		return err
	}

	// Check if loan is in valid state for disbursement
	if loanDetail.State != loan.StateInvested {
		return errors.New("loan is not in a valid state for disbursement")
	}

	// create disbursement and 	Update loan state after disbursement
	err = s.loanRepo.CreateDisbursement(ctx, disbursement)
	if err != nil {
		return err
	}

	return nil
}
