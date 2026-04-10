package entity

import "time"

// UserGroup represents a user group.
type UserGroup struct {
	ID          int64     `db:"id"`
	TenantID    int64     `db:"tenant_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// UserGroupMember links a user to a group.
type UserGroupMember struct {
	ID        int64     `db:"id"`
	GroupID   int64     `db:"group_id"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
