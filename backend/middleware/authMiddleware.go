package middleware

import (
	"net/http"
	"permit-app/helper"
	"permit-app/helper/apiresponse"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract token from Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header required", nil, nil)
			ctx.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization format", nil, nil)
			ctx.Abort()
			return
		}

		// Set token untuk VerifyToken
		ctx.Request.Header.Set("Authorization", authHeader)

		// Validate and parse token
		claims, err := helper.VerifyToken(ctx)
		if err != nil {
			apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", err, nil)
			ctx.Abort()
			return
		}

	// Set claims in context for downstream handlers
	// Convert numeric claims from float64 to int64 (JWT standard numeric type)
	if userID, ok := claims["user_id"]; ok {
		if uid, ok := userID.(float64); ok {
			ctx.Set("user_id", int64(uid))
		}
	}
	if domainID, ok := claims["domain_id"]; ok {
		if did, ok := domainID.(float64); ok {
			ctx.Set("domain_id", int64(did))
		}
	}
	if roleID, ok := claims["role_id"]; ok {
		if rid, ok := roleID.(float64); ok {
			ctx.Set("role_id", uint(rid))
		}
	}
	if username, ok := claims["username"]; ok {
		ctx.Set("username", username)
	}
	if email, ok := claims["email"]; ok {
		ctx.Set("email", email)
	}

	ctx.Next()
	}
}

// RequireRole checks if the authenticated user has the required role
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get role from context (set by AuthMiddleware)
		roleID, exists := ctx.Get("role_id")
		if !exists {
			apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil, nil)
			ctx.Abort()
			return
		}

		// In a production system, you would fetch the actual role name from database
		// For now, we'll use a simple role ID check
		// You can enhance this by creating a role cache or fetching from DB

		// For demonstration, assuming role IDs: 1=admin, 2=manager, 3=employee, 4=viewer
		roleMap := map[int64]string{
			1: "admin",
			2: "manager",
			3: "employee",
			4: "viewer",
		}

		userRole := roleMap[roleID.(int64)]
		
		// Check if user's role is in allowed roles
		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			apiresponse.Error(ctx, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions", nil, nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(ctx *gin.Context) (int64, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int64), true
}

// GetDomainIDFromContext retrieves domain ID from context
func GetDomainIDFromContext(ctx *gin.Context) (int64, bool) {
	domainID, exists := ctx.Get("domain_id")
	if !exists {
		return 0, false
	}
	return domainID.(int64), true
}

// GetRoleIDFromContext retrieves role ID from context
func GetRoleIDFromContext(ctx *gin.Context) (int64, bool) {
	roleID, exists := ctx.Get("role_id")
	if !exists {
		return 0, false
	}
	return roleID.(int64), true
}
