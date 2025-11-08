# Simulasi Test Sistem Realtime Tracking

## Deskripsi

Script ini mensimulasikan alur lengkap customer order dengan realtime tracking:
1. Customer membuat pesanan
2. Courier update lokasi secara real-time
3. Customer track lokasi via REST API dan WebSocket
4. Simulasi pergerakan kurir dari titik A ke titik B

## Requirements

```bash
pip install -r scripts/requirements_simulasi.txt
```

Atau install manual:
```bash
pip install requests websocket-client
```

## Cara Menjalankan

### 1. Pastikan Server Berjalan

```bash
# Pastikan server Go berjalan di localhost:8080
go run cmd/app/main.go
```

### 2. Jalankan Simulasi

```bash
# Windows
python scripts/simulasi_realtime_tracking.py

# Linux/Mac
python3 scripts/simulasi_realtime_tracking.py
```

Atau buat executable:
```bash
chmod +x scripts/simulasi_realtime_tracking.py
./scripts/simulasi_realtime_tracking.py
```

## Konfigurasi

Edit variabel di awal script untuk mengubah konfigurasi:

```python
BASE_URL = "http://localhost:8080"
WS_URL = "ws://localhost:8080/ws/chat"

CUSTOMER_EMAIL = "customer@test.com"
CUSTOMER_PASSWORD = "password123"
COURIER_EMAIL = "courier@test.com"
COURIER_PASSWORD = "password123"
```

## Output

Script akan menampilkan:
- ✅ Status registrasi dan login
- ✅ Status pembuatan pesanan
- ✅ Update lokasi courier secara real-time
- ✅ Tracking via REST API (polling)
- ✅ Tracking via WebSocket (real-time)
- ✅ Ringkasan simulasi

## Fitur Simulasi

### 1. REST API Tracking
- Polling lokasi setiap 3 detik
- Menampilkan koordinat, jarak, dan ETA

### 2. WebSocket Tracking
- Koneksi real-time ke server
- Menerima update lokasi secara real-time
- Menerima update status pesanan

### 3. Simulasi Pergerakan Courier
- Simulasi pergerakan dari pickup location ke destination
- Update lokasi setiap 2 detik
- Menghitung jarak dan ETA otomatis

## Troubleshooting

### Error: Connection refused
- Pastikan server Go berjalan di port 8080
- Cek konfigurasi BASE_URL

### Error: WebSocket connection failed
- Pastikan WebSocket handler sudah terdaftar di router
- Cek endpoint `/ws/chat` tersedia

### Error: Authentication failed
- Pastikan user customer dan courier sudah terdaftar
- Cek email dan password di konfigurasi

## Catatan

- Script ini menggunakan data dummy untuk simulasi
- Pastikan database sudah ter-setup dengan benar
- Untuk production, gunakan data real dan handle error dengan lebih baik

