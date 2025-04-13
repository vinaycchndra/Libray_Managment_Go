package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenericHandler struct {
}

func NewGenericHandler() *GenericHandler {
	return &GenericHandler{}
}

func (h *GenericHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Pong",
	})
}
