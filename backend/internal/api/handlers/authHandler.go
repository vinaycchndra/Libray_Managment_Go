package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/services"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/utils"
)

type AuthHandler struct {
	libraryService *services.LibraryService
}

func NewAuthHandler(libraryService *services.LibraryService) *AuthHandler {
	return &AuthHandler{
		libraryService: libraryService,
	}
}

func (a *AuthHandler) Register(c *gin.Context) {
	type requestBody struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	var request_body requestBody

	err := c.BindJSON(&request_body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "", "error": err.Error()})
		return
	}

	// password validation and hashing
	err = utils.ValidatePassword(request_body.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "", "error": err.Error()})
		return
	}

	registered_user, err := a.libraryService.RegisterUser(
		request_body.Name,
		request_body.Email,
		request_body.Password,
		request_body.PhoneNumber,
		false,
		false,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%v Successfully registered", registered_user.Email),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	type requestBody struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var request_body requestBody

	err := c.BindJSON(&request_body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "", "error": err.Error()})
		return
	}

	token, err := h.libraryService.LoginUser(request_body.Email, request_body.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "", "token": "", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"token":   token,
		"error":   "",
	})
}
