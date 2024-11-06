package server

import (
	"fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/HersheyPlus/go-auth/api/routes"
    "github.com/HersheyPlus/go-auth/config"
    "github.com/HersheyPlus/go-auth/database"
)


type Server struct {
    Config *config.Config
    Router *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
    return &Server{
        Config: cfg,
    }
}

func (s *Server) setUpRouter() {
	if s.Config.App.Environment == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware(s.Config.CORS))

	if s.Config.RateLimit.Enabled {
		router.Use(RateLimiter(s.Config.RateLimit))
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method not allowed"})
	})
	db := database.GetDB()
	routes.SetUpRoutes(router, db, s.Config)
	s.Router = router
}

func (s *Server) RunServer() error {
	s.setUpRouter()
	server := &http.Server{
        Addr:           fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port),
        Handler:        s.Router,
        ReadTimeout:    s.Config.Server.Timeout.Read,
        WriteTimeout:   s.Config.Server.Timeout.Write,
        IdleTimeout:    s.Config.Server.Timeout.Idle,
        MaxHeaderBytes: 1 << 20,
    }
	log.Printf("Server starting on %s:%s", s.Config.Server.Host, s.Config.Server.Port)
	return server.ListenAndServe()
}