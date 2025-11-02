# Refactoring Status: Clean Architecture → Domain-Based Architecture

## ✅ **Yang Sudah Selesai**

### 1. **Struktur Folder Baru Dibuat**
- ✅ `internal/users/` (auth, handler, service, repository, dto)
- ✅ `internal/orders/` (handler, service, repository, dto)
- ✅ `internal/payments/` (handler, service, repository, dto, legacy_payment)
- ✅ `internal/branches/` (handler, service, repository, dto)
- ✅ `internal/shared/` (config, database, middleware, notification, monitoring, utils, model, handlers)

### 2. **File-File Sudah Dipindah**
- ✅ Users domain files (auth, handler, repository, dto)
- ✅ Orders domain files (handler, service, repository, dto)
- ✅ Payments domain files (handler, service, repository, dto, legacy_payment)
- ✅ Branches domain files (handler, service, repository, dto)
- ✅ Shared resources (config, database, middleware, notification, monitoring, utils)
- ✅ Shared models (common, dashboard, file, membership, notification, rating, report, dll)
- ✅ Shared handlers (chat, dashboard, file, health, membership, notification, rating, report, swagger, websocket)

### 3. **Import Paths Diupdate (Sebagian)**
- ✅ Python script dibuat untuk auto-update imports
- ✅ 73 file sudah diproses oleh script
- ✅ Package names sudah diupdate
- ✅ Import paths dasar sudah diupdate

### 4. **Scripts Dibuat**
- ✅ `scripts/refactor-imports.py` - Smart import updater
- ✅ `scripts/update-imports.sh` - Bash script untuk update imports

## ⚠️ **Yang Masih Perlu Dikerjakan**

### 1. **Update routes.go**
**File:** `internal/shared/handlers/routes.go`

Masih ada beberapa referensi yang perlu diperbaiki:
- Update semua `core.*` menjadi `model.*` untuk shared models
- Update `branchHandler`, `orderHandler`, `paymentHandler` references
- Import handlers dari domain packages dengan benar

**Contoh yang perlu diperbaiki:**
```go
// OLD:
admin.Use(middleware.RoleMiddleware(core.RoleAdminPusat, core.RoleAdminCabang))
public.POST("/payments/midtrans/callback", paymentHandler.MidtransCallback)
c.JSON(http.StatusNotFound, core.CreateErrorResponse(...))

// NEW:
admin.Use(middleware.RoleMiddleware(model.RoleAdminPusat, model.RoleAdminCabang))
public.POST("/payments/midtrans/callback", paymentHdlr.MidtransCallback)
c.JSON(http.StatusNotFound, model.CreateErrorResponse(...))
```

### 2. **Update cmd/app/main.go**
**File:** `cmd/app/main.go`

Import paths perlu diupdate:
```go
// OLD:
import (
    "service/internal/config"
    "service/internal/database"
    "service/internal/delivery"
    "service/internal/middleware"
    "service/internal/monitoring"
    svc "service/internal/service"
    "service/internal/utils"
)

// NEW:
import (
    "service/internal/shared/config"
    "service/internal/shared/database"
    "service/internal/shared/handlers"
    "service/internal/shared/middleware"
    "service/internal/shared/monitoring"
    "service/internal/shared/utils"
)
```

Dan update fungsi:
```go
// OLD:
delivery.SetupRoutes(r)

// NEW:
handlers.SetupRoutes(r)
```

### 3. **Fix Import References di File-File yang Sudah Dipindah**

Beberapa file masih menggunakan `core.*` yang seharusnya:
- Domain DTOs: `dto.*`
- Shared models: `model.*`

**Perlu diupdate secara manual:**
- `internal/users/auth/auth_service.go` - semua `core.*` → `dto.*` atau `model.*`
- `internal/users/handler/auth_handler.go` - update imports
- `internal/orders/service/order_service.go` - update imports
- `internal/payments/service/payment_service.go` - update imports
- Semua handler files di shared/handlers

### 4. **Update Cross-Domain Dependencies**

File yang menggunakan multiple domains perlu update imports:
- `internal/orders/service/order_service.go` mungkin perlu import `users/dto`
- `internal/payments/service/payment_service.go` mungkin perlu import `orders/dto`
- dll

### 5. **Test Build**
Setelah semua import diupdate, perlu:
1. Test build: `go build ./cmd/app/main.go`
2. Fix any remaining import errors
3. Update go.mod jika perlu

## 📋 **Checklist Penyelesaian**

### Prioritas Tinggi 🔴
- [ ] Fix `routes.go` - Update semua handler references dan core → model
- [ ] Fix `cmd/app/main.go` - Update semua imports
- [ ] Fix shared handlers - Update imports untuk model
- [ ] Fix domain handlers - Update imports untuk dto dan cross-domain

### Prioritas Sedang 🟡
- [ ] Fix domain services - Update imports untuk repositories dan dto
- [ ] Fix domain repositories - Update imports untuk dto dan database
- [ ] Fix auth service - Update imports untuk users/dto
- [ ] Test build dan fix errors

### Prioritas Rendah 🟢
- [ ] Hapus folder lama setelah semua berfungsi
- [ ] Update dokumentasi README.md
- [ ] Update migration scripts jika ada

## 🔧 **Cara Menyelesaikan**

### Step 1: Fix routes.go
```bash
cd /home/il/service-go
# Edit internal/shared/handlers/routes.go
# Replace all `core.*` with `model.*`
# Replace all handler references dengan yang baru
```

### Step 2: Fix main.go
```bash
# Edit cmd/app/main.go
# Update all imports ke shared paths
# Update delivery.SetupRoutes ke handlers.SetupRoutes
```

### Step 3: Run Build & Fix Errors
```bash
go build ./cmd/app/main.go
# Fix semua error yang muncul satu per satu
```

### Step 4: Fix Remaining Imports
Gunakan script atau update manual:
```bash
python3 scripts/refactor-imports.py
# Atau update manual untuk file-file yang masih error
```

### Step 5: Final Test
```bash
go build ./...
go test ./...  # Jika ada tests
```

## 📝 **Catatan Penting**

1. **Shared Models vs Domain DTOs:**
   - Shared models (di `shared/model`): common.go, dashboard.go, file.go, membership.go, notification.go, rating.go, report.go, dll
   - Domain DTOs (di `domain/dto`): user.go, order.go, payment.go, branch.go

2. **Cross-Domain Dependencies:**
   - Orders domain bisa import `users/dto` untuk User reference
   - Payments domain bisa import `orders/dto` untuk Order reference
   - Tapi hindari circular dependencies

3. **Package Names:**
   - Domain handlers: `package handler`
   - Domain services: `package service`
   - Domain repositories: `package repository`
   - Domain DTOs: `package dto`
   - Shared: sesuai subdirectory (config, database, middleware, dll)

## 🎯 **Target Akhir**

Setelah semua selesai:
- ✅ Struktur folder domain-based lengkap
- ✅ Semua imports diupdate dan benar
- ✅ Build berhasil tanpa error
- ✅ Aplikasi bisa berjalan normal
- ✅ Folder lama bisa dihapus (backup dulu!)

## 📚 **Referensi**

Lihat `REFACTORING-GUIDE.md` untuk mapping lengkap file-file.

---

**Last Updated:** $(date)
**Status:** 60% Complete - Need manual fixes for routes.go, main.go, and cross-domain dependencies

