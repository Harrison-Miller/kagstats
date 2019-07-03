package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Harrison-Miller/kagstats/models"
)

const killsQuery = `SELECT kills.*,
	victim.ID "victim.ID", victim.username "victim.username", victim.charactername "victim.charactername", victim.clantag "victim.clantag",
	killer.ID "killer.ID", killer.username "killer.username", killer.charactername "killer.charactername", killer.clantag "killer.clantag" FROM kills
	INNER JOIN players as victim ON kills.victimID=victim.ID
	INNER JOIN players as killer ON kills.killerID=killer.ID `


func getKills(w http.ResponseWriter, r *http.Request) {
	var kills []models.Kill
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

	err := db.Select(&kills, killsQuery + "ORDER BY kills.ID DESC LIMIT ?,?", start, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	next := int(start)+len(kills)
	if len(kills) < int(limit) {
		next = -1
	}

	JSONResponse(w, struct {
		Limit int `json:"limit"`
		Start int `json:"start"`
		Size int `json:"size"`
		Next int `json:"next"`
		Kills []models.Kill `json:"kills"`
	}{
		Limit: int(limit),
		Start: int(start),
		Size: len(kills),
		Next: next,
		Kills: kills,
	})
}

func getPlayerKills(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, fmt.Sprintf("coud not get id: %v", err), http.StatusBadRequest)
		return
	}

	start, err := GetURLParam("start", 0, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse start: %v", err), http.StatusBadRequest)
		return
	}

	if start < 0 {
		http.Error(w, fmt.Sprintf("start must be >= 0"), http.StatusBadRequest)
		return
	}

	limit, err := GetURLParam("limit", 50, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse start: %v", err), http.StatusBadRequest)
		return
	}

	limit = int(Min(Max(int64(limit), 1), 50))

	player, err := models.GetPlayer(playerID, db)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not find player: %v", err), http.StatusInternalServerError)
		return
	}

	var kills []models.Kill
	err = db.Select(&kills, killsQuery + "WHERE killerID=? OR victimID=? ORDER BY kills.ID DESC LIMIT ?,?", playerID, playerID, start, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting kills: %v", err), http.StatusInternalServerError)
		return
	}

	next := int(start)+len(kills)
	if len(kills) < int(limit) {
		next = -1
	}

	JSONResponse(w, struct{
		Player models.Player `json:"player"`
		Limit int `json:"limit"`
		Start int `json:"start"`
		Size int `json:"size"`
		Next int `json:"next"`
		Kills []models.Kill `json:"kills"`
	}{
		Player: player,
		Limit: limit,
		Start: start,
		Size: len(kills),
		Next: next,
		Kills: kills,
	})

}

func getKill(w http.ResponseWriter, r *http.Request) {
	killID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	var kill models.Kill
	err = db.Get(&kill, killsQuery + "WHERE kills.ID=?", int64(killID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Kill not found: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &kill)
}
