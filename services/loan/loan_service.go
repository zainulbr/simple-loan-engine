package loan

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/libs/notification/mail"
	"github.com/zainulbr/simple-loan-engine/libs/report/pdf"
	"github.com/zainulbr/simple-loan-engine/libs/template"
	"github.com/zainulbr/simple-loan-engine/models/filemanager"
	"github.com/zainulbr/simple-loan-engine/models/loan"
	repoFile "github.com/zainulbr/simple-loan-engine/repositories/filemanager"
	loanRepository "github.com/zainulbr/simple-loan-engine/repositories/loan"
)

// Implementation
type loanService struct {
	loanRepo loanRepository.LoanRepository
	fileRepo repoFile.FileRepository
	basePath string // Directory untuk menyimpan file
}

// Constructor
func NewLoanService(loanRepo loanRepository.LoanRepository, fileRepo repoFile.FileRepository) LoanService {
	return &loanService{
		loanRepo: loanRepo,
		fileRepo: fileRepo,
		basePath: "./reports",
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

	// Check if principal == total all investment
	if (loanDetail.TotalInvestment + investment.Amount) == investment.Amount {
		// async no blocking
		go s.publishAggrementLatter(context.Background(), investment.LoanId)
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

func (s *loanService) getEmailInvestors(ctx context.Context, loanId uuid.UUID) ([]string, error) {
	emails, err := s.loanRepo.GetInvestorEmailsByLoanID(ctx, loanId)
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func (s *loanService) genReport(ctx context.Context, loanId uuid.UUID) (string, error) {
	pathFile := path.Join(s.basePath, loanId.String()+".pdf")
	pdf.NewService().GeneratePDF(pdf.ReportParam{
		FilePath: pathFile,
		Data:     pdf.Data{},
	})
	fileDetail := &filemanager.File{
		FileType:     ".pdf",
		Label:        loanId.String() + ".pdf",
		Location:     pathFile,
		LocationType: filemanager.LocationTypeLocal,
	}
	err := s.fileRepo.Create(ctx, fileDetail)
	if err != nil {
		return "", err

	}
	return fileDetail.FileID.String(), nil
}

func (s *loanService) genReportLink(id string) string {
	// TBD: set url configuratble
	return "http://localhost:8080/api/files/" + id
}

func (s *loanService) publishAggrementLatter(ctx context.Context, loanId uuid.UUID) error {
	emails, err := s.getEmailInvestors(ctx, loanId)
	if err != nil {
		return err
	}

	// call gen report
	fileId, err := s.genReport(ctx, loanId)
	if err != nil {
		return err
	}
	// get file link
	link := s.genReportLink(fileId)
	fmt.Println(link)
	for _, v := range emails {
		emailBody, err := template.TemplateEmailAgreement(template.EmailData{
			LoanID:       loanId.String(),
			InvestorName: v,
			AgreementURL: link,
		})
		if err != nil {
			fmt.Println(err)
			break
		}
		mail.Mail().Send([]string{v}, "Investment Agreement", emailBody)
	}
	return nil
}

func (s *loanService) TotalPayment(ctx context.Context, loanID string) (*loan.BorrowerPayment, error) {
	if loanID == "" {
		return nil, errors.New("loan_id is required")
	}

	totalPayment, err := s.loanRepo.GetTotalPaymentByLoanID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	return totalPayment, nil
}

func (s *loanService) GetInvestorProfit(ctx context.Context, loanID string) ([]loan.InvestorProfit, error) {
	if loanID == "" {
		return nil, errors.New("loan_id is required")
	}

	listInvestors, err := s.loanRepo.GetInvestorProfitList(ctx, loanID)
	if err != nil {
		return nil, err
	}
	return listInvestors, nil
}
