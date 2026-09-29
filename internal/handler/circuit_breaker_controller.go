package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
)

type CircuitBreakerController struct {
	circuitBreakerService *service.CircuitBreakerService
}

func NewCircuitBreakerController(circuitBreakerService *service.CircuitBreakerService) *CircuitBreakerController {
	return &CircuitBreakerController{circuitBreakerService: circuitBreakerService}
}

func (circuitBreakerController *CircuitBreakerController) List(c *gin.Context) {
	var dto model.CircuitBreakerDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	list, total, err := circuitBreakerController.circuitBreakerService.List(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OKPage(list, total, dto.Page, dto.Size, ""))
}

func (circuitBreakerController *CircuitBreakerController) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	cb, err := circuitBreakerController.circuitBreakerService.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(cb))
}

func (circuitBreakerController *CircuitBreakerController) Create(c *gin.Context) {
	var cb model.CircuitBreaker
	if err := c.ShouldBindJSON(&cb); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	if err := circuitBreakerController.circuitBreakerService.Create(c.Request.Context(), &cb); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (circuitBreakerController *CircuitBreakerController) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var cb model.CircuitBreaker
	if err := c.ShouldBindJSON(&cb); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	cb.ID = id
	if err := circuitBreakerController.circuitBreakerService.Update(c.Request.Context(), &cb); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (circuitBreakerController *CircuitBreakerController) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := circuitBreakerController.circuitBreakerService.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}
