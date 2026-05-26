package service

import (
	"slices"
	"strings"
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
)

type windowProvider interface {
	Add(query string)
	GetTopN(n int) []domain.TopEntry
}

type stopList interface {
	AddWords(words []string)
	Contains(query []string) bool
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
		sortSlice := s.queryToSortSlice(searchEvent.QueryText)
		if !s.stopList.Contains(sortSlice){
			s.window.Add(strings.Join(sortSlice, " "))
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


// TODO: sorted key breaks display order (ex: "nike air max" to "air max nike").
func (s *Service) queryToSortSlice(query string) []string {
	sl := strings.Fields(query)
	slices.Sort(sl)
	return sl
}