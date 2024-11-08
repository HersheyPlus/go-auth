package routes

import (
	"github.com/HersheyPlus/go-auth/api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/HersheyPlus/go-auth/config"
	"gorm.io/gorm"
	"github.com/HersheyPlus/go-auth/api/handlers"
)

func ProtectedRoutes(default_route *gin.RouterGroup, db *gorm.DB, cfg *config.Config){
	protected := default_route.Group("/protected")
	protected.Use(middlewares.AuthMiddleware(cfg))
	authHandler := handlers.NewAuthHandler(db, cfg)
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/logout", authHandler.Logout)
	}
}