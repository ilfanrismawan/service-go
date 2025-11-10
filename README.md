# iPhone Service API

Backend aplikasi POS (Point of Sales) untuk bisnis jasa service iPhone menggunakan bahasa Go (Golang) dengan Clean Architecture.

## 🎯 Deskripsi Umum

Aplikasi ini digunakan oleh perusahaan jasa service iPhone dengan 50 cabang di seluruh Indonesia. Setiap cabang memiliki kasir, teknisi, dan kurir untuk layanan antar-jemput iPhone. Pelanggan dapat melakukan order service secara online maupun langsung ke cabang. Sistem mendukung pembayaran tunai (cash) dan online (Midtrans, Gopay, QRIS, transfer bank).

## ⚙️ Spesifikasi Teknis

- **Bahasa:** Go (Golang)
- **Framework:** Gin
- **Database:** PostgreSQL
- **ORM:** GORM
- **Auth:** JWT (Access + Refresh Token)
- **Cache:** Redis
- **Payment Gateway:** Midtrans API (mock integrasi)
- **File Storage:** S3-compatible (MinIO)
- **Dokumentasi API:** Swagger (OpenAPI 3.0)
- **Deployment:** Docker + Docker Compose
- **Struktur Folder:** Clean Architecture / Hexagonal Architecture
- **Response Format:** JSON (RESTful API)

## 🧱 Struktur Folder (Clean Architecture)

```
/cmd/
├── app/main.go                  # Entry point aplikasi
├── migrate/main.go              # Database migration tool
├── seed/main.go                 # Database seeding tool
└── test-connections/main.go    # Connection testing tool

/internal/
├── modules/                     # Business modules (domain-driven)
│   ├── admin/                   # Admin dashboard
│   ├── branches/                # Branch management
│   ├── chat/                    # Chat system
│   ├── inventory/               # Inventory management
│   ├── media/                   # File upload & management
│   ├── membership/              # Membership system
│   ├── notification/           # Notification service
│   ├── orders/                  # Order management
│   ├── payments/                # Payment processing
│   ├── services/                # Service catalog
│   ├── tracking/                # Location tracking
│   └── users/                   # User & authentication
├── router/                      # API routing
└── shared/                      # Shared components
    ├── config/                  # Configuration
    ├── database/                # Database connection
    ├── handlers/                # Shared handlers (health, swagger, websocket)
    ├── middleware/              # Middleware (auth, CORS, logging, etc.)
    ├── model/                   # Shared models
    ├── monitoring/              # Metrics & monitoring
    └── utils/                   # Utility functions

/docs/                           # Swagger documentation (generated)
/migrations/                     # SQL migration files
/seed/                           # Database seed files
/scripts/                        # Utility scripts
/k8s/                            # Kubernetes deployment files
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make (optional, untuk menggunakan Makefile)
- swag (untuk generate Swagger docs) - akan diinstall otomatis oleh Makefile

### 1. Clone Repository

```bash
git clone <repository-url>
cd service-go
```

### 2. Setup Environment

```bash
# Copy environment template
cp env.example .env

# Edit .env file sesuai kebutuhan
nano .env
```

### 3. Start Services

```bash
# Menggunakan Makefile (recommended)
make start

# Atau manual
docker-compose up -d
```

### 4. Run Database Migrations

```bash
# Menggunakan Makefile
make migrate

# Atau manual
go run cmd/migrate/main.go
```

### 5. Seed Database (Optional)

```bash
# Menggunakan Makefile
make seed

# Atau manual
go run cmd/seed/main.go
```

### 6. Generate Swagger Documentation

```bash
# Menggunakan Makefile
make generate-swagger

# Atau manual
swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal
```

### 7. Run Application

```bash
# Menggunakan Makefile
make run

# Atau manual
go run cmd/app/main.go
```

### 8. Test API

```bash
# Health check
curl http://localhost:8080/health

# Test API endpoints
make test-api

