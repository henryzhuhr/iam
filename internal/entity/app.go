package entity

import "time"

// Application represents an application in the IAM system.
type Application struct {
	ID          int64     `db:"id"`
	TenantID    int64     `db:"tenant_id"`
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Status      int8      `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// UserAppAuthorization links a user to an application.
type UserAppAuthorization struct {
	ID        int64     `db:"id"`
	TenantID  int64     `db:"tenant_id"`
	UserID    int64     `db:"user_id"`
	AppID     int64     `db:"app_id"`
	CreatedAt time.Time `db:"created_at"`
}
