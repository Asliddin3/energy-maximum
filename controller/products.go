package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductController struct {
	*Handler
}

func (h *Handler) NewProductController(api *gin.RouterGroup) {
	product := &ProductController{h}
	prod := api.Group("product", h.DeserializeAdmin())
	{
		prod.POST("", product.CreateProduct)
		prod.PUT("/:id", product.UpdateProduct)
		prod.DELETE("/:id", product.DeleteProduct)
		prod.GET("/all", product.GetAllProducts)
		// prod.DELETE("/:id", product.)
		// prod.GET("", product.GetProducts)
		prod.POST("/media/:id", product.AddProductMedia)
		prod.DELETE("/media/:id", product.DeleteProductMedia)
		prod.POST("/parameter/:id", product.CreateProductParameter)

		// prod.POST("/recommend/", product.AddProductRecommend)
		// prod.DELETE("/recommend/:id", product.DeleteProductRecommend)
	}
	api.GET("/product", product.GetProducts)
	api.POST("/product/list", product.GetProductsByIds)
	api.GET("/product/:id", product.GetByID)

	// customProd := api.Group("product", h.DeserializeCustomer())
	{
		// customProd.POST("/favorite/:id", product.AddFavoriteProduct)
		// customProd.DELETE("/favorite/:id", product.DeleteFavoriteProduct)
		// customProd.GET("/favorite", product.GetFavoriteProduct)
	}
}

// @Summary		  Create new product
// @Description	   this api is create new product
// @Tags			Product
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	formData		models.ProductRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Products
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product [POST]
func (h *ProductController) CreateProduct(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.ProductRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create product", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathProducts, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return

		}

		body.Image = &name
	}
	product := models.Products{
		NameRu:         body.NameRu,
		NameUz:         body.NameUz,
		NameEn:         body.NameEn,
		Price:          body.Price,
		IsTop:          body.IsTop,
		IsNew:          body.IsNew,
		CountryID:      body.CountryID,
		Position:       body.Position,
		IsActive:       body.IsActive,
		DescriptionRu:  body.DescriptionRu,
		DescriptionUz:  body.DescriptionUz,
		DescriptionEn:  body.DescriptionEn,
		SeoTitle:       body.SeoTitle,
		Image:          body.Image,
		SeoDescription: body.SeoDescription,
		CreatedID:      &admin.Id,
		CreatedAt:      timeNow(),
	}
	if body.ParentID != 0 {
		product.ParentID = &body.ParentID
	}
	if body.BrandID != 0 {
		product.BrandID = &body.BrandID
	}
	err = h.db.Clauses(clause.Returning{}).Create(&product).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, product)
}

// @Summary		  Create product parameters
// @Description	   this api will create gotten parameters and delete other relations product parameters
// @Tags			Product
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param          	id    	path    int   true "product id"
// @Param			data 	body	models.ProductParamReq	true	"data body"
// @Success			201		{object}	models.Products
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product/parameter/{id} [POST]
func (h *ProductController) CreateProductParameter(c *gin.Context) {
	admin := h.GetAdmin(c)
	inputId := c.Param("id")
	id, _ := strconv.ParseInt(inputId, 10, 64)
	var body models.ProductParamReq
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	tr := h.db.Begin()
	err = tr.Delete(&models.ProductParameters{}, "product_id=?", id).Error
	if err != nil {
		tr.Rollback()
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to delete product params")
		return
	}
	params := make([]models.ProductParameters, len(body.Parameters))
	for i, param := range body.Parameters {
		params[i] = models.ProductParameters{
			ProductID:   int(id),
			ParameterID: param.ParameterID,
			ValRu:       param.ValRu,
			ValUz:       param.ValUz,
			ValEn:       param.ValEn,
		}
	}
	err = tr.Clauses(clause.Returning{}).Create(&params).Error
	if err != nil {
		tr.Rollback()
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to create ")
		return
	}
	err = tr.Model(&models.Products{}).Where("id=?", id).Updates(map[string]interface{}{
		"updated_id": admin.Id,
		"updated_at": timeNow(),
	}).Error
	if err != nil {
		tr.Rollback()
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update products ", err.Error())
		return
	}
	tr.Commit()
	c.JSON(http.StatusOK, params)

}

