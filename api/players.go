package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Harrison-Miller/kagstats/common/utils"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

const playersQuery = `SELECT * FROM players `

func getBasicPlayerInfoInstead(player *models.Player, w http.ResponseWriter, playerID int) error {
	err := db.Get(player, "SELECT * FROM players WHERE players.ID=?", int64(playerID))
	playerNotFoundError(w, err)
	return err
}

func playersError(w http.ResponseWriter, err error) {
	log.Printf("Could not get players: %v\n", err)
	http.Error(w, "Could not get players", http.StatusInternalServerError)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	var players []models.Player
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

	if v := r.URL.Query().Get("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "Could not parse limit", http.StatusBadRequest)
		}
		limit = Min(int64(l), limit)
	}


	err := db.Select(&players, playersQuery+"WHERE NOT players.statsBan ORDER BY players.ID DESC LIMIT ?,?", start, limit)
	if err != nil {
		playersError(w, err)
		return
	}

	next := int(start) + len(players)
	if len(players) < int(limit) {
		next = -1
	}

	JSONResponse(w, struct {
		Limit   int             `json:"limit"`
		Start   int             `json:"start"`
		Size    int             `json:"size"`
		Next    int             `json:"next"`
		Players []models.Player `json:"values"`
	}{
		Limit:   int(limit),
		Start:   int(start),
		Size:    len(players),
		Next:    next,
		Players: players,
	})
}

func searchPlayers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	search := vars["search"]
	search = "%" + search + "%"

	var showBanned = "NOT players.StatsBan AND"
	if v := r.URL.Query().Get("showbanned"); v == "true" {
		showBanned = ""
	}

	var players []models.Player
	err := db.Select(&players, playersQuery+"WHERE "+showBanned+" (lower(username) LIKE ? OR lower(charactername) LIKE ? OR lower(clantag) LIKE ?) LIMIT 100", search, search, search)
	if err != nil {
		playersError(w, err)
		return
	}

	JSONResponse(w, &players)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Could not parse id", http.StatusBadRequest)
		return
	}

	var player models.Player
	err = db.Get(&player, playersQuery+"WHERE players.ID=?", int64(playerID))
	if err != nil {
		playerNotFoundError(w, err)
		return
	}

	JSONResponse(w, &player)
}

func refreshPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	var player models.Player
	err = db.Get(&player, playersQuery+"WHERE players.ID=?", int64(playerID))
	if err != nil {
		playerNotFoundError(w, err)
		return
	}

	utils.GetPlayerAvatar(&player)
	utils.GetPlayerTier(&player)
	utils.GetPlayerInfo(&player)

	_, err = db.Exec("UPDATE players SET oldgold=?,registered=?,role=?,avatar=?,tier=? WHERE ID=?",
		player.OldGold, player.Registered, player.Role, player.Avatar, player.Tier, player.ID)
	if err != nil {
		playerNotFoundError(w, err)
		return
	}

	JSONResponse(w, &struct {
		Message string `json:"message"`
	}{
		Message: "success",
	})
}

type Captures struct {
	PlayerID int64 `json:"playerID" db:"-"`
	Captures int64 `json:"captures" db:"captures"`
}

func getCaptures(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "coud not get id", http.StatusBadRequest)
		return
	}

	var c Captures
	c.PlayerID = int64(playerID)

	err = db.Get(&c, "SELECT COUNT(*) as captures FROM flag_captures WHERE playerID=?", playerID)
	if err != nil {
		playerNotFoundError(w, err)
		return
	}

	JSONResponse(w, c)
}
