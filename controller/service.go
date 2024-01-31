package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServiceController struct {
	*Handler
}

func (h *Handler) NewServiceController(api *gin.RouterGroup) {
	service := &ServiceController{h}
	Ban := api.Group("service", h.DeserializeAdmin())
	{
		Ban.POST("", service.CreateService)
		Ban.PUT("/:id", service.UpdateService)
		Ban.GET("/:id", service.GetByID)
		// Ban.GET("/all", service.Getservice)
		Ban.DELETE("/:id", service.DeleteByID)
	}
	customBan := api.Group("service")
	{
		customBan.GET("", service.GetCustomerService)
	}
}

// @Summary		  Create new service
// @Description	   this api is create new service
// @Tags			Service
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	formData		models.ServiceRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Service
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/service [POST]
func (h *ServiceController) CreateService(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.ServiceRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathService, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}

		body.Image = name

	}
	category := models.Service{
		NameRu:        body.NameRu,
		NameUz:        body.NameUz,
		NameEn:        body.NameEn,
		Position:      body.Position,
		DescriptionRu: body.DescriptionRu,
		DescriptionUz: body.DescriptionUz,
		DescriptionEn: body.DescriptionEn,
		Image:         body.Image,
		CreatedID:     &admin.Id,
		CreatedAt:     timeNow(),
	}
	err = h.db.Clauses(clause.Returning{}).Create(&category).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, category)
}

// @Summary		  Update service
// @Description	   this api is Update service
// @Tags			Service
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "category id"
// @Param			data 	formData		models.ServiceRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Service
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/service/{id} [PUT]
func (h *ServiceController) UpdateService(c *gin.Context) {
	admin := h.GetAdmin(c)
	categoryId := c.Param("id")
	id, _ := strconv.ParseInt(categoryId, 10, 64)
	var body models.ServiceRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create service", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	columns := map[string]interface{}{}

	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathService, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}
		columns["image"] = name
	}
	category := models.Service{
		ID: int(id),
	}
	if body.DescriptionEn != "" {
		columns["description_en"] = body.DescriptionEn
	}
	if body.Position != nil {
		columns["position"] = body.Position
	}
	if body.DescriptionRu != "" {
		columns["description_ru"] = body.DescriptionRu
	}
	if body.DescriptionUz != "" {
		columns["description_uz"] = body.DescriptionUz
	}
	if body.NameEn != "" {
		columns["name_en"] = body.NameEn
	}
	if body.NameRu != "" {
		columns["name_ru"] = body.NameRu
	}
	if body.NameUz != "" {
		columns["name_uz"] = body.NameUz
	}
	columns["updated_at"] = timeNow()
	columns["updated_id"] = admin.Id
	err = h.db.Debug().Clauses(clause.Returning{}).Model(&category).
		Where("id=?", categoryId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save category", err.Error())
		return
	}
	c.JSON(http.StatusOK, category)
}

// @Summary		  Get service
// @Description	   this api is get service
// @Tags			Service
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Service
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/service [GET]
func (h *ServiceController) GetCustomerService(c *gin.Context) {
	var service []models.Service
	err := h.db.Model(&models.Service{}).Order("position NULLS LAST").Find(&service).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find service", err.Error())
		return
	}
	c.JSON(http.StatusOK, service)
}

// @Summary		  Update service
// @Description	   this api is Update service
// @Tags			Service
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "service id"
// @Success			201		{object}	models.Service
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/service/{id} [GET]
func (h *ServiceController) GetByID(c *gin.Context) {
	inputId := c.Param("id")
	var service models.Service
	err := h.db.First(&service, "id=?", inputId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found service")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, service)
}

// @Summary		  DELETE service
// @Description	   this api is to delete service
// @Tags			Service
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/service/{id} [DELETE]
func (h *ServiceController) DeleteByID(c *gin.Context) {
	// admin := h.GetAdmin(c)
	serviceId := c.Param("id")
	if serviceId == "" {
		newResponse(c, http.StatusBadRequest, "empty service id")
		return
	}
	err := h.db.Delete(&models.Service{}, "id=?", serviceId).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
