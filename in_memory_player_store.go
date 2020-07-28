package main

type InMemoryPlayerStore struct {
	store  map[string]int
	league []Player
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}, []Player{}}
}

func (i *InMemoryPlayerStore) GETPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) GETLeague() []Player {
	var league []Player
	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}
	return league
}

func (i *InMemoryPlayerStore) recordWin(name string) {
	i.store[name]++
}
