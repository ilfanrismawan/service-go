#!/usr/bin/env python3
"""
Stress Testing Script for iPhone Service API
Tests API with high concurrent load
"""

import requests
import concurrent.futures
import time
import statistics
from datetime import datetime
from typing import List, Tuple

BASE_URL = "http://localhost:8080"
API_VERSION = "/api/v1"

class Colors:
    GREEN = '\033[92m'
    RED = '\033[91m'
    YELLOW = '\033[93m'
    CYAN = '\033[96m'
    ENDC = '\033[0m'

def print_header(text: str):
    print(f"\n{Colors.CYAN}{'='*70}{Colors.ENDC}")
    print(f"{Colors.CYAN}{text:^70}{Colors.ENDC}")
    print(f"{Colors.CYAN}{'='*70}{Colors.ENDC}\n")

def print_success(msg: str):
    print(f"{Colors.GREEN}✓ {msg}{Colors.ENDC}")

def print_error(msg: str):
    print(f"{Colors.RED}✗ {msg}{Colors.ENDC}")

def print_info(msg: str):
    print(f"{Colors.YELLOW}→ {msg}{Colors.ENDC}")

def make_request(url: str) -> Tuple[int, float]:
    """Make a single request and return status code and time"""
    start = time.time()
    try:
        response = requests.get(url, timeout=10)
        elapsed = time.time() - start
        return response.status_code, elapsed
    except Exception as e:
        elapsed = time.time() - start
        print_error(f"Request failed: {e}")
        return 0, elapsed

def stress_test(endpoint: str, num_requests: int = 1000, workers: int = 50):
    """Run stress test on an endpoint"""
    print_header(f"Stress Testing: {endpoint}")
    print_info(f"Total Requests: {num_requests}")
    print_info(f"Concurrent Workers: {workers}")
    print_info("Starting stress test...\n")
    
    url = f"{BASE_URL}{endpoint}"
    results = []
    
    start_time = time.time()
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=workers) as executor:
        futures = [executor.submit(make_request, url) for _ in range(num_requests)]
        
        for future in concurrent.futures.as_completed(futures):
            status, elapsed = future.result()
            results.append((status, elapsed))
    
    total_time = time.time() - start_time
    
    # Calculate statistics
    response_times = [r[1] for r in results]
    status_codes = [r[0] for r in results]
    
    successful = sum(1 for s in status_codes if 200 <= s < 300)
    failed = len(results) - successful
    
    success_rate = (successful / len(results)) * 100
    
    stats = {
        "total_requests": len(results),
        "successful": successful,
        "failed": failed,
        "success_rate": success_rate,
        "avg_response_time": statistics.mean(response_times),
        "median_response_time": statistics.median(response_times),
        "min_response_time": min(response_times),
        "max_response_time": max(response_times),
        "p95_response_time": sorted(response_times)[int(len(response_times) * 0.95)] if response_times else 0,
        "p99_response_time": sorted(response_times)[int(len(response_times) * 0.99)] if response_times else 0,
        "requests_per_second": len(results) / total_time if total_time > 0 else 0,
        "total_time": total_time
    }
    
    # Print results
    print("\nResults:")
    print(f"  Total Requests: {stats['total_requests']}")
    print(f"  Successful: {stats['successful']} ({stats['success_rate']:.2f}%)")
    print(f"  Failed: {stats['failed']}")
    print(f"  Average Response Time: {stats['avg_response_time']:.3f}s")
    print(f"  Median Response Time: {stats['median_response_time']:.3f}s")
    print(f"  Min Response Time: {stats['min_response_time']:.3f}s")
    print(f"  Max Response Time: {stats['max_response_time']:.3f}s")
    print(f"  P95 Response Time: {stats['p95_response_time']崩溃:.3f}s")
    print(f"  P99 Response Time: {stats['p99_response_time']:.3f}s")
    print(f"  Requests Per Second: {stats['requests_per_second']:.2f}")
    print(f"  Total Time: {stats['total_time']:.2f}s")
    
    # Determine status
    if success_rate >= 99:
        print_success(f"Stress test PASSED - Success rate: {success_rate:.2f}%")
    elif success_rate >= 95:
        print_error(f"Stress test WARNING - Success rate: {success_rate:.2f}%")
    else:
        print_error(f"Stress test FAILED - Success rate: {success_rate:.2f}%")
    
    return stats

def main():
    print_header("IPHONE SERVICE API - STRESS TEST")
    
    # Check if server is running
    try:
        response = requests.get(f"{BASE_URL}/health", timeout=5)
        if response.status_code == 200:
            print_success("Server is running!")
        else:
            print_error("Server not responding properly")
            return
    except:
        print_error("Cannot connect to server at http://localhost:8080")
        return
    
    # Run stress tests on different endpoints
    endpoints_to_test = [
        ("/health", 1000, 50),
        ("/api/v1/branches", 500, 30),
        ("/api/v1/auth/login", 200, 10),  # Login endpoint with caution
    ]
    
    all_results = []
    
    for endpoint, num_requests, workers in endpoints_to_test:
        try:
            stats = stress_test(endpoint, num_requests, workers)
            all_results.append((endpoint, stats))
            time.sleep(2)  # Brief pause between tests
        except KeyboardInterrupt:
            print_error("\nStress test interrupted by user")
            break
        except Exception as e:
            print_error(f"Error testing {endpoint}: {e}")
    
    # Summary
    print_header("STRESS TEST SUMMARY")
    
    for endpoint, stats in all_results:
        print(f"\n{endpoint}:")
        print(f"  Success Rate: {stats['success_rate']:.2f}%")
        print(f"  Avg Response: {stats['avg_response_time']:.3f}s")
        print(f"  RPS: {stats['requests_per_second']:.2f}")
    
    # Overall assessment
    avg_success_rate = sum(s['success_rate'] for _, s in all_results) / len(all_results) if all_results else 0
    
    print("\n" + "="*70)
    if avg_success_rate >= 99:
        print_success(f"Overall Stress Test: PASSED (Avg Success: {avg_success_rate:.2f}%)")
    else:
        print_error(f"Overall Stress Test: NEEDS IMPROVEMENT (Avg Success: {avg_success_rate:.2f}%)")
    print("="*70 + "\n")

if __name__ == "__main__":
    main()

