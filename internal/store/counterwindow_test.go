package store

import (
	"testing"
	"time"
)

func TestCounterWindow_AddAndGetCount(t *testing.T) {
	w := NewCounterWindow(5*time.Minute, 30)

	w.Add("nike")
	w.Add("nike")
	w.Add("adidas")

	if got := w.GetCount("nike"); got != 2 {
		t.Fatalf("expected count 2 for nike, got %d", got)
	}
	if got := w.GetCount("adidas"); got != 1 {
		t.Fatalf("expected count 1 for adidas, got %d", got)
	}
	if got := w.GetCount("puma"); got != 0 {
		t.Fatalf("expected count 0 for unknown key, got %d", got)
	}
}

func TestCounterWindow_RotateClearsOldestBucket(t *testing.T) {
	w := NewCounterWindow(5*time.Minute, 2)

	w.Add("nike")

	// первый rotate выдвигает следующий бакет, второй rotate должен вычесть бакет с "nike"
	w.Rotate()
	w.Rotate()

	if got := w.GetCount("nike"); got != 0 {
		t.Fatalf("expected count 0 after full rotation, got %d", got)
	}
}
