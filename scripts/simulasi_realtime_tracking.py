#!/usr/bin/env python3
"""
Simulasi Test Sistem Realtime Tracking
========================================

Script ini mensimulasikan alur lengkap customer order dengan realtime tracking:
1. Customer membuat pesanan
2. Courier update lokasi secara real-time
3. Customer track lokasi via REST API dan WebSocket
4. Simulasi pergerakan kurir dari titik A ke titik B

Usage:
    python scripts/simulasi_realtime_tracking.py

Requirements:
    pip install requests websocket-client
"""

import requests
import json
import time
import threading
import websocket
from datetime import datetime
from typing import Optional, Dict, Any
import math

# Konfigurasi
BASE_URL = "http://localhost:8080"
WS_URL = "ws://localhost:8080/ws/chat"

# Data simulasi
CUSTOMER_EMAIL = "customer@test.com"
CUSTOMER_PASSWORD = "password123"
COURIER_EMAIL = "courier@test.com"
COURIER_PASSWORD = "password123"

# Koordinat simulasi (Jakarta)
PICKUP_LOCATION = {
    "latitude": -6.2088,  # Jakarta Pusat
    "longitude": 106.8456
}

DESTINATION_LOCATION = {
    "latitude": -6.1944,  # Jakarta Barat
    "longitude": 106.8229
}

# Variabel global
customer_token: Optional[str] = None
courier_token: Optional[str] = None
order_id: Optional[str] = None
order_number: Optional[str] = None
ws_connected = False
location_updates = []


class Colors:
    """ANSI color codes untuk terminal output"""
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


def print_header(text: str):
    """Print header dengan format"""
    print(f"\n{Colors.HEADER}{Colors.BOLD}{'='*60}{Colors.ENDC}")
    print(f"{Colors.HEADER}{Colors.BOLD}{text:^60}{Colors.ENDC}")
    print(f"{Colors.HEADER}{Colors.BOLD}{'='*60}{Colors.ENDC}\n")


def print_success(text: str):
    """Print success message"""
    print(f"{Colors.OKGREEN}✓ {text}{Colors.ENDC}")


def print_info(text: str):
    """Print info message"""
    print(f"{Colors.OKCYAN}ℹ {text}{Colors.ENDC}")


def print_warning(text: str):
    """Print warning message"""
    print(f"{Colors.WARNING}⚠ {text}{Colors.ENDC}")


def print_error(text: str):
    """Print error message"""
    print(f"{Colors.FAIL}✗ {text}{Colors.ENDC}")


def register_user(email: str, password: str, name: str, role: str = "pelanggan") -> bool:
    """Registrasi user baru"""
    url = f"{BASE_URL}/api/v1/auth/register"
    payload = {
        "name": name,
        "email": email,
        "password": password,
        "phone": "+6281234567890",
        "role": role
    }
    
    try:
        response = requests.post(url, json=payload)
        if response.status_code in [200, 201]:
            print_success(f"User {name} berhasil registrasi")
            return True
        else:
            print_warning(f"Registrasi {name} mungkin sudah ada: {response.text}")
            return True  # User mungkin sudah ada
    except Exception as e:
        print_error(f"Error registrasi {name}: {str(e)}")
        return False


def login(email: str, password: str) -> Optional[str]:
    """Login dan dapatkan access token"""
    url = f"{BASE_URL}/api/v1/auth/login"
    payload = {
        "email": email,
        "password": password
    }
    
    try:
        response = requests.post(url, json=payload)
        if response.status_code == 200:
            data = response.json()
            if data.get("status") == "success":
                token = data["data"]["access_token"]
                print_success(f"Login berhasil untuk {email}")
                return token
        print_error(f"Login gagal: {response.text}")
        return None
    except Exception as e:
        print_error(f"Error login: {str(e)}")
        return None


def get_branches(token: str) -> Optional[str]:
    """Ambil list branches dan return branch_id pertama"""
    url = f"{BASE_URL}/api/v1/branches"
    headers = {"Authorization": f"Bearer {token}"}
    
    try:
        response = requests.get(url, headers=headers)
        if response.status_code == 200:
            data = response.json()
            if data.get("status") == "success" and data.get("data"):
                branches = data["data"]
                if len(branches) > 0:
                    branch_id = branches[0].get("id")
                    print_success(f"Branch ditemukan: {branches[0].get('name')}")
                    return branch_id
        print_warning("Tidak ada branch tersedia")
        return None
    except Exception as e:
        print_error(f"Error get branches: {str(e)}")
        return None