# View API documentation
# Buka browser: http://localhost:8080/swagger/index.html
```

## 📋 Available Commands

```bash
make help          # Show all available commands
make build         # Build application
make run           # Run application locally
make test          # Run tests
make docker-up     # Start all services
make docker-down   # Stop all services
make clean         # Clean build artifacts
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:password@localhost:5432/iphone_service?sslmode=disable` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379/0` |
| `JWT_SECRET` | JWT secret key | `your-secret-key-change-this-in-production` |
| `JWT_EXPIRY` | JWT token expiry | `24h` |
| `REFRESH_EXPIRY` | Refresh token expiry | `168h` |
| `MIDTRANS_SERVER_KEY` | Midtrans server key | - |
| `MIDTRANS_CLIENT_KEY` | Midtrans client key | - |
| `MIDTRANS_IS_PRODUCTION` | Midtrans production mode | `false` |
| `S3_ENDPOINT` | S3-compatible storage endpoint | `http://localhost:9000` |
| `S3_ACCESS_KEY` | S3 access key | `minioadmin` |
| `S3_SECRET_KEY` | S3 secret key | `minioadmin` |
| `S3_BUCKET_NAME` | S3 bucket name | `iphone-service` |
| `S3_REGION` | S3 region | `us-east-1` |
| `FIREBASE_SERVER_KEY` | Firebase server key untuk FCM | - |
| `WHATSAPP_API_KEY` | WhatsApp API key | - |
| `GOOGLE_MAPS_API_KEY` | Google Maps API key untuk real-time tracking | - |
| `SENTRY_DSN` | Sentry DSN untuk error tracking | - |
| `RECONCILE_INTERVAL` | Payment reconciliation interval | `5m` |
| `RATE_LIMIT_REQUESTS` | Rate limit requests per window | `100` |
| `RATE_LIMIT_WINDOW` | Rate limit time window | `1m` |
| `TIMEZONE` | Application timezone | `Asia/Jakarta` |
| `DEFAULT_LANGUAGE` | Default language | `id-ID` |
| `CURRENCY` | Default currency | `IDR` |

### Database Configuration

Aplikasi menggunakan PostgreSQL dengan GORM untuk ORM. Database akan otomatis di-migrate saat startup.

### Redis Configuration

Redis digunakan untuk caching dan session management.

### File Storage

MinIO digunakan sebagai S3-compatible storage untuk menyimpan foto-foto service.

### Google Maps API Configuration

Sistem real-time tracking terintegrasi dengan Google Maps API untuk mendapatkan jarak dan waktu tempuh yang akurat berdasarkan rute jalan dan kondisi lalu lintas.

#### Setup Google Maps API

1. **Dapatkan API Key:**
   - Buka [Google Cloud Console](https://console.cloud.google.com/)
   - Buat project baru atau pilih project yang ada
   - Aktifkan **Distance Matrix API**
   - Buat API key di "Credentials"
   - Batasi API key untuk keamanan (opsional tapi direkomendasikan)

2. **Set Environment Variable:**
   ```bash
   export GOOGLE_MAPS_API_KEY=your-actual-api-key-here
   ```
   Atau tambahkan ke file `.env`:
   ```
   GOOGLE_MAPS_API_KEY=your-actual-api-key-here
   ```

3. **Fitur:**
   - ✅ Jarak berdasarkan rute jalan (bukan garis lurus)
   - ✅ Waktu tempuh aktual dengan mempertimbangkan lalu lintas
   - ✅ Fallback otomatis ke perhitungan Haversine jika API tidak tersedia
   - ✅ Real-time traffic data untuk ETA yang lebih akurat

**Catatan:** Jika `GOOGLE_MAPS_API_KEY` tidak di-set, sistem akan otomatis menggunakan perhitungan Haversine (jarak garis lurus) sebagai fallback.

## 🔒 Authentication & Authorization

### JWT Authentication

- **Access Token:** Expires in 24 hours
- **Refresh Token:** Expires in 7 days
- **Algorithm:** HS256

### User Roles

- `admin_pusat`: Admin pusat (full access)
- `admin_cabang`: Admin cabang (branch-specific access)
- `kasir`: Kasir (payment and order management)
- `teknisi`: Teknisi (service order management)
- `kurir`: Kurir (pickup and delivery)
- `pelanggan`: Pelanggan (customer access)

### API Endpoints

#### Authentication

```bash
POST /api/v1/auth/register          # Register new user
POST /api/v1/auth/login              # Login user
POST /api/v1/auth/refresh            # Refresh access token
POST /api/v1/auth/forgot-password    # Forgot password
POST /api/v1/auth/reset-password     # Reset password
GET  /api/v1/auth/profile            # Get user profile (protected)
PUT  /api/v1/auth/profile            # Update user profile (protected)
POST /api/v1/auth/change-password    # Change password (protected)
```

#### Example Request/Response

**Register User:**
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

**Login:**
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
  "phone": "081234567890",
      "role": "pelanggan",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "expires_in": 86400
  },
  "message": "Login successful",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 🧩 Fitur Utama

### ✅ Implemented

- [x] **Authentication System**
  - JWT-based authentication
  - Role-based access control (RBAC)
  - Password hashing dengan bcrypt
  - Refresh token mechanism

