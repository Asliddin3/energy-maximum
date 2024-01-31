package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NewsController struct {
	*Handler
}

func (h *Handler) NewNewsController(api *gin.RouterGroup) {
	news := &NewsController{h}
	nw := api.Group("news", h.DeserializeAdmin())
	{
		nw.POST("", news.CreateNews)
		nw.PUT("/:id", news.UpdateNews)
		nw.DELETE("/:id", news.DeleteNew)
	}
	api.GET("/news/:id", news.GetByID)
	api.GET("/news", news.GetNews)
}

// @Summary		  Create news
// @Description	   this api is create news
// @Tags			News
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	formData		models.NewsRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.News
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/news [POST]
func (h *NewsController) CreateNews(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.NewsRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create news", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathNews, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}

		body.Image = name

	}
	news := models.News{
		NameRu:        body.NameRu,
		NameUz:        body.NameUz,
		NameEn:        body.NameEn,
		DescriptionRu: body.DescriptionRu,
		DescriptionUz: body.DescriptionUz,
		DescriptionEn: body.DescriptionEn,
		Image:         body.Image,
		CreatedID:     &admin.Id,
		CreatedAt:     timeNow(),
	}
	err = h.db.Clauses(clause.Returning{}).Create(&news).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, news)
}

// @Summary		  	Update news
// @Description	   	this api is Update news
// @Tags			News
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "news id"
// @Param			phone 	formData		models.NewsRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.News
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/news/{id} [PUT]
func (h *NewsController) UpdateNews(c *gin.Context) {
	admin := h.GetAdmin(c)
	newsId := c.Param("id")
	var body models.NewsRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create news", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	columns := map[string]interface{}{}
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathBrand, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}
		columns["image"] = name
	}
	if body.DescriptionEn != "" {
		columns["description_en"] = body.DescriptionEn
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
	if body.Position != nil {
		columns["position"] = body.Position
	}
	// if body.IsActive != nil {
	// 	columns["is_active"] = body.IsActive
	// }
	columns["updated_at"] = timeNow()
	columns["updated_id"] = admin.Id
	news := &models.News{}
	err = h.db.Clauses(clause.Returning{}).Model(&news).
		Where("id=?", newsId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save news", err.Error())
		return
	}
	c.JSON(http.StatusOK, news)
}

// @Summary		  Get news
// @Description	   this api is get news
// @Tags			News
// @Accept			json
// @Produce			json
// @Param           data 	query   models.NewsFilter  true "filter"
// @Success			201		{object}	[]models.News
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/news [GET]
func (h *NewsController) GetNews(c *gin.Context) {
	var body models.NewsFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var news []models.News
	db := h.db.Model(&models.News{}).Order("position NULLS LAST")
	if body.Name != "" {
		field := fmt.Sprintf("%%%s%%", body.Name)
		db = db.Where("name_en LIKE ? OR name_ru LIKE ? OR name_uz LIKE ?", field, field, field)
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	err = db.Limit(body.PageSize).Offset((body.Page - 1) * body.PageSize).Find(&news).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find news", err.Error())
		return
	}
	c.JSON(http.StatusOK, news)
}

// @Summary		  Update news
// @Description	   this api is Update news
// @Tags			News
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "news id"
// @Success			201		{object}	models.News
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/news/{id} [GET]
func (h *NewsController) GetByID(c *gin.Context) {
	inputId := c.Param("id")
	var news models.News
	err := h.db.First(&news, "id=?", inputId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found news")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, news)
}

// @Summary		  DELETE news
// @Description	   this api is delete news
// @Tags			News
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "news id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/news/{id} [DELETE]
func (h *NewsController) DeleteNew(c *gin.Context) {
	inputId := c.Param("id")
	var news models.News
	err := h.db.Delete(&news, "id=?", inputId).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