// @Summary		  Update product
// @Description	   this api is Update product
// @Tags			Product
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "product id"
// @Param			data 	formData		models.ProductRequest	true	"data body"
// @Param           image_file   formData	file				false	"file"
// @Success			201		{object}	models.Products
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product/{id} [PUT]
func (h *ProductController) UpdateProduct(c *gin.Context) {
	admin := h.GetAdmin(c)
	productId := c.Param("id")
	var body models.ProductRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create product", err.Error())
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
	product := models.Products{}
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
	if body.Position != nil {
		columns["position"] = body.Position
	}
	if body.IsNew != nil {
		columns["is_new"] = body.IsNew
	}
	if body.IsTop != nil {
		columns["is_top"] = body.IsTop
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
	if body.ParentID != 0 {
		columns["parent_id"] = body.ParentID
	}
	if body.BrandID != 0 {
		columns["brand_id"] = body.BrandID
	}
	if body.CountryID != nil {
		columns["brand_id"] = body.CountryID
	}
	columns["updated_at"] = timeNow()
	columns["updated_id"] = admin.Id
	err = h.db.Clauses(clause.Returning{}).Model(&product).
		Where("id=?", productId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save product", err.Error())
		return
	}
	c.JSON(http.StatusOK, product)
}
func (h *ProductController) updateProduct(productId int, adminId int) error {
	columns := map[string]interface{}{
		"updatedAt":  timeNow(),
		"updated_id": adminId,
	}
	return h.db.Model(&models.Products{ID: productId}).Updates(columns).Error
}

// @Summary		  Add product media
// @Description	   this api is add media to product
// @Tags			Product
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path   int   true   "product id"
// @Param			data 	formData		models.ProductMediaRequest	true	"data body"
// @Param			file	formData	file				true	"file"
// @Success			201		{object}	models.Products
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product/media/{id} [POST]
func (h *ProductController) AddProductMedia(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.ProductMediaRequest
	inputId := c.Param("id")
	id, err := strconv.Atoi(inputId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "Invalid id")
		return
	}
	err = c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	file, err := c.FormFile("file")
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathProducts, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}

		body.Media = name
	} else {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to get file", err.Error())
		return
	}
	media := models.ProductMedia{
		ProductID: &id,
		Type:      body.Type,
		Position:  &body.Position,
		Media:     body.Media,
	}
	err = h.db.Create(&media).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to create media", err.Error())
		return
	}
	h.updateProduct(id, admin.Id)
	newResponse(c, http.StatusOK, "Successfully created")
}

// @Summary		  Add product media
// @Description	   this api is add media to product
// @Tags			Product
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path   int   true   "media id"
// @Success			201		{object}	models.Products
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product/media/{id} [DELETE]
func (h *ProductController) DeleteProductMedia(c *gin.Context) {
	admin := h.GetAdmin(c)
	inputId := c.Param("id")
	id, err := strconv.Atoi(inputId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "Invalid id")
		return
	}
	media := models.ProductMedia{ID: id}
	err = h.db.First(&media).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.db.Delete(&media).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to create media", err.Error())
		return
	}
	err = h.filesService.Delete(c.Request.Context(), models.FilePathProducts, media.Media)
	if err != nil {
		h.log.Error("failed to delete media", err.Error())
	}
	h.updateProduct(id, admin.Id)
	newResponse(c, http.StatusOK, "Successfully deleted")
}

// @Summary		  get product
// @Description	   this api is get product
// @Tags			Product
// @Accept			json
// @Produce			json
// @Param 			data   body       models.ProductsIds  false "product brand ids"
// @Success			201		{object}	models.ProductsList
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product/list [POST]
func (h *ProductController) GetProductsByIds(c *gin.Context) {
	var body models.ProductsIds
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var products []models.Products
	db := h.db.Debug().Model(&models.Products{}).Where("is_active=true AND deleted_at IS NULL ")
	if len(body.ProductsIds) > 0 {
		db = db.Where("id IN ?", body.ProductsIds)
	}
	err = db.Order("position NULLS LAST").Select("*").Find(&products).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find products", err.Error())
		return
	}
	c.JSON(http.StatusOK, products)
}

