package store

import (
	"container/heap"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
)

const (
	interval = 10 * time.Second
	countOfBuckets = 30
	maxN = 100
)


// Min Heap for Calculating Top N Queries

type minHeap []domain.TopEntry

func (h minHeap) Len() int { return len(h) }

func (h minHeap) Less(i, j int) bool { return h[i].Count < h[j].Count}

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


// Bucket

type Bucket struct {
	mtx sync.Mutex
	cnt map[string]int64
}

func NewBucket() *Bucket {
	return &Bucket{
		cnt: make(map[string]int64),
	}
}


// Window

type Window struct {
	cacheTop atomic.Pointer[[]domain.TopEntry]
	totalCounts map[string]int64
	tcMtx sync.Mutex
	buckets [countOfBuckets]*Bucket
	currBuck atomic.Int64
}

func NewWindow() *Window {
	buckets := [countOfBuckets]*Bucket{}
	for i := range buckets{
		buckets[i] = NewBucket()
	}
	return &Window{
		totalCounts: make(map[string]int64),
		buckets: buckets,
	}
}

func (w *Window) Add(query string) {
	w.tcMtx.Lock()
	w.totalCounts[query]++
	w.tcMtx.Unlock()

	buck := w.buckets[w.currBuck.Load()]
	buck.mtx.Lock()
	buck.cnt[query]++
	buck.mtx.Unlock()
}

func (w *Window) Ticker(ctx context.Context) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			next := (w.currBuck.Load() + 1) % int64(len(w.buckets))
			w.tcMtx.Lock()
			w.buckets[next].mtx.Lock()
			for str, count := range w.buckets[next].cnt{
				w.totalCounts[str] -= count
				if w.totalCounts[str] <= 0{
					delete(w.totalCounts, str)
				}
			}
			w.buckets[next].mtx.Unlock()

			snapshot := make(map[string]int64, len(w.totalCounts))
			for k,v := range w.totalCounts {
				snapshot[k] = v
			}

			w.tcMtx.Unlock()

			w.buckets[next] = NewBucket()
			w.currBuck.Store(next)

			w.calculateTopN(snapshot)
		}
	}
}

func (w *Window) GetTopN(n int) []domain.TopEntry {
	top := w.cacheTop.Load()
	if top == nil{
		return []domain.TopEntry{}
	}
	if len(*top) < n {
		return *top
	}

	return (*top)[:n]

}

func (w *Window) calculateTopN(totalCounts map[string]int64) {
	h := make(minHeap, 0, maxN)
	heap.Init(&h)
	for str, count := range totalCounts {
		x := domain.TopEntry{
			Query: str,
			Count: count,
		}
		heap.Push(&h, x)
		if h.Len() > maxN{
			heap.Pop(&h)
		}
	}

	res := make([]domain.TopEntry, h.Len())
	for i := len(res)-1 ; i >= 0; i--{
		res[i] = heap.Pop(&h).(domain.TopEntry)
	}

	w.cacheTop.Store(&res) 
}