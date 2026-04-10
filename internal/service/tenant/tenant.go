// Package tenant provides business logic for tenant management.
package tenant

import (
	"context"
	"fmt"
	"time"

	"iam/internal/dto/tenant"
	"iam/internal/entity"
	"iam/internal/repository"
)

// Service handles tenant-related business operations.
type Service struct {
	repo *repository.TenantRepository
}

// NewService creates a new tenant Service.
func NewService(repo *repository.TenantRepository) *Service {
	return &Service{repo: repo}
}

// Create creates a new tenant after validating name uniqueness.
func (s *Service) Create(ctx context.Context, req tenant.CreateTenantRequest) (*entity.Tenant, int, string, error) {
	exists, err := s.repo.CheckNameExists(ctx, req.Name, 0)
	if err != nil {
		return nil, 0, "", err
	}
	if exists {
		return nil, 409, "tenant name already exists", fmt.Errorf("tenant name '%s' already exists", req.Name)
	}

	var expireAt time.Time
	if req.ExpireAt != "" {
		expireAt, err = time.Parse("2006-01-02 15:04:05", req.ExpireAt)
		if err != nil {
			return nil, 400, "invalid expire_at format, expected 'YYYY-MM-DD HH:MM:SS'", err
		}
	}

	t := &entity.Tenant{
		Name:     req.Name,
		Status:   entity.TenantStatusActive,
		MaxUsers: req.MaxUsers,
		MaxApps:  req.MaxApps,
		ExpireAt: expireAt,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, 500, "failed to create tenant", err
	}

	return t, 0, "", nil
}

// GetByID retrieves a tenant by ID.
func (s *Service) GetByID(ctx context.Context, id int64) (*entity.Tenant, int, string, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, 500, "failed to get tenant", err
	}
	if t == nil {
		return nil, 404, "tenant not found", fmt.Errorf("tenant id=%d not found", id)
	}
	return t, 0, "", nil
}

// List retrieves tenants with pagination.
func (s *Service) List(ctx context.Context, page, pageSize int) ([]entity.Tenant, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.List(ctx, page, pageSize)
}

// Update updates an existing tenant.
func (s *Service) Update(ctx context.Context, id int64, req tenant.UpdateTenantRequest) (*entity.Tenant, int, string, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, 500, "failed to get tenant", err
	}
	if t == nil {
		return nil, 404, "tenant not found", fmt.Errorf("tenant id=%d not found", id)
	}

	var parseErr error
	if req.Name != "" {
		exists, err := s.repo.CheckNameExists(ctx, req.Name, id)
		if err != nil {
			return nil, 500, "failed to check tenant name", err
		}
		if exists {
			return nil, 409, "tenant name already exists", fmt.Errorf("tenant name '%s' already exists", req.Name)
		}
		t.Name = req.Name
	}
	if req.MaxUsers != nil {
		t.MaxUsers = *req.MaxUsers
	}
	if req.MaxApps != nil {
		t.MaxApps = *req.MaxApps
	}
	if req.ExpireAt != "" {
		t.ExpireAt, parseErr = time.Parse("2006-01-02 15:04:05", req.ExpireAt)
		if parseErr != nil {
			return nil, 400, "invalid expire_at format", parseErr
		}
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, 500, "failed to update tenant", err
	}

	return t, 0, "", nil
}

// Delete deletes a tenant by ID.
func (s *Service) Delete(ctx context.Context, id int64) (int, string, error) {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return 500, "failed to delete tenant", err
	}
	return 0, "", nil
}

// UpdateStatus updates a tenant's status.
func (s *Service) UpdateStatus(ctx context.Context, id int64, status int8) (int, string, error) {
	err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return 500, "failed to update tenant status", err
	}
	return 0, "", nil
}

// ToResponse converts a tenant entity to response DTO.
func ToResponse(t *entity.Tenant) tenant.TenantResponse {
	return tenant.TenantResponse{
		ID:        t.ID,
		Name:      t.Name,
		Status:    t.Status,
		MaxUsers:  t.MaxUsers,
		MaxApps:   t.MaxApps,
		ExpireAt:  t.ExpireAt,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// ToResponseList converts tenant entities to list response DTO.
func ToResponseList(items []entity.Tenant, total int64, page, pageSize int) tenant.TenantListResponse {
	respItems := make([]tenant.TenantResponse, len(items))
	for i, t := range items {
		respItems[i] = ToResponse(&t)
	}
	return tenant.TenantListResponse{
		Items:    respItems,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
