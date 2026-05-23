package service

import (
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

func (s *Service) Add(query string) {
	s.window.Add(query)
}

func (s *Service) GetTopN(n int) []domain.TopEntry {
	return s.window.GetTopN(n)
}
