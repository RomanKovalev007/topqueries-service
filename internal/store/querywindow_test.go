package store

import "testing"

func TestCalculateTopN_Order(t *testing.T) {
	qw := NewQueryWindow()

	counts := map[string]int64{
		"nike":    100,
		"adidas":  50,
		"puma":    200,
	}
	qw.calculateTopN(counts)

	top := qw.GetTopN(3)
	if len(top) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(top))
	}
	if top[0].Query != "puma" || top[0].Count != 200 {
		t.Fatalf("expected puma first, got %+v", top[0])
	}
	if top[1].Query != "nike" || top[1].Count != 100 {
		t.Fatalf("expected nike second, got %+v", top[1])
	}
}

func TestCalculateTopN_LimitN(t *testing.T) {
	qw := NewQueryWindow()

	counts := map[string]int64{"nike": 3, "adidas": 2, "puma": 1}
	qw.calculateTopN(counts)

	top := qw.GetTopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
}

func TestCalculateTopN_MaxN(t *testing.T) {
	qw := NewQueryWindow()

	counts := make(map[string]int64, maxN+2)
	for i := range maxN + 2 {
		counts[string(rune('a'+i%26))+string(rune('0'+i/26))] = int64(i)
	}
	qw.calculateTopN(counts)

	top := qw.GetTopN(maxN + 2)
	if len(top) != maxN {
		t.Fatalf("expected top capped at %d, got %d", maxN, len(top))
	}
}

func TestCalculateTopN_LessThanRequested(t *testing.T) {
	qw := NewQueryWindow()
	qw.calculateTopN(map[string]int64{"nike": 1})

	top := qw.GetTopN(10)
	if len(top) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(top))
	}
}

func TestCalculateTopN_Empty(t *testing.T) {
	qw := NewQueryWindow()
	qw.calculateTopN(map[string]int64{})

	top := qw.GetTopN(10)
	if len(top) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(top))
	}
}

func TestGetTopN_BeforeCalculate(t *testing.T) {
	qw := NewQueryWindow()
	top := qw.GetTopN(10)

	if top == nil {
		t.Fatal("expected non-nil slice before first calculate")
	}
	_ = top
}
