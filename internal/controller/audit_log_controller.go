package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/service"
)

type AuditLogController struct {
	auditLogService *service.AuditLogService
}

func NewAuditLogController(auditLogService *service.AuditLogService) *AuditLogController {
	return &AuditLogController{auditLogService: auditLogService}
}

func (ctrl *AuditLogController) GetAdminLogs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	logs, err := ctrl.auditLogService.GetAdminLogs(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
