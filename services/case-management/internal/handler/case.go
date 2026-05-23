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

	if err := h.db.Create(&cs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create case"})
	}

	return c.JSON(http.StatusCreated, cs)
}

func (h *CaseHandler) GetCases(c echo.Context) error {
	var cases []store.Case
	h.db.Preload("Alerts").Preload("Evidence").Find(&cases)
	return c.JSON(http.StatusOK, cases)
}

func (h *CaseHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/cases", h.CreateCase)
	e.GET("/cases", h.GetCases)
}
