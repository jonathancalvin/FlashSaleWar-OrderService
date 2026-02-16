package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model/converter"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/infrastructure/repository"
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

	CancelOrder(
		ctx context.Context, 
		req model.CancelOrderRequest,
	) (*model.OrderResponse, error)

	MarkOrderPaid(
		ctx context.Context,
		orderID string,
	) (*model.OrderResponse, error)
}

type orderUseCase struct {
	DB        *gorm.DB
	Log       *logrus.Logger

	OrderRepo repository.OrderRepository
	OutboxRepo repository.OutboxRepository
}

func NewOrderUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	orderRepo repository.OrderRepository,
	outboxRepo repository.OutboxRepository,
) OrderUseCase {
	return &orderUseCase{
		DB:        db,
		Log:       log,
		OrderRepo: orderRepo,
		OutboxRepo: outboxRepo,
	}
}

func (s *orderUseCase) CreateOrder(
	ctx context.Context,
	req model.CreateOrderRequest,
) (*model.OrderResponse, error) {

	s.Log.WithFields(logrus.Fields{
		"user_id":          req.UserUUID,
		"idempotency_key":  req.IdempotencyKey,
		"items_count":      len(req.Items),
	}).Info("create order requested")

	order, err := s.createOrderWithTx(ctx, req)
	if err != nil {
		s.Log.WithError(err).WithFields(logrus.Fields{
			"user_id":         req.UserUUID,
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
			req.UserUUID,
			req.IdempotencyKey,
		)

		if err == nil {
			s.Log.WithFields(logrus.Fields{
				"order_id":        existing.OrderID,
				"user_id":         req.UserUUID,
				"idempotency_key": req.IdempotencyKey,
			}).Warn("idempotency hit, returning existing order")

			result = existing
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 2. Prepare entities

		orderEntity := entity.NewOrder(
			req.UserUUID,
			req.IdempotencyKey,
			enum.StatusCreated,
			*enum.CalculateExpiryTime(enum.StatusCreated),
			req.Currency,
			req.TotalAmount,
		)

		for _, item := range req.Items {
			orderEntity.AddItem(
				item.SkuID,
				item.Quantity,
				item.Price,
				req.Currency,
			)
		}

		payload := model.OrderCreatedPayload{
			OrderID:  orderEntity.OrderID.String(),
			UserID:   orderEntity.UserID.String(),
			TotalAmount: orderEntity.TotalAmount,
			Currency:    orderEntity.Currency,
			Items:      converter.OrderItemsToPayloads(orderEntity),
			CreatedAt:  orderEntity.CreatedAt,
		}

		jsonPayload, _ := json.Marshal(payload)

		outboxEntity := entity.NewOutboxEvent(
			orderEntity.OrderID,
			enum.EventTypeOrderCreated,
			jsonPayload,
			"PENDING",
		)


		// 3. Persist
		if err := s.OrderRepo.Create(tx, orderEntity); err != nil {
			return err
		}

		if err := s.OutboxRepo.Create(tx, outboxEntity); err != nil {
			return err
		}

		result = orderEntity
		return nil
	})

	return result, err
}

func (s *orderUseCase) CancelOrder(
    ctx context.Context, 
    req model.CancelOrderRequest,
) (*model.OrderResponse, error) {
    var updatedOrder *entity.Order

	s.Log.WithFields(logrus.Fields{
		"order_id":        req.OrderID,
		"reason":           req.Reason,
	}).Info("cancel order requested")

    err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 1. Fetch Order dan validate ownership
        order, err := s.OrderRepo.FindByID(tx, req.OrderID)
        if err != nil {
            s.Log.WithError(err).WithField("order_id", req.OrderID).Error("failed to fetch order for cancellation")
            return err
        }

        // 2. Validate Transition
        if err := enum.ValidateTransition(order.Status, enum.StatusCancelled); err != nil {
            s.Log.WithFields(logrus.Fields{
                "order_id":    req.OrderID,
                "current_status": order.Status,
                "target_status":  enum.StatusCancelled,
            }).Warn("invalid order status transition for cancellation")
            return err
        }

        // 3. Update Status
        if err := s.OrderRepo.UpdateStatus(
            tx, 
            req.OrderID, 
            order.Status, 
            enum.StatusCancelled, 
            nil,
        ); err != nil {
            s.Log.WithError(err).WithField("order_id", req.OrderID).Error("failed to update order status to cancelled")
            return err
        }
        
        // 4. Reload updated order inside TX
        updatedOrder, err = s.OrderRepo.FindByID(tx, req.OrderID)
        if err != nil {
            return err
        }

        // 5. Create Outbox Event
        cancelPayload := model.OrderCancelledPayload{
            OrderID:     updatedOrder.OrderID.String(),
            CancelledAt: time.Now(),
            Reason:      req.Reason,
        }

        jsonPayload, err := json.Marshal(cancelPayload)
        if err != nil {
            s.Log.WithError(err).WithField("order_id", updatedOrder.OrderID).Error("failed to marshal cancel payload")
            return err
        }

        outboxEvent := entity.NewOutboxEvent(
            updatedOrder.OrderID,
            enum.EventTypeOrderCancelled,
            jsonPayload,
            "PENDING",
        )

        if err := s.OutboxRepo.Create(tx, outboxEvent); err != nil {
            s.Log.WithError(err).WithField("order_id", updatedOrder.OrderID).Error("failed to persist cancel outbox event")
            return err
        }

        s.Log.WithFields(logrus.Fields{
            "order_id": updatedOrder.OrderID,
            "user_id":  updatedOrder.UserID,
            "reason":   req.Reason,
        }).Info("order successfully cancelled and outbox event created")

        return nil
    })

    if err != nil {
        return nil, err
    }

    return converter.OrderToResponse(updatedOrder), nil
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

