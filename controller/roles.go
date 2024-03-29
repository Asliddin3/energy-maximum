package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RolesController struct {
	*Handler
}

func (h *Handler) NewRolesController(api *gin.RouterGroup) {
	role := &RolesController{
		Handler: h,
	}
	custom := api.Group("roles", h.DeserializeAdmin())
	{
		custom.GET("", role.GetRoles)
		custom.PUT("/:id", role.UpdateRole)
		custom.GET("/:id", role.GetRoleById)
		custom.POST("", role.CreateRole)
		custom.DELETE("/:id", role.DeleteRole)
		// custom.POST("/refresh", customer.RefreshToken)
	}
	roleItem := api.Group("role-items", h.DeserializeAdmin())
	{
		roleItem.GET("", role.GetRoleItems)
		roleItem.PUT("/list", role.RoleModuleItemsUpdate)
		roleItem.PUT("/:id", role.UpdateRoleItem)
		roleItem.GET("/:id", role.GetRoleItemById)
		roleItem.POST("", role.CreateRoleItem)
		roleItem.DELETE("/:id", role.DeleteRoleItem)
	}
}

// @Summary		  Update roles
// @Description	   this api is for Update roles
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	body		models.RolesRequest	true	"data body"
// @Success			201		{object}	models.Roles
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/roles/{id}  [PUT]
func (h *RolesController) UpdateRole(c *gin.Context) {
	currentUser := h.GetAdmin(c)
	var body models.RolesRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create role", err.Error())
		return
	}
	var customer models.Roles
	id := c.Param("id")
	err = h.db.First(&customer, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get customer", err.Error())
		return
	}
	columns := map[string]interface{}{
		"title":      body.Title,
		"comment":    body.Comment,
		"key":        body.Key,
		"is_active":  body.IsActive,
		"updated_at": timeNow(),
		"updated_id": currentUser.Id,
	}
	customer.UpdatedAt = timeNow()
	err = h.db.Clauses(clause.Returning{}).Model(&customer).UpdateColumns(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary		  Create role
// @Description	   this api is for Create role
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	body		models.RolesRequest	true	"data body"
// @Success			201		{object}	models.Roles
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/roles [POST]
func (h *RolesController) CreateRole(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.RolesRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create role", err.Error())
		return
	}

	customer := models.Roles{
		Title:     body.Title,
		Key:       body.Key,
		Comment:   body.Comment,
		IsActive:  body.IsActive,
		CreatedAt: timeNow(),
		CreatedID: &admin.Id,
	}
	err = h.db.Clauses(clause.Returning{}).Create(&customer).Error
	if err != nil {

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary		  Get by id
// @Description	   this api is for Get by id
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path  int  true  "role id"
// @Success			200		{object}	models.RoleWithModuleItems
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/roles/{id} [GET]
func (h *RolesController) GetRoleById(c *gin.Context) {
	id := c.Param("id")
	var roles models.Roles
	err := h.db.First(&roles, "id=?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "no such role")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get role", err.Error())
		return
	}
	var modulesKeys []string
	err = h.db.Model(&models.RoleItems{}).Select("module_item_key").Find(&modulesKeys).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.RoleWithModuleItems{
		Roles:          roles,
		ModuleItemKeys: modulesKeys,
	})
}

// RoleModuleItemsUpdate
//
//	@Summary		Update a Role module items
//	@Description	This api is for update Role module item it will delete all existing role module items and will create news
//	@Tags			Roles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		models.UpdateRoleItemsList	true	"data body"
//	@Success		200		{object}	[]models.RoleItems
//	@Failure		400,409	{object}	response
//	@Failure		500		{object}	response
//	@Router			/api/role-items/list [PUT]
func (h *RolesController) RoleModuleItemsUpdate(c *gin.Context) {
	currentUser := h.GetAdmin(c)

	var body models.UpdateRoleItemsList

	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	tx := h.db.Begin()

	err = tx.Delete(&models.RoleItems{}, "role_id=?", body.RoleID).Error
	if err != nil {
		tx.Rollback()
		h.log.Error("failed to delete role items", err.Error())
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	roleModuleItems := make([]models.RoleItems, 0)
	timeNow := time.Now()
	for _, key := range body.ModuleItemKeys {
		roleModuleItems = append(roleModuleItems,
			models.RoleItems{
				RoleID:        int(body.RoleID),
				ModuleItemKey: key,
				CreatedAt:     &timeNow,
				CreatedID:     &currentUser.Id,
			})
	}
	err = tx.Create(&roleModuleItems).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrEmptySlice) {
			tx.Rollback()

			newResponse(c, http.StatusInternalServerError, err.Error())
			h.log.Errorf("failed to create role module items %v", err)
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, roleModuleItems)
}

