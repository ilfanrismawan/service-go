#!/bin/bash
set -e

echo "🚀 Memulai refactor struktur ke domain-based (versi dengan auto import fix)..."
sleep 1

# --- 1. Buat struktur domain utama
echo "📁 Membuat folder domain utama..."
mkdir -p internal/{users,orders,payments,branches}/
for domain in users orders payments branches; do
  mkdir -p internal/$domain/{handler,service,repository,dto}
done

# --- 2. Buat shared utilities
echo "📁 Membuat folder shared..."
mkdir -p internal/shared/{config,database,middleware,notification,monitoring,utils,errors}

# --- 3. Pindahkan folder global yang sudah ada ke shared
echo "📦 Memindahkan folder global ke shared..."
for f in config database middleware notification monitoring utils; do
  if [ -d internal/$f ]; then
    mv internal/$f internal/shared/$f
  fi
done

# --- 4. Pindahkan layer-based ke domain (copy sementara)
echo "📦 Menyalin service, repository, delivery ke tiap domain..."
for folder in internal/service internal/repository internal/delivery; do
  if [ -d "$folder" ]; then
    cp -r $folder internal/orders/ 2>/dev/null || true
    cp -r $folder internal/payments/ 2>/dev/null || true
    cp -r $folder internal/users/ 2>/dev/null || true
    cp -r $folder internal/branches/ 2>/dev/null || true
  fi
done

# --- 5. Auth & Payment
if [ -d internal/auth ]; then
  echo "🔐 Memindahkan auth ke domain users..."
  mv internal/auth internal/users/auth
fi

if [ -d internal/payment ]; then
  echo "💳 Memindahkan payment ke domain payments..."
  mv internal/payment internal/payments/legacy_payment
fi

# --- 6. Core (model/entities)
if [ -d internal/core ]; then
  echo "📦 Menyalin core ke shared/model..."
  mkdir -p internal/shared/model
  cp -r internal/core/* internal/shared/model/
fi

# --- 7. Auto fix import paths
echo "🔍 Memperbarui import path di seluruh file Go..."
# Backup dulu
find . -type f -name "*.go" -exec cp {} {}.bak \;

# Replace layer-based import ke domain-based
find . -type f -name "*.go" -exec sed -i \
  -e 's|internal/service|internal/orders/service|g' \
  -e 's|internal/repository|internal/orders/repository|g' \
  -e 's|internal/delivery|internal/orders/handler|g' \
  -e 's|internal/payment|internal/payments|g' \
  -e 's|internal/auth|internal/users/auth|g' \
  {} +

echo "✅ Import path berhasil diperbarui (backup: *.go.bak)"

# --- 8. Log hasil
echo
echo "✅ Struktur baru berhasil dibuat!"
tree -L 3 internal | tee refactor_result.log

echo
echo "📋 Cek log di file refactor_result.log"
echo "⚠️ Backup file tersimpan dengan ekstensi .go.bak"
echo "👉 Setelah verifikasi, hapus file .bak dengan: find . -name '*.bak' -delete"
echo "👉 Jalankan: go mod tidy && go test ./..."
