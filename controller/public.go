package controller

import (
	"errors"
	"net/http"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PublicOfferController struct {
	*Handler
}

func (h *Handler) NewPublicOfferController(api *gin.RouterGroup) {
	brand := &PublicOfferController{h}
	br := api.Group("public-offer", h.DeserializeAdmin())
	{
		br.POST("/", brand.CreatePublicOffer)
		// br.PUT("/", brand.UpdatePublic)
		// br.GET("", brand.GetAbout)
		br.DELETE("/", brand.DeletePublic)
	}
	api.GET("/public-offer", brand.GetByID)
}

// @Summary		  Create new public offer
// @Description	   this api is create new public
// @Tags			PublicOffer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data    body		models.PublicOfferRequest	false	"data body"
// @Success			201		{object}	models.PublicOffer
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/public-offer/ [POST]
func (h *PublicOfferController) CreatePublicOffer(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.PublicOfferRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create about", err.Error())
		return
	}
	offer := models.PublicOffer{}
	err = h.db.First(&offer).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	about := models.PublicOffer{
		ID:            offer.ID,
		NameRu:        body.NameRu,
		NameUz:        body.NameUz,
		NameEn:        body.NameEn,
		DescriptionUz: body.DescriptionUz,
		DescriptionRu: body.DescriptionRu,
		DescriptionEn: body.DescriptionEn,
		CreatedID:     &admin.Id,
		CreatedAt:     timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Save(&about).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, about)
}

// // @Summary		  Update public offer
// // @Description	   this api is Update public offer
// // @Tags			PublicOffer
// // @Security		BearerAuth
// // @Accept			json
// // @Produce			json
// // @Param           id    path     string   true   "update id"
// // @Param			data 	body		models.PublicOfferRequest	true	"data body"
// // @Success			201		{object}	models.About
// // @Failure			400,409	{object}	response
// // @Failure			500		{object}	response
// // @Router			/api/public-offer/ [PUT]
// func (h *PublicOfferController) UpdatePublic(c *gin.Context) {
// 	admin := h.GetAdmin(c)
// 	var body models.PublicOfferRequest
// 	err := c.ShouldBind(&body)
// 	if err != nil {
// 		newResponse(c, http.StatusBadRequest, err.Error())
// 		h.log.Error("failed to create about", err.Error())
// 		return
// 	}
// 	offer := models.PublicOffer{}
// 	err = h.db.First(&offer).Error
// 	if err != nil {
// 		if !errors.Is(err, gorm.ErrRecordNotFound) {
// 			newResponse(c, http.StatusInternalServerError, err.Error())
// 		}
// 	}

// 	columns := map[string]interface{}{}

// 	about := models.PublicOffer{
// 		ID: int(offer.ID),
// 	}
// 	if body.DescriptionEn != "" {
// 		columns["description_en"] = body.DescriptionEn
// 	}
// 	if body.DescriptionRu != "" {
// 		columns["description_ru"] = body.DescriptionRu
// 	}
// 	if body.DescriptionUz != "" {
// 		columns["description_uz"] = body.DescriptionUz
// 	}
// 	if body.NameEn != "" {
// 		columns["name_en"] = body.NameEn
// 	}
// 	if body.NameRu != "" {
// 		columns["name_ru"] = body.NameRu
// 	}
// 	if body.NameUz != "" {
// 		columns["name_uz"] = body.NameUz
// 	}

// 	columns["updated_at"] = timeNow()
// 	columns["updated_id"] = admin.Id
// 	err = h.db.Debug().Clauses(clause.Returning{}).Model(&about).
// 		Where("id=?", offer.ID).Updates(columns).Error
// 	if err != nil {
// 		newResponse(c, http.StatusInternalServerError, err.Error())
// 		h.log.Error("failed to save about", err.Error())
// 		return
// 	}
// 	c.JSON(http.StatusOK, about)
// }

// @Summary		  Get public offer
// @Description	   this api is to get public offer
// @Tags			PublicOffer
// @Accept			json
// @Produce			json
// @Success			201		{object}	models.PublicOffer
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/public-offer [GET]
func (h *PublicOfferController) GetByID(c *gin.Context) {
	var about models.PublicOffer
	err := h.db.First(&about).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found public offer")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, about)
}

// @Summary		  Delete public offer
// @Description	   this api is to delete public offer
// @Tags			PublicOffer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/public-offer/ [DELETE]
func (h *PublicOfferController) DeletePublic(c *gin.Context) {
	err := h.db.Delete(&models.PublicOffer{}, "created_id IS NOT NUlL").Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
