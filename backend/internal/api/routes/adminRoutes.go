package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/api/handlers"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/api/middlewares"
)

func SetupAdminRoutes(router *gin.RouterGroup, handler *handlers.AdminHandler) {
	adminRouter := router.Group("/admin")
	adminRouter.Use(middlewares.AuthMiddleware())
	adminRouter.POST("/add-author", handler.InsertAuthor)
	adminRouter.GET("/get-author", handler.GetAuthor)
	adminRouter.GET("/get-book", handler.GetBookWithBookId)
	adminRouter.POST("/add-book", handler.InsertBook)
	adminRouter.PUT("/update-book/:book_id", handler.UpdateBook)
}
