package models

import "database/sql"

type Player struct {
	ID            int64  `json:"ID" db:"ID"`
	Username      string `json:"username"`
	Charactername string `json:"charactername"`
	Clantag       string `json:"clantag"`
}

type Server struct {
	ID   int64  `json:"ID" db:"ID"`
	Name string `json:"name"`
	Tags string `json:"tags"`
}

type Kill struct {
	ID          int64         `json:"ID" db:"ID"`
	KillerID    sql.NullInt64 `json:"killerID" db:"killerID"`
	VictimID    int64         `json:"victimID" db:"victimID"`
	AssistID    sql.NullInt64 `json:"assistID" db:"assistID"`
	KillerClass string        `json:"killerClass" db:"killerClass"`
	VictimClass string        `json:"victimClass" db:"victimClass"`
	Hitter      int64         `json:"hitter" db:"hitter"`
	Time        int64         `json:"time" db:"epoch"`
	ServerID    int64         `json:"serverID" db:"serverID"`
	TeamKill    bool          `json:"teamKill" db:"teamKill"`
}
