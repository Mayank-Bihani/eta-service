package services

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// Restaurant coordinates — in production this would come from a DB
var restaurantCoords = map[string][2]float64{
	"REST001": {28.6139, 77.2090}, // Delhi
	"REST002": {28.5355, 77.3910}, // Noida
	"REST003": {28.4595, 77.0266}, // Gurugram
}

type ETAService struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func NewETAService(db *sqlx.DB, redis *redis.Client) *ETAService {
	return &ETAService{DB: db, Redis: redis}
}

// Haversine formula — calculates distance between two lat/lng points in km
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371.0 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// Surge multiplier based on time of day
func getSurgeFactor() float64 {
	hour := time.Now().Hour()
	switch {
	case hour >= 12 && hour <= 14: // lunch rush
		return 1.4
	case hour >= 19 && hour <= 21: // dinner rush
		return 1.6
	default:
		return 1.0
	}
}

// Get queue depth from Redis, fall back to DB if cache miss
func (s *ETAService) getQueueDepth(ctx context.Context, restaurantID string) (int, error) {
	key := fmt.Sprintf("queue:%s", restaurantID)

	// Try Redis first
	val, err := s.Redis.Get(ctx, key).Result()
	if err == nil {
		depth, _ := strconv.Atoi(val)
		return depth, nil
	}

	// Cache miss — simulate fetching from DB and caching it
	// In production this would be a real DB query
	simulatedDepth := 3
	s.Redis.Set(ctx, key, simulatedDepth, 2*time.Minute)
	return simulatedDepth, nil
}

func (s *ETAService) CalculateETA(ctx context.Context, restaurantID string, deliveryLat, deliveryLng float64, itemCount int) (int, float64, int, float64, error) {
	// Get restaurant coordinates
	coords, exists := restaurantCoords[restaurantID]
	if !exists {
		return 0, 0, 0, 0, fmt.Errorf("restaurant %s not found", restaurantID)
	}

	// Calculate distance
	distance := haversine(coords[0], coords[1], deliveryLat, deliveryLng)

	// Base travel time: assume 20 km/h average speed in city
	travelTime := (distance / 20.0) * 60 // in minutes

	// Queue depth from Redis
	queueDepth, err := s.getQueueDepth(ctx, restaurantID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Each order in queue adds 4 minutes prep time
	prepTime := float64(queueDepth * 4)

	// Item count adds slight prep overhead
	itemOverhead := float64(itemCount) * 0.5

	// Surge factor
	surge := getSurgeFactor()

	// Final ETA
	eta := int((travelTime + prepTime + itemOverhead) * surge)

	return eta, surge, queueDepth, math.Round(distance*100) / 100, nil
}