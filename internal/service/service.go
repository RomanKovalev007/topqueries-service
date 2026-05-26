package service

import (
	"log/slog"
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
	window         windowProvider
	stopList       stopList
	rateLimiter    rateLimiter
	windowDuration time.Duration
	log            *slog.Logger
}

func NewService(window windowProvider, stopList stopList, rateLimiter rateLimiter, windowDuration time.Duration, log *slog.Logger) *Service {
	return &Service{
		window:         window,
		stopList:       stopList,
		rateLimiter:    rateLimiter,
		windowDuration: windowDuration,
		log:            log,
	}
}

func (s *Service) Add(searchEvent domain.SearchEvent) {
	if time.Since(searchEvent.TimeRequest) < s.windowDuration {
		if s.rateLimiter.Allow(searchEvent.UserID.String()) {
			searchEvent.QueryText = strings.ToLower(searchEvent.QueryText)
			words := strings.Fields(searchEvent.QueryText)
			if len(words) == 0 {
				metrics.EventsTotal.WithLabelValues("empty_query").Inc()
				return
			}
			if !s.stopList.Contains(words) {
				s.window.Add(strings.Join(words, " "))
				metrics.EventsTotal.WithLabelValues("completed").Inc()
			} else {
				metrics.EventsTotal.WithLabelValues("stoplist").Inc()
			}
		} else {
			metrics.EventsTotal.WithLabelValues("rate_limited").Inc()
			s.log.Warn("user rate limited", "user_id", searchEvent.UserID.String())
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
