package main

import (
	"github.com/Harrison-Miller/kagstats/api/docs"
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

	"github.com/swaggo/http-swagger"
	_ "github.com/Harrison-Miller/kagstats/api/docs"
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

// @title KAG Stats
// @description KAG Stats records statics of players in King Arthur's Gold such as kills, flag captures, map votes

// @license.name MIT
// @license.url https://github.com/Harrison-Miller/kagstats/blob/master/LICENSE

// @host kagstats.com
// @BasePath /api
// @Schemes https

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

	if value, ok := os.LookupEnv("AUTH_SECRET"); ok {
		AUTH_SECRET = value
	}

	version, _ = os.LookupEnv("VERSION")
	log.Printf("KAG Stats  %s\n", version)

	docs.SwaggerInfo.Version = version

	db, err = utils.ConnectToDatabase(config.DatabaseConnection, 10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database!")

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/players", GetPlayers).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}", GetPlayer).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}/kills", GetPlayerKills).Methods("GET")
	r.HandleFunc("/players/search/{search:.+}", SearchPlayers).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}/captures", GetCaptures).Methods("GET")

	r.HandleFunc("/servers", GetServers).Methods("GET")
	r.HandleFunc("/servers/{id:[0-9]+}", GetServer).Methods("GET")
	r.HandleFunc("/servers/{id:[0-9]+}/kills", GetServerKills).Methods("GET")

	r.HandleFunc("/kills", GetKills).Methods("GET")
	r.HandleFunc("/kills/{id:[0-9]+}", GetKill).Methods("GET")

	BasicStatsRoutes(r)
	NemesisRoutes(r)
	HitterRoutes(r)
	MonthlyStatsRoutes(r)
	MapsRoutes(r)


	protected := r.NewRoute().Subrouter()
	protected.Use(Verify)
	ClanRoutes(r, protected)
	AuthRoutes(r, protected)
	FollowingRoutes(r, protected)

	r.Use(LogHandler)

	log.Println("Ready to accept connections on ", config.API.Host)

	err = http.ListenAndServe(config.API.Host, r)
	if err != nil {
		log.Fatal(err)
	}
}
