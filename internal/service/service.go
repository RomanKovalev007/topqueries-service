package service

import (
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
)

type windowProvider interface{
	Add(query string)
	GetTopN(n int) []domain.TopEntry
}

type Service struct{
	window windowProvider
}

func NewService(window windowProvider) *Service {
	return &Service{
		window: window,
	}
}

func (s *Service) Add(searchEvent domain.SearchEvent) {
	if time.Since(searchEvent.TimeRequest) < 5 * time.Minute {
		s.window.Add(searchEvent.QueryText)
	}
}

func (s *Service) GetTopN(n int) []domain.TopEntry {
	return s.window.GetTopN(n)
}
