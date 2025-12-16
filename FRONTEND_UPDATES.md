# Frontend Updates untuk Penggabungan Aplikasi

## Overview
Update frontend untuk mendukung struktur baru dari backend dimana:
- User dapat memiliki multiple domains dengan role berbeda untuk setiap domain
- Role memiliki code, category, dan description
- User tidak lagi memiliki single role, tetapi role per domain (user_domain_roles)

## File yang Diupdate

### 1. Type Definitions

#### `frontend/src/types/role.ts`
**Changes:**
- ✅ Tambah field `code: string` - unique identifier untuk role
- ✅ Tambah field `category: string` - kategori role (Permit/Ticketing)
- ✅ Tambah field `description?: string` - deskripsi optional untuk role

```typescript
export interface Role {
  id: number;
  code: string;        // NEW
  name: string;
  category: string;    // NEW
  description?: string; // NEW
  created_at: string;
  updated_at: string;
}
```

#### `frontend/src/types/user.ts`
**Changes:**
- ✅ Tambah interface `UserDomainRole` untuk relasi user-domain-role
- ✅ Remove field `role_id` dari User interface
- ✅ Remove field `role?: Role` dari User interface
- ✅ Remove field `domains?: Domain[]` dari User interface
- ✅ Tambah field `domain_roles?: UserDomainRole[]` ke User interface
- ✅ Tambah interface `DomainRoleRequest` untuk request body
- ✅ Update `UserRequest` untuk menggunakan `domain_roles: DomainRoleRequest[]`
- ✅ Update `UserUpdateRequest` untuk menggunakan `domain_roles?: DomainRoleRequest[]`

```typescript
export interface UserDomainRole {
  user_id: number;
  domain_id: number;
  role_id: number;
  is_default: boolean;
  domain?: Domain;
  role?: Role;
}

export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  phone_number: string;
  nip: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  domain_roles?: UserDomainRole[]; // NEW
}
```

#### `frontend/src/types/loginUser.ts`
**Changes:**
- ✅ Tambah field `role_id?: number` ke LoginRequest
- ✅ Tambah field `default_domain_role?: UserDomainRole` ke LoginResponse

### 2. Services

#### `frontend/src/services/user.service.ts`
**Changes:**
- ✅ Rename `addDomain()` → `addDomainRole()` dengan parameter `roleId`
- ✅ Rename `removeDomain()` → `removeDomainRole()` dengan parameter `roleId`
- ✅ Rename `setDefaultDomain()` → `setDefaultDomainRole()` dengan parameter `roleId`
- ✅ Update endpoint dari `/users/:id/domains` → `/users/:id/domain-roles`

```typescript
async addDomainRole(id: number, domainId: number, roleId: number, isDefault: boolean = false)
async removeDomainRole(id: number, domainId: number, roleId: number)
async setDefaultDomainRole(id: number, domainId: number, roleId: number)
```

#### `frontend/src/services/auth.service.ts`
**Changes:**
- ✅ Store `default_domain_role` ke localStorage saat login
- ✅ Tambah method `getStoredDomainRole()` untuk retrieve default domain role
- ✅ Update `hasRole()` untuk menggunakan role code dari `default_domain_role.role.code`
- ✅ Clear `default_domain_role` saat logout

**Important:** `hasRole()` sekarang menerima array of role codes (bukan role names):
```typescript
// Before
authService.hasRole(["admin"])

// After
authService.hasRole(["ADMIN"])
```

### 3. Hooks

#### `frontend/src/hooks/useUserDomains.ts`
**Changes:**
- ✅ Extract domains dari `user.domain_roles` bukan `user.domains`
- ✅ Create unique domain map dari domain_roles
- ✅ Update `canAccessDomain()` untuk check `domain_roles` array

### 4. Pages - User Management

#### `frontend/src/app/(admin)/users/create/page.tsx`
**Changes:**
- ✅ Remove single `role_id` field dari formData
- ✅ Remove single `domain_ids` array dari formData
- ✅ Tambah `domain_roles` array ke formData
- ✅ Tambah helper functions:
  - `handleDomainRoleToggle()` - toggle domain selection
  - `handleRoleChange()` - set role untuk specific domain
  - `handleDefaultChange()` - set default domain-role
- ✅ Update validation untuk check role selected untuk semua domains
- ✅ Update UI untuk menampilkan domain-role assignment dalam satu section
- ✅ Show role dropdown untuk setiap selected domain
- ✅ Show radio button untuk set default domain-role
- ✅ Update `hasRole` check dari `["admin"]` → `["ADMIN"]`

**UI Changes:**
```
Old Structure:
- Role dropdown (single selection)
- Domain checkboxes (multiple selection)

New Structure:
- Domain & Role Assignments section:
  ├── Domain checkbox
  ├── Role dropdown (per domain)
  └── Default radio button (per domain-role)
```

#### `frontend/src/app/(admin)/users/[id]/edit/page.tsx`
**Changes:**
- ✅ Same changes as create page
- ✅ Update `loadUser()` untuk map `user.domain_roles` ke formData
- ✅ Update `hasRole` check dari `["admin"]` → `["ADMIN"]`

