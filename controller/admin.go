package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/Asliddin3/energy-maximum/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {
	*Handler
}

func (h *Handler) NewAdminController(api *gin.RouterGroup) {
	adminContr := &AdminController{h}
	admin := api.Group("admin")
	{
		admin.POST("", h.DeserializeAdmin(), adminContr.CreateAdmin)
		admin.POST("/auth", adminContr.AuthAdmin)
		admin.PUT("/:id", h.DeserializeAdmin(), adminContr.UpdateAdmin)
		admin.GET("", h.DeserializeAdmin(), adminContr.GetAdmins)
		admin.PUT("/me", h.DeserializeAdmin(), adminContr.UpdateMe)
		admin.GET("/me", h.DeserializeAdmin(), adminContr.GetMe)
		admin.GET("/:id", h.DeserializeAdmin(), adminContr.GetAdminById)
		admin.PUT("/activate/:id", h.DeserializeAdmin(), adminContr.ActivateAdmin)
		admin.PUT("/deactivate/:id", h.DeserializeAdmin(), adminContr.DeactivateAdmin)
		admin.DELETE("/:id", h.DeserializeAdmin(), adminContr.DeleteAdmin)
	}
}

// @Summary		  Create admin
// @Description	   this api is for create admin
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	body		models.AdminsCreateRequest	true	"data body"
// @Success			201		{object}	models.TokenResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin [POST]
func (h *AdminController) CreateAdmin(c *gin.Context) {
	var body models.AdminsCreateRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	hashed, err := h.hash.HashPassword(body.Password)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "failed to hash password")
		h.log.Error("error while hash password", logger.Error(err))
		return
	}
	admin := &models.Admins{
		Username:  body.Username,
		Password:  hashed,
		RoleID:    body.RoleID,
		CreatedAt: timeNow(),
		IsActive:  body.IsActive,
		// CreatedID: ,
	}
	err = h.db.Create(&admin).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique ") {
			newResponse(c, http.StatusBadRequest, "already exists username")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, admin)
}

// // @Summary		  Refresh access token
// // @Description	   this api is for refresh access token
// // @Tags			Admin
// // @Security		BearerAuth
// // @Accept			json
// // @Produce			json
// // @Param           refreshToken   query  string   true  "refresh token"
// // @Success			201		{object}	models.AccessToken
// // @Failure			400,409	{object}	response
// // @Failure			500		{object}	response
// // @Router			/api/admin/refresh [POST]
// func (h *AdminController) RefreshToken(c *gin.Context) {
// 	refreshToken := c.Query("refreshToken")
// 	sub, err := utils.ValidateToken(refreshToken, h.cfg.RefreshTokenPublicKey)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
// 		return
// 	}
// 	admin := &models.Admins{}
// 	err = h.db.First(admin, "id=? AND is_active=true AND deleted_at IS NUll", sub).Error
// 	if err != nil {
// 		newResponse(c, http.StatusInternalServerError, "failed to find customer")
// 		return
// 	}

// 	token, err := utils.CreateToken(time.Duration(time.Hour), admin.ID, h.cfg.AccessTokenPrivateKey)
// 	if err != nil {
// 		newResponse(c, http.StatusInternalServerError, err.Error())
// 		h.log.Error("error while token", logger.Error(err))
// 		return
// 	}
// 	c.JSON(http.StatusOK, models.AccessToken{AccessToken: token})
// }

// @Summary		  auth admin
// @Description	   this api is for auth admin
// @Tags			Admin
// @Accept			json
// @Produce			json
// @Param			data 	body		models.AdminAuth	true	"data body"
// @Success			201		{object}	models.TokenResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin/auth [POST]
func (h *AdminController) AuthAdmin(c *gin.Context) {
	var body models.AdminAuth
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to bind json", err.Error())
		return
	}

	admin := &models.Admins{}
	err = h.db.First(admin, "username=? AND is_active=true AND deleted_at IS NUll", body.Username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "wrong username or password")
			return
		}
		newResponse(c, http.StatusInternalServerError, "failed to find admin")
		return
	}
	err = h.hash.CheckPassword(admin.Password, body.Password)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "wrong password or username")
		return
	}
	token, err := utils.CreateToken(time.Duration(time.Hour*720), admin.ID, h.cfg.AccessTokenPrivateKey)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("error while token", logger.Error(err))
		return
	}
	var moduleItemKeys []string
	err = h.db.Model(&models.RoleItems{}).Where("role_id=?", admin.RoleID).Find(&moduleItemKeys).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// refreshToken, err := utils.CreateToken(time.Duration(time.Hour*720), admin.ID, h.cfg.RefreshTokenPrivateKey)
	// if err != nil {
	// 	newResponse(c, http.StatusInternalServerError, err.Error())
	// 	h.log.Error("error while token", logger.Error(err))
	// 	return
	// }

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:    token,
		ModuleItemKeys: moduleItemKeys,
	})
}

