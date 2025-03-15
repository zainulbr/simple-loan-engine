package loan

import (
	"github.com/gin-gonic/gin"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/middlewares"
	"github.com/zainulbr/simple-loan-engine/models/user"
	"github.com/zainulbr/simple-loan-engine/registry"
	"github.com/zainulbr/simple-loan-engine/repositories/filemanager"
	"github.com/zainulbr/simple-loan-engine/repositories/loan"
	fileService "github.com/zainulbr/simple-loan-engine/services/filemanager"
	loanService "github.com/zainulbr/simple-loan-engine/services/loan"
)

func NewLoan() registry.Router {
	// Initialize Repositories
	loanRepo := loan.NewLoanRepository(pgsql.DB())
	fileRepo := filemanager.NewFileRepository(pgsql.DB())
	// Initialize Services
	loanService := loanService.NewLoanService(loanRepo, fileRepo)
	fileManagerService := fileService.NewFileService(fileRepo, "./uploads")

	return &loanController{
		loanService:        loanService,
		fileManagerService: fileManagerService,
	}
}

func (c *loanController) RegisterRoutes(router *gin.RouterGroup) {

	group := router.Group("/loans")
	group.Use(middlewares.AuthorizeJWT())

	group.POST("",
		middlewares.RolePermission(user.RoleBorrower),
		c.CreateLoan,
	)

	group.GET("/:id",
		c.GetLoanDetail)

	group.GET("/:id/total-interest",
		middlewares.RolePermission(user.RoleFieldOfficer),
		c.GetTotalPayment)

	group.GET("/:id/profit-investor",
		middlewares.RolePermission(user.RoleFieldOfficer),
		c.GetInvestorProfitList)

	// TBD: Chanege permission to Field Validator
	group.POST("/:id/approve",
		middlewares.RolePermission(user.RoleFiledValidator),
		c.ApproveLoan)

	// TBD: Change permission to investor
	group.POST("/:id/invest",
		middlewares.RolePermission(user.RoleInvestor),
		c.CreateInvestment)
	// TBD: Change permission to Field Officer
	group.POST("/:id/disburse",
		middlewares.RolePermission(user.RoleFieldOfficer),
		c.CreateDisbursement)

}

func init() {
	registry.RegisterRouter(NewLoan)

}
