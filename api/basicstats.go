package main

import (
	"log"
	"net/http"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

const basicQuery = `SELECT basic_stats.*, p.ID "player.ID", p.username "player.username",
p.charactername "player.charactername", p.clantag "player.clantag", p.clanID "player.clanID", p.joinedClan "player.joinedClan", p.oldgold "player.oldgold",
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
	r.HandleFunc("/players/{id:[0-9]+}/basic", GetBasicStats).Methods("GET")
	r.HandleFunc("/players/lookup/{name:.+}", GetBasicStatsByName).Methods("GET")
	r.HandleFunc("/leaderboard", GetBasicLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/kills", GetKillsLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/archer", GetArcherLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/builder", GetBuilderLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/knight", GetKnightLeaderBoard).Methods("GET")
	r.HandleFunc("/status", GetStatus).Methods("GET")
}

func playerNotFoundError(w http.ResponseWriter, err error) {
	log.Printf("Player not found: %v\n", err)
	http.Error(w, "Player not found", http.StatusInternalServerError)
}

// GetBasicStats godoc
// @Tags Basic Stats
// @Summary Gets basics stats for a player given and ID
// @Produce json
// @Param id path integer true "Player ID"
// @Success 200 {object} BasicStats
// @Router /players/{id}/basic [get]
func GetBasicStats(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	var stats BasicStats

	err = db.Get(&stats, basicQuery+"WHERE basic_stats.playerID=?", int64(playerID))
	if err != nil {
		err = getBasicPlayerInfoInstead(&stats.Player, w, playerID)
		if err != nil {
			return
		}
	}

	JSONResponse(w, &stats)
}

// GetBasicStatsByName godoc
// @Tags Basic Stats
// @Summary gets basic stats for a player given their username
// @Produces json
// @Param username path string true "Username"
// @Success 200 {object} BasicStats
// @Router /players/lookup/{name} [get]
func GetBasicStatsByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerName := vars["name"]

	var stats BasicStats

	err := db.Get(&stats, basicQuery+"WHERE LOWER(p.username)=LOWER(?)", playerName)
	if err != nil {
		err = db.Get(&stats.Player, "SELECT * FROM players WHERE LOWER(players.username)=LOWER(?)", playerName)
		if err != nil {
			playerNotFoundError(w, err)
			return
		}
	}

	JSONResponse(w, &stats)
}

func leaderboardError(w http.ResponseWriter, err error) {
	log.Printf("Error fetching leaderboard: %v\n", err)
	http.Error(w, "Error fetching leaderboard", http.StatusInternalServerError)
}

type LeaderboardList struct {
	Size        int          `json:"size"`
	LeaderBoard []BasicStats `json:"leaderboard"`
}

// GetBasicLeaderBoard godoc
// @Tags Leaderboards
// @Summary gets the BasicStats for the top 20 players sorted by total kdr
// @Produce json
// @Success 200 {object} LeaderboardList
// @Router /leaderboard [get]
func GetBasicLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.total_kills >= ? AND basic_stats.total_deaths >= ? 
		ORDER BY (basic_stats.total_kills / basic_stats.total_deaths) DESC LIMIT 20`, config.API.KDGate, config.API.KDGate)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, LeaderboardList{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

// GetKillsLeaderBoard godoc
// @Tags Leaderboards
// @Summary gets the BasicStats for the top 20 players sorted by total kills
// @Produce json
// @Success 200 {object} LeaderboardList
// @Router /leaderboard/kills [get]
func GetKillsLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.total_kills >= ? AND basic_stats.total_deaths >= ? 
		ORDER BY basic_stats.total_kills DESC LIMIT 20`, config.API.KDGate, config.API.KDGate)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, LeaderboardList{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

// GetArcherLeaderBoard godoc
// @Tags Leaderboards
// @Summary gets the BasicStats for the top 20 players sorted by archer kdr
// @Produce json
// @Success 200 {object} LeaderboardList
// @Router /leaderboard/archer [get]
func GetArcherLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.archer_kills >= ? AND basic_stats.archer_deaths >= ? 
		ORDER BY (basic_stats.archer_kills / basic_stats.archer_deaths) DESC LIMIT 20`, config.API.ArcherGate, config.API.ArcherGate)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, LeaderboardList{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

// GetBuilderLeaderBoard godoc
// @Tags Leaderboards
// @Summary gets the BasicStats for the top 20 players sorted by builder kdr
// @Produce json
// @Success 200 {object} LeaderboardList
// @Router /leaderboard/builder [get]
func GetBuilderLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.builder_kills >= ? AND basic_stats.builder_deaths >= ? 
		ORDER BY (basic_stats.builder_kills / basic_stats.builder_deaths) DESC LIMIT 20`, config.API.BuilderGate, config.API.BuilderGate)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, LeaderboardList{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

// GetKnightLeaderBoard godoc
// @Tags Leaderboards
// @Summary gets the BasicStats for the top 20 players sorted by knight kdr
// @Produce json
// @Success 200 {object} LeaderboardList
// @Router /leaderboard/knight [get]
func GetKnightLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []BasicStats

	err := db.Select(&stats, basicQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND basic_stats.knight_kills >= ? AND basic_stats.knight_deaths >= ? 
		ORDER BY (basic_stats.knight_kills / basic_stats.knight_deaths) DESC LIMIT 20`, config.API.KnightGate, config.API.KnightGate)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, LeaderboardList{
		Size:        len(stats),
		LeaderBoard: stats,
	})
}

// GetStatus godoc
// @Tags Misc
// @Summary Show basic info about the KAG Stats site
// @Produce json
// @Success 200 {object} Status
// @Failure 500
// @Router /status [get]
func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status Status

	err := db.Get(&status, `SELECT (SELECT COUNT(id) FROM players) as players, (select ID from kills order by ID DESC limit 1) as kills, (SELECT COUNT(id) FROM servers) as servers`)
	if err != nil {
		log.Printf("Error fetching status: %v\n", err)
		http.Error(w, "Error fetching status", http.StatusInternalServerError)
		return
	}

	status.Version = version

	JSONResponse(w, &status)
}
