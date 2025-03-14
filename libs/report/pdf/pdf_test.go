package pdf

import (
	"os"
	"testing"
)

func TestGeneratePDF(t *testing.T) {
	filename := "test_loan_summary.pdf"
	input := ReportParam{
		FilePath: filename,
		Data: Data{
			Applicant:      "Test User",
			Investor:       "Test Investor",
			ROI:            "10%",
			Rate:           "4%",
			Duration:       "12 months",
			ImportantDates: []string{"Application Date: 01-Feb-2024", "Approval Date: 10-Feb-2024"},
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
