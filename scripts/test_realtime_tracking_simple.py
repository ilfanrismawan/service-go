#!/usr/bin/env python3
"""
Simple Test Script untuk Realtime Tracking
===========================================

Script sederhana untuk test realtime tracking dengan minimal dependencies.
Hanya menggunakan requests library (no WebSocket untuk simplicity).

Usage:
    python scripts/test_realtime_tracking_simple.py
"""

import requests
import json
import time
from datetime import datetime

BASE_URL = "http://localhost:8080"

def print_status(message, status="info"):
    """Print status message dengan format"""
    symbols = {
        "success": "✓",
        "error": "✗",
        "info": "ℹ",
        "warning": "⚠"
    }
    symbol = symbols.get(status, "•")
    print(f"{symbol} {message}")

def test_realtime_tracking():
    """Test realtime tracking dengan REST API"""
    print("\n" + "="*60)
    print("SIMPLE REALTIME TRACKING TEST")
    print("="*60 + "\n")
    
    # 1. Login Customer
    print_status("1. Login Customer...", "info")
    login_url = f"{BASE_URL}/api/v1/auth/login"
    login_data = {
        "email": "customer@test.com",
        "password": "password123"
    }
    
    try:
        response = requests.post(login_url, json=login_data)
        if response.status_code != 200:
            print_status(f"Login gagal: {response.text}", "error")
            return False
        
        data = response.json()
        if data.get("status") != "success":
            print_status(f"Login gagal: {data.get('message')}", "error")
            return False
        
        customer_token = data["data"]["access_token"]
        print_status("Login customer berhasil", "success")
    except Exception as e:
        print_status(f"Error login: {str(e)}", "error")
        return False
    
    # 2. Login Courier
    print_status("2. Login Courier...", "info")
    courier_data = {
        "email": "courier@test.com",
        "password": "password123"
    }
    
    try:
        response = requests.post(login_url, json=courier_data)
        if response.status_code != 200:
            print_status(f"Login courier gagal: {response.text}", "error")
            return False
        
        data = response.json()
        if data.get("status") != "success":
            print_status(f"Login courier gagal: {data.get('message')}", "error")
            return False
        
        courier_token = data["data"]["access_token"]
        print_status("Login courier berhasil", "success")
    except Exception as e:
        print_status(f"Error login courier: {str(e)}", "error")
        return False
    
    # 3. Buat Pesanan
    print_status("3. Buat Pesanan...", "info")
    order_url = f"{BASE_URL}/api/v1/orders"
    headers = {
        "Authorization": f"Bearer {customer_token}",
        "Content-Type": "application/json"
    }
    
    order_data = {
        "description": "iPhone layar retak - test tracking",
        "complaint": "Layar iPhone retak",
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
    
    try:
        response = requests.post(order_url, json=order_data, headers=headers)
        if response.status_code != 201:
            print_status(f"Buat pesanan gagal: {response.text}", "error")
            return False
        
        data = response.json()
        if data.get("status") != "success":
            print_status(f"Buat pesanan gagal: {data.get('message')}", "error")
            return False
        
        order_info = data["data"]
        order_id = order_info["id"]
        order_number = order_info["order_number"]
        print_status(f"Pesanan berhasil dibuat: {order_number}", "success")
        print_status(f"Order ID: {order_id}", "info")
    except Exception as e:
        print_status(f"Error buat pesanan: {str(e)}", "error")
        return False
    
    # 4. Update Lokasi Courier (Simulasi)
    print_status("4. Simulasi Update Lokasi Courier...", "info")
    location_url = f"{BASE_URL}/api/v1/orders/{order_id}/location"
    courier_headers = {
        "Authorization": f"Bearer {courier_token}",
        "Content-Type": "application/json"
    }
    
    # Simulasi beberapa update lokasi
    locations = [
        {"lat": -6.2088, "lon": 106.8456, "step": 1},
        {"lat": -6.2040, "lon": 106.8350, "step": 2},
        {"lat": -6.1990, "lon": 106.8280, "step": 3},
        {"lat": -6.1944, "lon": 106.8229, "step": 4},  # Destination
    ]
    
    for loc in locations:
        location_data = {
            "latitude": loc["lat"],
            "longitude": loc["lon"],
            "accuracy": 10.0,
            "speed": 30.0,
            "heading": 45.0
        }
        
        try:
            response = requests.post(location_url, json=location_data, headers=courier_headers)
            if response.status_code == 200:
                data = response.json()
                if data.get("status") == "success":
                    loc_info = data["data"]
                    distance = loc_info.get("distance", 0)
                    eta = loc_info.get("eta", 0)
                    print_status(f"Step {loc['step']}: Lokasi updated | "
                                f"Jarak: {distance:.2f} km | ETA: {eta} menit", "success")
            else:
                print_status(f"Update lokasi gagal: {response.text}", "warning")
        except Exception as e:
            print_status(f"Error update lokasi: {str(e)}", "error")
        
        time.sleep(2)  # Delay antar update
    
    # 5. Customer Track Lokasi
    print_status("5. Customer Track Lokasi...", "info")
    track_url = f"{BASE_URL}/api/v1/orders/{order_id}/location"
    customer_headers = {
        "Authorization": f"Bearer {customer_token}"
    }
    
    try:
        response = requests.get(track_url, headers=customer_headers)
        if response.status_code == 200:
            data = response.json()
            if data.get("status") == "success":
                location = data["data"]
                lat = location.get("latitude", 0)
                lon = location.get("longitude", 0)
                distance = location.get("distance", 0)
                eta = location.get("eta", 0)
                updated_at = location.get("updated_at", "")
                
                print_status(f"Lokasi saat ini: ({lat:.6f}, {lon:.6f})", "success")
                print_status(f"Jarak ke tujuan: {distance:.2f} km", "info")
                print_status(f"ETA: {eta} menit", "info")
                print_status(f"Last update: {updated_at}", "info")
            else:
                print_status(f"Track lokasi gagal: {data.get('message')}", "error")
        else:
            print_status(f"Track lokasi gagal: {response.text}", "error")
    except Exception as e:
        print_status(f"Error track lokasi: {str(e)}", "error")
    
    # 6. Lihat History Lokasi
    print_status("6. Lihat History Lokasi...", "info")
    history_url = f"{BASE_URL}/api/v1/orders/{order_id}/location/history?limit=10"
    
    try:
        response = requests.get(history_url, headers=customer_headers)
        if response.status_code == 200:
            data = response.json()
            if data.get("status") == "success":
                history = data["data"]
                print_status(f"History lokasi: {len(history)} records", "success")
                for i, loc in enumerate(history[:3], 1):  # Show first 3
                    lat = loc.get("latitude", 0)
                    lon = loc.get("longitude", 0)
                    timestamp = loc.get("timestamp", "")
                    print_status(f"  {i}. ({lat:.6f}, {lon:.6f}) - {timestamp}", "info")
            else:
                print_status(f"History lokasi gagal: {data.get('message')}", "error")
        else:
            print_status(f"History lokasi gagal: {response.text}", "error")
    except Exception as e:
        print_status(f"Error history lokasi: {str(e)}", "error")
    
    # 7. Summary
    print("\n" + "="*60)
    print_status("TEST SELESAI", "success")
    print_status(f"Order ID: {order_id}", "info")
    print_status(f"Order Number: {order_number}", "info")
    print("="*60 + "\n")
    
    return True

if __name__ == "__main__":
    try:
        success = test_realtime_tracking()
        if success:
            print_status("Semua test berhasil!", "success")
        else:
            print_status("Beberapa test gagal. Cek output di atas.", "error")
    except KeyboardInterrupt:
        print_status("\nTest dihentikan oleh user.", "warning")
    except Exception as e:
        print_status(f"Error: {str(e)}", "error")
        import traceback
        traceback.print_exc()

