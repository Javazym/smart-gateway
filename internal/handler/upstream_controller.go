package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
)

type UpstreamController struct {
	upstreamService *service.UpstreamService
}

func NewUpstreamController(upstreamService *service.UpstreamService) *UpstreamController {
	return &UpstreamController{upstreamService: upstreamService}
}

// 分组

func (upstreamController *UpstreamController) ListGroup(c *gin.Context) {
	var dto model.GroupDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	list, total, err := upstreamController.upstreamService.ListGroup(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OKPage(list, total, dto.Page, dto.Size, ""))
}

func (upstreamController *UpstreamController) GetGroup(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	group, err := upstreamController.upstreamService.GetGroup(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(group))
}

func (upstreamController *UpstreamController) CreateGroup(c *gin.Context) {
	var group model.UpstreamGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	if err := upstreamController.upstreamService.CreateGroup(c.Request.Context(), &group); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (upstreamController *UpstreamController) UpdateGroup(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var group model.UpstreamGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	group.ID = id
	if err := upstreamController.upstreamService.UpdateGroup(c.Request.Context(), &group); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (upstreamController *UpstreamController) DeleteGroup(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := upstreamController.upstreamService.DeleteGroup(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

// 节点

func (upstreamController *UpstreamController) ListUpstream(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Query("group_id"), 10, 64)
	list, err := upstreamController.upstreamService.ListUpstream(c.Request.Context(), groupID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(list))
}

func (upstreamController *UpstreamController) GetUpstream(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	upstream, err := upstreamController.upstreamService.GetUpstream(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(upstream))
}

func (upstreamController *UpstreamController) CreateUpstream(c *gin.Context) {
	var upstream model.Upstream
	if err := c.ShouldBindJSON(&upstream); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	if err := upstreamController.upstreamService.CreateUpstream(c.Request.Context(), &upstream); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (upstreamController *UpstreamController) UpdateUpstream(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var upstream model.Upstream
	if err := c.ShouldBindJSON(&upstream); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	upstream.ID = id
	if err := upstreamController.upstreamService.UpdateUpstream(c.Request.Context(), &upstream); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (upstreamController *UpstreamController) DeleteUpstream(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := upstreamController.upstreamService.DeleteUpstream(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}
