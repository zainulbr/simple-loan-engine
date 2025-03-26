package pdf

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type ReportParam struct {
	FilePath string
	Data     Data
}

type Data struct {
	Description         string
	LoanNumber          string
	BorrowerName        string
	PrincipalAmount     float64
	Duration            int
	RatePerYear         float64
	TotalPrincipal      float64
	InvestorName        string
	InvestmentAmount    int
	ReturnOfInvestment  float64
	TotalInvestorProfit float64
	SubmissionDate      time.Time
	ApprovalDate        time.Time
}

type pdfService struct{}

// NewService creates a new PDF service
func NewService() PDFService {
	return &pdfService{}
}

// GeneratePDF generates a PDF file
func (s *pdfService) GeneratePDF(param ReportParam) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Buat halaman baru
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Agreement Letter")
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 10, param.Data.Description)
	pdf.Ln(8)
	pdf.Cell(190, 10, fmt.Sprintf("Loan Number : %s", param.Data.LoanNumber))
	pdf.Ln(12)

	// Principal Details
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Loan Principal Details")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, fmt.Sprintf("Borrower Name: %s", param.Data.BorrowerName))
	pdf.Ln(6)
	pdf.Cell(190, 10, fmt.Sprintf("Principal Amount %.2f ", param.Data.PrincipalAmount))
	pdf.Ln(6)
	pdf.Cell(190, 10, fmt.Sprintf("Duration(month): %d ", param.Data.Duration))
	pdf.Ln(6)
	pdf.Cell(190, 10, fmt.Sprintf("Rate per year: %.2f ", param.Data.RatePerYear))
	pdf.Ln(6)
	pdf.Cell(190, 10, fmt.Sprintf("Total Prinsipal Should Pay: %.2f", param.Data.TotalPrincipal))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Loan Investor Details")
	pdf.Ln(8)

	// Investor Details
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, fmt.Sprintf("Investor Name: %s", param.Data.InvestorName))
	pdf.Ln(6)
	pdf.Cell(190, 10, fmt.Sprintf("Investment Amount: %d", param.Data.InvestmentAmount))
	pdf.Ln(6)
	pdf.Cell(190, 10, fmt.Sprintf("Return Of Investment: %.2f", param.Data.ReturnOfInvestment))
	pdf.Ln(6)
	pdf.Cell(190, 10, fmt.Sprintf("Total Investor Profit: %.2f", param.Data.TotalInvestorProfit))
	pdf.Ln(10)

	// Tanggal
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, fmt.Sprintf("Submission Date: %s", param.Data.SubmissionDate.Format("2006-01-02")))
	pdf.Ln(6)
	pdf.Cell(190, 10, fmt.Sprintf("Approval Date: %s", param.Data.ApprovalDate.Format("2006-01-02")))
	pdf.Ln(10)

	return pdf.OutputFileAndClose(param.FilePath)
}
