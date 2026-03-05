package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/theBaffo/subito-backend-challenge/internal/domain"
	"github.com/theBaffo/subito-backend-challenge/internal/handler"
	"github.com/theBaffo/subito-backend-challenge/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ---------------------------------------------------------------------------
// Mock service
// ---------------------------------------------------------------------------

type mockOrderService struct{ mock.Mock }

func (m *mockOrderService) CreateOrder(ctx context.Context, req service.CreateOrderRequest) (*domain.Order, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *mockOrderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func setupRouter(svc *mockOrderService) *gin.Engine {
	r := gin.New()
	h := handler.NewOrderHandler(svc)
	v1 := r.Group("/v1")
	{
		v1.POST("/orders", h.CreateOrder)
		v1.GET("/orders/:id", h.GetOrder)
	}
	return r
}

func makeOrder() *domain.Order {
	return &domain.Order{
		ID: "ord-abc12345",
		Items: []domain.OrderItem{
			{
				ProductID:  "prod-001",
				Name:       "Laptop Stand",
				Quantity:   2,
				UnitPrice:  decimal.NewFromFloat(39.99),
				TotalPrice: decimal.NewFromFloat(79.98),
				VATRate:    decimal.NewFromFloat(0.22),
				VATAmount:  decimal.NewFromFloat(14.42),
			},
		},
		TotalPrice: decimal.NewFromFloat(79.98),
		TotalVAT:   decimal.NewFromFloat(14.42),
		Status:     domain.StatusConfirmed,
		CreatedAt:  time.Now().UTC(),
	}
}

// ---------------------------------------------------------------------------
// POST /v1/orders
// ---------------------------------------------------------------------------

func TestCreateOrder_ValidRequest_Returns201(t *testing.T) {
	svc := &mockOrderService{}
	router := setupRouter(svc)

	expectedOrder := makeOrder()
	svc.On("CreateOrder", mock.Anything, mock.AnythingOfType("service.CreateOrderRequest")).
		Return(expectedOrder, nil)

	body := `{"items":[{"product_id":"prod-001","quantity":2}]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp domain.Order
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, expectedOrder.ID, resp.ID)
	assert.Equal(t, string(domain.StatusConfirmed), string(resp.Status))
}

func TestCreateOrder_InvalidJSON_Returns400(t *testing.T) {
	svc := &mockOrderService{}
	router := setupRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/orders", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_EmptyItems_Returns422(t *testing.T) {
	svc := &mockOrderService{}
	router := setupRouter(svc)

	svc.On("CreateOrder", mock.Anything, mock.AnythingOfType("service.CreateOrderRequest")).
		Return(nil, domain.ErrEmptyOrder)

	body := `{"items":[]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCreateOrder_InvalidQuantity_Returns422(t *testing.T) {
	svc := &mockOrderService{}
	router := setupRouter(svc)

	svc.On("CreateOrder", mock.Anything, mock.AnythingOfType("service.CreateOrderRequest")).
		Return(nil, domain.ErrInvalidQuantity)

	body := `{"items":[{"product_id":"prod-001","quantity":0}]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCreateOrder_ProductNotFound_Returns404(t *testing.T) {
	svc := &mockOrderService{}
	router := setupRouter(svc)

	svc.On("CreateOrder", mock.Anything, mock.AnythingOfType("service.CreateOrderRequest")).
		Return(nil, domain.ErrProductNotFound)

	body := `{"items":[{"product_id":"prod-999","quantity":1}]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateOrder_ResponseContainsRequiredFields(t *testing.T) {
	svc := &mockOrderService{}
	router := setupRouter(svc)

	svc.On("CreateOrder", mock.Anything, mock.AnythingOfType("service.CreateOrderRequest")).
		Return(makeOrder(), nil)

	body := `{"items":[{"product_id":"prod-001","quantity":2}]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))

	// Verify all required fields are present in the response
	assert.Contains(t, resp, "id", "response must include order ID")
	assert.Contains(t, resp, "total_price", "response must include total price")
	assert.Contains(t, resp, "total_vat", "response must include total VAT")
	assert.Contains(t, resp, "items", "response must include items")
	assert.Contains(t, resp, "status", "response must include status")

	items, ok := resp["items"].([]interface{})
	require.True(t, ok)
	require.Len(t, items, 1)

	item := items[0].(map[string]interface{})
	assert.Contains(t, item, "unit_price")
	assert.Contains(t, item, "total_price")
	assert.Contains(t, item, "vat_rate")
	assert.Contains(t, item, "vat_amount")
}

// ---------------------------------------------------------------------------
// GET /v1/orders/:id
// ---------------------------------------------------------------------------

func TestGetOrder_ExistingOrder_Returns200(t *testing.T) {
	svc := &mockOrderService{}
	router := setupRouter(svc)

	expected := makeOrder()
	svc.On("GetOrder", mock.Anything, "ord-abc12345").Return(expected, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/orders/ord-abc12345", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.Order
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, expected.ID, resp.ID)
}

func TestGetOrder_NonExistentOrder_Returns404(t *testing.T) {
	svc := &mockOrderService{}
	router := setupRouter(svc)

	svc.On("GetOrder", mock.Anything, "ord-nope").Return(nil, domain.ErrOrderNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/orders/ord-nope", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Contains(t, resp, "error")
}
