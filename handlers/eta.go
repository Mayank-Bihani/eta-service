package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mayank-Bihani/eta-service/models"
	"github.com/Mayank-Bihani/eta-service/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

)

type ETAHandler struct {
	Service *services.ETAService
	DB      *sqlx.DB
	Redis   *redis.Client
}

func NewETAHandler(service *services.ETAService, db *sqlx.DB, redis *redis.Client) *ETAHandler {
	return &ETAHandler{Service: service, DB: db, Redis: redis}
}

func (h *ETAHandler) GetETA(c *gin.Context) {
	var req models.ETARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eta, surge, queueDepth, distance, err := h.Service.CalculateETA(
		c.Request.Context(),
		req.RestaurantID,
		req.DeliveryLat,
		req.DeliveryLng,
		req.ItemCount,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Save order to DB
	_, err = h.DB.ExecContext(c.Request.Context(), `
		INSERT INTO orders (restaurant_id, delivery_lat, delivery_lng, item_count, estimated_eta)
		VALUES ($1, $2, $3, $4, $5)`,
		req.RestaurantID, req.DeliveryLat, req.DeliveryLng, req.ItemCount, eta,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save order"})
		return
	}

	c.JSON(http.StatusOK, models.ETAResponse{
		EstimatedETA: eta,
		RestaurantID: req.RestaurantID,
		SurgeFactor:  surge,
		QueueDepth:   queueDepth,
		DistanceKm:   distance,
	})
}

func (h *ETAHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	err := h.DB.GetContext(c.Request.Context(), &order, `
		SELECT id, restaurant_id, delivery_lat, delivery_lng, item_count, estimated_eta, created_at
		FROM orders WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *ETAHandler) SetQueueDepth(c *gin.Context) {
	var payload struct {
		RestaurantID string `json:"restaurant_id" binding:"required"`
		Depth        int    `json:"depth" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := fmt.Sprintf("queue:%s", payload.RestaurantID)
	h.Service.Redis.Set(c.Request.Context(), key, payload.Depth, 2*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"message":       "queue depth updated",
		"restaurant_id": payload.RestaurantID,
		"depth":         payload.Depth,
	})
}