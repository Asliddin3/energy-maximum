package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AnalogController struct {
	*Handler
}

func (h *Handler) NewAnalogController(api *gin.RouterGroup) {
	analog := &AnalogController{h}
	analogProd := api.Group("analog", h.DeserializeAdmin())
	{
		analogProd.POST("", analog.CreateAnalog)
		analogProd.PUT("/:id", analog.UpdateAnalog)
		analogProd.GET("/:id", analog.GetByID)
		analogProd.GET("", analog.GetAnalogs)
		analogProd.DELETE("/:id", analog.DeleteAnalog)
		analogProd.PUT("/products/:id", analog.AddAnalogProducts)
		analogProd.DELETE("/products/:id", analog.DeleteAnalogProducts)
		analogProd.GET("/products/:id", analog.GetProductAnalogs)
	}
}

// @Summary		  Create analog
// @Description	   this api is create analog
// @Tags			Analogs
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 		body		models.AnalogRequest	true	"data body"
// @Success			201		{object}	models.Analog
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/analog [POST]
func (h *AnalogController) CreateAnalog(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.AnalogRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create news", err.Error())
		return
	}

	news := models.Analog{
		NameRu:    body.NameRu,
		NameUz:    body.NameUz,
		NameEn:    body.NameEn,
		CreatedID: &admin.Id,
		CreatedAt: timeNow(),
	}
	err = h.db.Clauses(clause.Returning{}).Create(&news).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, news)
}

// @Summary		  	Update Analog
// @Description	   	this api is Update Analog
// @Tags			Analogs
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "Analog id"
// @Param			data 	query		models.AnalogRequest	true	"data body"
// @Success			201		{object}	models.Analog
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/analog/{id} [PUT]
func (h *AnalogController) UpdateAnalog(c *gin.Context) {
	admin := h.GetAdmin(c)
	analogId := c.Param("id")
	var body models.AnalogRequest
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create news", err.Error())
		return
	}
	columns := map[string]interface{}{}

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
	analog := &models.Analog{}
	err = h.db.Clauses(clause.Returning{}).Model(&analog).
		Where("id=?", analogId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save analog", err.Error())
		return
	}
	c.JSON(http.StatusOK, analog)
}

// @Summary		  Get analogs
// @Description	   this api is get analogs
// @Tags			Analogs
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Analog
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/analog [GET]
func (h *AnalogController) GetAnalogs(c *gin.Context) {
	var analogs []models.Analog
	err := h.db.Find(&analogs).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get analogs")
		h.log.Error("failed to get analogs", err.Error())
		return
	}
	c.JSON(http.StatusOK, analogs)
}

// @Summary		  Get analog
// @Description	   this api is get to analogs
// @Tags			Analogs
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "analog id"
// @Success			201		{object}	models.AnalogResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/analog/{id} [GET]
func (h *AnalogController) GetByID(c *gin.Context) {
	inputId := c.Param("id")
	var analog models.Analog
	err := h.db.First(&analog, "id=?", inputId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found analog")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var products []models.Products
	err = h.db.Table("products AS p").Select("p.*").Joins("INNER JOIN analog_product AS pa ON pa.product_id=p.id ").
		Where("pa.analog_id=?", analog.ID).Find(&products).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get analog products")
		h.log.Error("failed to get products", err.Error())
		return
	}
	c.JSON(http.StatusOK, models.AnalogResponse{
		Analog: &analog,
		Items:  products,
	})
}

// @Summary		  DELETE news
// @Description	   this api is delete news
// @Tags			Analogs
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "news id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/analog/{id} [DELETE]
func (h *AnalogController) DeleteAnalog(c *gin.Context) {
	inputId := c.Param("id")
	var analog models.Analog
	err := h.db.Delete(&analog, "id=?", inputId).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Add analog products
// @Description	   this api is  Add analog products
// @Tags			Analogs
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "product id"
// @Success			201		{object}	[]models.Products
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/analog/products/{id} [GET]
func (h *AnalogController) GetProductAnalogs(c *gin.Context) {
	inputId := c.Param("id")
	id, _ := strconv.ParseInt(inputId, 10, 64)
	var products []models.Products
	err := h.db.Table("products AS p").Select("p.*").Joins(`INNER JOIN analog_product AS ap ON ap.product_id = p.id 
	AND ap.analog_id IN (SELECT analog_id FROM analog_product WHERE product_id=?)`, id).Where("ap.product_id!=?", id).Find(&products).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get product analog", err.Error())
		return
	}
	c.JSON(http.StatusOK, products)
}

// @Summary		  Add analog products
// @Description	   this api is  Add analog products
// @Tags			Analogs
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "news id"
// @Param           products  query   array true "new products add"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/analog/products/{id} [PUT]
func (h *AnalogController) AddAnalogProducts(c *gin.Context) {
	inputId := c.Param("id")
	id, _ := strconv.ParseInt(inputId, 10, 64)
	products := c.Query("products")
	if products == "" {
		newResponse(c, http.StatusBadRequest, "empty products")
		return
	}
	prodArr := strings.Split(products, ",")
	items := make([]models.AnalogProduct, len(prodArr))
	for i, prod := range prodArr {
		prodId, _ := strconv.ParseInt(prod, 10, 64)
		items[i] = models.AnalogProduct{
			AnalogID:  int(id),
			ProductID: int(prodId),
		}
	}
	err := h.db.Create(&items).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to create items", err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Add analog products
// @Description	   this api is  Add analog products
// @Tags			Analogs
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "news id"
// @Param           products  query   array true "new products add"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/analog/products/{id} [DELETE]
func (h *AnalogController) DeleteAnalogProducts(c *gin.Context) {
	inputId := c.Param("id")
	id, _ := strconv.ParseInt(inputId, 10, 64)
	products := c.Query("products")
	if products == "" {
		newResponse(c, http.StatusBadRequest, "empty products")
		return
	}
	prodArr := strings.Split(products, ",")

	err := h.db.Delete(&models.AnalogProduct{}, "analog_id=? AND product_id IN (?)", id, prodArr).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to delete items", err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
