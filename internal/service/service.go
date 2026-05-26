package service

import (
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
)

type windowProvider interface {
	Add(query string)
	GetTopN(n int) []domain.TopEntry
}

type stopList interface {
	AddWords(words []string)
	Contains(word string) bool
	RemoveWords(words []string)
}

type Service struct{
	window windowProvider
	stopList stopList
}

func NewService(window windowProvider, stopList stopList) *Service {
	return &Service{
		window: window,
		stopList: stopList,
	}
}

func (s *Service) Add(searchEvent domain.SearchEvent) {
	if time.Since(searchEvent.TimeRequest) < 5 * time.Minute {
		if !s.stopList.Contains(searchEvent.QueryText){
			s.window.Add(searchEvent.QueryText)
		}
	}
}

func (s *Service) GetTopN(n int) []domain.TopEntry {
	return s.window.GetTopN(n)
}

func (s *Service) AddStopWords(words []string) {
	s.stopList.AddWords(words)
}

func (s *Service) RemoveStopWords(words []string) {
	s.stopList.RemoveWords(words)
}