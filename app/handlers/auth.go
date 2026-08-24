package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"city-pulse/app/services"
)

type registerRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Username string `json:"username" binding:"required,min=2,max=30"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := h.userSvc.Register(h.ctx, req.Email, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kayıt başarısız"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kayıt başarılı",
		"user": gin.H{
			"id":       user.ID.Hex(),
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	token, user, err := h.userSvc.Login(h.ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCreds) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "giriş başarısız"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func formatValidationError(err error) string {
	return err.Error()
}
