package service

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
	"github.com/RomanKovalev007/topqueries-service/internal/store"
	"github.com/google/uuid"
)

func setupService() *Service {
	window := store.NewQueryWindow(5*time.Minute, 30, 100)
	stopList := store.NewStopList([]string{})
	rateLimiter := store.NewRateLimiter(time.Minute, 12, 100)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewService(window, stopList, rateLimiter, 5*time.Minute, log)
}

func BenchmarkAdd(b *testing.B) {
	svc := setupService()
	event := domain.SearchEvent{
		QueryText:   "кроссовки nike",
		UserID:      uuid.New(),
		TimeRequest: time.Now(),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			svc.Add(event)
		}
	})
}

func BenchmarkGetTopN(b *testing.B) {
	svc := setupService()

	queries := map[string]int{
		"кроссовки nike":       1000,
		"платье летнее":        850,
		"набор посуды":         100,
		"зарядка для ноутбука": 240,
		"чехол iphone":         600,
		"магнитофон":           450,
		"зубная щетка":         390,
		"паста арахисовая":     350,
		"шлем гоночный":        330,
		"холодильник":          330,
		"шапка зимняя":         600,
		"повербанк":            120,
	}

	for q, n := range queries {
		for range n {
			svc.Add(domain.SearchEvent{
				QueryText:   q,
				UserID:      uuid.New(),
				TimeRequest: time.Now(),
			})
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			svc.GetTopN(10)
		}
	})
}
