package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zainulbr/simple-loan-engine/controllers"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/middlewares"
	"github.com/zainulbr/simple-loan-engine/repositories/filemanager"
	fileService "github.com/zainulbr/simple-loan-engine/services/filemanager"
)

func RegisterFileRoutes(router *gin.RouterGroup) {

	group := router.Group("/files")
	group.Use(middlewares.AuthorizeJWT())

	// Initialize Repositories
	fileRepo := filemanager.NewFileRepository(pgsql.DB())
	// Initialize Services
	fileManagerService := fileService.NewFileService(fileRepo, "./uploads")

	// Initialize Controllers
	fileController := controllers.NewFileController(
		fileManagerService,
	)

	// TBD: need one time token for email
	group.GET("/:id",
		fileController.Get)

}
