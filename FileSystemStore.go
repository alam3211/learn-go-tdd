package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileSystemStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemStore(database *os.File) (*FileSystemStore, error) {
	database.Seek(0, 0)
	league, err := NewLeague(database)

	if err != nil {
		return nil, fmt.Errorf("error storing league %s,%v", database.Name(), err)
	}

	return &FileSystemStore{
		json.NewEncoder(&Tape{database}), league,
	}, nil
}

func (f *FileSystemStore) GETLeague() League {
	return f.league
}

func (f *FileSystemStore) GETPlayerScore(name string) int {
	player, _ := f.league.Find(name)

	if player != nil {
		return player.Wins
	}
	return 0
}

func (f *FileSystemStore) recordWin(name string) {
	player, _ := f.league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}
	f.database.Encode(f.league)
}
