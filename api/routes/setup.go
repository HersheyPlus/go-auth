package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/HersheyPlus/go-auth/config"
	"gorm.io/gorm"
)

func SetUpRoutes(router *gin.Engine, db *gorm.DB, config *config.Config){}