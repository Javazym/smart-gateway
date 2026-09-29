package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"smart-gateway/internal/util"
)

// parseID 解析路径中的 :id 参数为 uint64。
func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, "无效的ID"))
		return 0, false
	}
	return id, true
}

// respondError 统一处理 service 层错误：记录不存在返回 404，其余返回 500。
func respondError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, util.Error(http.StatusNotFound, "记录不存在"))
		return
	}
	log.Printf("internal error: %v", err)
	c.JSON(http.StatusInternalServerError, util.Error(http.StatusInternalServerError, "服务器内部错误"))
}
