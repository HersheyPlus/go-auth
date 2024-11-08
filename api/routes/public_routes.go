package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/HersheyPlus/go-auth/api/handlers"
    "gorm.io/gorm"
	 "github.com/HersheyPlus/go-auth/config"
)

func PublicRoutes(default_route *gin.RouterGroup, db *gorm.DB, cfg *config.Config){

	authHandler := handlers.NewAuthHandler(db, cfg)
	public := default_route.Group("/public")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}
}