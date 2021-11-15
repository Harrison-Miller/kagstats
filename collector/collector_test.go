package main

import (
	"fmt"
	"github.com/Harrison-Miller/kagstats/collector/fixtures"
	"github.com/Harrison-Miller/rcon"
	"log"
	"testing"
	"time"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/stretchr/testify/require"
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

func TestCollector(t *testing.T) {
	c := Collector{
		logger: log.Default(),
	}
	db = &fixtures.DatabaseMock{
		UpdatePlayerInfoFunc: func(player *models.Player) error {
			player.ID = 1
			return nil
		},
	}
	players = make(map[string]models.Player)
	updater = NewPlayerInfoUpdater(db)

	t.Run("player join and leave", func(t *testing.T) {
		c.OnPlayerJoined(rcon.Message{
			Text: "",
			Args: map[string]string{
				"object": `{"charactername": "Foo", "username": "Foo"}`,
			},
		}, nil)
		require.Equal(t, c.playerCount, 1)
		require.NotNil(t, players["Foo"])
		require.Equal(t, "Foo", players["Foo"].Username)
		c.OnPlayerLeave(rcon.Message{
			Text: "",
			Args: map[string]string{
				"object": `{"charactername": "Foo", "username": "Foo"}`,
			},
		}, nil)
		require.Equal(t, c.playerCount, 0)
		require.NotNil(t, players["Foo"])
	})

	t.Run("joining with new charactername updates the cache", func(t *testing.T) {
		c.OnPlayerJoined(rcon.Message{
			Text: "",
			Args: map[string]string{
				"object": `{"charactername": "bar", "username": "Foo"}`,
			},
		}, nil)
		require.Equal(t, c.playerCount, 1)
		require.NotNil(t, players["Foo"])
		require.Equal(t, "Foo", players["Foo"].Username)
		require.Equal(t, "bar", players["Foo"].Charactername)
	})

}
