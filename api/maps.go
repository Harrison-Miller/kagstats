package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type MapBasics struct {
	MapName string  `json:"mapName" db:"mapName"`
	Average float32 `json:"average"`
	Matches int64   `json:"matches"`
}

func MapsRoutes(r *mux.Router) {
	r.HandleFunc("/maps", getMaps).Methods("GET")
}

func getMaps(w http.ResponseWriter, r *http.Request) {
	var m []MapBasics
	err := db.Select(&m, `SELECT mapName, ROUND((AVG(ticks)/30)/60) AS average, COUNT(mapName) AS matches FROM map_stats group by mapName`)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get basic map stats: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, m)
}
