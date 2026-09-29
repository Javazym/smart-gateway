package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
)

type RequestLogController struct {
	requestLogService *service.RequestLogService
}

func NewRequestLogController(requestLogService *service.RequestLogService) *RequestLogController {
	return &RequestLogController{requestLogService: requestLogService}
}

func (requestLogController *RequestLogController) List(c *gin.Context) {
	var dto model.RequestLogDto
	if err := c.ShouldBind(&dto); err != nil {
		return
	}
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	list, total, err := requestLogController.requestLogService.List(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OKPage(list, total, dto.Page, dto.Size, ""))
}
