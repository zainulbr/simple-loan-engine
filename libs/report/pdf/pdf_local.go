package pdf

import (
	"github.com/jung-kurt/gofpdf"
)

type ReportParam struct {
	FilePath string
	Data     Data
}

type Data struct {
	Applicant      string
	Investor       string
	ROI            string
	Rate           string
	Duration       string
	ImportantDates []string
}

type pdfService struct{}

// NewService creates a new PDF service
func NewService() PDFService {
	return &pdfService{}
}

// GeneratePDF generates a PDF file
func (s *pdfService) GeneratePDF(param ReportParam) error {
	data := param.Data
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Loan Application Summary")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Applicant Name: ")
	pdf.Cell(0, 10, data.Applicant)
	pdf.Ln(8)

	pdf.Cell(40, 10, "Investor: ")
	pdf.Cell(0, 10, data.Investor)
	pdf.Ln(8)

	pdf.Cell(40, 10, "ROI: ")
	pdf.Cell(0, 10, data.ROI)
	pdf.Ln(8)

	pdf.Cell(40, 10, "Rate: ")
	pdf.Cell(0, 10, data.Rate)
	pdf.Ln(8)

	pdf.Cell(40, 10, "Loan Duration: ")
	pdf.Cell(0, 10, data.Duration)
	pdf.Ln(8)

	pdf.Cell(40, 10, "Important Dates: ")
	pdf.Ln(6)
	for _, date := range data.ImportantDates {
		pdf.Cell(10, 10, "- ")
		pdf.Cell(0, 10, date)
		pdf.Ln(6)
	}

	return pdf.OutputFileAndClose(param.FilePath)
}
