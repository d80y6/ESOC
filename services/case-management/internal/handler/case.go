package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/omniguard/case-management/internal/store"
	"github.com/omniguard/libs/auth"
	"gorm.io/gorm"
)

type CaseHandler struct {
	db *gorm.DB
}

func NewCaseHandler(db *gorm.DB) *CaseHandler {
	return &CaseHandler{db: db}
}

func (h *CaseHandler) CreateCase(c echo.Context) error {
	tenantID := auth.GetEchoTenantID(c)
	if tenantID == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "missing tenant context"})
	}

	var cs store.Case
	if err := c.Bind(&cs); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	cs.TenantID = tenantID
	if err := h.db.Create(&cs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create case"})
	}

	return c.JSON(http.StatusCreated, cs)
}

func (h *CaseHandler) GetCases(c echo.Context) error {
	tenantID := auth.GetEchoTenantID(c)
	if tenantID == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "missing tenant context"})
	}

	var cases []store.Case
	if err := h.db.Where("tenant_id = ?", tenantID).Preload("Alerts").Preload("Evidence").Find(&cases).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch cases"})
	}
	return c.JSON(http.StatusOK, cases)
}

func (h *CaseHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/cases", h.CreateCase)
	e.GET("/cases", h.GetCases)
}
