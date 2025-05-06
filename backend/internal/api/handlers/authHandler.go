package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/services"
)

type AuthHandler struct {
	libraryService *services.LibraryService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully registered",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
	})
}
