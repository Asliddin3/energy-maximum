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

type AboutController struct {
	*Handler
}

func (h *Handler) NewAboutController(api *gin.RouterGroup) {
	brand := &AboutController{h}
	br := api.Group("about", h.DeserializeAdmin())
	{
		br.POST("/", brand.CreateAbout)
		br.PUT("/:id", brand.UpdateAbout)
		// br.GET("", brand.GetAbout)
		br.GET("/:id", brand.GetByID)
		br.DELETE("/:id", brand.DeleteAbout)
	}
	api.GET("/about", brand.GetAbout)
}

// @Summary		  Create new about
// @Description	   this api is create new about
// @Tags			About
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data    formData		models.AboutRequest	false	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.About
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/about/ [POST]
func (h *AboutController) CreateAbout(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.AboutRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create about", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathAbout, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}
		body.Image = name
	}
	about := models.About{
		NameRu:        body.NameRu,
		NameUz:        body.NameUz,
		NameEn:        body.NameEn,
		Image:         body.Image,
		DescriptionUz: body.DescriptionUz,
		DescriptionRu: body.DescriptionRu,
		DescriptionEn: body.DescriptionEn,
		Position:      body.Position,
		Type:          body.Type,
		CreatedID:     &admin.Id,
		CreatedAt:     timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Create(&about).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, about)
}

// @Summary		  Update about
// @Description	   this api is Update about
// @Tags			About
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "update id"
// @Param			data 	formData		models.AboutRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.About
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/about/{id} [PUT]
func (h *AboutController) UpdateAbout(c *gin.Context) {
	admin := h.GetAdmin(c)
	inputId := c.Param("id")
	var body models.AboutRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create about", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	columns := map[string]interface{}{}

	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathAbout, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}
		columns["image"] = name
	}
	id, _ := strconv.ParseInt(inputId, 10, 64)

	about := models.About{
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
	if body.Type != "" {
		columns["type"] = body.Type
	}
	columns["updated_at"] = timeNow()
	columns["updated_id"] = admin.Id
	err = h.db.Debug().Clauses(clause.Returning{}).Model(&about).
		Where("id=?", inputId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save about", err.Error())
		return
	}
	c.JSON(http.StatusOK, about)
}

// @Summary		    Get about
// @Description	    this api is to get about
// @Tags			About
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Country
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/about  [GET]
func (h *AboutController) GetAbout(c *gin.Context) {
	var countries []models.About
	err := h.db.Model(&models.About{}).Order("position NULLS LAST").Find(&countries).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, countries)
}

// @Summary		  Get about
// @Description	   this api is to get about
// @Tags			About
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "id"
// @Success			201		{object}	models.About
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/about/{id} [GET]
func (h *AboutController) GetByID(c *gin.Context) {
	brandId := c.Param("id")
	var about models.About
	err := h.db.First(&about, "id=?", brandId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found about")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, about)
}

// @Summary		  Delete about
// @Description	   this api is to delete about
// @Tags			About
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/about/{id} [DELETE]
func (h *AboutController) DeleteAbout(c *gin.Context) {
	id := c.Param("id")
	err := h.db.Delete(&models.About{}, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
