package repository

import (
	"errors"
	"time"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/entity"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Interface for OrderRepository
type OrderRepository interface {
	BaseRepository[entity.Order]

	FindByID(
		db *gorm.DB,
		orderID string,
	) (*entity.Order, error)

	FindByIdempotencyKey(
		db *gorm.DB,
		userID string,
		idempotencyKey string,
	) (*entity.Order, error)

	UpdateStatus(
		tx *gorm.DB,
		orderID string,
		from enum.OrderStatus,
		to enum.OrderStatus,
	) error

	FindExpired(
		db *gorm.DB,
		now time.Time,
		limit int,
	) ([]entity.Order, error)
}

// Implementation of OrderRepository
type orderRepository struct {
	Repository[entity.Order]
	Log *logrus.Logger
}

func NewOrderRepository(log *logrus.Logger) OrderRepository {
	return &orderRepository{
		Log:  log,
	}
}

func (r *orderRepository) FindByID(
	db *gorm.DB,
	orderID string,
) (*entity.Order, error) {

	var order entity.Order
	err := db.
		Where("order_id = ?", orderID).
		Preload("Items").
		Take(&order).Error

	return &order, err
}

func (r *orderRepository) UpdateStatus(
	tx *gorm.DB,
	orderID string,
	from enum.OrderStatus,
	to enum.OrderStatus,
) error {

	result := tx.
		Model(&entity.Order{}).
		Where("order_id = ? AND status = ?", orderID, from).
		Update("status", to)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.Log.WithFields(logrus.Fields{
			"order_id":    orderID,
			"from_status": from,
			"to_status":   to,
		}).Error("no rows affected when updating order status")
		var ErrInvalidOrderTransition = errors.New("invalid order state transition")
		return ErrInvalidOrderTransition
	}
	return nil
}

func (r *orderRepository) FindByIdempotencyKey(
	db *gorm.DB,
	userID string,
	idempotencyKey string,
) (*entity.Order, error) {

	var order entity.Order
	err := db.
		Where("user_id = ? AND idempotency_key = ?", userID, idempotencyKey).
		Preload("Items").
		Take(&order).Error

	return &order, err
}

func (r *orderRepository) FindExpired(
	db *gorm.DB,
	now time.Time,
	limit int,
) ([]entity.Order, error) {

	var orders []entity.Order
	err := db.
		Where("status IN ? AND expires_at <= ?", enum.ExpirableOrderStatuses, now).
		Order("expires_at ASC").
		Limit(limit).
		Find(&orders).Error

	return orders, err
}

func (r *orderRepository) Delete(tx *gorm.DB, entity *entity.Order) error {
	if err := tx.Delete(entity).Error; err != nil {
		r.Log.WithFields(logrus.Fields{
			"order_id": entity.OrderID,	
		}).Error("failed to delete order")
		return err
	}
	return nil
}
