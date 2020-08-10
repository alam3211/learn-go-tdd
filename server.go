package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	GETPlayerScore(name string) int
	recordWin(name string)
	GETLeague() League
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router
	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		leagueTable := p.store.GETLeague()
		json.NewEncoder(w).Encode(leagueTable)
	}
	w.WriteHeader(http.StatusOK)

}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
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
