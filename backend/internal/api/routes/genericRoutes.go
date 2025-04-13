package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/api/handlers"
)

func SetupGenericRoutes(router *gin.RouterGroup, handler *handlers.GenericHandler) {
	router.GET("/ping", handler.Ping)
}
