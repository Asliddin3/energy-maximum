package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Asliddin3/energy-maximum/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModuleController struct {
	*Handler
}

func (h *Handler) NewModuleController(rg *gin.RouterGroup) {
	controller := ModuleController{h}
	router := rg.Group("modules", h.DeserializeAdmin())
	router.POST("", controller.ModuleCreate)
	router.PUT("/:id", controller.ModuleUpdate)
	router.GET("", controller.ModulesGet)
	router.GET("/:id", controller.ModuleGetByID)
	router.DELETE("/:id", controller.ModuleDelete)
	routerItems := rg.Group("modules-items", h.DeserializeAdmin())
	routerItems.POST("", controller.ModuleItemsCreate)
	routerItems.PUT("/:id", controller.ModuleItemUpdate)
	routerItems.GET("", controller.ModulesItemsGet)
	routerItems.GET("/:id", controller.ModuleItemGetByID)
	routerItems.DELETE("/:id", controller.ModuleItemDelete)
}

// ModuleCreate
// @Summary			Create a Module.
// @Description		This API to create a Module.
// @Tags			  Module
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data	body	    models.CreateModuleInput	true	"data body"
// @Success			201		{object}	models.Modules
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/modules [POST]
func (h *ModuleController) ModuleCreate(c *gin.Context) {

	var body models.CreateModuleInput

	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	module := models.Modules{
		Name:        body.Name,
		Description: body.Description,
	}

	err = h.db.Select(
		"name",
		"description",
	).Create(&module).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update", err.Error())
		return
	}

	c.JSON(http.StatusCreated, module)
}

// ModuleUpdate
// @Summary		Update an Module.
// @Description	This API to update an Module.
// @Tags			  Module
// @Security		BearerAuth
// @Accept			json
// @Produce		json
// @Param			id		path		int		true	"id for update Module"	Format(id)
// @Param			data	body		models.UpdateModuleInput	true	"data body"
// @Success		  200		{object}	models.Modules
// @Failure		   400,409	{object}	response
// @Failure		   500		{object}	response
// @Router			/api/modules/{id} [PUT]
func (h *ModuleController) ModuleUpdate(c *gin.Context) {

	var body models.UpdateModuleInput

	inputID := c.Param("id")

	ID, err := strconv.ParseUint(inputID, 10, 32)
	if inputID == "" || err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = c.ShouldBind(&body)
	if err != nil || inputID == "" {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	module := models.Modules{
		ID: uint32(ID),
	}

	columns := map[string]interface{}{
		"name":        body.Name,
		"description": body.Description,
	}

	err = h.db.Clauses(clause.Returning{}).Model(&module).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update menu", err.Error())
		return
	}

	c.JSON(http.StatusOK, module)
}

// ModulesGet
// @Summary			Get all Modules
// @Description		This api for get Modules
// @Tags			  Module
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			200	{object}	[]models.Modules
// @Failure			500	{object}	response
// @Router			/api/modules [GET]
func (h *ModuleController) ModulesGet(c *gin.Context) {
	var modules []models.Modules

	err := h.db.Model(models.Modules{}).Find(&modules).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get modules", err.Error())
		return
	}

	c.JSON(http.StatusOK, modules)
}

// ModuleGetByID
//
//	@Summary		Get module by id.
//	@Description	This API to get module by id.
//	@Tags			  Module
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"id for get module"	Format(id)
//	@Success		200		{object}	models.Modules
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Router			/api/modules/{id} [GET]
func (h *ModuleController) ModuleGetByID(c *gin.Context) {

	var module models.Modules

	inputID := c.Param("id")

	ID, err := strconv.ParseUint(inputID, 10, 32)
	if inputID == "" || err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.db.Model(models.Modules{}).First(&module, "id = ?", ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get by id", err.Error())
		return
	}

	c.JSON(http.StatusOK, module)
}

// ModuleDelete
// @Summary		Get module by id.
// @Description	This API to get module by id.
// @Tags	    Module
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Param		id		path		int	true	"id for get module"	Format(id)
// @Success		200		{object}	response
// @Failure		400,404	{object}	response
// @Failure		500		{object}	response
// @Router		/api/modules/{id} [DELETE]
func (h *ModuleController) ModuleDelete(c *gin.Context) {

	inputID := c.Param("id")

	ID, err := strconv.ParseUint(inputID, 10, 32)
	if inputID == "" || err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result := h.db.Delete(&models.Modules{
		ID: uint32(ID),
	})
	if result.Error != nil {
		newResponse(c, http.StatusInternalServerError, result.Error.Error())
		h.log.Error("failed to delete modules", result.Error.Error())
		return
	} else if result.RowsAffected == 0 {
		newResponse(c, http.StatusBadRequest, "not found")
		h.log.Error("failed to delete modules")
		return
	}

	c.JSON(http.StatusOK, response{"success"})
}

