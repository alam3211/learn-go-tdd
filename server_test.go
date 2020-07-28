package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
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
		}, nil,
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

	t.Run("it return expected Response /league", func(t *testing.T) {
		expectedLeague := []Player{{"Alam", 20}}
		store := NewInMemoryPlayerStore()
		store.store = map[string]int{"Alam": 20}
		server := NewPlayerServer(store)

		request := newGetLeagueRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := getLeagueFromBodyResponse(t, response.Body)
		assertContentTypes(t, response, "application/json")
		assertHTTPCodes(t, response.Code, http.StatusOK)
		assertLeague(t, got, expectedLeague)
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

func newGetLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func assertHTTPCodes(t *testing.T, got int, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("HTTP code got %q, expected %q", got, expected)
	}
}

func assertContentTypes(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	t.Helper()
	if r.Result().Header.Get("content-type") != expected {
		t.Fatalf("Header content-type got %s, expected 'application/json'", r.Result().Header.Get("content-type"))
	}
}

func assertLeague(t *testing.T, got []Player, expected []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Error comparing JSON response between %v and %v", got, expected)
	}
}

func assertBodyResponse(t *testing.T, got string, expected string) {
	t.Helper()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

func getLeagueFromBodyResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Error decoding JSON as the errors are %s", err)
	}

	return
}
