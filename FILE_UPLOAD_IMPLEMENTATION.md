# File Upload Feature - Implementation Summary

## Overview
Implemented complete file upload functionality for permit documents with drag-and-drop support, allowing users to upload, store, and download document attachments.

## Backend Changes

### 1. Database Schema (`backend/database/migration.sql`)
Added 4 new columns to `permits` table:
- `doc_file_name VARCHAR(255)` - Original filename
- `doc_file_path VARCHAR(500)` - Server storage path
- `doc_file_size BIGINT` - File size in bytes
- `doc_file_type VARCHAR(100)` - MIME type

### 2. Data Model (`backend/model/permit.go`)
Updated structs to include file metadata:
- `Permit` - Added 4 file fields with GORM tags
- `PermitRequest` - Added file fields for create requests
- `PermitUpdateRequest` - Added file fields for update requests
- `PermitResponse` - Added file fields for API responses

### 3. File Storage Helper (`backend/helper/fileStorage.go`)
Created new helper with functions:
- `SaveFile()` - Saves uploaded file with unique filename
- `DeleteFile()` - Removes file from storage
- `ValidateFileType()` - Validates file extension
- `GetMimeType()` - Returns MIME type from extension

**Configuration:**
- Allowed types: PDF, DOC, DOCX, JPG, JPEG, PNG
- Max file size: 10MB
- Storage path: `file/permits/`
- Filename format: `YYYYMMDDHHMMSS_<uuid>.<ext>`

### 4. Controller (`backend/controller/permitController/permitController.go`)
Added 2 new endpoints:
- `UploadDocument()` - POST `/permits/:id/upload`
  - Accepts multipart/form-data with "file" field
  - Validates file type and size
  - Updates permit record with file metadata
  
- `DownloadDocument()` - GET `/permits/:id/download`
  - Streams file to client
  - Sets proper Content-Type and Content-Disposition headers

### 5. Service (`backend/service/permitService/permitService.go`)
Enhanced service layer:
- `HandleFileUpload()` - Manages file upload process
  - Deletes old file if exists
  - Saves new file
  - Updates database with file metadata
  - Rollback on error
  
- Updated `DeletePermit()` - Now deletes associated file
- Updated `toResponse()` - Includes file fields in response

### 6. Routes (`backend/routes/route.go`)
Added routes to permit group:
```go
permit.POST("/:id/upload", permitCtrl.UploadDocument)
permit.GET("/:id/download", permitCtrl.DownloadDocument)
```

### 7. Dependencies (`backend/go.mod`)
Added package:
- `github.com/google/uuid v1.6.0` - For unique filename generation

## Frontend Changes

### 1. FileUpload Component (`frontend/src/components/form/FileUpload.tsx`)
Created reusable component with features:
- **Drag-and-drop zone** with visual feedback
- **File validation** (type, size)
- **File preview** (name, size display)
- **Remove file** capability
- **Error handling** with user-friendly messages
- **Dark mode** support

Props:
```typescript
interface FileUploadProps {
  onFileSelect: (file: File | null) => void;
  currentFile?: { name: string; size: number; type: string };
  accept?: string;
  maxSize?: number;
  disabled?: boolean;
}
```

### 2. Permit Service (`frontend/src/services/permit.service.ts`)
Added methods:
- `uploadDocument(id, file)` - Uploads file using FormData
- `downloadDocument(id)` - Downloads file and triggers browser download

### 3. Permit Create Form (`frontend/src/app/(admin)/permits/create/page.tsx`)
Integrated file upload:
- Added `uploadedFile` state
- Added FileUpload component after document fields
- Modified `handleSubmit` to upload file after permit creation
- Shows allowed file types hint

## Usage Flow

### Creating Permit with File:
1. User fills permit form
2. User drags/drops or clicks to select file
3. File validates (type, size)
4. User submits form
5. Backend creates permit record
6. Backend uploads file if selected
7. File saved to `file/permits/` directory
8. Database updated with file metadata
9. User redirected to permits list

### Downloading File:
1. User clicks download button (to be added to UI)
2. Backend retrieves file path from database
3. Backend streams file with proper headers
4. Browser downloads file with original filename

## Security Features
- File type validation (whitelist)
- File size limit (10MB)
- Unique filename generation (prevents collisions)
- Database transaction rollback on upload failure
- Old file cleanup on replacement/deletion

## File Storage Structure
```
backend/
  file/
    permits/
      20250605143022_a1b2c3d4-e5f6-7890-abcd-ef1234567890.pdf
      20250605143045_b2c3d4e5-f6a7-8901-bcde-f12345678901.jpg
```

## Next Steps

### Suggested Enhancements:
1. **Edit Form Integration** - Add file upload to permit edit page
2. **File Download Button** - Add download button to permit view/list
3. **File Preview** - Show existing file info when editing
4. **File Replace** - Allow replacing file during edit
5. **Cloud Storage** - Integrate with AWS S3 or Azure Blob Storage
6. **Image Thumbnails** - Generate thumbnails for image files
7. **Virus Scanning** - Integrate with antivirus API
8. **File Compression** - Compress large files before storage

### Database Migration:
Run this SQL to add columns to existing database:
```sql
ALTER TABLE permits ADD COLUMN doc_file_name VARCHAR(255);
ALTER TABLE permits ADD COLUMN doc_file_path VARCHAR(500);
ALTER TABLE permits ADD COLUMN doc_file_size BIGINT;
ALTER TABLE permits ADD COLUMN doc_file_type VARCHAR(100);
```

Or rerun full migration:
```bash
psql -U postgres -d permit_db -f backend/database/migration.sql
```

## Testing Checklist
- [ ] Upload PDF file (< 10MB)
- [ ] Upload image file (JPG, PNG)
- [ ] Upload DOC/DOCX file
- [ ] Try uploading invalid file type (should fail)
- [ ] Try uploading file > 10MB (should fail)
- [ ] Download uploaded file
- [ ] Verify file integrity after download
- [ ] Upload file via drag-and-drop
- [ ] Remove selected file before submit
- [ ] Edit permit and replace file
- [ ] Delete permit (verify file is deleted)
- [ ] Check dark mode styling
- [ ] Test on mobile/tablet

## API Endpoints

### Upload Document
```http
POST /permits/:id/upload
Content-Type: multipart/form-data

file: <binary>
```

Response:
```json
{
  "message": "File uploaded successfully",
  "data": {
    "id": 1,
    "name": "Sample Permit",
    "doc_file_name": "contract.pdf",
    "doc_file_size": 2048576,
    "doc_file_type": "application/pdf",
    ...
  }
}
```

### Download Document
```http
GET /permits/:id/download
```

Response: Binary file stream with headers:
- `Content-Type`: File MIME type
- `Content-Disposition`: attachment; filename=original_name.pdf

## Dependencies Added
- Backend: `github.com/google/uuid v1.6.0`
- Frontend: None (using built-in browser APIs)

## Files Modified/Created

### Backend (7 files):
- ✅ `database/migration.sql` - Modified
- ✅ `model/permit.go` - Modified
- ✅ `helper/fileStorage.go` - Created
- ✅ `controller/permitController/permitController.go` - Modified
- ✅ `service/permitService/permitService.go` - Modified
- ✅ `routes/route.go` - Modified
- ✅ `go.mod` - Modified

### Frontend (3 files):
- ✅ `components/form/FileUpload.tsx` - Created
- ✅ `services/permit.service.ts` - Modified
- ✅ `app/(admin)/permits/create/page.tsx` - Modified

## Build Status
✅ Backend builds successfully (`go build`)
✅ All TypeScript changes compile
✅ No import errors
✅ Ready for testing
