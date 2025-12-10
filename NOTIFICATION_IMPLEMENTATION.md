# Notification System Implementation

## Overview
Sistem notifikasi lengkap untuk mengingatkan user tentang permit yang akan expired. Notifikasi dikirim melalui:
1. **In-App Notification** - Notifikasi di dalam aplikasi web
2. **Email Notification** - Notifikasi dikirim ke email

## Features

### Backend Features
- ✅ Email notification dengan SMTP (Gmail)
- ✅ In-app notification system
- ✅ Scheduler otomatis untuk check permit expiry
- ✅ Notifikasi 30 hari sebelum expiry (reminder)
- ✅ Notifikasi 7 hari sebelum expiry (warning)
- ✅ Notifikasi saat permit expired
- ✅ Real-time notification counting
- ✅ Mark as read/unread functionality
- ✅ Notification history

### Frontend Features
- ✅ Real-time notification bell with badge counter
- ✅ Notification dropdown preview
- ✅ Full notification page
- ✅ Mark as read functionality
- ✅ Delete notification
- ✅ Auto-refresh notification count every 30 seconds

## Configuration

### Backend Configuration (.env)
```env
# Mail Configuration
MAIL_HOST=smtp.gmail.com
MAIL_PORT=465
MAIL_USERNAME=rstd@rs-triadipa.com
MAIL_PASSWORD=amhrccoknmdssiww
MAIL_FROM_ADDRESS=rstd@rs-triadipa.com
MAIL_FROM_NAME=Permit Management System
```

### Email Setup (Gmail)
1. Enable 2-Factor Authentication di Google Account
2. Generate App Password:
   - Go to: https://myaccount.google.com/apppasswords
   - Generate password untuk "Mail"
   - Gunakan password tersebut sebagai `MAIL_PASSWORD`

## Database Migration

Jalankan SQL migration untuk membuat tabel notifications:
```bash
cd backend/database
psql -U postgres -d permit_app -f migration_notifications.sql
```

Atau migration akan otomatis dijalankan saat aplikasi start (via GORM AutoMigrate).

## API Endpoints

### Get User Notifications
```
GET /notifications?page=1&limit=10
Authorization: Bearer {token}
```

### Get Unread Notifications
```
GET /notifications/unread
Authorization: Bearer {token}
```

### Get Unread Count
```
GET /notifications/unread/count
Authorization: Bearer {token}
```

### Mark As Read
```
POST /notifications/read
Authorization: Bearer {token}
Content-Type: application/json

{
  "notification_ids": [1, 2, 3]
}
```

### Mark All As Read
```
POST /notifications/read/all
Authorization: Bearer {token}
```

### Delete Notification
```
DELETE /notifications/{id}
Authorization: Bearer {token}
```

## Notification Types

1. **expiry_reminder** - 30 hari sebelum expired
   - Badge: Blue
   - Dikirim ke: Permit Responsible, Document Responsible, Admin, Manager (domain)

2. **expiry_warning** - 7 hari sebelum expired
   - Badge: Orange
   - Dikirim ke: Permit Responsible, Document Responsible, Admin, Manager (domain)

3. **expired** - Permit sudah expired
   - Badge: Red
   - Dikirim ke: Permit Responsible, Document Responsible, Admin, Manager (domain)

## Notification Recipients

Notifikasi akan dikirim ke:
1. **Permit Responsible** - User yang bertanggung jawab atas permit
2. **Document Responsible** - User yang bertanggung jawab atas dokumen
3. **Admin** - Semua user dengan role "admin"
4. **Manager** - User dengan role "manager" pada domain yang sama dengan permit

## Scheduler

Scheduler berjalan otomatis setiap hari jam **08:00 pagi** untuk:
- Check semua permit yang akan expired
- Kirim notifikasi in-app dan email
- Log aktivitas ke console

### Testing Scheduler
Untuk testing, Anda bisa mengubah scheduler interval di `main.go`:
```go
// Ganti dari:
notificationScheduler.Start()

// Menjadi (untuk testing setiap 5 menit):
notificationScheduler.StartWithInterval(5 * time.Minute)
```

