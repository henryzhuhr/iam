// Package repository provides data access layer for the IAM application.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"iam/internal/entity"
)

// TenantRepository provides data access for tenant entities.
type TenantRepository struct {
	db *sql.DB
}

// NewTenantRepository creates a new TenantRepository.
func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create inserts a new tenant.
func (r *TenantRepository) Create(ctx context.Context, t *entity.Tenant) error {
	query := `INSERT INTO tenants (name, status, max_users, max_apps, expire_at) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		t.Name, t.Status, t.MaxUsers, t.MaxApps, r.nullableTime(t.ExpireAt))
	if err != nil {
		return fmt.Errorf("create tenant: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	t.ID = id
	return nil
}

// GetByID retrieves a tenant by ID.
func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*entity.Tenant, error) {
	query := `SELECT id, name, status, max_users, max_apps, expire_at, created_at, updated_at FROM tenants WHERE id = ?`
	var t entity.Tenant
	var expireAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.Status, &t.MaxUsers, &t.MaxApps,
		&expireAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get tenant by id: %w", err)
	}
	if expireAt.Valid {
		t.ExpireAt = expireAt.Time
	}
	return &t, nil
}

// List retrieves tenants with pagination.
func (r *TenantRepository) List(ctx context.Context, page, pageSize int) ([]entity.Tenant, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenants`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count tenants: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, status, max_users, max_apps, expire_at, created_at, updated_at FROM tenants ORDER BY id DESC LIMIT ? OFFSET ?`,
		pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list tenants: %w", err)
	}
	defer rows.Close()

	var tenants []entity.Tenant
	for rows.Next() {
		var t entity.Tenant
		var expireAt sql.NullTime
		if err := rows.Scan(&t.ID, &t.Name, &t.Status, &t.MaxUsers, &t.MaxApps, &expireAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan tenant: %w", err)
		}
		if expireAt.Valid {
			t.ExpireAt = expireAt.Time
		}
		tenants = append(tenants, t)
	}

	return tenants, total, nil
}

// Update updates an existing tenant.
func (r *TenantRepository) Update(ctx context.Context, t *entity.Tenant) error {
	query := `UPDATE tenants SET name=?, status=?, max_users=?, max_apps=?, expire_at=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query,
		t.Name, t.Status, t.MaxUsers, t.MaxApps, r.nullableTime(t.ExpireAt), t.ID)
	if err != nil {
		return fmt.Errorf("update tenant: %w", err)
	}
	return nil
}

// Delete deletes a tenant by ID.
func (r *TenantRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tenants WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete tenant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tenant not found: id=%d", id)
	}
	return nil
}

// UpdateStatus updates a tenant's status.
func (r *TenantRepository) UpdateStatus(ctx context.Context, id int64, status int8) error {
	result, err := r.db.ExecContext(ctx, `UPDATE tenants SET status=? WHERE id=?`, status, id)
	if err != nil {
		return fmt.Errorf("update tenant status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tenant not found: id=%d", id)
	}
	return nil
}

// CheckNameExists checks if a tenant name already exists.
func (r *TenantRepository) CheckNameExists(ctx context.Context, name string, excludeID int64) (bool, error) {
	var count int64
	if excludeID > 0 {
		if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenants WHERE name = ? AND id != ?`, name, excludeID).Scan(&count); err != nil {
			return false, fmt.Errorf("check name exists: %w", err)
		}
	} else {
		if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenants WHERE name = ?`, name).Scan(&count); err != nil {
			return false, fmt.Errorf("check name exists: %w", err)
		}
	}
	return count > 0, nil
}

// nullableTime returns nil if t is zero, otherwise returns t.
func (r *TenantRepository) nullableTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}
