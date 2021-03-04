package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

const killsQuery = `SELECT kills.*,
	victim.ID "victim.ID", victim.username "victim.username", victim.charactername "victim.charactername", victim.clantag "victim.clantag", vc.name "victim.clan_info.name", victim.clanID "victim.clanID", victim.joinedClan "victim.joinedClan",
	victim.avatar "victim.avatar", victim.oldgold "victim.oldgold", victim.registered "victim.registered", victim.role "victim.role", victim.tier "victim.tier",
	victim.gold "victim.gold", victim.silver "victim.silver", victim.gold "victim.gold", victim.participation "victim.participation",
	victim.github "victim.github", victim.community "victim.community", victim.mapmaker "victim.mapmaker", victim.moderation "victim.moderation", victim.leaderboardBan "victim.leaderboardBan", victim.statsBan "victim.statsBan",
	killer.ID "killer.ID", killer.username "killer.username", killer.charactername "killer.charactername", killer.clantag "killer.clantag", kc.name "killer.clan_info.name", killer.clanID "killer.clanID", killer.joinedClan "killer.joinedClan",
	killer.avatar "killer.avatar", killer.oldgold "killer.oldgold", killer.registered "killer.registered", killer.role "killer.role", killer.tier "killer.tier",
	killer.gold "killer.gold", killer.silver "killer.silver", killer.gold "killer.gold", killer.participation "killer.participation",
	killer.github "killer.github", killer.community "killer.community", killer.mapmaker "killer.mapmaker", killer.moderation "killer.moderation", killer.leaderboardBan "killer.leaderboardBan", killer.statsBan "killer.statsBan",
	server.ID "server.ID", server.name "server.name" FROM kills
	INNER JOIN players as victim ON kills.victimID=victim.ID
	INNER JOIN players as killer ON kills.killerID=killer.ID
	INNER JOIN servers as server ON kills.serverID=server.ID 
	LEFT JOIN clan_info as vc ON victim.clanID=vc.ID
	LEFT JOIN clan_info as kc ON killer.clanID=kc.ID
`

func killsError(w http.ResponseWriter, err error) {
	log.Printf("Error getting kills: %v\n", err)
	http.Error(w, "Error getting kills", http.StatusInternalServerError)
}

type KillsList struct {
	Limit int           `json:"limit"`
	Start int           `json:"start"`
	Size  int           `json:"size"`
	Next  int           `json:"next"`
	Kills []models.Kill `json:"values"`
}

// GetKills godoc
// @Tags Kills
// @Summary gets kills sorted by the most recent
// @Produce json
// @Param start query int false "Start offset for listing kills"
// @Param limit query int false "Number of response to return"
// @Success 200 {object} KillsList
// @Router /kills [get]
func GetKills(w http.ResponseWriter, r *http.Request) {
	var kills []models.Kill
	var start int64
	var limit int64 = 100

	if v := r.URL.Query().Get("start"); v != "" {
		s, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "Could not parse start", http.StatusBadRequest)
		}

		if s < 0 {
			http.Error(w, "start < 0 is not valid", http.StatusBadRequest)
			return
		}
		start = int64(s)
	}

	err := db.Select(&kills, `SELECT kills.*,
		victim.ID "victim.ID", victim.username "victim.username", victim.charactername "victim.charactername", victim.clantag "victim.clantag", victim.clanID "victim.clanID", vc.name "victim.clan_info.name",
		victim.avatar "victim.avatar", victim.oldgold "victim.oldgold", victim.registered "victim.registered", victim.role "victim.role", victim.tier "victim.tier",
		victim.gold "victim.gold", victim.silver "victim.silver", victim.gold "victim.gold", victim.participation "victim.participation",
		victim.github "victim.github", victim.community "victim.community", victim.mapmaker "victim.mapmaker", victim.moderation "victim.moderation", victim.leaderboardBan "victim.leaderboardBan", victim.statsBan "victim.statsBan",
		killer.ID "killer.ID", killer.username "killer.username", killer.charactername "killer.charactername", killer.clantag "killer.clantag", killer.clanID "killer.clanID", kc.name "killer.clan_info.name",
		killer.avatar "killer.avatar", killer.oldgold "killer.oldgold", killer.registered "killer.registered", killer.role "killer.role", killer.tier "killer.tier",
		killer.gold "killer.gold", killer.silver "killer.silver", killer.gold "killer.gold", killer.participation "killer.participation",
		killer.github "killer.github", killer.community "killer.community", killer.mapmaker "killer.mapmaker", killer.moderation "killer.moderation", killer.leaderboardBan "killer.leaderboardBan", killer.statsBan "killer.statsBan",
		server.ID "server.ID", server.name "server.name"
		FROM (SELECT * FROM kills ORDER by ID DESC LIMIT ?,?) AS kills
		INNER JOIN players AS victim ON kills.victimID=victim.ID
		INNER JOIN players AS killer ON kills.killerID=killer.ID
		INNER JOIN servers AS server ON kills.serverID=server.ID
		LEFT JOIN clan_info as vc ON victim.clanID=vc.ID
		LEFT JOIN clan_info as kc ON killer.clanID=kc.ID
		WHERE NOT victim.statsBan AND NOT killer.statsBan`, start, limit)
	if err != nil {
		killsError(w, err)
		return
	}

	next := int(start) + len(kills)
	if len(kills) < int(limit) {
		next = -1
	}

	JSONResponse(w, KillsList{
		Limit: int(limit),
		Start: int(start),
		Size:  len(kills),
		Next:  next,
		Kills: kills,
	})
}

