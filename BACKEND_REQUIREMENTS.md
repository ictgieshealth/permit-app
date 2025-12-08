# Backend API Requirements for Menu Management & Role-Based Access Control

## Overview
Sistem menu management dengan role-based access control telah dibuat di frontend. Backend perlu menyediakan endpoint API untuk mendukung fitur-fitur berikut:

## 1. Menu Management API

### Database Schema - Table: `menus`
```sql
CREATE TABLE menus (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    path VARCHAR(255) NOT NULL UNIQUE,
    icon VARCHAR(50),
    parent_id INTEGER REFERENCES menus(id) ON DELETE CASCADE,
    order_index INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_menus_parent_id ON menus(parent_id);
CREATE INDEX idx_menus_is_active ON menus(is_active);
CREATE INDEX idx_menus_order_index ON menus(order_index);
```

### Database Schema - Table: `menu_roles` (Junction Table)
```sql
CREATE TABLE menu_roles (
    id SERIAL PRIMARY KEY,
    menu_id INTEGER NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(menu_id, role_id)
);

CREATE INDEX idx_menu_roles_menu_id ON menu_roles(menu_id);
CREATE INDEX idx_menu_roles_role_id ON menu_roles(role_id);
```

### API Endpoints

#### 1. GET /api/menus
**Description:** Get all menus with pagination and filters
**Query Parameters:**
- `name` (optional): Filter by menu name
- `path` (optional): Filter by menu path
- `is_active` (optional): Filter by active status (true/false)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Dashboard",
      "path": "/",
      "icon": "GridIcon",
      "parent_id": null,
      "order_index": 1,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "roles": [
        {
          "id": 1,
          "name": "Admin"
        }
      ],
      "children": []
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 50
  }
}
```

#### 2. GET /api/menus/:id
**Description:** Get menu by ID with roles
**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "name": "Dashboard",
    "path": "/",
    "icon": "GridIcon",
    "parent_id": null,
    "order_index": 1,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "roles": [
      {
        "id": 1,
        "name": "Admin"
      }
    ]
  }
}
```

