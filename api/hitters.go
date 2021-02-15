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
	r.HandleFunc("/players/{id:[0-9]+}/hitters", GetHitters).Methods("GET")
}

type HittersList struct {
	MyPlayer models.Player `json:"player"`
	Size     int           `json:"size"`
	Hitters  []Hitter      `json:"hitters"`
}

// GetHitters godoc
// @Tags Detailed Stats
// @Summary Returns the top five weapons (hitters) used by the player
// @Produce json
// @Param id path int true "PlayerID"
// @Success 200 {object} HittersList
// @Router /players/{id}/hitters [get]
func GetHitters(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
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

	JSONResponse(w, HittersList{
		MyPlayer: player,
		Size:     len(h),
		Hitters:  h,
	})
}
