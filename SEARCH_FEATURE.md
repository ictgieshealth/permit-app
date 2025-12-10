# Permit Search Feature

## Overview

Fitur search untuk Permit Management yang memungkinkan pencarian permits berdasarkan berbagai field secara global dari search bar di header aplikasi.

## Backend Implementation

### Endpoint

```
GET /permits/search?q={query}&page={page}&limit={limit}
```

**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `q` (required) - Search query string
- `page` (optional) - Page number, default: 1
- `limit` (optional) - Items per page, default: 20
- `domain_id` (optional) - Filter by domain
- `division_id` (optional) - Filter by division
- `permit_type_id` (optional) - Filter by permit type
- `status` (optional) - Filter by status (active, expired, inactive)

**Search Fields:**
Pencarian akan dilakukan pada field-field berikut:
- `name` - Nama permit/equipment/services
- `permit_no` - Nomor permit
- `doc_name` - Nama dokumen
- `doc_number` - Nomor dokumen

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "name": "Izin Operasional Alat Medis",
      "permit_no": "001/PER/2024",
      "domain": { ... },
      "division": { ... },
      "permit_type": { ... },
      "responsible_person": { ... },
      "status": "active",
      ...
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 45
  },
  "message": "Search results retrieved successfully"
}
```

### Repository Layer

**File:** `backend/repo/permitRepository/permitRepository.go`

```go
func (r *permitRepository) Search(query string, filter *model.PermitListRequest) ([]model.Permit, int64, error)
```

- Melakukan pencarian dengan `LIKE` query pada multiple fields
- Mendukung additional filters (domain, division, permit_type, status)
- Pagination support
- Preload all relations (Domain, Division, PermitType, ResponsiblePerson, ResponsibleDocPerson)

### Service Layer

**File:** `backend/service/permitService/permitService.go`

```go
func (s *permitService) SearchPermits(query string, filter *model.PermitListRequest) ([]model.PermitResponse, int64, error)
```

- Memanggil repository search
- Convert results ke PermitResponse format

### Controller Layer

**File:** `backend/controller/permitController/permitController.go`

```go
func (c *PermitController) Search(ctx *gin.Context)
```

- Validate query parameter `q` (required)
- Parse pagination parameters
- Call service layer
- Return formatted response

### Routes

**File:** `backend/routes/route.go`

```go
permit.GET("/search", permitCtrl.Search)
```

⚠️ **Important:** Route `/search` harus didefinisikan **SEBELUM** `/:id` untuk menghindari conflict routing.

## Frontend Implementation

### Search Bar Component

**File:** `frontend/src/layout/AppHeader.tsx`

**Features:**
- Search input dengan placeholder "Search permits by name, number, or document..."
- Keyboard shortcut: `Cmd/Ctrl + K` untuk focus ke search bar
- Submit dengan Enter key
- Redirect ke `/permits?search={query}` setelah submit
- Clear input setelah submit

**Implementation:**
```tsx
const [searchQuery, setSearchQuery] = useState("");
const router = useRouter();

const handleSearch = (e: React.FormEvent) => {
  e.preventDefault();
  if (searchQuery.trim()) {
    router.push(`/permits?search=${encodeURIComponent(searchQuery.trim())}`);
    setSearchQuery("");
  }
};
```

### Permits Page

**File:** `frontend/src/app/(admin)/permits/page.tsx`

**Features:**
- Deteksi search query dari URL parameters
- Auto-load results ketika ada search query
- Show "Search results for {query}" di subtitle
- Button "Clear Search" untuk kembali ke list normal
- Hide filters ketika dalam mode search
- Pagination support untuk search results

**Implementation:**
```tsx
const searchParams = useSearchParams();
const searchQuery = searchParams.get("search");

useEffect(() => {
  loadPermits();
}, [page, filters, searchQuery]);

