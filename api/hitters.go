package main

import (
	"log"
	"net/http"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

type Hitter struct {
	PlayerID int64  `json:"-" db:"playerID"`
	Hitter   int64  `json:"hitter"`
	Kills    int64  `json:"kills"`
	Name     string `json:"name"`
}

func HitterRoutes(r *mux.Router) {
	r.HandleFunc("/players/{id:[0-9]+}/hitters", getHitters).Methods("GET")
}

func getHitters(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "coud not get id", http.StatusBadRequest)
		return
	}

	var player models.Player
	err = db.Get(&player, "SELECT * FROM players WHERE ID=?", playerID)
	if err != nil {
		playerNotFoundError(w, err)
		return
	}

	h := []Hitter{}
	err = db.Select(&h, `SELECT * FROM top_hitters AS hitters WHERE hitters.playerID=? ORDER BY hitters.kills DESC LIMIT 5`, playerID)
	if err != nil {
		log.Printf("Could not find hitters for player: %v\n", err)
		http.Error(w, "Could not find hitters for player", http.StatusInternalServerError)
		return
	}

	for i, hitter := range h {
		h[i].Name = models.HitterName(hitter.Hitter)
	}

	JSONResponse(w, struct {
		MyPlayer models.Player `json:"player"`
		Size     int           `json:"size"`
		Hitters  []Hitter      `json:"hitters"`
	}{
		MyPlayer: player,
		Size:     len(h),
		Hitters:  h,
	})
}
