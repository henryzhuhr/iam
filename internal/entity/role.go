package entity

import "time"

// Role represents a role in the IAM system.
type Role struct {
	ID          int64     `db:"id"`
	TenantID    int64     `db:"tenant_id"`
	Name        string    `db:"name"`
	Code        string    `db:"code"`
	Type        int8      `db:"type"` // 1=system 2=custom
	Status      int8      `db:"status"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Role type constants.
const (
	RoleTypeSystem int8 = 1
	RoleTypeCustom int8 = 2
)

// Role status constants.
const (
	RoleStatusActive   int8 = 1
	RoleStatusDisabled int8 = 2
)

// RolePermission links a role to a permission.
type RolePermission struct {
	ID           int64     `db:"id"`
	RoleID       int64     `db:"role_id"`
	PermissionID int64     `db:"permission_id"`
	DataScope    string    `db:"data_scope"`
	CreatedAt    time.Time `db:"created_at"`
}

// UserRole links a user to a role.
type UserRole struct {
	ID        int64     `db:"id"`
	TenantID  int64     `db:"tenant_id"`
	UserID    int64     `db:"user_id"`
	RoleID    int64     `db:"role_id"`
	AppCode   string    `db:"app_code"`
	CreatedAt time.Time `db:"created_at"`
}

// RoleConstraint represents SoD (Separation of Duties) constraints.
type RoleConstraint struct {
	ID        int64     `db:"id"`
	TenantID  int64     `db:"tenant_id"`
	Type      int8      `db:"type"` // 1=static SoD 2=dynamic SoD
	RoleA     int64     `db:"role_a"`
	RoleB     int64     `db:"role_b"`
	CreatedAt time.Time `db:"created_at"`
}
