package main

import (
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/Harrison-Miller/kagstats/common/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/felixge/httpsnoop"
)

var db *sqlx.DB
var config configs.Config
var username = "admin"
var password = "admin1234!"
var host = ":8080"
var apiHost = "http://localhost/api/"
var prefix = ""

func JSONResponse(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(i)
}

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		m := httpsnoop.CaptureMetrics(next, w, r)
		log.Printf("%s - %s %v %d %dms\n", r.RemoteAddr, r.Method, r.URL, m.Code, m.Duration/time.Millisecond)
	})
}

func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="login"`)
			http.Error(w, "Unathorized", 401)
			return
		}

		handler(w, r)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "unable to load index.html", http.StatusInternalServerError)
		return
	}

	t.Execute(w, struct {
		APIHost string
		Prefix  string
	}{
		APIHost: apiHost,
		Prefix:  prefix,
	})
}

func player(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/player.html")
	if err != nil {
		http.Error(w, "unable to load player.html", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	t.Execute(w, struct {
		APIHost  string
		Prefix   string
		PlayerId string
	}{
		APIHost:  apiHost,
		Prefix:   prefix,
		PlayerId: vars["id"],
	})
}

type SaveParam struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	LeaderboardBan bool   `json:"leaderboardBan"`
	MonthlyLeaderboardBan bool `json:"monthlyLeaderboardBan"`
	StatsBan       bool   `json:"statsBan"`
	Notes          string `json:"notes"`
}

func save(w http.ResponseWriter, r *http.Request) {
	var params SaveParam
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.Username == "" || params.ID == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete user from database
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE players SET username=?,leaderboardBan=?,monthlyLeaderboardBan=?,statsBan=?,notes=? WHERE ID=?", params.Username, params.LeaderboardBan, params.MonthlyLeaderboardBan, params.StatsBan, params.Notes, params.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func delete(w http.ResponseWriter, r *http.Request) {
	var params SaveParam
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write changes to the database
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM players WHERE ID=?", params.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func recalculate(w http.ResponseWriter, r *http.Request) {
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT key_name FROM indexer_info WHERE key_name NOT LIKE '%_version'")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		keys = append(keys, key)

	}

	for _, key := range keys {
		_, err := tx.Exec("UPDATE indexer_info SET value=0 WHERE key_name=?", key)
		if err != nil {
			log.Println("error reseting counter")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("DELETE FROM " + key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type NotesParams struct {
	ID    string `json:"id" db:"ID"`
	Notes string `json:"notes"`
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId := vars["id"]

	var n NotesParams

	err := db.Get(&n, "SELECT ID,notes FROM players WHERE ID=?", playerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &n)
}

type AltInfo struct {
	Avatar string `json:"avatar"`
	IP string `json:"ip" db:"lastIP"`
	ClanTag string `json:"clantag" db:"clantag"`
	CharacterName string `json:"charactername"`
	Username string `json:"username" db:"username"`
	PlayerID int64 `json:"id" db:"ID"`
	LeaderboardBan bool   `json:"leaderboardBan" db:"leaderboardBan"`
	MonthlyLeaderboardBan bool `json:"monthlyLeaderboardBan" db:"monthlyLeaderboardBan"`
	StatsBan       bool   `json:"statsBan" db:"statsBan"`
}

type LastIP struct {
	IP string `json:"ip" db:"lastIP"`
}

func getAlts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID := vars["id"]

	var ip LastIP
	err := db.Get(&ip, `SELECT lastIP FROM players WHERE players.ID=?`, playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alts := []AltInfo{}
	if ip.IP == "" {
		JSONResponse(w, &alts)
		return
	}

	err = db.Select(&alts, `SELECT avatar,lastIP,clantag,charactername,username,ID,leaderboardBan,monthlyLeaderboardBan,statsBan
		FROM players WHERE players.lastIP=? AND players.ID!=?`, ip.IP, playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &alts)
}

type AltCountInfo struct {
	Avatar string `json:"avatar"`
	ClanTag string `json:"clantag" db:"clantag"`
	IP string `json:"ip" db:"lastIP"`
	CharacterName string `json:"charactername"`
	Username string `json:"username" db:"username"`
	PlayerID int64 `json:"id" db:"ID"`
	Alts int64 `json:"altsCount"`
}

func getTopAlts(w http.ResponseWriter, r *http.Request) {
	ArcherGate := 40
	BuilderGate := 20
	KnightGate := 75

	archerAlts := []AltCountInfo{}
	err := db.Select(&archerAlts, 	`SELECT p.avatar,p.clantag,p.lastIP,p.charactername,p.username,p.ID FROM basic_stats 
	INNER JOIN players as p ON basic_stats.playerID=p.ID 
	WHERE NOT p.leaderboardBan AND NOT p.statsBan 
	AND basic_stats.archer_kills >= ? AND basic_stats.archer_deaths >= ? 
	ORDER BY (basic_stats.archer_kills / basic_stats.archer_deaths) DESC LIMIT 20`, ArcherGate, ArcherGate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, alt := range archerAlts {
		if alt.IP != "" {
			var count int64
			db.Get(&count, `SELECT COUNT(ID) FROM players WHERE players.lastIP=? AND players.ID!=?`, alt.IP, alt.PlayerID)
			archerAlts[i].Alts = count
		}
	}

	builderAlts := []AltCountInfo{}
	err = db.Select(&builderAlts, 	`SELECT p.avatar,p.clantag,p.lastIP,p.charactername,p.username,p.ID FROM basic_stats 
	INNER JOIN players as p ON basic_stats.playerID=p.ID 
	WHERE NOT p.leaderboardBan AND NOT p.statsBan 
	AND basic_stats.builder_kills >= ? AND basic_stats.builder_deaths >= ? 
	ORDER BY (basic_stats.builder_kills / basic_stats.builder_deaths) DESC LIMIT 20`, BuilderGate, BuilderGate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, alt := range builderAlts {
		if alt.IP != "" {
			var count int64
			db.Get(&count, `SELECT COUNT(ID) FROM players WHERE players.lastIP=? AND players.ID!=?`, alt.IP, alt.PlayerID)
			builderAlts[i].Alts = count
		}
	}

	knightAlts := []AltCountInfo{}
	err = db.Select(&knightAlts, 	`SELECT p.avatar,p.clantag,p.lastIP,p.charactername,p.username,p.ID FROM basic_stats 
	INNER JOIN players as p ON basic_stats.playerID=p.ID 
	WHERE NOT p.leaderboardBan AND NOT p.statsBan 
	AND basic_stats.knight_kills >= ? AND basic_stats.knight_deaths >= ? 
	ORDER BY (basic_stats.knight_kills / basic_stats.knight_deaths) DESC LIMIT 20`, KnightGate, KnightGate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, alt := range knightAlts {
		if alt.IP != "" {
			var count int64
			db.Get(&count, `SELECT COUNT(ID) FROM players WHERE players.lastIP=? AND players.ID!=?`, alt.IP, alt.PlayerID)
			knightAlts[i].Alts = count
		}
	}

	alts := []AltCountInfo{}
	for _, alt := range archerAlts {
		if alt.Alts > 0 {
			alts = append(alts, alt)
		}
	}

	for _, alt := range knightAlts {
		if alt.Alts > 0 {
			alts = append(alts, alt)
		}
	}

	for _, alt := range builderAlts {
		if alt.Alts > 0 {
			alts = append(alts, alt)
		}
	}

	JSONResponse(w, &alts)

}

func main() {
	config, _ = configs.Get()
	if value, ok := os.LookupEnv("ADMIN_DB"); ok {
		config.DatabaseConnection = value
	}

	if value, ok := os.LookupEnv("USER"); ok {
		username = value
	}

	if value, ok := os.LookupEnv("PASS"); ok {
		password = value
	}

	if value, ok := os.LookupEnv("HOST"); ok {
		host = value
	}

	if value, ok := os.LookupEnv("API_HOST"); ok {
		apiHost = value
	}

	if value, ok := os.LookupEnv("PREFIX"); ok {
		prefix = value
	}

	/*TODO:
	* paginate nemesis/bullied
	* collector info/status
	* DELETE FROM players WHERE username REGEXP '^.*~[0-9]+';
	 */

	var err error
	db, err = utils.ConnectToDatabase(config.DatabaseConnection, 10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database!")

	r := mux.NewRouter()
	r.HandleFunc("/", BasicAuth(index)).Methods("GET")
	r.HandleFunc("/player/{id:[0-9]+}", BasicAuth(player)).Methods("GET")
	r.HandleFunc("/save", BasicAuth(save)).Methods("POST")
	r.HandleFunc("/delete", BasicAuth(delete)).Methods("POST")
	r.HandleFunc("/recalculate", BasicAuth(recalculate)).Methods("POST")
	r.HandleFunc("/notes/{id:[0-9]+}", BasicAuth(getNotes)).Methods("GET")
	r.HandleFunc("/alts/{id:[0-9]+}", BasicAuth(getAlts)).Methods("GET")
	r.HandleFunc("/leaderboardalts", BasicAuth(getTopAlts)).Methods("GET")

	r.Use(LogHandler)

	err = http.ListenAndServe(host, r)
	if err != nil {
		log.Fatal(err)
	}
}
