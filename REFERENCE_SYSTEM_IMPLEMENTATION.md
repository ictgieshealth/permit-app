# Reference System Implementation

## Overview
This document describes the implementation of the hierarchical reference management system for the permit application. The system enables dynamic categorization of data across multiple modules (Task, Project, Permit).

## Database Schema

### Tables Created

#### 1. modules
Stores application modules that categorize different functional areas.

```sql
CREATE TABLE modules (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Sample Data:**
- **Task Module** (code: `task`) - Task management functionality
- **Project Module** (code: `project`) - Project management functionality  
- **Permit Module** (code: `permit`) - Permit management functionality

#### 2. reference_categories
Categorizes references within each module.

```sql
CREATE TABLE reference_categories (
    id SERIAL PRIMARY KEY,
    module_id INTEGER NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Sample Data:**
- Task Status, Task Priority, Task Type, Task Stack (under Task module)
- Project Status, Project Type (under Project module)
- Other Categories (under Permit module)

#### 3. references
Individual reference values for each category.

```sql
CREATE TABLE references (
    id SERIAL PRIMARY KEY,
    reference_category_id INTEGER NOT NULL REFERENCES reference_categories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Sample Data:**
- Task Status: To Do, In Progress, Done
- Task Priority: Low, Medium, High, Critical  
- Task Type: Bug Fix, Feature, Enhancement, Research, Documentation
- Project Status: Planning, Active, Paused, Completed, Cancelled
- And more...

## Backend Implementation

### Architecture Layers

#### 1. Models (`backend/model/`)
- **module.go**: Module entity and DTOs
  - `Module` - database entity
  - `ModuleRequest` - creation/update request
  - `ModuleResponse` - API response
  - `ModuleListRequest` - filtering/pagination request

- **referenceCategory.go**: Category entity and DTOs
  - `ReferenceCategory` - database entity with Module relation
  - `ReferenceCategoryRequest` - creation/update request
  - `ReferenceCategoryResponse` - API response with nested Module
  - `ReferenceCategoryListRequest` - filtering/pagination request

- **reference.go**: Reference entity and DTOs
  - `Reference` - database entity with ReferenceCategory relation
  - `ReferenceRequest` - creation/update request
  - `ReferenceResponse` - API response with nested Category and Module
  - `ReferenceListRequest` - filtering/pagination request

#### 2. Repositories (`backend/repo/`)
- **moduleRepository**: Module data access
  - `Create`, `FindAll`, `FindByID`, `FindByCode`, `Update`, `Delete`
  - Filtering: code (ILIKE), name (ILIKE), is_active (exact)
  
- **referenceCategoryRepository**: Category data access
  - `Create`, `FindAll`, `FindByID`, `FindByModuleID`, `Update`, `Delete`
  - Automatic Module relation preloading
  
- **referenceRepository**: Reference data access
  - `Create`, `FindAll`, `FindByID`, `FindByCategoryID`, `FindByModuleID`, `Update`, `Delete`
  - Nested preloading: Category → Module
  - Complex queries with JOIN for module_id filtering

#### 3. Services (`backend/service/`)
- **moduleService**: Module business logic
  - Validation: Code uniqueness check, module existence
  - Returns: Array of responses + total count for pagination
  
- **referenceCategoryService**: Category business logic
  - Cross-validation: Verifies module exists before creating category
  - GetCategoriesByModuleID for filtered access
  - Returns: Array of responses + total count
  
- **referenceService**: Reference business logic
  - Validation: Category existence validation
  - Multi-level queries: Filter by category or module
  - Returns: Array of responses + total count

#### 4. Controllers (`backend/controller/`)
All controllers follow consistent patterns:
- **moduleController**: HTTP handlers for modules
- **referenceCategoryController**: HTTP handlers for categories
- **referenceController**: HTTP handlers for references

Each provides:
- `Create` - POST endpoint with JSON body validation
- `GetAll` - GET with query parameter filtering and pagination
- `GetByID` - GET by path parameter ID
- `Update` - PUT endpoint with ID + JSON body
- `Delete` - DELETE endpoint with empty struct response

#### 5. Routes (`backend/routes/route.go`)
Routes are grouped under protected middleware:

**Module Routes:**
```
POST   /modules
GET    /modules
GET    /modules/:id
PUT    /modules/:id
DELETE /modules/:id
GET    /modules/:module_id/categories
GET    /modules/:module_id/references
```

**Reference Category Routes:**
```
POST   /reference-categories
GET    /reference-categories
GET    /reference-categories/:id
PUT    /reference-categories/:id
DELETE /reference-categories/:id
GET    /reference-categories/:category_id/references
```

**Reference Routes:**
```
POST   /references
GET    /references
GET    /references/:id
PUT    /references/:id
DELETE /references/:id
```

## API Usage Examples

### 1. Get All Modules
```bash
GET /modules?page=1&limit=10&is_active=true
```

Response:
```json
{
  "message": "Modules retrieved successfully",
  "data": [
    {
      "id": 1,
      "code": "task",
      "name": "Task Module",
      "is_active": true,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 3
  }
}
```

### 2. Get Categories by Module
```bash
GET /modules/1/categories
```

Response:
```json
{
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": 1,
      "module_id": 1,
      "name": "Task Status",
      "is_active": true,
      "module": {
        "id": 1,
        "code": "task",
        "name": "Task Module"
      }
    }
  ]
}
```

### 3. Get References by Category
```bash
GET /reference-categories/1/references
```

Response:
```json
{
  "message": "References retrieved successfully",
  "data": [
    {
      "id": 1,
      "reference_category_id": 1,
      "name": "To Do",
      "is_active": true,
      "category": {
        "id": 1,
        "name": "Task Status",
        "module": {
          "id": 1,
          "code": "task",
          "name": "Task Module"
        }
      }
    }
  ]
}
```

### 4. Create New Module
```bash
POST /modules
Content-Type: application/json

