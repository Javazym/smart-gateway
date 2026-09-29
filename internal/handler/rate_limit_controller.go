package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
)

type RateLimitController struct {
	rateLimitService *service.RateLimitService
}

func NewRateLimitController(rateLimitService *service.RateLimitService) *RateLimitController {
	return &RateLimitController{rateLimitService: rateLimitService}
}

func (rateLimitController *RateLimitController) List(c *gin.Context) {
	var dto model.RateLimitPolicyDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	list, total, err := rateLimitController.rateLimitService.List(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OKPage(list, total, dto.Page, dto.Size, ""))
}

func (rateLimitController *RateLimitController) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	policy, err := rateLimitController.rateLimitService.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(policy))
}

func (rateLimitController *RateLimitController) Create(c *gin.Context) {
	var policy model.RateLimitPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	if err := rateLimitController.rateLimitService.Create(c.Request.Context(), &policy); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (rateLimitController *RateLimitController) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var policy model.RateLimitPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	policy.ID = id
	if err := rateLimitController.rateLimitService.Update(c.Request.Context(), &policy); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (rateLimitController *RateLimitController) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := rateLimitController.rateLimitService.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}
