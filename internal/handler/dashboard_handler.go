package handler

import (
	"github.com/britinogn/bizkeeper/internal/services"
	"github.com/britinogn/bizkeeper/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
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

	summary, err := h.dashboardService.GetDashboardSummary(c.Request.Context(), parsedUserID)
	if err != nil {
		response.InternalServerError(c, "something went wrong")
		return
	}

	response.OK(c, "dashboard summary fetched successfully", summary)
}

func (h *DashboardHandler) GetPriceHistory(c *gin.Context) {
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

	history, err := h.dashboardService.GetPriceHistory(c.Request.Context(), parsedUserID)
	if err != nil {
		response.InternalServerError(c, "something went wrong")
		return
	}

	response.OK(c, "price history fetched successfully", history)
}


// admin dashboard only
func (h *DashboardHandler) GetAdminDashboard(c *gin.Context) {
	stats, err := h.dashboardService.GetAdminDashboard(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "something went wrong")
		return
	}

	response.OK(c, "admin dashboard fetched successfully", stats)
}