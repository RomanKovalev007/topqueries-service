package store

import (
	"container/heap"
	"context"
	"sync/atomic"
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
)

// Min Heap for Calculating Top N Queries

type minHeap []domain.TopEntry

func (h minHeap) Len() int { return len(h) }

func (h minHeap) Less(i, j int) bool { return h[i].Count < h[j].Count }

func (h minHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x any) {
	*h = append(*h, x.(domain.TopEntry))
}

func (h *minHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

type QueryWindow struct {
	cacheTop atomic.Pointer[[]domain.TopEntry]
	window   *CounterWindow
	maxTopN  int
}

func NewQueryWindow(duration time.Duration, bucketCount int, maxTopN int) *QueryWindow {
	window := NewCounterWindow(duration, bucketCount)
	return &QueryWindow{
		window:  window,
		maxTopN: maxTopN,
	}
}

func (w *QueryWindow) Add(query string) {
	w.window.Add(query)
}

func (w *QueryWindow) Ticker(ctx context.Context) {
	ticker := time.NewTicker(w.window.GetBucketInterval())
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			snapshot := w.window.RotateAndSnapshot()
			w.calculateTopN(snapshot)
		}
	}
}

func (w *QueryWindow) GetTopN(n int) []domain.TopEntry {
	top := w.cacheTop.Load()
	if top == nil {
		return []domain.TopEntry{}
	}
	if len(*top) < n {
		return *top
	}
	return (*top)[:n]
}

func (w *QueryWindow) calculateTopN(totalCounts map[string]int64) {
	h := make(minHeap, 0, w.maxTopN)
	for str, count := range totalCounts {
		heap.Push(&h, domain.TopEntry{Query: str, Count: count})
		if h.Len() > w.maxTopN {
			heap.Pop(&h)
		}
	}

	res := make([]domain.TopEntry, h.Len())
	for i := len(res) - 1; i >= 0; i-- {
		res[i] = heap.Pop(&h).(domain.TopEntry)
	}
	w.cacheTop.Store(&res)
}
