package test_usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	uc := setupOrderUseCase()

	req := model.CreateOrderRequest{
		UserID:         uuid.New().String(),
		IdempotencyKey: "idem-123",
		Currency:       "IDR",
		Items: []model.CreateOrderItem{
			{
				SkuID:    "sku-1",
				Quantity: 2,
				Price:    10000,
			},
		},
	}

	ctx := context.Background()

	// First call
	res1, err := uc.CreateOrder(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res1)
	assert.Equal(t, enum.StatusCreated, res1.Status)
	assert.NotZero(t, res1.ExpiredAt)

	// Second call (idempotency)
	res2, err := uc.CreateOrder(ctx, req)
	assert.NoError(t, err)

	assert.Equal(t, res1.ID, res2.ID)
}
