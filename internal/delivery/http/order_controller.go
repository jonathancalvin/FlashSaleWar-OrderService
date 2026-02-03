package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/application"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/domainerr"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/enum"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/domain/model"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/shared/util"
)

type OrderController struct {
	OrderUC application.OrderUseCase
	Log     *logrus.Logger
	Validate *validator.Validate
}

func NewOrderController(
	orderUC application.OrderUseCase,
	log *logrus.Logger,
	validator *validator.Validate,
) *OrderController {
	return &OrderController{
		OrderUC: orderUC,
		Log:     log,
		Validate: validator,
	}
}

func (h *OrderController) CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    string(enum.ErrorInvalidRequest),
					Message: "Invalid JSON format",
				},
			},
		)
		return
	}

	idempotencyKey := c.GetHeader("X-Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    string(enum.ErrorInvalidRequest),
					Message: "Idempotency-Key header is required",
				},
			},
		)
		return
	}

	req.IdempotencyKey = idempotencyKey

	userUUID, err := util.StringToUUID(req.UserID)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    string(enum.ErrorInvalidRequest),
					Message: "user_id must be a valid UUID",
				},
			},
		)
		return
	}

	req.UserUUID = userUUID

	if err := h.Validate.Struct(req); err != nil {
		h.Log.WithError(err).Warn("validation failed")
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    string(enum.ErrorValidationFailed),
					Message: h.formatValidationError(err.(validator.ValidationErrors)),
				},
			},
		)
		return
	}

	resp, err := h.OrderUC.CreateOrder(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		model.SuccessResponse[model.OrderResponse]{Data: *resp},
	)
}

func (h *OrderController) CancelOrder(c *gin.Context) {
	h.updateOrderStatus(c, enum.StatusCancelled)
}

func (h *OrderController) updateOrderStatus(c *gin.Context, status enum.OrderStatus) {
	orderID := c.Param("orderID")

	if orderID == "" {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    string(enum.ErrorInvalidRequest),
					Message: "orderID is required",
				},
			},
		)
		return
	}

	if _, err := util.StringToUUID(orderID); err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    string(enum.ErrorInvalidRequest),
					Message: "orderID must be a valid UUID",
				},
			},
		)
		return
	}

	resp, err := h.OrderUC.UpdateOrderStatus(
		c.Request.Context(),
		orderID,
		status,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse[model.OrderResponse]{Data: *resp})
}

func (h *OrderController) handleError(c *gin.Context, err error) {
	var (
		httpStatus int
		errorCode  enum.ErrorCode
		message    string
	)

	switch err{
	case gorm.ErrRecordNotFound,
		domainerr.ErrOrderNotFound:

		httpStatus = http.StatusNotFound
		errorCode = enum.ErrorOrderNotFound
		message = err.Error()

	case domainerr.ErrInvalidOrderTransition:

		httpStatus = http.StatusConflict
		errorCode = enum.ErrorOrderInvalidState
		message = err.Error()

	case domainerr.ErrOrderExpired:

		httpStatus = http.StatusConflict
		errorCode = enum.ErrorOrderAlreadyExpired
		message = err.Error()

	default:
		h.Log.WithError(err).Error("unexpected handler error")

		httpStatus = http.StatusInternalServerError
		errorCode = enum.ErrorInternal
		message = "internal server error"
	}

	c.JSON(
		httpStatus,
		model.ErrorResponse{
			Error: model.ErrorBody{
				Code:    string(errorCode),
				Message: message,
			},
		},
	)
}

func (h *OrderController) formatValidationError(ve validator.ValidationErrors) string {
    var errors []string
    for _, fe := range ve {
        errors = append(errors, fmt.Sprintf("field %s is %s", fe.Field(), fe.Tag()))
    }
    return strings.Join(errors, ", ")
}

