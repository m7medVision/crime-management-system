package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/auth"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
	"gorm.io/gorm"
)

func BasicAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		username, password, err := auth.ParseBasicAuth(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication credentials"})
			return
		}

		userRepo := repository.NewUserRepository(db)
		user, err := userRepo.GetByUsername(username)
		if err != nil || !auth.CheckPasswordHash(password, user.Password) || !user.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials or inactive account"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func RequireRole(roles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		user, ok := userInterface.(*model.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		allowed := false
		for _, role := range roles {
			if user.Role == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}

func RequireClearance(minClearance model.ClearanceLevel) gin.HandlerFunc {
	clearanceLevels := map[model.ClearanceLevel]int{
		model.ClearanceLow:      1,
		model.ClearanceMedium:   2,
		model.ClearanceHigh:     3,
		model.ClearanceCritical: 4,
	}

	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		user, ok := userInterface.(*model.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		userClearance := clearanceLevels[user.ClearanceLevel]
		requiredClearance := clearanceLevels[minClearance]

		if userClearance < requiredClearance {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient clearance level"})
			return
		}

		c.Next()
	}
}

func AuditLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log only write operations
		method := c.Request.Method
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete || method == http.MethodPatch {
			// Extract path and user info for audit logging
			path := c.Request.URL.Path
			userInterface, exists := c.Get("user")
			userID := uint(0)
			if exists {
				if user, ok := userInterface.(*model.User); ok {
					userID = user.ID
				}
			}

			// In a real application, you would save this to your audit log repository
			action := "unknown"
			switch method {
			case http.MethodPost:
				action = "create"
			case http.MethodPut, http.MethodPatch:
				action = "update"
			case http.MethodDelete:
				action = "delete"
			}

			// This is just for demonstration
			// In reality, you would call the audit log repository to save this information
			_ = gin.H{
				"user_id":  userID,
				"action":   action,
				"path":     path,
				"status":   c.Writer.Status(),
				"ip":       c.ClientIP(),
				"response": c.Errors.String(),
			}
		}
	}
}