def create_order(token: str, branch_id: Optional[str] = None) -> Optional[Dict[str, Any]]:
    """Buat pesanan baru"""
    url = f"{BASE_URL}/api/v1/orders"
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "description": "iPhone layar retak - simulasi test",
        "complaint": "Layar iPhone retak setelah jatuh",
        "item_model": "iPhone 13 Pro",
        "item_color": "Blue",
        "item_serial": "IMEI123456789",
        "item_type": "iPhone",
        "pickup_address": "Jl. Kebon Jeruk No. 45, Jakarta Barat",
        "pickup_latitude": PICKUP_LOCATION["latitude"],
        "pickup_longitude": PICKUP_LOCATION["longitude"],
        "estimated_cost": 500000,
        "estimated_duration": 120
    }
    
    if branch_id:
        payload["branch_id"] = branch_id
    
    try:
        response = requests.post(url, json=payload, headers=headers)
        if response.status_code == 201:
            data = response.json()
            if data.get("status") == "success":
                order_data = data["data"]
                print_success(f"Pesanan berhasil dibuat: {order_data.get('order_number')}")
                return order_data
        print_error(f"Error create order: {response.text}")
        return None
    except Exception as e:
        print_error(f"Error create order: {str(e)}")
        return None


def assign_courier(order_id: str, courier_token: str) -> bool:
    """Assign courier ke pesanan"""
    url = f"{BASE_URL}/api/v1/orders/{order_id}/assign-courier"
    headers = {
        "Authorization": f"Bearer {courier_token}",
        "Content-Type": "application/json"
    }
    
    # Get courier user ID dari token (simplified - in real app, decode JWT)
    # Untuk simulasi, kita asumsikan courier_id sudah diketahui
    # Atau bisa ambil dari profile endpoint
    
    try:
        # Get courier profile untuk ambil user_id
        profile_url = f"{BASE_URL}/api/v1/auth/profile"
        profile_response = requests.get(profile_url, headers=headers)
        if profile_response.status_code == 200:
            profile_data = profile_response.json()
            courier_id = profile_data["data"]["id"]
            
            # Assign courier
            params = {"courier_id": courier_id}
            response = requests.post(url, params=params, headers=headers)
            if response.status_code == 200:
                print_success(f"Courier berhasil di-assign ke pesanan")
                return True
        print_warning(f"Assign courier gagal: {response.text if 'response' in locals() else 'Unknown'}")
        return False
    except Exception as e:
        print_error(f"Error assign courier: {str(e)}")
        return False


