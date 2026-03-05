package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/theBaffo/subito-backend-challenge/internal/domain"
	"github.com/theBaffo/subito-backend-challenge/internal/handler"
)

// ---------------------------------------------------------------------------
// Mock product service
// ---------------------------------------------------------------------------

type mockProductService struct{ mock.Mock }

func (m *mockProductService) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *mockProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func setupProductRouter(svc *mockProductService) *gin.Engine {
	r := gin.New()
	h := handler.NewProductHandler(svc)
	v1 := r.Group("/v1")
	{
		v1.GET("/products", h.ListProducts)
		v1.GET("/products/:id", h.GetProduct)
	}
	return r
}

var sampleProducts = []domain.Product{
	{
		ID:         "prod-001",
		Name:       "Laptop Stand",
		GrossPrice: decimal.NewFromFloat(39.99),
		VATRate:    domain.VATStandard,
		Category:   "accessories",
	},
	{
		ID:         "prod-002",
		Name:       "Coffee Beans",
		GrossPrice: decimal.NewFromFloat(14.90),
		VATRate:    domain.VATReduced,
		Category:   "food",
	},
}

func TestListProducts_Returns200WithProducts(t *testing.T) {
	svc := &mockProductService{}
	router := setupProductRouter(svc)

	svc.On("GetAllProducts", mock.Anything).Return(sampleProducts, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	products, ok := resp["products"].([]interface{})
	require.True(t, ok)
	assert.Len(t, products, 2)
}

func TestGetProduct_ExistingProduct_Returns200(t *testing.T) {
	svc := &mockProductService{}
	router := setupProductRouter(svc)

	svc.On("GetProduct", mock.Anything, "prod-001").Return(&sampleProducts[0], nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/products/prod-001", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.Product
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "prod-001", resp.ID)
}

func TestGetProduct_NonExistent_Returns404(t *testing.T) {
	svc := &mockProductService{}
	router := setupProductRouter(svc)

	svc.On("GetProduct", mock.Anything, "prod-999").Return(nil, domain.ErrProductNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/products/prod-999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
