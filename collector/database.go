package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var nextID = 1

func UpdatePlayerInfo(player *models.Player) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO players (username, charactername, clantag, lastIP) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE username=?,charactername=?,clantag=?,lastIP=?",
		player.Username, player.Charactername, player.Clantag, player.IP, player.Username, player.Charactername, player.Clantag, player.IP)
	if err != nil {
		return errors.Wrap(err, "error updating/creating player")
	}

	err = tx.Get(player, "SELECT * FROM players WHERE username=?", player.Username)
	if err != nil {
		return errors.Wrap(err, "error getting player ID")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error commiting player update")
	}

	return nil
}

func UpdateServerInfo(server *models.Server) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO servers (name, description, gamemode, address, port, tags, status) VALUE (?,?,?,?,?,?, true)
		ON DUPLICATE KEY UPDATE name=?, description=?, gamemode=?, address=?, port=?, tags=?, status=true`,
		server.Name, server.Description, server.Gamemode, server.Address, server.Port, server.Tags,
		server.Name, server.Description, server.Gamemode, server.Address, server.Port, server.Tags)
	if err != nil {
		return errors.Wrap(err, "error creating/updating server info")
	}

	row := tx.QueryRow("SELECT ID FROM servers WHERE name=?", server.Name)
	err = row.Scan(&server.ID)
	if err != nil {
		return errors.Wrap(err, "error getting server ID")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing server info")
	}
	return nil
}

func UpdateServerStatus(server models.Server, status bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE servers SET status=? WHERE ID=?`, status, server.ID)
	if err != nil {
		return errors.Wrap(err, "error setting server status")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing server status")
	}
	return nil
}

func InitDB() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS players (
		ID INT PRIMARY KEY AUTO_INCREMENT,
		username varchar(30) NOT NULL UNIQUE,
		charactername varchar(30) NOT NULL,
		clantag varchar(30) NOT NULL
		) CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin`)
	if err != nil {
		return errors.Wrap(err, "error creating player table")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS servers (
		ID INTEGER PRIMARY KEY AUTO_INCREMENT,
		name varchar(255) NOT NULL UNIQUE,
		description varchar(255) NOT NULL,
		gamemode varchar(30) NOT NULL,
		address varchar(30) NOT NULL,
		port varchar(30) NOT NULL,
		tags varchar(1000)
	)`)
	if err != nil {
		return errors.Wrap(err, "error creating servers table")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS kills (
		ID INTEGER PRIMARY KEY AUTO_INCREMENT,
		killerID INT NOT NULL,
		victimID INT NOT NULL,
		killerClass ENUM('archer', 'builder', 'knight', 'other', 'none') DEFAULT 'none',
		victimClass ENUM('archer', 'builder', 'knight', 'other') DEFAULT 'archer' NOT NULL,
		hitter INT DEFAULT 0,
		epoch INT NOT NULL,
		serverID INT NOT NULL,
		teamKill BOOLEAN DEFAULT false,
		FOREIGN KEY(killerID) REFERENCES players(ID),
		FOREIGN KEY(victimID) REFERENCES players(ID),
		FOREIGN KEY(serverID) REFERENCES servers(ID)
	)`)
	if err != nil {
		return errors.Wrap(err, "error creating kills table")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS stats_info (
		key_name VARCHAR(30) PRIMARY KEY,
		value INT NOT NULL	
	)`)
	if err != nil {
		return err
	}

	db.Exec("INSERT INTO stats_info (key_name, value) VALUES(?,?)", "database_version", 0)

	err = RunMigrations(db)
	if err != nil {
		return errors.Wrap(err, "error running migrations")
	}
	return nil
}

func Commit() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//https://stackoverflow.com/questions/21108084/how-to-insert-multiple-data-at-once
	values := []interface{}{}
	sqlStr := "INSERT INTO kagstats.kills (killerID, victimID, killerClass, victimClass, hitter, epoch, serverID, teamKill) VALUES "
	for _, v := range uncommitted {
		sqlStr += "(?,?,?,?,?,?,?,?),"
		values = append(values, v.KillerID, v.VictimID,
			v.KillerClass, v.VictimClass, v.Hitter,
			v.Time, v.ServerID, v.TeamKill)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	stmnt, err := tx.Prepare(sqlStr)
	if err != nil {
		return errors.Wrap(err, "error preparing bulk load statement")
	}
	defer stmnt.Close()

	_, err = stmnt.Exec(values...)
	if err != nil {
		return errors.Wrap(err, "error executing bulk load")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "error commiting bulk load")
	}

	uncommitted = make([]models.Kill, 0, 100)
	return nil
}

func CommitFlagCapture(capture models.FlagCapture) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO flag_captures (playerID,ticks) VALUES (?,?)",
		capture.PlayerID, capture.Ticks)

	if err != nil {
		return errors.Wrap(err, "error inserting flag capture")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing flag capture")
	}

	return nil
}

func CommitMapStats(stats models.MapStats) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO map_stats (mapName,ticks) VALUES (?,?)",
		stats.MapName, stats.Ticks)

	if err != nil {
		return errors.Wrap(err, "error inserting map stats")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing map stats")
	}

	return nil
}

func CommitMapVotes(votes models.MapVotes) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO map_votes (map1Name,map1Votes,map2Name,map2Votes,randomVotes) VALUES (?,?,?,?,?)",
		votes.Map1Name, votes.Map1Votes, votes.Map2Name, votes.Map2Votes, votes.RandomVotes)

	if err != nil {
		return errors.Wrap(err, "error inserting map votes")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing map votes")
	}

	return nil
}

func RunMigration(version int64, migrations func(db *sqlx.DB) error, db *sqlx.DB) error {
	db.Exec("INSERT INTO stats_info (key_name, value) VALUES(?,?)", "database_version", version)
	row := db.QueryRow("SELECT value FROM stats_info WHERE key_name=?", "database_version")
	var currentVersion int64
	err := row.Scan(&currentVersion)
	if err != nil {
		return errors.Wrap(err, "error getting current database version")
	}

	if currentVersion < version {
		log.Printf("Running migrations for version %d", version)
		err = migrations(db)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error running migrations for version %d", version))
		}

		_, err = db.Exec("UPDATE stats_info SET value=? WHERE key_name=?", version, "database_version")
	} else {
		log.Printf("Skipping migrations for version %d", version)
	}

	return nil
}
