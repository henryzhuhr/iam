//go:build integration

package integration

import (
	"database/sql"
	"fmt"
	"testing"

	"iam/internal/dto/tenant"
	"iam/internal/repository"
	tenantsvc "iam/internal/service/tenant"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"root", "rootpassword", "mysql", 3306, "iam")

	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)

	// Clean tenants table
	_, err = db.Exec("DELETE FROM tenants")
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM tenants")
		db.Close()
	})

	return db
}

func setupService(t *testing.T) *tenantsvc.Service {
	t.Helper()
	db := setupTestDB(t)
	repo := repository.NewTenantRepository(db)
	return tenantsvc.NewService(repo)
}

func TestTenant_CreateAndGet(t *testing.T) {
	svc := setupService(t)

	// Create tenant
	createReq := tenant.CreateTenantRequest{
		Name:     "Test Tenant",
		MaxUsers: 100,
		MaxApps:  5,
	}
	t2, code, msg, err := svc.Create(t.Context(), createReq)
	require.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Empty(t, msg)
	assert.Equal(t, "Test Tenant", t2.Name)
	assert.Equal(t, int64(100), int64(t2.MaxUsers))
	assert.True(t, t2.ID > 0)

	// Get tenant by ID
	got, code, msg, err := svc.GetByID(t.Context(), t2.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Equal(t, "Test Tenant", got.Name)
	assert.Equal(t, t2.ID, got.ID)
}

func TestTenant_CreateDuplicateName(t *testing.T) {
	svc := setupService(t)

	req := tenant.CreateTenantRequest{Name: "Dup Tenant", MaxUsers: 10, MaxApps: 1}

	// First create
	_, code, _, err := svc.Create(t.Context(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, code)

	// Second create (duplicate name)
	_, code, msg, err := svc.Create(t.Context(), req)
	assert.Error(t, err)
	assert.Equal(t, 409, code)
	assert.Contains(t, msg, "already exists")
}

func TestTenant_GetNotFound(t *testing.T) {
	svc := setupService(t)

	_, code, msg, err := svc.GetByID(t.Context(), 999999)
	assert.Error(t, err)
	assert.Equal(t, 404, code)
	assert.Contains(t, msg, "not found")
}

func TestTenant_UpdateAndDelete(t *testing.T) {
	svc := setupService(t)

	// Create
	createReq := tenant.CreateTenantRequest{Name: "Update Test", MaxUsers: 50, MaxApps: 3}
	t2, code, _, err := svc.Create(t.Context(), createReq)
	require.NoError(t, err)
	assert.Equal(t, 0, code)

	// Update
	maxUsers := 200
	updateReq := tenant.UpdateTenantRequest{Name: "Updated Name", MaxUsers: &maxUsers}
	updated, code, _, err := svc.Update(t.Context(), t2.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, 200, updated.MaxUsers)

	// Update status
	code, _, err = svc.UpdateStatus(t.Context(), t2.ID, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, code)

	// Delete
	code, _, err = svc.Delete(t.Context(), t2.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, code)

	// Verify deleted
	_, code, _, err = svc.GetByID(t.Context(), t2.ID)
	assert.Error(t, err)
	assert.Equal(t, 404, code)
}

func TestTenant_List(t *testing.T) {
	svc := setupService(t)

	// Create multiple tenants
	for i := 0; i < 3; i++ {
		_, _, _, err := svc.Create(t.Context(), tenant.CreateTenantRequest{
			Name:     fmt.Sprintf("List Tenant %d", i),
			MaxUsers: 10,
			MaxApps:  1,
		})
		require.NoError(t, err)
	}

	items, total, err := svc.List(t.Context(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, items, 3)
}
