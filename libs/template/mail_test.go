package template

import (
	"fmt"
	"log"
	"testing"
)

func Test(t *testing.T) {
	// Data yang akan dimasukkan ke dalam template
	data := EmailData{
		InvestorName: "John Doe",
		LoanID:       "12345-67890",
		AgreementURL: "https://yourplatform.com/agreement/12345.pdf",
	}

	// Generate email template
	emailBody, err := TemplateEmailAgreement(data)
	if err != nil {
		log.Fatal("Error generating email:", err)
	}

	// Cetak hasil email
	fmt.Println(emailBody)
}
