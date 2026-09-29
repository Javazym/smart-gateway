package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
)

type RouteController struct {
	routeService *service.RouteService
}

func NewRouteController(routeService *service.RouteService) *RouteController {
	return &RouteController{routeService: routeService}
}

func (routeController *RouteController) List(c *gin.Context) {
	var dto model.RouteDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	list, total, err := routeController.routeService.List(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OKPage(list, total, dto.Page, dto.Size, ""))
}

func (routeController *RouteController) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	route, err := routeController.routeService.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(route))
}

func (routeController *RouteController) Create(c *gin.Context) {
	var route model.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	if err := routeController.routeService.Create(c.Request.Context(), &route); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (routeController *RouteController) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var route model.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	route.ID = id
	if err := routeController.routeService.Update(c.Request.Context(), &route); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (routeController *RouteController) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := routeController.routeService.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}
