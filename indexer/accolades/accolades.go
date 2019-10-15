package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/Harrison-Miller/kagstats/common/utils"
	. "github.com/Harrison-Miller/kagstats/indexer"
)

const ACCOLADE_URL = "https://raw.githubusercontent.com/transhumandesign/kag-base/master/Rules/CommonScripts/accolade_data.cfg"

var db *sqlx.DB

type Accolades struct {
	GoldCount          int
	SilverCount        int
	BronzeCount        int
	ParticipationCount int
	Github             bool
	Community          bool
	MapMaker           bool
	Moderation         bool
}

func ParseAccolades() (error, int) {
	resp, err := http.Get(ACCOLADE_URL)
	if err != nil {
		return err, 0
	}

	defer resp.Body.Close()

	var username string
	var accolades Accolades

	var reader = bufio.NewReader(resp.Body)

	var usercount = 0
	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			err = UpdateUser(username, accolades)
			if err != nil {
				return err, usercount
			}
			usercount++
			break
		}

		// Strip comments
		commentIndex := strings.Index(line, "#")
		if commentIndex != -1 {
			line = line[:commentIndex]
		}

		// Parse for username (key)
		equalIndex := strings.Index(line, "=")
		if equalIndex != -1 {
			// Update the last user we parsed
			if username != "" {
				err = UpdateUser(username, accolades)
				if err != nil {
					return err, usercount
				}
				usercount++
			}
			username = strings.TrimSpace(line[:equalIndex])
			accolades = Accolades{}

			// Strip username and equal
			line = line[equalIndex+1:]
		}

		// Parse for accolades (values)
		var prev = 0
		var i = strings.Index(line, ";")

		for ; i > -1; i = strings.Index(line[i+1:], ";") {
			a := strings.TrimSpace(line[prev:i])
			prev = i

			fields := strings.Fields(a)
			name := fields[0]
			var count int
			if len(fields) > 1 {
				count, _ = strconv.Atoi(fields[1])
			}

			if name == "gold" {
				accolades.GoldCount = count
			} else if name == "silver" {
				accolades.SilverCount = count
			} else if name == "bronze" {
				accolades.BronzeCount = count
			} else if name == "participation" {
				accolades.ParticipationCount = count
			} else if name == "github" {
				accolades.Github = true
			} else if name == "community" {
				accolades.Community = true
			} else if name == "map" {
				accolades.MapMaker = true
			} else if name == "moderation" {
				accolades.Moderation = true
			}

		}
	}

	return nil, usercount
}

func UpdateUser(username string, accolades Accolades) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO players (username, charactername, clantag, gold, silver, bronze, participation, 
		github, community, mapmaker, moderation) 
		VALUES(?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE 
		username=?,gold=?,silver=?,bronze=?,participation=?,github=?,community=?,mapmaker=?,moderation=?`,
		username, username, "", accolades.GoldCount, accolades.SilverCount, accolades.BronzeCount, accolades.ParticipationCount,
		accolades.Github, accolades.Community, accolades.MapMaker, accolades.Moderation,
		username, accolades.GoldCount, accolades.SilverCount, accolades.BronzeCount, accolades.ParticipationCount,
		accolades.Github, accolades.Community, accolades.MapMaker, accolades.Moderation)

	if err != nil {
		return errors.Wrap(err, "error updating player with accolades")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error commiting player accolades")
	}

	return nil
}

func main() {
	log.Println("Starting Accolade Parser")

	config, err := ReadConfig()
	if err != nil {
		log.Fatal("Error reading accolades parser configuration: ", err)
	}

	db, err = utils.ConnectToDatabase(config.DatabaseConnection, 10)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database!")

	for {
		log.Println("Parsing Accolades")
		err, count := ParseAccolades()
		log.Printf("Updated %d players accolades\n", count)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(24 * time.Hour)
	}
}