// @Summary		  get product
// @Description	   this api is get product
// @Tags			Product
// @Accept			json
// @Produce			json
// @Param           data    query    	models.ProductsFilter   true   "product filter"
// @Param 			brandId query       array  false "product brand ids"
// @Param 			countryId query       array  false "product country ids"
// @Success			201		{object}	models.ProductsList
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product [GET]
func (h *ProductController) GetProducts(c *gin.Context) {
	var body models.ProductsFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	brandId := c.Query("brandId")
	countryId := c.Query("countryId")
	var products []models.Products
	db := h.db.Debug().Model(&models.Products{}).Where("is_active=true AND deleted_at IS NULL ")
	if brandId != "" {
		arr := strings.Split(brandId, ",")
		db = db.Where("brand_id IN ?", arr)
	}
	if body.ParentID != 0 {
		db = db.Where("parent_id=?", body.ParentID)
	}
	if body.PriceFrom != 0 {
		db = db.Where("price>=?", body.PriceFrom)
	}
	if body.PriceTo != 0 {
		db = db.Where("price<=?", body.PriceTo)
	}
	if body.IsNew != nil {
		db = db.Where("is_new=?", body.IsNew)
	}
	if body.IsTop != nil {
		db = db.Where("is_top=?", body.IsTop)
	}
	if countryId != "" {
		arr := strings.Split(countryId, ",")

		db = db.Where("country_id IN ?", arr)
	}
	if body.MultiSearch != "" {
		field := fmt.Sprintf("%%%s%%", body.MultiSearch)
		db = db.Where(`LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?) `, field, field)
	}
	var count int
	err = db.Select("COUNT(*)").Scan(&count).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Errorf("failed to get products %v", err)
		return
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	db = db.Offset((body.Page - 1) * body.PageSize).Limit(body.PageSize)
	err = db.Order("position NULLS LAST").Select("*").Find(&products).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find products", err.Error())
		return
	}
	c.JSON(http.StatusOK, models.ProductsList{
		Products: products,
		Page:     body.Page,
		PageSize: body.PageSize,
		Count:    count,
	})
}

// @Summary		  get product
// @Description	   this api is get product
// @Tags			Product
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           data    query    	models.ProductsFilter   true   "product filter"
// @Param 			brandId query       array  false "product brand ids"
// @Param 			countryId query       array  false "product country ids"
// @Success			201		{object}	models.ProductsList
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product/all [GET]
func (h *ProductController) GetAllProducts(c *gin.Context) {
	var body models.ProductsFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	brandId := c.Query("brandId")
	countryId := c.Query("countryId")
	var products []models.Products
	db := h.db.Debug().Model(&models.Products{})
	if brandId != "" {
		arr := strings.Split(brandId, ",")
		db = db.Where("brand_id IN ?", arr)
	}
	if body.ParentID != 0 {
		db = db.Where("parent_id=?", body.ParentID)
	}
	if body.PriceFrom != 0 {
		db = db.Where("price>=?", body.PriceFrom)
	}
	if body.PriceTo != 0 {
		db = db.Where("price<=?", body.PriceTo)
	}
	if body.IsNew != nil {
		db = db.Where("is_new=?", body.IsNew)
	}
	if body.IsTop != nil {
		db = db.Where("is_top=?", body.IsTop)
	}
	if countryId != "" {
		arr := strings.Split(countryId, ",")

		db = db.Where("country_id IN ?", arr)
	}
	if body.MultiSearch != "" {
		field := fmt.Sprintf("%%%s%%", body.MultiSearch)
		db = db.Where(`LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?) `, field, field)
	}
	var count int
	err = db.Select("COUNT(*)").Scan(&count).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Errorf("failed to get products %v", err)
		return
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	db = db.Offset((body.Page - 1) * body.PageSize).Limit(body.PageSize)
	err = db.Order("position NULLS LAST").Select("*").Find(&products).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find products", err.Error())
		return
	}
	c.JSON(http.StatusOK, models.ProductsList{
		Products: products,
		Page:     body.Page,
		PageSize: body.PageSize,
		Count:    count,
	})
}

