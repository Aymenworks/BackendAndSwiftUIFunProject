package entities

import (
	"time"

	"gorm.io/gorm"
)

type Tip struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Name      string         `json:"name"`
	ImagePath string         `json:"image_path"`
}

type Tips []*Tip
