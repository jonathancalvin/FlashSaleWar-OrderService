package application

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model/converter"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/repository"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/shared/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUseCase interface {
	CreateOrder(
		ctx context.Context,
		req model.CreateOrderRequest,
	) (*model.OrderResponse, error)

	UpdateOrderStatus(
		ctx context.Context,
		orderID string,
		to enum.OrderStatus,
	) (*model.OrderResponse, error)
}

type orderUseCase struct {
	DB        *gorm.DB
	Log       *logrus.Logger
	Validator *validator.Validate

	OrderRepo repository.OrderRepository
}

func NewOrderUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validator *validator.Validate,
	orderRepo repository.OrderRepository,
) OrderUseCase {
	return &orderUseCase{
		DB:        db,
		Log:       log,
		Validator: validator,
		OrderRepo: orderRepo,
	}
}

func (s *orderUseCase) CreateOrder(
	ctx context.Context,
	req model.CreateOrderRequest,
) (*model.OrderResponse, error) {

	s.Log.WithFields(logrus.Fields{
		"user_id":          req.UserID,
		"idempotency_key":  req.IdempotencyKey,
		"items_count":      len(req.Items),
	}).Info("create order requested")

	order, err := s.createOrderWithTx(ctx, req)
	if err != nil {
		s.Log.WithError(err).WithFields(logrus.Fields{
			"user_id":         req.UserID,
			"idempotency_key": req.IdempotencyKey,
		}).Error("failed to create order")

		return nil, err
	}

	s.Log.WithFields(logrus.Fields{
		"order_id": order.OrderID,
		"user_id":  order.UserID,
		"status":   order.Status,
	}).Info("order created successfully")

	return converter.OrderToResponse(order), nil
}

func (s *orderUseCase) createOrderWithTx(
	ctx context.Context,
	req model.CreateOrderRequest,
) (*entity.Order, error) {

	var result *entity.Order

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Idempotency check
		existing, err := s.OrderRepo.FindByIdempotencyKey(
			tx,
			req.UserID,
			req.IdempotencyKey,
		)

		if err == nil {
			s.Log.WithFields(logrus.Fields{
				"order_id":        existing.OrderID,
				"user_id":         req.UserID,
				"idempotency_key": req.IdempotencyKey,
			}).Warn("idempotency hit, returning existing order")

			result = existing
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 2. Build aggregate
		userID, err := util.StringToUUID(req.UserID)
		if err != nil {
			s.Log.Error("Invalid User UUID format")
			return err
    	}

		order := entity.NewOrder(
			userID,
			req.IdempotencyKey,
			enum.StatusCreated,
			*enum.CalculateExpiryTime(enum.StatusCreated),
			req.Currency,
			req.TotalAmount,
		)

		for _, item := range req.Items {
			order.AddItem(
				item.SkuID,
				item.Quantity,
				item.Price,
				req.Currency,
			)
		}

		// 3. Persist
		if err := s.OrderRepo.Create(tx, order); err != nil {
			return err
		}

		result = order
		return nil
	})

	return result, err
}

func (s *orderUseCase) UpdateOrderStatus(
	ctx context.Context,
	orderID string,
	to enum.OrderStatus,
) (*model.OrderResponse, error) {

	var updatedOrder *entity.Order

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		order, err := s.OrderRepo.FindByID(tx, orderID)
		if err != nil {
			s.Log.WithError(err).
				WithField("order_id", orderID).
				Error("failed to fetch order")
			return err
		}

		from := order.Status

		// 1. Validate transition
		if err := enum.ValidateTransition(from, to); err != nil {
			s.Log.WithFields(logrus.Fields{
				"order_id": orderID,
				"from":     from,
				"to":       to,
			}).Warn("invalid order status transition")
			return err
		}

		// 2. Determine expired_at
		newExpiredAt := enum.CalculateExpiryTime(to)

		// 3. Update status
		if err := s.OrderRepo.UpdateStatus(
			tx,
			orderID,
			from,
			to,
			newExpiredAt,
		); err != nil {
			s.Log.WithError(err).
				WithField("order_id", orderID).
				Error("failed to update order status")
			return err
		}

		// 4. Reload updated order inside TX
		updatedOrder, err = s.OrderRepo.FindByID(tx, orderID)
		if err != nil {
			return err
		}

		s.Log.WithFields(logrus.Fields{
			"order_id": orderID,
			"from":     from,
			"to":       to,
		}).Info("order status updated")

		return nil
	})

	if err != nil {
		return nil, err
	}

	return converter.OrderToResponse(updatedOrder), nil
}

