package main

import (
	"strings"
	"time"

	"github.com/Harrison-Miller/kagstats/common/models"
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

	_, err = tx.Exec("INSERT INTO servers (name, description, gamemode, tags) VALUE (?,?,?,?) ON DUPLICATE KEY UPDATE name=?, description=?, gamemode=?, tags=?",
		server.Name, server.Description, server.Gamemode, server.Tags,
		server.Name, server.Description, server.Gamemode, server.Tags)
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

	_, err = tx.Exec("INSERT INTO events (playerID,Type,Time,ServerID) VALUES(?,?,?,?)",
		playerID, eventType, time.Now().Unix(), serverID)

	if err != nil {
		return errors.Wrap(err, "error inserting event")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing event")
	}
	return nil
}