## Frontend Usage

### 1. Import NotificationProvider
Sudah diimplementasikan di `app/(admin)/layout.tsx`:
```tsx
import { NotificationProvider } from "@/context/NotificationContext";

<NotificationProvider>
  {/* Your app content */}
</NotificationProvider>
```

### 2. Use NotificationBell Component
Sudah diimplementasikan di `AppHeader.tsx`:
```tsx
import NotificationBell from "@/components/header/NotificationBell";

<NotificationBell />
```

### 3. Use Notification Hook
Di component manapun:
```tsx
import { useNotifications } from "@/context/NotificationContext";

function MyComponent() {
  const { notifications, unreadCount, markAsRead } = useNotifications();
  
  return (
    <div>
      <p>Unread: {unreadCount}</p>
      {/* ... */}
    </div>
  );
}
```

## Files Created/Modified

### Backend
**Created:**
- `backend/helper/email.go` - Email service
- `backend/model/notification.go` - Notification model
- `backend/repo/notificationRepository/notificationRepository.go` - Repository
- `backend/service/notificationService/notificationService.go` - Service
- `backend/controller/notificationController/notificationController.go` - Controller
- `backend/scheduler/scheduler.go` - Notification scheduler
- `backend/database/migration_notifications.sql` - SQL migration

**Modified:**
- `backend/main.go` - Added scheduler initialization
- `backend/routes/route.go` - Added notification endpoints
- `backend/repo/permitRepository/permitRepository.go` - Added FindExpiringPermits
- `backend/repo/userRepository/userRepository.go` - Added user search by role
- `backend/.env` - Added mail configuration

### Frontend
**Created:**
- `frontend/src/services/notificationService.ts` - API service
- `frontend/src/context/NotificationContext.tsx` - Context provider
- `frontend/src/components/Header/NotificationBell.tsx` - Bell component
- `frontend/src/app/(admin)/notifications/page.tsx` - Full notifications page

**Modified:**
- `frontend/src/app/(admin)/layout.tsx` - Added NotificationProvider
- `frontend/src/layout/AppHeader.tsx` - Added NotificationBell
- `frontend/package.json` - Added date-fns dependency

## Testing

### 1. Test Email Configuration
```bash
# Di backend directory
go run tools/test_email.go
```

### 2. Test Notification API
```bash
# Get notifications
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8000/notifications

# Get unread count
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8000/notifications/unread/count
```

### 3. Test Frontend
1. Login ke aplikasi
2. Check notification bell di header
3. Click bell untuk melihat notifications
4. Click "Lihat Semua Notifikasi" untuk halaman lengkap

## Troubleshooting

### Email tidak terkirim
1. Check SMTP credentials di .env
2. Pastikan App Password sudah benar
3. Check firewall untuk port 465
4. Check log backend untuk error details

### Notification tidak muncul
1. Check apakah scheduler berjalan (lihat log backend)
2. Check apakah ada permit yang akan expired
3. Check database table notifications
4. Check browser console untuk error

### Badge count tidak update
1. Refresh halaman
2. Check network tab untuk API calls
3. Pastikan token masih valid

## Best Practices

1. **Production:** Gunakan dedicated email server, bukan Gmail
2. **Security:** Simpan MAIL_PASSWORD di environment variable, jangan commit ke git
3. **Performance:** Limit notification query dengan pagination
4. **Testing:** Set scheduler interval pendek untuk testing
5. **Monitoring:** Monitor email delivery rate dan notification performance

## Future Enhancements

- [ ] Push notifications (Firebase)
- [ ] SMS notifications
- [ ] Notification preferences per user
- [ ] Notification templates customization
- [ ] Notification delivery reports
- [ ] Retry mechanism untuk failed emails
- [ ] Notification grouping
- [ ] Sound alerts untuk notifications

## Support

Untuk pertanyaan atau issue, hubungi:
- Development Team
- Email: rstd@rs-triadipa.com
