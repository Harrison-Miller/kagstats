package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/kagstats/common/utils"
	. "github.com/Harrison-Miller/kagstats/indexer"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func ProcessPlayers(currentIndex int, maxRows int, db *sqlx.DB) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "error starting transaction")
	}
	defer tx.Rollback()

	var players []models.Player
	err = db.Select(&players, "SELECT * FROM players LIMIT ?,?", currentIndex, maxRows)
	if err != nil {
		return 0, errors.Wrap(err, "error getting players")
	}

	for _, p := range players {
		err = utils.GetPlayerAvatar(&p)
		if err != nil {
			log.Println(err)
		}

		err = utils.GetPlayerTier(&p)
		if err != nil {
			log.Println(err)
		}

		utils.GetPlayerInfo(&p)
		if err != nil {
			log.Println(err)
		}

		_, err = db.Exec("UPDATE players SET oldgold=?,registered=?,role=?,avatar=?,tier=? WHERE ID=?",
			p.OldGold, p.Registered, p.Role, p.Avatar, p.Tier, p.ID)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("error updating player %s with api info", p.Username))
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "error commiting updated player info")
	}

	return len(players), nil
}

func main() {
	log.Println("Starting API cache")

	config, err := ReadConfig()
	if err != nil {
		log.Fatal("Error reading api cache configuration: ", err)
	}

	db, err := utils.ConnectToDatabase(config.DatabaseConnection, 10)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database!")

	fmt.Printf("Batch Size: %d\n", config.Indexer.BatchSize)
	fmt.Printf("Processing Interval: %s\n", config.Indexer.Interval)

	currentIndex := 0
	for {
		processed, err := ProcessPlayers(currentIndex, config.Indexer.BatchSize, db)
		if err != nil {
			log.Println(err)
		} else if processed != 0 {
			log.Printf("Processed %d rows\n", processed)
		}

		currentIndex += processed
		time.Sleep(config.Indexer.IntervalDuration)
	}
}
