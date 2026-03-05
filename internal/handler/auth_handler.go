package handler

import (
	"net/http"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/britinogn/bizkeeper/internal/services"
	"github.com/britinogn/bizkeeper/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	ctx := c.Request.Context()
	user, err := h.authService.Register(ctx, &req)
	if err != nil {
		switch err {
		case services.ErrEmailAlreadyRegistered:
			response.Error(c, http.StatusConflict, err.Error())
		case services.ErrMissingRequiredFields:
			response.BadRequest(c, err.Error())
		case services.ErrWeakPassword:
			response.BadRequest(c, err.Error())
		case services.ErrInvalidInput:
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.Created(c, "account created successfully", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			response.Unauthorized(c, err.Error())
		case services.ErrMissingLoginFields:
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.OK(c, "login successful", gin.H{"user": user, "token": token})
}



func (h *AuthHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	user, err := h.authService.UpdateUser(c.Request.Context(), userID.(string), &req)
	if err != nil {
		switch err {
		case services.ErrWeakPassword:
			response.BadRequest(c, err.Error())
		case services.ErrInvalidInput:
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.OK(c, "user updated successfully", user)
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	if err := h.authService.DeleteUser(c.Request.Context(), userID.(string)); err != nil {
		response.InternalServerError(c, "something went wrong")
		return
	}

	response.OK(c, "account deleted successfully", nil)
}

func (h *AuthHandler) GetProfile (c *gin.Context){
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			response.NotFound(c, err.Error())
		default:
			response.InternalServerError(c, "something went wrong")
		}
		return
	}

	response.OK(c, "profile fetched successfully", user)
}