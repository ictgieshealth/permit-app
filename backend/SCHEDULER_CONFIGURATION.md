# Notification Scheduler Configuration

## Environment Variables

Konfigurasi scheduler notification dapat diatur melalui file `.env`:

### Mode Scheduler

```env
SCHEDULER_MODE=production
```

**Pilihan:**
- `production` - Scheduler berjalan setiap hari pada waktu yang ditentukan
- `testing` - Scheduler berjalan dengan interval waktu tertentu (untuk development/testing)

### Production Mode Configuration

Untuk mode `production`, atur jam dan menit kapan scheduler harus berjalan:

```env
SCHEDULER_HOUR=8
SCHEDULER_MINUTE=0
```

**SCHEDULER_HOUR:**
- Jam kapan scheduler harus berjalan (0-23)
- Default: `8` (jam 8 pagi)
- Contoh: `14` untuk jam 2 siang

**SCHEDULER_MINUTE:**
- Menit kapan scheduler harus berjalan (0-59)
- Default: `0` (tepat jam)
- Contoh: `30` untuk 30 menit setelah jam

### Testing Mode Configuration

Untuk mode `testing`, atur interval berapa menit scheduler harus berjalan:

```env
SCHEDULER_INTERVAL_MINUTES=5
```

**SCHEDULER_INTERVAL_MINUTES:**
- Interval dalam menit untuk menjalankan scheduler
- Default: `5` (setiap 5 menit)
- Berguna untuk testing dan development

## Contoh Konfigurasi

### Production - Jam 8 Pagi

```env
SCHEDULER_MODE=production
SCHEDULER_HOUR=8
SCHEDULER_MINUTE=0
```

Scheduler akan berjalan setiap hari jam 08:00

### Production - Jam 2 Siang 30 Menit

```env
SCHEDULER_MODE=production
SCHEDULER_HOUR=14
SCHEDULER_MINUTE=30
```

Scheduler akan berjalan setiap hari jam 14:30

### Testing - Setiap 5 Menit

```env
SCHEDULER_MODE=testing
SCHEDULER_INTERVAL_MINUTES=5
```

Scheduler akan berjalan setiap 5 menit (berguna untuk development)

### Testing - Setiap 1 Jam

```env
SCHEDULER_MODE=testing
SCHEDULER_INTERVAL_MINUTES=60
```

Scheduler akan berjalan setiap 1 jam

## Cara Mengubah Konfigurasi

1. Edit file `.env` di folder `backend/`
2. Ubah nilai variable sesuai kebutuhan
3. Restart aplikasi backend
4. Scheduler akan otomatis menggunakan konfigurasi baru

## Log Scheduler

Ketika aplikasi berjalan, Anda akan melihat log seperti:

**Production Mode:**
```
Notification Scheduler: Started successfully
Notification Scheduler: Next run scheduled at 2025-12-10 08:00:00 (in 10h15m30s)
```

**Testing Mode:**
```
Notification Scheduler: Started with 5m0s interval
Notification Scheduler: Initial check for expiring permits
```

## Notes

- Pastikan nilai `SCHEDULER_HOUR` antara 0-23
- Pastikan nilai `SCHEDULER_MINUTE` antara 0-59
- Untuk production, sebaiknya gunakan `SCHEDULER_MODE=production`
- Untuk testing/development, gunakan `SCHEDULER_MODE=testing`
- Scheduler akan otomatis menjalankan initial check saat aplikasi pertama kali start
