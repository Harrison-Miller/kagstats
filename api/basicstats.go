package main

import (
	"fmt"
	"net/http"

	"github.com/Harrison-Miller/kagstats/models"
	"github.com/gorilla/mux"
)

type BasicStats struct {
	models.Player `json:"player"`
	PlayerID      int64 `json:"-" db:"playerID"`
	Suicides      int64 `json:"suicides"`
	TeamKills     int64 `json:"teamKills"`
	ArcherKills   int64 `json:"archerKills" db:"archer_kills"`
	ArcherDeaths  int64 `json:"archerDeaths" db:"archer_deaths"`
	BuilderKills  int64 `json:"builderKills" db:"builder_kills"`
	BuilderDeaths int64 `json:"builderDeaths" db:"builder_deaths"`
	KnightKills   int64 `json:"kngihtKills" db:"knight_kills"`
	KnightDeaths  int64 `json:"knightDeaths" db:"knight_deaths"`
	OtherKills    int64 `json:"otherKills" db:"other_kills"`
	OtherDeaths   int64 `json:"otherDeaths" db:"other_deaths"`
}

func BasicStatsRoutes(r *mux.Router) {
	r.HandleFunc("/basic/{id:[0-9]+}", getBasicStats).Methods("GET")
	r.HandleFunc("/basic/leaderboard", getBasicLeaderBoard).Methods("GET")
	r.HandleFunc("/basic/leaderboard/archer", getArcherLeaderBoard).Methods("GET")
	r.HandleFunc("/basic/leaderboard/builder", getBuilderLeaderBoard).Methods("GET")
	r.HandleFunc("/basic/leaderboard/knight", getKnightLeaderBoard).Methods("GET")
}

func getBasicStats(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	var stats BasicStats

	err = db.Get(&stats, "SELECT * FROM basic_stats INNER JOIN players ON basic_stats.playerID=players.ID WHERE basic_stats.playerID=?", int64(playerID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Player not found: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &stats)
}

func getBasicLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, "SELECT * FROM basic_stats INNER JOIN players ON basic_stats.playerID=players.ID ORDER BY ((basic_stats.archer_kills + basic_stats.builder_kills + basic_stats.knight_kills) / (basic_stats.archer_deaths + basic_stats.builder_deaths + basic_stats.knight_deaths)) DESC LIMIT 20")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size int `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size: len(stats),
		LeaderBoard: stats,
	})
}

func getArcherLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, "SELECT * FROM basic_stats INNER JOIN players ON basic_stats.playerID=players.ID ORDER BY (basic_stats.archer_kills / basic_stats.archer_deaths) DESC LIMIT 20")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size int `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size: len(stats),
		LeaderBoard: stats,
	})
}

func getBuilderLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, "SELECT * FROM basic_stats INNER JOIN players ON basic_stats.playerID=players.ID ORDER BY (basic_stats.builder_kills / basic_stats.builder_deaths) DESC LIMIT 20")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size int `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size: len(stats),
		LeaderBoard: stats,
	})
}

func getKnightLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, "SELECT * FROM basic_stats INNER JOIN players ON basic_stats.playerID=players.ID ORDER BY (basic_stats.knight_kills / basic_stats.knight_deaths) DESC LIMIT 20")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size int `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size: len(stats),
		LeaderBoard: stats,
	})
}
