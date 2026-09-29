package handler

import (
	"net/http"
	"smart-gateway/internal/util"

	"github.com/gin-gonic/gin"
)

// Health 返回网关健康状态。
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, util.OK(gin.H{"status": "up"}))
}
