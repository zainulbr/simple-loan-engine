package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zainulbr/simple-loan-engine/controllers"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/middlewares"
	"github.com/zainulbr/simple-loan-engine/models/user"
	"github.com/zainulbr/simple-loan-engine/repositories/loan"
	loanService "github.com/zainulbr/simple-loan-engine/services/loan"
)

func RegisterLoanRoutes(router *gin.RouterGroup) {

	group := router.Group("/loans")
	group.Use(middlewares.AuthorizeJWT())

	// Initialize Repositories
	loanRepo := loan.NewLoanRepository(pgsql.DB())

	// Initialize Services
	loanService := loanService.NewLoanService(loanRepo)

	// Initialize Controllers
	loanController := controllers.NewLoanController(loanService)

	group.POST("",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.CreateLoan,
	)

	group.GET("/:id",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.GetLoanDetail)

	group.POST("/:id/approve",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.ApproveLoan)

	group.POST("/:id/invest",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.CreateInvestment)

	group.POST("/:id/disburse",
		middlewares.RolePermission(user.RoleBorrower),
		loanController.CreateDisbursement)

}
