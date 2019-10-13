package main

import (
	"fmt"
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
	r.HandleFunc("/players/{id:[0-9]+}/nemesis", getNemesis).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}/bullied", getBullied).Methods("GET")
}

func getNemesis(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, fmt.Sprintf("coud not get id: %v", err), http.StatusBadRequest)
		return
	}

	var n Nemesis
	err = db.Get(&n, `SELECT * FROM nemesis AS n INNER JOIN players 
		ON n.nemesisID=players.ID WHERE n.playerID=? AND n.deaths >= ? ORDER BY n.deaths DESC LIMIT 1`, playerID, config.API.NemesisGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not find nemeses for player: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, n)
}

func getBullied(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, fmt.Sprintf("coud not get id: %v", err), http.StatusBadRequest)
		return
	}

	var b []Nemesis
	err = db.Select(&b, `SELECT players.*, n1.playerID as playerID, n1.nemesisID as nemesisID, n1.deaths as deaths 
		FROM nemesis as n1 INNER JOIN (SELECT playerID, MAX(deaths) as deaths FROM nemesis GROUP BY playerID) AS n2 
		ON n1.playerID = n2.playerID AND n1.deaths = n2.deaths AND n1.deaths >= ? AND n1.nemesisID = ? INNER JOIN players ON n1.playerID=players.ID
	`, config.API.NemesisGate, playerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not find bullied players for player: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Bullied []Nemesis `json:"bullied"`
	}{
		Bullied: b,
	})
}
