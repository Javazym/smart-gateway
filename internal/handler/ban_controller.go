package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
)

type BanController struct {
	banService *service.BanService
}

func NewBanController(banService *service.BanService) *BanController {
	return &BanController{banService: banService}
}

func (banController *BanController) List(c *gin.Context) {
	var dto model.BanDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	list, total, err := banController.banService.List(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OKPage(list, total, dto.Page, dto.Size, ""))
}

func (banController *BanController) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ban, err := banController.banService.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(ban))
}

func (banController *BanController) Create(c *gin.Context) {
	var ban model.Ban
	if err := c.ShouldBindJSON(&ban); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	if err := banController.banService.Create(c.Request.Context(), &ban); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (banController *BanController) Unban(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := banController.banService.Unban(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (banController *BanController) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := banController.banService.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}