// @Summary		    Get roles
// @Description	    this api is to get roles
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           data    query   models.RoleFilter  true "filter"
// @Success			201		{object}	[]models.Roles
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/roles  [GET]
func (h *RolesController) GetRoles(c *gin.Context) {
	var body models.RoleFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var countries []models.Roles
	db := h.db.Model(&models.Roles{})
	if body.Name != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", body.Name))
	}
	// if !body.All {
	// 	db = db.Where("is_active=true")
	// }
	if !body.WithDelete {
		db = db.Where("is_deleted=false")
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}

	err = db.Limit(body.PageSize).Offset((body.Page - 1) * body.PageSize).Find(&countries).Error
	if err != nil {
		h.log.Error("failed to get roles", err.Error())
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, countries)
}

// @Summary		  Delete Roles
// @Description	   this api is to delete Roles
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/roles/{id} [DELETE]
func (h *RolesController) DeleteRole(c *gin.Context) {
	admin := h.GetAdmin(c)
	id := c.Param("id")
	columns := map[string]interface{}{
		"is_deleted": true,
		"deleted_at": timeNow(),
		"deleted_id": admin.Id,
	}
	err := h.db.Model(&models.Roles{}).Where("id=?", id).Updates(columns).Error
	if err != nil {
		h.log.Error("failed to delete role", err.Error())
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Update roles
// @Description	   this api is for Update roles
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			id 		path     int  true "id"
// @Param			data 	body		models.RoleItemRequest	true	"data body"
// @Success			201		{object}	models.RoleItems
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/role-items/{id} [PUT]
func (h *RolesController) UpdateRoleItem(c *gin.Context) {
	admin := h.GetAdmin(c)
	id := c.Param("id")
	var body models.RoleItemRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create role", err.Error())
		return
	}
	var customer models.RoleItems
	columns := map[string]interface{}{
		"role_id":         body.RoleID,
		"module_item_key": body.Key,
		"updated_id":      admin.Id,
		"updated_at":      timeNow(),
	}
	err = h.db.Clauses(clause.Returning{}).Model(&customer).Where("id=?", id).UpdateColumns(columns).Error
	if err != nil {
		h.log.Error("failed to update role item", err.Error())
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary		  Create role
// @Description	   this api is for Create role
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	body		models.RoleItemRequest	true	"data body"
// @Success			201		{object}	models.RoleItems
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/role-items [POST]
func (h *RolesController) CreateRoleItem(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.RoleItemRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create role", err.Error())
		return
	}

	customer := models.RoleItems{
		RoleID:        body.RoleID,
		ModuleItemKey: body.Key,
		CreatedAt:     timeNow(),
		CreatedID:     &admin.Id,
	}
	err = h.db.Clauses(clause.Returning{}).Create(&customer).Error
	if err != nil {
		h.log.Error("failed to create role item", err.Error())
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary		  Get by id
// @Description	   this api is for Get by id
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path  int  true  "role id"
// @Success			200		{object}	models.RoleItems
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/role-items/{id} [GET]
func (h *RolesController) GetRoleItemById(c *gin.Context) {
	id := c.Param("id")
	var admins models.Roles
	err := h.db.First(&admins, "id=?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "no such role")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get role", err.Error())
		return
	}
	c.JSON(http.StatusOK, admins)
}

// @Summary		    Get roles
// @Description	    this api is to get roles
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           data    query   models.RoleItemFilter  true "filter"
// @Success			201		{object}	[]models.RoleItems
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/role-items  [GET]
func (h *RolesController) GetRoleItems(c *gin.Context) {
	var body models.RoleItemFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var countries []models.RoleItems
	db := h.db.Model(&models.RoleItems{})
	if body.Name != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", body.Name))
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	err = db.Limit(body.PageSize).Offset((body.Page - 1) * body.PageSize).Find(&countries).Error
	if err != nil {
		h.log.Error("failed to get role items", err.Error())
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, countries)
}

// @Summary		  Delete Roles
// @Description	   this api is to delete Roles
// @Tags			Roles
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/role-items/{id} [DELETE]
func (h *RolesController) DeleteRoleItem(c *gin.Context) {
	id := c.Param("id")
	columns := map[string]interface{}{
		"is_deleted": true,
	}
	err := h.db.Model(&models.RoleItems{}).Where("id=?", id).Updates(columns).Error
	if err != nil {
		h.log.Error("failed to delete role item", err.Error())
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}
