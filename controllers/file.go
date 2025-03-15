package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/services/filemanager"
)

type fileHandler struct {
	fileService filemanager.FileService
}

func NewFileController(fileService filemanager.FileService) *fileHandler {
	return &fileHandler{fileService: fileService}
}

// Get serves file
// Defulat serve as preview
func (h *fileHandler) Get(c *gin.Context) {
	// Extract file ID from URL
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get file path & MIME type
	filePath, mimeType, err := h.fileService.PreviewFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if c.Query("download") == "true" {
		// TBD: Serve file as attachment (force download)
		// c.Header("Content-Disposition", "attachment; filename="+path.Base(filePath))
		// c.Header("Content-Type", "application/octet-stream") // Force download
		// c.FileAttachment(filePath, path.Base(filePath))
		// return
	}

	// Set MIME type and serve the file
	c.Header("Content-Type", mimeType)
	c.File(filePath)
}
