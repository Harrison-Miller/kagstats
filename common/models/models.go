package models

import (
	"github.com/jmoiron/sqlx"
)

type Player struct {
	ID            int64  `json:"id" db:"ID"`
	Username      string `json:"username"`
	Charactername string `json:"characterName"`
	Clantag       string `json:"clanTag"`
}

func GetPlayer(playerID int, db *sqlx.DB) (Player, error) {
	var player Player
	err := db.Get(&player, "SELECT * FROM players WHERE ID=?", playerID)
	return player, err
}

type Server struct {
	ID          int64  `json:"id" db:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Gamemode    string `json:"gameMode"`
	Tags        string `json:"tags"`
}

type Killer = Player

type Kill struct {
	ID          int64  `json:"id" db:"ID"`
	KillerID    int64  `json:"-" db:"killerID"`
	VictimID    int64  `json:"-" db:"victimID"`
	KillerClass string `json:"killerClass" db:"killerClass"`
	VictimClass string `json:"victimClass" db:"victimClass"`
	Hitter      int64  `json:"hitter" db:"hitter"`
	Time        int64  `json:"time" db:"epoch"`
	ServerID    int64  `json:"serverId" db:"serverID"`
	TeamKill    bool   `json:"teamKill" db:"teamKill"`

	Player `json:"victim" db:"victim,prefix=victim."`
	Killer `json:"killer" db:"killer,prefix=killer."`
}

type Event struct {
	ID       int64  `json:"id" db:"ID"`
	PlayerID int64  `json:"-" db:"playerID"`
	Type     string `json:"type"`
	Time     int64  `json:"time"`
	ServerID int64  `json:"serverId" db:"serverID"`

	Player `json:"player" db:"player,prefix=player."`
}
