# Analisis On-Demand Service & Real-Time Tracking

## 📊 Status Saat Ini

### ✅ Yang Sudah Ada

#### 1. **On-Demand Service Support (Sebagian)**
- ✅ Field `PickupAddress`, `PickupLatitude`, `PickupLongitude` di ServiceOrder
- ✅ Field `ServiceLocation` untuk lokasi layanan
- ✅ Field `RequiresLocation` di ServiceCatalog (default: true)
  - Jika `false`, berarti service datang ke customer
- ✅ Field `RequiresPickup` dan `RequiresDelivery` di ServiceCatalog
- ✅ Support ServiceProvider dengan location (Latitude, Longitude)

#### 2. **Real-Time Communication (Sebagian)**
- ✅ WebSocket handler untuk real-time communication
- ✅ Support chat messages real-time
- ✅ Support typing indicators
- ✅ Room-based messaging (per order)
- ✅ Broadcast messages ke semua client di room

### ❌ Yang Belum Ada

#### 1. **On-Demand Service Logic**
- ❌ Belum ada logic khusus untuk handle on-demand service
- ❌ Belum ada field untuk current location provider/courier
- ❌ Belum ada logic untuk menentukan apakah service datang ke customer atau customer datang ke provider
- ❌ Belum ada field untuk ETA (Estimated Time of Arrival)

#### 2. **Real-Time Location Tracking**
- ❌ Belum ada message type untuk location tracking di WebSocket
- ❌ Belum ada model untuk store location history
- ❌ Belum ada endpoint untuk update location real-time
- ❌ Belum ada logic untuk broadcast location updates
- ❌ Belum ada field untuk current location courier/provider
- ❌ Belum ada real-time ETA calculation

## 🎯 Analisis Kebutuhan

### 1. On-Demand Service (Service Datang ke Customer)

**Karakteristik:**
- Service provider/courier datang ke lokasi customer
- Customer perlu tahu kapan provider akan datang (ETA)
- Customer perlu track lokasi provider secara real-time
- Provider perlu update lokasi secara berkala

**Contoh Layanan:**
- Massage panggilan
- Home cleaning service
- On-demand beauty service (nail art, eyelash extension)
- On-demand repair service (handphone repair)

### 2. Real-Time Tracking

**Karakteristik:**
- Provider/Courier update location secara berkala (setiap 5-10 detik)
- Customer bisa melihat lokasi provider di map real-time
- System calculate ETA berdasarkan current location
- Broadcast location updates ke customer via WebSocket

## 🔧 Rekomendasi Implementasi

### 1. **Location Tracking Model** (PRIORITAS TINGGI)

```go
// LocationTracking represents real-time location tracking for orders
type LocationTracking struct {
    ID              uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    OrderID         uuid.UUID   `gorm:"type:uuid;not null;index"`
    UserID          uuid.UUID   `gorm:"type:uuid;not null"` // Courier/Provider ID
    Latitude        float64     `gorm:"type:decimal(10,6);not null"`
    Longitude       float64     `gorm:"type:decimal(10,6);not null"`
    Accuracy        float64     `gorm:"type:decimal(10,2)"` // GPS accuracy in meters
    Speed           float64     `gorm:"type:decimal(10,2)"` // Speed in km/h
    Heading         float64     `gorm:"type:decimal(5,2)"`  // Direction in degrees
    Timestamp       time.Time   `gorm:"not null;index"`
    CreatedAt       time.Time
}

// CurrentLocation represents current location of courier/provider
type CurrentLocation struct {
    UserID          uuid.UUID   `gorm:"type:uuid;primary_key"`
    OrderID         uuid.UUID   `gorm:"type:uuid;not null;index"`
    Latitude        float64     `gorm:"type:decimal(10,6);not null"`
    Longitude       float64     `gorm:"type:decimal(10,6);not null"`
    Accuracy        float64     `gorm:"type:decimal(10,2)"`
    Speed           float64     `gorm:"type:decimal(10,2)"`
    Heading         float64     `gorm:"type:decimal(5,2)"`
    UpdatedAt       time.Time   `gorm:"not null"`
}
```

### 2. **Update ServiceOrder Model**

```go
// Tambahkan field untuk on-demand service
type ServiceOrder struct {
    // ... existing fields ...
    
    // On-demand service fields
    IsOnDemand       bool        `gorm:"default:false"` // Service datang ke customer
    CurrentLatitude  *float64    `gorm:"type:decimal(10,6)"` // Current location courier/provider
    CurrentLongitude *float64    `gorm:"type:decimal(10,6)"` // Current location courier/provider
    ETA              *int        `gorm:"default:0"` // Estimated time of arrival in minutes
    LastLocationUpdate *time.Time `gorm:"type:timestamp"` // Last location update time
}
```

