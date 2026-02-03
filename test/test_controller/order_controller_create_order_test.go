package test_controller

import (
	"net/http"
	"testing"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder_Success(t *testing.T) {
	r, mockUC := setupTestController()

	reqBody := map[string]any{
		"user_id": "550e8400-e29b-41d4-a716-446655440000",
		"currency": "IDR",
		"total_amount": 10000,
		"items": []map[string]any{
			{
				"sku_id": "sku-1",
				"quantity": 1,
				"price": 10000,
			},
		},
	}

	mockUC.On(
		"CreateOrder",
		mock.Anything,
		mock.AnythingOfType("model.CreateOrderRequest"),
	).Return(&model.OrderResponse{
		ID:     "order-1",
		Status: enum.StatusCreated,
	}, nil).Once()

	w := performRequest(
		r,
		http.MethodPost,
		"/order",
		reqBody,
		map[string]string{
			"X-Idempotency-Key": "idem-123",
		},
	)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateOrder_InvalidJSON(t *testing.T) {
	r, mockUC := setupTestController()

	invalidJSON := "{invalid-json"

	w := performRequest(
		r,
		http.MethodPost,
		"/order",
		invalidJSON,
		map[string]string{
			"X-Idempotency-Key": "idem-123",
		},
	)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "CreateOrder")
}

func TestCreateOrder_ValidationFailed(t *testing.T) {
	r, mockUC := setupTestController()

	reqBody := map[string]any{
		"user_id": "wrong-uuid",
		"currency": "IDR",
		"total_amount": 10000,
		"items": []map[string]any{
			{
				"sku_id": "sku-1",
				"quantity": 1,
				"price": 10000,
			},
		},
	}

	w := performRequest(
		r,
		http.MethodPost,
		"/order",
		reqBody,
		map[string]string{
			"X-Idempotency-Key": "idem-123",
		},
	)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "CreateOrder")
}

func TestCreateOrder_MissingIdempotencyKey(t *testing.T) {
	r, mockUC := setupTestController()

	reqBody := map[string]any{
		"user_id": "550e8400-e29b-41d4-a716-446655440000",
		"currency": "IDR",
		"total_amount": 10000,
		"items": []map[string]any{
			{
				"sku_id": "sku-1",
				"quantity": 1,
				"price": 10000,
			},
		},
	}

	w := performRequest(
		r,
		http.MethodPost,
		"/order",
		reqBody,
		nil,
	)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "CreateOrder")
}
