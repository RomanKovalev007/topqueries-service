package service

import (
	"testing"
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
	"github.com/google/uuid"
)

type mockWindow struct {
	added []string
}

func (m *mockWindow) Add(q string)                    { m.added = append(m.added, q) }
func (m *mockWindow) GetTopN(n int) []domain.TopEntry { return nil }

type mockStopList struct {
	contains bool
}

func (m *mockStopList) Contains(q []string) bool { return m.contains }
func (m *mockStopList) AddWords(w []string)       {}
func (m *mockStopList) RemoveWords(w []string)    {}

type mockRateLimiter struct {
	allow bool
}

func (m *mockRateLimiter) Allow(id string) bool { return m.allow }

func freshEvent() domain.SearchEvent {
	return domain.SearchEvent{
		QueryText:   "кроссовки nike",
		UserID:      uuid.New(),
		TimeRequest: time.Now(),
	}
}

func TestAdd_Completed(t *testing.T) {
	w := &mockWindow{}
	svc := NewService(w, &mockStopList{contains: false}, &mockRateLimiter{allow: true}, 5*time.Minute)

	svc.Add(freshEvent())

	if len(w.added) != 1 {
		t.Fatalf("expected 1 added query, got %d", len(w.added))
	}
}

func TestAdd_OutdatedEvent(t *testing.T) {
	w := &mockWindow{}
	svc := NewService(w, &mockStopList{contains: false}, &mockRateLimiter{allow: true}, 5*time.Minute)

	event := freshEvent()
	event.TimeRequest = time.Now().Add(-6 * time.Minute)
	svc.Add(event)

	if len(w.added) != 0 {
		t.Fatal("outdated event should not be added to window")
	}
}

func TestAdd_RateLimited(t *testing.T) {
	w := &mockWindow{}
	svc := NewService(w, &mockStopList{contains: false}, &mockRateLimiter{allow: false}, 5*time.Minute)

	svc.Add(freshEvent())

	if len(w.added) != 0 {
		t.Fatal("rate limited event should not be added to window")
	}
}

func TestAdd_StopList(t *testing.T) {
	w := &mockWindow{}
	svc := NewService(w, &mockStopList{contains: true}, &mockRateLimiter{allow: true}, 5*time.Minute)

	svc.Add(freshEvent())

	if len(w.added) != 0 {
		t.Fatal("stop-listed event should not be added to window")
	}
}

func TestAdd_QueryTrimmed(t *testing.T) {
	w := &mockWindow{}
	svc := NewService(w, &mockStopList{contains: false}, &mockRateLimiter{allow: true}, 5*time.Minute)

	e := freshEvent()
	e.QueryText = "  кроссовки nike  "
	svc.Add(e)

	if len(w.added) != 1 {
		t.Fatalf("expected 1 add, got %d", len(w.added))
	}
	if w.added[0] != "кроссовки nike" {
		t.Fatalf("expected trimmed query, got %q", w.added[0])
	}
}

func TestAdd_WhitespaceOnlyQuery(t *testing.T) {
	w := &mockWindow{}
	svc := NewService(w, &mockStopList{contains: false}, &mockRateLimiter{allow: true}, 5*time.Minute)

	e := freshEvent()
	e.QueryText = "   "
	svc.Add(e)

	if len(w.added) != 0 {
		t.Fatal("whitespace-only query should not be added to window")
	}
}
