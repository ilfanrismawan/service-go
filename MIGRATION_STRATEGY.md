# Migration Strategy untuk Multi-Service Refactoring

## 📋 Overview

Dokumen ini menjelaskan strategi migrasi untuk transformasi aplikasi dari single-service (iPhone repair) ke multi-service (Gojek-like).

## 🔄 Perubahan Database

### 1. Tabel Baru
- `service_categories` - Kategori layanan
- `service_catalogs` - Katalog layanan
- `service_providers` - Provider/merchant
- `provider_services` - Relasi many-to-many provider dan service

### 2. Perubahan Tabel `service_orders`

#### Field Baru (Nullable)
- `service_catalog_id` (UUID, nullable) - Reference ke ServiceCatalog
- `service_provider_id` (UUID, nullable) - Reference ke ServiceProvider
- `service_name` (VARCHAR) - Nama layanan dari catalog
- `item_model`, `item_color`, `item_serial`, `item_type` - Generic item fields
- `appointment_date`, `appointment_time` - Untuk layanan yang butuh appointment
- `service_location` (TEXT) - Lokasi layanan dilakukan
- `metadata` (JSONB) - Metadata spesifik per layanan

#### Field yang Diubah (Made Optional)
- `branch_id` - Dari NOT NULL menjadi nullable (karena bisa dari provider)
- `iphone_model`, `iphone_color`, `iphone_imei`, `iphone_type` - Dari NOT NULL menjadi nullable
- `pickup_address`, `pickup_location`, `pickup_latitude`, `pickup_longitude` - Dari NOT NULL menjadi nullable
- `complaint` - Dari NOT NULL menjadi nullable

## 🔧 Migration Steps

### Step 1: Backup Database
```bash
pg_dump -U postgres service_db > backup_before_migration.sql
```

### Step 2: Run Auto Migration
```bash
go run cmd/migrate/main.go
```

AutoMigrate akan:
- Menambahkan kolom baru (nullable)
- Mengubah constraint existing kolom menjadi nullable
- Membuat foreign key untuk `service_catalog_id` dan `service_provider_id`
- Membuat indexes baru

### Step 3: Data Migration
Migration script akan otomatis:
1. Map iPhone fields ke generic item fields untuk existing orders
2. Set default `service_name` untuk existing orders
3. Set default `service_location` dari branch address
4. Set default `service_type` jika kosong

### Step 4: Verify Migration
```sql
-- Check existing orders still have data
SELECT COUNT(*) FROM service_orders WHERE item_model IS NOT NULL;

-- Check new fields are populated
SELECT COUNT(*) FROM service_orders WHERE service_name IS NOT NULL;

-- Verify backward compatibility
SELECT id, iphone_model, item_model, service_name 
FROM service_orders 
LIMIT 10;
```

## ✅ Backward Compatibility

### 1. API Compatibility
- **Legacy API tetap berfungsi**: Request dengan `branch_id` dan field iPhone masih didukung
- **New API**: Request dengan `service_catalog_id` untuk multi-service

### 2. Data Compatibility
- **Existing orders**: Tetap bisa di-read dan di-update
- **Field mapping**: iPhone fields otomatis di-map ke item fields
- **Default values**: Service name dan location di-set otomatis

### 3. Code Compatibility
- **SetAliasFields()**: Method untuk map field legacy ke field baru
- **Conditional logic**: Service layer handle kedua flow (legacy dan new)

## 🚨 Rollback Strategy

Jika migration gagal atau perlu rollback:

### Step 1: Restore Database
```bash
psql -U postgres service_db < backup_before_migration.sql
```

### Step 2: Revert Code
```bash
git checkout <previous-commit>
```

### Step 3: Rebuild
```bash
go build ./cmd/app
```

## 📊 Migration Checklist

- [x] Backup database
- [x] Run AutoMigrate untuk tabel baru
- [x] Run AutoMigrate untuk update ServiceOrder
- [x] Run data migration script
- [x] Verify indexes created
- [x] Test API dengan data existing
- [x] Test API dengan data baru
- [x] Verify backward compatibility
- [x] Update documentation

## 🔍 Post-Migration Verification

### 1. Check Table Structure
```sql
\d service_orders
```

### 2. Check Data Integrity
```sql
-- All orders should have either branch_id or service_provider_id
SELECT COUNT(*) FROM service_orders 
WHERE branch_id IS NULL AND service_provider_id IS NULL;

-- All orders should have service_name
SELECT COUNT(*) FROM service_orders WHERE service_name IS NULL OR service_name = '';
```

### 3. Check Indexes
```sql
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'service_orders';
```

### 4. Test API Endpoints
- [ ] GET /api/v1/orders (existing orders)
- [ ] POST /api/v1/orders (legacy flow)
- [ ] POST /api/v1/orders (new flow with service_catalog_id)
- [ ] GET /api/v1/services/catalogs
- [ ] GET /api/v1/services/providers

## 📝 Notes

1. **No Data Loss**: Semua data existing tetap aman karena field lama tidak dihapus
2. **Gradual Migration**: Bisa migrate ke multi-service secara bertahap
3. **Dual Support**: Sistem support kedua flow (legacy dan new) secara bersamaan
4. **Zero Downtime**: Migration bisa dilakukan tanpa downtime jika dilakukan dengan benar

## 🎯 Next Steps

Setelah migration berhasil:
1. Update client applications untuk menggunakan API baru
2. Migrate existing orders ke ServiceCatalog (optional)
3. Add ServiceCatalog entries untuk layanan baru
4. Onboard ServiceProviders untuk layanan baru
5. Deprecate legacy API endpoints (setelah semua client migrated)

