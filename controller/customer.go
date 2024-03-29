package controller

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/Asliddin3/energy-maximum/pkg/sms"
	"github.com/Asliddin3/energy-maximum/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerController struct {
	*Handler
	Sms             *sms.Sms
	DEVELOPER_PHONE string
}

func (h *Handler) NewCustomerController(api *gin.RouterGroup, sms *sms.Sms) {
	customer := &CustomerController{
		Handler:         h,
		Sms:             sms,
		DEVELOPER_PHONE: "998995117361",
	}
	custom := api.Group("customer")
	{
		custom.POST("/register", h.DeserializeCustomer(), customer.Register)
		custom.POST("/login", customer.Login)
		custom.POST("/check-code", customer.checkCode)
		custom.POST("/send-code", customer.SendCode)
		custom.GET("", h.DeserializeAdmin(), customer.GetCustomers)
		custom.PUT("", h.DeserializeCustomer(), customer.UpdateCustomer)
		custom.GET("/me", h.DeserializeCustomer(), customer.GetMe)
		custom.GET("/:id", h.DeserializeAdmin(), customer.GetCustomerById)
		custom.POST("", h.DeserializeAdmin(), customer.CreateCustomerByAdmin)
		// custom.POST("/refresh", customer.RefreshToken)
		custom.PUT("/password", h.DeserializeCustomer(), customer.UpdatePassword)
	}
}

// @Summary		  Send code to customer
// @Description	   this api is for sending code to a customer
// @Tags			Customer
// @Accept			json
// @Produce			json
// @Param			data 	body		models.CustomerCode	true	"data body"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer/send-code [POST]
func (h *CustomerController) SendCode(c *gin.Context) {
	var body models.CustomerCode
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to bind body", err.Error())
		return
	}
	phone := formatPhone(body.Phone)
	if phone == "" {
		newResponse(c, http.StatusUnauthorized, "Неправильный номер телефона!")
		return
	}
	e := time.Now()
	number, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, "Неправильный номер телефона!")
		h.log.Error("failed to parse phone", logger.Error(err))
		return
	}
	code, err := h.GetCodes(int(number))
	if err != nil && err != gorm.ErrRecordNotFound {
		newResponse(c, http.StatusInternalServerError, "failed to get codes")
		h.log.Error("Error while find codes", logger.Error(err))
		return
	}
	if err != gorm.ErrRecordNotFound {
		s := code.CreatedAt
		second := int(e.Unix()) - int(s.Unix())
		if second < 60 {
			newResponse(c, http.StatusUnauthorized,
				fmt.Sprintf("Попробуйте еще раз через %d секунд", 60-second))
			return
		}
	}
	sendCode := genRandNum(99999, 999999)
	if phone == h.DEVELOPER_PHONE {
		sendCode = 997361
	}
	err = h.Sms.SendCode(phone, fmt.Sprintf("Your code from pribor -%d ", sendCode))
	if err != nil {
		// newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("Error while send code", logger.Error(err))
		// return
	}
	fmt.Println("message sended")

	err = h.db.Create(&models.Codes{
		Phone:     int(number),
		Code:      int(sendCode),
		CreatedAt: timeNow(),
	}).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("Error while create codes", logger.Error(err))
		return
	}

	newResponse(c, http.StatusOK, "success")
}
func genRandNum(min, max int64) int64 {
	bg := big.NewInt(max - min)
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}
	return n.Int64() + min
}
func formatPhone(phone string) string {
	phone = strings.ReplaceAll(phone, "+", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, " ", "")
	if len([]rune(phone)) != 12 {
		return ""
	}
	return phone
}
func (h *CustomerController) GetCodes(phone int) (*models.Codes, error) {
	var code models.Codes
	err := h.db.Model(&models.Codes{}).Order("created_at DESC").Limit(1).
		First(&code, "phone=?", phone).Error
	return &code, err
}

// @Summary		  Login customer
// @Description	   this api is for login customer
// @Tags			Customer
// @Accept			json
// @Produce			json
// @Param			data 	body		models.CustomerLogin	true	"data body"
// @Success         200     {object}    models.TokenResponse
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer/login [POST]
func (h *CustomerController) Login(c *gin.Context) {
	var body models.CustomerLogin
	err := c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	phone := formatPhone(body.Phone)
	if phone == "" {
		newResponse(c, http.StatusUnauthorized, "Неправильный номер телефона!")
		return
	}
	var customer models.Customer
	err = h.db.First(&customer, "phone=?", phone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found user")
			return
		}
		newResponse(c, http.StatusInternalServerError, "failed to get customer")
		h.log.Error("failed to get customer", err.Error())
		return
	}
	err = h.hash.CheckPassword(customer.Password, body.Password)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "wrong password or username")
		return
	}
	key := fmt.Sprintf("customer:%d", customer.ID)
	token, err := utils.CreateToken(time.Duration(time.Hour*720), key, h.cfg.AccessTokenPrivateKey)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("error while token", logger.Error(err))
		return
	}
	// refreshToken, err := utils.CreateToken(time.Duration(time.Hour*720), customer.ID, h.cfg.RefreshTokenPrivateKey)
	// if err != nil {
	// 	newResponse(c, http.StatusInternalServerError, err.Error())
	// 	h.log.Error("error while token", logger.Error(err))
	// 	return
	// }

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken: token,
	})
}

