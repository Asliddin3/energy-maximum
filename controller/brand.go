package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BrandController struct {
	*Handler
}

func (h *Handler) NewBrandController(api *gin.RouterGroup) {
	brand := &BrandController{h}
	br := api.Group("brand", h.DeserializeAdmin())
	{
		br.POST("/", brand.CreateBrand)
		br.PUT("/:id", brand.UpdateBrand)
		br.GET("/all", brand.GetBrands)
		br.DELETE("/:id", brand.DeleteByID)
	}
	customerBr := api.Group("brand")
	{
		customerBr.GET("", brand.GetCustomerBrands)
		customerBr.GET("/letter", brand.GetBrandsByLetter)
		customerBr.GET("/:id", brand.GetByID)
	}
}

// @Summary		  Create new brand
// @Description	   this api is create new brand
// @Tags			Brand
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	formData		models.BrandRequest	true	"data body"
// @Param			image_file	formData	file		false	"file"
// @Success			201		{object}	models.Brand
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/brand/ [POST]
func (h *BrandController) CreateBrand(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.BrandRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create brand", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	var image string
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathBrand, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}

		image = name
	}
	fmt.Println("name ru", body.NameEn)

	brand := models.Brand{
		NameRu:         body.NameRu,
		NameUz:         body.NameUz,
		NameEn:         body.NameEn,
		Letter:         strings.ToUpper(body.Letter),
		DescriptionRu:  body.DescriptionRu,
		DescriptionUz:  body.DescriptionUz,
		Position:       body.Position,
		DescriptionEn:  body.DescriptionEn,
		SeoTitle:       body.SeoTitle,
		IsActive:       body.IsActive,
		Image:          image,
		SeoDescription: body.SeoDescription,
		CreatedID:      &admin.Id,
		CreatedAt:      timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Create(&brand).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, brand)
}

// @Summary		  Update brand
// @Description	   this api is Update brand
// @Tags			Brand
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "update id"
// @Param			data 	formData		models.BrandRequest	true	"data body"
// @Param			image_file	formData	file		false	"file"
// @Success			201		{object}	models.Brand
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/brand/{id} [PUT]
func (h *BrandController) UpdateBrand(c *gin.Context) {
	admin := h.GetAdmin(c)
	brandId := c.Param("id")
	var body models.BrandRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
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
	if body.SeoDescription != "" {
		columns["seo_description"] = body.SeoDescription
	}
	if body.SeoTitle != "" {
		columns["seo_title"] = body.SeoTitle
	}
	if body.IsActive != nil {
		columns["is_active"] = body.IsActive
	}
	if body.Position != nil {
		columns["position"] = body.Position
	}
	if body.Letter != "" {
		columns["letter"] = strings.ToUpper(body.Letter)
	}
	columns["updated_at"] = timeNow()
	columns["updated_id"] = admin.Id
	brand := &models.Brand{}
	err = h.db.Clauses(clause.Returning{}).Model(&brand).
		Where("id=?", brandId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save brand", err.Error())
		return
	}
	c.JSON(http.StatusOK, brand)
}

// @Summary		  Get brands
// @Description	   this api is to get brands
// @Tags			Brand
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Brand
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/brand/all  [GET]
func (h *BrandController) GetBrands(c *gin.Context) {
	var brands []models.Brand
	err := h.db.Model(&models.Brand{}).Find(&brands).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, brands)
}

// @Summary		  Get brands
// @Description	   this api is to get brands
// @Tags			Brand
// @Accept			json
// @Produce			json
// @Param           data    query   models.BrandFilter  true "brand filter"
// @Success			201		{object}	[]models.Brand
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/brand  [GET]
func (h *BrandController) GetCustomerBrands(c *gin.Context) {
	var body models.BrandFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var brands []models.Brand
	db := h.db.Model(&models.Brand{})
	if body.Name != "" {
		field := fmt.Sprintf("%%%s%%", body.Name)
		db = db.Where("name_ru LIKE ? OR name_uz LIKE ? OR name_en LIKE ?", field, field, field)
	}
	if body.Letter != "" {
		db = db.Where("letter=? ", body.Letter)
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	err = db.Order("position NULLS LAST").Limit(body.PageSize).Offset((body.Page-1)*body.PageSize).
		Find(&brands, "is_active=true").Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, brands)
}

// @Summary		  Get brands by letter
// @Description	   this api is to get brands by letter
// @Tags			Brand
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.BrandsByLetter
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/brand/letter  [GET]
func (h *BrandController) GetBrandsByLetter(c *gin.Context) {
	var body models.BrandFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var letters []string
	err = h.db.Model(&models.Brand{}).Order("letter").Group("letter").Select("UPPER(letter)").Find(&letters).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := make([]models.BrandsByLetter, len(letters))
	for i, letter := range letters {
		var brands []models.Brand
		err = h.db.Model(&models.Brand{}).Order("position NULLS LAST").Find(&brands, "is_active=true AND UPPER(letter)=?", letter).Error
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		res[i] = models.BrandsByLetter{
			Letter: letter,
			Brands: brands,
		}

	}

	c.JSON(http.StatusOK, res)
}

// @Summary		  Get brands
// @Description	   this api is to get brands
// @Tags			Brand
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "id"
// @Success			201		{object}	[]models.Brand
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/brand/{id} [GET]
func (h *BrandController) GetByID(c *gin.Context) {
	brandId := c.Param("id")
	var brands models.Brand
	err := h.db.First(&brands, "id=?", brandId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found brand")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, brands)
}

// @Summary		  DELETE brands
// @Description	   this api is to delete brands
// @Tags			Brand
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/brand/{id} [DELETE]
func (h *BrandController) DeleteByID(c *gin.Context) {
	// admin := h.GetAdmin(c)
	brandId := c.Param("id")
	if brandId == "" {
		newResponse(c, http.StatusBadRequest, "empty brand id")
		return
	}
	err := h.db.Delete(&models.Brand{}, "id=?", brandId).Error
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			newResponse(c, http.StatusBadRequest, "this brand has relation to product")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