- [x] **User Management**
  - User registration dan login
  - Profile management
  - Password reset functionality

- [x] **Clean Architecture**
  - Domain models dan entities
  - Repository pattern
  - Service layer
  - Delivery layer (HTTP handlers)

- [x] **Database Integration**
  - PostgreSQL dengan GORM
  - Auto migration
  - Redis untuk caching

- [x] **Docker Support**
  - Docker Compose setup
  - PostgreSQL, Redis, MinIO
  - Development environment

- [x] **Branch Management**
  - CRUD operations untuk cabang
  - Pencarian cabang terdekat berdasarkan geolocation
  - Branch statistics dan management

- [x] **Order Service System**
  - Order creation dan management
  - Status tracking dengan real-time updates
  - Photo upload untuk service
  - Order assignment (courier & technician)
  - Order reports dan analytics

- [x] **Payment System**
  - Midtrans integration (mock untuk development)
  - Multiple payment methods (Bank Transfer, GoPay, QRIS, Cash)
  - Invoice generation
  - Payment reconciliation
  - Payment status tracking

- [x] **Notification System**
  - Email notifications
  - WhatsApp integration (template-based)
  - Push notifications (FCM)
  - Order status notifications
  - Payment notifications

- [x] **Service Catalog System**
  - Dynamic service catalog management
  - Service categories dan pricing
  - Service availability management
  - Service metadata support

- [x] **Location Tracking**
  - Real-time location tracking
  - Courier location updates
  - WebSocket support untuk real-time updates
  - Location history

- [x] **Chat System**
  - Real-time chat dengan WebSocket
  - Chat per order
  - Message history

- [x] **Membership System**
  - 4-tier membership (Bronze, Silver, Gold, Platinum)
  - Points system
  - Auto-upgrade berdasarkan spending
  - Membership benefits

- [x] **Reporting System**
  - Monthly reports
  - Yearly reports
  - Revenue analytics
  - Branch performance reports

- [x] **Rating System**
  - Order ratings
  - Service ratings
  - Average rating calculation

- [x] **File Management**
  - File upload ke S3-compatible storage (MinIO)
  - Image compression
  - Order photos
  - User avatars

- [x] **Dashboard**
  - Admin dashboard
  - Overview statistics
  - Order statistics
  - Revenue statistics
  - Branch statistics

- [x] **Monitoring & Metrics**
  - Prometheus metrics
  - Health checks (liveness & readiness)
  - Request logging
  - Performance monitoring
  - Sentry integration untuk error tracking

## 🐳 Docker Services

| Service | Port | Description |
|---------|------|-------------|
| **app** | 8080 | Main application |
| **postgres** | 5432 | PostgreSQL database |
| **redis** | 6379 | Redis cache |
| **minio** | 9000, 9001 | S3-compatible storage |

## 📊 API Response Format

