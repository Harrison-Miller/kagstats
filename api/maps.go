package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type MapBasics struct {
	MapName string  `json:"mapName" db:"mapName"`
	Average float32 `json:"average"`
	Stddev  float32 `json:"stddev" db:"stddev"`
	Matches int64   `json:"matches"`
	Ballots int64   `json:"ballots"`
	Votes   int64   `json:"votes"`
	Wins    int64   `json:"wins"`
}

func MapsRoutes(r *mux.Router) {
	r.HandleFunc("/maps", getMaps).Methods("GET")
}

func getMaps(w http.ResponseWriter, r *http.Request) {
	var m []MapBasics
	err := db.Select(&m, `SELECT map_stats.mapName, average, stddev, matches, ballots, votes, wins 
		FROM (SELECT mapName, ROUND((AVG(ticks)/30)/60) AS average, ROUND(STDDEV((ticks/30)/60)) AS stddev, COUNT(mapName) AS matches 
		FROM map_stats GROUP BY mapName) as map_stats JOIN map_vote_stats ON map_stats.mapName=map_vote_stats.mapName;
	`)
	if err != nil {
		log.Printf("Could not get map stats: %v\n", err)
		http.Error(w, "could not get basic map stats", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, m)
}
