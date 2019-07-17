package models

import (
	"github.com/jmoiron/sqlx"
)

type Player struct {
	ID            int64  `json:"ID" db:"ID"`
	Username      string `json:"username"`
	Charactername string `json:"charactername"`
	Clantag       string `json:"clantag"`
}

func GetPlayer(playerID int, db *sqlx.DB) (Player, error) {
	var player Player
	err := db.Get(&player, "SELECT * FROM players WHERE ID=?", playerID)
	return player, err
}

type Server struct {
	ID   int64  `json:"ID" db:"ID"`
	Name string `json:"name"`
	Tags string `json:"tags"`
}

type Killer = Player

type Kill struct {
	ID          int64  `json:"ID" db:"ID"`
	KillerID    int64  `json:"-" db:"killerID"`
	VictimID    int64  `json:"-" db:"victimID"`
	KillerClass string `json:"killerClass" db:"killerClass"`
	VictimClass string `json:"victimClass" db:"victimClass"`
	Hitter      int64  `json:"hitter" db:"hitter"`
	Time        int64  `json:"time" db:"epoch"`
	ServerID    int64  `json:"serverID" db:"serverID"`
	TeamKill    bool   `json:"teamKill" db:"teamKill"`

	Player `json:"victim" db:"victim,prefix=victim."`
	Killer `json:"killer" db:"killer,prefix=killer."`
}
