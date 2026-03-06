package routes

import (
	"github.com/britinogn/bizkeeper/internal/handler"
	"github.com/britinogn/bizkeeper/internal/middleware"
	"github.com/gin-gonic/gin"
)

func DashboardRoutes(rg *gin.RouterGroup, dashboardHandler *handler.DashboardHandler) {
	dashboard := rg.Group("/dashboard")
	{
		dashboard.GET("/summary", dashboardHandler.GetDashboardSummary)
		dashboard.GET("/price-history", dashboardHandler.GetPriceHistory)
		dashboard.GET("/admin", middleware.AdminOnly("admin"), dashboardHandler.GetAdminDashboard)

	}
}

