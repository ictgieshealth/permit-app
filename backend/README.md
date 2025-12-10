# Permit Management System - Backend

Sistem manajemen perizinan perusahaan yang dibangun dengan Go (Golang), Gin Framework, dan GORM.

## 📋 Deskripsi

Aplikasi backend untuk mengelola perizinan di perusahaan dengan fitur lengkap CRUD (Create, Read, Update, Delete) untuk:
- **Division**: Manajemen divisi/departemen perusahaan
- **Permit Type**: Manajemen jenis-jenis perizinan
- **Permit**: Manajemen data perizinan aktual

## 🚀 Teknologi

- **Go 1.21+**: Programming language
- **Gin**: Web framework
- **GORM**: ORM library
- **PostgreSQL**: Database
- **Air**: Live reload untuk development

## 📁 Struktur Proyek

```
backend/
├── controller/          # HTTP request handlers
│   ├── divisionController/
│   ├── permitTypeController/
│   └── permitController/
├── database/           # Database configuration & migrations
│   ├── db.go
│   └── migration.sql
├── helper/             # Utility functions
│   ├── apiresponse/   # API response helpers
│   ├── apiRequest/    # API request helpers
│   └── validator.go   # Request validation
├── middleware/         # HTTP middlewares
│   ├── auth.go
│   └── authorization.go
├── model/             # Data models
│   ├── division.go
│   ├── permitType.go
│   ├── permit.go
│   └── user.go
├── repo/              # Repository layer (database operations)
│   ├── divisionRepository/
│   ├── permitTypeRepository/
│   └── permitRepository/
├── routes/            # Route definitions
│   └── route.go
├── service/           # Business logic layer
│   ├── divisionService/
│   ├── permitTypeService/
│   └── permitService/
├── main.go           # Application entry point
└── go.mod           # Go dependencies
```

## 🔧 Instalasi

### Prerequisites
- Go 1.21 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi
- Git

### Langkah Instalasi

1. **Clone repository**
```bash
git clone <repository-url>
cd backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup Database**
```bash
# Buat database PostgreSQL
createdb permit_management

# Jalankan migration
psql -U postgres -d permit_management -f database/migration.sql
```

4. **Setup environment variables**
Buat file `.env` dengan konfigurasi berikut:
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=permit_management
DB_SSL_MODE=disable
JWT_SECRET=your_jwt_secret_key
CORS_ALLOW_ORIGINS=http://localhost:3000
```

5. **Jalankan aplikasi**
```bash
# Development mode dengan live reload
air

# Atau production mode
go run main.go
```

## 📚 API Documentation

Dokumentasi lengkap API tersedia di [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

### Quick Reference

#### Division Endpoints
- `POST /divisions` - Create division
- `GET /divisions` - Get all divisions
- `GET /divisions/:id` - Get division by ID
- `PUT /divisions/:id` - Update division
- `DELETE /divisions/:id` - Delete division

#### Permit Type Endpoints
- `POST /permit-types` - Create permit type
- `GET /permit-types` - Get all permit types
- `GET /permit-types/:id` - Get permit type by ID
- `PUT /permit-types/:id` - Update permit type
- `DELETE /permit-types/:id` - Delete permit type

#### Permit Endpoints
- `POST /permits` - Create permit
- `GET /permits` - Get all permits
- `GET /permits/:id` - Get permit by ID
- `PUT /permits/:id` - Update permit
- `DELETE /permits/:id` - Delete permit

## 🗄️ Database Schema

### Tables

1. **divisions**
   - `id`: Primary key (BIGSERIAL)
   - `code`: Unique code (VARCHAR)
   - `name`: Division name (VARCHAR)
   - `created_at`: Timestamp
   - `updated_at`: Timestamp

2. **permit_types**
   - `id`: Primary key (BIGSERIAL)
   - `division_id`: Foreign key to divisions
   - `category`: Permit category (VARCHAR)
   - `name`: Permit type name (VARCHAR)
   - `risk_point`: Risk level (VARCHAR)
   - `default_application_type`: Default application type (VARCHAR)
   - `default_validity_period`: Default validity period (VARCHAR)
   - `notes`: Additional notes (TEXT)
   - `created_at`: Timestamp
   - `updated_at`: Timestamp

3. **permits**
   - `id`: Primary key (BIGSERIAL)
   - `permit_type_id`: Foreign key to permit_types
   - `application_type`: Type of application (VARCHAR)
   - `permit_no`: Unique permit number (VARCHAR)
   - `effective_date`: Start date (TIMESTAMP)
   - `expiry_date`: End date (TIMESTAMP)
   - `effective_term`: Duration (VARCHAR)
   - `responsible_person`: Person in charge (VARCHAR)
   - `doc_name`: Document name (VARCHAR)
   - `doc_number`: Document number (VARCHAR)
   - `status`: Permit status (VARCHAR)
   - `created_at`: Timestamp
   - `updated_at`: Timestamp

## 🔒 Authentication & Authorization

Module User dipertahankan sebagai contoh implementasi autentikasi dan autorisasi. Untuk menggunakannya, implementasikan middleware auth di routes sesuai kebutuhan.

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## 📝 Development Guidelines

### Code Structure
- **Controller**: Handle HTTP requests/responses
- **Service**: Business logic implementation
- **Repository**: Database operations
- **Model**: Data structures

### Naming Conventions
- Use camelCase for variable and function names
- Use PascalCase for struct and interface names
- Use descriptive names

### Error Handling
- Always return appropriate HTTP status codes
- Provide clear error messages
- Log errors for debugging

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.