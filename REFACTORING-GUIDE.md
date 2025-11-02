# Refactoring Guide: Clean Architecture → Domain-Based Architecture

## 📋 Mapping File Struktur Baru

### Users Domain
- `internal/auth/auth_service.go` → `internal/users/auth/auth_service.go`
- `internal/delivery/auth_handler.go` → `internal/users/handler/auth_handler.go`
- `internal/repository/user_repository.go` → `internal/users/repository/user_repository.go`
- `internal/core/user.go` → `internal/users/dto/user.go`

### Orders Domain
- `internal/delivery/order_handler.go` → `internal/orders/handler/order_handler.go`
- `internal/service/order_service.go` → `internal/orders/service/order_service.go`
- `internal/repository/service_order_repository.go` → `internal/orders/repository/order_repository.go`
- `internal/core/service_order.go` → `internal/orders/dto/order.go`

### Payments Domain
- `internal/delivery/payment_handler.go` → `internal/payments/handler/payment_handler.go`
- `internal/service/payment_service.go` → `internal/payments/service/payment_service.go`
- `internal/repository/payment_repository.go` → `internal/payments/repository/payment_repository.go`
- `internal/payment/payment_service.go` → `internal/payments/legacy_payment/payment_service.go`
- `internal/core/payment.go` → `internal/payments/dto/payment.go`

### Branches Domain
- `internal/delivery/branch_handler.go` → `internal/branches/handler/branch_handler.go`
- `internal/service/branch_service.go` → `internal/branches/service/branch_service.go`
- `internal/repository/branch_repository.go` → `internal/branches/repository/branch_repository.go`
- `internal/core/branch.go` → `internal/branches/dto/branch.go`

### Shared Resources
- `internal/config/*` → `internal/shared/config/*`
- `internal/database/*` → `internal/shared/database/*`
- `internal/middleware/*` → `internal/shared/middleware/*`
- `internal/notification/*` → `internal/shared/notification/*`
- `internal/monitoring/*` → `internal/shared/monitoring/*`
- `internal/utils/*` → `internal/shared/utils/*`
- `internal/core/*` (kecuali yang sudah dipindah) → `internal/shared/model/*`
- `internal/delivery/*` (handler cross-domain) → `internal/shared/handlers/*`

## 🔄 Import Path Updates Needed

### Package Name Changes:
1. `service/internal/core` → 
   - `service/internal/users/dto` (untuk user)
   - `service/internal/orders/dto` (untuk order)
   - `service/internal/payments/dto` (untuk payment)
   - `service/internal/branches/dto` (untuk branch)
   - `service/internal/shared/model` (untuk shared models)

2. `service/internal/repository` → 
   - `service/internal/users/repository`
   - `service/internal/orders/repository`
   - `service/internal/payments/repository`
   - `service/internal/branches/repository`

3. `service/internal/service` →
   - `service/internal/users/service` (jika ada)
   - `service/internal/orders/service`
   - `service/internal/payments/service`
   - `service/internal/branches/service`

4. `service/internal/delivery` → 
   - `service/internal/users/handler`
   - `service/internal/orders/handler`
   - `service/internal/payments/handler`
   - `service/internal/branches/handler`
   - `service/internal/shared/handlers`

5. `service/internal/auth` → `service/internal/users/auth`
6. `service/internal/config` → `service/internal/shared/config`
7. `service/internal/database` → `service/internal/shared/database`
8. `service/internal/middleware` → `service/internal/shared/middleware`
9. `service/internal/utils` → `service/internal/shared/utils`
10. `service/internal/payment` → `service/internal/payments/legacy_payment`

## ✅ Next Steps

1. Update all import paths in moved files
2. Update package declarations
3. Create service layer untuk users jika diperlukan
4. Update routes.go dengan import paths baru
5. Update cmd/app/main.go dengan import paths baru
6. Test build setelah semua import diupdate
7. Hapus folder lama setelah semua sudah dipindah

