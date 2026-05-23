package store

import (
	"time"

	"gorm.io/gorm"
)

type Case struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Severity    string         `json:"severity"`
	Status      string         `json:"status"` // open, investigating, resolved, closed
	Assignee    string         `json:"assignee"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Alerts      []Alert        `gorm:"foreignKey:CaseID" json:"alerts"`
	Evidence    []Evidence     `gorm:"foreignKey:CaseID" json:"evidence"`
}

type Alert struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	CaseID    *uint  `json:"case_id"`
	ExternalID string `json:"external_id"` // Reference to Alerting service ID
	Title      string `json:"title"`
}

type Evidence struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CaseID    uint      `json:"case_id"`
	Type      string    `json:"type"` // ip, hash, url, file, etc
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Case{}, &Alert{}, &Evidence{})
}
