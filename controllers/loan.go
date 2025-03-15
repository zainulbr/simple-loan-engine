package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/middlewares"
	"github.com/zainulbr/simple-loan-engine/models/loan"
	models "github.com/zainulbr/simple-loan-engine/models/loan"
	services "github.com/zainulbr/simple-loan-engine/services/loan"
)

// loanController struct
type loanController struct {
	loanService services.LoanService
}

// NewLoanController creates a new instance of loanController
func NewLoanController(loanService services.LoanService) *loanController {
	return &loanController{
		loanService: loanService,
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

	request := loan.LoanApproval{LoanId: loanID}

	userId, ok := c.getUserId(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	request.ApprovedBy = userId

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.loanService.ApproveLoan(ctx.Request.Context(), &request)
	if err != nil {
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

	var disbursement loan.LoanDisbursement
	if err := ctx.ShouldBindJSON(&disbursement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	disbursement.LoanId = loanID
	disbursement.DisbursementBy = userId

	err = c.loanService.CreateDisbursement(ctx.Request.Context(), &disbursement)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan disbursed successfully"})
}