#### `frontend/src/app/(admin)/users/page.tsx`
**Changes:**
- ✅ Update table header dari "Role" + "Domains" → "Domain & Role Assignments"
- ✅ Update table cell untuk display domain_roles dengan format:
  - `DomainCode: RoleCode ⭐` (⭐ untuk default)
  - Tooltip menampilkan full name: "DomainName - RoleName (Default)"
- ✅ Update colspan dari 7 → 6 (merge 2 columns jadi 1)
- ✅ Update `hasRole` check dari `["admin"]` → `["ADMIN"]`

**Display Format:**
```
Old: Role Badge | Domain1, Domain2, Domain3
New: D1:ADMIN⭐, D2:PERMIT_MANAGER
```

### 5. Components yang TIDAK Perlu Diupdate

✅ Menu management pages - sudah menggunakan role_ids (tidak terpengaruh)
✅ Permit management pages - tidak mengakses user role structure
✅ Domain management pages - tidak terpengaruh
✅ Division management pages - tidak terpengaruh
✅ Notification components - tidak terpengaruh

## Breaking Changes

### 1. Role Checking
**Before:**
```typescript
authService.hasRole(["admin", "manager"])
```

**After:**
```typescript
authService.hasRole(["ADMIN", "PERMIT_MANAGER"])
```

### 2. User Object Structure
**Before:**
```typescript
user.role_id
user.role.name
user.domains
```

**After:**
```typescript
user.domain_roles[0].role_id
user.domain_roles[0].role.code
user.domain_roles.map(dr => dr.domain)
```

### 3. API Endpoints
**Before:**
```
POST /users/:id/domains
DELETE /users/:id/domains/:domain_id
PUT /users/:id/domains/:domain_id/set-default
```

**After:**
```
POST /users/:id/domain-roles
DELETE /users/:id/domain-roles/:domain_id/:role_id
PUT /users/:id/domain-roles/:domain_id/:role_id/set-default
```

## Testing Checklist

### User Management
- [ ] Create user dengan multiple domain-role assignments
- [ ] Verify role dropdown hanya muncul untuk selected domains
- [ ] Set different roles untuk different domains
- [ ] Set default domain-role (hanya 1 yang bisa default)
- [ ] Edit user dan update domain-role assignments
- [ ] Remove domain dari user
- [ ] Verify validation: semua selected domain harus punya role
- [ ] Verify validation: minimal 1 domain-role assignment
- [ ] Verify validation: harus ada 1 default domain-role

### Authentication & Authorization
- [ ] Login dengan user yang punya multiple domain-roles
- [ ] Verify default domain-role tersimpan di localStorage
- [ ] Verify hasRole() check menggunakan role code
- [ ] Verify ADMIN role code dapat akses admin pages
- [ ] Test logout dan verify default_domain_role cleared

### Display & UI
- [ ] User list page menampilkan domain-role assignments dengan benar
- [ ] Domain codes dan role codes tampil dengan format: "D1:ADMIN⭐"
- [ ] Tooltip menampilkan full information
- [ ] Default domain ditandai dengan ⭐
- [ ] Create/Edit form menampilkan domain-role UI dengan benar

### Hooks & Context
- [ ] useUserDomains() extract domains dari domain_roles
- [ ] canAccessDomain() check domain_roles array
- [ ] AuthContext store dan retrieve user dengan domain_roles

## Migration Notes

### Database
Pastikan backend migration sudah dijalankan:
```sql
-- user_domain_roles table sudah dibuat
-- roles table sudah ada code, category, description
-- users table sudah tidak ada role_id
```

### Role Codes Mapping
```
Old Role Name → New Role Code
---------------------------------
admin         → ADMIN
permit_manager → PERMIT_MANAGER
permit_employee → PERMIT_EMPLOYEE
ticketing_employee → TICKETING_EMPLOYEE
ticketing_manager → TICKETING_MANAGER
ticketing_developer → TICKETING_DEVELOPER
```

### Categories
- `Permit` - untuk roles terkait permit management
- `Ticketing` - untuk roles terkait ticketing system

## Next Steps

1. ✅ Update backend API selesai
2. ✅ Update frontend types selesai
3. ✅ Update frontend services selesai
4. ✅ Update frontend pages selesai
5. ✅ No compilation errors
6. ⏳ Manual testing diperlukan
7. ⏳ Integration testing dengan backend
8. ⏳ Deploy ke staging environment

## Important Notes

1. **Role Code vs Name**: Sekarang menggunakan role CODE untuk authorization checks, bukan role NAME
2. **Multiple Roles**: User bisa punya role berbeda di domain berbeda (e.g., ADMIN di Permit, EMPLOYEE di Ticketing)
3. **Default Domain-Role**: Saat login, default domain-role digunakan untuk context awal
4. **Backward Compatibility**: Tidak ada - ini breaking change yang memerlukan database migration

## Troubleshooting

### Error: "user.role is undefined"
**Solution:** Update code untuk menggunakan `user.domain_roles` bukan `user.role`

### Error: "hasRole is not working"
**Solution:** Pastikan menggunakan role CODE (uppercase) bukan role name

### Error: "Cannot read property 'domain' of undefined"
**Solution:** Check apakah domain_roles sudah di-load dengan proper eager loading

### UI not showing domain-role selections
**Solution:** Verify domains dan roles sudah di-load via `loadDomains()` dan `loadRoles()`