// @Summary		  Update product
// @Description	   this api is Update product
// @Tags			Product
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "product id"
// @Success			201		{object}	models.ProductResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product/{id} [GET]
func (h *ProductController) GetByID(c *gin.Context) {
	inputId := c.Param("id")
	var product models.Products
	err := h.db.Model(&models.Products{}).Preload("Country").Preload("Parent").
		First(&product, "id=?", inputId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found product")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var media []models.ProductMedia
	err = h.db.Find(&media, "product_id=?", inputId).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get media")
		h.log.Error("failed to get product media", err.Error())
		return
	}
	var parameters []models.ProductParameters
	err = h.db.Model(&models.ProductParameters{}).Preload("Parameter").Find(&parameters, "product_id=?", inputId).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get media")
		h.log.Error("failed to get product media", err.Error())
		return
	}
	c.JSON(http.StatusOK, models.ProductResponse{
		Products:   &product,
		Media:      media,
		Parameters: parameters,
	})
}

// @Summary		  delete product from favorites
// @Description	   this api is for delete product from favorites
// @Tags			Product
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path    int   true  "product id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/product/{id} [DELETE]
func (h *ProductController) DeleteProduct(c *gin.Context) {
	admin := h.GetAdmin(c)
	inputId := c.Param("id")
	if inputId == "" {
		newResponse(c, http.StatusBadRequest, "empty id")
		return
	}
	err := h.db.Model(&models.Products{}).Where("id=?", inputId).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_id": admin.Id,
	}).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to delete customer favorites", err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// // @Summary		  Add product to favorites
// // @Description	   this api is for add product to favorites
// // @Tags			Product
// // @Security		BearerAuth
// // @Accept			json
// // @Produce			json
// // @Param           id    path    int   true  "product id"
// // @Success			201		{object}	response
// // @Failure			400,409	{object}	response
// // @Failure			500		{object}	response
// // @Router			/api/product/favorite/{id} [POST]
// func (h *ProductController) AddFavoriteProduct(c *gin.Context) {
// 	customer := h.GetCustomer(c)
// 	inputId := c.Param("id")
// 	if inputId == "" {
// 		newResponse(c, http.StatusBadRequest, "empty id")
// 		return
// 	}
// 	id, err := strconv.ParseInt(inputId, 10, 64)
// 	if err != nil {
// 		newResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	prodId := int(id)
// 	err = h.db.Create(&models.CustomerFavorites{
// 		CustomerID: &customer.Id,
// 		ProductID:  &prodId,
// 	}).Error
// 	if err != nil {
// 		newResponse(c, http.StatusInternalServerError, err.Error())
// 		h.log.Error("failed to add customer favorites", err.Error())
// 		return
// 	}
// 	c.JSON(http.StatusOK, response{"success"})
// }

// // @Summary		  get customer favorites products
// // @Description	   this api is for get customer favorites products
// // @Tags			Product
// // @Security		BearerAuth
// // @Accept			json
// // @Produce			json
// // @Success			201		{object}	[]models.Products
// // @Failure			400,409	{object}	response
// // @Failure			500		{object}	response
// // @Router			/api/product/favorite  [GET]
// func (h *ProductController) GetFavoriteProduct(c *gin.Context) {
// 	customer := h.GetCustomer(c)
// 	var favorites []models.Products
// 	err := h.db.Table("products AS p").Select("p.*").Joins("INNER JOIN customer_favorites AS cf ON p.id=cf.product_id").
// 		Where("cf.customer_id=?", customer.Id).Find(&favorites).Error
// 	if err != nil {
// 		newResponse(c, http.StatusInternalServerError, "failed to get favorites")
// 		h.log.Error("failed to get favorites", logger.Error(err))
// 		return
// 	}
// 	c.JSON(http.StatusOK, favorites)
// }