const loadPermits = async () => {
  if (searchQuery) {
    const response = await permitService.search(searchQuery, { page, limit });
    // ...
  } else {
    const response = await permitService.getAll({ page, limit, ...filters });
    // ...
  }
};
```

### Service Layer

**File:** `frontend/src/services/permit.service.ts`

```typescript
async search(
  query: string,
  params: { page?: number; limit?: number } = {}
): Promise<ApiResponse<Permit[]> & { meta?: { ... } }>
```

- Encode query parameter dengan `encodeURIComponent`
- Support pagination
- Return response dengan meta information

## Usage Examples

### From Header Search Bar

1. User ketik query di search bar (misalnya: "Izin Alat")
2. Press Enter atau klik submit
3. Redirect ke `/permits?search=Izin%20Alat`
4. Page load dengan hasil pencarian

### API Call Example

**Request:**
```bash
GET /permits/search?q=alat&page=1&limit=10
Authorization: Bearer {token}
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 5,
      "name": "Izin Operasional Alat Radiologi",
      "permit_no": "005/RAD/2024",
      "domain": {
        "id": 1,
        "name": "Medical Equipment"
      },
      "status": "active"
    },
    {
      "id": 12,
      "name": "Izin Penggunaan Alat Sterilisasi",
      "permit_no": "012/STE/2024",
      "domain": {
        "id": 2,
        "name": "Laboratory Equipment"
      },
      "status": "active"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 2
  },
  "message": "Search results retrieved successfully"
}
```

## Search Behavior

### Case Insensitive
Search dilakukan case-insensitive (tidak membedakan huruf besar/kecil).

### Partial Match
Search menggunakan `LIKE %query%`, sehingga akan match:
- Di awal string: "Izin" akan match "Izin Alat Medis"
- Di tengah string: "Alat" akan match "Izin Alat Medis"
- Di akhir string: "Medis" akan match "Izin Alat Medis"

### Multiple Fields
Satu query akan mencari di semua field:
- name OR permit_no OR doc_name OR doc_number

### With Filters
Search bisa dikombinasikan dengan filters:
```
GET /permits/search?q=alat&status=active&domain_id=1
```

## UI/UX Features

### Search Indicator
- Subtitle berubah: "Search results for {query}"
- Button "Clear Search" muncul untuk reset
- Filter form disembunyikan saat search active

### Empty State
Jika tidak ada hasil:
```tsx
{permits.length === 0 && !loading && (
  <div className="text-center py-8 text-gray-500">
    {searchQuery 
      ? `No permits found for "${searchQuery}"`
      : "No permits found"
    }
  </div>
)}
```

### Keyboard Shortcuts
- `Cmd/Ctrl + K` - Focus search bar
- `Enter` - Submit search
- `Esc` - Clear search input

## Performance Considerations

### Database Indexes
Untuk performance optimal, tambahkan indexes di database:
```sql
CREATE INDEX idx_permits_name ON permits(name);
CREATE INDEX idx_permits_permit_no ON permits(permit_no);
CREATE INDEX idx_permits_doc_name ON permits(doc_name);
CREATE INDEX idx_permits_doc_number ON permits(doc_number);
```

### Pagination
- Default limit: 20 items
- Maksimal bisa disesuaikan via query parameter
- Client-side pagination dengan page navigation

## Testing

### Manual Testing Steps

1. **Basic Search**
   - Ketik "izin" di search bar
   - Verify redirect ke `/permits?search=izin`
   - Verify hasil yang relevan muncul

2. **No Results**
   - Ketik query yang tidak ada (misalnya: "xyzabc123")
   - Verify empty state muncul

3. **Clear Search**
   - Lakukan search
   - Klik "Clear Search"
   - Verify kembali ke list normal

4. **Keyboard Shortcut**
   - Press `Cmd/Ctrl + K`
   - Verify cursor di search bar

5. **Pagination**
   - Search dengan banyak hasil
   - Navigate ke page 2, 3, dst
   - Verify pagination bekerja

## Troubleshooting

### Search Not Working
1. Cek backend running: `air` atau `go run main.go`
2. Cek route registration di `routes/route.go`
3. Cek search endpoint urutan (harus sebelum `/:id`)
4. Cek authorization token valid

### No Results
1. Cek query tidak kosong/whitespace
2. Cek database memiliki data yang match
3. Cek case sensitivity di database
4. Check logs untuk SQL query yang dijalankan

### Frontend Not Redirecting
1. Cek `useRouter` dari `next/navigation`
2. Cek `searchQuery.trim()` tidak kosong
3. Check browser console untuk errors

## Future Enhancements

Potensi improvement untuk fitur search:

1. **Advanced Search**
   - Search dengan multiple terms (AND/OR logic)
   - Date range search
   - Regex support

2. **Search Suggestions**
   - Autocomplete dropdown
   - Recent searches
   - Popular searches

3. **Search Filters in Results**
   - Apply filters tanpa clear search
   - Combine search + filters

4. **Export Search Results**
   - Export hasil pencarian ke Excel
   - Download filtered data

5. **Search Analytics**
   - Track popular search terms
   - Improve search relevance

6. **Full Text Search**
   - PostgreSQL Full Text Search
   - Better ranking algorithm
   - Support untuk Indonesian language
