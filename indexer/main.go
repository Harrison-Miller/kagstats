package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const defaultInterval = 5 * time.Second
const defaultBatchSize int64 = 100

var interval = defaultInterval
var batchSize = defaultBatchSize

type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/", c.User, c.Password, c.Address)
}

type Config struct {
	Database         DatabaseConfig `json:"database"`
	IndexerBatchSize int64          `json:"indexerBatchSize"`
	IndexerInterval  string         `json:"indexerInterval"`
}

func initDatabase(db *sql.DB) {
	_, err := db.Exec("USE kagstats")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ingest_info (
			key_name VARCHAR(30) PRIMARY KEY,
			value INT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec("INSERT INTO kagstats.ingest_info (key_name, value) VALUES(?, ?)", "basic_stats_index", 0)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS basic_stats (
			playerID INT PRIMARY KEY,
			suicides INT NOT NULL,
			teamkills INT NOT NULL,
			archer_kills INT NOT NULL,
			archer_deaths INT NOT NULL,
			builder_kills INT NOT NULL,
			builder_deaths INT NOT NULL,
			knight_kills INT NOT NULL,
			knight_deaths INT NOT NULL,
			other_kills INT NOT NULL,
			other_deaths INT NOT NULL,
			FOREIGN KEY(playerID) REFERENCES players(ID)
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

type StatsUpdate struct {
	playerID      int64
	suicides      int64
	teamkills     int64
	archerKills   int64
	archerDeaths  int64
	builderKills  int64
	builderDeaths int64
	knightKills   int64
	knightDeaths  int64
	otherKills    int64
	otherDeaths   int64
}

func indexStats(db *sql.DB) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	var index int64
	row := tx.QueryRow("SELECT value FROM kagstats.ingest_info WHERE key_name=?", "basic_stats_index")
	err = row.Scan(&index)
	newIndex := index

	if err != nil {
		return 0, err // ingest_index probably not created in table
	}

	rows, err := tx.Query("SELECT ID, killerID, victimID, killerClass, victimClass, teamKill FROM kagstats.kills WHERE ID>? AND ID<?",
		index, index+batchSize+1)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var updates map[int64]*StatsUpdate = make(map[int64]*StatsUpdate)

	for rows.Next() {
		var (
			id          int64
			killerID    sql.NullInt64
			victimID    int64
			killerClass string
			victimClass string
			teamKill    bool
		)

		if err := rows.Scan(&id, &killerID, &victimID, &killerClass, &victimClass, &teamKill); err != nil {
			tx.Rollback()
			return 0, err
		}

		if killerID.Valid {
			// add stats for killer
			var killer *StatsUpdate
			if k, ok := updates[killerID.Int64]; ok {
				killer = k
			} else {
				killer = new(StatsUpdate)
				updates[killerID.Int64] = killer
			}

			killer.playerID = killerID.Int64

			if teamKill {
				killer.teamkills++
			} else {
				switch killerClass {
				case "archer":
					killer.archerKills++
				case "builder":
					killer.builderKills++
				case "knight":
					killer.knightKills++
				default:
					killer.otherKills++
				}
			}
		}

		var victim *StatsUpdate
		if v, ok := updates[victimID]; ok {
			victim = v
		} else {
			victim = new(StatsUpdate)
			updates[victimID] = victim
		}

		victim.playerID = victimID

		if !killerID.Valid {
			victim.suicides++
		}

		switch victimClass {
		case "archer":
			victim.archerDeaths++
		case "builder":
			victim.builderDeaths++
		case "knight":
			victim.knightDeaths++
		default:
			victim.otherDeaths++
		}

		newIndex = id
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return 0, err
	}

	stmt, err := tx.Prepare(`INSERT INTO kagstats.basic_stats (playerID, suicides, teamkills, archer_kills,
		archer_deaths, builder_kills, builder_deaths, knight_kills, knight_deaths, other_kills, other_deaths)
		VALUES (?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE suicides=suicides+?,teamkills=teamkills+?,
		archer_kills=archer_kills+?,archer_deaths=archer_deaths+?,builder_kills=builder_kills+?,
		builder_deaths=builder_deaths+?,knight_kills=knight_kills+?,knight_deaths=knight_deaths+?,
		other_kills=other_kills+?,other_deaths=other_deaths+?`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	defer stmt.Close()

	for _, u := range updates {
		_, err = stmt.Exec(u.playerID, u.suicides, u.teamkills, u.archerKills, u.archerDeaths, u.builderKills,
			u.builderDeaths, u.knightKills, u.knightDeaths, u.otherKills, u.otherDeaths, u.suicides, u.teamkills,
			u.archerKills, u.archerDeaths, u.builderKills, u.builderDeaths, u.knightKills, u.knightDeaths,
			u.otherKills, u.otherDeaths)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	_, err = tx.Exec("UPDATE kagstats.ingest_info SET value=? WHERE key_name='basic_stats_index' AND value=?", newIndex, index)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()

	return newIndex - index, err

}

func main() {
	configPath := "settings.json"
	if value, ok := os.LookupEnv("KAGSTATS_CONFIG"); ok {
		configPath = value
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		log.Fatal(err)
	}

	interval, err = time.ParseDuration(config.IndexerInterval)
	if err != nil {
		interval = defaultInterval
	}

	db, err := sql.Open("mysql", config.Database.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	initDatabase(db)

	for {
		processed, err := indexStats(db)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Processed %d kills\n", processed)
		}
		time.Sleep(interval)
	}

}
