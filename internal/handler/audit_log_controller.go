package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
)

type AuditLogController struct {
	auditLogService *service.AuditLogService
}

func NewAuditLogController(auditLogService *service.AuditLogService) *AuditLogController {
	return &AuditLogController{auditLogService: auditLogService}
}

func (auditLogController *AuditLogController) List(c *gin.Context) {
	var dto model.AuditLogDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	list, total, err := auditLogController.auditLogService.List(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OKPage(list, total, dto.Page, dto.Size, ""))
}
