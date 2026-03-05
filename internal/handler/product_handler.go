package handler

import (
	"context"
	"net/http"

	"github.com/theBaffo/subito-backend-challenge/internal/domain"
	"github.com/gin-gonic/gin"
)

// productServicer is the interface the handler depends on.
// Depending on an interface (not the concrete service) keeps the handler testable.
type productServicer interface {
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
}

// ProductHandler handles HTTP requests for the /v1/products resource.
type ProductHandler struct {
	service productServicer
}

// NewProductHandler constructs a ProductHandler.
func NewProductHandler(svc productServicer) *ProductHandler {
	return &ProductHandler{service: svc}
}

// ListProducts godoc
// GET /v1/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("failed to retrieve products"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProduct godoc
// GET /v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	product, err := h.service.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(domainErrorToStatus(err), errorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, product)
}
