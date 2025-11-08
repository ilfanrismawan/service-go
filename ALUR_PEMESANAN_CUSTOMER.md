# Alur Proses Customer dalam Mengorder Jasa

## 📋 Daftar Isi
1. [Overview](#overview)
2. [Tahap 1: Registrasi & Login](#tahap-1-registrasi--login)
3. [Tahap 2: Pilih Layanan](#tahap-2-pilih-layanan)
4. [Tahap 3: Buat Pesanan](#tahap-3-buat-pesanan)
5. [Tahap 4: Tracking Pesanan](#tahap-4-tracking-pesanan)
6. [Tahap 5: Pembayaran](#tahap-5-pembayaran)
7. [Tahap 6: Penyelesaian](#tahap-6-penyelesaian)
8. [Status Pesanan](#status-pesanan)
9. [API Endpoints](#api-endpoints)

---

## Overview

Sistem ini memungkinkan customer untuk memesan jasa service iPhone dengan berbagai fitur:
- ✅ Multi-service support (berbagai jenis layanan)
- ✅ On-demand service (layanan datang ke customer)
- ✅ Real-time tracking (tracking lokasi kurir/provider secara real-time)
- ✅ WebSocket untuk update real-time
- ✅ Notifikasi status pesanan
- ✅ Chat dengan kurir/teknisi

---

## Tahap 1: Registrasi & Login

### 1.1 Registrasi Customer

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "phone": "+6281234567890",
  "role": "pelanggan"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "pelanggan"
  },
  "message": "User registered successfully"
}
```

### 1.2 Login

**Endpoint:** `POST /api/v1/auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "access_token": "jwt_token",
    "refresh_token": "refresh_token",
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "pelanggan"
    },
    "expires_in": 3600
  }
}
```

**Catatan:** Simpan `access_token` untuk digunakan di header Authorization pada request berikutnya:
```
Authorization: Bearer <access_token>
```

---

## Tahap 2: Pilih Layanan

### 2.1 Lihat Katalog Layanan

**Endpoint:** `GET /api/v1/services/catalog`

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "name": "Service iPhone Screen Repair",
      "description": "Perbaikan layar iPhone",
      "base_price": 500000,
      "estimated_duration": 120,
      "requires_location": false,
      "requires_appointment": false,
      "is_active": true
    }
  ]
}
```

### 2.2 Lihat Provider Layanan (Opsional)

**Endpoint:** `GET /api/v1/services/providers?catalog_id=<catalog_id>`

### 2.3 Lihat Cabang (Jika Layanan di Lokasi)

**Endpoint:** `GET /api/v1/branches`

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "name": "Cabang Jakarta Pusat",
      "address": "Jl. Sudirman No. 123",
      "latitude": -6.2088,
      "longitude": 106.8456,
      "phone": "+6281234567890"
    }
  ]
}
```

---

## Tahap 3: Buat Pesanan

### 3.1 Buat Pesanan Baru

**Endpoint:** `POST /api/v1/orders`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

**Untuk Layanan On-Demand (datang ke customer):**
```json
{
  "service_catalog_id": "uuid",
  "service_provider_id": "uuid",
  "description": "iPhone saya layarnya retak",
  "complaint": "Layar iPhone retak setelah jatuh",
  "item_model": "iPhone 13 Pro",
  "item_color": "Blue",
  "item_serial": "IMEI123456789",
  "item_type": "iPhone",
  "pickup_address": "Jl. Kebon Jeruk No. 45, Jakarta Barat",
  "pickup_latitude": -6.1944,
  "pickup_longitude": 106.8229,
  "estimated_cost": 500000,
  "estimated_duration": 120
}
```

**Untuk Layanan di Lokasi Provider/Branch:**
```json
{
  "service_catalog_id": "uuid",
  "branch_id": "uuid",
  "description": "iPhone saya layarnya retak",
  "complaint": "Layar iPhone retak setelah jatuh",
  "item_model": "iPhone 13 Pro",
  "item_color": "Blue",
  "item_serial": "IMEI123456789",
  "item_type": "iPhone",
  "appointment_date": "2025-01-15",
  "appointment_time": "14:00",
  "estimated_cost": 500000,
  "estimated_duration": 120
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order_uuid",
    "order_number": "ORD-20250115-001234",
    "customer_id": "customer_uuid",
    "service_catalog_id": "catalog_uuid",
    "service_name": "Service iPhone Screen Repair",
    "status": "pending_pickup",
    "is_on_demand": true,
    "pickup_address": "Jl. Kebon Jeruk No. 45, Jakarta Barat",
    "service_location": "Jl. Kebon Jeruk No. 45, Jakarta Barat",
    "estimated_cost": 500000,
    "estimated_duration": 120,
    "created_at": "2025-01-15T10:00:00Z"
  },
  "message": "Order created successfully"
}
```

**Catatan Penting:**
- Simpan `order_number` dan `id` untuk tracking
- Status awal: `pending_pickup`
- Untuk on-demand service, kurir akan datang ke alamat pickup

---

## Tahap 4: Tracking Pesanan

### 4.1 Lihat Detail Pesanan

**Endpoint:** `GET /api/v1/orders/{order_id}`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order_uuid",
    "order_number": "ORD-20250115-001234",
    "status": "on_pickup",
    "courier": {
      "id": "courier_uuid",
      "name": "Budi Kurir",
      "phone": "+6281234567890"
    },
    "current_latitude": -6.1944,
    "current_longitude": 106.8229,
    "eta": 15,
    "last_location_update": "2025-01-15T10:30:00Z"
  }
}
```

