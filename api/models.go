package main

type Nemesis struct {
	PlayerID  int64 `json:"playerID" db:"playerID"`
	NemesisID int64 `json:"nemesisID" db:"nemesisID"`
	Deaths    int64 `json:"deaths"`
}

type Hitters struct {
	PlayerID int64 `json:"playerID" db:"playerID"`
	Hitter   int64 `json:"hitter"`
	Kills    int64 `json:"kills"`
}
