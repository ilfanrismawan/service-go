# Troubleshooting: Container Running Tapi Tidak Bisa Diakses

## 🔴 Masalah
Container Docker berjalan (terlihat di `sudo docker ps`) tapi tidak bisa diakses di browser `http://localhost:8080`.

## 🔍 Diagnosis

Jalankan diagnosis script:
```bash
# Dengan sudo
sudo ./scripts/diagnose-docker.sh

# Atau via Makefile
DOCKER_SUDO=sudo make docker-diagnose
```

## ✅ Solusi Umum

### 1. Cek Status Container
```bash
DOCKER_SUDO=sudo make docker-ps
# atau
sudo docker ps
```

### 2. Cek Logs Container
```bash
DOCKER_SUDO=sudo make docker-logs-app
# atau
sudo docker logs iphone_service_app --tail 50
```

### 3. Cek Port Mapping
```bash
sudo docker port iphone_service_app
# Harus menunjukkan: 8080/tcp -> 0.0.0.0:8080
```

### 4. Restart Container
```bash
DOCKER_SUDO=sudo make docker-restart
# atau
sudo docker restart iphone_service_app
```

### 5. Rebuild Container
```bash
DOCKER_SUDO=sudo make docker-up-build
```

## 🐛 Masalah Umum dan Solusi

### Masalah 1: Port Mapping Tidak Bekerja
**Gejala:** Container running tapi port 8080 tidak terbuka di host

**Solusi:**
```bash
# Cek docker-compose.yml, pastikan:
ports:
  - "8080:8080"

# Restart dengan rebuild
DOCKER_SUDO=sudo make docker-up-build
```

### Masalah 2: Aplikasi Crash di Dalam Container
**Gejala:** Container start lalu langsung stop/restart

**Solusi:**
```bash
# Cek logs untuk error
DOCKER_SUDO=sudo make docker-logs-app

# Common issues:
# - Database connection failed
# - Missing environment variables
# - Application error
```

### Masalah 3: Aplikasi Tidak Listen di Port 8080
**Gejala:** Container running tapi app tidak listen di port yang benar

**Solusi:**
```bash
# Cek di dalam container
DOCKER_SUDO=sudo make docker-exec CMD="netstat -tlnp"

# Pastikan app listen di 0.0.0.0:8080 (bukan 127.0.0.1:8080)
```

### Masalah 4: Firewall Blocking
**Gejala:** Port tidak accessible dari luar

**Solusi:**
```bash
# Cek firewall
sudo ufw status

# Jika aktif, allow port 8080
sudo ufw allow 8080
```

### Masalah 5: Port Sudah Digunakan
**Gejala:** Error "port already in use"

**Solusi:**
```bash
# Cek apa yang menggunakan port 8080
sudo lsof -i :8080
# atau
sudo netstat -tlnp | grep 8080

# Stop service yang menggunakan port tersebut
# atau ubah port di docker-compose.yml
```

## 🔧 Langkah-langkah Debugging

### Step 1: Cek Container Status
```bash
sudo docker ps -a | grep app
```

**Expected:** Container status "Up" dan port mapping terlihat

### Step 2: Cek Logs
```bash
sudo docker logs iphone_service_app --tail 50
```

**Cari:**
- Error messages
- "Server listening on :8080"
- Database connection errors
- Application startup errors

### Step 3: Test dari Dalam Container
```bash
sudo docker exec iphone_service_app curl http://localhost:8080/health
```

**Expected:** Response JSON dengan status

### Step 4: Test dari Host
```bash
curl http://localhost:8080/health
```

**Expected:** Response JSON dengan status

### Step 5: Cek Port Binding
```bash
sudo docker port iphone_service_app
```

**Expected:** `8080/tcp -> 0.0.0.0:8080`

## 📋 Checklist

- [ ] Container status "Up" (tidak "Exited" atau "Restarting")
- [ ] Port mapping benar di docker-compose.yml
- [ ] Logs tidak menunjukkan error
- [ ] App listen di 0.0.0.0:8080 (bukan 127.0.0.1:8080)
- [ ] Database connection OK
- [ ] Port 8080 tidak digunakan aplikasi lain
- [ ] Firewall tidak block port 8080

## 🆘 Jika Masih Tidak Bisa

1. **Stop semua container:**
   ```bash
   DOCKER_SUDO=sudo make docker-down
   ```

2. **Rebuild dari scratch:**
   ```bash
   DOCKER_SUDO=sudo make docker-up-build
   ```

3. **Cek logs real-time:**
   ```bash
   DOCKER_SUDO=sudo make docker-logs-app
   ```

4. **Masuk ke container untuk debug:**
   ```bash
   DOCKER_SUDO=sudo make docker-exec CMD="sh"
   # Di dalam container:
   # ps aux
   # netstat -tlnp
   # curl http://localhost:8080/health
   ```

## 📞 Informasi untuk Debug

Saat melaporkan masalah, sertakan:
- Output dari `sudo docker ps`
- Output dari `sudo docker logs iphone_service_app --tail 50`
- Output dari `sudo docker port iphone_service_app`
- Output dari `curl -v http://localhost:8080/health`

