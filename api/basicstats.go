package main

import (
	"fmt"
	"net/http"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

const basicQuery = `SELECT basic_stats.*, p.ID "player.ID", p.username "player.username",
p.charactername "player.charactername", p.clantag "player.clantag", p.oldgold "player.oldgold",
p.registered "player.registered", p.role "player.role", p.avatar "player.avatar", p.tier "player.tier",
p.gold "player.gold", p.silver "player.silver", p.bronze "player.bronze", p.participation "player.participation",
p.github "player.github", p.community "player.community", p.mapmaker "player.mapmaker", p.moderation "player.moderation",
p.leaderboardBan "player.leaderboardBan", p.statsBan "player.statsBan"
 FROM basic_stats 
 INNER JOIN players as p ON basic_stats.playerID=p.ID `

type BasicStats struct {
	models.Player `json:"player" db:"player,prefix=player."`
	PlayerID      int64 `json:"-" db:"playerID"`
	Suicides      int64 `json:"suicides"`
	TeamKills     int64 `json:"teamKills"`
	ArcherKills   int64 `json:"archerKills" db:"archer_kills"`
	ArcherDeaths  int64 `json:"archerDeaths" db:"archer_deaths"`
	BuilderKills  int64 `json:"builderKills" db:"builder_kills"`
	BuilderDeaths int64 `json:"builderDeaths" db:"builder_deaths"`
	KnightKills   int64 `json:"knightKills" db:"knight_kills"`
	KnightDeaths  int64 `json:"knightDeaths" db:"knight_deaths"`
	OtherKills    int64 `json:"otherKills" db:"other_kills"`
	OtherDeaths   int64 `json:"otherDeaths" db:"other_deaths"`
	TotalKills    int64 `json:"totalKills" db:"total_kills"`
	TotalDeaths   int64 `json:"totalDeaths" db:"total_deaths"`
}

func BasicStatsRoutes(r *mux.Router) {
	r.HandleFunc("/players/{id:[0-9]+}/basic", getBasicStats).Methods("GET")
	r.HandleFunc("/players/lookup/{name:.+}", getBasicStatsByName).Methods("GET")
	r.HandleFunc("/leaderboard", getBasicLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/kills", getKillsLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/archer", getArcherLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/builder", getBuilderLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/knight", getKnightLeaderBoard).Methods("GET")
	r.HandleFunc("/status", getStatus).Methods("GET")
}

func getBasicStats(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	var stats BasicStats

	err = db.Get(&stats, basicQuery+"WHERE basic_stats.playerID=?", int64(playerID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Player not found: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &stats)
}

func getBasicStatsByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerName := vars["name"]

	var stats BasicStats

	err := db.Get(&stats, basicQuery+"WHERE LOWER(p.username)=LOWER(?)", playerName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Player not found: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &stats)
}

func getBasicLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.total_kills >= ? AND basic_stats.total_deaths >= ? 
		ORDER BY (basic_stats.total_kills / basic_stats.total_deaths) DESC LIMIT 20`, config.API.KDGate, config.API.KDGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size        int          `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

func getKillsLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.total_kills >= ? AND basic_stats.total_deaths >= ? 
		ORDER BY basic_stats.total_kills DESC LIMIT 20`, config.API.KDGate, config.API.KDGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size        int          `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

func getArcherLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.archer_kills >= ? AND basic_stats.archer_deaths >= ? 
		ORDER BY (basic_stats.archer_kills / basic_stats.archer_deaths) DESC LIMIT 20`, config.API.ArcherGate, config.API.ArcherGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size        int          `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

func getBuilderLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.builder_kills >= ? AND basic_stats.builder_deaths >= ? 
		ORDER BY (basic_stats.builder_kills / basic_stats.builder_deaths) DESC LIMIT 20`, config.API.BuilderGate, config.API.BuilderGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size        int          `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

func getKnightLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.knight_kills >= ? AND basic_stats.knight_deaths >= ? 
		ORDER BY (basic_stats.knight_kills / basic_stats.knight_deaths) DESC LIMIT 20`, config.API.KnightGate, config.API.KnightGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size        int          `json:"size"`
		LeaderBoard []BasicStats `json:"leaderboard"`
	}{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	var status Status

	err := db.Get(&status, `SELECT (SELECT COUNT(id) FROM players) as players, (select ID from kills order by ID DESC limit 1) as kills, (SELECT COUNT(id) FROM servers) as servers`)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching status: %v", err), http.StatusInternalServerError)
		return
	}

	status.Version = version

	JSONResponse(w, &status)
}
