package models

import "time"

type Order struct {
	ID             int       `db:"id" json:"id"`
	RestaurantID   string    `db:"restaurant_id" json:"restaurant_id"`
	DeliveryLat    float64   `db:"delivery_lat" json:"delivery_lat"`
	DeliveryLng    float64   `db:"delivery_lng" json:"delivery_lng"`
	ItemCount      int       `db:"item_count" json:"item_count"`
	EstimatedETA   int       `db:"estimated_eta" json:"estimated_eta"` // in minutes
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type ETARequest struct {
	RestaurantID string  `json:"restaurant_id" binding:"required"`
	DeliveryLat  float64 `json:"delivery_lat" binding:"required"`
	DeliveryLng  float64 `json:"delivery_lng" binding:"required"`
	ItemCount    int     `json:"item_count" binding:"required"`
}

type ETAResponse struct {
	EstimatedETA  int     `json:"estimated_eta_minutes"`
	RestaurantID  string  `json:"restaurant_id"`
	SurgeFactor   float64 `json:"surge_factor"`
	QueueDepth    int     `json:"queue_depth"`
	DistanceKm    float64 `json:"distance_km"`
}