// @Summary		  Register a new customer
// @Description	   this api is for register new customer
// @Tags			Customer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	body		models.CustomerRegister	true	"data body"
// @Success			200		{object}	 models.Customer
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer/register [POST]
func (h *CustomerController) Register(c *gin.Context) {
	customer := h.GetCustomer(c)
	var body models.CustomerRegister
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "failed to bind JSON")
		h.log.Error("Error while bind to json AccountCreate", logger.Error(err))
		return
	}
	if body.Password == "" {
		newResponse(c, http.StatusUnauthorized, "Empty password")
		return
	}

	hashed, err := h.hash.HashPassword(body.Password)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "failed to hash password")
		h.log.Error("error while hash password", logger.Error(err))
		return
	}
	columns := map[string]interface{}{
		"name":     body.Name,
		"password": hashed,
	}
	custom := models.Customer{
		ID: customer.Id,
	}
	err = h.db.Clauses(clause.Returning{}).Model(&custom).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update customer", logger.Error(err))
		return
	}
	// var customer

	c.JSON(http.StatusOK, custom)
}

// @Summary		  Update customer password
// @Description	   this api is for update customer password
// @Tags			Customer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			password  query	string 	true	"customer password"
// @Success			200		{object}	 models.TokenResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer/password [PUT]
func (h *CustomerController) UpdatePassword(c *gin.Context) {
	customer := h.GetCustomer(c)
	password := c.Param("password")
	if password != "" {
		newResponse(c, http.StatusBadRequest, "empty password")
		return
	}
	hashed, err := h.hash.HashPassword(password)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "failed to hash password")
		h.log.Error("error while hash password", logger.Error(err))
		return
	}
	err = h.db.Model(&models.Customer{
		ID: customer.Id,
	}).UpdateColumn("password", hashed).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to update password", logger.Error(err))
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  	Check customer code
// @Description	   	this api is for checking customer code
// @Tags			Customer
// @Accept			json
// @Produce			json
// @Param			data 	body		models.CheckCodeRequest	true	"data body"
// @Success			200		{object}	 models.TokenResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer/check-code [POST]
func (h *CustomerController) checkCode(c *gin.Context) {
	var body models.CheckCodeRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	phone := formatPhone(body.Phone)
	if phone == "" {
		newResponse(c, http.StatusUnauthorized, "Неправильный номер телефона!")
		return
	}
	number, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, "Неправильный номер телефона!")
		h.log.Error("failed to parse phone", logger.Error(err))
		return
	}
	code, err := h.GetCodes(int(number))
	if err != nil && err != gorm.ErrRecordNotFound {
		newResponse(c, http.StatusInternalServerError, "failed to get codes")
		h.log.Error("Error while find codes", logger.Error(err))
		return
	}
	codeInt, err := strconv.Atoi(body.Code)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, "Неправильный код!")
		return
	}
	if code.Code != codeInt {
		newResponse(c, http.StatusUnauthorized, "Неправильный код!")
		return
	}
	var customer models.Customer
	if body.IsRegistered {
		err = h.db.First(&customer, "phone=?", phone).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newResponse(c, http.StatusNotFound, "not found customer by this phone number")
				return
			}
			newResponse(c, http.StatusInternalServerError, err.Error())
			h.log.Error("failed to find customer", err.Error())
			return
		}
	} else {
		customer = models.Customer{
			Phone:     phone,
			CreatedAt: timeNow(),
			LastVisit: timeNow(),
		}
		err = h.db.Create(&customer).Error
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				newResponse(c, http.StatusBadRequest, "already exists customer")
				return
			}
			newResponse(c, http.StatusInternalServerError, "failed to create customer")
			h.log.Error("failed to create customer", err.Error())
			return
		}
	}
	token, err := utils.CreateToken(time.Duration(time.Hour*720), customer.ID, h.cfg.AccessTokenPrivateKey)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("error while token", logger.Error(err))
		return
	}
	// refreshToken, err := utils.CreateToken(time.Duration(time.Hour*720), customer.ID, h.cfg.RefreshTokenPrivateKey)
	// if err != nil {
	// 	newResponse(c, http.StatusInternalServerError, err.Error())
	// 	h.log.Error("error while token", logger.Error(err))
	// 	return
	// }
	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken: token,
	})
}

