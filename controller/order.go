package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Asliddin3/energy-maximum/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderController struct {
	*Handler
}

func (h *Handler) NewOrderController(api *gin.RouterGroup) {
	order := &OrderController{h}
	orderHandler := api.Group("order", h.DeserializeCustomer())
	{
		orderHandler.POST("/", order.CreateOrder)
		// orderHandler.PUT("/:id", order.UpdateOrder)
		orderHandler.GET("", order.GetOrders)
		orderHandler.GET("/:id", order.GetByID)
	}
	api.POST("/order/applicant", order.CreateOrderApplicant)

	adminHandler := api.Group("order", h.DeserializeAdmin())
	{
		adminHandler.GET("/applicant", order.GetOrdersApplicant)
		adminHandler.POST("/:id", order.CreateOrderByAdmin)
		adminHandler.PUT("/:id", order.UpdateOrderByAdmin)
		adminHandler.GET("/all", order.GetAllOrders)
		adminHandler.PUT("/cancel/:id", order.CancelOrder)
		adminHandler.PUT("/finish/:id", order.FinishOrder)
		adminHandler.DELETE("/:id", order.DeleteOrder)
		adminHandler.DELETE("/applicant/:id", order.DeleteOrderApplicantByID)

	}
}

// @Summary		  Create new order applicant
// @Description	   this api is create new order applicant
// @Tags			Order
// @Accept			json
// @Produce			json
// @Param			data 	body		models.OrderApplicantRequest	false	"data body"
// @Success			201		{object}	models.OrderApplicant
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/applicant [POST]
func (h *OrderController) CreateOrderApplicant(c *gin.Context) {
	var body models.OrderApplicantRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	order := models.OrderApplicant{
		FullName:  body.FullName,
		Phone:     body.Phone,
		Message:   body.Message,
		CreatedAt: timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Create(&order).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, order)
}

// @Summary		  Get customer orders applicants
// @Description	   this api is to get customer orders applicants
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			filter   query   models.OrderApplicantFilter  true "filter"
// @Success			201		{object}	[]models.CustomOrderResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/applicant  [GET]
func (h *OrderController) GetOrdersApplicant(c *gin.Context) {
	var body models.OrderApplicantFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var orders []models.OrderApplicant
	db := h.db.Debug().Model(&models.OrderApplicant{})

	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	if body.Phone != "" {
		db = db.Where("phone LIKE ?", fmt.Sprintf("%%%s%%", body.Phone))
	}

	if body.FullName != "" {
		db = db.Where("LOWER(full_name) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", body.FullName))
	}
	var count int
	err = db.Model(&models.OrderApplicant{}).Select("COUNT(*)").Scan(&count).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get orders status")
		h.log.Error("failed to get orders status", err.Error())
		return
	}

	err = db.Debug().Model(&models.OrderApplicant{}).Select("*").Limit(body.PageSize).Offset((body.Page - 1) * body.PageSize).Find(&orders).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res := models.OrderApplicantResponse{
		Page:      body.Page,
		PageSize:  body.PageSize,
		Applicant: orders,
		Count:     count,
	}
	c.JSON(http.StatusOK, res)
}

