package handler

import (
    "github.com/britinogn/bizkeeper/internal/services"
    "github.com/britinogn/bizkeeper/pkg/response"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type ReorderHandler struct {
    reorderService *services.ReorderService
}

func NewReorderHandler(reorderService *services.ReorderService) *ReorderHandler {
    return &ReorderHandler{reorderService: reorderService}
}

func (h *ReorderHandler) GetReorderReminders(c *gin.Context) {
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

    reminders, err := h.reorderService.GetReorderReminders(c.Request.Context(), parsedUserID)
    if err != nil {
        response.InternalServerError(c, "something went wrong")
        return
    }

    response.OK(c, "reorder reminders fetched successfully", reminders)
}