# API Documentation

Dokumentasi lengkap untuk iPhone Service POS API.

## 📋 Daftar Isi

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Base URL](#base-url)
4. [Response Format](#response-format)
5. [Error Handling](#error-handling)
6. [Endpoints](#endpoints)
7. [WebSocket](#websocket)
8. [Rate Limiting](#rate-limiting)
9. [Examples](#examples)

---

## Overview

iPhone Service POS API adalah RESTful API untuk sistem Point of Sales jasa service iPhone dengan dukungan 50 cabang di seluruh Indonesia.

### Fitur Utama

- ✅ Authentication & Authorization dengan JWT
- ✅ Multi-service support (Service Catalog)
- ✅ Order Management dengan real-time tracking
- ✅ Payment Processing (Midtrans, Cash, Bank Transfer, GoPay, QRIS)
- ✅ Membership System (4-tier: Bronze, Silver, Gold, Platinum)
- ✅ Location Tracking dengan WebSocket
- ✅ Real-time Chat dengan WebSocket
- ✅ Notification System (Email, WhatsApp, FCM)
- ✅ File Management (S3-compatible storage)
- ✅ Reporting & Analytics
- ✅ Rating System
- ✅ Dashboard untuk Admin

### Teknologi

- **Framework:** Gin (Go)
- **Database:** PostgreSQL dengan GORM
- **Cache:** Redis
- **File Storage:** S3-compatible (MinIO)
- **Payment Gateway:** Midtrans (mock untuk development)
- **Documentation:** Swagger/OpenAPI 3.0

---

## Authentication

API menggunakan JWT (JSON Web Token) untuk authentication.

### Token Types

- **Access Token:** Expires in 24 hours
- **Refresh Token:** Expires in 7 days

### Authentication Header

Setiap request ke protected endpoint harus menyertakan header:

```
Authorization: Bearer <access_token>
```

### Authentication Flow

1. **Register/Login** → Dapatkan `access_token` dan `refresh_token`
2. **Use Access Token** → Sertakan di header untuk setiap request
3. **Token Expired** → Gunakan `refresh_token` untuk mendapatkan `access_token` baru
4. **Logout** → Revoke `refresh_token`

### Endpoints

- `POST /api/v1/auth/register` - Register user baru
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout (revoke token)
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password
- `GET /api/v1/auth/profile` - Get user profile (protected)
- `PUT /api/v1/auth/profile` - Update user profile (protected)
- `POST /api/v1/auth/change-password` - Change password (protected)
- `PUT /api/v1/auth/fcm-token` - Update FCM token (protected)

---

## Base URL

- **Development:** `http://localhost:8080`
- **Production:** `https://api.iphoneservice.com`

---

## Response Format

### Success Response

```json
{
  "status": "success",
  "data": {
    // Response data
  }
  },
  "message": "Operation successful",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Error Response

```json
{
  "status": "error",
  "error": "error_code",
  "message": "Error message",
  "details": {
    // Additional error details
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Paginated Response

```json
{
  "status": "success",
  "data": [
    // Array of items
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  },
  "message": "Data retrieved successfully",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

## Error Handling

### HTTP Status Codes

- `200 OK` - Request berhasil
- `201 Created` - Resource berhasil dibuat
- `400 Bad Request` - Request tidak valid
- `401 Unauthorized` - Tidak terautentikasi
- `403 Forbidden` - Tidak memiliki akses
- `404 Not Found` - Resource tidak ditemukan
- `422 Unprocessable Entity` - Validation error
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

### Error Codes

- `validation_error` - Validation error
- `unauthorized` - Authentication required
- `forbidden` - Insufficient permissions
- `not_found` - Resource not found
- `internal_error` - Internal server error
- `rate_limit_exceeded` - Rate limit exceeded

---

## Endpoints

### Public Endpoints (No Authentication Required)

#### Authentication

- `POST /api/v1/auth/register` - Register user baru
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout (revoke token)
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password

#### Branches

- `GET /api/v1/branches` - Get all branches
- `GET /api/v1/branches/nearest` - Get nearest branches (requires lat, lng query params)
- `GET /api/v1/branches/:id` - Get branch by ID

#### Service Catalog

- `GET /api/v1/services/catalog` - Get service catalog
- `GET /api/v1/services/catalog/:id` - Get service by ID

#### Payment Callback

- `POST /api/v1/payments/midtrans/callback` - Midtrans webhook callback (signature verified)

---

### Protected Endpoints (Authentication Required)

#### User Profile

- `GET /api/v1/auth/profile` - Get user profile
- `PUT /api/v1/auth/profile` - Update user profile
- `POST /api/v1/auth/change-password` - Change password
- `PUT /api/v1/auth/fcm-token` - Update FCM token

#### Orders

- `POST /api/v1/orders` - Create new order
- `GET /api/v1/orders` - Get user orders (with pagination)
- `GET /api/v1/orders/:id` - Get order by ID
- `PUT /api/v1/orders/:id/status` - Update order status
- `PUT /api/v1/orders/:id/assign-courier` - Assign courier to order
- `PUT /api/v1/orders/:id/assign-technician` - Assign technician to order

#### Location Tracking

- `POST /api/v1/tracking/update` - Update location (for courier)
- `GET /api/v1/tracking/order/:orderId` - Get current tracking for order
- `GET /api/v1/tracking/history/:orderId` - Get tracking history for order

#### Payments

- `POST /api/v1/payments/create-invoice` - Create payment invoice
- `POST /api/v1/payments/process` - Process payment
- `GET /api/v1/payments/:id` - Get payment by ID
- `GET /api/v1/payments/order/:orderId` - Get payments by order ID

#### Notifications

- `GET /api/v1/notifications` - Get user notifications (with pagination)
- `PUT /api/v1/notifications/:id/read` - Mark notification as read
- `POST /api/v1/notifications` - Send notification (admin only)
- `POST /api/v1/notifications/order/:orderId/status` - Send order status notification
- `POST /api/v1/notifications/order/:orderId/payment` - Send payment notification

#### Files

- `POST /api/v1/files/upload` - Upload file
- `POST /api/v1/files/orders/photo` - Upload order photo
- `POST /api/v1/files/users/avatar` - Upload user avatar
- `GET /api/v1/files/url` - Get file URL
- `GET /api/v1/files/list` - List files
- `DELETE /api/v1/files/delete` - Delete file

#### Chat

- `GET /api/v1/chat/orders/:orderId` - Get chat messages for order
- `POST /api/v1/chat/orders/:orderId` - Send chat message

#### Dashboard

- `GET /api/v1/dashboard/overview` - Get dashboard overview
- `GET /api/v1/dashboard/orders` - Get order statistics
- `GET /api/v1/dashboard/revenue` - Get revenue statistics
- `GET /api/v1/dashboard/branches` - Get branch statistics

#### Membership

- `GET /api/v1/membership` - Get user membership
- `POST /api/v1/membership` - Create membership
- `PUT /api/v1/membership` - Update membership
- `POST /api/v1/membership/redeem-points` - Redeem points
- `POST /api/v1/membership/subscribe` - Subscribe to membership
- `POST /api/v1/membership/cancel` - Cancel membership subscription
- `POST /api/v1/membership/trial` - Start trial membership
- `GET /api/v1/membership/tiers` - Get membership tiers
- `POST /api/v1/membership/upgrade` - Upgrade membership
- `GET /api/v1/membership/usage` - Get membership usage

#### Reports

- `GET /api/v1/reports/current-month` - Get current month report
- `GET /api/v1/reports/monthly` - Get monthly report (with month, year query params)
- `GET /api/v1/reports/yearly` - Get yearly report (with year query param)
- `GET /api/v1/reports/summary` - Get report summary

#### Ratings

- `POST /api/v1/ratings` - Create rating
- `GET /api/v1/ratings` - Get ratings (with pagination)
- `GET /api/v1/ratings/average` - Get average rating
- `GET /api/v1/ratings/:id` - Get rating by ID
- `PUT /api/v1/ratings/:id` - Update rating
- `DELETE /api/v1/ratings/:id` - Delete rating

---

### Admin Endpoints (Admin Role Required)

#### Branch Management

- `POST /api/v1/admin/branches` - Create branch
- `PUT /api/v1/admin/branches/:id` - Update branch
- `DELETE /api/v1/admin/branches/:id` - Delete branch
- `GET /api/v1/admin/branches` - Get all branches (admin view)

#### User Management

- `GET /api/v1/admin/users` - Get all users (with pagination)
- `GET /api/v1/admin/users/:id` - Get user by ID
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user

#### Order Management

- `GET /api/v1/admin/orders` - Get all orders (with pagination and filters)
- `PUT /api/v1/admin/orders/:id` - Update order
- `DELETE /api/v1/admin/orders/:id` - Delete order

#### Payment Management

- `GET /api/v1/admin/payments` - Get all payments (with pagination and filters)
- `PUT /api/v1/admin/payments/:id` - Update payment

#### Dashboard

- `GET /api/v1/admin/dashboard` - Get admin dashboard

#### Membership Management

- `GET /api/v1/admin/membership/list` - Get all memberships (with pagination)
- `GET /api/v1/admin/membership/stats` - Get membership statistics
- `GET /api/v1/admin/membership/top-spenders` - Get top spenders

#### Service Catalog

- Protected routes untuk manage service catalog (CRUD operations)

---

### Cashier Endpoints (Cashier Role Required)

#### Orders

- `GET /api/v1/cashier/orders` - Get cashier orders
- `PUT /api/v1/cashier/orders/:id/status` - Update order status
- `POST /api/v1/cashier/orders/:id/payment` - Process payment for order

#### Branch Orders

- `GET /api/v1/cashier/branches/:id/orders` - Get orders for specific branch

---

### Technician Endpoints (Technician Role Required)

#### Orders

- `GET /api/v1/technician/orders` - Get technician orders
- `PUT /api/v1/technician/orders/:id/status` - Update order status
- `POST /api/v1/technician/orders/:id/photo` - Upload order photo

#### Chat

- `GET /api/v1/technician/chat/orders/:orderId` - Get chat messages for order
- `POST /api/v1/technician/chat/orders/:orderId` - Send chat message

---

### Courier Endpoints (Courier Role Required)

#### Orders

- `GET /api/v1/courier/orders` - Get courier orders
- `PUT /api/v1/courier/orders/:id/status` - Update order status
- `POST /api/v1/courier/orders/:id/photo` - Upload order photo

#### Jobs

- `GET /api/v1/courier/jobs` - Get available jobs
- `POST /api/v1/courier/jobs/:id/accept` - Accept job

---

## WebSocket

### Chat WebSocket

**Endpoint:** `ws://localhost:8080/ws/chat`

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/chat?token=<access_token>');
```

**Message Format:**
```json
{
  "type": "message",
  "order_id": "order-uuid",
  "message": "Hello",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

**Event Types:**
- `message` - New chat message
- `location_update` - Location update (for tracking)
- `order_status_update` - Order status update
- `payment_update` - Payment status update

---

## Rate Limiting

API menggunakan rate limiting untuk mencegah abuse.

### Default Limits

- **Rate Limit:** 100 requests per minute
- **Window:** 1 minute

### Rate Limit Headers

Response headers:
- `X-RateLimit-Limit` - Maximum requests allowed
- `X-RateLimit-Remaining` - Remaining requests in window
- `X-RateLimit-Reset` - Time when rate limit resets

### Rate Limit Exceeded

When rate limit is exceeded, API returns:
- **Status Code:** `429 Too Many Requests`
- **Response:**
```json
{
  "status": "error",
  "error": "rate_limit_exceeded",
  "message": "Rate limit exceeded. Please try again later.",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

## Examples

### Example: Register User

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123",
    "full_name": "John Doe",
    "phone": "081234567890",
    "role": "pelanggan"
  }'
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "customer@example.com",
    "full_name": "John Doe",
    "phone": "081234567890",
    "role": "pelanggan",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "User registered successfully",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Example: Login

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "customer@example.com",
      "full_name": "John Doe",
      "role": "pelanggan"
    },
    "expires_in": 86400
  },
  "message": "Login successful",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Example: Create Order (Protected)

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "service_id": "service-uuid",
    "branch_id": "branch-uuid",
    "customer_name": "John Doe",
    "customer_phone": "081234567890",
    "customer_address": "Jl. Example No. 123",
    "customer_latitude": -6.2088,
    "customer_longitude": 106.8456,
    "description": "iPhone screen repair",
    "pickup_type": "on_demand"
  }'
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid",
    "order_number": "ORD-20240101-001",
    "service_id": "service-uuid",
    "branch_id": "branch-uuid",
    "customer_id": "customer-uuid",
    "status": "pending",
    "total_amount": 500000,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "Order created successfully",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Example: Get Orders (Protected)

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/orders?page=1&limit=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "order-uuid",
      "order_number": "ORD-20240101-001",
      "status": "in_progress",
      "total_amount": 500000,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  },
  "message": "Orders retrieved successfully",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

## Additional Resources

- **Swagger UI:** http://localhost:8080/swagger/index.html
- **API Docs:** http://localhost:8080/docs
- **API Info:** http://localhost:8080/api-docs
- **Health Check:** http://localhost:8080/health
- **Metrics:** http://localhost:8080/metrics

---

## Support

Untuk pertanyaan atau bantuan:
- **Email:** support@iphoneservice.com
- **Website:** https://iphoneservice.com
- **Documentation:** https://docs.iphoneservice.com

---

**Last Updated:** 2024-01-01

