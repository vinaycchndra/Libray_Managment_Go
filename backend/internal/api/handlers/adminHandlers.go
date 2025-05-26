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

func (h *AdminHandler) GetBookWithBookId(c *gin.Context) {
	type requestBody struct {
		BookId int `json:"book_id"`
	}
	var request_body requestBody

	err := c.BindJSON(&request_body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "", "error": err.Error()})
		return
	}

	book, err := h.libraryService.GetBook(request_body.BookId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

type authorRequestBody struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	About string `json:"about,omitempty"`
}

func (h *AdminHandler) InsertAuthor(c *gin.Context) {
	var request_body authorRequestBody

	err := c.BindJSON(&request_body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "", "error": err.Error()})
		return
	}

	if request_body.Name == "" || request_body.About == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "", "error": "name and about is mandatory to insert the author."})
		return
	}

	author, err := h.libraryService.InsertAuthor(request_body.Name, request_body.About)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, author)
}

func (h *AdminHandler) GetAuthor(c *gin.Context) {
	var request_body authorRequestBody

	err := c.BindJSON(&request_body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "", "error": err.Error()})
		return
	}

	if request_body.ID != 0 {
		authors, err := h.libraryService.GetAuthor(request_body.ID, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, authors)
		return
	} else {
		authors, err := h.libraryService.GetAuthor(0, request_body.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, authors)
		return
	}

}
