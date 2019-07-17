package main

import (
	"log"
	"strings"
	"time"

	"github.com/Harrison-Miller/kagstats/common/configs"

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

func UpdateJoinTime(player models.Player) error {
	log.Printf("%s join at %v", player.Username, time.Now())
	return nil
}

func UpdateLeaveTime(player models.Player) error {
	log.Printf("%s left at %v", player.Username, time.Now())
	return nil
}

func UpdateServerInfo(server configs.ServerConfig) (models.Server, error) {
	tx, err := db.Begin()
	if err != nil {
		return models.Server{}, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO servers (name, tags) VALUE (?,?) ON DUPLICATE KEY UPDATE name=?, tags=?",
		server.Name, server.TagsString(), server.Name, server.TagsString())
	if err != nil {
		return models.Server{}, errors.Wrap(err, "error creating/updating server info")
	}

	s := models.Server{
		Name: server.Name,
		Tags: server.TagsString(),
	}
	row := tx.QueryRow("SELECT ID FROM servers WHERE name=?", server.Name)
	err = row.Scan(&s.ID)
	if err != nil {
		return s, errors.Wrap(err, "error getting server ID")
	}

	if err := tx.Commit(); err != nil {
		return s, errors.Wrap(err, "error committing server info")
	}
	return s, nil
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
