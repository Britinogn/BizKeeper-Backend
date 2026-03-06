package handler

import (
	"fmt"
	"time"

	"github.com/britinogn/bizkeeper/internal/services"
	"github.com/britinogn/bizkeeper/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExportHandler struct {
	exportService *services.ExportService
}

func NewExportHandler(exportService *services.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

func (h *ExportHandler) Export(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	parsedUserID, err := uuid.Parse(userID.(string))
	if err != nil {
		response.Unauthorized(c, "invalid user ID")
		return
	}

	format := c.Query("format")
	rangeStr := c.Query("range")

	if format == "" || rangeStr == "" {
		response.BadRequest(c, "format and range are required")
		return
	}

	if format != "csv" && format != "pdf" {
		response.BadRequest(c, "format must be csv or pdf")
		return
	}

	switch format {
	case "csv":
		data, err := h.exportService.ExportCSV(c.Request.Context(), parsedUserID, rangeStr)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		filename := fmt.Sprintf("bizkeeper-export-%s.csv", time.Now().Format("2006-01-02"))
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(200, "text/csv", data)

	case "pdf":
		data, err := h.exportService.ExportPDF(c.Request.Context(), parsedUserID, rangeStr)
		if err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		filename := fmt.Sprintf("bizkeeper-export-%s.pdf", time.Now().Format("2006-01-02"))
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(200, "application/pdf", data)
	}
}