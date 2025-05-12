package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/api/handlers"
)

func SetupAuthRoutes(router *gin.RouterGroup, handler *handlers.AuthHandler) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", handler.Register)
		// router.POST("/login", handler.Login)
	}
}
