package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/Harrison-Miller/kagstats/common/models"
)

func TestIsNewPlayer(t *testing.T) {
	player := models.Player{}
	player.Registered = "2011-11-09 05:20:22"

	ret := isNewPlayer(&player)

	if ret {
		fmt.Errorf("Expected false")
	}

	newPlayerRegistered := time.Now().Add(-6 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	newPlayer := models.Player{}
	newPlayer.Registered = newPlayerRegistered

	ret2 := isNewPlayer(&newPlayer)
	if !ret2 {
		fmt.Errorf("Expected true")
	}
}