// // @Summary		  Refresh access token
// // @Description	   this api is for refresh access token
// // @Tags			Customer
// // @Security		BearerAuth
// // @Accept			json
// // @Produce			json
// // @Param           refreshToken   query  string   true  "refresh token"
// // @Success			201		{object}	models.AccessToken
// // @Failure			400,409	{object}	response
// // @Failure			500		{object}	response
// // @Router			/api/customer/refresh [POST]
// func (h *CustomerController) RefreshToken(c *gin.Context) {
// 	// var refreshToken string
// 	refreshToken := c.Query("refreshToken")
// 	sub, err := utils.ValidateToken(refreshToken, h.cfg.RefreshTokenPublicKey)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
// 		return
// 	}
// 	customer := &models.Customer{}
// 	err = h.db.First(customer, "id=?", sub).Error
// 	if err != nil {
// 		newResponse(c, http.StatusInternalServerError, "failed to find customer")
// 		return
// 	}

// 	token, err := utils.CreateToken(time.Duration(time.Hour), customer.ID, h.cfg.AccessTokenPrivateKey)
// 	if err != nil {
// 		newResponse(c, http.StatusInternalServerError, err.Error())
// 		h.log.Error("error while token", logger.Error(err))
// 		return
// 	}
// 	c.JSON(http.StatusOK, models.AccessToken{AccessToken: token})
// }

// @Summary		  Update customer
// @Description	   this api is for Update customer
// @Tags			Customer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	body		models.CustomerRequest	true	"data body"
// @Success			201		{object}	models.Customer
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer [PUT]
func (h *CustomerController) UpdateCustomer(c *gin.Context) {
	currentUser := h.GetCustomer(c)
	var body models.CustomerRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	var customer models.Customer
	err = h.db.First(&customer, "id=?", currentUser.Id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get customer", err.Error())
		return
	}
	if body.Name != "" {
		customer.Name = body.Name
	}
	if body.Email != "" {
		customer.Email = body.Email
	}
	if body.Birthday != "" {
		customer.Birthday = body.Birthday
	}

	customer.UpdatedAt = timeNow()
	err = h.db.Save(&customer).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary		  Create customer
// @Description	   this api is for Create customer by the phone if already exists return existing customer
// @Tags			Customer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	body		models.CustomerRegisterByAdmin	true	"data body"
// @Success			201		{object}	models.Customer
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer [POST]
func (h *CustomerController) CreateCustomerByAdmin(c *gin.Context) {
	admin := h.GetAdmin(c)
	var body models.CustomerRegisterByAdmin
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}
	phone := formatPhone(body.Phone)
	if phone == "" {
		newResponse(c, http.StatusUnauthorized, "Неправильный номер телефона!")
		return
	}

	customer := models.Customer{
		Phone:     phone,
		Name:      body.Name,
		CreatedAt: timeNow(),
		CreatedID: &admin.Id,
	}
	err = h.db.Clauses(clause.Returning{}).Create(&customer).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			err = h.db.First(&customer, "phone=?", phone).Error
			if err != nil {
				newResponse(c, http.StatusInternalServerError, "failed to get customer by this phone")
				return
			}
			c.JSON(http.StatusOK, customer)
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary		  Get me
// @Description	   this api is for Get me
// @Tags			Customer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			200		{object}	models.Customer
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer/me [GET]
func (h *CustomerController) GetMe(c *gin.Context) {
	admin := h.GetCustomer(c)
	var admins models.Customer
	err := h.db.First(&admins, "id=?", admin.Id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get customer", err.Error())
		return
	}
	c.JSON(http.StatusOK, admins)
}

// @Summary		  Get by id
// @Description	   this api is for Get by id
// @Tags			Customer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path  int  true  "customer id"
// @Success			200		{object}	models.Customer
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer/{id} [GET]
func (h *CustomerController) GetCustomerById(c *gin.Context) {
	id := c.Param("id")
	var admins models.Customer
	err := h.db.First(&admins, "id=?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "no such customer")
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get customer", err.Error())
		return
	}
	c.JSON(http.StatusOK, admins)
}

// @Summary		  Get customers
// @Description	   this api is for Get me
// @Tags			Customer
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Success			200		{object}	[]models.Customer
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/customer   [GET]
func (h *CustomerController) GetCustomers(c *gin.Context) {
	var customers []models.Customer
	err := h.db.Find(&customers).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error("failed to get customer", err.Error())
		return
	}
	c.JSON(http.StatusOK, customers)
}
