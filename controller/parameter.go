package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ParameterController struct {
	*Handler
}

func (h *Handler) NewParameterController(api *gin.RouterGroup) {
	brand := &ParameterController{h}
	br := api.Group("parameter", h.DeserializeAdmin())
	{
		br.POST("/", brand.CreateParameter)
		br.PUT("/:id", brand.UpdateParameter)
		br.GET("", brand.GetProductParameter)
		br.GET("/:id", brand.GetByID)
		br.DELETE("/:id", brand.DeleteParameter)
	}
}

// @Summary		  Create new parameter
// @Description	   this api is create new parameter
// @Tags			Parameter
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data    body		models.ParametersRequest	false	"data body"
// @Success			201		{object}	models.Parameters
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/parameter/ [POST]
func (h *ParameterController) CreateParameter(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.ParametersRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	parameter := models.Parameters{
		NameRu:    body.NameRu,
		NameUz:    body.NameUz,
		NameEn:    body.NameEn,
		Position:  body.Position,
		CreatedID: &admin.Id,
		CreatedAt: timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Create(&parameter).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, parameter)
}

// @Summary		  Update parameter
// @Description	   this api is Update parameter
// @Tags			Parameter
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "update id"
// @Param			data 	body		models.ParametersRequest	true	"data body"
// @Success			201		{object}	models.Parameters
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/parameter/{id} [PUT]
func (h *ParameterController) UpdateParameter(c *gin.Context) {
	admin := h.GetAdmin(c)
	id := c.Param("id")
	var body models.ParametersRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	parameter := models.Parameters{}
	err = h.db.First(&parameter, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to create parameter", err.Error())
		return
	}
	if body.NameEn != "" {
		parameter.NameEn = body.NameEn
	}
	if body.NameRu != "" {
		parameter.NameRu = body.NameRu
	}
	if body.NameUz != "" {
		parameter.NameUz = body.NameUz
	}
	if body.Position != nil {
		parameter.Position = body.Position
	}
	parameter.CreatedAt = timeNow()
	parameter.CreatedID = &admin.Id
	err = h.db.Save(&parameter).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save brand", err.Error())
		return
	}
	c.JSON(http.StatusOK, parameter)
}

// @Summary		  get product parameter
// @Description	   this api is for get product parameter
// @Tags			Parameter
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           data    query   models.ParameterFilter  true "filter parameters"
// @Success			201		{object}	[]models.Parameters
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/parameter [GET]
func (h *ParameterController) GetProductParameter(c *gin.Context) {
	var body models.ParameterFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	db := h.db.Model(&models.Parameters{}).Order("position NULLS LAST")
	if body.Name != "" {
		filed := fmt.Sprintf("%%%s%%", body.Name)
		filed = strings.ToLower(filed)
		db = db.Where("name_ru LIKE ? OR name_en LIKE ? OR name_uz LIKE ?", filed, filed, filed)
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	if !body.WithDeleted {
		db = db.Where("is_deleted=false")
	}
	var count int64
	err = db.Count(&count).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get count ")
		h.log.Error("failed to get count", err.Error())
		return
	}
	var params []models.Parameters
	err = db.Limit(body.PageSize).Offset((body.Page - 1) * body.PageSize).Find(&params).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.ParameterResponse{
		Parameters: params,
		Page:       body.Page,
		PageSize:   body.PageSize,
		Count:      int(count),
	})
}

// @Summary		  Get Parameter
// @Description	   this api is to get Parameter
// @Tags			Parameter
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	models.Parameters
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/parameter/{id} [GET]
func (h *ParameterController) GetByID(c *gin.Context) {
	paramId := c.Param("id")
	var parameter models.Parameters
	err := h.db.First(&parameter, "id=?", paramId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found parameter")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, parameter)
}

// @Summary		  Delete parameter
// @Description	   this api is to delete parameter
// @Tags			Parameter
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/parameter/{id} [DELETE]
func (h *ParameterController) DeleteParameter(c *gin.Context) {
	id := c.Param("id")
	err := h.db.Model(&models.Parameters{}).Where("id=?", id).UpdateColumn("is_deleted", true).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
