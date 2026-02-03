package test_controller

import (
	"net/http"
	"testing"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/domainerr"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCancelOrder_Success(t *testing.T) {
	r, mockUC := setupTestController()

	orderID := "550e8400-e29b-41d4-a716-446655440000"

	mockUC.
		On(
			"UpdateOrderStatus",
			mock.Anything,
			orderID,
			enum.StatusCancelled,
		).
		Return(&model.OrderResponse{
			ID:     orderID,
			Status: enum.StatusCancelled,
		}, nil).
		Once()

	w := performRequest(
		r,
		http.MethodPost,
		"/order/"+orderID+"/cancel",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCancelOrder_InvalidOrderID(t *testing.T) {
	r, mockUC := setupTestController()

	w := performRequest(
		r,
		http.MethodPost,
		"/order/not-a-uuid/cancel",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "UpdateOrderStatus")
}

func TestCancelOrder_NotFound(t *testing.T) {
	r, mockUC := setupTestController()

	orderID := "550e8400-e29b-41d4-a716-446655440000"

	mockUC.
		On(
			"UpdateOrderStatus",
			mock.Anything,
			orderID,
			enum.StatusCancelled,
		).
		Return(nil, domainerr.ErrOrderNotFound).
		Once()

	w := performRequest(
		r,
		http.MethodPost,
		"/order/"+orderID+"/cancel",
		nil,
		nil,
	)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}