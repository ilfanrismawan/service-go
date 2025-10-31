# API Usage Examples

## Membership Subscription System

### Get Available Membership Tiers
```bash
curl -X GET http://localhost:8080/api/v1/membership/tiers
```

### Start Trial Membership
```bash
curl -X POST "http://localhost:8080/api/v1/membership/trial?tier=premium" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Subscribe to Membership
```bash
curl -X POST http://localhost:8080/api/v1/membership/subscribe \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "tier": "premium",
    "subscription_type": "monthly",
    "payment_method": "gopay",
    "auto_renew": true
  }'
```

### Get Membership Details
```bash
curl -X GET http://localhost:8080/api/v1/membership \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Upgrade Membership
```bash
curl -X POST "http://localhost:8080/api/v1/membership/upgrade?tier=vip&subscription_type=yearly" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Cancel Subscription
```bash
curl -X POST http://localhost:8080/api/v1/membership/cancel \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "No longer needed"
  }'
```

### Get Membership Usage
```bash
curl -X GET http://localhost:8080/api/v1/membership/usage \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Redeem Points
```bash
curl -X POST "http://localhost:8080/api/v1/membership/redeem-points?points=1000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Membership Tier Examples

### Basic Tier (Rp 50,000/month)
- 5% discount on all services
- Basic support
- Mobile app access
- Order tracking

### Premium Tier (Rp 100,000/month)
- 10% discount on all services
- Priority support
- 1 free service per month
- 2 free pickups per month
- Free diagnostics
- 30 days extended warranty
- Priority queue
- Mobile app premium features

### VIP Tier (Rp 200,000/month)
- 15% discount on all services
- VIP support with dedicated manager
- 3 free services per month
- 5 free pickups per month
- Free diagnostics
- 60 days extended warranty
- Exclusive offers & promotions
- Priority queue & fast track
- Free device cleaning
- Concierge service

### Elite Tier (Rp 350,000/month)
- 20% discount on all services
- Elite support with personal assistant
- 5 free services per month
- 10 free pickups per month
- Free diagnostics
- 90 days extended warranty
- Exclusive offers & early access
- Highest priority queue
- Free device cleaning & maintenance
- White-glove concierge service
- Home service visits
- Exclusive events & workshops
- Lifetime warranty on repairs

## Example Responses

### Get Membership Tiers Response
```json
{
  "status": "success",
  "data": [
    {
      "tier": "basic",
      "discount_percentage": 5.0,
      "monthly_price": 50000,
      "yearly_price": 500000,
      "points_multiplier": 1.0,
      "max_free_services": 0,
      "max_free_pickups": 0,
      "priority_support": false,
      "extended_warranty_days": 0,
      "exclusive_offers": false,
      "free_diagnostics": false,
      "benefits": [
        "5% discount on all services",
        "Basic support",
        "Mobile app access",
        "Order tracking"
      ]
    }
  ],
  "message": "Membership tiers retrieved successfully"
}
```

### Subscribe to Membership Response
```json
{
  "status": "success",
  "data": {
    "membership_id": "123e4567-e89b-12d3-a456-426614174000",
    "tier": "premium",
    "subscription_type": "monthly",
    "status": "active",
    "price": 100000,
    "next_billing_date": "2024-02-01T00:00:00Z",
    "payment_url": "/payment/subscription/123e4567-e89b-12d3-a456-426614174000",
    "message": "Subscription created successfully"
  },
  "message": "Subscription created successfully"
}
```

### Get Membership Details Response
```json
{
  "status": "success",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "456e7890-e89b-12d3-a456-426614174000",
    "tier": "premium",
    "status": "active",
    "subscription_type": "monthly",
    "discount_percentage": 10.0,
    "points": 1500,
    "total_spent": 500000,
    "orders_count": 8,
    "monthly_price": 100000,
    "yearly_price": 1000000,
    "current_price": 100000,
    "joined_at": "2024-01-01T00:00:00Z",
    "expires_at": "2024-02-01T00:00:00Z",
    "next_billing_date": "2024-02-01T00:00:00Z",
    "auto_renew": true,
    "trial_ends_at": null,
    "benefits": [
      "10% discount on all services",
      "Priority support",
      "1 free service per month",
      "2 free pickups per month",
      "Free diagnostics",
      "30 days extended warranty",
      "Priority queue",
      "Mobile app premium features"
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "Membership retrieved successfully"
}
```

## Monthly Reports

### Get Current Month Report
```bash
curl -X GET http://localhost:8080/api/v1/reports/current-month \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Specific Month Report
```bash
curl -X GET "http://localhost:8080/api/v1/reports/monthly?year=2024&month=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Yearly Report
```bash
curl -X GET "http://localhost:8080/api/v1/reports/yearly?year=2024" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Payment with Midtrans Mock

### Process Payment
```bash
curl -X POST http://localhost:8080/api/v1/payments/process \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": "payment-uuid",
    "payment_method": "gopay"
  }'
```

### Midtrans Webhook Callback
```bash
curl -X POST http://localhost:8080/api/v1/payments/midtrans/callback \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "INV-20240101-000123",
    "transaction_id": "abc123",
    "status_code": "200",
    "gross_amount": "500000.00",
    "transaction_status": "settlement",
    "signature_key": "<sha512(order_id+status_code+gross_amount+server_key)>"
  }'
```

### Payment Methods Available
- `cash` - Pembayaran tunai (langsung paid)
- `bank_transfer` - Transfer bank (Virtual Account)
- `gopay` - GoPay (Redirect URL)
- `qris` - QRIS (QR String)
- `mandiri_echannel` - Mandiri E-Channel (Bill Key)

## Admin Endpoints

### Get Membership Statistics
```bash
curl -X GET http://localhost:8080/api/v1/admin/membership/stats \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### Get Top Spenders
```bash
curl -X GET "http://localhost:8080/api/v1/admin/membership/top-spenders?limit=10" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### List All Memberships
```bash
curl -X GET "http://localhost:8080/api/v1/admin/membership/list?page=1&limit=10&tier=gold" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```
