package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func NewLeague(f io.ReadSeeker) (League, error) {
	var league []Player
	err := json.NewDecoder(f).Decode(&league)

	if err != nil {
		err = fmt.Errorf("Problem passing %v", err)
	}
	return league, err
}

func (l League) Find(name string) (*Player, error) {

	for i := range l {
		if l[i].Name == name {
			return &l[i], nil
		}
	}
	err := fmt.Errorf("Player not found")
	return nil, err
}
