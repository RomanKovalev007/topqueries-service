package service

import (
	"strings"
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
	"github.com/RomanKovalev007/topqueries-service/internal/metrics"
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

type rateLimiter interface {
	Allow(userID string) bool
}

type Service struct {
	window      windowProvider
	stopList    stopList
	rateLimiter rateLimiter
}

func NewService(window windowProvider, stopList stopList, rateLimiter rateLimiter) *Service {
	return &Service{
		window:      window,
		stopList:    stopList,
		rateLimiter: rateLimiter,
	}
}

func (s *Service) Add(searchEvent domain.SearchEvent) {
	if time.Since(searchEvent.TimeRequest) < 5*time.Minute {
		if s.rateLimiter.Allow(searchEvent.UserID.String()) {
			words := strings.Fields(searchEvent.QueryText)
			if !s.stopList.Contains(words) {
				s.window.Add(strings.TrimSpace(searchEvent.QueryText))
				metrics.EventsTotal.WithLabelValues("complete").Inc()
			} else {
				metrics.EventsTotal.WithLabelValues("stoplist").Inc()
			}
		} else {
			metrics.EventsTotal.WithLabelValues("rate_limited").Inc()
		}
	} else {
		metrics.EventsTotal.WithLabelValues("outdated").Inc()
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
