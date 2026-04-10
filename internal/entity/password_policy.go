package entity

import "time"

// PasswordPolicy represents per-tenant password requirements.
type PasswordPolicy struct {
	ID               int64     `db:"id"`
	TenantID         int64     `db:"tenant_id"`
	MinLength        int       `db:"min_length"`
	RequireUppercase int8      `db:"require_uppercase"`
	RequireLowercase int8      `db:"require_lowercase"`
	RequireDigit     int8      `db:"require_digit"`
	RequireSpecial   int8      `db:"require_special"`
	HistoryCount     int       `db:"history_count"`
	ExpireDays       int       `db:"expire_days"`
	MaxLoginAttempts int       `db:"max_login_attempts"`
	LockoutMinutes   int       `db:"lockout_minutes"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// PasswordHistory tracks previous password hashes.
type PasswordHistory struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
