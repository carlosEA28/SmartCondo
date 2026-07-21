package server

import (
	"net/http"

	"github.com/carlosEA28/smartcondo/internal/config"
	providers "github.com/carlosEA28/smartcondo/internal/providers/aws"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/carlosEA28/smartcondo/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	userRepository repositories.UserRepository
}

func New(cfg *config.Config, db *gorm.DB, userRepository repositories.UserRepository) *Server {
	return &Server{
		config:         cfg,
		db:             db,
		userRepository: userRepository,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Add middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	router.GET("/health", s.healthCheck)

	awsProvider := providers.NewAwsProvider(s.config)
	userService := services.NewUserService(s.userRepository, awsProvider)
	userHandler := newUserHandler(userService)
	router.POST("/users", userHandler.create)
	router.GET("/users", userHandler.list)
	router.GET("/users/:id", userHandler.getByID)
	router.PUT("/users/:id", userHandler.update)
	router.DELETE("/users/:id", userHandler.delete)

	s.registerDocsRoutes(router)

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
