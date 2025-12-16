package permitController

import (
	"net/http"
	"permit-app/helper"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/permitService"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type PermitController struct {
	service permitService.PermitService
}

func NewPermitController(service permitService.PermitService) *PermitController {
	return &PermitController{service: service}
}

func (c *PermitController) Create(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")
	
	var req model.PermitRequest

	// Handle multipart/form-data separately
	if strings.Contains(contentType, "multipart/form-data") {
		var formReq model.PermitFormRequest
		if err := ctx.ShouldBindWith(&formReq, binding.FormMultipart); err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
			return
		}
		
		// Parse dates
		effectiveDate, err := time.Parse("2006-01-02", formReq.EffectiveDate)
		if err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid effective_date format", err, nil)
			return
		}
		
		expiryDate, err := time.Parse("2006-01-02", formReq.ExpiryDate)
		if err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid expiry_date format", err, nil)
			return
		}
		
		// Convert to PermitRequest
		req = model.PermitRequest{
			DomainID:               formReq.DomainID,
			DivisionID:             formReq.DivisionID,
			PermitTypeID:           formReq.PermitTypeID,
			Name:                   formReq.Name,
			ApplicationType:        formReq.ApplicationType,
			PermitNo:               formReq.PermitNo,
			EffectiveDate:          helper.Date{Time: effectiveDate},
			ExpiryDate:             helper.Date{Time: expiryDate},
			EffectiveTerm:          formReq.EffectiveTerm,
			ResponsiblePersonID:    formReq.ResponsiblePersonID,
			ResponsibleDocPersonID: formReq.ResponsibleDocPersonID,
			DocName:                formReq.DocName,
			DocNumber:              formReq.DocNumber,
			Status:                 formReq.Status,
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := ctx.ShouldBindWith(&req, binding.Form); err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
			return
		}
	} else {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
			return
		}
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	// Validate domain_id from JWT token
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}
	
	if req.DomainID != domainID.(int64) {
		apiresponse.Error(ctx, http.StatusForbidden, "FORBIDDEN", "Cannot create permit in different domain", nil, nil)
		return
	}

	permit, err := c.service.CreatePermit(&req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create permit", err, nil)
		return
	}

	// Handle file upload if present
	file, err := ctx.FormFile("file")
	if err == nil && file != nil {
		permit, err = c.service.HandleFileUpload(permit.ID, file)
		if err != nil {
			apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Permit created but file upload failed", err, nil)
			return
		}
	}

	apiresponse.Created(ctx, permit, "Permit created successfully", nil)
}

func (c *PermitController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	permit, err := c.service.GetPermitByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Permit not found", err, nil)
		return
	}

	apiresponse.OK(ctx, permit, "Permit retrieved successfully", nil)
}

func (c *PermitController) GetAll(ctx *gin.Context) {
	var filter model.PermitListRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	// Extract domain_id from JWT token context
	if domainID, exists := ctx.Get("domain_id"); exists && domainID != nil {
		did := domainID.(int64)
		filter.DomainID = &did
	}

	permits, total, err := c.service.GetAllPermits(&filter)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve permits", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, permits, "Permits retrieved successfully", meta)
}

func (c *PermitController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Search query is required", nil, nil)
		return
	}

	// Build filter from search query
	var filter model.PermitListRequest
	filter.Page = 1
	filter.Limit = 20
	
	// Try to bind other query params (page, limit, etc)
	if err := ctx.ShouldBindQuery(&filter); err == nil {
		// Set search term to multiple fields
		filter.Name = query
		filter.PermitNo = query
	}

	// Extract domain_id from JWT token context
	if domainID, exists := ctx.Get("domain_id"); exists && domainID != nil {
		did := domainID.(int64)
		filter.DomainID = &did
	}

	permits, total, err := c.service.SearchPermits(query, &filter)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to search permits", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, permits, "Search results retrieved successfully", meta)
}

func (c *PermitController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	contentType := ctx.GetHeader("Content-Type")
	var req model.PermitUpdateRequest

	// Handle multipart/form-data separately
	if strings.Contains(contentType, "multipart/form-data") {
		var formReq model.PermitFormUpdateRequest
		if err := ctx.ShouldBindWith(&formReq, binding.FormMultipart); err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
			return
		}
		
		// Parse dates
		effectiveDate, err := time.Parse("2006-01-02", formReq.EffectiveDate)
		if err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid effective_date format", err, nil)
			return
		}
		
		expiryDate, err := time.Parse("2006-01-02", formReq.ExpiryDate)
		if err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid expiry_date format", err, nil)
			return
		}
		
		// Convert to PermitUpdateRequest
		req = model.PermitUpdateRequest{
			DomainID:               formReq.DomainID,
			DivisionID:             formReq.DivisionID,
			PermitTypeID:           formReq.PermitTypeID,
			Name:                   formReq.Name,
			ApplicationType:        formReq.ApplicationType,
			PermitNo:               formReq.PermitNo,
			EffectiveDate:          helper.Date{Time: effectiveDate},
			ExpiryDate:             helper.Date{Time: expiryDate},
			EffectiveTerm:          formReq.EffectiveTerm,
			ResponsiblePersonID:    formReq.ResponsiblePersonID,
			ResponsibleDocPersonID: formReq.ResponsibleDocPersonID,
			DocName:                formReq.DocName,
			DocNumber:              formReq.DocNumber,
			Status:                 formReq.Status,
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := ctx.ShouldBindWith(&req, binding.Form); err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
			return
		}
	} else {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
			return
		}
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	// Validate domain_id from JWT token
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}
	
	if req.DomainID != domainID.(int64) {
		apiresponse.Error(ctx, http.StatusForbidden, "FORBIDDEN", "Cannot update permit in different domain", nil, nil)
		return
	}

	permit, err := c.service.UpdatePermit(id, &req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update permit", err, nil)
		return
	}

	// Handle file upload if present
	file, err := ctx.FormFile("file")
	if err == nil && file != nil {
		permit, err = c.service.HandleFileUpload(permit.ID, file)
		if err != nil {
			apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Permit updated but file upload failed", err, nil)
			return
		}
	}

	apiresponse.OK(ctx, permit, "Permit updated successfully", nil)
}

func (c *PermitController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	err = c.service.DeletePermit(id)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete permit", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Permit deleted successfully", nil)
}

func (c *PermitController) UploadDocument(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "No file uploaded", err, nil)
		return
	}

	permit, err := c.service.HandleFileUpload(id, file)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to upload file", err, nil)
		return
	}

	apiresponse.OK(ctx, permit, "File uploaded successfully", nil)
}

func (c *PermitController) DownloadDocument(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	permit, err := c.service.GetPermitByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Permit not found", err, nil)
		return
	}

	if permit.DocFilePath == nil || *permit.DocFilePath == "" {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "No document found for this permit", nil, nil)
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+*permit.DocFileName)

	if permit.DocFileType != nil {
		ctx.Header("Content-Type", *permit.DocFileType)
	}

	ctx.File(*permit.DocFilePath)
}

func (c *PermitController) PreviewDocument(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	permit, err := c.service.GetPermitByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Permit not found", err, nil)
		return
	}

	if permit.DocFilePath == nil || *permit.DocFilePath == "" {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "No document found for this permit", nil, nil)
		return
	}

	// Set Content-Disposition to inline for preview
	ctx.Header("Content-Disposition", "inline; filename="+*permit.DocFileName)

	if permit.DocFileType != nil {
		ctx.Header("Content-Type", *permit.DocFileType)
	}

	ctx.File(*permit.DocFilePath)
}
