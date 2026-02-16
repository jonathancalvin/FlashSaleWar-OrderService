package test_controller

import (
	"bytes"
	"encoding/json"
	"io"
	nethttp "net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/delivery/http"
	"github.com/jonathancalvin/FlashSaleWar-OrderService/test/mock"
	"github.com/sirupsen/logrus"
)

func setupTestController() (*gin.Engine, *mock.MockOrderUseCase) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	mockUC := new(mock.MockOrderUseCase)

	ctrl := http.NewOrderController(mockUC, logrus.New(), validator.New())

	r.POST("/order", ctrl.CreateOrder)
	r.POST("/order/:orderID/cancel", func(c *gin.Context) {
		ctrl.CancelOrder(c)
	})

	return r, mockUC
}

func performRequest(
	r nethttp.Handler,
	method, path string,
	body any,
	headers map[string]string,
) *httptest.ResponseRecorder {

	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
