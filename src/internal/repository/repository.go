package repository

import (
	"team/transmitter/internal/domain"

	"gorm.io/gorm"
)

type Messages struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Messages {
	return &Messages{
		db: db,
	}
}

func (m *Messages) GetMessages(msg domain.AnomalyMessage) {
	m.db.Model(&domain.AnomalyMessage{}).Create(&msg)
}