### 4.2 Real-time Tracking via REST API

**Endpoint:** `GET /api/v1/orders/{order_id}/location`

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "location_uuid",
    "order_id": "order_uuid",
    "user_id": "courier_uuid",
    "latitude": -6.1944,
    "longitude": 106.8229,
    "accuracy": 10.5,
    "speed": 30.0,
    "heading": 45.0,
    "eta": 15,
    "distance": 5.2,
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

### 4.3 Real-time Tracking via WebSocket

**Endpoint:** `WS /ws/chat?order_id={order_id}`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/chat?order_id=order_uuid', {
  headers: {
    'Authorization': 'Bearer <access_token>'
  }
});
```

**Menerima Update Lokasi:**
```json
{
  "type": "location",
  "order_id": "order_uuid",
  "user_id": "courier_uuid",
  "timestamp": "2025-01-15T10:30:00Z",
  "data": {
    "latitude": -6.1944,
    "longitude": 106.8229,
    "accuracy": 10.5,
    "speed": 30.0,
    "heading": 45.0,
    "eta": 15,
    "distance": 5.2
  }
}
```

**Menerima Update Status:**
```json
{
  "type": "status",
  "order_id": "order_uuid",
  "timestamp": "2025-01-15T10:35:00Z",
  "data": {
    "status": "in_service",
    "message": "Pesanan sedang dalam proses perbaikan"
  }
}
```

### 4.4 Lihat History Lokasi

**Endpoint:** `GET /api/v1/orders/{order_id}/location/history?limit=50`

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "tracking_uuid",
      "order_id": "order_uuid",
      "user_id": "courier_uuid",
      "latitude": -6.1944,
      "longitude": 106.8229,
      "accuracy": 10.5,
      "speed": 30.0,
      "heading": 45.0,
      "timestamp": "2025-01-15T10:30:00Z"
    }
  ]
}
```

### 4.5 Lihat ETA (Estimated Time of Arrival)

**Endpoint:** `GET /api/v1/orders/{order_id}/eta`

**Response:**
```json
{
  "status": "success",
  "data": {
    "eta": 15,
    "distance": 5.2
  }
}
```

---

## Tahap 5: Pembayaran

### 5.1 Buat Invoice

**Endpoint:** `POST /api/v1/payments/create-invoice`

**Request Body:**
```json
{
  "order_id": "order_uuid",
  "payment_method": "midtrans"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "invoice_id": "invoice_uuid",
    "order_id": "order_uuid",
    "amount": 500000,
    "payment_url": "https://app.midtrans.com/snap/v2/...",
    "status": "pending"
  }
}
```

### 5.2 Proses Pembayaran

**Endpoint:** `POST /api/v1/payments/process`

**Request Body:**
```json
{
  "invoice_id": "invoice_uuid",
  "payment_method": "midtrans",
  "payment_data": {
    "token": "midtrans_token"
  }
}
```

### 5.3 Cek Status Pembayaran

**Endpoint:** `GET /api/v1/payments/order/{order_id}`

---

## Tahap 6: Penyelesaian

### 6.1 Lihat Status Pesanan

**Endpoint:** `GET /api/v1/orders/{order_id}`

**Status yang mungkin:**
- `pending_pickup` - Menunggu pickup
- `on_pickup` - Kurir sedang dalam perjalanan untuk pickup
- `in_service` - Pesanan sedang dalam proses perbaikan
- `ready` - Pesanan selesai dan siap diambil/dikirim
- `delivered` - Pesanan sudah dikirim/diterima
- `completed` - Pesanan selesai
- `cancelled` - Pesanan dibatalkan

