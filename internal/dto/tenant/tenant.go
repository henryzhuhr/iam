// Package tenant provides DTO definitions for the tenant module.
package tenant

import "time"

// CreateTenantRequest is the request body for creating a tenant.
type CreateTenantRequest struct {
	Name     string `json:"name"`
	MaxUsers int    `json:"max_users"`
	MaxApps  int    `json:"max_apps"`
	ExpireAt string `json:"expire_at,omitempty"`
}

// UpdateTenantRequest is the request body for updating a tenant.
type UpdateTenantRequest struct {
	Name     string `json:"name,omitempty"`
	MaxUsers *int   `json:"max_users,omitempty"`
	MaxApps  *int   `json:"max_apps,omitempty"`
	ExpireAt string `json:"expire_at,omitempty"`
}

// TenantResponse is the response body for a single tenant.
type TenantResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    int8      `json:"status"`
	MaxUsers  int       `json:"max_users"`
	MaxApps   int       `json:"max_apps"`
	ExpireAt  time.Time `json:"expire_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TenantListResponse is the response body for listing tenants.
type TenantListResponse struct {
	Items    []TenantResponse `json:"items"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}
