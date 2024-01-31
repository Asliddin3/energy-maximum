package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VacancyController struct {
	*Handler
}

func (h *Handler) NewVacancyController(api *gin.RouterGroup) {
	vacancy := &VacancyController{h}
	vac := api.Group("vacancy", h.DeserializeAdmin())
	{
		vac.POST("", vacancy.CreateVacancy)
		vac.PUT("/:id", vacancy.UpdateVacancy)
		vac.DELETE("/:id", vacancy.Delete)
		vac.GET("/all", vacancy.GetVacancy)
		vac.GET("/applicant/:id", vacancy.GetApplicantById)
		vac.GET("/applicant", vacancy.GetApplicant)
		vac.PUT("/applicant/reject/:id", vacancy.RejectedApplicant)
		vac.PUT("/applicant/accept/:id", vacancy.AcceptApplicant)
		vac.PUT("/applicant/save/:id", vacancy.SaveApplicant)
		vac.DELETE("/applicant/:id", vacancy.DeleteApplicantById)
	}
	api.GET("/vacancy", vacancy.GetCustomerVacancy)
	api.GET("/vacancy/:id", vacancy.GetByID)
	api.POST("/vacancy/applicant", vacancy.CreateApplicant)
}

// @Summary		  Create vacancy
// @Description	   this api is create vacancy
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	formData		models.VacancyRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Vacancy
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy [POST]
func (h *VacancyController) CreateVacancy(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.VacancyRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create vacancy", err.Error())
		return
	}
	file, err := c.FormFile("image_file")
	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathVacancy, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}

		body.Image = name

	}
	vacancy := models.Vacancy{
		NameRu:           body.NameRu,
		NameUz:           body.NameUz,
		NameEn:           body.NameEn,
		ResponsibilityRu: body.ResponsibilityRu,
		ResponsibilityUz: body.ResponsibilityUz,
		ResponsibilityEn: body.ResponsibilityEn,
		Region:           body.Region,
		RequirementRu:    body.RequirementRu,
		RequirementUz:    body.RequirementUz,
		RequirementEn:    body.RequirementEn,
		DescriptionRu:    body.DescriptionRu,
		DescriptionUz:    body.DescriptionUz,
		DescriptionEn:    body.DescriptionEn,
		TypeUz:           body.TypeUz,
		TypeEn:           body.TypeEn,
		TypeRu:           body.TypeRu,
		IsActive:         body.IsActive,
		Image:            body.Image,
		CreatedID:        &admin.Id,
		CreatedAt:        timeNow(),
	}
	err = h.db.Clauses(clause.Returning{}).Create(&vacancy).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, vacancy)
}

// @Summary		  Update vacancy
// @Description	   this api is Update vacancy
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "vacancy id"
// @Param			data 	formData		models.VacancyRequest	true	"data body"
// @Param			image_file	formData	file				false	"file"
// @Success			201		{object}	models.Vacancy
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/{id} [PUT]
func (h *VacancyController) UpdateVacancy(c *gin.Context) {
	admin := h.GetAdmin(c)
	vacancyId := c.Param("id")
	var body models.VacancyRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create vacancy", err.Error())
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
	if body.Region != "" {
		columns["region"] = body.Region
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
	if body.RequirementEn != "" {
		columns["requirement_en"] = body.RequirementEn
	}
	if body.RequirementRu != "" {
		columns["requirement_ru"] = body.RequirementRu
	}
	if body.RequirementUz != "" {
		columns["requirement_uz"] = body.RequirementUz
	}
	if body.ResponsibilityEn != "" {
		columns["responsibility_en"] = body.ResponsibilityEn
	}
	if body.ResponsibilityRu != "" {
		columns["responsibility_ru"] = body.ResponsibilityRu
	}
	if body.TypeEn != "" {
		columns["type_uz"] = body.TypeEn
	}
	if body.TypeRu != "" {
		columns["type_ru"] = body.TypeRu
	}
	if body.TypeUz != "" {
		columns["type_uz"] = body.TypeUz
	}
	if body.IsActive != nil {
		columns["is_active"] = body.IsActive
	}
	columns["updated_at"] = timeNow()
	columns["updated_id"] = admin.Id
	vacancy := &models.Vacancy{}
	err = h.db.Clauses(clause.Returning{}).Model(&vacancy).
		Where("id=?", vacancyId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to save vacancy", err.Error())
		return
	}
	c.JSON(http.StatusOK, vacancy)
}

// @Summary		  Get all vacancy
// @Description	   this api is get all vacancy
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Vacancy
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/all [GET]
func (h *VacancyController) GetVacancy(c *gin.Context) {
	var vacancy []models.Vacancy
	err := h.db.Find(&vacancy, "is_deleted=false").Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find vacancy", err.Error())
		return
	}
	c.JSON(http.StatusOK, vacancy)
}

// @Summary		  Get vacancy
// @Description	   this api is get vacancy
// @Tags			Vacancy
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Vacancy
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy [GET]
func (h *VacancyController) GetCustomerVacancy(c *gin.Context) {
	var vacancy []models.Vacancy
	err := h.db.Find(&vacancy, "is_active=true AND is_deleted=false").Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to find vacancy", err.Error())
		return
	}
	c.JSON(http.StatusOK, vacancy)
}

