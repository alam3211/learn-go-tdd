package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{store}
	player := "Alam"

	server.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertHTTPCodes(t, response.Code, http.StatusOK)
	assertBodyResponse(t, response.Body.String(), "3")
}

func TestGETPlayers(t *testing.T) {
	storeScores := StubPlayerStore{
		map[string]int{
			"Alam":  20,
			"Dimas": 10,
		}, nil,
	}
	server := &PlayerServer{&storeScores}

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
	storeScores := StubPlayerStore{
		map[string]int{},
		nil,
	}
	server := &PlayerServer{&storeScores}

	t.Run("it records win when POST", func(t *testing.T) {
		request := newPostScoreRequest("Alam")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertHTTPCodes(t, response.Code, http.StatusAccepted)
		if len(storeScores.winCalls) != 1 {
			t.Errorf("got %d, expected %d", len(storeScores.winCalls), 1)

		}
	})
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
