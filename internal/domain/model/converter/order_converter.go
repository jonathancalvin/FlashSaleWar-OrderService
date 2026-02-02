package converter

import (
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/shared/util"
)

func OrderToResponse(order *entity.Order) *model.OrderResponse {
	return &model.OrderResponse{
		ID:          util.UUIDToString(order.OrderID),
		UserID:      util.UUIDToString(order.UserID),
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Currency:    order.Currency,
		CreatedAt:   order.CreatedAt.Unix(),
		UpdatedAt:   order.UpdatedAt.Unix(),
		ExpiredAt: 	 order.ExpiredAt.Unix(),
	}
}

func OrdersToResponses(orders []entity.Order) []*model.OrderResponse {
	responses := make([]*model.OrderResponse, 0, len(orders))
	for i := range orders {
		responses = append(responses, OrderToResponse(&orders[i]))
	}
	return responses
}