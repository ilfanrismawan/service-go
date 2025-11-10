#!/usr/bin/env python3
"""
Advanced Comprehensive API Testing Script
Covers: Protected endpoints, Admin endpoints, RBAC, Error scenarios, Performance, Security
"""

import requests
import json
import time
from datetime import datetime
from typing import Dict, List, Any
import sys
import hashlib

# Configuration
BASE_URL = "http://localhost:8080"
API_VERSION = "/api/v1"
TIMEOUT = 10

# Test results
test_results: List[Dict[str, Any]] = []
admin_token: str = ""
user_token: str = ""
admin_id: str = ""
user_id: str = ""

class Colors:
    """ANSI color codes"""
    GREEN = '\033[92m'
    RED = '\033[91m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    CYAN = '\033[96m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'

def print_header(text: str):
    print(f"\n{'='*70}")
    print(f"{text:^70}")
    print(f"{'='*70}\n")

def print_test(desc: str):
    print(f"{Colors.CYAN}→ {desc}{Colors.ENDC}")

def print_success(msg: str):
    print(f"{Colors.GREEN}✓ {msg}{Colors.ENDC}")

def print_error(msg: str):
    print(f"{Colors.RED}✗ {msg}{Colors.ENDC}")

def print_warning(msg: str):
    print(f"{Colors.YELLOW}⚠ {msg}{Colors.ENDC}")

def record_test(name: str, endpoint: str, method: str, status: int, passed: bool, msg: str, rt: float, data: Any = None):
    global test_results
    test_results.append({
        "test_name": name,
        "endpoint": endpoint,
        "method": method,
        "status_code": status,
        "passed": passed,
        "message": msg,
        "response_time": rt,
        "timestamp": datetime.now().isoformat(),
        "response": data
    })

def make_request(method: str, endpoint: str, headers: Dict = None, data: Dict = None, params: Dict = None, token: str = None) -> tuple:
    url = f"{BASE_URL}{endpoint}"
    
    if headers is None:
        headers = {"Content-Type": "application/json"}
    
    if token:
        headers["Authorization"] = f"Bearer {token}"
    
    start_time = time.time()
    
    try:
        if method == "GET":
            response = requests.get(url, headers=headers, params=params, timeout=TIMEOUT)
        elif method == "POST":
            response = requests.post(url, headers=headers, json=data, timeout=TIMEOUT)
        elif method == "PUT":
            response = requests.put(url, headers=headers, json=data, timeout=TIMEOUT)
        elif method == "DELETE":
            response = requests.delete(url, headers=headers, timeout=TIMEOUT)
        else:
            return 0, {}, 0
        
        response_time = time.time() - start_time
        return response.status_code, response.json() if response.text else {}, response_time
    except requests.exceptions.RequestException as e:
        response_time = time.time() - start_time
        return 0, {"error": str(e)}, response_time

# ==================== SETUP TEST USERS ====================
def setup_test_users():
    """Setup test users with different roles"""
    global admin_token, user_token, admin_id, user_id
    
    print_header("SETUP TEST USERS")
    
    # Register admin user
    print_test("Registering admin user")
    admin_email = f"admin_test_{int(time.time())}@test.com"
    admin_data = {
        "email": admin_email,
        "password": "Admin123!",
        "full_name": "Admin User",
        "phone": f"081234{int(time.time())%100000}",
        "role": "admin_pusat"
    }
    
    status, response, rt = make_request("POST", f"{API_VERSION}/auth/register", data=admin_data)
    if status in [200, 201]:
        print_success("Admin user registered")
    elif status == 409:
        print_warning("Admin user already exists")
    else:
        print_error(f"Failed to register admin: {status}")
    
    # Login as admin
    print_test("Logging in as admin")
    status, response, rt = make_request("POST", f"{API_VERSION}/auth/login", data={
        "email": admin_email,
        "password": "Admin123!"
    })
    
    if status == 200 and "data" in response and "access_token" in response["data"]:
        admin_token = response["data"]["access_token"]
        admin_id = response["data"].get("user", {}).get("id", "")
        print_success("Admin login successful")
        record_test("Admin Login", f"{API_VERSION}/auth/login", "POST", status, True, "Admin logged in", rt, response)
    else:
        print_error("Admin login failed")
        record_test("Admin Login", f"{API_VERSION}/auth/login", "POST", status, False, "Admin login failed", rt, response)
    
    # Register regular user
    print_test("Registering regular user")
    user_email = f"user_test_{int(time.time())}@test.com"
    user_data = {
        "email": user_email,
        "password": "User123!",
        "full_name": "Regular User",
        "phone": f"081235{int(time.time())%100000}",
        "role": "pelanggan"
    }
    
    status, response, rt = make_request("POST", f"{API_VERSION}/auth/register", data=user_data)
    if status in [200, 201]:
        print_success("Regular user registered")
    elif status == 409:
        print_warning("User already exists")
    
    # Login as regular user
    print_test("Logging in as regular user")
    status, response, rt = make_request("POST", f"{API_VERSION}/auth/login", data={
        "email": user_email,
        "password": "User123!"
    })
    
    if status == 200 and "data" in response and "access_token" in response["data"]:
        user_token = response["data"]["access_token"]
        user_id = response["data"].get("user", {}).get("id", "")
        print_success("User login successful")
        record_test("User Login", f"{API_VERSION}/auth/login", "POST", status, True, "User logged in", rt, response)
    else:
        print_error("User login failed")
        record_test("User Login", f"{API_VERSION}/auth/login", "POST", status, False, "User login failed", rt, response)

# ==================== TEST PROTECTED ENDPOINTS ====================
def test_protected_endpoints():
    """Test endpoints that require authentication"""
    print_header("TESTING PROTECTED ENDPOINTS")
    
    if not user_token:
        print_error("No user token available, skipping protected endpoints test")
        return
    
    # Test Get Profile
    print_test("GET /auth/profile")
    status, response, rt = make_request("GET", f"{API_VERSION}/auth/profile", token=user_token)
    passed = status == 200
    record_test("Get Profile", f"{API_VERSION}/auth/profile", "GET", status, passed, 
                "Profile retrieved" if passed else "Failed to get profile", rt, response)
    if passed:
        print_success(f"Status: {status}")
    else:
        print_error(f"Status: {status}")
    
    # Test Get Orders
    print_test("GET /orders")
    status, response, rt = make_request("GET", f"{API_VERSION}/orders", token=user_token)
    passed = status == 200
    record_test("Get Orders", f"{API_VERSION}/orders", "GET", status, passed,
                "Orders retrieved" if passed else "Failed to get orders", rt, response)
    if passed:
        print_success(f"Status: {status}")
    else:
        print_error(f"Status: {status}")
    
    # Test Get Membership
    print_test("GET /membership")
    status, response, rt = make_request("GET", f"{API_VERSION}/membership", token=user_token)
    passed = status in [200, 404]
    record_test("Get Membership", f"{API_VERSION}/membership", "GET", status, passed,
                "Membership retrieved" if passed else "Failed", rt, response)
    if passed:
        print_success(f"Status: {status}")
    
    # Test Get Reports
    print_test("GET /reports/current-month")
    status, response, rt = make_request("GET", f"{API_VERSION}/reports/current-month", token=user_token)
    passed = status == 200
    record_test("Get Current Month Report", f"{API_VERSION}/reports/current-month", "GET", status, passed,
                "Report retrieved" if passed else "Failed", rt, response)
    if passed:
        print_success(f"Status: {status}")
    else:
        print_error(f"Status: {status}")

# ==================== TEST RBAC ====================
def test_rbac():
    """Test Role-Based Access Control"""
    print_header("TESTING ROLE-BASED ACCESS CONTROL")
    
    if not admin_token or not user_token:
        print_error("Missing tokens for RBAC testing")
        return
    
    # Test 1: User should NOT access admin endpoints
    print_test("User accessing admin endpoint (should fail)")
    status, response, rt = make_request("GET", f"{API_VERSION}/admin/users", token=user_token)
    passed = status in [401, 403, 404]
    record_test("RBAC: User Cannot Access Admin", f"{API_VERSION}/admin/users", "GET", status, passed,
                "Access denied (expected)" if passed else "Security issue: User accessed admin endpoint", rt, response)
    if passed:
        print_success(f"Status: {status} - Access properly denied")
    else:
        print_error(f"Status: {status} - Security issue!")
    
    # Test 2: User CAN access their own endpoints
    print_test("User accessing their own protected endpoint")
    status, response, rt = make_request("GET", f"{API_VERSION}/auth/profile", token=user_token)
    passed = status == 200
    record_test("RBAC: User Can Access Own Endpoints", f"{API_VERSION}/auth/profile", "GET", status, passed,
                "Access granted" if passed else "Access denied unexpectedly", rt, response)
    if passed:
        print_success(f"Status: {status} - Access granted")
    else:
        print_error(f"Status: {status} - Access denied unexpectedly")
    
    # Test 3: Admin CAN access admin endpoints
    if admin_token:
        print_test("Admin accessing admin endpoint")
        status, response, rt = make_request("GET", f"{API_VERSION}/admin/dashboard", token=admin_token)
        passed = status in [200, 404]  # 404 if not implemented, but 401/403 would be security issue
        record_test("RBAC: Admin Can Access Admin Endpoints", f"{API_VERSION}/admin/dashboard", "GET", status, passed,
                    "Access granted" if passed else "Access denied", rt, response)
        if passed:
            print_success(f"Status: {status}")
        else:
            print_error(f"Status: {status}")

# ==================== TEST ERROR SCENARIOS ====================
def test_error_scenarios():
    """Test error handling and edge cases"""
    print_header("TESTING ERROR SCENARIOS")
    
    # Test 1: Invalid credentials
    print_test("Login with invalid credentials")
    status, response, rt = make_request("POST", f"{API_VERSION}/auth/login", data={
        "email": "invalid@test.com",
        "password": "wrongpass"
    })
    passed = status == 401
    record_test("Error: Invalid Login", f"{API_VERSION}/auth/login", "POST", status, passed,
                "Properly rejected" if passed else "Should return 401", rt, response)
    if passed:
        print_success(f"Status: {status} - Properly rejected")
    else:
        print_error(f"Status: {status}")
    
    # Test 2: Missing required fields
    print_test("Register with missing fields")
    status, response, rt = make_request("POST", f"{API_VERSION}/auth/register", data={
        "email": "incomplete@test.com"
        # Missing required fields
    })
    passed = status in [400, 422]
    record_test("Error: Missing Fields", f"{API_VERSION}/auth/register", "POST", status, passed,
                "Validation error returned" if passed else "Should return 400", rt, response)
    if passed:
        print_success(f"Status: {status} - Validation error")
    else:
        print_error(f"Status: {status}")
    
    # Test 3: Access without token
    print_test("Access protected endpoint without token")
    status, response, rt = make_request("GET", f"{API_VERSION}/orders")
    passed = status in [401, 403]
    record_test("Error: No Token", f"{API_VERSION}/orders", "GET", status, passed,
                "Unauthorized (expected)" if passed else "Should require auth", rt, response)
    if passed:
        print_success(f"Status: {status} - Unauthorized")
    else:
        print_error(f"Status: {status} - Should require auth")
    
    # Test 4: Invalid token
    print_test("Access with invalid token")
    status, response, rt = make_request("GET", f"{API_VERSION}/orders", token="invalid.token.here")
    passed = status in [401, 403]
    record_test("Error: Invalid Token", f"{API_VERSION}/orders", "GET", status, passed,
                "Private rejected" if passed else "Should reject invalid token", rt, response)
    if passed:
        print_success(f"Status: {status} - Invalid token rejected")
    else:
        print_error(f"Status: {status}")
    
    # Test 5: Non-existent resource
    print_test("Get non-existent order")
    if user_token:
        status, response, rt = make_request("GET", f"{API_VERSION}/orders/00000000-0000-0000-0000-000000000000", token=user_token)
        passed = status == 404
        record_test("Error: 404 Not Found", f"{API_VERSION}/orders/:id", "GET", status, passed,
                    "Not found (expected)" if passed else "Should return 404", rt, response)
        if passed:
            print_success(f"Status: {status} - Proper 404")
        else:
            print_error(f"Status: {status}")

# ==================== PERFORMANCE TESTING ====================
def performance_test():
    """Test API performance"""
    print_header("PERFORMANCE TESTING")
    
    # Test response times for critical endpoints
    endpoints_to_test = [
        ("GET /health", "/health", "GET"),
        ("GET /branches", f"{API_VERSION}/branches", "GET"),
        ("GET /orders", f"{API_VERSION}/orders", "GET"),
    ]
    
    for name, endpoint, method in endpoints_to_test:
        print_test(f"Performance: {name}")
        times = []
        
        for i in range(5):
            token = user_token if "orders" in endpoint else None
            status, response, rt = make_request(method, endpoint, token=token)
            times.append(rt)
            time.sleep(0.1)  # Small delay between requests
        
        avg_time = sum(times) / len(times)
        max_time = max(times)
        min_time = min(times)
        
        print_success(f"Avg: {avg_time:.3f}s, Min: {min_time:.3f}s, Max: {max_time:.3f}s")
        
        # Determine performance rating
        if avg_time < 0.1:
            rating = "Excellent"
        elif avg_time < 0.3:
            rating = "Good"
        elif avg_time < 0.5:
            rating = "Acceptable"
        else:
            rating = "Slow"
        
        record_test(f"Perf: {name}", endpoint, method, 200, avg_time < 1.0, 
                    f"Avg: {avg_time:.3f}s ({rating})", avg_time, {"times": times, "min": min_time, "max": max_time})

# ==================== SECURITY TESTING ====================
def security_test():
    """Test security aspects"""
    print_header("SECURITY TESTING")
    
    # Test 1: SQL Injection attempt
    print_test("SQL Injection attempt in email")
    status, response, rt = make_request("POST", f"{API_VERSION}/auth/login", data={
        "email": "admin@test.com'; DROP TABLE users; --",
        "password": "test"
    })
    passed = status in [401, 400, 422]  # Should reject or return 401, not 500
    record_test("Security: SQL Injection", f"{API_VERSION}/auth/login", "POST", status, passed,
                "Properly handled" if passed else "Potential SQL injection vulnerability", rt, response)
    if passed:
        print_success(f"Status: {status} - Properly handled")
    else:
        print_error(f"Status: {status} - Potential vulnerability!")
    
    # Test 2: XSS attempt
    print_test("XSS attempt in registration")
    status, response, rt = make_request("POST", f"{API_VERSION}/auth/register", data={
        "email": f"xss_{int(time.time())}@test.com",
        "password": "Test123!",
        "full_name": "<script>alert('xss')</script>",
        "phone": "081234567890",
        "role": "pelanggan"
    })
    # Should either reject or sanitize
    passed = status in [400, 422, 201]
    record_test("Security: XSS Attempt", f"{API_VERSION}/auth/register", "POST", status, passed,
                "Handled properly" if passed else "Should sanitize input", rt, response)
    if passed:
        print_success(f"Status: {status} - Properly handled")
    else:
        print_error(f"Status: {status}")
    
    # Test 3: Token validation
    print_test("Access with malformed token")
    status, response, rt = make_request("GET", f"{API_VERSION}/orders", token="malformed.token")
    passed = status in [401, 403]
    record_test("Security: Malformed Token", f"{API_VERSION}/orders", "GET", status, passed,
                "Rejected" if passed else "Should reject malformed token", rt, response)
    if passed:
        print_success(f"Status: {status} - Malformed token rejected")
    else:
        print_error(f"Status: {status}")
    
    # Test 4: Sensitive data exposure
    print_test("Check if password is in response")
    if user_token:
        status, response, rt = make_request("GET", f"{API_VERSION}/auth/profile", token=user_token)
        if status == 200 and "password" not in json.dumps(response).lower():
            print_success("Password not exposed in response")
        else:
            print_error("Potential password exposure!")

# ==================== GENERATE REPORT ====================
def generate_report():
    """Generate comprehensive report"""
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    
    report = f"""# Comprehensive API Test Report
Generated: {timestamp}

## Summary

"""
    
    total = len(test_results)
    passed = sum(1 for t in test_results if t["passed"])
    failed = total - passed
    
    report += f"""
- **Total Tests:** {total}
- **Passed:** {passed} ({passed/total*100:.1f}%)
- **Failed:** {failed} ({failed/total*100:.1f}%)
- **Success Rate:** {passed/total*100:.1f}%

## Test Results

"""
    
    for result in test_results:
        status_icon = "✅" if result["passed"] else "❌"
        report += f"""### {status_icon} {result['test_name']}

- **Endpoint:** {result['endpoint']}
- **Method:** {result['method']}
- **Status:** {result['status_code']}
- **Time:** {result['response_time']:.3f}s
- **Message:** {result['message']}

"""
    
    # Save report
    with open("test-reports/advanced-test-report.md", "w", encoding="utf-8") as f:
        f.write(report)
    
    print_success("Advanced test report generated: test-reports/advanced-test-report.md")

# ==================== MAIN ====================
def main():
    print_header("ADVANCED COMPREHENSIVE API TEST SUITE")
    
    # Check server
    try:
        response = requests.get(f"{BASE_URL}/health", timeout=5)
        if response.status_code == 200:
            print_success("Server is running!")
        else:
            print_error("Server not responding properly")
            sys.exit(1)
    except:
        print_error("Cannot connect to server")
        sys.exit(1)
    
    # Run all tests
    setup_test_users()
    test_protected_endpoints()
    test_rbac()
    test_error_scenarios()
    performance_test()
    security_test()
    
    # Generate report
    generate_report()
    
    # Summary
    print_header("TEST SUMMARY")
    total = len(test_results)
    passed = sum(1 for t in test_results if t["passed"])
    avg_time = sum(t["response_time"] for t in test_results) / total if total > 0 else 0
    
    print(f"Total Tests: {total}")
    print(f"Passed: {passed} ({passed/total*100:.1f}%)")
    print(f"Failed: {total - passed}")
    print(f"Average Response Time: {avg_time:.3f}s")
    
    print_header("TESTING COMPLETE")

if __name__ == "__main__":
    main()

