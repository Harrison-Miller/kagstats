package main

import (
	"log"
	"os"
	"time"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/kagstats/common/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var config configs.Config
var db *sqlx.DB
var players map[string]models.Player
var kills chan models.Kill
var uncommitted []models.Kill

func commitTimer(notify chan bool) {
	for {
		time.Sleep(config.CommitIntervalDuration)
		notify <- true
	}
}

func main() {
	log.SetPrefix("[main] ")
	var err error
	config, err = configs.Get()
	if err != nil {
		log.Fatal(err)
	}

	if value, ok := os.LookupEnv("DB"); ok {
		config.DatabaseConnection = value
	}

	if value, ok := os.LookupEnv("MONITOR_HOST"); ok {
		config.Monitoring.Host = value
	}

	db, err = utils.ConnectToDatabase(config.DatabaseConnection, 10)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database!")
	err = InitDB()
	if err != nil {
		log.Fatal(err)
	}

	players = make(map[string]models.Player)
	kills = make(chan models.Kill)
	uncommitted = make([]models.Kill, 0, 100)

	for _, c := range config.Servers {
		go Collect(c)
		time.Sleep(1 * time.Second)
	}

	notify := make(chan bool)
	go commitTimer(notify)

	for {
		select {
		case kill := <-kills:
			uncommitted = append(uncommitted, kill)
			size := len(uncommitted)
			if size >= config.BatchSize {
				err := Commit()
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Committed %d new kills\n", size)
				}
			}
		case <-notify:
			size := len(uncommitted)
			if size > 0 {
				err := Commit()
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Committed %d new kills\n", size)
				}
			}
		}
	}

}
