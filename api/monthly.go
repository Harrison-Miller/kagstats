package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const monthlyQuery = `SELECT monthly_stats.*, p.ID "player.ID", p.username "player.username",
p.charactername "player.charactername", p.clantag "player.clantag", p.oldgold "player.oldgold",
p.registered "player.registered", p.role "player.role", p.avatar "player.avatar", p.tier "player.tier",
e.type "player.lastEvent.type", e.time "player.lastEvent.time", e.serverID "player.lastEvent.serverID",
p.gold "player.gold", p.silver "player.silver", p.bronze "player.bronze", p.participation "player.participation",
p.github "player.github", p.community "player.community", p.mapmaker "player.mapmaker", p.moderation "player.moderation",
p.leaderboardBan "player.leaderboardBan", p.statsBan "player.statsBan"
 FROM monthly_stats 
INNER JOIN players as p ON monthly_stats.playerID=p.ID 
INNER JOIN events as e ON p.lastEventID=e.ID `

type MonthlyStats struct {
	Year  int64 `json:"year"`
	Month int64 `json:"month"`
	BasicStats
}

func MonthlyStatsRoutes(r *mux.Router) {
	r.HandleFunc("/leaderboard/monthly/archer", getMonthlyArcherLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/monthly/builder", getMonthlyBuilderLeaderBoard).Methods("GET")
	r.HandleFunc("/leaderboard/monthly/knight", getMonthlyKnightLeaderBoard).Methods("GET")
}

func getYearMonth(r *http.Request) (int, int) {
	t := time.Now()
	year, _ := GetURLParam("year", t.Year(), r)
	month, _ := GetURLParam("month", int(t.Month()), r)
	return year, month
}

func getMonthlyArcherLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []MonthlyStats

	year, month := getYearMonth(r)

	err := db.Select(&stats, monthlyQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND monthly_stats.year=? AND monthly_stats.month=? AND monthly_stats.archer_kills >= ? AND monthly_stats.archer_deaths >= ? 
		ORDER BY (monthly_stats.archer_kills / monthly_stats.archer_deaths) DESC LIMIT 20`, year, month, config.API.ArcherGate, config.API.ArcherGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size        int            `json:"size"`
		Year        int            `json:"year"`
		Month       int            `json:"month"`
		LeaderBoard []MonthlyStats `json:"leaderboard"`
	}{
		Size:        len(stats),
		Year:        year,
		Month:       month,
		LeaderBoard: stats,
	})
}

func getMonthlyBuilderLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []MonthlyStats

	year, month := getYearMonth(r)

	err := db.Select(&stats, monthlyQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND monthly_stats.year=? AND monthly_stats.month=? AND monthly_stats.builder_kills >= ? AND monthly_stats.builder_deaths >= ? 
		ORDER BY (monthly_stats.builder_kills / monthly_stats.builder_deaths) DESC LIMIT 20`, year, month, config.API.BuilderGate, config.API.BuilderGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size        int            `json:"size"`
		Year        int            `json:"year"`
		Month       int            `json:"month"`
		LeaderBoard []MonthlyStats `json:"leaderboard"`
	}{
		Size:        len(stats),
		Year:        year,
		Month:       month,
		LeaderBoard: stats,
	})
}

func getMonthlyKnightLeaderBoard(w http.ResponseWriter, r *http.Request) {
	var stats []MonthlyStats

	year, month := getYearMonth(r)

	err := db.Select(&stats, monthlyQuery+`WHERE NOT p.leaderboardBan AND NOT p.statsBan AND monthly_stats.year=? AND monthly_stats.month=? AND monthly_stats.knight_kills >= ? AND monthly_stats.knight_deaths >= ? 
		ORDER BY (monthly_stats.knight_kills / monthly_stats.knight_deaths) DESC LIMIT 20`, year, month, config.API.KnightGate, config.API.KnightGate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching leader board: %v", err), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, struct {
		Size        int            `json:"size"`
		Year        int            `json:"year"`
		Month       int            `json:"month"`
		LeaderBoard []MonthlyStats `json:"leaderboard"`
	}{
		Size:        len(stats),
		Year:        year,
		Month:       month,
		LeaderBoard: stats,
	})
}
