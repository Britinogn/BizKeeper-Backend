package routes

import (
	"github.com/britinogn/bizkeeper/internal/handler"
	"github.com/gin-gonic/gin"
)

// func PurchaseRoutes(rg *gin.RouterGroup, purchaseHandler *handler.PurchaseHandler) {
// 	purchase := rg.Group("/purchase")
// 	{
// 		purchase.POST("/", purchaseHandler.CreatePurchaseSession)
// 		purchase.GET("/", purchaseHandler.ListPurchaseSessions)
// 		purchase.GET("/:id", purchaseHandler.GetPurchaseSessionByID)
// 		purchase.PUT("/:id", purchaseHandler.UpdatePurchaseSession)
// 		purchase.DELETE("/:id", purchaseHandler.DeletePurchaseSession)
// 	}

// 	items := rg.Group("/purchase/:sessionId/items")
// 	{
// 		items.PUT("/:itemId", purchaseHandler.UpdateProductItem)
// 		items.DELETE("/:itemId", purchaseHandler.DeleteProductItem)
// 	}
// }


func PurchaseRoutes(rg *gin.RouterGroup, purchaseHandler *handler.PurchaseHandler) {
	purchase := rg.Group("/purchases")
	{
		purchase.POST("", purchaseHandler.CreatePurchaseSession)
		purchase.GET("", purchaseHandler.ListPurchaseSessions)
		purchase.GET("/:id", purchaseHandler.GetPurchaseSessionByID)
		purchase.PUT("/:id", purchaseHandler.UpdatePurchaseSession)
		purchase.DELETE("/:id", purchaseHandler.DeletePurchaseSession)

		// Product item routes
		purchase.PUT("/:id/items/:itemId", purchaseHandler.UpdateProductItem)
		purchase.DELETE("/:id/items/:itemId", purchaseHandler.DeleteProductItem)
	}
}