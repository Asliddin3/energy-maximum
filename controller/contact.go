package controller

import (
	"errors"
	"net/http"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContactController struct {
	*Handler
}

func (h *Handler) NewContactController(api *gin.RouterGroup) {
	brand := &ContactController{h}
	br := api.Group("contact")
	{
		br.POST("/", h.DeserializeAdmin(), brand.CreateContact)
		br.PUT("/:id", h.DeserializeAdmin(), brand.UpdateContact)
		br.GET("", brand.GetContacts)
		br.GET("/:id", brand.GetByID)
		br.GET("/main", brand.GetMainFilial)
		br.DELETE("/:id", h.DeserializeAdmin(), brand.DeleteContact)
	}
}

// @Summary		  Create new contact
// @Description	   this api is create new contact
// @Tags			Contact
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data    body		models.ContactRequest	false	"data body"
// @Success			201		{object}	models.Contact
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/contact/ [POST]
func (h *ContactController) CreateContact(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.ContactRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	contact := models.Contact{
		NameRu:       body.NameRu,
		NameUz:       body.NameUz,
		NameEn:       body.NameEn,
		Address:      body.Address,
		IsMain:       body.IsMain,
		Phone:        body.Phone,
		WorkingHours: body.WorkingHours,
		Email:        body.Email,
		Longitude:    body.Longitude,
		Latitude:     body.Latitude,
		CreatedID:    &admin.Id,
		CreatedAt:    timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Create(&contact).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, contact)
}

// @Summary		  Update contact
// @Description	   this api is Update contact
// @Tags			Contact
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path      	string   true   "update id"
// @Param			data 	body		models.ContactRequest	true	"data body"
// @Success			201		{object}	models.Contact
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/contact/{id} [PUT]
func (h *ContactController) UpdateContact(c *gin.Context) {
	admin := h.GetAdmin(c)
	id := c.Param("id")
	var body models.ContactRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	country := models.Contact{}
	err = h.db.First(&country, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to create country", err.Error())
		return
	}
	if body.Address != "" {
		country.Address = body.Address
	}
	if body.Phone != "" {
		country.Phone = body.Phone
	}
	if body.NameEn != "" {
		country.NameEn = body.NameEn
	}
	if body.NameRu != "" {
		country.NameRu = body.NameRu
	}
	if body.NameUz != "" {
		country.NameUz = body.NameUz
	}
	if body.Email != "" {
		country.Email = body.Email
	}
	if body.WorkingHours != "" {
		country.WorkingHours = body.WorkingHours
	}
	if body.IsMain != nil {
		country.IsMain = body.IsMain
	}
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

// @Summary		    Get contacts
// @Description	    this api is to get contacts
// @Tags			Contact
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Contact
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/contact  [GET]
func (h *ContactController) GetContacts(c *gin.Context) {
	var countries []models.Contact
	err := h.db.Find(&countries).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, countries)
}

// @Summary		  Get contacts
// @Description	   this api is to get contacts
// @Tags			Contact
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "id"
// @Success			201		{object}	models.Contact
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/contact/{id} [GET]
func (h *ContactController) GetByID(c *gin.Context) {
	contactId := c.Param("id")
	var contacts models.Contact
	err := h.db.First(&contacts, "id=?", contactId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found brand")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, contacts)
}

// @Summary		  Get contacts
// @Description	   this api is to get contacts
// @Tags			Contact
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Contact
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/contact/main [GET]
func (h *ContactController) GetMainFilial(c *gin.Context) {
	var contacts []models.Contact
	err := h.db.Find(&contacts, "is_main=true").Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, contacts)
}

// @Summary		  Delete country
// @Description	   this api is to delete country
// @Tags			Contact
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/contact/{id} [DELETE]
func (h *ContactController) DeleteContact(c *gin.Context) {
	id := c.Param("id")
	err := h.db.Delete(&models.Contact{}, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
