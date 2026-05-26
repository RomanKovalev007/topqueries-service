package store

import "sync"

type StopList struct {
	mtx sync.RWMutex
	words map[string]struct{}
}

func NewStopList(words []string) *StopList {
	wds := make(map[string]struct{}, len(words))
	for _, w := range words {
		wds[w] = struct{}{}
	}
	return &StopList{
		words: wds,
	}
}

func (sl *StopList) Contains(word string) bool {
	sl.mtx.RLock()
	_, ok := sl.words[word]
	sl.mtx.RUnlock()
	return ok
}

func (sl *StopList) AddWords(words []string) {
	sl.mtx.Lock()
	for _, w := range words{
		sl.words[w] = struct{}{}
	}
	sl.mtx.Unlock()
}

func (sl *StopList) RemoveWords(words []string) {
	sl.mtx.Lock()
	for _, w := range words {
		delete(sl.words, w)
	}
	sl.mtx.Unlock()
}