// @Summary		  delete by id order applicant
// @Description	   this api is to delete by id order applicant
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	response
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/applicant/{id} [DELETE]
func (h *OrderController) DeleteOrderApplicantByID(c *gin.Context) {
	orderId := c.Param("id")
	var order models.OrderApplicant
	err := h.db.Delete(&order, "id=?", orderId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found order")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, response{"success"})
}

// @Summary		  Create new order
// @Description	   this api is create new order
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			data 	body		models.OrderRequest	false	"data body"
// @Success			201		{object}	models.OrderResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order [POST]
func (h *OrderController) CreateOrder(c *gin.Context) {
	customer := h.GetCustomer(c)
	var body models.OrderRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	order := models.Orders{
		Description: body.Description,
		CustomerID:  customer.Id,
		Total:       body.Total,
		CreatedAt:   timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Create(&order).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	orderItems := make([]models.OrderItems, len(body.Items))
	for i, item := range body.Items {
		orderItems[i] = models.OrderItems{
			Price:  item.Price,
			Amount: item.Amount,
			ItemId: item.ItemID,
		}
	}
	fmt.Println("orderItems: ", orderItems)
	if len(orderItems) > 0 {
		err = h.db.Clauses(clause.Returning{}).Create(&orderItems).Error
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			h.log.Error("failed to create order", err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, models.OrderResponse{
		Orders: &order,
		Items:  orderItems,
	})
}

// @Summary		  Get customer orders
// @Description	   this api is to get customer orders
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			filter   query   models.OrderFilter  true "filter"
// @Success			201		{object}	[]models.CustomOrderResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order  [GET]
func (h *OrderController) GetOrders(c *gin.Context) {
	customer := h.GetCustomer(c)
	var body models.OrderFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var orders []models.Orders
	db := h.db.Debug().Model(&models.Orders{}).Where("customer_id=?", customer.Id)

	// if body.DateFrom != "" {
	// 	db = db.Where("created_at>=?", body.DateFrom)
	// }
	// if body.DateTo != "" {
	// 	db = db.Where("created_at<=?", body.DateTo)
	// }
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	if body.Status != nil {
		db = db.Where("status=?", body.Status)
	}
	var statusCounts []struct {
		Status int
		Count  int
	}

	err = db.Debug().Model(&models.Orders{}).
		Select("status, COUNT(*) as count").
		Group("status").Where("customer_id=?", customer.Id).
		Scan(&statusCounts).
		Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get orders status")
		h.log.Error("failed to get orders status", err.Error())
		return
	}
	res := models.CustomOrderResponse{
		Page:     body.Page,
		PageSize: body.PageSize,
	}
	for _, status := range statusCounts {
		if status.Status == 0 {
			res.ActiveCount = status.Count
		} else if status.Status == 1 {
			res.FinishedCount = status.Count
		} else if status.Status == 2 {
			res.CancelledCount = status.Count
		}
		res.Count += status.Count
	}
	err = db.Debug().Select("*").Limit(body.PageSize).Offset((body.Page - 1) * body.PageSize).Find(&orders).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res.Orders = orders
	c.JSON(http.StatusOK, res)
}

// @Summary		  Get all orders for admin
// @Description	   this api is to get all orders from admin
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param			filter   query   models.AdminOrderFilter  true "filter"
// @Success			201		{object}	[]models.Orders
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/all  [GET]
func (h *OrderController) GetAllOrders(c *gin.Context) {
	var body models.AdminOrderFilter
	err := c.ShouldBindQuery(&body)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var orders []models.Orders
	db := h.db.Model(&models.Orders{})
	if body.CustomerID != 0 {
		db = db.Where("customer_id=?", body.CustomerID)
	}
	if body.DateFrom != "" {
		db = db.Where("created_at>=?", body.DateFrom)
	}
	if body.DateTo != "" {
		db = db.Where("created_at<=?", body.DateTo)
	}
	if body.Status != nil {
		db = db.Where("status=?", body.Status)
	}
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	err = db.Limit(body.Page).Offset((body.PageSize - 1) * body.PageSize).Find(&orders).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, orders)
}

// @Summary		  Create new order by admin
// @Description	   this api is create new order admin
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path     int   true   "customer id"
// @Param			data 	body		models.OrderRequest	false	"data body"
// @Success			201		{object}	models.OrderResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/{id} [POST]
func (h *OrderController) CreateOrderByAdmin(c *gin.Context) {
	// customer := h.GetCustomer(c)
	inputId := c.Param("id")
	id, err := strconv.ParseInt(inputId, 10, 64)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var body models.OrderRequest
	err = c.ShouldBindJSON(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	order := models.Orders{
		Description: body.Description,
		CustomerID:  int(id),
		Total:       body.Total,
		CreatedAt:   timeNow(),
	}
	err = h.db.Debug().Clauses(clause.Returning{}).Create(&order).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	orderItems := make([]models.OrderItems, len(body.Items))
	for i, item := range body.Items {
		orderItems[i] = models.OrderItems{
			Price:  item.Price,
			Amount: item.Amount,
			ItemId: item.ItemID,
		}
	}
	fmt.Println("orderItems: ", orderItems)
	if len(orderItems) > 0 {
		err = h.db.Clauses(clause.Returning{}).Create(&orderItems).Error
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			h.log.Error("failed to create order", err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, models.OrderResponse{
		Orders: &order,
		Items:  orderItems,
	})
}

// @Summary		  Update order by admin
// @Description	   this api is Update order admin
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id     path     int   true   "order id"
// @Param			data 	formData		models.OrderRequest	false	"data body"
// @Success			201		{object}	models.OrderResponse
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/{id} [PUT]
func (h *OrderController) UpdateOrderByAdmin(c *gin.Context) {
	admin := h.GetAdmin(c)
	inputId := c.Param("id")
	id, err := strconv.ParseInt(inputId, 10, 64)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var body models.OrderUpdateRequest
	err = c.ShouldBind(&body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to create category", err.Error())
		return
	}

	order := models.Orders{
		ID: int(id),
	}
	columns := map[string]interface{}{
		"updated_at": timeNow(),
		"updated_id": admin.Id,
	}
	if body.Description != "" {
		columns["description"] = body.Description
	}
	if body.CustomerId != 0 {
		columns["customer_id"] = body.CustomerId
	}
	if body.Total != nil {
		columns["total"] = body.Total
	}
	tr := h.db.Begin()
	err = tr.Debug().Clauses(clause.Returning{}).Model(&order).
		Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		tr.Rollback()
		return
	}
	err = tr.Delete(&models.OrderItems{}, "order_id=?", id).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		tr.Rollback()
		return
	}
	orderItems := make([]models.OrderItems, len(body.Items))
	for i, item := range body.Items {
		orderItems[i] = models.OrderItems{
			Price:  item.Price,
			Amount: item.Amount,
			ItemId: item.ItemID,
		}
	}
	fmt.Println("orderItems: ", orderItems)
	if len(orderItems) > 0 {
		err = tr.Clauses(clause.Returning{}).Create(&orderItems).Error
		if err != nil {
			tr.Rollback()
			newResponse(c, http.StatusInternalServerError, err.Error())
			h.log.Error("failed to create order", err.Error())
			return
		}
	}
	tr.Commit()
	c.JSON(http.StatusOK, models.OrderResponse{
		Orders: &order,
		Items:  orderItems,
	})
}

// @Summary		  Get by id order
// @Description	   this api is to get by id order
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	[]models.Orders
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/{id} [GET]
func (h *OrderController) GetByID(c *gin.Context) {
	orderId := c.Param("id")
	var order models.Orders
	err := h.db.First(&order, "id=?", orderId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newResponse(c, http.StatusBadRequest, "not found order")
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, order)
}

// @Summary		  Cancel  order
// @Description	   this api is for cancel order this api change status to 3
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id    path     int   true   "id"
// @Success			201		{object}	[]models.Orders
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/cancel/{id} [PUT]
func (h *OrderController) CancelOrder(c *gin.Context) {
	admin := h.GetAdmin(c)
	orderId := c.Param("id")
	var order models.Orders
	columns := map[string]interface{}{
		"status":     3,
		"updated_at": timeNow(),
		"updated_id": admin.Id,
	}
	err := h.db.Clauses(clause.Returning{}).Model(&order).Where("id=?", orderId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error("failed to update order", err.Error())
		return
	}
	c.JSON(http.StatusOK, order)
}

// @Summary		 	Finish Order
// @Description	   	this api is finish order this api change status to 2
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id   	path     int   true   "id"
// @Success			201		{object}	models.Orders
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/finish/{id} [PUT]
func (h *OrderController) FinishOrder(c *gin.Context) {
	admin := h.GetAdmin(c)
	orderId := c.Param("id")
	var order models.Orders
	columns := map[string]interface{}{
		"status":     2,
		"updated_at": timeNow(),
		"updated_id": admin.Id,
	}
	err := h.db.Clauses(clause.Returning{}).Model(&order).Where("id=?", orderId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to update order")
		h.log.Error("failed to update order", err.Error())
		return
	}
	c.JSON(http.StatusOK, order)
}

// @Summary		 	Delete Order
// @Description	   	this api is delete order this api change status to 2
// @Tags			Order
// @Security		BearerAuth
// @Accept			json
// @Produce			json
// @Param           id   	path     int   true   "id"
// @Success			201		{object}	models.Orders
// @Failure			400,409	{object}	response
// @Failure			500		{object}	response
// @Router			/api/order/{id} [DELETE]
func (h *OrderController) DeleteOrder(c *gin.Context) {
	admin := h.GetAdmin(c)
	orderId := c.Param("id")
	var order models.Orders
	columns := map[string]interface{}{
		"deleted_at": timeNow(),
		"deleted_id": admin.Id,
	}
	err := h.db.Clauses(clause.Returning{}).Model(&order).Where("id=?", orderId).Updates(columns).Error
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to update order")
		h.log.Error("failed to update order", err.Error())
		return
	}
	c.JSON(http.StatusOK, order)
}