### 3. **WebSocket Message Types untuk Tracking**

```go
// Tambahkan message types di WebSocket handler
type Message struct {
    Type      string      `json:"type"` // "chat", "typing", "location", "eta", "status"
    OrderID   string      `json:"order_id"`
    UserID    string      `json:"user_id"`
    Content   string      `json:"content"`
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data,omitempty"`
}

// LocationUpdateData untuk location tracking
type LocationUpdateData struct {
    Latitude    float64   `json:"latitude"`
    Longitude   float64   `json:"longitude"`
    Accuracy    float64   `json:"accuracy,omitempty"`
    Speed       float64   `json:"speed,omitempty"`
    Heading     float64   `json:"heading,omitempty"`
    ETA         int       `json:"eta,omitempty"` // in minutes
}
```

### 4. **API Endpoints Baru**

```go
// Update location real-time (untuk courier/provider)
POST /api/v1/orders/:id/location
{
    "latitude": -6.2088,
    "longitude": 106.8456,
    "accuracy": 10.5,
    "speed": 45.2,
    "heading": 90.0
}

// Get current location (untuk customer)
GET /api/v1/orders/:id/location

// Get location history
GET /api/v1/orders/:id/location/history

// Get ETA
GET /api/v1/orders/:id/eta
```

### 5. **WebSocket Handler Update**

```go
// Tambahkan handler untuk location tracking
func (h *WebSocketHandler) handleLocationMessage(client *Client, message *Message) {
    // Parse location data
    locationData, ok := message.Data.(map[string]interface{})
    if !ok {
        return
    }
    
    // Update location di database
    // Broadcast ke customer di room yang sama
    h.broadcastToRoomExcluding(message.OrderID, client.Conn, Message{
        Type:      "location",
        OrderID:   message.OrderID,
        UserID:    client.UserID.String(),
        Timestamp: time.Now(),
        Data:      locationData,
    })
}
```

## 📋 Rencana Implementasi

### Phase 1: Location Tracking Model (1-2 hari)
1. ✅ Buat model LocationTracking
2. ✅ Buat model CurrentLocation
3. ✅ Buat migration
4. ✅ Buat repository untuk location tracking

### Phase 2: API Endpoints (1-2 hari)
1. ✅ Buat endpoint untuk update location
2. ✅ Buat endpoint untuk get current location
3. ✅ Buat endpoint untuk get location history
4. ✅ Buat endpoint untuk calculate ETA

### Phase 3: WebSocket Integration (1-2 hari)
1. ✅ Update WebSocket handler untuk support location messages
2. ✅ Implementasi broadcast location updates
3. ✅ Implementasi ETA calculation real-time

### Phase 4: On-Demand Service Logic (1-2 hari)
1. ✅ Update ServiceCatalog untuk support on-demand
2. ✅ Update OrderService untuk handle on-demand
3. ✅ Update logic untuk determine service location

### Phase 5: Testing & Optimization (1-2 hari)
1. ✅ Unit testing
2. ✅ Integration testing
3. ✅ Performance optimization
4. ✅ Documentation

## 🎯 Kesimpulan

### ✅ **On-Demand Service: SEBAGIAN SUPPORT**

**Yang sudah ada:**
- Field location di ServiceOrder
- Field RequiresLocation di ServiceCatalog
- Support ServiceProvider dengan location

**Yang perlu ditambahkan:**
- Logic untuk handle on-demand service
- Field untuk current location provider
- Field untuk ETA
- Logic untuk determine service location

### ❌ **Real-Time Tracking: BELUM SUPPORT**

**Yang sudah ada:**
- WebSocket infrastructure
- Real-time communication (chat)

**Yang perlu ditambahkan:**
- Location tracking model
- API endpoints untuk location updates
- WebSocket message type untuk location
- Broadcast location updates
- ETA calculation real-time

## 📝 Rekomendasi Prioritas

1. **PRIORITAS TINGGI**: Location Tracking Model & API
2. **PRIORITAS TINGGI**: WebSocket Location Updates
3. **PRIORITAS SEDANG**: On-Demand Service Logic
4. **PRIORITAS SEDANG**: ETA Calculation
5. **PRIORITAS RENDAH**: Location History

## 🚀 Estimasi Waktu

**Total: 5-10 hari** untuk implementasi lengkap on-demand service dan real-time tracking.