// @Summary		  Update admin for superuser
// @Description	   this api is for Update admin for superuser
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path    int    true    "admin  id"
// @Param			data 	query		models.AdminsCreateRequest	true	"data body"
// @Success			201		{object}	models.TokenResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin/{id} [PUT]
func (h *AdminController) UpdateAdmin(c *gin.Context) {
	// user := h.GetAdmin(c)
	// if !user.IsSuperuser {
	// 	newResponse(c, http.StatusForbidden, "yor are not superuser")
	// 	return
	// }
	inputId := c.Param("id")
	id, err := strconv.Atoi(inputId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "Invalid id")
		return
	}
	var body models.AdminsCreateRequest
	err = c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	hashed, err := h.hash.HashPassword(body.Password)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "failed to hash password")
		h.log.Error("error while hash password", logger.Error(err))
		return
	}
	admin := &models.Admins{}
	err = h.db.First(&admin, "id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get admin")
		return
	}
	admin.Username = body.Username
	admin.Password = hashed
	admin.RoleID = body.RoleID
	if body.IsActive != nil {
		admin.IsActive = body.IsActive
	}
	admin.UpdatedAt = timeNow()
	err = h.db.Save(&admin).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique ") {
			newResponse(c, http.StatusBadRequest, "already exists username")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, admin)
}

// @Summary		  Update admin
// @Description	   this api is for Update admin
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	query		models.AdminsRequest	true	"data body"
// @Success			201		{object}	models.TokenResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin/me [PUT]
func (h *AdminController) UpdateMe(c *gin.Context) {
	user := h.GetAdmin(c)
	var body models.AdminsRequest
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	hashed, err := h.hash.HashPassword(body.Password)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "failed to hash password")
		h.log.Error("error while hash password", logger.Error(err))
		return
	}
	admin := &models.Admins{}
	err = h.db.First(&admin, "id=?", user.Id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get admin")
		return
	}
	admin.Username = body.Username
	admin.Password = hashed
	admin.UpdatedAt = timeNow()
	err = h.db.Save(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			newResponse(c, http.StatusBadRequest, "already exists username")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, admin)
}

// @Summary		  Get admins
// @Description	   this api is for Get admins
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			201		{object}	[]models.Admins
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin [GET]
func (h *AdminController) GetAdmins(c *gin.Context) {
	var admins []models.Admins
	err := h.db.Find(&admins).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get admin", err.Error())
		return
	}
	c.JSON(http.StatusOK, admins)
}

// @Summary		  Deactivate admin
// @Description	   this api is for deactivating admin
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path   int    true  "admin id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin/deactivate/{id} [PUT]
func (h *AdminController) DeactivateAdmin(c *gin.Context) {
	// currentUser := h.GetAdmin(c)
	// if !currentUser.IsSuperuser {
	// 	newResponse(c, http.StatusForbidden, "your are not superuser")
	// 	return
	// }
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "missing id")
		return
	}
	err := h.db.Model(&models.Admins{}).Where("id = ?", id).UpdateColumn("is_active", false).Error
	if err != nil {
		h.log.Error("failed to update admins", err.Error())
		newResponse(c, http.StatusInternalServerError, "failed to update admin status")
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Activate admin
// @Description	   this api is for activating admin
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path   int    true  "admin id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin/activate/{id} [PUT]
func (h *AdminController) ActivateAdmin(c *gin.Context) {
	// currentUser := h.GetAdmin(c)
	// if !currentUser.IsSuperuser {
	// 	newResponse(c, http.StatusForbidden, "your are not superuser")
	// 	return
	// }
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "missing id")
		return
	}
	err := h.db.Model(&models.Admins{}).Where("id = ?", id).UpdateColumn("is_active", true).Error
	if err != nil {
		h.log.Error("failed to update admins", err.Error())
		newResponse(c, http.StatusInternalServerError, "failed to update admin status")
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Delete admin
// @Description	   this api is for DELETE admin
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path   int    true  "admin id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin/{id} [DELETE]
func (h *AdminController) DeleteAdmin(c *gin.Context) {
	// currentUser := h.GetAdmin(c)
	// if !currentUser.IsSuperuser {
	// 	newResponse(c, http.StatusForbidden, "your are not superuser")
	// 	return
	// }
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "missing id")
		return
	}
	columns := map[string]interface{}{
		"deleted_at": timeNow(),
		"is_active":  false,
	}
	err := h.db.Model(&models.Admins{}).Where("id = ?", id).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to update admin status")
		return
	}

	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Get me
// @Description	   this api is for Get me
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			200		{object}	models.AdminResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin/me [GET]
func (h *AdminController) GetMe(c *gin.Context) {
	admin := h.GetAdmin(c)
	var admins models.Admins
	err := h.db.First(&admins, "id=?", admin.Id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get admin", err.Error())
		return
	}
	var moduleItemKeys []string
	err = h.db.Model(&models.RoleItems{}).Where("role_id=?", admins.RoleID).Find(&moduleItemKeys).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.AdminResponse{
		Admins:         admins,
		ModuleItemKeys: moduleItemKeys,
	})
}

// @Summary		  Get admin  by id from superuser
// @Description	   this api is for Get admin by id from superuser
// @Tags			Admin
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path  int  true  "admin id"
// @Success			200		{object}	models.Admins
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/admin/{id} [GET]
func (h *AdminController) GetAdminById(c *gin.Context) {
	// admin := h.GetAdmin(c)
	id := c.Param("id")
	var admins models.Admins
	// if admin.IsSuperuser {
	err := h.db.First(&admins, "id=?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusNotFound, "no such admin")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get admin", err.Error())
		return
	}
	// } else {
	// 	newResponse(c, http.StatusForbidden, "you are not superuser")
	// 	return
	// }

	c.JSON(http.StatusOK, admins)
}
