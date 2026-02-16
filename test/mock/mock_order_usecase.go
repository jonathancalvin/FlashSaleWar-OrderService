package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
)

type MockOrderUseCase struct {
	mock.Mock
}

func (m *MockOrderUseCase) CreateOrder(
	ctx context.Context,
	req model.CreateOrderRequest,
) (*model.OrderResponse, error) {

	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.OrderResponse), args.Error(1)
}

func (m *MockOrderUseCase) UpdateOrderStatus(
	ctx context.Context,
	orderID string,
	status enum.OrderStatus,
) (*model.OrderResponse, error) {

	args := m.Called(ctx, orderID, status)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.OrderResponse), args.Error(1)
}

func (m *MockOrderUseCase) CancelOrder(
	ctx context.Context,
	req model.CancelOrderRequest,
) (*model.OrderResponse, error) {
	
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.OrderResponse), args.Error(1)
}

func (m *MockOrderUseCase) MarkOrderPaid(
	ctx context.Context,
	orderID string,
) (*model.OrderResponse, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*model.OrderResponse), args.Error(1)
}