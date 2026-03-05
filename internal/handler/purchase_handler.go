package handler

import (
	"log"
	"strconv"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/britinogn/bizkeeper/internal/services"
	"github.com/britinogn/bizkeeper/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PurchaseHandler struct {
	purchaseService *services.PurchaseService
}

func NewPurchaseHandler(purchaseService *services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{purchaseService: purchaseService}
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}
	parsed, err := uuid.Parse(userID.(string))
	if err != nil {
		return uuid.Nil, false
	}
	return parsed, true
}

func (h *PurchaseHandler) CreatePurchaseSession(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req model.PurchaseSession
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	err := h.purchaseService.CreatePurchaseSession(c.Request.Context(), userID, &req)
	if err != nil {
		switch err {
		case services.ErrInvalidSession:
			response.BadRequest(c, err.Error())
		case services.ErrNoProductItems:
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.Created(c, "purchase session created successfully", req)
}

func (h *PurchaseHandler) GetPurchaseSessionByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid session ID")
		return
	}

	session, err := h.purchaseService.GetPurchaseSessionByID(c.Request.Context(), userID, sessionID)
	if err != nil {
		switch err {
		case services.ErrSessionNotFound:
			response.NotFound(c, err.Error())
		case services.ErrUnauthorizedSession:
			response.Unauthorized(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.OK(c, "purchase session fetched successfully", session)
}

func (h *PurchaseHandler) ListPurchaseSessions(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}

	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil || parsed <= 0 {
			response.BadRequest(c, "invalid limit parameter")
			return
		}
		limit = parsed
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsed, err := strconv.Atoi(offsetStr)
		if err != nil || parsed < 0 {
			response.BadRequest(c, "invalid offset parameter")
			return
		}
		offset = parsed
	}

	sessions, err := h.purchaseService.ListPurchaseSessions(c.Request.Context(), userID, limit, offset)
	if err != nil {
		response.InternalServerError(c, "something went wrong")
		return
	}

	response.OK(c, "purchase sessions fetched successfully", sessions)
}

func (h *PurchaseHandler) UpdatePurchaseSession(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid session ID")
		return
	}

	var req model.PurchaseSession
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	err = h.purchaseService.UpdatePurchaseSession(c.Request.Context(), userID, sessionID, &req)
	if err != nil {
		switch err {
		case services.ErrSessionNotFound:
			response.NotFound(c, err.Error())
		case services.ErrUnauthorizedSession:
			response.Unauthorized(c, err.Error())
		case services.ErrInvalidSession:
			response.BadRequest(c, err.Error())
		case services.ErrNoProductItems:
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.OK(c, "purchase session updated successfully", req)
}

func (h *PurchaseHandler) DeletePurchaseSession(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid session ID")
		return
	}

	err = h.purchaseService.DeletePurchaseSession(c.Request.Context(), userID, sessionID)
	if err != nil {
		switch err {
		case services.ErrSessionNotFound:
			response.NotFound(c, err.Error())
		case services.ErrUnauthorizedSession:
			response.Unauthorized(c, err.Error())
		default:
			// response.InternalServerError(c, "something went wrong")
			log.Println("Delete error:", err)
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.OK(c, "purchase session deleted successfully", nil)
}



func (h *PurchaseHandler) UpdateProductItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid session ID")
		return
	}

	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		response.BadRequest(c, "invalid item ID")
		return
	}

	var req model.ProductItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	item, err := h.purchaseService.UpdateProductItem(c.Request.Context(), userID, sessionID, itemID, &req)
	if err != nil {
		switch err {
		case services.ErrSessionNotFound:
			response.NotFound(c, err.Error())
		case services.ErrItemNotFound:
			response.NotFound(c, err.Error())
		case services.ErrUnauthorizedSession:
			response.Unauthorized(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.OK(c, "product item updated successfully", item)
}

func (h *PurchaseHandler) DeleteProductItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid session ID")
		return
	}

	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		response.BadRequest(c, "invalid item ID")
		return
	}

	err = h.purchaseService.DeleteProductItem(c.Request.Context(), userID, sessionID, itemID)
	if err != nil {
		switch err {
		case services.ErrSessionNotFound:
			response.NotFound(c, err.Error())
		case services.ErrItemNotFound:
			response.NotFound(c, err.Error())
		case services.ErrUnauthorizedSession:
			response.Unauthorized(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.OK(c, "product item deleted successfully", nil)
}