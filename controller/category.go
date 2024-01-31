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

type CategoryController struct {
	*Handler
}

func (h *Handler) NewCategoryController(api *gin.RouterGroup) {
	category := &CategoryController{h}
	cate := api.Group("category", h.DeserializeAdmin())
	{
		cate.POST("", category.CreateCategory)
		cate.PUT("/:id", category.UpdateCategory)
		cate.GET("/:id", category.GetByID)
		cate.GET("/all", category.GetCategory)
		cate.DELETE("/:id", category.DeleteByID)
		cate.GET("/addition/:id", category.GetProductAddition)
		cate.POST("/addition", category.AddProductAddition)
		cate.DELETE("/addition", category.DeleteProductAddition)

	}
	customCate := api.Group("category")
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
// @Param			data 	formData		models.CategoryRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Category
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/category [POST]
func (h *CategoryController) CreateCategory(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.CategoryRequest
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
	category := models.Category{
		NameRu:         body.NameRu,
		NameUz:         body.NameUz,
		NameEn:         body.NameEn,
		DescriptionRu:  body.DescriptionRu,
		DescriptionUz:  body.DescriptionUz,
		DescriptionEn:  body.DescriptionEn,
		Position:       body.Position,
		SeoTitle:       body.SeoTitle,
		IsActive:       body.IsActive,
		Image:          body.Image,
		SeoDescription: body.SeoDescription,
		CreatedID:      &admin.Id,
		CreatedAt:      timeNow(),
	}
	if body.Position != nil {
		category.Position = body.Position
	}
	if body.ParentID != 0 {
		category.CategoryID = &body.ParentID
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
// @Param			data 	formData		models.CategoryRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Category
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/category/{id} [PUT]
func (h *CategoryController) UpdateCategory(c *gin.Context) {
	admin := h.GetAdmin(c)
	categoryId := c.Param("id")
	id, _ := strconv.ParseInt(categoryId, 10, 64)
	var body models.CategoryRequest
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
	category := models.Category{
		ID: int(id),
	}
	if body.DescriptionEn != "" {
		columns["description_en"] = body.DescriptionEn
	}
	if body.Position != nil {
		columns["position"] = body.Position
	}
	if body.ParentID != 0 {
		columns["category_id"] = body.ParentID
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
	if body.SeoDescription != "" {
		columns["seo_description"] = body.SeoDescription
	}
	if body.SeoTitle != "" {
		columns["seo_title"] = body.SeoTitle
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
// @Success			201		{object}	[]models.Category
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/category/all [GET]
func (h *CategoryController) GetCategory(c *gin.Context) {
	var category []models.Category
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
// @Param           data  query    models.CategoryFilter true "category filter"
// @Success			201		{object}	[]models.Category
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/category [GET]
func (h *CategoryController) GetCustomerCategory(c *gin.Context) {
	var body models.CategoryFilter
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
	var category []models.Category
	db := h.db.Debug().Model(&models.Category{}).Order("position NULLS LAST")
	if body.Name != "" {
		field := fmt.Sprintf("%%%s%%", body.Name)
		db = db.Where("name_ru LIKE ? OR name_uz  LIKE ? OR name_en LIKE ?", field)
	}
	if body.ParentID != nil {
		db = db.Where("parent_id=?", body.ParentID)
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
// @Success			201		{object}	models.Category
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/category/{id} [GET]
func (h *CategoryController) GetByID(c *gin.Context) {
	inputId := c.Param("id")
	var category models.Category
	err := h.db.First(&category, "id=?", inputId).Error
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
// @Router			/api/category/{id} [DELETE]
func (h *CategoryController) DeleteByID(c *gin.Context) {
	// admin := h.GetAdmin(c)
	categoryId := c.Param("id")
	if categoryId == "" {
		newResponse(c, http.StatusBadRequest, "empty category id")
		return
	}
	err := h.db.Delete(&models.Category{}, "id=?", categoryId).Error
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

// @Summary		  Add category addition
// @Description	   this api is to add category addition
// @Tags			Category
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           body    body    models.ProductAdditionRequest true "add products addition"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/category/addition [POST]
func (h *CategoryController) AddProductAddition(c *gin.Context) {
	var body models.ProductAdditionRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.db.Create(&models.ProductAdditions{
		ProductCategoryID:  body.ProductCategoryID,
		AdditionCategoryID: body.AdditionCategoryID,
	}).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			newResponse(c, http.StatusBadRequest, "this category has relation with addition")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		 Get category additions
// @Description	   this api is to get category addition
// @Tags			Category
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     	path   int  true  "category id"
// @Success			201		{object}	[]models.Products
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/category/addition/{id} [GET]
func (h *CategoryController) GetProductAddition(c *gin.Context) {
	id := c.Param("id")
	var products []models.Products
	err := h.db.Table("products AS p").Select("p.*").Joins("INNER JOIN product_additions AS pa ON pa.product_category_id=?", id).
		Where("p.parent_id=pa.addition_category_id AND p.is_active=true AND p.deleted_at IS NUll").Find(&products).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, products)
}

// @Summary		  DELETE category addition
// @Description	   this api is to delete category addition
// @Tags			Category
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           body    query    models.ProductAdditionRequest true "delete products addition"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/category/addition [DELETE]
func (h *CategoryController) DeleteProductAddition(c *gin.Context) {
	var body models.ProductAdditionRequest
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.db.Delete(&models.ProductAdditions{}, "product_category_id=? AND addition_category_id=?",
		body.ProductCategoryID, body.AdditionCategoryID).Error
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			newResponse(c, http.StatusBadRequest, "this category has relation with addition")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
