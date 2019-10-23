package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Harrison-Miller/kagstats/common/utils"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var nextID = 1

func UpdatePlayerInfo(player *models.Player) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO players (username, charactername, clantag) VALUES (?,?,?) ON DUPLICATE KEY UPDATE username=?,charactername=?,clantag=?",
		player.Username, player.Charactername, player.Clantag, player.Username, player.Charactername, player.Clantag)
	if err != nil {
		return errors.Wrap(err, "error updating/creating player")
	}

	row := tx.QueryRow("SELECT ID FROM players WHERE username=?", player.Username)
	err = row.Scan(&player.ID)
	if err != nil {
		return errors.Wrap(err, "error getting player ID")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error commiting player update")
	}

	return nil
}

func UpdateJoinTime(playerID int64, serverID int64) error {
	return AddEvent(playerID, "joined", serverID)
}

func UpdateLeaveTime(playerID int64, serverID int64) error {
	return AddEvent(playerID, "left", serverID)
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS events (
		ID INTEGER PRIMARY KEY AUTO_INCREMENT,
		playerID INT NOT NULL,
		type varchar(30) NOT NULL,
		time INT NOT NULL,
		serverID INT NOT NULL,
		FOREIGN KEY(playerID) REFERENCES players(ID),
		FOREIGN KEY(serverID) REFERENCES servers(ID)
	)`)
	if err != nil {
		return errors.Wrap(err, "error creating events table")
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

func AddEvent(playerID int64, eventType string, serverID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO events (playerID,Type,Time,ServerID) VALUES(?,?,?,?)",
		playerID, eventType, utils.NowAsUnixMilliseconds(), serverID)

	if err != nil {
		return errors.Wrap(err, "error inserting event")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "error getting last inserted id")
	}

	_, err = tx.Exec("UPDATE players SET lastEventID=? WHERE ID=?", id, playerID)
	if err != nil {
		return errors.Wrap(err, "error updating player lastEventID")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing event")
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
