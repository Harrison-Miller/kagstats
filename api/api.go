package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/pkg/errors"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/Harrison-Miller/kagstats/common/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var config configs.Config
var version string

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		m := httpsnoop.CaptureMetrics(next, w, r)
		log.Printf("%s - %s %v %d %dms\n", r.RemoteAddr, r.Method, r.URL, m.Code, m.Duration/time.Millisecond)
	})
}

func main() {
	config, _ = configs.Get()
	if value, ok := os.LookupEnv("API_HOST"); ok {
		config.API.Host = value
	}

	if value, ok := os.LookupEnv("API_DB"); ok {
		config.DatabaseConnection = value
	}

	var err error
	if value, ok := os.LookupEnv("KD_GATE"); ok {
		config.API.KDGate, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(errors.Wrap(err, "could convert KD_GATE to int"))
		}
	}

	if value, ok := os.LookupEnv("ARCHER_GATE"); ok {
		config.API.ArcherGate, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(errors.Wrap(err, "could convert ARCHER_GATE to int"))
		}
	}

	if value, ok := os.LookupEnv("BUILDER_GATE"); ok {
		config.API.BuilderGate, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(errors.Wrap(err, "could convert BUILDER_GATE to int"))
		}
	}

	if value, ok := os.LookupEnv("KNIGHT_GATE"); ok {
		config.API.KnightGate, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(errors.Wrap(err, "could convert KNIGHT_GATE to int"))
		}
	}

	if value, ok := os.LookupEnv("NEMESIS_GATE"); ok {
		config.API.NemesisGate, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(errors.Wrap(err, "could convert NEMESIS_GATE to int"))
		}
	}

	version, _ = os.LookupEnv("VERSION")
	log.Printf("KAG Stats  %s\n", version)

	db, err = utils.ConnectToDatabase(config.DatabaseConnection, 10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database!")

	r := mux.NewRouter()

	r.HandleFunc("/players", getPlayers).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}", getPlayer).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}/kills", getPlayerKills).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}/events", getPlayerEvents).Methods("GET")
	r.HandleFunc("/players/search/{search:.+}", searchPlayers).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}/refresh", refreshPlayer).Methods("GET")

	r.HandleFunc("/servers", getServers).Methods("GET")
	r.HandleFunc("/servers/{id:[0-9]+}", getServer).Methods("GET")
	r.HandleFunc("/servers/{id:[0-9]+}/events", getServerEvents).Methods("GET")
	r.HandleFunc("/servers/{id:[0-9]+}/kills", getServerKills).Methods("GET")

	r.HandleFunc("/kills", getKills).Methods("GET")
	r.HandleFunc("/kills/{id:[0-9]+}", getKill).Methods("GET")

	BasicStatsRoutes(r)
	NemesisRoutes(r)
	HitterRoutes(r)

	r.Use(LogHandler)

	log.Println("Ready to accept connections on ", config.API.Host)

	err = http.ListenAndServe(config.API.Host, r)
	if err != nil {
		log.Fatal(err)
	}
}
