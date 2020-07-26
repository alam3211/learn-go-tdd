package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GETPlayerScore(name string) int
	recordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GETPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) recordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodGet:
		p.showScore(w, player)
	case http.MethodPost:
		p.storeScore(w, player)
	}
}

func (p *PlayerServer) storeScore(w http.ResponseWriter, name string) {
	p.store.recordWin(name)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) showScore(w http.ResponseWriter, name string) {
	score := p.store.GETPlayerScore(name)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}
