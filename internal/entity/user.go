package entity

import "time"

// User represents a user in the IAM system.
type User struct {
	ID                int64      `db:"id"`
	TenantID          int64      `db:"tenant_id"`
	Email             string     `db:"email"`
	Phone             string     `db:"phone"`
	PasswordHash      string     `db:"password_hash"`
	Status            int8       `db:"status"` // 1=active 2=disabled 3=locked
	MFAEnabled        int8       `db:"mfa_enabled"`
	MFASecret         string     `db:"mfa_secret"`
	LastLoginAt       *time.Time `db:"last_login_at"`
	PasswordChangedAt *time.Time `db:"password_changed_at"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

// User status constants.
const (
	UserStatusActive   int8 = 1
	UserStatusDisabled int8 = 2
	UserStatusLocked   int8 = 3
)
