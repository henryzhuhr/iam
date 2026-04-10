// Package tenant provides HTTP handlers for tenant management.
package tenant

import (
	"encoding/json"
	"net/http"
	"strconv"

	"iam/internal/constant"
	"iam/internal/dto/tenant"
	"iam/internal/middleware"
	tenantsvc "iam/internal/service/tenant"
)

// Handler handles HTTP requests for tenant operations.
type Handler struct {
	svc *tenantsvc.Service
}

// NewHandler creates a new tenant Handler.
func NewHandler(svc *tenantsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// Create handles POST /api/v1/tenants
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req tenant.CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, constant.CodeAuthFailed, "invalid request body")
		return
	}

	t, httpStatus, msg, err := h.svc.Create(r.Context(), req)
	if err != nil {
		middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
		return
	}

	middleware.WriteSuccess(w, middleware.ToResponse(t))
}

// Get handles GET /api/v1/tenants/:id
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, constant.CodeTenantNotFound, "invalid tenant id")
		return
	}

	t, httpStatus, msg, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
		return
	}

	middleware.WriteSuccess(w, middleware.ToResponse(t))
}

// List handles GET /api/v1/tenants
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	items, total, err := h.svc.List(r.Context(), page, pageSize)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, constant.CodeInternalError, "failed to list tenants")
		return
	}

	middleware.WriteSuccess(w, tenantsvc.ToResponseList(items, total, page, pageSize))
}

// Update handles PUT /api/v1/tenants/:id
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, constant.CodeTenantNotFound, "invalid tenant id")
		return
	}

	var req tenant.UpdateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, constant.CodeAuthFailed, "invalid request body")
		return
	}

	t, httpStatus, msg, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
		return
	}

	middleware.WriteSuccess(w, middleware.ToResponse(t))
}

// Delete handles DELETE /api/v1/tenants/:id
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, constant.CodeTenantNotFound, "invalid tenant id")
		return
	}

	httpStatus, msg, err := h.svc.Delete(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
		return
	}

	middleware.WriteSuccess(w, nil)
}

// UpdateStatus handles PUT /api/v1/tenants/:id/status
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, constant.CodeTenantNotFound, "invalid tenant id")
		return
	}

	var req struct {
		Status int8 `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, constant.CodeAuthFailed, "invalid request body")
		return
	}

	httpStatus, msg, err := h.svc.UpdateStatus(r.Context(), id, req.Status)
	if err != nil {
		middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
		return
	}

	middleware.WriteSuccess(w, nil)
}
