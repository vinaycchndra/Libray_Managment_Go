package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/services"
)

type AdminHandler struct {
	libraryService *services.LibraryService
}

func NewAdminHandler(libService *services.LibraryService) *AdminHandler {
	return &AdminHandler{
		libraryService: libService,
	}
}

func (h *AdminHandler) getBook(c *gin.Context) {
	type requestBody struct {
		BookId int `json:"book_id"`
	}

	var request_body requestBody

	err := c.BindJSON(&request_body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "", "error": err.Error()})
		return
	}

}
