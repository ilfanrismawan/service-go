# API Test Script Documentation

## 📋 Overview

Script `api_tester.py` adalah script automated testing untuk iPhone Service API yang menguji berbagai endpoint dan menghasilkan laporan lengkap.

## 🚀 Quick Start

```bash
# Install dependencies
pip install requests

# Jalankan test
python scripts/api_tester.py
```

## 📊 Output

Setelah menjalankan test, laporan akan otomatis dibuat di folder `test-reports/`:

- `api-test-report.html` - HTML report (visual)
- `api-test-report.md` - Markdown report
- `LAPORAN-FINAL-TEST.md` - Comprehensive report in Indonesian
- `TEST-RESULTS-SUMMARY.txt` - Visual ASCII summary

## 🔧 Configuration

Edit variabel di awal file `api_tester.py` untuk mengubah konfigurasi:

```python
BASE_URL = "http://localhost:8080"
API_VERSION = "/api/v1"
TIMEOUT = 10  # seconds
```

## 📝 Test Coverage

Script ini menguji:

### ✅ Health Checks
- `/health` - Main health check
- `/health/live` - Liveness probe
- `/health/ready` - Readiness probe

### ✅ Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### ✅ Public Endpoints
- `GET /api/v1/branches` - List branches
- `GET /api/v1/branches/:id` - Get branch by ID
- `GET /api/v1/branches/nearest` - Get nearest branches

### ✅ Protected Endpoints (requires auth)
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders` - List orders
- `GET /api/v1/membership` - Get membership
- `GET /api/v1/membership/tiers` - Get membership tiers
- `GET /api/v1/reports/current-month` - Current month report
- `GET /api/v1/reports/monthly` - Monthly report

## 🎯 Features

- ✅ Automatic authentication
- ✅ JWT token management
- ✅ Response time measurement
- ✅ Color-coded terminal output
- ✅ Multiple report formats
- ✅ Error handling
- ✅ Test retry logic

## 📈 Example Output

```
======================================================================
            iPhone Service API - Comprehensive Test Suite
======================================================================

Base URL: http://localhost:8080
API Version: /api/v1

Checking server availability...
✓ Server is running!

======================================================================
                          HEALTH CHECK TESTS
======================================================================

→ Testing: GET /health
✓ Status: 200, Response: Health check completed
...
```

## 🔍 Troubleshooting

### Server not running
```bash
# Start Docker services
docker-compose up -d

# Check if server is running
curl http://localhost:8080/health
```

### Authentication fails
- Script will automatically try multiple credential combinations
- Check server logs for authentication issues
- Ensure database has been migrated

### Tests timeout
- Increase `TIMEOUT` value in script
- Check server performance
- Check network connectivity

## 🛠️ Customization

### Add New Test

Edit `api_tester.py` and add new test function:

```python
def test_new_endpoint():
    print_header("NEW ENDPOINT TESTS")
    
    print_test("GET /api/v1/new-endpoint")
    status, response, rt = make_request("GET", f"{API_VERSION}/new-endpoint")
    passed = status == 200
    
    record_test("New Endpoint", f"{API_VERSION}/new-endpoint", "GET", 
                status, passed, "Test message", rt, response)
    
    if passed:
        print_success(f"Status: {status}")
    else:
        print_error(f"Status: {status}")

# Call in main()
test_new_endpoint()
```

### Modify Test Credentials

Edit test credentials in `test_auth()` function:

```python
test_credentials = [
    {"email": "your-test@email.com", "password": "password123", "register": True},
    # Add more credentials...
]
```

## 📚 Requirements

- Python 3.7+
- requests library
- API server running on configured port

## 📞 Support

Untuk bantuan lebih lanjut:
- Lihat dokumentasi API di `README.md`
- Check API examples di `API_EXAMPLES.md`
- Review test reports di `test-reports/`

---

**Version:** 1.0  
**Last Updated:** October 27, 2025

