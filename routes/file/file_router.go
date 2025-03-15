package file

import (
	"github.com/gin-gonic/gin"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/middlewares"
	"github.com/zainulbr/simple-loan-engine/registry"
	"github.com/zainulbr/simple-loan-engine/repositories/filemanager"
	fileService "github.com/zainulbr/simple-loan-engine/services/filemanager"
)

func NewFile() registry.Router {
	// Initialize Repositories
	fileRepo := filemanager.NewFileRepository(pgsql.DB())
	// Initialize Services
	fileManagerService := fileService.NewFileService(fileRepo, "./uploads")

	return &fileHandler{fileService: fileManagerService}
}

func (h *fileHandler) RegisterRoutes(router *gin.RouterGroup) {

	group := router.Group("/files")
	group.Use(middlewares.AuthorizeJWT())

	// TBD: need one time token for email
	group.GET("/:id",
		h.Get)

}

func init() {
	registry.RegisterRouter(NewFile)

}
