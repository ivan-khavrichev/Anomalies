package service

import "team/transmitter/internal/domain"

type MessagesRepository interface {
	GetMessages(domain.AnomalyMessage)
}

type Messages struct {
	repo MessagesRepository
}

func NewMessages(repo MessagesRepository) *Messages {
	return &Messages{
		repo: repo,
	}
}

func (m *Messages) GetMessages(msg domain.AnomalyMessage) {
	m.repo.GetMessages(msg)
}