func (s *orderUseCase) MarkOrderPaid(
	ctx context.Context,
	orderID string,
) (*model.OrderResponse, error) {

	var updatedOrder *entity.Order

	s.Log.WithFields(logrus.Fields{
		"order_id": orderID,
	}).Info("mark order paid requested")

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Fetch order
		order, err := s.OrderRepo.FindByID(tx, orderID)
		if err != nil {
			s.Log.WithError(err).
				WithField("order_id", orderID).
				Error("failed to fetch order for mark paid")
			return err
		}

		// 2. Validate transition
		if err := enum.ValidateTransition(order.Status, enum.StatusPaid); err != nil {
			s.Log.WithFields(logrus.Fields{
				"order_id":      orderID,
				"current_status": order.Status,
				"target_status":  enum.StatusPaid,
			}).Warn("invalid order status transition to paid")
			return err
		}

		// 3. Optional: validate expiry
		if !order.ExpiredAt.IsZero() && time.Now().After(order.ExpiredAt) {
			s.Log.WithFields(logrus.Fields{
				"order_id": orderID,
				"expired_at": order.ExpiredAt,
			}).Warn("cannot mark order paid because order already expired")
			return errors.New("order expired")
		}

		// 4. Update status
		if err := s.OrderRepo.UpdateStatus(
			tx,
			orderID,
			order.Status,
			enum.StatusPaid,
			nil,
		); err != nil {
			s.Log.WithError(err).
				WithField("order_id", orderID).
				Error("failed to update order status to paid")
			return err
		}

		// 5. Reload updated order
		updatedOrder, err = s.OrderRepo.FindByID(tx, orderID)
		if err != nil {
			return err
		}

		// 6. Create outbox event
		payload := model.OrderPaidPayload{
			OrderID: updatedOrder.OrderID.String(),
			PaidAt:  time.Now(),
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			s.Log.WithError(err).
				WithField("order_id", orderID).
				Error("failed to marshal order paid payload")
			return err
		}

		outboxEvent := entity.NewOutboxEvent(
			updatedOrder.OrderID,
			enum.EventTypeOrderPaid,
			jsonPayload,
			"PENDING",
		)

		if err := s.OutboxRepo.Create(tx, outboxEvent); err != nil {
			s.Log.WithError(err).
				WithField("order_id", orderID).
				Error("failed to persist order paid outbox event")
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.Log.WithFields(logrus.Fields{
		"order_id": updatedOrder.OrderID,
		"user_id":  updatedOrder.UserID,
		"status":   updatedOrder.Status,
	}).Info("order successfully marked as paid and outbox event created")

	return converter.OrderToResponse(updatedOrder), nil
}