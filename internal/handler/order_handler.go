package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theBaffo/subito-backend-challenge/internal/domain"
	"github.com/theBaffo/subito-backend-challenge/internal/service"
)

// orderServicer is the interface the handler depends on.
type orderServicer interface {
	CreateOrder(ctx context.Context, req service.CreateOrderRequest) (*domain.Order, error)
	GetOrder(ctx context.Context, id string) (*domain.Order, error)
}

// OrderHandler handles HTTP requests for the /v1/orders resource.
type OrderHandler struct {
	service orderServicer
}

// NewOrderHandler constructs an OrderHandler.
func NewOrderHandler(svc orderServicer) *OrderHandler {
	return &OrderHandler{service: svc}
}

// CreateOrder godoc
// POST /v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid request body: "+err.Error()))
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), req)

	if err != nil {
		c.JSON(domainErrorToStatus(err), errorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder godoc
// GET /v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	order, err := h.service.GetOrder(c.Request.Context(), c.Param("id"))

	if err != nil {
		c.JSON(domainErrorToStatus(err), errorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, order)
}
