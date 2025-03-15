package template

import (
	"bytes"
	"html/template"
)

type EmailData struct {
	InvestorName string
	LoanID       string
	AgreementURL string
}

// TemplateEmailAgreement function for render template email
func TemplateEmailAgreement(data EmailData) (string, error) {
	// Template format HTML
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Investment Agreement</title>
		</head>
		<body>
			<p>Dear {{ .InvestorName }},</p>
			<p>Congratulations! The investment for Loan ID <strong>{{ .LoanID }}</strong> has been fully funded.</p>
			<p>You can review and download the investment agreement using the link below:</p>
			<p><a href="{{ .AgreementURL }}" target="_blank">Download Agreement PDF</a></p>
			<p>Thank you for your trust in our platform.</p>
			<br>
			<p>Best Regards,</p>
			<p><strong>Your Loan Platform Team</strong></p>
		</body>
		</html>
	`

	// Parse template
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	// Render template with data
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
