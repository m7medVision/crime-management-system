package dto

type AuditLogDTO struct {
	UserID     uint   `json:"userId"`
	Action     string `json:"action"`
	EntityType string `json:"entityType"`
	EntityID   uint   `json:"entityId"`
	OldValue   string `json:"oldValue"`
	NewValue   string `json:"newValue"`
	IPAddress  string `json:"ipAddress"`
}
