# Comprehensive API Test Report
Generated: 2025-10-27 11:08:35

## Summary


- **Total Tests:** 20
- **Passed:** 19 (95.0%)
- **Failed:** 1 (5.0%)
- **Success Rate:** 95.0%

## Test Results

### ✅ Admin Login

- **Endpoint:** /api/v1/auth/login
- **Method:** POST
- **Status:** 200
- **Time:** 0.089s
- **Message:** Admin logged in

### ✅ User Login

- **Endpoint:** /api/v1/auth/login
- **Method:** POST
- **Status:** 200
- **Time:** 0.093s
- **Message:** User logged in

### ✅ Get Profile

- **Endpoint:** /api/v1/auth/profile
- **Method:** GET
- **Status:** 200
- **Time:** 0.013s
- **Message:** Profile retrieved

### ✅ Get Orders

- **Endpoint:** /api/v1/orders
- **Method:** GET
- **Status:** 200
- **Time:** 0.027s
- **Message:** Orders retrieved

### ✅ Get Membership

- **Endpoint:** /api/v1/membership
- **Method:** GET
- **Status:** 404
- **Time:** 0.015s
- **Message:** Membership retrieved

### ✅ Get Current Month Report

- **Endpoint:** /api/v1/reports/current-month
- **Method:** GET
- **Status:** 200
- **Time:** 0.035s
- **Message:** Report retrieved

### ✅ RBAC: User Cannot Access Admin

- **Endpoint:** /api/v1/admin/users
- **Method:** GET
- **Status:** 403
- **Time:** 0.016s
- **Message:** Access denied (expected)

### ✅ RBAC: User Can Access Own Endpoints

- **Endpoint:** /api/v1/auth/profile
- **Method:** GET
- **Status:** 200
- **Time:** 0.009s
- **Message:** Access granted

### ✅ RBAC: Admin Can Access Admin Endpoints

- **Endpoint:** /api/v1/admin/dashboard
- **Method:** GET
- **Status:** 200
- **Time:** 0.023s
- **Message:** Access granted

### ✅ Error: Invalid Login

- **Endpoint:** /api/v1/auth/login
- **Method:** POST
- **Status:** 401
- **Time:** 0.015s
- **Message:** Properly rejected

### ✅ Error: Missing Fields

- **Endpoint:** /api/v1/auth/register
- **Method:** POST
- **Status:** 400
- **Time:** 0.010s
- **Message:** Validation error returned

### ✅ Error: No Token

- **Endpoint:** /api/v1/orders
- **Method:** GET
- **Status:** 401
- **Time:** 0.010s
- **Message:** Unauthorized (expected)

### ✅ Error: Invalid Token

- **Endpoint:** /api/v1/orders
- **Method:** GET
- **Status:** 401
- **Time:** 0.034s
- **Message:** Private rejected

### ✅ Error: 404 Not Found

- **Endpoint:** /api/v1/orders/:id
- **Method:** GET
- **Status:** 404
- **Time:** 0.011s
- **Message:** Not found (expected)

### ✅ Perf: GET /health

- **Endpoint:** /health
- **Method:** GET
- **Status:** 200
- **Time:** 0.011s
- **Message:** Avg: 0.011s (Excellent)

### ✅ Perf: GET /branches

- **Endpoint:** /api/v1/branches
- **Method:** GET
- **Status:** 200
- **Time:** 0.009s
- **Message:** Avg: 0.009s (Excellent)

### ✅ Perf: GET /orders

- **Endpoint:** /api/v1/orders
- **Method:** GET
- **Status:** 200
- **Time:** 0.012s
- **Message:** Avg: 0.012s (Excellent)

### ✅ Security: SQL Injection

- **Endpoint:** /api/v1/auth/login
- **Method:** POST
- **Status:** 400
- **Time:** 0.009s
- **Message:** Properly handled

### ❌ Security: XSS Attempt

- **Endpoint:** /api/v1/auth/register
- **Method:** POST
- **Status:** 409
- **Time:** 0.010s
- **Message:** Should sanitize input

### ✅ Security: Malformed Token

- **Endpoint:** /api/v1/orders
- **Method:** GET
- **Status:** 401
- **Time:** 0.007s
- **Message:** Rejected