#### 3. GET /api/menus/user
**Description:** Get menus accessible by current logged-in user based on their role
**Authentication:** Required (JWT Token)
**Business Logic:**
- Ambil role_id dari user yang sedang login
- Query semua menus yang memiliki relasi dengan role_id tersebut di table menu_roles
- Return menu dengan struktur hierarchical (parent-children)
- Hanya return menu dengan is_active = true
- Sort by order_index ASC

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Dashboard",
      "path": "/",
      "icon": "GridIcon",
      "parent_id": null,
      "order_index": 1,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "children": []
    },
    {
      "id": 2,
      "name": "Permits",
      "path": "/permits",
      "icon": "ListIcon",
      "parent_id": null,
      "order_index": 2,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "children": []
    }
  ]
}
```

#### 4. POST /api/menus
**Description:** Create new menu
**Request Body:**
```json
{
  "name": "Domains",
  "path": "/domains",
  "icon": "PageIcon",
  "parent_id": null,
  "order_index": 5,
  "role_ids": [1, 2, 3]
}
```

**Business Logic:**
- Validate path is unique
- Create menu record
- Create records in menu_roles table for each role_id

**Response:**
```json
{
  "code": 201,
  "message": "Menu created successfully",
  "data": {
    "id": 10,
    "name": "Domains",
    "path": "/domains",
    "icon": "PageIcon",
    "parent_id": null,
    "order_index": 5,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "roles": [
      { "id": 1, "name": "Admin" }
    ]
  }
}
```

#### 5. PUT /api/menus/:id
**Description:** Update menu
**Request Body:**
```json
{
  "name": "Domains Management",
  "path": "/domains",
  "icon": "PageIcon",
  "parent_id": null,
  "order_index": 5,
  "role_ids": [1, 2]
}
```

**Business Logic:**
- Update menu record
- Delete all existing menu_roles for this menu_id
- Create new menu_roles records for provided role_ids

**Response:**
```json
{
  "code": 200,
  "message": "Menu updated successfully",
  "data": {
    "id": 10,
    "name": "Domains Management",
    "path": "/domains",
    "icon": "PageIcon",
    "parent_id": null,
    "order_index": 5,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 6. DELETE /api/menus/:id
**Description:** Delete menu (cascade delete menu_roles)
**Response:**
```json
{
  "code": 200,
  "message": "Menu deleted successfully",
  "data": null
}
```

#### 7. POST /api/menus/:id/roles
**Description:** Assign roles to menu (alternative endpoint)
**Request Body:**
```json
{
  "role_ids": [1, 2, 3]
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Roles assigned successfully",
  "data": null
}
```

## 2. Update Existing User API

### Modify GET /api/auth/profile
**Add domain information to response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "full_name": "Administrator",
    "role_id": 1,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "role": {
      "id": 1,
      "name": "Admin"
    },
    "domains": [
      {
        "id": 1,
        "code": "FIN",
        "name": "Finance",
        "description": "Finance Domain",
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**Note:** Jika user tidak memiliki domain (admin super), return array kosong `"domains": []`

### Update GET /api/roles
**Add is_active filter:**
**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page
- `is_active` (optional): Filter by active status (true/false)

## 3. Sample Data untuk Testing

### Insert Sample Menus:
```sql
INSERT INTO menus (name, path, icon, parent_id, order_index, is_active) VALUES
('Dashboard', '/', 'GridIcon', NULL, 1, true),
('Permits', '/permits', 'ListIcon', NULL, 2, true),
('Domains', '/domains', 'PageIcon', NULL, 3, true),
('Divisions', '/divisions', 'BoxCubeIcon', NULL, 4, true),
('Permit Types', '/permit-types', 'TableIcon', NULL, 5, true),
('Users', '/users', 'UserCircleIcon', NULL, 6, true),
('Menus', '/menus', 'PlugInIcon', NULL, 7, true);
```

### Assign Menus to Admin Role (role_id = 1):
```sql
INSERT INTO menu_roles (menu_id, role_id) VALUES
(1, 1), -- Dashboard
(2, 1), -- Permits
(3, 1), -- Domains
(4, 1), -- Divisions
(5, 1), -- Permit Types
(6, 1), -- Users
(7, 1); -- Menus
```

## 4. Icon Mapping (Frontend Reference)

Icon names yang bisa digunakan (sudah di-map di frontend):
- `GridIcon` - Dashboard icon
- `ListIcon` - List/document icon
- `PageIcon` - Page/file icon
- `BoxCubeIcon` - Box/cube icon
- `UserCircleIcon` - User icon
- `CalenderIcon` - Calendar icon
- `PieChartIcon` - Chart icon
- `PlugInIcon` - Plugin/settings icon
- `TableIcon` - Table icon

## 5. Business Rules

1. **Menu Access Control:**
   - User hanya dapat melihat menu yang role-nya memiliki akses
   - Menu tidak aktif (is_active = false) tidak ditampilkan di sidebar
   - Menu di-sort berdasarkan order_index ascending

2. **Domain Filtering:**
   - Jika user memiliki 1 domain, domain input otomatis ter-select dan disabled
   - Jika user memiliki > 1 domain, user dapat memilih domain
   - Jika user tidak memiliki domain (admin super), dapat melihat semua domain
   - Form yang memiliki input domain harus di-filter: permits, divisions

3. **Hierarchical Menu:**
   - Menu dapat memiliki parent (support 1 level nesting)
   - Parent menu ditampilkan di "Menu" section
   - Non-parent menu bisa di "Master Data" section jika path-nya termasuk master data

4. **Validation:**
   - Path menu harus unique
   - Minimal 1 role harus di-assign ke menu
   - Order index harus angka positif

## 6. Error Handling

Standard error response:
```json
{
  "code": 400,
  "message": "Validation error",
  "data": {
    "errors": {
      "path": "Path already exists",
      "role_ids": "At least one role must be selected"
    }
  }
}
```

Status codes:
- 200: Success
- 201: Created
- 400: Bad Request / Validation Error
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## 7. Migration Priority

1. **Phase 1 (Critical):**
   - Create menus table
   - Create menu_roles table
   - Implement GET /api/menus/user (untuk sidebar)
   - Update GET /api/auth/profile (tambah domains)

2. **Phase 2 (High):**
   - Implement CRUD endpoints untuk menus
   - Seed initial menu data
   - Assign menus to existing roles

3. **Phase 3 (Medium):**
   - Add is_active filter to roles API
   - Testing & validation

## Notes untuk Backend Developer

1. Gunakan JOIN untuk mengambil data roles ketika query menu
2. Implement soft delete jika diperlukan (is_active flag)
3. Pastikan cascade delete pada foreign key menu_roles
4. Index pada kolom yang sering di-query (parent_id, is_active, order_index)
5. Validate circular reference pada parent_id (menu tidak boleh jadi parent dari dirinya sendiri)
6. Cache menu data jika perlu untuk performa (GET /api/menus/user dipanggil setiap page load)

## Testing Checklist

- [ ] Create menu with roles
- [ ] Update menu and change roles
- [ ] Delete menu (verify cascade delete menu_roles)
- [ ] Get user menus by role
- [ ] Filter menus by name, path, is_active
- [ ] Pagination works correctly
- [ ] Hierarchical menu structure (parent-children)
- [ ] Domain filtering in forms
- [ ] User with single domain auto-select
- [ ] User with multiple domains can choose
- [ ] Super admin sees all domains
