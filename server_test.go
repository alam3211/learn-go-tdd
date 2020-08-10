package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestFileSystemStore(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
		{"Name":"Alam","Wins":10},
		{"Name":"Dimas","Wins":15},
		{"Name":"Tamtam","Wins":20}
	]`)

	defer cleanDatabase()
	store := NewFileSystemStore(database)

	t.Run("store /league from a reader", func(t *testing.T) {
		got := store.GETLeague()

		expected := []Player{
			{"Alam", 10},
			{"Dimas", 15},
			{"Tamtam", 20},
		}
		assertLeague(t, got, expected)

		got = store.GETLeague()
		assertLeague(t, got, expected)
	})

	t.Run("get player score", func(t *testing.T) {
		got := store.GETPlayerScore("Tamtam")

		expected := 20

		assertScore(t, got, expected)
	})

	t.Run("store wins for existing player", func(t *testing.T) {
		store.recordWin("Alam")

		got := store.GETPlayerScore("Alam")
		expected := 11
		assertScore(t, got, expected)
	})

	t.Run("store wins for new player", func(t *testing.T) {
		store.recordWin("Deka")
		got := store.GETPlayerScore("Deka")
		expected := 1
		assertScore(t, got, expected)
	})
}

func TestGETPlayers(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
		{"Name":"Alam","Wins":20},
		{"Name":"Dimas","Wins":10},
		{"Name":"Tamtam","Wins":20}
	]`)

	defer cleanDatabase()
	store := NewFileSystemStore(database)
	server := NewPlayerServer(store)

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
	database, cleanDatabase := createTempFile(t, "")

	defer cleanDatabase()
	store := NewFileSystemStore(database)
	server := NewPlayerServer(store)

	t.Run("it records win when POST", func(t *testing.T) {
		player := "Alam"

		request := newPostScoreRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertHTTPCodes(t, response.Code, http.StatusAccepted)
		if len(store.GETLeague()) != 1 {
			t.Errorf("got %d, expected %d", len(store.GETLeague()), 1)

		}
	})
}

func TestLeague(t *testing.T) {

	t.Run("it return expected Response /league", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name":"Alam","Wins":20}
		]`)

		defer cleanDatabase()
		store := NewFileSystemStore(database)
		server := NewPlayerServer(store)

		expectedLeague := []Player{{"Alam", 20}}

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

func assertScore(t *testing.T, got int, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("There's an error on %v", err)
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

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
