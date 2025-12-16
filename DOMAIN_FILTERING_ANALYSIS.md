# Domain Filtering Analysis Report

> **Date:** December 15, 2025  
> **Purpose:** Comprehensive analysis of domain filtering implementation across backend and frontend

---

## 📊 Executive Summary

### Current Status

| Layer | Status | Coverage |
|-------|--------|----------|
| **Backend Services** | ✅ **COMPLETE** | 4/4 services support domain_id |
| **Backend Repositories** | ✅ **COMPLETE** | All repos filter by domain_id |
| **Backend Controllers** | ❌ **MISSING** | 0/4 controllers extract domain from context |
| **Frontend Services** | ✅ **COMPLETE** | 4/4 services send domain_id |

### Critical Gap Identified

**Controllers DO NOT extract `domain_id` from JWT token context and pass it to services!**

Currently, domain filtering only works if frontend explicitly sends `domain_id` in query parameters. The JWT token contains `domain_id`, but controllers are not using it.

---

## 🔍 Detailed Analysis

### 1. Backend Controllers Analysis

#### Controllers Requiring Domain Filtering

| Controller | Model Has domain_id | Extracts from Context | Passes to Service | Status |
|------------|--------------------|-----------------------|-------------------|--------|
| **divisionController** | ✅ Yes | ❌ **NO** | ❌ **NO** | ⚠️ **CRITICAL** |
| **permitController** | ✅ Yes | ❌ **NO** | ❌ **NO** | ⚠️ **CRITICAL** |
| **projectController** | ✅ Yes | ❌ **NO** | ❌ **NO** | ⚠️ **CRITICAL** |
| **userController** | ✅ Yes (via relation) | ❌ **NO** | ❌ **NO** | ⚠️ **CRITICAL** |
| **notificationController** | ⚠️ Indirect | 🟡 User-only | 🟡 User-only | ⚠️ **PARTIAL** |

#### What's Missing in Controllers

```go
// CURRENT CODE (WRONG - Missing domain extraction)
func (c *DivisionController) GetAll(ctx *gin.Context) {
    var filter model.DivisionListRequest
    ctx.ShouldBindQuery(&filter)
    divisions, total, err := c.service.GetAllDivisions(&filter)
    // ...
}

// SHOULD BE (CORRECT - With domain extraction)
func (c *DivisionController) GetAll(ctx *gin.Context) {
    var filter model.DivisionListRequest
    ctx.ShouldBindQuery(&filter)
    
    // Extract domain_id from JWT token context
    domainID, exists := ctx.Get("domain_id")
    if exists && domainID != nil {
        did := domainID.(int64)
        filter.DomainID = &did  // Set domain filter from token
    }
    
    divisions, total, err := c.service.GetAllDivisions(&filter)
    // ...
}
```

#### Affected Controller Methods

**divisionController:**
- ❌ `GetAll()` - Line 60-80
- ❌ `Create()` - Line 21-40
- ❌ `GetByID()` - Line 42-58
- ❌ `Update()` - Line 82-108
- ❌ `Delete()` - Line 110-124

**permitController:**
- ❌ `GetAll()` - Line 125-146
- ❌ `Search()` - Line 148-175
- ❌ `Create()` - Line 25-104
- ❌ `GetByID()` - Line 106-123
- ❌ `Update()` - Line 177-263
- ❌ `Delete()` - Line 265-279

**projectController:**
- ❌ `GetProjects()` - Line 69-93
- ❌ `GetProjectByID()` - Line 103-119
- ❌ `CreateProject()` - Line 34-53
- ❌ `UpdateProject()` - Line 179-200
- ❌ `DeleteProject()` - Line 212-228
- ❌ `ChangeProjectStatus()` - Line 240-270

**userController:**
- ❌ `GetAll()` - Line 169-190
- ❌ `Register()` - Line 127-145
- ❌ `GetByID()` - Line 147-167
- ❌ `Update()` - Line 192-219

---

### 2. Backend Services Analysis

#### Services Supporting Domain Filtering

| Service | Interface | GetAll Signature | domain_id Type | Filter Pattern |
|---------|-----------|------------------|----------------|----------------|
| **divisionService** | ✅ Complete | `GetAllDivisions(filter *model.DivisionListRequest)` | `*int64` | Direct field |
| **permitService** | ✅ Complete | `GetAllPermits(filter *model.PermitListRequest)` | `*int64` | Direct field |
| **projectService** | ✅ Complete | `GetProjects(req model.ProjectListRequest)` | `int64` | Map-based |
| **userService** | ✅ Complete | `GetAllUsers(filter *model.UserListRequest)` | `*int64` | JOIN relation |

#### Service Implementation Examples

