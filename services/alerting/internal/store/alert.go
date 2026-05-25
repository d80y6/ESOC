package store

import (
	"time"

	"gorm.io/gorm"
)

type Alert struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"index" json:"tenant_id"`
	RuleID      string         `json:"rule_id"`
	RuleName    string         `json:"rule_name"`
	Severity    string         `json:"severity"`
	Status      string         `json:"status"` // open, in-progress, closed
	Description string         `json:"description"`
	EventData   string         `gorm:"type:text" json:"event_data"` // JSON string
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Alert{})
}
