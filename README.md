# HyperLocal Delivery ETA Service

A high-performance Go microservice that computes real-time delivery ETAs for a hyperlocal delivery platform — inspired by systems like Swiggy and Zomato.

## Architecture
POST /api/order/eta
│
▼
[Gin Router]
│
▼
[ETA Handler]
│
├──▶ [Redis] ← restaurant queue depth (cache)
│         └── cache miss → simulate DB fetch → write to Redis
├──▶ [Distance Calculator] ← Haversine formula
└──▶ [Surge Calculator] ← time-of-day multiplier
│
▼
Save order → [PostgreSQL]
│
▼
Return ETA (JSON)

## Performance (k6 Load Test — 500 concurrent users)

| Metric | Result |
|---|---|
| Throughput | 2000 RPS |
| p95 Latency | 105ms |
| Error Rate | 0% |
| Total Requests | 300,000 |

## Tech Stack

- **Language:** Go (Gin framework)
- **Database:** PostgreSQL (order persistence)
- **Cache:** Redis (restaurant queue depth)
- **Containerization:** Docker + Docker Compose
- **Load Testing:** k6

## API Endpoints

### POST /api/order/eta
Compute delivery ETA for an order.

**Request:**
```json
{
  "restaurant_id": "REST001",
  "delivery_lat": 28.7041,
  "delivery_lng": 77.1025,
  "item_count": 3
}
```

**Response:**
```json
{
  "estimated_eta_minutes": 56,
  "restaurant_id": "REST001",
  "surge_factor": 1.4,
  "queue_depth": 3,
  "distance_km": 14.44
}
```

### GET /api/order/:id
Fetch a saved order by ID.

### POST /api/restaurant/queue
Update restaurant queue depth in Redis.

**Request:**
```json
{
  "restaurant_id": "REST001",
  "depth": 8
}
```

### GET /health
Service health check.

## How ETA is Calculated
ETA = (travel_time + prep_time + item_overhead) × surge_factor
travel_time   = Haversine distance / 20 km/h average city speed
prep_time     = queue_depth × 4 minutes per order
item_overhead = item_count × 0.5 minutes
surge_factor  = 1.4 during lunch (12–2pm), 1.6 during dinner (7–9pm), 1.0 otherwise

## Running Locally

**Prerequisites:** Docker, Docker Compose

```bash
git clone git@github.com:Mayank-Bihani/eta-service.git
cd eta-service
docker-compose up --build
```

Service runs on `http://localhost:8080`

## Project Structure
eta-service/
├── main.go
├── config/        # env config loader
├── db/            # PostgreSQL + Redis connections
├── handlers/      # HTTP request handlers
├── models/        # data structs
├── services/      # core ETA business logic
└── router/        # route definitions