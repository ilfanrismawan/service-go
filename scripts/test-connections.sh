#!/bin/sh

echo "Testing Redis connection..."
nc -zv redis 6379 || echo "Redis connection failed"

echo "\nTesting PostgreSQL connection..."
nc -zv postgres 5432 || echo "PostgreSQL connection failed"

echo "\nTesting MinIO connection..."
nc -zv minio 9000 || echo "MinIO connection failed"