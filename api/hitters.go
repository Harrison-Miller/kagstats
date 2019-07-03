package main


import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/Harrison-Miller/kagstats/models"
)

type Hitter struct {
  PlayerID int64 `json:"-" db:"playerID"`
  Hitter int64 `json:"hitter"`
  Kills int64 `json:"kills"`
}

func HitterRoutes(r *mux.Router) {
  r.HandleFunc("/hitters/{id:[0-9]+}", getHitters).Methods("GET")
}

func getHitters(w http.ResponseWriter, r *http.Request) {
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

  var h []Hitter
  err = db.Select(&h, `SELECT * FROM top_hitters AS hitters WHERE hitters.playerID=? ORDER BY hitters.kills DESC LIMIT 5`, playerID)
  if err != nil {
    http.Error(w, fmt.Sprintf("could not find hitters for player: %v", err), http.StatusInternalServerError)
    return
  }

  JSONResponse(w, struct{
      MyPlayer models.Player `json:"player"`
      Size int `json:"size"`
      Hitters []Hitter `json:"hitters"`
    }{
      MyPlayer: player,
      Size: len(h),
      Hitters: h,
    })
}
