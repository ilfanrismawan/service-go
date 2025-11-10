# Google Maps API Integration

Dokumentasi lengkap tentang integrasi Google Maps API dengan sistem real-time tracking.

## 📋 Overview

Sistem real-time tracking telah terintegrasi dengan Google Maps Distance Matrix API untuk memberikan:
- Jarak aktual berdasarkan rute jalan (bukan garis lurus)
- Waktu tempuh yang akurat dengan mempertimbangkan kondisi lalu lintas
- ETA (Estimated Time of Arrival) yang lebih presisi

## 🔧 Setup

### 1. Mendapatkan Google Maps API Key

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Buat project baru atau pilih project yang ada
3. Aktifkan **Distance Matrix API**:
   - Pergi ke "APIs & Services" > "Library"
   - Cari "Distance Matrix API"
   - Klik "Enable"
4. Buat API Key:
   - Pergi ke "APIs & Services" > "Credentials"
   - Klik "Create Credentials" > "API Key"
   - Copy API key yang dihasilkan
5. (Opsional) Batasi API Key:
   - Klik pada API key yang baru dibuat
   - Di "API restrictions", pilih "Restrict key"
   - Pilih "Distance Matrix API"
   - Di "Application restrictions", batasi berdasarkan IP atau referrer sesuai kebutuhan

### 2. Konfigurasi Environment Variable

Tambahkan API key ke file `.env`:

```bash
GOOGLE_MAPS_API_KEY=your-actual-api-key-here
```

Atau set sebagai environment variable:

```bash
export GOOGLE_MAPS_API_KEY=your-actual-api-key-here
```

### 3. Verifikasi Konfigurasi

Setelah API key di-set, sistem akan otomatis menggunakan Google Maps API untuk perhitungan jarak dan ETA. Jika API key tidak tersedia atau terjadi error, sistem akan otomatis fallback ke perhitungan Haversine.

## 🏗️ Arsitektur

### Komponen

1. **GoogleMapsService** (`internal/modules/tracking/service/google_maps_service.go`)
   - Service untuk berinteraksi dengan Google Maps Distance Matrix API
   - Method utama: `GetDistanceAndDuration()`
   - Mengembalikan `RouteInfo` dengan distance, duration, dan duration_in_traffic

2. **LocationTrackingService** (`internal/modules/tracking/service/location_tracking_service.go`)
   - Service utama untuk location tracking
   - Menggunakan `GoogleMapsService` untuk perhitungan jarak dan ETA
   - Fallback ke Haversine jika Google Maps API tidak tersedia

### Flow

```
Location Update Request
    ↓
LocationTrackingService.UpdateLocation()
    ↓
GoogleMapsService.GetDistanceAndDuration()
    ↓
[Success] → Use Google Maps API result (distance + duration_in_traffic)
[Error]   → Fallback to Haversine calculation
    ↓
Save to database & broadcast via WebSocket
```

## 📊 API Response Structure

### Google Maps Distance Matrix API Response

```json
{
  "status": "OK",
  "origin_addresses": ["Jakarta, Indonesia"],
  "destination_addresses": ["Bandung, Indonesia"],
  "rows": [
    {
      "elements": [
        {
          "status": "OK",
          "distance": {
            "value": 150000,  // in meters
            "text": "150 km"
          },
          "duration": {
            "value": 7200,    // in seconds
            "text": "2 hours"
          },
          "duration_in_traffic": {
            "value": 9000,    // in seconds (with traffic)
            "text": "2 hours 30 mins"
          }
        }
      ]
    }
  ]
}
```

### RouteInfo (Internal)

```go
type RouteInfo struct {
    Distance          float64 // Distance in kilometers
    Duration          int     // Duration in minutes
    DurationInTraffic int     // Duration in traffic in minutes (if available)
    Status            string  // Status of the API call
}
```

## 🔄 Fallback Mechanism

Sistem memiliki mekanisme fallback yang robust:

1. **Jika API key tidak di-set:**
   - Sistem langsung menggunakan perhitungan Haversine
   - Tidak ada error, hanya log warning

