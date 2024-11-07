package routes

import (
	"github.com/HersheyPlus/go-auth/api/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProtectedRoutes(default_route *gin.RouterGroup, db *gorm.DB){
	protected := default_route.Group("/protected")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "profile"})
		})
	}
}