// @Summary		  Update vacancy
// @Description	   this api is Update vacancy
// @Tags			Vacancy
// @Accept			json
// @Produce			json
// @Param           id    path     string   true   "vacancy id"
// @Success			201		{object}	models.Vacancy
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/{id} [GET]
func (h *VacancyController) GetByID(c *gin.Context) {
	inputId := c.Param("id")
	var vacancy models.Vacancy
	err := h.db.First(&vacancy, "id=?", inputId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found vacancy")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, vacancy)
}

// @Summary		  DELETE vacancy
// @Description	   this api is DELETE vacancy
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    	path     int   true   "vacancy id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/{id} [DELETE]
func (h *VacancyController) Delete(c *gin.Context) {
	inputId := c.Param("id")
	var vacancy models.Vacancy
	err := h.db.Debug().Model(&vacancy).Where("id=?", inputId).Update("is_deleted", true).Error
	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	newResponse(c, http.StatusBadRequest, "not found vacancy")
		// 	return
		// }
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Create applicant to vacancy
// @Description	   this api is create applicant to vacancy
// @Tags			Vacancy
// @Accept			json
// @Produce			json
// @Param			data 	formData		models.ApplicantRequest	true	"data body"
// @Param           resume_file   formData	file				true	"file"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/applicant [POST]
func (h *VacancyController) CreateApplicant(c *gin.Context) {
	var body models.ApplicantRequest
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	file, err := c.FormFile("resume_file")

	if err == nil {
		name, err := h.filesService.Save(c.Request.Context(), models.File{Path: models.FilePathApplicant, File: file})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save image")
			h.log.Error("error while save file ", logger.Error(err))
			return
		}
		body.Resume = name
	}
	applicant := models.Applicant{
		Name:        body.Name,
		Phone:       body.Phone,
		Description: body.Description,
		VacancyID:   &body.VacancyId,
		Resume:      body.Resume,
		CreatedAt:   timeNow(),
	}
	err = h.db.Clauses(clause.Returning{}).Create(&applicant).Error
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			newResponse(c, http.StatusBadRequest, "not found vacancy")
			return
		}
		newResponse(c, http.StatusInternalServerError, "failed to create applicant")
		h.log.Error("failed to create applicant", err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Create applicant to vacancy
// @Description	   this api is get applicant of vacancy statuses(0--new , 1-accept,2-saved,3-rejected)
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           data     query   models.ApplicantFilter  true "applicant filter"
// @Success			201		{object}	[]models.Applicant
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/applicant [GET]
func (h *VacancyController) GetApplicant(c *gin.Context) {
	var body models.ApplicantFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var applicants []models.Applicant
	db := h.db.Debug().Model(&models.Applicant{})
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	if body.VacancyId != 0 {
		db = db.Where("vacancy_id = ?", body.VacancyId)
	}
	if body.Status != nil {
		db = db.Where("status = ?", body.Status)
	}
	err = db.Limit(body.PageSize).Offset((body.Page - 1) * body.PageSize).Find(&applicants).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get applicant")
		h.log.Error("Failed to get application", err.Error())
		return
	}
	c.JSON(http.StatusOK, applicants)
}

// @Summary		 	Reject applicant
// @Description	   this api is for change status applicant to 3
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    	path     int   true   "applicant id"
// @Success			201		{object}	models.Applicant
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/applicant/reject/{id} [PUT]
func (h *VacancyController) RejectedApplicant(c *gin.Context) {
	id := c.Param("id")
	applicant := models.Applicant{}
	err := h.db.Clauses(clause.Returning{}).Model(&applicant).Where("id=?", id).
		UpdateColumn("status", 4).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to update status")
		h.log.Error("failed to update status", err.Error())
		return
	}
	c.JSON(http.StatusOK, applicant)
}

// @Summary		 	Accept applicant
// @Description	   this api is for change status applicant to 1
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    	path     int   true   "applicant id"
// @Success			201		{object}	models.Applicant
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/applicant/accept/{id} [PUT]
func (h *VacancyController) AcceptApplicant(c *gin.Context) {
	id := c.Param("id")
	applicant := models.Applicant{}
	err := h.db.Clauses(clause.Returning{}).Model(&applicant).Where("id=?", id).
		UpdateColumn("status", 2).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to update status")
		h.log.Error("failed to update status", err.Error())
		return
	}
	c.JSON(http.StatusOK, applicant)
}

// @Summary		  Save applicant
// @Description	   this api is change applicant status to 2
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    	path     int   true   "applicant id"
// @Success			201		{object}	models.Applicant
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/applicant/save/{id} [PUT]
func (h *VacancyController) SaveApplicant(c *gin.Context) {
	id := c.Param("id")
	applicant := models.Applicant{}
	err := h.db.Clauses(clause.Returning{}).Model(&applicant).Where("id=?", id).
		UpdateColumn("status", 3).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to update status")
		h.log.Error("failed to update status", err.Error())
		return
	}
	c.JSON(http.StatusOK, applicant)
}

// @Summary		  Create applicant to vacancy
// @Description	   this api is create applicant to vacancy
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    	path     int   true   "vacancy id"
// @Success			201		{object}	models.Applicant
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/applicant/{id} [GET]
func (h *VacancyController) GetApplicantById(c *gin.Context) {
	id := c.Param("id")
	var applicants models.Applicant
	err := h.db.First(&applicants, "id=?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found applicant")
			return
		}
		newResponse(c, http.StatusInternalServerError, "failed to applicant")
		h.log.Error("Failed to get application", err.Error())
		return
	}
	c.JSON(http.StatusOK, applicants)
}

// @Summary		  Create applicant to vacancy
// @Description	   this api is create applicant to vacancy
// @Tags			Vacancy
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    	path     int   true   "vacancy id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/vacancy/applicant/{id} [DELETE]
func (h *VacancyController) DeleteApplicantById(c *gin.Context) {
	id := c.Param("id")
	var applicants models.Applicant
	err := h.db.Delete(&applicants, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to applicant")
		h.log.Error("Failed to get application", err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