### 6.2 Berikan Rating

**Endpoint:** `POST /api/v1/ratings`

**Request Body:**
```json
{
  "order_id": "order_uuid",
  "rating": 5,
  "comment": "Pelayanan sangat baik, teknisi profesional"
}
```

### 6.3 Lihat Pesanan Saya

**Endpoint:** `GET /api/v1/orders/my`

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "order_uuid",
      "order_number": "ORD-20250115-001234",
      "status": "completed",
      "service_name": "Service iPhone Screen Repair",
      "estimated_cost": 500000,
      "actual_cost": 500000,
      "created_at": "2025-01-15T10:00:00Z"
    }
  ]
}
```

---

## Status Pesanan

| Status | Deskripsi | Aksi Customer |
|--------|-----------|---------------|
| `pending_pickup` | Menunggu kurir/provider untuk pickup | Tunggu konfirmasi |
| `on_pickup` | Kurir/provider sedang dalam perjalanan | Track lokasi real-time |
| `in_service` | Pesanan sedang diperbaiki | Tunggu notifikasi selesai |
| `ready` | Pesanan selesai, siap diambil/dikirim | Tunggu pengiriman atau ambil di lokasi |
| `delivered` | Pesanan sudah dikirim/diterima | Konfirmasi penerimaan |
| `completed` | Pesanan selesai | Berikan rating |
| `cancelled` | Pesanan dibatalkan | - |

---

## API Endpoints Summary

### Authentication
- `POST /api/v1/auth/register` - Registrasi
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/profile` - Lihat profil
- `PUT /api/v1/auth/profile` - Update profil

### Orders
- `POST /api/v1/orders` - Buat pesanan
- `GET /api/v1/orders` - List pesanan (dengan filter)
- `GET /api/v1/orders/my` - Pesanan saya
- `GET /api/v1/orders/{id}` - Detail pesanan
- `GET /api/v1/orders/number/{order_number}` - Detail pesanan by order number

### Tracking
- `GET /api/v1/orders/{id}/location` - Lokasi saat ini
- `GET /api/v1/orders/{id}/location/history` - History lokasi
- `GET /api/v1/orders/{id}/eta` - ETA
- `WS /ws/chat?order_id={id}` - WebSocket real-time tracking

### Payments
- `POST /api/v1/payments/create-invoice` - Buat invoice
- `POST /api/v1/payments/process` - Proses pembayaran
- `GET /api/v1/payments/order/{order_id}` - Status pembayaran

### Services
- `GET /api/v1/services/catalog` - Katalog layanan
- `GET /api/v1/services/providers` - List provider
- `GET /api/v1/branches` - List cabang

### Chat
- `GET /api/v1/chat/orders/{order_id}` - Chat messages
- `POST /api/v1/chat/orders/{order_id}` - Kirim pesan

### Ratings
- `POST /api/v1/ratings` - Berikan rating
- `GET /api/v1/ratings` - List rating

---

## Tips untuk Customer

1. **Simpan Order Number**: Simpan order number untuk tracking mudah
2. **Aktifkan Notifikasi**: Pastikan notifikasi aktif untuk update real-time
3. **Gunakan WebSocket**: Untuk tracking real-time yang lebih akurat
4. **Siapkan Lokasi**: Pastikan alamat pickup akurat dengan koordinat GPS
5. **Monitor Status**: Cek status pesanan secara berkala
6. **Chat Support**: Gunakan fitur chat untuk komunikasi dengan kurir/teknisi

---

## Contoh Flow Lengkap

```
1. Customer registrasi/login
   ↓
2. Customer lihat katalog layanan
   ↓
3. Customer pilih layanan dan buat pesanan
   ↓
4. Sistem assign kurir (jika on-demand)
   ↓
5. Customer track lokasi kurir via WebSocket/REST API
   ↓
6. Kurir pickup item
   ↓
7. Status berubah ke "in_service"
   ↓
8. Teknisi perbaiki item
   ↓
9. Status berubah ke "ready"
   ↓
10. Customer bayar invoice
    ↓
11. Kurir kirim kembali (jika on-demand)
    ↓
12. Status berubah ke "delivered"
    ↓
13. Customer konfirmasi penerimaan
    ↓
14. Status berubah ke "completed"
    ↓
15. Customer berikan rating
```

---

**Dokumen ini dibuat untuk membantu customer memahami alur pemesanan jasa secara lengkap.**

