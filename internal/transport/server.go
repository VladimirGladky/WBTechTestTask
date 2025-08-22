package transport

import (
	"WBTechTestTask/internal/config"
	"WBTechTestTask/internal/models"
	"WBTechTestTask/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OderServer struct {
	cfg     *config.Config
	ctx     context.Context
	service service.OrderServiceInterface
}

func New(cfg *config.Config, ctx context.Context, service service.OrderServiceInterface) *OderServer {
	return &OderServer{
		cfg:     cfg,
		ctx:     ctx,
		service: service,
	}
}

func (s *OderServer) Start() error {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/create", s.createOrder())
		api.GET("/order/:orderId", s.getOrder())
	}
	return nil
}

func (s *OderServer) createOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
				return
			}
		}()
		if c.Request.Method != http.MethodPost {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		id, err := s.service.Create(order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"orderId": id})
	}
}

func (s *OderServer) getOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
				return
			}
		}()
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}
		id := c.Param("orderId")
		order, err := s.service.GetOrder(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.Order{
			OrderId:           order.OrderId,
			TrackNumber:       order.TrackNumber,
			Entry:             order.Entry,
			Delivery:          order.Delivery,
			Payment:           order.Payment,
			Items:             order.Items,
			Locale:            order.Locale,
			InternalSignature: order.InternalSignature,
			CustomerId:        order.CustomerId,
			DeliveryService:   order.DeliveryService,
			Shardkey:          order.Shardkey,
			SmId:              order.SmId,
			DateCreated:       order.DateCreated,
			OofShard:          order.OofShard,
		})
	}
}
