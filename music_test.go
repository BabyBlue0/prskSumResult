package main

import (
	"fmt"
	"testing"
)

func TestGetAllSongTitle(t *testing.T) {
	allSongs, err := getAllSongTitle()
	if err != nil {
		t.Fatal("getAllSongTitle is failed.")
	}

	for _, s := range allSongs {
		fmt.Println(s)
	}
}
