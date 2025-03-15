package loan

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/models/loan"
	"github.com/zainulbr/simple-loan-engine/settings"

	"github.com/go-pg/pg/v10"
)

func setupTestDB(t *testing.T) *pg.DB {
	config := settings.Settings{
		Conn: settings.ConnectionSettings{
			Postgres: settings.PostgresOption{
				URI:     "postgres://admin:admin@localhost:5432/loan-db?sslmode=disable",
				Enabled: true,
			},
		},
	}

	if err := pgsql.Open(&config); err != nil {
		t.Error(err)
	}

	return pgsql.DB()
}

func TestLoan(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	proposedby, _ := uuid.Parse("f4fcbdce-d833-48d1-bc98-af05aa29ad1d")
	approvedBy, _ := uuid.Parse("e6df5aef-090b-4ec0-974f-03f58abc00b2")
	investedBy1, _ := uuid.Parse("99385f18-5655-432c-a383-163e9e634b0d")
	// investedBy2, _ := uuid.Parse("f4fcbdce-d833-48d1-bc98-af05aa29ad1d")
	disbursedBy, _ := uuid.Parse("c76a1d43-5758-4e05-9f08-a16e991bc656")

	fileId, _ := uuid.Parse("32141712-93f5-4b53-a0f1-b266c1da60de")

	data := &loan.Loan{
		Amount:        1000000,
		ProposedBy:    proposedby,
		DurationMonth: 12,
		Description:   "Test Loan",
	}
	svc := NewLoanRepository(db)

	err := svc.CreateLoan(context.Background(), data)
	assert.NoError(t, err)
	assert.NotEmpty(t, data.LoanId)
	t.Logf("Created  ID: %s", data.LoanId)

	resp, err := svc.GetLoanDetail(context.Background(), data.LoanId)
	assert.NoError(t, err)
	assert.Equal(t, resp.Amount, data.Amount)
	assert.Equal(t, resp.ProposedBy, data.ProposedBy)
	assert.Equal(t, resp.DurationMonth, data.DurationMonth)
	assert.Equal(t, resp.LoanId, data.LoanId)
	assert.NotEmpty(t, resp.CreatedAt)
	assert.NotEmpty(t, resp.UpdatedAt)
	t.Logf("Retrieved : %+v", resp)

	// approval loan
	err = svc.Approve(context.Background(),
		&loan.LoanApproval{
			LoanId:       data.LoanId,
			VisitedFile:  fileId,
			ApprovedBy:   approvedBy,
			ApprovalDate: time.Now(),
			Rate:         0.1,
		})

	assert.NoError(t, err)

	// Test Investment 1x total investment
	err = svc.CreateInvestment(context.Background(),
		&loan.LoanInvestment{
			LoanId:     data.LoanId,
			Amount:     1000000,
			InvestedBy: investedBy1,
		})

	assert.NoError(t, err)

	// Test re invest again should be error
	err = svc.CreateInvestment(context.Background(),
		&loan.LoanInvestment{
			LoanId:     data.LoanId,
			Amount:     1000000,
			InvestedBy: investedBy1,
		})

	assert.Error(t, err)

	// Test Investment 1x total investment
	err = svc.CreateDisbursement(context.Background(),
		&loan.LoanDisbursement{
			LoanId:          data.LoanId,
			DisbursementBy:  disbursedBy,
			DisbursedFile:   fileId,
			DisbursmentDate: time.Now(),
		})

	assert.NoError(t, err)

}

func TestLoanMultiInvestor(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	proposedby, _ := uuid.Parse("f4fcbdce-d833-48d1-bc98-af05aa29ad1d")
	approvedBy, _ := uuid.Parse("e6df5aef-090b-4ec0-974f-03f58abc00b2")
	investedBy1, _ := uuid.Parse("99385f18-5655-432c-a383-163e9e634b0d")
	investedBy2, _ := uuid.Parse("f4fcbdce-d833-48d1-bc98-af05aa29ad1d")
	disbursedBy, _ := uuid.Parse("c76a1d43-5758-4e05-9f08-a16e991bc656")

	fileId, _ := uuid.Parse("32141712-93f5-4b53-a0f1-b266c1da60de")

	data := &loan.Loan{
		Amount:        1000000,
		ProposedBy:    proposedby,
		DurationMonth: 12,
		Description:   "Test Loan",
	}
	svc := NewLoanRepository(db)

	err := svc.CreateLoan(context.Background(), data)
	assert.NoError(t, err)
	assert.NotEmpty(t, data.LoanId)
	t.Logf("Created  ID: %s", data.LoanId)

	resp, err := svc.GetLoanDetail(context.Background(), data.LoanId)
	assert.NoError(t, err)
	assert.Equal(t, resp.Amount, data.Amount)
	assert.Equal(t, resp.ProposedBy, data.ProposedBy)
	assert.Equal(t, resp.DurationMonth, data.DurationMonth)
	assert.Equal(t, resp.LoanId, data.LoanId)
	assert.NotEmpty(t, resp.CreatedAt)
	assert.NotEmpty(t, resp.UpdatedAt)
	t.Logf("Retrieved : %+v", resp)

	// approval loan
	err = svc.Approve(context.Background(),
		&loan.LoanApproval{
			LoanId:       data.LoanId,
			VisitedFile:  fileId,
			ApprovedBy:   approvedBy,
			ApprovalDate: time.Now(),
			Rate:         0.1,
		})

	assert.NoError(t, err)

	// Test Investment 1x total investment
	err = svc.CreateInvestment(context.Background(),
		&loan.LoanInvestment{
			LoanId:     data.LoanId,
			Amount:     500000,
			InvestedBy: investedBy1,
		})

	assert.NoError(t, err)

	// Test Investment 1x total investment
	err = svc.CreateInvestment(context.Background(),
		&loan.LoanInvestment{
			LoanId:     data.LoanId,
			Amount:     500000,
			InvestedBy: investedBy2,
		})

	assert.NoError(t, err)

	// Test Investment 1x total investment
	err = svc.CreateDisbursement(context.Background(),
		&loan.LoanDisbursement{
			LoanId:          data.LoanId,
			DisbursementBy:  disbursedBy,
			DisbursedFile:   fileId,
			DisbursmentDate: time.Now(),
		})

	assert.NoError(t, err)

}
