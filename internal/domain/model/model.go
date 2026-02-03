package model

type WebResponse[T any] struct {
	Data   T             `json:"data"`
	Errors string        `json:"errors,omitempty"`
}

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}