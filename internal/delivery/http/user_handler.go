package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/islamyakin/otel-propagation-monorepo/internal/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// This handler is only accessible by admins (enforced by middleware)
	users, err := h.userUseCase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