### Success Response
```json
{
  "status": "success",
  "data": { ... },
  "message": "Operation successful",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Error Response
```json
{
  "status": "error",
  "error": "validation_error",
  "message": "Validation failed",
  "details": { ... },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Paginated Response
```json
{
  "status": "success",
  "data": [ ... ],
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

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Test API endpoints
make test-api
```

## 🔧 Development

### Code Quality

```bash
# Format code
make format

# Run linter
make lint
```

### Database Management

```bash
# Run migrations
make migrate

# Seed database
make seed
```

## 🚀 Deployment

### Production Build

   ```bash
# Build for production
make prod-build
```

### Environment Setup

1. Copy `env.example` to `.env`
2. Update production values
3. Set `ENVIRONMENT=production`
4. Update `JWT_SECRET` dengan nilai yang aman
5. Configure database dan Redis URLs

## 📚 API Documentation

### Response Format

API menggunakan format JSON dengan struktur response standar:

**Success Response:**
```json
{
  "status": "success",
  "data": {...},
  "message": "Operation completed successfully",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

**Error Response:**
```json
{
  "status": "error",
  "error": "error_code",
  "message": "Error message",
  "details": {...},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

**Paginated Response:**
```json
{
  "status": "success",
  "data": [...],
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

### Swagger Documentation

API documentation tersedia melalui Swagger UI:

- **Swagger UI:** http://localhost:8080/swagger/index.html
- **API Docs (Redirect):** http://localhost:8080/docs
- **API Info:** http://localhost:8080/api-docs

Untuk generate dokumentasi Swagger:

```bash
# Menggunakan Makefile
make generate-swagger

# Atau manual
swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal
```

**Catatan:** Swagger documentation di-generate dari annotations di code. Pastikan untuk mengupdate annotations setelah menambah atau mengubah endpoint.

### Endpoint Utama

#### Public Endpoints (No Authentication)
- **Authentication:** `/api/v1/auth/register`, `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/auth/logout`, `/api/v1/auth/forgot-password`, `/api/v1/auth/reset-password`
- **Branches:** `GET /api/v1/branches`, `GET /api/v1/branches/nearest`, `GET /api/v1/branches/:id`
- **Service Catalog:** `GET /api/v1/services/catalog`, `GET /api/v1/services/catalog/:id`
- **Payment Callback:** `POST /api/v1/payments/midtrans/callback` (webhook Midtrans, signature verified)

#### Protected Endpoints (Authentication Required)
- **User Profile:** `GET /api/v1/auth/profile`, `PUT /api/v1/auth/profile`, `POST /api/v1/auth/change-password`, `PUT /api/v1/auth/fcm-token`
- **Orders:** `POST /api/v1/orders`, `GET /api/v1/orders`, `GET /api/v1/orders/:id`, `PUT /api/v1/orders/:id/status`, `PUT /api/v1/orders/:id/assign-courier`, `PUT /api/v1/orders/:id/assign-technician`
- **Location Tracking:** `POST /api/v1/orders/:id/location`, `GET /api/v1/orders/:id/location`, `GET /api/v1/orders/:id/location/history`, `POST /api/v1/orders/:id/eta`, `GET /api/v1/orders/:id/eta` (dengan integrasi Google Maps API)
- **Payments:** `POST /api/v1/payments/create-invoice`, `POST /api/v1/payments/process`, `GET /api/v1/payments/:id`, `GET /api/v1/payments/order/:orderId`
- **Notifications:** `GET /api/v1/notifications`, `PUT /api/v1/notifications/:id/read`, `POST /api/v1/notifications`
- **Files:** `POST /api/v1/files/upload`, `POST /api/v1/files/orders/photo`, `POST /api/v1/files/users/avatar`, `GET /api/v1/files/url`, `GET /api/v1/files/list`, `DELETE /api/v1/files/delete`
- **Chat:** `GET /api/v1/chat/orders/:orderId`, `POST /api/v1/chat/orders/:orderId`
- **Dashboard:** `GET /api/v1/dashboard/overview`, `GET /api/v1/dashboard/orders`, `GET /api/v1/dashboard/revenue`, `GET /api/v1/dashboard/branches`
- **Membership:** `GET /api/v1/membership`, `POST /api/v1/membership`, `PUT /api/v1/membership`, `POST /api/v1/membership/redeem-points`, `POST /api/v1/membership/subscribe`, `POST /api/v1/membership/cancel`, `POST /api/v1/membership/trial`, `GET /api/v1/membership/tiers`, `POST /api/v1/membership/upgrade`, `GET /api/v1/membership/usage`
- **Reports:** `GET /api/v1/reports/current-month`, `GET /api/v1/reports/monthly`, `GET /api/v1/reports/yearly`, `GET /api/v1/reports/summary`
- **Ratings:** `POST /api/v1/ratings`, `GET /api/v1/ratings`, `GET /api/v1/ratings/average`, `GET /api/v1/ratings/:id`, `PUT /api/v1/ratings/:id`, `DELETE /api/v1/ratings/:id`

#### Admin Endpoints (Admin Role Required)
- **Branch Management:** `POST /api/v1/admin/branches`, `PUT /api/v1/admin/branches/:id`, `DELETE /api/v1/admin/branches/:id`, `GET /api/v1/admin/branches`
- **User Management:** `GET /api/v1/admin/users`, `GET /api/v1/admin/users/:id`, `PUT /api/v1/admin/users/:id`, `DELETE /api/v1/admin/users/:id`
- **Order Management:** `GET /api/v1/admin/orders`, `PUT /api/v1/admin/orders/:id`, `DELETE /api/v1/admin/orders/:id`
- **Payment Management:** `GET /api/v1/admin/payments`, `PUT /api/v1/admin/payments/:id`
- **Dashboard:** `GET /api/v1/admin/dashboard`
- **Membership Management:** `GET /api/v1/admin/membership/list`, `GET /api/v1/admin/membership/stats`, `GET /api/v1/admin/membership/top-spenders`
- **Service Catalog:** Protected routes untuk manage service catalog

#### Cashier Endpoints (Cashier Role Required)
- **Orders:** `GET /api/v1/cashier/orders`, `PUT /api/v1/cashier/orders/:id/status`, `POST /api/v1/cashier/orders/:id/payment`
- **Branch Orders:** `GET /api/v1/cashier/branches/:id/orders`

#### Technician Endpoints (Technician Role Required)
- **Orders:** `GET /api/v1/technician/orders`, `PUT /api/v1/technician/orders/:id/status`, `POST /api/v1/technician/orders/:id/photo`
- **Chat:** `GET /api/v1/technician/chat/orders/:orderId`, `POST /api/v1/technician/chat/orders/:orderId`

#### Courier Endpoints (Courier Role Required)
- **Orders:** `GET /api/v1/courier/orders`, `PUT /api/v1/courier/orders/:id/status`, `POST /api/v1/courier/orders/:id/photo`
- **Jobs:** `GET /api/v1/courier/jobs`, `POST /api/v1/courier/jobs/:id/accept`

#### WebSocket
- **Chat WebSocket:** `GET /ws/chat`

#### System Endpoints
- **Health Check:** `GET /health`, `GET /health/live`, `GET /health/ready`
- **Metrics:** `GET /metrics` (Prometheus)
- **API Documentation:** `GET /swagger/index.html`, `GET /docs`, `GET /api-docs`

### Fitur Baru

#### 🎖️ Sistem Membership
- **4 Tier Membership:** Bronze (5%), Silver (10%), Gold (15%), Platinum (20%)
- **Sistem Poin:** Earn points dari setiap transaksi
- **Auto Upgrade:** Otomatis upgrade tier berdasarkan spending dan jumlah order
- **Redeem Points:** Tukar poin untuk diskon
- **Benefits:** Berbagai benefit sesuai tier (priority support, free pickup, dll)

#### 📊 Laporan Bulanan
- **Monthly Report:** Laporan komprehensif per bulan
- **Yearly Report:** Laporan tahunan dengan data bulanan
- **Growth Metrics:** Perbandingan dengan bulan sebelumnya
- **Analytics:** Revenue by branch, payment methods, service types
- **Top Performers:** Top services dan branches

#### 💳 Integrasi Payment Midtrans (Mock)
- **Multiple Payment Methods:** Bank Transfer, GoPay, QRIS, Mandiri E-Channel
- **Cash Payment:** Support pembayaran tunai
- **Mock Implementation:** Simulasi lengkap untuk development
- **Payment Status Tracking:** Real-time status update
 - **Webhook Security:** Verifikasi `signature_key = sha512(order_id+status_code+gross_amount+server_key)`

#### 🔐 Keamanan Token
- **Refresh Rotation & Revoke:** Refresh token di-rotasi saat refresh; token lama di-blacklist (Redis). Endpoint `POST /api/v1/auth/logout` untuk revoke manual.
- **Token Blacklist:** Menggunakan Redis untuk menyimpan blacklisted tokens
- **JWT Validation:** Middleware untuk validasi JWT token di setiap protected endpoint

#### 🧪 CI Coverage Gate
- **Coverage Minimal:** Build CI gagal bila coverage < 75% (workflow CI/CD dan API Testing)

#### 🔄 WebSocket Support
- **Real-time Chat:** WebSocket endpoint untuk chat real-time
- **Real-time Tracking:** WebSocket untuk update lokasi kurir secara real-time
- **Connection Management:** Automatic reconnection dan connection pooling

#### 📊 Monitoring & Observability
- **Prometheus Metrics:** Expose metrics di `/metrics` endpoint
- **Health Checks:** Liveness dan readiness probes untuk Kubernetes
- **Structured Logging:** Logging dengan logrus untuk observability
- **Sentry Integration:** Error tracking dengan Sentry (optional)
- **Request ID:** Unique request ID untuk tracing
- **Performance Logging:** Logging untuk performance monitoring

#### 🛡️ Security Features
- **CORS:** Configurable CORS middleware
- **Rate Limiting:** Rate limiting untuk mencegah abuse
- **Security Headers:** Security headers middleware
- **HTTPS Enforcement:** HTTPS redirect di production
- **Input Validation:** Request validation dengan validator
- **SQL Injection Protection:** Menggunakan GORM untuk prepared statements

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

Jika mengalami masalah atau memiliki pertanyaan:

1. Check [Issues](../../issues) untuk masalah yang sudah diketahui
2. Create new issue dengan detail yang jelas
3. Untuk pertanyaan umum, gunakan [Discussions](../../discussions)

## 📞 Contact

- **Email:** support@iphoneservice.com
- **Website:** https://iphoneservice.com
- **Documentation:** https://docs.iphoneservice.com

---

**Happy Coding! 🚀**