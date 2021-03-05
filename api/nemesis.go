package main

import (
	"log"
	"net/http"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

type Nemesis struct {
	PlayerID      int64 `json:"-" db:"playerID"`
	NemesisID     int64 `json:"-" db:"nemesisID"`
	Deaths        int64 `json:"deaths"`
	models.Player `json:"nemesis"`
}

func NemesisRoutes(r *mux.Router) {
	r.HandleFunc("/players/{id:[0-9]+}/nemesis", GetNemesis).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}/bullied", GetBullied).Methods("GET")
}

type NemesisResp struct {
	Deaths int `json:"deaths"`
	Nemesis interface{} `json:"nemesis"`
}

// GetNemesis godoc
// @Tags Detailed Stats
// @Summary gets the nemesis of the given player
// @Description the nemesis of a player is the player who has killed them the most
// @Produce json
// @Param id path int true "Player ID"
// @Success 200 {object} NemesisResp
// @Router /players/{id}/nemesis [get]
func GetNemesis(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	var n Nemesis
	err = db.Get(&n, `SELECT n.*, players.*, clan_info.name "clan_info.name" FROM nemesis AS n 
		INNER JOIN players ON n.nemesisID=players.ID 
		LEFT JOIN clan_info ON players.clanID=clan_info.ID
		WHERE n.playerID=? AND n.deaths >= ? ORDER BY n.deaths DESC LIMIT 1`, playerID, config.API.NemesisGate)
	if err != nil {
		JSONResponse(w, NemesisResp{
			Deaths: 0,
			Nemesis: nil,
		})
		return
	}

	JSONResponse(w, n)
}

type BulliedList struct {
	Bullied []Nemesis `json:"bullied"`
}

// GetBullied godoc
// @Tags Detailed Stats
// @Summary returns a list of players bullied by the given player
// @Description a player is bullied by the given player if the given player appears as their nemesis
// @Produce json
// @Param id path int true "Player ID"
// @Success 200 {object} BulliedList
// @Router /players/{id}/bullied [get]
func GetBullied(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	b := []Nemesis{}
	err = db.Select(&b, `SELECT players.*, clan_info.name "clan_info.name", n1.playerID as playerID, n1.nemesisID as nemesisID, n1.deaths as deaths FROM nemesis as n1
		INNER JOIN (SELECT playerID, MAX(deaths) as deaths FROM nemesis GROUP BY playerID) AS n2 ON n1.playerID = n2.playerID AND n1.deaths = n2.deaths AND n1.deaths >= ? AND n1.nemesisID = ?
		INNER JOIN players ON n1.playerID=players.ID
		LEFT JOIN clan_info ON players.clanID=clan_info.ID
	`, config.API.NemesisGate, playerID)
	if err != nil {
		log.Printf("Could not find bullied players for player: %v\n", err)
		http.Error(w, "Could not find bullied players for player", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, BulliedList{
		Bullied: b,
	})
}
