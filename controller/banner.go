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

type BannerController struct {
	*Handler
}

func (h *Handler) NewBannerController(api *gin.RouterGroup) {
	banner := &BannerController{h}
	Ban := api.Group("banner", h.DeserializeAdmin())
	{
		Ban.POST("", banner.CreateBanner)
		Ban.PUT("/:id", banner.UpdateBanner)
		Ban.GET("/:id", banner.GetByID)
		// Ban.GET("/all", banner.GetBanner)
		Ban.DELETE("/:id", banner.DeleteByID)
	}
	customBan := api.Group("banner")
	{
		customBan.GET("", banner.GetCustomerBanner)
	}
}

// @Summary		  Create new banner
// @Description	   this api is create new banner
// @Tags			Banner
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	formData		models.BannerRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Banner
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/banner [POST]
func (h *BannerController) CreateBanner(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.BannerRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathBanner, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}

		body.Image = name

	}
	category := models.Banner{
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

// @Summary		  Update banner
// @Description	   this api is Update banner
// @Tags			Banner
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "category id"
// @Param			data 	formData		models.BannerRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Banner
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/banner/{id} [PUT]
func (h *BannerController) UpdateBanner(c *gin.Context) {
	admin := h.GetAdmin(c)
	categoryId := c.Param("id")
	id, _ := strconv.ParseInt(categoryId, 10, 64)
	var body models.BannerRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	columns := map[string]interface{}{}

	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathBanner, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}
		columns["image"] = name
	}
	category := models.Banner{
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

// @Summary		  Get banner
// @Description	   this api is get banner
// @Tags			Banner
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Banner
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/banner [GET]
func (h *BannerController) GetCustomerBanner(c *gin.Context) {
	var banner []models.Banner
	err := h.db.Model(&models.Banner{}).Order("position NULLS LAST").Find(&banner).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find banner", err.Error())
		return
	}
	c.JSON(http.StatusOK, banner)
}

// @Summary		  Update banner
// @Description	   this api is Update banner
// @Tags			Banner
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "banner id"
// @Success			201		{object}	models.Banner
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/banner/{id} [GET]
func (h *BannerController) GetByID(c *gin.Context) {
	inputId := c.Param("id")
	var banner models.Banner
	err := h.db.First(&banner, "id=?", inputId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found banner")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, banner)
}

// @Summary		  DELETE banner
// @Description	   this api is to delete banner
// @Tags			Banner
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/banner/{id} [DELETE]
func (h *BannerController) DeleteByID(c *gin.Context) {
	// admin := h.GetAdmin(c)
	bannerId := c.Param("id")
	if bannerId == "" {
		newResponse(c, http.StatusBadRequest, "empty banner id")
		return
	}
	err := h.db.Delete(&models.Banner{}, "id=?", bannerId).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
