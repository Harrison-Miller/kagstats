package main

import (
	"github.com/Harrison-Miller/kagstats/collector/database"
	"log"
	"os"
	"time"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/kagstats/common/utils"
	_ "github.com/go-sql-driver/mysql"
)

var config configs.Config
var players map[string]models.Player
var kills chan models.Kill
var uncommitted []models.Kill
var updater *PlayerInfoUpdater
var db database.Database

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
		log.Println(err)
	}

	if value, ok := os.LookupEnv("DB"); ok {
		config.DatabaseConnection = value
	}

	if value, ok := os.LookupEnv("MONITOR_HOST"); ok {
		config.Monitoring.Host = value
	}

	version, _ := os.LookupEnv("VERSION")
	log.Printf("KAG Stats  %s\n", version)

	sqlDB, err := utils.ConnectToDatabase(config.DatabaseConnection, 10)
	if err != nil {
		log.Fatal(err)
	}

	db = database.NewSQLDatabase(sqlDB)

	updater = NewPlayerInfoUpdater(db)

	log.Println("Connected to the database!")
	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	players = make(map[string]models.Player)
	kills = make(chan models.Kill, 1000)
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
				err := db.Commit(uncommitted)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Committed %d new kills\n", size)
					uncommitted = make([]models.Kill, 0, 100)
				}
			}
		case <-notify:
			size := len(uncommitted)
			if size > 0 {
				err := db.Commit(uncommitted)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Committed %d new kills\n", size)
					uncommitted = make([]models.Kill, 0, 100)
				}
			}
		}
	}

}
