package http

import "github.com/gin-gonic/gin"

type OrderHandler struct{}

func (h *OrderHandler) Create(c *gin.Context) {
	// map HTTP → use case
}
