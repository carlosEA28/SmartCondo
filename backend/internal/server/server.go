package server

import (
	"net/http"

	"github.com/carlosEA28/smartcondo/internal/config"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/carlosEA28/smartcondo/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
}

func New(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		config: cfg,
		db:     db,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Add middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	router.GET("/health", s.healthCheck)

	userRepository := repositories.NewGormUserRepository(s.db)
	userService := services.NewUserService(userRepository)
	userHandler := newUserHandler(userService)
	router.POST("/users", userHandler.create)

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}
