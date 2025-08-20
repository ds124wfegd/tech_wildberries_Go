package transport

import (
	"net/http"

	"github.com/ds124wfegd/tech_wildberries_Go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// initialization routing
func (h *OrderHandler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("./internal/web/templates/*")

	order := router.Group("/order")
	{
		order.GET("/:order_uid", h.GetOrderByID)

		// route for html-page
		router.GET("/", func(c *gin.Context) {
			c.HTML(200, "1.html", gin.H{
				"title": "Поиск заказа",
			})
		})

	}
	return router
}
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	// recieve parametr from string
	orderUID := c.Query("order_uid")
	if orderUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uid parameter is required"})
		return
	}

	// recieves order from service
	order, err := (h.service).GetByUID(c.Request.Context(), orderUID)
	if err != nil {
		logrus.Printf("Error getting order %s: %v", orderUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// send answer
	c.JSON(http.StatusOK, order)
}
