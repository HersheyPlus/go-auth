package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/HersheyPlus/go-auth/config"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config){
	r.GET("/health")
	default_route := r.Group(cfg.App.API.Prefix + "/" + cfg.App.API.Version)
	PublicRoutes(default_route, db, cfg)
	ProtectedRoutes(default_route, db)
}