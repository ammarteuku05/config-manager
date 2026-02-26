package domain

import (
	"time"
)

type GlobalConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Config    string `gorm:"type:text"` // JSON format
	Version   string
	CreatedAt time.Time
}