**divisionService:**
```go
// Interface
GetAllDivisions(filter *model.DivisionListRequest) ([]model.DivisionResponse, int64, error)

// Model
type DivisionListRequest struct {
    DomainID *int64 `form:"domain_id"`  // ✅ Accepts domain_id
    Code     string `form:"code"`
    Name     string `form:"name"`
    Page     int    `form:"page"`
    Limit    int    `form:"limit"`
}

// Repository uses it
if filter.DomainID != nil {
    query = query.Where("domain_id = ?", *filter.DomainID)
}
```

**projectService:**
```go
// Uses map-based filters
if req.DomainID > 0 {
    filters["domain_id"] = req.DomainID
}
projects, total, err := s.projectRepo.FindAll(filters, req.Page, req.Limit)
```

**userService:**
```go
// Filters via JOIN on UserDomainRole
if filter.DomainID != nil {
    query = query.Joins("JOIN user_domain_roles ON user_domain_roles.user_id = users.id").
        Where("user_domain_roles.domain_id = ?", *filter.DomainID)
}
```

---

### 3. Backend Repositories Analysis

All repositories properly implement domain filtering when domain_id is provided:

| Repository | Filtering Logic | Status |
|------------|----------------|--------|
| **divisionRepository** | `WHERE domain_id = ?` | ✅ Complete |
| **permitRepository** | `WHERE domain_id = ?` | ✅ Complete |
| **projectRepository** | `WHERE domain_id = ?` | ✅ Complete |
| **userRepository** | `JOIN user_domain_roles WHERE domain_id = ?` | ✅ Complete |

---

### 4. Frontend Services Analysis

| Service | Method | Sends domain_id | How | Status |
|---------|--------|----------------|-----|--------|
| **division.service.ts** | `getAll()` | ✅ Yes | Query param | ✅ Complete |
| **permit.service.ts** | `getAll()` | ✅ Yes | Query param | ✅ Complete |
| **project.service.ts** | `getAll()` | ✅ Yes | Query param | ✅ Complete |
| **user.service.ts** | `getAll()` | ✅ Yes | Query param | ✅ Complete |

#### Frontend Implementation Pattern

```typescript
// Example from project.service.ts
async getAll(params: {
  domain_id?: number;
  project_status_id?: number;
  // ...
}) {
  const queryParams = new URLSearchParams();
  if (params.domain_id) {
    queryParams.append('domain_id', params.domain_id.toString());
  }
  // ...
  return apiClient.get(`/projects?${queryParams.toString()}`);
}
```

**Issue:** Frontend is manually sending `domain_id` from `authService.getStoredDomain()?.id`, but this is redundant since the JWT token already contains this information.

---

## 🎯 Data Models with Domain Relationships

### Models Requiring Domain Filtering

| Model | Domain Field | Relationship | Notes |
|-------|--------------|--------------|-------|
| **Division** | `domain_id int64` | Direct | MUST filter by domain |
| **Permit** | `domain_id int64` | Direct | MUST filter by domain |
| **Project** | `domain_id int64` | Direct | MUST filter by domain |
| **User** | Via `UserDomainRole` | Indirect | MUST filter by domain via JOIN |

### Models NOT Requiring Domain Filtering

| Model | Reason |
|-------|--------|
| **PermitType** | Global resource (references Division which has domain) |
| **Menu** | Global resource (same for all domains) |
| **Role** | Global resource (same for all domains) |
| **Module** | Global resource |
| **Reference** | Global resource |
| **ReferenceCategory** | Global resource |
| **Notification** | User-scoped (references Permit which has domain) |

---

## ⚠️ Security Implications

### Current Vulnerabilities

1. **No Automatic Domain Isolation**
   - Controllers don't enforce domain context from JWT
   - Frontend must explicitly send domain_id
   - If frontend forgets, data leakage across domains is possible

2. **Manual Domain Filtering is Error-Prone**
   - Each page component must remember to include domain_id
   - Easy to miss in new features
   - No compile-time safety

3. **Inconsistent Authorization**
   - User can potentially request data from other domains by manipulating query params
   - JWT contains domain_id but controllers don't validate against it

### Example Vulnerability

```typescript
// Frontend code (VULNERABLE)
// If developer forgets to add domain_id:
const projects = await projectService.getAll({
  // Missing: domain_id: currentDomain.id
  status: true
});
// Result: Returns ALL projects from ALL domains!
```

---

## ✅ Recommended Solution

### Phase 1: Update Controllers (CRITICAL)

Update all controllers to extract domain_id from JWT token context:

```go
// Example: divisionController.go
func (c *DivisionController) GetAll(ctx *gin.Context) {
    var filter model.DivisionListRequest
    if err := ctx.ShouldBindQuery(&filter); err != nil {
        apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
        return
    }
    
    // CRITICAL: Extract domain_id from JWT token
    if domainID, exists := ctx.Get("domain_id"); exists && domainID != nil {
        did := domainID.(int64)
        filter.DomainID = &did  // Override with token's domain
    }
    
    divisions, total, err := c.service.GetAllDivisions(&filter)
    // ... rest of code
}
```

