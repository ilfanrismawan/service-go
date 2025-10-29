#!/bin/bash

# Generate Swagger documentation
echo "Generating Swagger documentation..."

# Install swag if not installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate docs
swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal

echo "Swagger documentation generated successfully!"
echo "API Documentation available at: http://localhost:8080/swagger/index.html"
echo "API Docs available at: http://localhost:8080/docs"