{
  "code": "inventory",
  "name": "Inventory Module",
  "is_active": true
}
```

### 5. Filter References
```bash
GET /references?module_id=1&is_active=true&page=1&limit=20
```

## Data Relationships

```
modules (1)
  ↓
reference_categories (N)
  ↓
references (N)
```

- One Module has many Reference Categories
- One Reference Category has many References
- Cascade delete: Deleting a module removes its categories and references
- Cascade delete: Deleting a category removes its references

## Query Features

### Filtering Support
- **Modules**: Filter by code, name, is_active
- **Categories**: Filter by module_id, name, is_active
- **References**: Filter by reference_category_id, module_id, name, is_active

### Pagination
All list endpoints support:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 10000)
- Returns total count in meta

### ILIKE Search
String fields (code, name) use case-insensitive ILIKE matching for flexible searching.

## Validation

### Request Validation
- JSON body binding for POST/PUT requests
- Query parameter binding for GET requests
- Struct validation using go-playground/validator

### Business Logic Validation
- Module code uniqueness check
- Module existence check before creating categories
- Category existence check before creating references

## Error Handling

Standardized error responses:
- `400 Bad Request` - Invalid input or validation failure
- `404 Not Found` - Resource doesn't exist
- `500 Internal Server Error` - Database or server errors

Error response format:
```json
{
  "message": "Validation failed",
  "error": {
    "code": "BAD_REQUEST",
    "details": "Field 'code' is required" // Only in non-production
  },
  "trace_id": "abc-123"
}
```

## Testing

### Build Verification
```bash
cd backend
go build -o main.exe .
```

### Manual Testing
1. Start the backend server
2. Use Postman collection (if available)
3. Test CRUD operations for each entity
4. Verify cascade deletes
5. Test filtering and pagination
6. Verify nested relation loading

## Future Enhancements

1. **Caching**: Add Redis caching for frequently accessed references
2. **Soft Delete**: Implement soft delete instead of hard delete
3. **Versioning**: Add version tracking for reference changes
4. **Audit Log**: Track who created/modified references
5. **Bulk Operations**: Support bulk create/update/delete
6. **Import/Export**: CSV/JSON import/export functionality
7. **Frontend Integration**: Create UI for managing references
8. **Reference Ordering**: Add display_order field for custom sorting
9. **Reference Metadata**: Support additional JSON metadata field
10. **Reference Dependencies**: Track relationships between references

## Migration Notes

### Database Migration
Run the migration SQL in [backend/database/migration.sql](backend/database/migration.sql) to create tables and seed initial data.

### Code Changes
- Added 3 models with DTOs
- Added 3 repositories with filtering
- Added 3 services with validation
- Added 3 controllers with full CRUD
- Updated routes configuration
- All compilation errors resolved

## Summary

✅ Database schema created with proper relations and constraints
✅ Complete backend implementation following layered architecture
✅ Full CRUD operations for all three entities
✅ Filtering and pagination support
✅ Nested relation loading (Category → Module, Reference → Category → Module)
✅ Proper error handling and validation
✅ API routes configured with authentication middleware
✅ Sample data seeded for testing
✅ Successful compilation with no errors

The reference system is now ready for use and can be extended for Task, Project, and Permit management features.
