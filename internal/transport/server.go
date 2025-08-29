package transport

import (
	"WBTechTestTask/internal/config"
	"WBTechTestTask/internal/models"
	"WBTechTestTask/internal/service"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")
	{
		api.POST("/create", s.createOrder())
		api.GET("/order/:orderId", s.getOrder())
		api.GET("/stresstest/:orderId", s.stressTest())
	}
	return router.Run(s.cfg.Host + ":" + s.cfg.Port)
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
		var order *models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
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
			OrderUid:          order.OrderUid,
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

func (s *OderServer) stressTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("orderId")
		iterations := 100000

		start := time.Now()
		_, err := s.service.GetOrder(orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		firstTime := time.Since(start)

		cacheStart := time.Now()
		for i := 0; i < iterations; i++ {
			_, err := s.service.GetOrder(orderID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		cacheTime := time.Since(cacheStart)
		avgCacheTime := cacheTime / time.Duration(iterations)

		c.JSON(http.StatusOK, gin.H{
			"first_request_time": firstTime.String(),
			"cache_requests":     iterations,
			"total_cache_time":   cacheTime.String(),
			"average_cache_time": avgCacheTime.String(),
		})
	}
}
