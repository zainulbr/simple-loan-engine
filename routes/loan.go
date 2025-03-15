package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zainulbr/simple-loan-engine/controllers"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/middlewares"
	"github.com/zainulbr/simple-loan-engine/models/user"
	"github.com/zainulbr/simple-loan-engine/repositories/filemanager"
	"github.com/zainulbr/simple-loan-engine/repositories/loan"
	fileService "github.com/zainulbr/simple-loan-engine/services/filemanager"
	loanService "github.com/zainulbr/simple-loan-engine/services/loan"
)

func RegisterLoanRoutes(router *gin.RouterGroup) {

	group := router.Group("/loans")
	group.Use(middlewares.AuthorizeJWT())

	// Initialize Repositories
	loanRepo := loan.NewLoanRepository(pgsql.DB())
	fileRepo := filemanager.NewFileRepository(pgsql.DB())
	// Initialize Services
	loanService := loanService.NewLoanService(loanRepo, fileRepo)
	fileManagerService := fileService.NewFileService(fileRepo, "./uploads")

	// Initialize Controllers
	loanController := controllers.NewLoanController(
		loanService,
		fileManagerService,
	)

	group.POST("",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.CreateLoan,
	)

	group.GET("/:id",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.GetLoanDetail)

	// TBD: Chanege permission to Field Validator
	group.POST("/:id/approve",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.ApproveLoan)

	// TBD: Change permission to investor
	group.POST("/:id/invest",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.CreateInvestment)
	// TBD: Change permission to Field Officer
	group.POST("/:id/disburse",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.CreateDisbursement)

}
