package pdf

import (
	"os"
	"testing"
	"time"
)

func TestGeneratePDF(t *testing.T) {
	filename := "test_loan_summary.pdf"
	input := ReportParam{
		FilePath: filename,
		Data: Data{
			Description:         "TBD Description",
			LoanNumber:          "TBD Loan Number",
			BorrowerName:        "TBD Borrower Name",
			PrincipalAmount:     1000000,
			Duration:            12,
			RatePerYear:         10, // in percentage
			TotalPrincipal:      1100000,
			InvestorName:        "TBD Investor Name",
			InvestmentAmount:    1000000,
			ReturnOfInvestment:  5, // in percentage
			TotalInvestorProfit: 1050000,
			SubmissionDate:      time.Now().AddDate(0, 0, -7),
			ApprovalDate:        time.Now().AddDate(0, 0, -12),
		},
	}

	err := NewService().GeneratePDF(input)
	if err != nil {
		t.Fatalf("Failed to generate PDF: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("PDF file was not created: %s", input.FilePath)
	}

	// Clean up
	os.Remove(input.FilePath)
}
