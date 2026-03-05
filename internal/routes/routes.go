package routes

import (
	"net/http"

	"github.com/britinogn/bizkeeper/internal/handler"
	"github.com/britinogn/bizkeeper/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	purchaseHandler *handler.PurchaseHandler,
	dashboardHandler *handler.DashboardHandler,
) {
	api := router.Group("/api")

	// Health
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "BizKeeper API is running",
		})
	})

	// Public
	public := api.Group("")

	// Protected
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	// _ = protected // add this line
	// protected.Use(middleware.AdminOnly())

	// Register separated routes
	AuthRoutes(public, authHandler)
	UserRoutes(protected, authHandler)
	PurchaseRoutes(protected, purchaseHandler)
	DashboardRoutes(protected, dashboardHandler)
}
