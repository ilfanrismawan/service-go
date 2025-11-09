# Fix: Docker Build Error - Go Version Mismatch

## 🔴 Error yang Terjadi

```
go: golang.org/x/sync@v0.12.0: module golang.org/x/sync@v0.12.0 requires go >= 1.23.0 (running go 1.21.13; GOTOOLCHAIN=local)
```

## ✅ Solusi

Dockerfile diupdate dari Go 1.21 ke Go 1.23 untuk kompatibilitas dengan dependencies.

### Perubahan:
```diff
- FROM golang:1.21-alpine AS builder
+ FROM golang:1.23-alpine AS builder
```

## 🚀 Build Ulang

Setelah perbaikan, build ulang dengan:

```bash
# Jika menggunakan sudo
DOCKER_SUDO=sudo make docker-up-build

# Atau jika user sudah di docker group
make docker-up-build
```

## 📝 Catatan

- Go 1.23 backward compatible dengan Go 1.18 (yang dideklarasikan di go.mod)
- Docker build akan menggunakan Go 1.23 untuk compile, yang kompatibel dengan semua dependencies
- Aplikasi yang dihasilkan tetap dapat berjalan di environment Go 1.18+

## 🔍 Verifikasi

Setelah build berhasil, cek dengan:
```bash
make docker-ps
make docker-logs-app
```

