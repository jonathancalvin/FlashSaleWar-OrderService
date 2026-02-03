package route

import (
	"github.com/jonathancalvin/FlashSaleWar-OrderService/internal/delivery/http"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App *gin.Engine
	Log *logrus.Logger
	OrderController *http.OrderController
}

func (c *RouteConfig) Setup() {
	api := c.App.Group("/api/v1")
	{
		orders := api.Group("/orders")
		{
			orders.POST("/order", c.OrderController.CreateOrder)
			orders.POST("/:orderID/cancel", c.OrderController.CancelOrder)
		}
	}
}