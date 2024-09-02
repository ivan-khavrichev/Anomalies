package domain

import (
	"time"
)

type AnomalyMessage struct {
	SessionId string    `gorm:"sessionID"`
	Frequency float64   `gorm:"frequency"`
	Timestamp time.Time `gorm:"timestamp"`
}
