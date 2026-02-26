package domain

import (
	"time"
)

type Agent struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
}