// GetPlayerKills godoc
// @Tags Players, Kills
// @Summary gets kills for player sorted by most recent
// @Produce json
// @Param id path int true "Player ID"
// @Param start query int false "start offset for returning kills"
// @Param limit query int false "number of responses to return"
// @Success 200 {object} KillsList
// @Router /players/{id}/kills [get]
func GetPlayerKills(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	start, err := GetURLParam("start", 0, r)
	if err != nil {
		http.Error(w, "could not parse start", http.StatusBadRequest)
		return
	}

	if start < 0 {
		http.Error(w, "start must be >= 0", http.StatusBadRequest)
		return
	}

	limit, err := GetURLParam("limit", 50, r)
	if err != nil {
		http.Error(w, "could not parse limit", http.StatusBadRequest)
		return
	}

	limit = int(Min(Max(int64(limit), 1), 50))

	_, err = models.GetPlayer(playerID, db)
	if err != nil {
		playerNotFoundError(w, err)
		return
	}

	var showBanned = "NOT victim.statsBan AND NOT killer.statsBan AND"
	if v := r.URL.Query().Get("showbanned"); v == "true" {
		showBanned = ""
	}

	var kills []models.Kill
	err = db.Select(&kills, killsQuery+"WHERE "+showBanned+" (killerID=? OR victimID=?) ORDER BY kills.ID DESC LIMIT ?,?", playerID, playerID, start, limit)
	if err != nil {
		killsError(w, err)
		return
	}

	next := int(start) + len(kills)
	if len(kills) < int(limit) {
		next = -1
	}

	JSONResponse(w, KillsList{
		Limit:  limit,
		Start:  start,
		Size:   len(kills),
		Next:   next,
		Kills:  kills,
	})

}

// GetKill godoc
// @Tags Kills
// @Summary gets a specific kill by id
// @Produce json
// @Param id path int true "Kill ID"
// @Success 200 {object} models.Kill
// @Router /kills/{id} [get]
func GetKill(w http.ResponseWriter, r *http.Request) {
	killID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "Could not get id", http.StatusBadRequest)
		return
	}

	var kill models.Kill
	err = db.Get(&kill, killsQuery+"WHERE NOT victim.statsBan AND NOT killer.statsBan AND kills.ID=?", int64(killID))
	if err != nil {
		killsError(w, err)
		return
	}

	JSONResponse(w, &kill)
}

// GetServerKills godoc
// @Tags Servers, Kills
// @Summary lists the most recent kills on a server
// @Produce json
// @Param id path int true "Server ID"
// @Success 200 {object} []models.Kill
// @Router /servers/{id}/kills [get]
func GetServerKills(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	var kills []models.Kill
	err = db.Select(&kills, killsQuery+"WHERE NOT victim.statsBan AND NOT killer.statsBan AND serverID=? LIMIT 100", serverID)
	if err != nil {
		killsError(w, err)
		return
	}

	JSONResponse(w, &kills)
}
