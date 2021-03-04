package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const monthlyQuery = `SELECT monthly_stats.*, p.ID "player.ID", p.username "player.username",
p.charactername "player.charactername", p.clantag "player.clantag", p.oldgold "player.oldgold",
p.registered "player.registered", p.role "player.role", p.avatar "player.avatar", p.tier "player.tier",
p.gold "player.gold", p.silver "player.silver", p.bronze "player.bronze", p.participation "player.participation",
p.github "player.github", p.community "player.community", p.mapmaker "player.mapmaker", p.moderation "player.moderation",
p.leaderboardBan "player.leaderboardBan", p.statsBan "player.statsBan"
 FROM monthly_stats 
INNER JOIN players as p ON monthly_stats.playerID=p.ID `

type MonthlyStats struct {
	Year  int64 `json:"year"`
	Month int64 `json:"month"`
	BasicStats
}

func MonthlyStatsRoutes(r *mux.Router) {
	r.HandleFunc("/leaderboard/monthly/archer", GetMonthlyArcherLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/monthly/builder", GetMonthlyBuilderLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/monthly/knight", GetMonthlyKnightLeaderBoard).Methods("GET")

	r.HandleFunc("/players/{id:[0-9]+}/monthly", GetMonthlyStats).Methods("GET")
}

func getYearMonth(r *http.Request) (int, int) {
	t := time.Now()
	year, _ := GetURLParam("year", t.Year(), r)
	month, _ := GetURLParam("month", int(t.Month()), r)
	return year, month
}

type MonthlyLeaderboardList struct {
	Size        int            `json:"size"`
	Year        int            `json:"year"`
	Month       int            `json:"month"`
	LeaderBoard []MonthlyStats `json:"leaderboard"`
}

// GetMonthlyArcherLeaderBoard godoc
// @Tags Leaderboards
// @Summary returns the top 20 players for a given month sorted by archer kdr
// @Produce json
// @Param year query int false "year - defaults to current year"
// @Param month query int false "month - defaults to current month"
// @Success 200 {object} MonthlyLeaderboardList
// @Router /leaderboard/monthly/archer [get]
func GetMonthlyArcherLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []MonthlyStats

	year, month := getYearMonth(r)

	err := db.Select(&stats, monthlyQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND monthly_stats.year=? AND monthly_stats.month=? AND monthly_stats.archer_kills >= ? AND monthly_stats.archer_deaths >= ? 
		ORDER BY (monthly_stats.archer_kills / monthly_stats.archer_deaths) DESC LIMIT 20`, year, month, config.API.ArcherGate, config.API.ArcherGate)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, MonthlyLeaderboardList{
		Size:        len(stats),
		Year:        year,
		Month:       month,
		LeaderBoard: stats,
	})
}

// GetMonthlyBuilderLeaderBoard godoc
// @Tags Leaderboards
// @Summary returns the top 20 players for a given month sorted by builder kdr
// @Produce json
// @Param year query int false "year - defaults to current year"
// @Param month query int false "month - defaults to current month"
// @Success 200 {object} MonthlyLeaderboardList
// @Router /leaderboard/monthly/builder [get]
func GetMonthlyBuilderLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []MonthlyStats

	year, month := getYearMonth(r)

	err := db.Select(&stats, monthlyQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND monthly_stats.year=? AND monthly_stats.month=? AND monthly_stats.builder_kills >= ? AND monthly_stats.builder_deaths >= ? 
		ORDER BY (monthly_stats.builder_kills / monthly_stats.builder_deaths) DESC LIMIT 20`, year, month, config.API.BuilderGate, config.API.BuilderGate)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, MonthlyLeaderboardList{
		Size:        len(stats),
		Year:        year,
		Month:       month,
		LeaderBoard: stats,
	})
}

// GetMonthlyKnightLeaderBoard godoc
// @Tags Leaderboards
// @Summary returns the top 20 players for a given month sorted by knight kdr
// @Produce json
// @Param year query int false "year - defaults to current year"
// @Param month query int false "month - defaults to current month"
// @Success 200 {object} MonthlyLeaderboardList
// @Router /leaderboard/monthly/knight [get]
func GetMonthlyKnightLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []MonthlyStats

	year, month := getYearMonth(r)

	err := db.Select(&stats, monthlyQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND monthly_stats.year=? AND monthly_stats.month=? AND monthly_stats.knight_kills >= ? AND monthly_stats.knight_deaths >= ? 
		ORDER BY (monthly_stats.knight_kills / monthly_stats.knight_deaths) DESC LIMIT 20`, year, month, config.API.KnightGate, config.API.KnightGate)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, MonthlyLeaderboardList{
		Size:        len(stats),
		Year:        year,
		Month:       month,
		LeaderBoard: stats,
	})
}

func GetMonthlyStats(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	var stats []MonthlyStats
	err = db.Select(&stats, monthlyQuery +`WHERE p.ID=? ORDER BY monthly_stats.year DESC LIMIT 12`, playerID)
	if err != nil {
		playerNotFoundError(w, err)
		return
	}

	JSONResponse(w, &stats)
}
