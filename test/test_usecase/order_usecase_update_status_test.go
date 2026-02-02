package test_usecase

import (
	"context"
	"testing"
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateOrderStatus_ValidTransition(t *testing.T) {
	uc := setupOrderUseCase()

	ctx := context.Background()

	// Create order
	createRes, err := uc.CreateOrder(ctx, model.CreateOrderRequest{
		UserID:         "user-1",
		IdempotencyKey: "idem-1",
		Currency:       "IDR",
		Items: []model.CreateOrderItem{
			{SkuID: "sku-1", Quantity: 1, Price: 1000},
		},
	})
	require.NoError(t, err)

	now := time.Now()

	// Update status
	updated, err := uc.UpdateOrderStatus(
		ctx,
		createRes.ID,
		enum.StatusReserved,
	)
	assert.NoError(t, err)

	assert.Equal(t, enum.StatusReserved, updated.Status)

	// expires_at must be in the future
	assert.True(t, updated.ExpiredAt > now.Unix())
}

func TestUpdateOrderStatus_InvalidTransition(t *testing.T) {
	uc := setupOrderUseCase()

	ctx := context.Background()

	createRes, err := uc.CreateOrder(ctx, model.CreateOrderRequest{
		UserID:         "user-1",
		IdempotencyKey: "idem-2",
		Currency:       "IDR",
		Items: []model.CreateOrderItem{
			{SkuID: "sku-1", Quantity: 1, Price: 1000},
		},
	})
	require.NoError(t, err)

	// CREATED → PAID (invalid)
	_, err = uc.UpdateOrderStatus(
		ctx,
		createRes.ID,
		enum.StatusPaid,
	)

	assert.Error(t, err)
}
