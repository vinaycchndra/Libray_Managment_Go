package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/api/handlers"
)

func LoginRoutes(router *gin.RouterGroup, handler *handlers.AdminHandler) {
	router.POST("/login", handler.GetBookWithBookId)
}
