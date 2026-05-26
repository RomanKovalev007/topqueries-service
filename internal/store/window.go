package store

import (
	"maps"
	"sync"
	"sync/atomic"
	"time"
)

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


// CounterWindow

type CounterWindow struct {
	totalCounts map[string]int64
	tcMtx sync.Mutex
	buckets []*Bucket
	currBuck atomic.Int64
	bucketInterval time.Duration
}

func NewCounterWindow(timeForWindow time.Duration, countOfBuckets int) *CounterWindow {
	buckets := make([]*Bucket, countOfBuckets)
	for i := range buckets{
		buckets[i] = NewBucket()
	}
	return &CounterWindow{
		totalCounts: make(map[string]int64),
		buckets: buckets,
		bucketInterval: timeForWindow / time.Duration(countOfBuckets),
	}
}

func (w *CounterWindow) GetBucketInterval() time.Duration {
	return w.bucketInterval
}

func (w *CounterWindow) Add(entry string) {
	w.tcMtx.Lock()
	w.totalCounts[entry]++
	w.tcMtx.Unlock()

	buck := w.buckets[w.currBuck.Load()]
	buck.mtx.Lock()
	buck.cnt[entry]++
	buck.mtx.Unlock()
}

func (w *CounterWindow) GetCount(key string) int64 {
	w.tcMtx.Lock()
	defer w.tcMtx.Unlock()
	return w.totalCounts[key]
}


func (w *CounterWindow) Rotate() {
	w.tcMtx.Lock()
	defer w.tcMtx.Unlock()

	w.rotate()
}

func (w *CounterWindow) RotateAndSnapshot() map[string]int64 {
	w.tcMtx.Lock()
	defer w.tcMtx.Unlock()
	
	w.rotate()
	snapshot := make(map[string]int64, len(w.totalCounts))
	maps.Copy(snapshot, w.totalCounts)
	return snapshot
}

func (w *CounterWindow) rotate() {
	next := (w.currBuck.Load() + 1) % int64(len(w.buckets))
	w.buckets[next].mtx.Lock()
	for str, count := range w.buckets[next].cnt {
		w.totalCounts[str] -= count
		if w.totalCounts[str] <= 0 {
			delete(w.totalCounts, str)
		}
	}
	w.buckets[next].mtx.Unlock()
	w.buckets[next] = NewBucket()
	w.currBuck.Store(next)
}