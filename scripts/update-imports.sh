#!/bin/bash
# Script to update import paths after domain-based refactoring

echo "🔧 Updating import paths for domain-based architecture..."

# Update users domain files
echo "📝 Updating users domain..."
find internal/users -name "*.go" -type f -exec sed -i 's|service/internal/core|service/internal/users/dto|g' {} \;
find internal/users -name "*.go" -type f -exec sed -i 's|service/internal/repository|service/internal/users/repository|g' {} \;
find internal/users -name "*.go" -type f -exec sed -i 's|service/internal/utils|service/internal/shared/utils|g' {} \;

# Update orders domain files
echo "📝 Updating orders domain..."
find internal/orders -name "*.go" -type f -exec sed -i 's|service/internal/core|service/internal/orders/dto|g' {} \;
find internal/orders -name "*.go" -type f -exec sed -i 's|service/internal/service|service/internal/orders/service|g' {} \;
find internal/orders -name "*.go" -type f -exec sed -i 's|service/internal/repository|service/internal/orders/repository|g' {} \;
find internal/orders -name "*.go" -type f -exec sed -i 's|service/internal/utils|service/internal/shared/utils|g' {} \;

# Update payments domain files
echo "📝 Updating payments domain..."
find internal/payments -name "*.go" -type f -exec sed -i 's|service/internal/core|service/internal/payments/dto|g' {} \;
find internal/payments -name "*.go" -type f -exec sed -i 's|service/internal/service|service/internal/payments/service|g' {} \;
find internal/payments -name "*.go" -type f -exec sed -i 's|service/internal/repository|service/internal/payments/repository|g' {} \;
find internal/payments -name "*.go" -type f -exec sed -i 's|service/internal/payment|service/internal/payments/legacy_payment|g' {} \;
find internal/payments -name "*.go" -type f -exec sed -i 's|service/internal/utils|service/internal/shared/utils|g' {} \;

# Update branches domain files
echo "📝 Updating branches domain..."
find internal/branches -name "*.go" -type f -exec sed -i 's|service/internal/core|service/internal/branches/dto|g' {} \;
find internal/branches -name "*.go" -type f -exec sed -i 's|service/internal/service|service/internal/branches/service|g' {} \;
find internal/branches -name "*.go" -type f -exec sed -i 's|service/internal/repository|service/internal/branches/repository|g' {} \;
find internal/branches -name "*.go" -type f -exec sed -i 's|service/internal/utils|service/internal/shared/utils|g' {} \;

# Update shared resources
echo "📝 Updating shared resources..."
find internal/shared -name "*.go" -type f -exec sed -i 's|service/internal/config|service/internal/shared/config|g' {} \;
find internal/shared -name "*.go" -type f -exec sed -i 's|service/internal/database|service/internal/shared/database|g' {} \;
find internal/shared -name "*.go" -type f -exec sed -i 's|service/internal/middleware|service/internal/shared/middleware|g' {} \;
find internal/shared -name "*.go" -type f -exec sed -i 's|service/internal/utils|service/internal/shared/utils|g' {} \;

echo "✅ Import paths updated!"

