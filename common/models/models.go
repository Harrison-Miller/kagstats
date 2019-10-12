package models

import (
	"github.com/jmoiron/sqlx"
)

type Player struct {
	ID            int64  `json:"id" db:"ID"`
	Username      string `json:"username"`
	Charactername string `json:"characterName"`
	Clantag       string `json:"clanTag"`
	ServerID      int64

	//Information cached from api.kag2d.com
	OldGold    bool   `json:"oldGold"`
	Registered string `json:"registered"`
	Role       int64  `json:"role"`
	Avatar     string `json:"avatar"`
	Tier       int64  `json:"tier"`

	LastEvent     string `json:"lastEvent" db:"lastEvent"`
	LastEventTime int64  `json:"lastEventTime" db:"lastEventTime"`
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
	Address     string `json:"address"`
	Port        string `json:"port"`
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
	Server `json:"server" db:"server,prefix=server."`
}

type Event struct {
	ID       int64  `json:"id" db:"ID"`
	PlayerID int64  `json:"-" db:"playerID"`
	Type     string `json:"type"`
	Time     int64  `json:"time"`
	ServerID int64  `json:"serverId" db:"serverID"`

	Player `json:"player" db:"player,prefix=player."`
}
