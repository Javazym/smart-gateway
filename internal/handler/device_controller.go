package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smart-gateway/internal/model"
	"smart-gateway/internal/service"
	"smart-gateway/internal/util"
)

type DeviceController struct {
	deviceService *service.DeviceService
}

func NewDeviceController(deviceService *service.DeviceService) *DeviceController {
	return &DeviceController{deviceService: deviceService}
}

func (deviceController *DeviceController) List(c *gin.Context) {
	var dto model.DeviceDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	list, total, err := deviceController.deviceService.List(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OKPage(list, total, dto.Page, dto.Size, ""))
}

func (deviceController *DeviceController) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	device, err := deviceController.deviceService.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(device))
}

func (deviceController *DeviceController) UpdateStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.Error(http.StatusBadRequest, err.Error()))
		return
	}
	if err := deviceController.deviceService.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}

func (deviceController *DeviceController) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := deviceController.deviceService.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, util.OK(nil))
}
