package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGETPlayers(t *testing.T) {
	storeScores := InMemoryPlayerStore{
		map[string]int{
			"Alam":  20,
			"Dimas": 10,
		},
	}
	server := NewPlayerServer(&storeScores)

	t.Run("returns Alam's score", func(t *testing.T) {
		request := newGetScoreRequest("Alam")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertBodyResponse(t, response.Body.String(), "20")
	})

	t.Run("returns Dimas's score", func(t *testing.T) {
		request := newGetScoreRequest("Dimas")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertBodyResponse(t, response.Body.String(), "10")

	})

	t.Run("returns 404 for missing player", func(t *testing.T) {
		request := newGetScoreRequest("Alpha")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertHTTPCodes(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	storeScores := NewInMemoryPlayerStore()
	server := NewPlayerServer(storeScores)

	t.Run("it records win when POST", func(t *testing.T) {
		player := "Alam"

		request := newPostScoreRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertHTTPCodes(t, response.Code, http.StatusAccepted)
		if len(storeScores.store) != 1 {
			t.Errorf("got %d, expected %d", len(storeScores.store), 1)

		}
	})
}

func TestLeague(t *testing.T) {
	storeScores := NewInMemoryPlayerStore()
	server := NewPlayerServer(storeScores)

	t.Run("it return status OK", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertHTTPCodes(t, response.Code, http.StatusOK)
	},
	)
}

func newPostScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertHTTPCodes(t *testing.T, got int, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("HTTP code got %q, expected %q", got, expected)
	}
}

func assertBodyResponse(t *testing.T, got string, expected string) {
	t.Helper()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}
