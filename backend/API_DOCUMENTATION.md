# Permit Management System API Documentation

Sistem manajemen perizinan perusahaan dengan fitur CRUD untuk Division, Permit Type, dan Permit.

## Base URL
```
http://localhost:{PORT}
```

## Endpoints

### Division API

#### 1. Create Division
**POST** `/divisions`

**Request Body:**
```json
{
  "code": "IT",
  "name": "Information Technology"
}
```

**Response (201):**
```json
{
  "status": "success",
  "message": "Division created successfully",
  "data": {
    "id": 1,
    "code": "IT",
    "name": "Information Technology",
    "created_at": "2025-11-25T10:00:00Z",
    "updated_at": "2025-11-25T10:00:00Z"
  }
}
```

#### 2. Get All Divisions
**GET** `/divisions`

**Query Parameters:**
- `code` (optional): Filter by code
- `name` (optional): Filter by name
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response (200):**
```json
{
  "status": "success",
  "message": "Divisions retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "IT",
      "name": "Information Technology",
      "created_at": "2025-11-25T10:00:00Z",
      "updated_at": "2025-11-25T10:00:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

#### 3. Get Division by ID
**GET** `/divisions/:id`

**Response (200):**
```json
{
  "status": "success",
  "message": "Division retrieved successfully",
  "data": {
    "id": 1,
    "code": "IT",
    "name": "Information Technology",
    "created_at": "2025-11-25T10:00:00Z",
    "updated_at": "2025-11-25T10:00:00Z"
  }
}
```

#### 4. Update Division
**PUT** `/divisions/:id`

**Request Body:**
```json
{
  "code": "IT-DEPT",
  "name": "IT Department"
}
```

**Response (200):**
```json
{
  "status": "success",
  "message": "Division updated successfully",
  "data": {
    "id": 1,
    "code": "IT-DEPT",
    "name": "IT Department",
    "created_at": "2025-11-25T10:00:00Z",
    "updated_at": "2025-11-25T11:00:00Z"
  }
}
```

#### 5. Delete Division
**DELETE** `/divisions/:id`

**Response (200):**
```json
{
  "status": "success",
  "message": "Division deleted successfully"
}
```

---

### Permit Type API

#### 1. Create Permit Type
**POST** `/permit-types`

**Request Body:**
```json
{
  "division_id": 1,
  "category": "Operational",
  "name": "Building Permit",
  "risk_point": "High",
  "default_application_type": "New Application",
  "default_validity_period": "1 year",
  "notes": "Required for construction activities"
}
```

**Response (201):**
```json
{
  "status": "success",
  "message": "Permit type created successfully",
  "data": {
    "id": 1,
    "division_id": 1,
    "category": "Operational",
    "name": "Building Permit",
    "risk_point": "High",
    "default_application_type": "New Application",
    "default_validity_period": "1 year",
    "notes": "Required for construction activities",
    "created_at": "2025-11-25T10:00:00Z",
    "updated_at": "2025-11-25T10:00:00Z",
    "division": {
      "id": 1,
      "code": "IT",
      "name": "Information Technology",
      "created_at": "2025-11-25T10:00:00Z",
      "updated_at": "2025-11-25T10:00:00Z"
    }
  }
}
```

#### 2. Get All Permit Types
**GET** `/permit-types`

**Query Parameters:**
- `division_id` (optional): Filter by division ID
- `category` (optional): Filter by category
- `name` (optional): Filter by name
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response (200):**
```json
{
  "status": "success",
  "message": "Permit types retrieved successfully",
  "data": [
    {
      "id": 1,
      "division_id": 1,
      "category": "Operational",
      "name": "Building Permit",
      "risk_point": "High",
      "default_application_type": "New Application",
      "default_validity_period": "1 year",
      "notes": "Required for construction activities",
      "created_at": "2025-11-25T10:00:00Z",
      "updated_at": "2025-11-25T10:00:00Z",
      "division": {
        "id": 1,
        "code": "IT",
        "name": "Information Technology",
        "created_at": "2025-11-25T10:00:00Z",
        "updated_at": "2025-11-25T10:00:00Z"
      }
    }
  ],
  "pagination": {
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

#### 3. Get Permit Type by ID
**GET** `/permit-types/:id`

#### 4. Update Permit Type
**PUT** `/permit-types/:id`

#### 5. Delete Permit Type
**DELETE** `/permit-types/:id`

---

### Permit API

#### 1. Create Permit
**POST** `/permits`

**Request Body:**
```json
{
  "permit_type_id": 1,
  "application_type": "New Application",
  "permit_no": "PRM-2025-001",
  "effective_date": "2025-01-01T00:00:00Z",
  "expiry_date": "2026-01-01T00:00:00Z",
  "effective_term": "1 year",
  "responsible_person": "John Doe",
  "doc_name": "Building Construction Permit",
  "doc_number": "DOC-2025-001",
  "status": "active"
}
```

**Response (201):**
```json
{
  "status": "success",
  "message": "Permit created successfully",
  "data": {
    "id": 1,
    "permit_type_id": 1,
    "application_type": "New Application",
    "permit_no": "PRM-2025-001",
    "effective_date": "2025-01-01T00:00:00Z",
    "expiry_date": "2026-01-01T00:00:00Z",
    "effective_term": "1 year",
    "responsible_person": "John Doe",
    "doc_name": "Building Construction Permit",
    "doc_number": "DOC-2025-001",
    "status": "active",
    "created_at": "2025-11-25T10:00:00Z",
    "updated_at": "2025-11-25T10:00:00Z",
    "permit_type": {
      "id": 1,
      "division_id": 1,
      "category": "Operational",
      "name": "Building Permit",
      "risk_point": "High",
      "default_application_type": "New Application",
      "default_validity_period": "1 year",
      "notes": "Required for construction activities",
      "created_at": "2025-11-25T10:00:00Z",
      "updated_at": "2025-11-25T10:00:00Z",
      "division": {
        "id": 1,
        "code": "IT",
        "name": "Information Technology",
        "created_at": "2025-11-25T10:00:00Z",
        "updated_at": "2025-11-25T10:00:00Z"
      }
    }
  }
}
```

#### 2. Get All Permits
**GET** `/permits`

**Query Parameters:**
- `permit_type_id` (optional): Filter by permit type ID
- `application_type` (optional): Filter by application type
- `permit_no` (optional): Filter by permit number
- `responsible_person` (optional): Filter by responsible person
- `status` (optional): Filter by status (active, expired, revoked)
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response (200):**
```json
{
  "status": "success",
  "message": "Permits retrieved successfully",
  "data": [
    {
      "id": 1,
      "permit_type_id": 1,
      "application_type": "New Application",
      "permit_no": "PRM-2025-001",
      "effective_date": "2025-01-01T00:00:00Z",
      "expiry_date": "2026-01-01T00:00:00Z",
      "effective_term": "1 year",
      "responsible_person": "John Doe",
      "doc_name": "Building Construction Permit",
      "doc_number": "DOC-2025-001",
      "status": "active",
      "created_at": "2025-11-25T10:00:00Z",
      "updated_at": "2025-11-25T10:00:00Z",
      "permit_type": {
        "id": 1,
        "division_id": 1,
        "category": "Operational",
        "name": "Building Permit",
        "risk_point": "High",
        "default_application_type": "New Application",
        "default_validity_period": "1 year",
        "notes": "Required for construction activities",
        "created_at": "2025-11-25T10:00:00Z",
        "updated_at": "2025-11-25T10:00:00Z",
        "division": {
          "id": 1,
          "code": "IT",
          "name": "Information Technology",
          "created_at": "2025-11-25T10:00:00Z",
          "updated_at": "2025-11-25T10:00:00Z"
        }
      }
    }
  ],
  "pagination": {
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

#### 3. Get Permit by ID
**GET** `/permits/:id`

#### 4. Update Permit
**PUT** `/permits/:id`

#### 5. Delete Permit
**DELETE** `/permits/:id`

---

## Error Responses

### 400 Bad Request
```json
{
  "status": "error",
  "message": "Invalid request body"
}
```

### 404 Not Found
```json
{
  "status": "error",
  "message": "Division not found"
}
```

### 500 Internal Server Error
```json
{
  "status": "error",
  "message": "Internal server error message"
}
```

---

## Database Schema

### Divisions Table
```sql
CREATE TABLE divisions (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Permit Types Table
```sql
CREATE TABLE permit_types (
    id BIGSERIAL PRIMARY KEY,
    division_id BIGINT REFERENCES divisions(id),
    category VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    risk_point VARCHAR(50),
    default_application_type VARCHAR(100),
    default_validity_period VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Permits Table
```sql
CREATE TABLE permits (
    id BIGSERIAL PRIMARY KEY,
    permit_type_id BIGINT NOT NULL REFERENCES permit_types(id),
    application_type VARCHAR(100) NOT NULL,
    permit_no VARCHAR(100) UNIQUE NOT NULL,
    effective_date TIMESTAMP NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    effective_term VARCHAR(100),
    responsible_person VARCHAR(255),
    doc_name VARCHAR(255),
    doc_number VARCHAR(100),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Status Codes

- **200**: OK - Request succeeded
- **201**: Created - Resource created successfully
- **400**: Bad Request - Invalid request parameters
- **404**: Not Found - Resource not found
- **500**: Internal Server Error - Server error occurred