2. **Jika API call gagal:**
   - Error ditangkap dan di-log
   - Sistem fallback ke perhitungan Haversine
   - Response tetap dikembalikan dengan data dari Haversine

3. **Jika API rate limit tercapai:**
   - Error ditangkap
   - Fallback ke Haversine
   - Disarankan untuk upgrade quota atau implement caching

## 💰 Pricing & Quotas

### Google Maps Distance Matrix API Pricing

- **Free tier:** $200 credit per bulan (setara dengan ~40,000 requests)
- **Per request:** $0.005 (setelah free tier)
- **Traffic data:** Tidak ada biaya tambahan

### Best Practices untuk Menghemat Cost

1. **Caching:**
   - Cache hasil perhitungan untuk rute yang sama dalam waktu singkat
   - Update cache setiap 5-10 menit untuk rute yang sama

2. **Rate Limiting:**
   - Implement rate limiting untuk mencegah abuse
   - Monitor usage melalui Google Cloud Console

3. **Optimasi Request:**
   - Hanya panggil API saat lokasi berubah signifikan
   - Gunakan batch request jika memungkinkan

## 🧪 Testing

### Test dengan API Key

```bash
# Set API key
export GOOGLE_MAPS_API_KEY=your-api-key

# Run application
go run cmd/app/main.go

# Test location update
curl -X POST http://localhost:8080/api/v1/orders/{order_id}/location \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": -6.2088,
    "longitude": 106.8456,
    "speed": 50,
    "heading": 90
  }'
```

### Test tanpa API Key (Fallback)

```bash
# Unset API key
unset GOOGLE_MAPS_API_KEY

# Run application
go run cmd/app/main.go

# System akan otomatis menggunakan Haversine
```

## 📝 Logging

Sistem akan mencatat log untuk:
- Success: Google Maps API call berhasil
- Warning: API key tidak di-set, menggunakan fallback
- Error: API call gagal, menggunakan fallback

Contoh log:
```
[INFO] Using Google Maps API for distance calculation
[WARN] Google Maps API key not configured, using Haversine fallback
[ERROR] Google Maps API error: REQUEST_DENIED - Invalid API key, using Haversine fallback
```

## 🔒 Security

1. **Jangan commit API key ke repository:**
   - Gunakan environment variables
   - Tambahkan `.env` ke `.gitignore`

2. **Batasi API key:**
   - Restrict by IP (untuk production)
   - Restrict by referrer (untuk web apps)
   - Enable hanya API yang diperlukan

3. **Monitor usage:**
   - Set up billing alerts di Google Cloud Console
   - Monitor API usage secara berkala

## 🐛 Troubleshooting

### API Key Invalid

**Error:** `REQUEST_DENIED - Invalid API key`

**Solusi:**
- Pastikan API key sudah di-copy dengan benar
- Pastikan Distance Matrix API sudah di-enable
- Cek apakah API key sudah di-restrict (mungkin perlu di-allow untuk IP server)

### Quota Exceeded

**Error:** `OVER_QUERY_LIMIT`

**Solusi:**
- Cek usage di Google Cloud Console
- Upgrade quota jika diperlukan
- Implement caching untuk mengurangi request

### No Route Found

**Error:** `NOT_FOUND` atau `ZERO_RESULTS`

**Solusi:**
- Pastikan koordinat valid (latitude: -90 to 90, longitude: -180 to 180)
- Cek apakah lokasi dapat diakses via jalan (bukan di tengah laut/gunung)

## 📚 Referensi

- [Google Maps Distance Matrix API Documentation](https://developers.google.com/maps/documentation/distance-matrix)
- [Google Cloud Console](https://console.cloud.google.com/)
- [API Pricing](https://developers.google.com/maps/billing-and-pricing/pricing)

## 🔄 Update History

- **2024-01-XX:** Initial integration dengan Google Maps Distance Matrix API
- **2024-01-XX:** Added fallback mechanism ke Haversine
- **2024-01-XX:** Added traffic data support (duration_in_traffic)


