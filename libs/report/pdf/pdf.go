package pdf

// PDFService is an interface to generate PDF
type PDFService interface {
	GeneratePDF(param ReportParam) error
}
