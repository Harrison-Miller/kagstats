package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	var err error
	attempts := 0
	for {
		db, err = sqlx.Connect("mysql", "root:password@tcp(mysql:3306)/kagstats")
		if err != nil {
			log.Printf("%v\n", fmt.Errorf("Couldn't connect to database: %v", err))
		} else {
			break
		}

		attempts=attempts+1
		if attempts > 10 {
			log.Fatal("Could not connect to database after 10 attempts")
		}
		time.Sleep(5 * time.Second)
	}


	r := mux.NewRouter()

	r.HandleFunc("/players", getPlayers).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}", getPlayer).Methods("GET")

	r.HandleFunc("/servers", getServers).Methods("GET")
	r.HandleFunc("/servers/{id:[0-9]+}", getServer).Methods("GET")

	r.HandleFunc("/kills", getKills).Methods("GET")
	r.HandleFunc("/kills/{id:[0-9]+}", getKill).Methods("GET")

	BasicStatsRoutes(r)
	NemesisRoutes(r)
	HitterRoutes(r)
	
	http.Handle("/", r)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
