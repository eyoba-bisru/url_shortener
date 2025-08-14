package models

import "gorm.io/gorm"

type URL struct {
	gorm.Model

	OriginalURL string `gorm:"uniqueIndex;not null" json:"original_url"`

	ShortCode string `gorm:"unique;not null" json:"short_code"`

	VisitCount uint `gorm:"default:0" json:"visit_count"`
}
