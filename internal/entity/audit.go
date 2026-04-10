package entity

import "time"

// AuditLog represents an audit log entry.
type AuditLog struct {
	ID           int64     `db:"id"`
	TenantID     int64     `db:"tenant_id"`
	UserID       int64     `db:"user_id"`
	Action       string    `db:"action"`
	ResourceType string    `db:"resource_type"`
	ResourceID   *int64    `db:"resource_id"`
	Detail       string    `db:"detail"` // JSON string
	IP           string    `db:"ip"`
	CreatedAt    time.Time `db:"created_at"`
}

// LoginLog represents a login attempt record.
type LoginLog struct {
	ID         int64     `db:"id"`
	TenantID   int64     `db:"tenant_id"`
	UserID     *int64    `db:"user_id"` // nil on failed login
	Email      string    `db:"email"`
	Status     int8      `db:"status"` // 1=success 2=failed 3=MFA pending
	FailReason string    `db:"fail_reason"`
	LoginType  string    `db:"login_type"`
	IP         string    `db:"ip"`
	UserAgent  string    `db:"user_agent"`
	CreatedAt  time.Time `db:"created_at"`
}
