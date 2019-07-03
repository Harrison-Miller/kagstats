package main

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/Harrison-Miller/kagstats/models"
)

type Nemesis struct {
  PlayerID int64 `json:"-" db:"playerID"`
  NemesisID int64 `json:"-" db:"nemesisID"`
  Deaths int64 `json:"deaths"`
  models.Player `json:"nemesis"`
}

func NemesisRoutes(r *mux.Router) {
  r.HandleFunc("/players/{id:[0-9]+}/nemesis", getNemesis).Methods("GET")
}

func getNemesis(w http.ResponseWriter, r *http.Request) {
  playerID, err := GetIntURLArg("id", r)
  if err != nil {
    http.Error(w, fmt.Sprintf("coud not get id: %v", err), http.StatusBadRequest)
    return
  }

  var player models.Player
  err = db.Get(&player, "SELECT * FROM players WHERE ID=?", playerID)
  if err != nil {
    http.Error(w, fmt.Sprintf("could not find player: %v", err), http.StatusInternalServerError)
    return
  }

  var n []Nemesis
  err = db.Select(&n, `SELECT * FROM nemesis AS n INNER JOIN players ON n.nemesisID=players.ID WHERE n.playerID=? ORDER BY n.deaths DESC LIMIT 3`, playerID)
  if err != nil {
    http.Error(w, fmt.Sprintf("could not find nemeses for player: %v", err), http.StatusInternalServerError)
    return
  }

  JSONResponse(w, struct{
      MyPlayer models.Player `json:"player"`
      Size int `json:"size"`
      Nemeses []Nemesis `json:"nemeses"`
    }{
      MyPlayer: player,
      Size: len(n),
      Nemeses: n,
    })
}
