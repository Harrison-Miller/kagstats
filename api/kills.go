package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

const killsQuery = `SELECT kills.*,
	victim.ID "victim.ID", victim.username "victim.username", victim.charactername "victim.charactername", victim.clantag "victim.clantag",
	victim.avatar "victim.avatar", victim.oldgold "victim.oldgold", victim.registered "victim.registered", victim.role "victim.role", victim.tier "victim.tier",
	victim.gold "victim.gold", victim.silver "victim.silver", victim.gold "victim.gold", victim.participation "victim.participation",
	victim.github "victim.github", victim.community "victim.community", victim.mapmaker "victim.mapmaker", victim.moderation "victim.moderation", victim.leaderboardBan "victim.leaderboardBan", victim.statsBan "victim.statsBan",
	killer.ID "killer.ID", killer.username "killer.username", killer.charactername "killer.charactername", killer.clantag "killer.clantag",
	killer.avatar "killer.avatar", killer.oldgold "killer.oldgold", killer.registered "killer.registered", killer.role "killer.role", killer.tier "killer.tier",
	killer.gold "killer.gold", killer.silver "killer.silver", killer.gold "killer.gold", killer.participation "killer.participation",
	killer.github "killer.github", killer.community "killer.community", killer.mapmaker "killer.mapmaker", killer.moderation "killer.moderation", killer.leaderboardBan "killer.leaderboardBan", killer.statsBan "killer.statsBan",
	server.ID "server.ID", server.name "server.name" FROM kills
	INNER JOIN players as victim ON kills.victimID=victim.ID
	INNER JOIN players as killer ON kills.killerID=killer.ID
	INNER JOIN servers as server ON kills.serverID=server.ID `

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

	err := db.Select(&kills, `SELECT kills.*,
		victim.ID "victim.ID", victim.username "victim.username", victim.charactername "victim.charactername", victim.clantag "victim.clantag",
		victim.avatar "victim.avatar", victim.oldgold "victim.oldgold", victim.registered "victim.registered", victim.role "victim.role", victim.tier "victim.tier",
		victim.gold "victim.gold", victim.silver "victim.silver", victim.gold "victim.gold", victim.participation "victim.participation",
		victim.github "victim.github", victim.community "victim.community", victim.mapmaker "victim.mapmaker", victim.moderation "victim.moderation", victim.leaderboardBan "victim.leaderboardBan", victim.statsBan "victim.statsBan",
		killer.ID "killer.ID", killer.username "killer.username", killer.charactername "killer.charactername", killer.clantag "killer.clantag",
		killer.avatar "killer.avatar", killer.oldgold "killer.oldgold", killer.registered "killer.registered", killer.role "killer.role", killer.tier "killer.tier",
		killer.gold "killer.gold", killer.silver "killer.silver", killer.gold "killer.gold", killer.participation "killer.participation",
		killer.github "killer.github", killer.community "killer.community", killer.mapmaker "killer.mapmaker", killer.moderation "killer.moderation", killer.leaderboardBan "killer.leaderboardBan", killer.statsBan "killer.statsBan",
		server.ID "server.ID", server.name "server.name"
		FROM (SELECT * FROM kills ORDER by ID DESC LIMIT ?,?) AS kills
		INNER JOIN players AS victim ON kills.victimID=victim.ID
		INNER JOIN players AS killer ON kills.killerID=killer.ID
		INNER JOIN servers AS server ON kills.serverID=server.ID
		WHERE NOT victim.statsBan AND NOT killer.statsBan`, start, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	next := int(start) + len(kills)
	if len(kills) < int(limit) {
		next = -1
	}

	JSONResponse(w, struct {
		Limit int           `json:"limit"`
		Start int           `json:"start"`
		Size  int           `json:"size"`
		Next  int           `json:"next"`
		Kills []models.Kill `json:"values"`
	}{
		Limit: int(limit),
		Start: int(start),
		Size:  len(kills),
		Next:  next,
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

	var showBanned = "NOT victim.statsBan AND NOT killer.statsBan AND"
	if v := r.URL.Query().Get("showbanned"); v == "true" {
		showBanned = ""
	}

	var kills []models.Kill
	err = db.Select(&kills, killsQuery+"WHERE "+showBanned+" (killerID=? OR victimID=?) ORDER BY kills.ID DESC LIMIT ?,?", playerID, playerID, start, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting kills: %v", err), http.StatusInternalServerError)
		return
	}

	next := int(start) + len(kills)
	if len(kills) < int(limit) {
		next = -1
	}

	JSONResponse(w, struct {
		Player models.Player `json:"player"`
		Limit  int           `json:"limit"`
		Start  int           `json:"start"`
		Size   int           `json:"size"`
		Next   int           `json:"next"`
		Kills  []models.Kill `json:"values"`
	}{
		Player: player,
		Limit:  limit,
		Start:  start,
		Size:   len(kills),
		Next:   next,
		Kills:  kills,
	})

}

func getKill(w http.ResponseWriter, r *http.Request) {
	killID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	var kill models.Kill
	err = db.Get(&kill, killsQuery+"WHERE NOT victim.statsBan AND NOT killer.statsBan AND kills.ID=?", int64(killID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Kill not found: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &kill)
}

func getServerKills(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	var kills []models.Kill
	err = db.Select(&kills, killsQuery+"WHERE NOT victim.statsBan AND NOT player.statsBan AND serverID=? LIMIT 100", serverID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting kills: %s", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &kills)
}
