package database

import (
	"fmt"
	"strings"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var nextID = 1

//go:generate moq -pkg fixtures -out ../fixtures/database.go . Database
type Database interface {
	InitDB() error
	UpdatePlayerInfo(player *models.Player) error
	UpdateServerInfo(server *models.Server) error
	UpdateServerStatus(server models.Server, status bool) error
	Commit(kills []models.Kill) error
	CommitFlagCapture(capture models.FlagCapture) error
	CommitMapStats(stats models.MapStats) error
	CommitMapVotes(votes models.MapVotes) error
	CommitPlayer(player *models.Player) error
}

type SQLDatabase struct {
	db *sqlx.DB
}

func NewSQLDatabase(db *sqlx.DB) *SQLDatabase {
	return &SQLDatabase{db: db}
}

func (d *SQLDatabase) UpdatePlayerInfo(player *models.Player) error {
	tx, err := d.db.Beginx()
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

func (d *SQLDatabase) UpdateServerInfo(server *models.Server) error {
	tx, err := d.db.Begin()
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

func (d *SQLDatabase) UpdateServerStatus(server models.Server, status bool) error {
	tx, err := d.db.Begin()
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

func (d *SQLDatabase) InitDB() error {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS players (
		ID INT PRIMARY KEY AUTO_INCREMENT,
		username varchar(30) NOT NULL UNIQUE,
		charactername varchar(30) NOT NULL,
		clantag varchar(30) NOT NULL
		) CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin`)
	if err != nil {
		return errors.Wrap(err, "error creating player table")
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS servers (
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

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS kills (
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

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS stats_info (
		key_name VARCHAR(30) PRIMARY KEY,
		value INT NOT NULL	
	)`)
	if err != nil {
		return err
	}

	d.db.Exec("INSERT INTO stats_info (key_name, value) VALUES(?,?)", "database_version", 0)

	err = d.RunMigrations()
	if err != nil {
		return errors.Wrap(err, "error running migrations")
	}
	return nil
}

func (d *SQLDatabase) Commit(kills []models.Kill) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//https://stackoverflow.com/questions/21108084/how-to-insert-multiple-data-at-once
	values := []interface{}{}
	sqlStr := "INSERT INTO kagstats.kills (killerID, victimID, killerClass, victimClass, hitter, epoch, serverID, teamKill) VALUES "
	for _, v := range kills {
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
	return nil
}

func (d *SQLDatabase) CommitFlagCapture(capture models.FlagCapture) error {
	tx, err := d.db.Begin()
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

func (d *SQLDatabase) CommitMapStats(stats models.MapStats) error {
	tx, err := d.db.Begin()
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

func (d *SQLDatabase) CommitMapVotes(votes models.MapVotes) error {
	tx, err := d.db.Begin()
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

func (d *SQLDatabase) CommitPlayer(player *models.Player) error {
	tx, err := d.db.Begin()
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	defer tx.Rollback()


	_, err = tx.Exec("UPDATE players SET oldgold=?,registered=?,role=?,avatar=?,tier=? WHERE ID=?",
		player.OldGold, player.Registered, player.Role, player.Avatar, player.Tier, player.ID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error updating player %s with api info", player.Username))
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error commiting updated player info")
	}
	return nil
}
