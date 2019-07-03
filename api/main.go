package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	stats "github.com/Harrison-Miller/kagstats/common/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	config, _ := stats.Get()
	if value, ok := os.LookupEnv("API_HOST"); ok {
		config.API.Host = value
	}

	if value, ok := os.LookupEnv("API_DB"); ok {
		config.DatabaseConnection = value
	}

	var err error
	db, err = sqlx.Connect("mysql", config.DatabaseConnection)
	if err != nil {
		log.Fatal(fmt.Errorf("Couldn't connect to database: %v", err))
	}

	r := mux.NewRouter()

	r.HandleFunc("/players", getPlayers).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}", getPlayer).Methods("GET")

	r.HandleFunc("/servers", getServers).Methods("GET")
	r.HandleFunc("/servers/{id:[0-9]+}", getServer).Methods("GET")

	r.HandleFunc("/kills", getKills).Methods("GET")
	r.HandleFunc("/kills/{id:[0-9]+}", getKill).Methods("GET")

	BasicStatsRoutes(r)

	http.Handle("/", r)

	err = http.ListenAndServe(config.API.Host, r)
	if err != nil {
		log.Fatal(err)
	}
}
