# Panduan Docker untuk iPhone Service API

## Setup Awal

### 1. Setup Docker Permission (Pilih salah satu)

**Opsi A: Tambahkan user ke docker group (Recommended)**
```bash
sudo usermod -aG docker $USER
# Logout dan login lagi untuk apply perubahan
```

**Opsi B: Gunakan sudo untuk setiap command**
```bash
# Gunakan DOCKER_SUDO=sudo sebelum command
DOCKER_SUDO=sudo make docker-up-build
DOCKER_SUDO=sudo make docker-logs
DOCKER_SUDO=sudo make docker-down

# Atau export sebagai environment variable untuk session
export DOCKER_SUDO=sudo
make docker-up-build
make docker-logs
```

### 2. Cek Setup Docker
```bash
make docker-help
```

## Menjalankan Aplikasi

### Build dan Start Semua Services
```bash
make docker-up-build
```

Ini akan:
- Build Docker image untuk aplikasi Go
- Start PostgreSQL database
- Start Redis cache
- Start MinIO (S3-compatible storage)
- Start aplikasi Go API

### Lihat Logs
```bash
# Semua services
make docker-logs

# Hanya aplikasi
make docker-logs-app
```

### Stop Services
```bash
make docker-down
```

### Restart Services
```bash
make docker-restart
```

## Akses Services

Setelah services running:
- **API**: http://localhost:8080
- **MinIO Console**: http://localhost:9001
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## Perintah Lainnya

### Build Image Saja
```bash
make docker-build
```

### Rebuild Image (No Cache)
```bash
make docker-rebuild
```

### Lihat Container Status
```bash
make docker-ps
```

### Masuk ke Container
```bash
make docker-exec CMD="sh"
```

### Stop dan Hapus Volumes
```bash
make docker-down-volumes
```

## Troubleshooting

### Permission Denied
Jika mendapat error permission, pastikan user sudah di docker group:
```bash
groups | grep docker
```

Jika tidak ada, tambahkan:
```bash
sudo usermod -aG docker $USER
# Logout dan login
```

### Port Already in Use
Jika port sudah digunakan, stop service yang menggunakan port tersebut atau ubah port di `docker-compose.yml`

### Build Error
Jika build error, coba rebuild tanpa cache:
```bash
make docker-rebuild
```

## Catatan Penting

1. **Docker menggunakan Go compiler standar** (bukan gccgo), jadi semua fitur Go termasuk generics akan bekerja dengan baik
2. **Database akan persist** di volume `postgres_data` bahkan setelah container di-stop
3. **Untuk reset database**, gunakan `make docker-down-volumes` (hati-hati, akan hapus semua data)

