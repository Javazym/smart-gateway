package handler

import (
	"net/http"
	"smart-gateway/internal/util"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
)

type AuthController struct {
	AuthService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (authController *AuthController) Login(c *gin.Context) {
	var req model.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	resp, err := authController.AuthService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.Error(http.StatusUnauthorized, err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.OK(resp))
}

func (authController *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, util.OK(nil))
}

func (authController *AuthController) AddUser(c *gin.Context) {
	var req model.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	if err := authController.AuthService.CreateUser(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}
