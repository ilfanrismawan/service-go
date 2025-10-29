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
/cmd/app/main.go                 # Entry point aplikasi
/internal/
├── config/                      # Konfigurasi aplikasi
├── core/                        # Domain entities dan business logic
├── service/                     # Business logic layer
├── repository/                   # Data access layer
├── delivery/                    # API handlers (HTTP)
├── middleware/                  # Middleware (auth, CORS, etc.)
├── auth/                        # Authentication service
├── notification/                 # Notification service
├── payment/                     # Payment service
└── utils/                       # Utility functions
/docs/                           # Swagger documentation
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional, untuk menggunakan Makefile)

### 1. Clone Repository

```bash
git clone <repository-url>
cd service
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

### 4. Run Application

```bash
# Menggunakan Makefile
make run

# Atau manual
go run cmd/app/main.go
```

### 5. Test API

```bash
# Health check
curl http://localhost:8080/health

# Test API endpoints
make test-api
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

### Database Configuration

Aplikasi menggunakan PostgreSQL dengan GORM untuk ORM. Database akan otomatis di-migrate saat startup.

### Redis Configuration

Redis digunakan untuk caching dan session management.

### File Storage

MinIO digunakan sebagai S3-compatible storage untuk menyimpan foto-foto service.

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

### 🚧 In Progress

- [ ] **Branch Management**
  - CRUD operations untuk cabang
  - Pencarian cabang terdekat
  - Geolocation support

- [ ] **Order Service System**
  - Order creation dan management
  - Status tracking
  - Photo upload untuk service

- [ ] **Payment System**
  - Midtrans integration
  - Multiple payment methods
  - Invoice generation

- [ ] **Notification System**
  - Email notifications
  - WhatsApp integration
  - Push notifications

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

API menggunakan format JSON dengan struktur response standar:

```json
{
  "status": "success",
  "data": {...},
  "message": "Operation completed successfully"
}
```

### Swagger Documentation

API documentation tersedia melalui Swagger UI:

- **Swagger UI:** http://localhost:8080/swagger/index.html
- **API Docs:** http://localhost:8080/docs
- **API Info:** http://localhost:8080/api-docs

Untuk generate dokumentasi Swagger:

   ```bash
# Menggunakan Makefile
make generate-swagger

# Atau manual
swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal
```

### Endpoint Utama

- **Authentication:** `/api/v1/auth/*`
- **Branches:** `/api/v1/branches/*`
- **Orders:** `/api/v1/orders/*`
- **Payments:** `/api/v1/payments/*`
- **Notifications:** `/api/v1/notifications/*`
- **Files:** `/api/v1/files/*`
- **Chat:** `/api/v1/chat/*`
- **Dashboard:** `/api/v1/dashboard/*`
- **Membership:** `/api/v1/membership/*` (NEW)
- **Reports:** `/api/v1/reports/*` (NEW)

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