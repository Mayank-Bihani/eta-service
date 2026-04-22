import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '30s', target: 50 },   // ramp up to 50 users
    { duration: '1m', target: 500 },   // ramp up to 500 users
    { duration: '30s', target: 500 },  // hold at 500
    { duration: '30s', target: 0 },    // ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% of requests under 500ms
    errors: ['rate<0.01'],             // error rate under 1%
  },
};

export default function () {
  const payload = JSON.stringify({
    restaurant_id: `REST00${Math.floor(Math.random() * 3) + 1}`,
    delivery_lat: 28.5 + Math.random() * 0.5,
    delivery_lng: 77.0 + Math.random() * 0.5,
    item_count: Math.floor(Math.random() * 5) + 1,
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post('http://localhost:8080/api/order/eta', payload, params);

  const success = check(res, {
    'status is 200': (r) => r.status === 200,
    'has eta': (r) => JSON.parse(r.body).estimated_eta_minutes > 0,
  });

  errorRate.add(!success);
  sleep(0.1);
}