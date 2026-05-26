package store

import "testing"

func TestStopList_Contains(t *testing.T) {
	sl := NewStopList([]string{"водка", "сигареты"})

	if !sl.Contains([]string{"мягкая", "водка"}) {
		t.Fatal("expected true for query containing stop word")
	}
	if sl.Contains([]string{"сигарета", "мальбаро"}) {
		t.Fatal("expected false for query without stop words")
	}
}

func TestStopList_AddWords(t *testing.T) {
	sl := NewStopList([]string{})
	sl.AddWords([]string{"казино"})

	if !sl.Contains([]string{"казино"}) {
		t.Fatal("expected added word to be detected")
	}
}

func TestStopList_RemoveWords(t *testing.T) {
	sl := NewStopList([]string{"казино"})
	sl.RemoveWords([]string{"казино"})

	if sl.Contains([]string{"казино"}) {
		t.Fatal("expected removed word to not be detected")
	}
}

func TestStopList_Contains_EmptyQuery(t *testing.T) {
	sl := NewStopList([]string{"водка"})

	if sl.Contains([]string{}) {
		t.Fatal("expected false for empty query")
	}
}
