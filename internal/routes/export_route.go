package routes

import (
	"github.com/britinogn/bizkeeper/internal/handler"
	"github.com/gin-gonic/gin"
)

func ExportRoutes(rg *gin.RouterGroup, exportHandler *handler.ExportHandler) {
	export := rg.Group("/export")
	{
		export.GET("", exportHandler.Export)
	}
}