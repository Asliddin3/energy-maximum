package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PagesController struct {
	*Handler
}

func (h *Handler) NewPagesController(api *gin.RouterGroup) {
	category := &PagesController{h}
	cate := api.Group("pages", h.DeserializeAdmin())
	{
		cate.POST("", category.CreateCategory)
		cate.PUT("/:id", category.UpdateCategory)
		cate.GET("/:id", category.GetByID)
		cate.GET("/all", category.GetCategory)
		cate.DELETE("/:id", category.DeleteByID)

	}
	customCate := api.Group("pages")
	{
		customCate.GET("", category.GetCustomerCategory)
	}
}

// @Summary		  Create new category
// @Description	   this api is create new category
// @Tags			Category
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	formData		models.PagesRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Pages
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/pages [POST]
func (h *PagesController) CreateCategory(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.PagesRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathCategory, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}

		body.Image = name

	}
	category := models.Pages{
		NameRu:           body.NameRu,
		NameUz:           body.NameUz,
		NameEn:           body.NameEn,
		Url:              body.Url,
		DescriptionRu:    body.DescriptionRu,
		DescriptionUz:    body.DescriptionUz,
		DescriptionEn:    body.DescriptionEn,
		Position:         body.Position,
		SeoTitleRu:       body.SeoTitleRu,
		SeoTitleEn:       body.SeoTitleEn,
		SeoTitleUz:       body.SeoTitleUz,
		IsActive:         body.IsActive,
		Image:            body.Image,
		SeoDescriptionRu: body.SeoDescriptionRu,
		SeoDescriptionEn: body.SeoDescriptionEn,
		SeoDescriptionUz: body.SeoDescriptionUz,
		CreatedID:        &admin.Id,
		CreatedAt:        timeNow(),
	}
	if body.Position != nil {
		category.Position = body.Position
	}
	err = h.db.Clauses(clause.Returning{}).Create(&category).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, category)
}

// @Summary		  Update category
// @Description	   this api is Update category
// @Tags			Category
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "category id"
// @Param			data 	formData		models.PagesRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Pages
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/pages/{id} [PUT]
func (h *PagesController) UpdateCategory(c *gin.Context) {
	admin := h.GetAdmin(c)
	categoryId := c.Param("id")
	id, _ := strconv.ParseInt(categoryId, 10, 64)
	var body models.PagesRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	columns := map[string]interface{}{}

	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathCategory, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}
		columns["image"] = name
	}
	category := models.Pages{
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
	if body.Url != "" {
		columns["url"] = body.Url
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
	if body.SeoDescriptionRu != "" {
		columns["seo_description_ru"] = body.SeoDescriptionRu
	}
	if body.SeoDescriptionEn != "" {
		columns["seo_description_en"] = body.SeoDescriptionEn
	}
	if body.SeoDescriptionUz != "" {
		columns["seo_description_uz"] = body.SeoDescriptionUz
	}

	if body.SeoTitleRu != "" {
		columns["seo_title_ru"] = body.SeoTitleRu
	}
	if body.SeoTitleUz != "" {
		columns["seo_title_uz"] = body.SeoTitleUz
	}
	if body.SeoTitleEn != "" {
		columns["seo_title_en"] = body.SeoTitleEn
	}

	if body.IsActive != nil {
		columns["is_active"] = body.IsActive
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

// @Summary		  Get category
// @Description	   this api is get category
// @Tags			Category
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Pages
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/pages/all [GET]
func (h *PagesController) GetCategory(c *gin.Context) {
	var category []models.Pages
	err := h.db.Find(&category).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find category", err.Error())
		return
	}
	c.JSON(http.StatusOK, category)
}

// @Summary		  Get category
// @Description	   this api is get category
// @Tags			Category
// @Accept			json
// @Produce			json
// @Param           data  query    models.PagesFilter true "category filter"
// @Success			201		{object}	[]models.Pages
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/pages [GET]
func (h *PagesController) GetCustomerCategory(c *gin.Context) {
	var body models.PagesFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	var category []models.Pages
	db := h.db.Debug().Model(&models.Pages{}).Order("position NULLS LAST")
	if body.Name != "" {
		field := fmt.Sprintf("%%%s%%", body.Name)
		db = db.Where("name_ru LIKE ? OR name_uz  LIKE ? OR name_en LIKE ?", field)
	}

	fmt.Println("size", body.PageSize)
	err = db.Limit(body.PageSize).Offset((body.Page-1)*body.PageSize).Find(&category, "is_active=true").Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get category")
		h.log.Error("failed to find category", err.Error())
		return
	}
	c.JSON(http.StatusOK, category)
}

// @Summary		  Update category
// @Description	   this api is Update category
// @Tags			Category
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "category id"
// @Success			201		{object}	models.Pages
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/pages/{id} [GET]
func (h *PagesController) GetByID(c *gin.Context) {
	inputId := c.Param("id")
	var category models.Pages
	err := h.db.Model(&models.Pages{}).Preload("Created").
		Preload("Updated").Preload("Deleted").First(&category, "id=?", inputId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found category")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, category)
}

// @Summary		  DELETE category
// @Description	   this api is to delete category
// @Tags			Category
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/pages/{id} [DELETE]
func (h *PagesController) DeleteByID(c *gin.Context) {
	// admin := h.GetAdmin(c)
	categoryId := c.Param("id")
	if categoryId == "" {
		newResponse(c, http.StatusBadRequest, "empty category id")
		return
	}
	err := h.db.Delete(&models.Pages{}, "id=?", categoryId).Error
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			newResponse(c, http.StatusBadRequest, "this category has relation to product")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
