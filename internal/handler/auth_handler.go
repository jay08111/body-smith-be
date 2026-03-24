package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"body-smith-be/internal/model"
	"body-smith-be/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			respondError(c, http.StatusUnauthorized, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to login")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Me(c *gin.Context) {
	emailValue, exists := c.Get("userEmail")
	if !exists {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	email, ok := emailValue.(string)
	if !ok || email == "" {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.authService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to load user")
		return
	}
	if user == nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": model.UserSummary{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
