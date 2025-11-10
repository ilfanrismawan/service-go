#!/bin/bash

# Script untuk diagnose masalah Docker container
# Bisa dijalankan dengan atau tanpa sudo

# Detect docker command
if command -v docker >/dev/null 2>&1 && docker ps >/dev/null 2>&1; then
    DOCKER_CMD="docker"
    DOCKER_COMPOSE_CMD="docker compose"
elif sudo docker ps >/dev/null 2>&1; then
    DOCKER_CMD="sudo docker"
    DOCKER_COMPOSE_CMD="sudo docker compose"
else
    echo "❌ Error: Docker tidak bisa diakses"
    exit 1
fi

echo "🔍 DIAGNOSIS DOCKER CONTAINER"
echo "=============================="
echo ""

# Check docker compose
echo "1. Checking Docker Compose status..."
$DOCKER_COMPOSE_CMD ps 2>&1
echo ""

# Check port mapping
echo "2. Checking port 8080 mapping..."
if $DOCKER_COMPOSE_CMD ps 2>&1 | grep -q "8080"; then
    echo "✅ Port 8080 ditemukan di docker compose"
    $DOCKER_COMPOSE_CMD ps 2>&1 | grep "8080"
else
    echo "❌ Port 8080 tidak ditemukan"
fi
echo ""

# Check app container logs
echo "3. Checking app container logs (last 30 lines)..."
$DOCKER_COMPOSE_CMD logs app --tail 30 2>&1
echo ""

# Check container status
echo "4. Container status detail..."
$DOCKER_COMPOSE_CMD ps app 2>&1
echo ""

# Check if app is listening inside container
echo "5. Checking if app is listening inside container..."
$DOCKER_COMPOSE_CMD exec -T app netstat -tlnp 2>/dev/null | grep 8080 || \
$DOCKER_COMPOSE_CMD exec -T app ss -tlnp 2>/dev/null | grep 8080 || \
$DOCKER_CMD exec iphone_service_app netstat -tlnp 2>/dev/null | grep 8080 || \
$DOCKER_CMD exec iphone_service_app ss -tlnp 2>/dev/null | grep 8080 || \
echo "⚠️  Tidak bisa cek dari dalam container (mungkin container tidak running atau app crash)"
echo ""

# Check host port
echo "6. Checking host port 8080..."
if ss -tlnp 2>/dev/null | grep -q ":8080"; then
    echo "✅ Port 8080 terbuka di host"
    ss -tlnp 2>/dev/null | grep ":8080"
else
    echo "❌ Port 8080 tidak terbuka di host"
    echo "   Ini berarti port mapping tidak bekerja atau container tidak running"
fi
echo ""

# Test connection
echo "7. Testing HTTP connection..."
if curl -s http://localhost:8080/health >/dev/null 2>&1; then
    echo "✅ HTTP connection OK"
    curl -s http://localhost:8080/health | head -5
else
    echo "❌ HTTP connection FAILED"
    echo "   Response:"
    curl -v http://localhost:8080/health 2>&1 | head -10
fi
echo ""

# Check container inspect
echo "8. Checking container port bindings..."
$DOCKER_CMD inspect iphone_service_app 2>/dev/null | grep -A 10 "Ports" || \
echo "⚠️  Tidak bisa inspect container"
echo ""

echo "=============================="
echo "💡 SOLUSI YANG MUNGKIN:"
echo "1. Restart container: DOCKER_SUDO=sudo make docker-restart"
echo "2. Cek logs: DOCKER_SUDO=sudo make docker-logs-app"
echo "3. Rebuild: DOCKER_SUDO=sudo make docker-up-build"
echo "=============================="

