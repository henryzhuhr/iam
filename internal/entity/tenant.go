// Package entity provides database entity models for the IAM application.
package entity

import "time"

// Tenant represents a tenant in the IAM system.
type Tenant struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Status    int8      `db:"status"` // 1=active 2=disabled 3=expired
	MaxUsers  int       `db:"max_users"`
	MaxApps   int       `db:"max_apps"`
	ExpireAt  time.Time `db:"expire_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Tenant status constants.
const (
	TenantStatusActive   int8 = 1
	TenantStatusDisabled int8 = 2
	TenantStatusExpired  int8 = 3
)