Apply this pattern to:
- ✅ divisionController: `GetAll`, `Create`, `Update`, `Delete`
- ✅ permitController: `GetAll`, `Search`, `Create`, `Update`, `Delete`
- ✅ projectController: `GetProjects`, `Create`, `Update`, `Delete`, `ChangeProjectStatus`
- ✅ userController: `GetAll`, `Register`, `Update`

### Phase 2: Simplify Frontend (OPTIONAL)

Remove manual domain_id from frontend since backend now enforces it:

```typescript
// Before (redundant):
const projects = await projectService.getAll({
  domain_id: currentDomain.id,  // Remove this
  status: true
});

// After (cleaner):
const projects = await projectService.getAll({
  status: true  // Backend auto-adds domain_id from JWT
});
```

### Phase 3: Add Validation (RECOMMENDED)

Add validation to prevent domain_id override attempts:

```go
// In controller
func (c *DivisionController) Create(ctx *gin.Context) {
    var req model.DivisionRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        apiresponse.BadRequest(ctx, /* ... */)
        return
    }
    
    // Extract domain from JWT token
    domainID, exists := ctx.Get("domain_id")
    if !exists || domainID == nil {
        apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
        return
    }
    
    // VALIDATE: Ensure request domain matches token domain
    if req.DomainID != domainID.(int64) {
        apiresponse.Error(ctx, http.StatusForbidden, "FORBIDDEN", "Cannot create resource in different domain", nil, nil)
        return
    }
    
    division, err := c.service.CreateDivision(&req)
    // ...
}
```

---

## 📋 Implementation Checklist

### High Priority (Security Critical)

- [ ] **divisionController.go**
  - [ ] Update `GetAll()` to extract domain_id from context
  - [ ] Update `Create()` to validate domain_id
  - [ ] Update `Update()` to validate domain_id
  - [ ] Update `Delete()` to validate domain_id

- [ ] **permitController.go**
  - [ ] Update `GetAll()` to extract domain_id from context
  - [ ] Update `Search()` to extract domain_id from context
  - [ ] Update `Create()` to validate domain_id
  - [ ] Update `Update()` to validate domain_id
  - [ ] Update `Delete()` to validate domain_id

- [ ] **projectController.go**
  - [ ] Update `GetProjects()` to extract domain_id from context
  - [ ] Update `CreateProject()` to validate domain_id
  - [ ] Update `UpdateProject()` to validate domain_id
  - [ ] Update `DeleteProject()` to validate domain_id
  - [ ] Update `ChangeProjectStatus()` to validate domain_id

- [ ] **userController.go**
  - [ ] Update `GetAll()` to extract domain_id from context
  - [ ] Update `Register()` to validate domain_id in UserDomainRoles
  - [ ] Update `Update()` to validate domain_id in UserDomainRoles

### Medium Priority (Code Quality)

- [ ] Create helper function for domain extraction
- [ ] Add unit tests for domain filtering
- [ ] Update API documentation with domain context behavior

### Low Priority (Optimization)

- [ ] Remove redundant domain_id from frontend service calls
- [ ] Add TypeScript types for domain-aware API responses
- [ ] Performance test domain filtering queries

---

## 🔐 JWT Token Structure

Current JWT token payload (set by authMiddleware):

```go
// In authMiddleware.go
ctx.Set("user_id", int64(uid))       // ✅ Available
ctx.Set("domain_id", int64(did))     // ✅ Available
ctx.Set("role_id", uint(rid))        // ✅ Available
ctx.Set("username", username)        // ✅ Available
ctx.Set("email", email)              // ✅ Available
```

**All necessary information is already in the context**, controllers just need to use it!

---

## 📚 References

- Auth Middleware: `backend/middleware/authMiddleware.go` (Lines 40-60)
- JWT Helper: `backend/helper/jwt.go`
- Service Layer: All services in `backend/service/*/`
- Repository Layer: All repos in `backend/repo/*/`
- Frontend Services: `frontend/src/services/*.service.ts`

---

## 💡 Conclusion

**The backend infrastructure is 90% ready for domain filtering:**
- ✅ Services accept domain_id
- ✅ Repositories filter by domain_id
- ✅ JWT contains domain_id
- ✅ AuthMiddleware extracts domain_id

**The missing 10% is critical:**
- ❌ Controllers don't extract domain_id from context
- ❌ No automatic domain enforcement
- ❌ Relies on frontend to send domain_id (security risk)

**Solution:** Update 4 controllers (~20 methods total) to extract `domain_id` from Gin context and enforce domain boundaries. This is a **high-priority security fix**.
