#!/bin/bash

# Script to setup Docker permissions for current user

set -e

echo "🔧 Setting up Docker permissions..."
echo ""

# Check if user is already in docker group
if groups | grep -q docker; then
    echo "✅ User $(whoami) is already in docker group"
    echo "   You may need to logout and login again for changes to take effect"
    exit 0
fi

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    echo "❌ Please don't run this script as root"
    echo "   Run it as a regular user and it will use sudo when needed"
    exit 1
fi

# Add user to docker group
echo "📝 Adding user $(whoami) to docker group..."
sudo usermod -aG docker $(whoami)

echo ""
echo "✅ User $(whoami) has been added to docker group"
echo ""
echo "⚠️  IMPORTANT: You need to logout and login again for the changes to take effect"
echo ""
echo "After logging in again, verify with:"
echo "  groups | grep docker"
echo ""
echo "Then you can use docker without sudo:"
echo "  make docker-up-build"

