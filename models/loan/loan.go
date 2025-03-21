package loan

import (
	"time"

	"github.com/google/uuid"
)

type LoanState string

const (
	StateProposed  LoanState = "proposed"
	StateApproved  LoanState = "approved"
	StateInvested  LoanState = "invested"
	StateRejected  LoanState = "rejected"
	StateDisbursed LoanState = "disbursed"
)

type Loan struct {
	LoanId        uuid.UUID `json:"loan_id,omitempty"`
	Description   string    `json:"description,omitempty" binding:"required"`
	ProposedBy    uuid.UUID `json:"proposed_by,omitempty"`
	Amount        float64   `json:"amount,omitempty" binding:"required"`
	DurationMonth int       `json:"duration_month,omitempty" binding:"required"`
	Rate          float64   `json:"rate,omitempty"`
	State         string    `json:"state,omitempty"`
	ApprovalDate  time.Time `json:"approval_date,omitempty"`
	AggrementFile uuid.UUID `json:"aggrement_file,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type LoanDetail struct {
	LoanId           uuid.UUID `json:"loan_id,omitempty"`
	Description      string    `json:"description,omitempty"`
	ProposedBy       uuid.UUID `json:"proposed_by,omitempty"`
	Amount           float64   `json:"amount,omitempty"`
	Rate             float64   `json:"rate,omitempty"` // Interest rate
	DurationMonth    int       `json:"duration_month,omitempty"`
	ApprovalDate     time.Time `json:"approval_date,omitempty"`
	DisbursementDate time.Time `json:"disbursement_date,omitempty"`
	State            LoanState `json:"state,omitempty"`
	AggrementFile    string    `json:"aggrement_file,omitempty"` // Draft / Signed aggrement file
	ValidationFile   string    `json:"-,omitempty"`
	TotalInvestment  float64   `json:"total_investment,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

type LoanApproval struct {
	ApprovalId   uuid.UUID `json:"approval_id,omitempty"`
	LoanId       uuid.UUID `json:"loan_id,omitempty"`
	ApprovedBy   uuid.UUID `json:"approved_by,omitempty"`
	ApprovalDate time.Time `json:"approval_date,omitempty" binding:"required"`
	VisitedFile  uuid.UUID `json:"visited_file,omitempty" binding:"required"` // File of visited location
	Rate         float64   `json:"rate,omitempty" binding:"required"`         // Interest rate
}

type LoanInvestment struct {
	InvestmentId uuid.UUID `json:"investment_id,omitempty"`
	LoanId       uuid.UUID `json:"loan_id,omitempty"`
	InvestedBy   uuid.UUID `json:"invested_by,omitempty"`
	Amount       float64   `json:"amount,omitempty" binding:"required"`
	ROI          float64   `json:"roi,omitempty"` // Return of Investment
}

type LoanDisbursement struct {
	DisbursmentId   uuid.UUID `json:"disbursment_id,omitempty"`
	LoanId          uuid.UUID `json:"loan_id,omitempty"`
	DisbursementBy  uuid.UUID `json:"disbursement_by,omitempty"`
	DisbursedFile   uuid.UUID `json:"disbursed_file,omitempty" binding:"required"` // Signed aggrement file
	DisbursmentDate time.Time `json:"disbursment_date,omitempty" binding:"required"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

type InvestorProfit struct {
	LoanID        string  `json:"loan_id"`
	EmailInvestor string  `json:"email_investor"`
	TotalProfit   float64 `json:"total_profit"`
	ROI           float64 `json:"roi"`
}

type BorrowerPayment struct {
	LoanID        string  `json:"loan_id"`
	Principal     float64 `json:"principal"`
	TotalInterest float64 `json:"total_interest"`
	TotalPayment  float64 `json:"total_payment"`
}