// ModuleItemsCreate
// @Summary			Create an Module item.
// @Description		This API to create an module item.
// @Tags			  ModuleItems
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data	body	    models.CreateModuleItemInput	true	"data body"
// @Success			201		{object}	models.ModuleItems
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/modules-items [POST]
func (h *ModuleController) ModuleItemsCreate(c *gin.Context) {

	var body models.CreateModuleItemInput

	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.db.WithContext(c.Request.Context()).Table("module_items").Create(&body).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update", err.Error())
		return
	}
	module := models.ModuleItems{}
	err = h.db.First(&module, "key=?", body.Key).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update", err.Error())
		return
	}

	c.JSON(http.StatusCreated, module)
}

// ModuleItemUpdate
// @Summary		Update a module item.
// @Description	This API to update an module item.
// @Tags			  ModuleItems
// @Security		BearerAuth
// @Accept			json
// @Produce		json
// @Param			id		path		int		true	"id for update Module item"	Format(id)
// @Param			data	body		models.UpdateModuleItemInput	true	"data body"
// @Success		200		{object}	models.ModuleItems
// @Failure		400,409	{object}	response
// @Failure		500		{object}	response
// @Router			/api/modules-items/{id} [PUT]
func (h *ModuleController) ModuleItemUpdate(c *gin.Context) {

	var body models.UpdateModuleItemInput

	inputID := c.Param("id")

	ID, err := strconv.ParseUint(inputID, 10, 32)
	if inputID == "" || err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Errorf("failed to  bind body: %v", err)
		return
	}
	id := uint32(ID)
	moduleItem := models.ModuleItems{
		ID: id,
	}

	columns := map[string]interface{}{
		"name":        body.Name,
		"description": body.Description,
		"module_id":   body.ModuleID,
		"key":         body.Key,
		"end_point":   body.EndPoint,
		"method":      body.Method,
	}

	err = h.db.WithContext(c.Request.Context()).Clauses(clause.Returning{}).
		Preload("Module").Model(&moduleItem).Where("id=?", id).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update", err.Error())
		return
	}

	c.JSON(http.StatusOK, moduleItem)
}

// ModulesItemsGet
//
//	@Summary		Get all Modules items
//	@Description	This api for get Modules items
//	@Tags			  ModuleItems
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param 			moduleId query uint false "Module ID"
//	@Success		200	{object}	[]models.ModuleItems
//	@Failure		500	{object}	response
//	@Router			/api/modules-items [GET]
func (h *ModuleController) ModulesItemsGet(c *gin.Context) {
	var filter models.GetModuleItemFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update", err.Error())
		h.log.Errorf("failed to bind JSON: %v", err)
		return
	}
	var modules []models.ModuleItems

	db := h.db.Model(models.ModuleItems{}).Preload("Module")
	if filter.ModuleID != nil {
		db = db.Where("module_id", filter.ModuleID)
	}
	err := db.Find(&modules).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Errorf("failed to	find module items %v", err.Error())
		return
	}

	c.JSON(http.StatusOK, modules)
}

// ModuleItemGetByID
//
//	@Summary		Get module by id.
//	@Description	This API to get module item by id.
//	@Tags			  ModuleItems
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"id for get module"	Format(id)
//	@Success		200		{object}	models.ModuleItems
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Router			/api/modules-items/{id} [GET]
func (h *ModuleController) ModuleItemGetByID(c *gin.Context) {
	var moduleItems models.ModuleItems

	inputID := c.Param("id")

	ID, err := strconv.ParseUint(inputID, 10, 32)
	if inputID == "" || err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.db.Model(models.ModuleItems{}).Preload("Module").First(&moduleItems, "id = ?", ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Errorf("failed to get module items: %v", err)
		return
	}

	c.JSON(http.StatusOK, moduleItems)
}

// ModuleItemDelete
//
//	@Summary		Delete module item by id.
//	@Description	This API to delete module item by id.
//	@Tags			  ModuleItems
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"id for delete module item"	Format(id)
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Router			/api/modules-items/{id} [DELETE]
func (h *ModuleController) ModuleItemDelete(c *gin.Context) {
	inputID := c.Param("id")

	ID, err := strconv.ParseUint(inputID, 10, 32)
	if inputID == "" || err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id := uint32(ID)
	result := h.db.WithContext(c.Request.Context()).Delete(&models.ModuleItems{
		ID: id,
	})
	if result.Error != nil {
		newResponse(c, http.StatusInternalServerError, result.Error.Error())
		h.log.Error("failed to update", result.Error.Error())
		return
	} else if result.RowsAffected == 0 {
		newResponse(c, http.StatusBadRequest, "not found")
		h.log.Error("failed to delete module item")
		return
	}

	c.JSON(http.StatusOK, response{"success"})
}
