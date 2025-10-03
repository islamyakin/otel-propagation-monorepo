package route

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/islamyakin/otel-propagation-monorepo/internal/delivery/http"
	"github.com/islamyakin/otel-propagation-monorepo/internal/delivery/http/middleware"
	"github.com/islamyakin/otel-propagation-monorepo/internal/usecase"
)

type Handler struct {
	AuthHandler *httpHandler.AuthHandler
	TodoHandler *httpHandler.TodoHandler
	UserHandler *httpHandler.UserHandler
}

func NewHandler(
	authUseCase usecase.AuthUseCase,
	todoUseCase usecase.TodoUseCase,
	userUseCase usecase.UserUseCase,
) *Handler {
	return &Handler{
		AuthHandler: httpHandler.NewAuthHandler(authUseCase),
		TodoHandler: httpHandler.NewTodoHandler(todoUseCase),
		UserHandler: httpHandler.NewUserHandler(userUseCase),
	}
}

func SetupRoutes(router *gin.Engine, handler *Handler, authUseCase usecase.AuthUseCase) {
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		v1.POST("/register", handler.AuthHandler.Register)
		v1.POST("/login", handler.AuthHandler.Login)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(authUseCase))
		{
			// User profile
			protected.GET("/profile", handler.AuthHandler.GetProfile)

			// Todo routes for users
			protected.POST("/todos", handler.TodoHandler.Create)
			protected.GET("/todos", handler.TodoHandler.GetUserTodos)
			protected.GET("/todos/:id", handler.TodoHandler.GetByID)
			protected.PUT("/todos/:id", handler.TodoHandler.Update)
			protected.DELETE("/todos/:id", handler.TodoHandler.Delete)
			protected.PATCH("/todos/:id/status", handler.TodoHandler.UpdateStatus)

			// Admin only routes
			admin := protected.Group("")
			admin.Use(middleware.AdminOnly())
			{
				// Admin can see all users
				admin.GET("/users", handler.UserHandler.GetAllUsers)

				// Admin can see all todos
				admin.GET("/admin/todos", handler.TodoHandler.GetAllTodos)
			}
		}
	}
}
