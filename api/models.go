package main

type Hitters struct {
	PlayerID int64 `json:"playerID" db:"playerID"`
	Hitter   int64 `json:"hitter"`
	Kills    int64 `json:"kills"`
}
