# Panduan Pengujian Simulasi Real-Time Tracking

## 📋 Overview

Sistem real-time tracking memungkinkan customer untuk melacak lokasi kurir secara real-time melalui REST API dan WebSocket. Script simulasi ini menguji seluruh alur dari pembuatan order hingga tracking lokasi.

## 🚀 Quick Start

### 1. Pastikan Server Berjalan

**Opsi A: Menggunakan Docker (Recommended)**
```bash
# Build dan start semua services
DOCKER_SUDO=sudo make docker-up-build

# Atau jika user sudah di docker group
make docker-up-build

# Cek status
make docker-ps

# Lihat logs
make docker-logs-app
```

**Opsi B: Menjalankan Lokal**
```bash
# Pastikan database dan dependencies sudah setup
make run
```

### 2. Setup Dependencies Python

**Opsi A: Otomatis (Recommended)**
```bash
make test-simulasi-setup
```

**Opsi B: Manual - Virtual Environment (RECOMMENDED)**
```bash
# Install python3-venv (jika belum ada)
sudo apt install python3-venv

# Buat virtual environment
python3 -m venv venv

# Install dependencies
venv/bin/pip install requests websocket-client
```

**Opsi C: Manual - User Install**
```bash
python3 -m pip install --user requests websocket-client
```

**Opsi D: Manual - System Install (Tidak Disarankan)**
```bash
python3 -m pip install --break-system-packages requests websocket-client
```

**Catatan:** Python 3.12+ menggunakan PEP 668 yang mencegah install package secara global. Gunakan virtual environment atau `--user` flag.

### 3. Jalankan Simulasi

**Simple Test (REST API only):**
```bash
make test-simulasi-simple
# atau
make test-simulasi
```

**Full Test (REST API + WebSocket):**
```bash
make test-simulasi-full
```

## 📝 Detail Simulasi

### Simple Test (`test-simulasi-simple`)

Menguji tracking menggunakan REST API polling:
- ✅ Login customer dan courier
- ✅ Membuat pesanan baru
- ✅ Update lokasi courier (simulasi pergerakan)
- ✅ Customer track lokasi via REST API
- ✅ Lihat history lokasi

**Durasi:** ~15-20 detik

### Full Test (`test-simulasi-full`)

Menguji tracking menggunakan REST API + WebSocket:
- ✅ Semua fitur dari simple test
- ✅ Koneksi WebSocket real-time
- ✅ Menerima update lokasi via WebSocket
- ✅ Menerima update status pesanan
- ✅ Simulasi pergerakan kurir yang lebih detail

**Durasi:** ~30-60 detik

## 🔧 Konfigurasi

Edit file script untuk mengubah konfigurasi:

**`scripts/test_realtime_tracking_simple.py`** atau **`scripts/simulasi_realtime_tracking.py`**:

```python
BASE_URL = "http://localhost:8080"
WS_URL = "ws://localhost:8080/ws/chat"

CUSTOMER_EMAIL = "customer@test.com"
CUSTOMER_PASSWORD = "password123"
COURIER_EMAIL = "courier@test.com"
COURIER_PASSWORD = "password123"
```

## 📊 Output yang Diharapkan

### Simple Test Output:
```
============================================================
SIMPLE REALTIME TRACKING TEST
============================================================

ℹ 1. Login Customer...
✓ Login customer berhasil
ℹ 2. Login Courier...
✓ Login courier berhasil
ℹ 3. Buat Pesanan...
✓ Pesanan berhasil dibuat: ORD-20241109-001
ℹ Order ID: 550e8400-e29b-41d4-a716-446655440000
ℹ 4. Simulasi Update Lokasi Courier...
✓ Step 1: Lokasi updated | Jarak: 2.45 km | ETA: 5 menit
✓ Step 2: Lokasi updated | Jarak: 1.23 km | ETA: 3 menit
...
ℹ 5. Customer Track Lokasi...
✓ Lokasi saat ini: (-6.194400, 106.822900)
ℹ Jarak ke tujuan: 0.00 km
ℹ ETA: 0 menit
```

### Full Test Output:
```
============================================================
SIMULASI REALTIME TRACKING
============================================================

✓ User customer@test.com registered successfully
✓ User courier@test.com registered successfully
✓ Customer logged in successfully
✓ Courier logged in successfully
✓ Order created: ORD-20241109-001
✓ WebSocket connected
✓ Location update received via WebSocket
...
```

## 🐛 Troubleshooting

### Error: Server tidak berjalan
```bash
# Cek apakah server berjalan
curl http://localhost:8080/health

# Jika tidak, start server:
make docker-up-build  # atau make run
```

### Error: Python dependencies tidak terinstall
```bash
# Install manual
python3 -m pip install --user requests websocket-client

# Atau dengan sudo (jika perlu)
sudo pip3 install requests websocket-client
```

### Error: Authentication failed
- Pastikan user customer dan courier sudah terdaftar di database
- Script akan mencoba register otomatis jika user belum ada
- Cek email dan password di konfigurasi script

### Error: WebSocket connection failed
- Pastikan WebSocket handler sudah terdaftar di router
- Cek endpoint `/ws/chat` atau `/ws/tracking` tersedia
- Pastikan server mendukung WebSocket

### Error: Order creation failed
- Pastikan database sudah ter-setup
- Cek apakah ada branch yang aktif
- Pastikan user customer memiliki role yang benar

## 📈 Monitoring

### Lihat Logs Server

**Docker:**
```bash
make docker-logs-app
```

**Local:**
```bash
# Logs akan muncul di terminal dimana server dijalankan
```

### Test Endpoint Manual

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@test.com","password":"password123"}'
```

## 🎯 Test Scenarios

### Scenario 1: Basic Tracking
1. Customer membuat order
2. Courier update lokasi sekali
3. Customer cek lokasi via REST API

### Scenario 2: Continuous Tracking
1. Customer membuat order
2. Courier update lokasi setiap 2 detik (simulasi pergerakan)
3. Customer track lokasi via REST API (polling setiap 3 detik)
4. Lihat history lokasi

### Scenario 3: Real-Time WebSocket
1. Customer membuat order
2. Customer connect ke WebSocket
3. Courier update lokasi
4. Customer menerima update via WebSocket secara real-time

## 📚 File Terkait

- `scripts/test_realtime_tracking_simple.py` - Simple test script
- `scripts/simulasi_realtime_tracking.py` - Full test dengan WebSocket
- `internal/modules/tracking/` - Implementation tracking system
- `scripts/README_SIMULASI.md` - Dokumentasi original

## 🔗 API Endpoints

- `POST /api/v1/orders/:id/location` - Update lokasi (Courier)
- `GET /api/v1/orders/:id/location` - Get current location (Customer)
- `GET /api/v1/orders/:id/location/history` - Get location history
- `POST /api/v1/orders/:id/eta` - Calculate ETA
- `GET /api/v1/orders/:id/eta` - Get ETA
- `WS /ws/chat` - WebSocket untuk real-time updates

## ✅ Checklist Sebelum Testing

- [ ] Server berjalan di http://localhost:8080
- [ ] Database sudah ter-setup dan berjalan
- [ ] Python dependencies terinstall (requests, websocket-client)
- [ ] User customer dan courier tersedia (atau akan dibuat otomatis)
- [ ] Port 8080 tidak digunakan aplikasi lain