def calculate_distance(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
    """Hitung jarak antara dua koordinat (Haversine formula) dalam km"""
    R = 6371  # Radius bumi dalam km
    
    dlat = math.radians(lat2 - lat1)
    dlon = math.radians(lon2 - lon1)
    
    a = math.sin(dlat/2) * math.sin(dlat/2) + \
        math.cos(math.radians(lat1)) * math.cos(math.radians(lat2)) * \
        math.sin(dlon/2) * math.sin(dlon/2)
    
    c = 2 * math.atan2(math.sqrt(a), math.sqrt(1-a))
    distance = R * c
    
    return distance


def interpolate_location(start_lat: float, start_lon: float, 
                         end_lat: float, end_lon: float, 
                         progress: float) -> tuple:
    """Interpolasi lokasi antara start dan end berdasarkan progress (0-1)"""
    lat = start_lat + (end_lat - start_lat) * progress
    lon = start_lon + (end_lon - start_lon) * progress
    return lat, lon


def update_courier_location(order_id: str, courier_token: str, 
                            latitude: float, longitude: float,
                            speed: float = 30.0, heading: float = 0.0) -> bool:
    """Update lokasi courier"""
    url = f"{BASE_URL}/api/v1/orders/{order_id}/location"
    headers = {
        "Authorization": f"Bearer {courier_token}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "latitude": latitude,
        "longitude": longitude,
        "accuracy": 10.0,
        "speed": speed,
        "heading": heading
    }
    
    try:
        response = requests.post(url, json=payload, headers=headers)
        if response.status_code == 200:
            data = response.json()
            if data.get("status") == "success":
                location_data = data["data"]
                distance = location_data.get("distance", 0)
                eta = location_data.get("eta", 0)
                print_info(f"Lokasi updated: ({latitude:.6f}, {longitude:.6f}) | "
                          f"Jarak: {distance:.2f} km | ETA: {eta} menit")
                return True
        return False
    except Exception as e:
        print_error(f"Error update location: {str(e)}")
        return False


def get_current_location(order_id: str, customer_token: str) -> Optional[Dict[str, Any]]:
    """Ambil lokasi saat ini dari customer"""
    url = f"{BASE_URL}/api/v1/orders/{order_id}/location"
    headers = {"Authorization": f"Bearer {customer_token}"}
    
    try:
        response = requests.get(url, headers=headers)
        if response.status_code == 200:
            data = response.json()
            if data.get("status") == "success":
                return data["data"]
        return None
    except Exception as e:
        print_error(f"Error get location: {str(e)}")
        return None


def simulate_courier_movement(order_id: str, courier_token: str):
    """Simulasi pergerakan courier dari pickup location ke destination"""
    print_header("SIMULASI PERGERAKAN COURIER")
    
    # Start location (slightly away from pickup)
    start_lat = PICKUP_LOCATION["latitude"] - 0.01
    start_lon = PICKUP_LOCATION["longitude"] - 0.01
    
    # Destination
    dest_lat = DESTINATION_LOCATION["latitude"]
    dest_lon = DESTINATION_LOCATION["longitude"]
    
    # Simulasi pergerakan dalam 10 langkah
    steps = 10
    current_lat, current_lon = start_lat, start_lon
    
    print_info("Courier mulai perjalanan...")
    
    for step in range(steps + 1):
        progress = step / steps
        
        # Interpolasi lokasi
        current_lat, current_lon = interpolate_location(
            start_lat, start_lon, dest_lat, dest_lon, progress
        )
        
        # Hitung heading (arah pergerakan)
        if step < steps:
            next_lat, next_lon = interpolate_location(
                start_lat, start_lon, dest_lat, dest_lon, (step + 1) / steps
            )
            heading = math.degrees(math.atan2(
                next_lon - current_lon, next_lat - current_lat
            ))
        else:
            heading = 0.0
        
        # Update lokasi
        update_courier_location(
            order_id, courier_token, 
            current_lat, current_lon,
            speed=30.0, heading=heading
        )
        
        # Simpan update
        location_updates.append({
            "timestamp": datetime.now().isoformat(),
            "latitude": current_lat,
            "longitude": current_lon,
            "step": step,
            "progress": progress * 100
        })
        
        # Delay antar update
        time.sleep(2)
    
    print_success("Courier telah sampai di tujuan!")


def on_ws_message(ws, message):
    """Handle WebSocket message"""
    global ws_connected
    try:
        data = json.loads(message)
        msg_type = data.get("type", "unknown")
        
        if msg_type == "location":
            location_data = data.get("data", {})
            lat = location_data.get("latitude", 0)
            lon = location_data.get("longitude", 0)
            eta = location_data.get("eta", 0)
            distance = location_data.get("distance", 0)
            
            print_info(f"[WebSocket] Update Lokasi: ({lat:.6f}, {lon:.6f}) | "
                      f"Jarak: {distance:.2f} km | ETA: {eta} menit")
        
        elif msg_type == "status":
            status_data = data.get("data", {})
            status = status_data.get("status", "unknown")
            print_info(f"[WebSocket] Update Status: {status}")
        
        elif msg_type == "pong":
            print_info("[WebSocket] Pong received")
        
        else:
            print_info(f"[WebSocket] Message: {msg_type}")
    
    except Exception as e:
        print_error(f"Error parsing WebSocket message: {str(e)}")


def on_ws_error(ws, error):
    """Handle WebSocket error"""
    print_error(f"WebSocket error: {str(error)}")


def on_ws_close(ws, close_status_code, close_msg):
    """Handle WebSocket close"""
    global ws_connected
    ws_connected = False
    print_warning("WebSocket connection closed")


def on_ws_open(ws):
    """Handle WebSocket open"""
    global ws_connected
    ws_connected = True
    print_success("WebSocket connected!")
    
    # Send ping
    ping_msg = {
        "type": "ping",
        "timestamp": datetime.now().isoformat()
    }
    ws.send(json.dumps(ping_msg))


def start_websocket_tracking(order_id: str, customer_token: str):
    """Start WebSocket connection untuk real-time tracking"""
    print_header("WEBSOCKET REAL-TIME TRACKING")
    
    ws_url = f"{WS_URL}?order_id={order_id}"
    
    # Create WebSocket connection
    ws = websocket.WebSocketApp(
        ws_url,
        header={"Authorization": f"Bearer {customer_token}"},
        on_message=on_ws_message,
        on_error=on_ws_error,
        on_close=on_ws_close,
        on_open=on_ws_open
    )
    
    # Run WebSocket in separate thread
    ws_thread = threading.Thread(target=ws.run_forever)
    ws_thread.daemon = True
    ws_thread.start()
    
    # Wait for connection
    timeout = 5
    elapsed = 0
    while not ws_connected and elapsed < timeout:
        time.sleep(0.1)
        elapsed += 0.1
    
    if ws_connected:
        print_success("WebSocket ready untuk tracking!")
        return ws
    else:
        print_warning("WebSocket connection timeout")
        return None


def simulate_rest_api_tracking(order_id: str, customer_token: str, duration: int = 30):
    """Simulasi tracking via REST API polling"""
    print_header("REST API TRACKING (POLLING)")
    
    start_time = time.time()
    poll_count = 0
    
    while time.time() - start_time < duration:
        location = get_current_location(order_id, customer_token)
        
        if location:
            lat = location.get("latitude", 0)
            lon = location.get("longitude", 0)
            eta = location.get("eta", 0)
            distance = location.get("distance", 0)
            updated_at = location.get("updated_at", "")
            
            poll_count += 1
            print_info(f"[Poll #{poll_count}] Lokasi: ({lat:.6f}, {lon:.6f}) | "
                      f"Jarak: {distance:.2f} km | ETA: {eta} menit | "
                      f"Updated: {updated_at}")
        else:
            print_warning("Lokasi belum tersedia")
        
        time.sleep(3)  # Poll setiap 3 detik


def main():
    """Main function untuk simulasi"""
    global customer_token, courier_token, order_id, order_number
    
    print_header("SIMULASI TEST SISTEM REALTIME TRACKING")
    print_info("Memulai simulasi...")
    
    # 1. Setup: Registrasi dan Login
    print_header("TAHAP 1: SETUP USER")
    
    # Registrasi customer
    register_user(CUSTOMER_EMAIL, CUSTOMER_PASSWORD, "Customer Test", "pelanggan")
    time.sleep(1)
    
    # Registrasi courier
    register_user(COURIER_EMAIL, COURIER_PASSWORD, "Courier Test", "kurir")
    time.sleep(1)
    
    # Login customer
    customer_token = login(CUSTOMER_EMAIL, CUSTOMER_PASSWORD)
    if not customer_token:
        print_error("Gagal login customer. Pastikan server berjalan.")
        return
    
    # Login courier
    courier_token = login(COURIER_EMAIL, COURIER_PASSWORD)
    if not courier_token:
        print_error("Gagal login courier. Pastikan server berjalan.")
        return
    
    # 2. Buat Pesanan
    print_header("TAHAP 2: BUAT PESANAN")
    
    branch_id = get_branches(customer_token)
    order_data = create_order(customer_token, branch_id)
    
    if not order_data:
        print_error("Gagal membuat pesanan. Cek koneksi server.")
        return
    
    order_id = order_data.get("id")
    order_number = order_data.get("order_number")
    
    print_success(f"Order ID: {order_id}")
    print_success(f"Order Number: {order_number}")
    
    # 3. Assign Courier
    print_header("TAHAP 3: ASSIGN COURIER")
    
    if assign_courier(order_id, courier_token):
        print_success("Courier berhasil di-assign")
    else:
        print_warning("Assign courier gagal, lanjutkan simulasi...")
    
    # 4. Start WebSocket Tracking (Customer)
    print_header("TAHAP 4: WEBSOCKET TRACKING")
    
    ws = start_websocket_tracking(order_id, customer_token)
    time.sleep(2)
    
    # 5. Start REST API Tracking in separate thread
    print_header("TAHAP 5: REST API TRACKING")
    
    rest_thread = threading.Thread(
        target=simulate_rest_api_tracking,
        args=(order_id, customer_token, 30)
    )
    rest_thread.daemon = True
    rest_thread.start()
    
    # 6. Simulasi Pergerakan Courier
    print_header("TAHAP 6: SIMULASI PERGERAKAN")
    
    courier_thread = threading.Thread(
        target=simulate_courier_movement,
        args=(order_id, courier_token)
    )
    courier_thread.daemon = True
    courier_thread.start()
    
    # Wait for simulation to complete
    courier_thread.join(timeout=30)
    rest_thread.join(timeout=30)
    
    # 7. Summary
    print_header("RINGKASAN SIMULASI")
    
    print_success(f"Total location updates: {len(location_updates)}")
    print_success(f"Order ID: {order_id}")
    print_success(f"Order Number: {order_number}")
    
    if location_updates:
        print_info("\nLocation Update History:")
        for i, update in enumerate(location_updates[:5], 1):  # Show first 5
            print_info(f"  {i}. Step {update['step']}: "
                      f"({update['latitude']:.6f}, {update['longitude']:.6f}) | "
                      f"Progress: {update['progress']:.1f}%")
        if len(location_updates) > 5:
            print_info(f"  ... dan {len(location_updates) - 5} update lainnya")
    
    print_success("\nSimulasi selesai!")
    print_info("Tekan Ctrl+C untuk keluar...")
    
    # Keep WebSocket alive
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print_info("\nMenghentikan simulasi...")
        if ws:
            ws.close()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print_info("\nSimulasi dihentikan oleh user.")
    except Exception as e:
        print_error(f"Error: {str(e)}")
        import traceback
        traceback.print_exc()

