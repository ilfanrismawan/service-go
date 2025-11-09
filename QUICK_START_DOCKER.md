# Quick Start: Menjalankan Go App dengan Docker

## 🔴 Error: Permission Denied

Jika Anda mendapatkan error:
```
permission denied while trying to connect to the Docker daemon socket
```

## ✅ Solusi Cepat

### Opsi 1: Setup Docker Permissions (RECOMMENDED - Hanya sekali)

```bash
# Jalankan script setup
make docker-setup

# Logout dan login lagi untuk apply perubahan group
# Setelah login lagi, verifikasi:
groups | grep docker

# Jika sudah ada "docker" di output, lanjutkan:
make docker-up-build
```

### Opsi 2: Gunakan Sudo (Alternatif)

```bash
# Untuk satu command
DOCKER_SUDO=sudo make docker-up-build

# Atau export untuk session
export DOCKER_SUDO=sudo
make docker-up-build
make docker-logs
make docker-down
```

## 📋 Langkah-langkah Lengkap

1. **Cek setup Docker:**
   ```bash
   make docker-help
   ```

2. **Setup permissions (jika belum):**
   ```bash
   make docker-setup
   # Logout dan login lagi
   ```

3. **Build dan start services:**
   ```bash
   make docker-up-build
   ```

4. **Lihat logs:**
   ```bash
   make docker-logs
   # atau hanya app
   make docker-logs-app
   ```

5. **Akses aplikasi:**
   - API: http://localhost:8080
   - MinIO Console: http://localhost:9001

6. **Stop services:**
   ```bash
   make docker-down
   ```

## 🆘 Troubleshooting

### User sudah di docker group tapi masih error?
- Pastikan sudah logout dan login lagi setelah menjalankan `make docker-setup`
- Verifikasi dengan: `groups | grep docker`
- Jika masih error, coba: `newgrp docker` (untuk apply group tanpa logout)

### Ingin reset semua?
```bash
make docker-down-volumes  # Stop dan hapus volumes
```

### Lihat semua perintah Docker:
```bash
make help | grep docker
```

