package main

type InMemoryPlayerStore struct {
	store map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

func (i *InMemoryPlayerStore) GETPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) recordWin(name string) {
	i.store[name]++
}
