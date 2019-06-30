package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Harrison-Miller/kagstats/models"
	"github.com/gorilla/mux"
)

func getPlayers(w http.ResponseWriter, r *http.Request) {
	var players []models.Player
	var start int64
	var limit int64 = 100

	if v := r.URL.Query().Get("start"); v != "" {
		s, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not parse start: %v", err), http.StatusBadRequest)
		}

		if s < 0 {
			http.Error(w, "start < 0 is not valid", http.StatusBadRequest)
			return
		}
		start = int64(s)
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not parse limit: %v", err), http.StatusBadRequest)
		}
		limit = Min(int64(l), limit)
	}

	err := db.Select(&players, "SELECT * FROM players LIMIT ?,?", start, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	next := int(start)+len(players)
	if len(players) < int(limit) {
		next = -1
	}

	JSONResponse(w, struct {
		Limit int `json:"limit"`
		Start int `json:"start"`
		Size int `json:"size"`
		Next int `json:"next"`
		Players []models.Player `json:"players"`
	}{
		Limit: int(limit),
		Start: int(start),
		Size: len(players),
		Next: next,
		Players: players,
	})
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	var player models.Player
	err = db.Get(&player, "SELECT * FROM players WHERE ID=?", int64(playerID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Player not found: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &player)
}
