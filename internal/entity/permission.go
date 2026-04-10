package entity

import "time"

// Permission represents a permission definition.
type Permission struct {
	ID          int64     `db:"id"`
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	Resource    string    `db:"resource"`
	Action      string    `db:"action"`
	AppCode     string    `db:"app_code"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
