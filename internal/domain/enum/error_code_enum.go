package enum

type ErrorCode string

const (
	// --- General ---
	ErrorInternal        ErrorCode = "INTERNAL_ERROR"
	ErrorInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrorUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrorForbidden      ErrorCode = "FORBIDDEN"
	ErrorNotFound       ErrorCode = "NOT_FOUND"
	ErrorConflict       ErrorCode = "CONFLICT"

	// --- Order domain ---
	ErrorOrderNotFound        ErrorCode = "ORDER_NOT_FOUND"
	ErrorOrderInvalidState   ErrorCode = "ORDER_INVALID_STATE"
	ErrorOrderAlreadyExpired ErrorCode = "ORDER_ALREADY_EXPIRED"
	ErrorOrderIdempotencyHit ErrorCode = "ORDER_IDEMPOTENT_REPLAY"

	// --- Validation ---
	ErrorValidationFailed ErrorCode = "VALIDATION_FAILED"
)
