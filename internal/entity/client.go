package entity

import "time"

// Client represents an internal client (service-to-service).
type Client struct {
	ID            int64     `db:"id"`
	ClientID      string    `db:"client_id"`
	AccessKey     string    `db:"access_key"`
	SecretKeyHash string    `db:"secret_key_hash"`
	Name          string    `db:"name"`
	AllowedScopes string    `db:"allowed_scopes"` // JSON string
	Status        int8      `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
