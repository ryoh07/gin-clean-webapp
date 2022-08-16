package handler

import (
	"context"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryoh07/gin-clean-webapp/common/dto"
	"github.com/ryoh07/gin-clean-webapp/usecase"
)

type IUserHandler interface {
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	CreateUser(c *gin.Context)
	DeleteUser(cc *gin.Context)
}

type UserHandler struct {
	usecase.UserInputPort
}

func NewUserHandler(srv usecase.UserInputPort) IUserHandler {
	return &UserHandler{srv}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userId := c.Param("id")
	user, _ := h.UserInputPort.GetUser(userId)
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := context.Background()
	user := dto.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.UserInputPort.UpdateUser(ctx, &user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := context.Background()
	user := dto.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.UserInputPort.CreateUser(ctx, &user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx := context.Background()
	userId := c.Param("userid")

	h.UserInputPort.DeleteUser(ctx, userId)
}
