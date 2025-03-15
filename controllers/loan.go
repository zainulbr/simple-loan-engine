package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/middlewares"
	"github.com/zainulbr/simple-loan-engine/models/filemanager"
	"github.com/zainulbr/simple-loan-engine/models/loan"
	models "github.com/zainulbr/simple-loan-engine/models/loan"
	fmServices "github.com/zainulbr/simple-loan-engine/services/filemanager"
	services "github.com/zainulbr/simple-loan-engine/services/loan"
)

// loanController struct
type loanController struct {
	loanService        services.LoanService
	fileManagerService fmServices.FileService
}

// NewLoanController creates a new instance of loanController
func NewLoanController(loanService services.LoanService,
	fileManagerService fmServices.FileService) *loanController {
	return &loanController{
		loanService:        loanService,
		fileManagerService: fileManagerService,
	}
}

func (c *loanController) getUserId(ctx *gin.Context) (uuid.UUID, bool) {
	userIdString, ok := middlewares.GetClaim(ctx)["loan.user_id"].(string)
	if !ok {
		return uuid.Nil, false
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.Nil, false
	}
	return userId, true

}

// Create Loan (POST /loans)
func (c *loanController) CreateLoan(ctx *gin.Context) {
	var loan models.Loan
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.getUserId(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	loan.ProposedBy = userId
	createdLoan, err := c.loanService.CreateLoan(ctx.Request.Context(), &loan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdLoan)
}

// Approve Loan (POST /loans/:id/approve)
func (c *loanController) ApproveLoan(ctx *gin.Context) {
	loanID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	userId, ok := c.getUserId(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	approvedDateStr := ctx.PostForm("approval_date")
	approvedRateStr := ctx.PostForm("rate")

	if approvedDateStr == "" || approvedRateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "approval_date & rate are required"})
		return
	}

	// Parse approved_date (menggunakan format ISO 8601 atau RFC3339)
	approvedDate, err := time.Parse(time.RFC3339, approvedDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid approved_date format, expected RFC3339 (e.g., 2024-03-13T15:04:05Z)"})
		return
	}

	// Parse approved_date (menggunakan format ISO 8601 atau RFC3339)
	approvalRate, err := strconv.ParseFloat(approvedRateStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid rate format, expected float (e.g., 0.1)"})
		return
	}

	// Handle file upload
	file, header, err := ctx.Request.FormFile("visited_file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "visited_file is required"})
		return
	}
	defer file.Close()

	// Validate file format
	if err := c.fileManagerService.ValidateFileFormat(file, header); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save file to disk and detail on db
	// TBD: Need rollback delete when approval loan failed
	fileDetail, err := c.fileManagerService.UploadFile(ctx.Request.Context(),
		file,
		header,
		filemanager.LocationTypeLocal,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	request := loan.LoanApproval{LoanId: loanID}
	request.ApprovedBy = userId
	request.ApprovalDate = approvedDate
	request.VisitedFile = fileDetail.FileID
	request.Rate = approvalRate
	err = c.loanService.ApproveLoan(ctx.Request.Context(), &request)
	if err != nil {
		// TBD: refactore validate status first
		go c.fileManagerService.DeleteFile(context.Background(), fileDetail.FileID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan approved successfully"})
}

// Get Loan Detail (GET /loans/:id)
func (c *loanController) GetLoanDetail(ctx *gin.Context) {
	loanID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	loanDetail, err := c.loanService.GetLoanDetail(ctx.Request.Context(), loanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loanDetail)
}

// // Create Investment (POST /loans/:id/invest)
func (c *loanController) CreateInvestment(ctx *gin.Context) {
	loanID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	userId, ok := c.getUserId(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var investment loan.LoanInvestment
	if err := ctx.ShouldBindJSON(&investment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	investment.LoanId = loanID
	investment.InvestedBy = userId

	err = c.loanService.CreateInvestment(ctx.Request.Context(), &investment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Investment created successfully"})
}

// Create Disbursement (POST /loans/:id/disburse)
func (c *loanController) CreateDisbursement(ctx *gin.Context) {
	loanID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	userId, ok := c.getUserId(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	disbursementDateStr := ctx.PostForm("disbursment_date")

	if disbursementDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "disbursment_date are required"})
		return
	}

	// Parse approved_date (menggunakan format ISO 8601 atau RFC3339)
	disbursmentDate, err := time.Parse(time.RFC3339, disbursementDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid approved_date format, expected RFC3339 (e.g., 2024-03-13T15:04:05Z)"})
		return
	}

	// Handle file upload
	file, header, err := ctx.Request.FormFile("disbursed_file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "disbursed_file is required"})
		return
	}
	defer file.Close()

	// Validate file format
	if err := c.fileManagerService.ValidateFileFormat(file, header); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save file to disk and detail on db
	// TBD: Need rollback delete when approval loan failed
	fileDetail, err := c.fileManagerService.UploadFile(ctx.Request.Context(),
		file,
		header,
		filemanager.LocationTypeLocal,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	disbursement := loan.LoanDisbursement{LoanId: loanID}
	disbursement.DisbursementBy = userId
	disbursement.DisbursmentDate = disbursmentDate
	disbursement.DisbursedFile = fileDetail.FileID

	err = c.loanService.CreateDisbursement(ctx.Request.Context(), &disbursement)
	if err != nil {
		// TBD: refactore validate status first
		go c.fileManagerService.DeleteFile(context.Background(), fileDetail.FileID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan disbursed successfully"})
}

// GetTotalPaymentHandler (GET /loans/:id/total-interest)
func (c *loanController) GetTotalPayment(ctx *gin.Context) {
	loanID := ctx.Param("id")

	totalPayment, err := c.loanService.TotalPayment(ctx.Request.Context(), loanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": totalPayment})
}

// GetInvestorProfitList (GET /loans/:id/profit-investor)
func (lc *loanController) GetInvestorProfitList(c *gin.Context) {
	loanID := c.Param("id")

	profits, err := lc.loanService.GetInvestorProfit(c.Request.Context(), loanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profits)
}
