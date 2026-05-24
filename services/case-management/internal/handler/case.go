package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/omniguard/case-management/internal/store"
	"gorm.io/gorm"
)

type CaseHandler struct {
	db *gorm.DB
}

func NewCaseHandler(db *gorm.DB) *CaseHandler {
	return &CaseHandler{db: db}
}

func (h *CaseHandler) CreateCase(c echo.Context) error {
	var cs store.Case
	if err := c.Bind(&cs); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "missing tenant context"})
	}
	cs.TenantID = tenantID

	if err := h.db.Create(&cs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create case"})
	}

	return c.JSON(http.StatusCreated, cs)
}

func (h *CaseHandler) GetCases(c echo.Context) error {
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "missing tenant context"})
	}

	var cases []store.Case
	h.db.Where("tenant_id = ?", tenantID).Preload("Alerts").Preload("Evidence").Find(&cases)
	return c.JSON(http.StatusOK, cases)
}

func (h *CaseHandler) GetCase(c echo.Context) error {
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "missing tenant context"})
	}

	id := c.Param("id")
	var cs store.Case
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tenantID).Preload("Alerts").Preload("Evidence").First(&cs).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "case not found"})
	}

	return c.JSON(http.StatusOK, cs)
}

func (h *CaseHandler) UpdateCase(c echo.Context) error {
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "missing tenant context"})
	}

	id := c.Param("id")
	var cs store.Case
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cs).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "case not found"})
	}

	if err := c.Bind(&cs); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	h.db.Save(&cs)
	return c.JSON(http.StatusOK, cs)
}

func (h *CaseHandler) DeleteCase(c echo.Context) error {
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "missing tenant context"})
	}

	id := c.Param("id")
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&store.Case{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete case"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CaseHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/cases", h.CreateCase)
	e.GET("/cases", h.GetCases)
	e.GET("/cases/:id", h.GetCase)
	e.PUT("/cases/:id", h.UpdateCase)
	e.DELETE("/cases/:id", h.DeleteCase)
}
