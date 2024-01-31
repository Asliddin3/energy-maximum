package controller

import (
	"errors"
	"net/http"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CountryController struct {
	*Handler
}

func (h *Handler) NewCountryController(api *gin.RouterGroup) {
	brand := &CountryController{h}
	br := api.Group("country", h.DeserializeAdmin())
	{
		br.POST("/", brand.CreateCountry)
		br.PUT("/:id", brand.UpdateCountry)
		br.DELETE("/:id", brand.DeleteCountry)
	}
	api.GET("/country", brand.GetCountries)
	api.GET("/country/:id", brand.GetByID)

}

// @Summary		  Create new country
// @Description	   this api is create new country
// @Tags			Country
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data    body		models.CountryRequest	false	"data body"
// @Success			201		{object}	models.Country
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/country/ [POST]
func (h *CountryController) CreateCountry(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.CountryRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	country := models.Country{
		NameRu:    body.NameRu,
		NameUz:    body.NameUz,
		NameEn:    body.NameEn,
		CreatedID: &admin.Id,
		CreatedAt: timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Create(&country).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, country)
}

// @Summary		  Update country
// @Description	   this api is Update country
// @Tags			Country
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "update id"
// @Param			data 	query		models.CountryRequest	true	"data body"
// @Success			201		{object}	models.Country
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/country/{id} [PUT]
func (h *CountryController) UpdateCountry(c *gin.Context) {
	admin := h.GetAdmin(c)
	id := c.Param("id")
	var body models.CountryRequest
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	country := models.Country{}
	err = h.db.First(&country, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to create country", err.Error())
		return
	}
	country.NameRu = body.NameRu
	country.NameEn = body.NameEn
	country.NameUz = body.NameUz
	country.CreatedAt = timeNow()
	country.CreatedID = &admin.Id
	err = h.db.Save(&country).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save brand", err.Error())
		return
	}
	c.JSON(http.StatusOK, country)
}

// @Summary		    Get countries
// @Description	    this api is to get countries
// @Tags			Country
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Country
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/country  [GET]
func (h *CountryController) GetCountries(c *gin.Context) {
	var countries []models.Country
	err := h.db.Find(&countries).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, countries)
}

// @Summary		  Get country
// @Description	   this api is to get country
// @Tags			Country
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "id"
// @Success			201		{object}	models.Country
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/country/{id} [GET]
func (h *CountryController) GetByID(c *gin.Context) {
	brandId := c.Param("id")
	var brands models.Country
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

// @Summary		  Delete country
// @Description	   this api is to delete country
// @Tags			Country
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/country/{id} [DELETE]
func (h *CountryController) DeleteCountry(c *gin.Context) {
	id := c.Param("id")
	err := h.db.Delete(&models.Country{}, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
