package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Harrison-Miller/kagstats/common/configs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	config, _ := configs.Get()
	if value, ok := os.LookupEnv("API_HOST"); ok {
		config.API.Host = value
	}

	if value, ok := os.LookupEnv("API_DB"); ok {
		config.DatabaseConnection = value
	}

	var err error
	attempts := 0
	for {
		db, err = sqlx.Connect("mysql", config.DatabaseConnection)
		if err != nil {
			log.Printf("%v\n", fmt.Errorf("Couldn't connect to database: %v", err))
		} else {
			break
		}

		attempts = attempts + 1
		if attempts > 10 {
			log.Fatal("Could not connect to database after 10 attempts")
		}
		time.Sleep(5 * time.Second)
	}

	r := mux.NewRouter()

	apiPath := `/api`

	r.HandleFunc(apiPath+"/players", getPlayers).Methods("GET")
	r.HandleFunc(apiPath+"/players/{id:[0-9]+}", getPlayer).Methods("GET")
	r.HandleFunc(apiPath+"/players/{id:[0-9]+}/kills", getPlayerKills).Methods("GET")

	r.HandleFunc(apiPath+"/servers", getServers).Methods("GET")
	r.HandleFunc(apiPath+"/servers/{id:[0-9]+}", getServer).Methods("GET")

	r.HandleFunc(apiPath+"/kills", getKills).Methods("GET")
	r.HandleFunc(apiPath+"/kills/{id:[0-9]+}", getKill).Methods("GET")

	BasicStatsRoutes(r)
	NemesisRoutes(r)
	HitterRoutes(r)

	http.Handle("/", r)

	err = http.ListenAndServe(config.API.Host, r)
	if err != nil {
		log.Fatal(err)
	}
}
