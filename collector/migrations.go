package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// add database migrations here
func RunMigrations(db *sqlx.DB) error {
	err := RunMigration(1, APICacheChanges, db)
	if err != nil {
		return err
	}

	err = RunMigration(2, TeamKillChanges, db)
	if err != nil {
		return err
	}

	err = RunMigration(3, BumpNameLimit, db)
	if err != nil {
		return err
	}

	err = RunMigration(4, SawKillBuilderOnly, db)
	if err != nil {
		return err
	}

	err = RunMigration(5, AddLastEventToPlayers, db)
	if err != nil {
		return err
	}

	err = RunMigration(6, AddAccolades, db)
	if err != nil {
		return err
	}

	err = RunMigration(7, AddServerStatus, db)
	if err != nil {
		return err
	}

	return nil
}

func AddColumn(table string, name string, props string, defaultVal string, db *sqlx.DB) error {
	stmnt := fmt.Sprintf("ALTER TABLE %s ADD %s %s DEFAULT %s", table, name, props, defaultVal)
	_, err := db.Exec(stmnt)
	return err
}

func APICacheChanges(db *sqlx.DB) error {
	err := AddColumn("players", "oldgold", "BOOLEAN NOT NULL", "FALSE", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "registered", "VARCHAR(100) NOT NULL", "''", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "role", "INT NOT NULL", "0", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "avatar", "VARCHAR(255) NOT NULL", "''", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "tier", "INT NOT NULL", "0", db)
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE events MODIFY time BIGINT UNSIGNED NOT NULL")
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE events SET time=time*1000")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE kills MODIFY epoch BIGINT UNSIGNED NOT NULL")
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE kills SET epoch=epoch*1000")
	if err != nil {
		return err
	}

	return nil
}

func TeamKillChanges(db *sqlx.DB) error {
	_, err := db.Exec("UPDATE kills SET teamKill=0 WHERE killerID=victimID")
	return err
}

func BumpNameLimit(db *sqlx.DB) error {
	_, err := db.Exec("ALTER TABLE players MODIFY username varchar(255) NOT NULL UNIQUE")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE players MODIFY charactername varchar(255) NOT NULL")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE players MODIFY clantag varchar(255) NOT NULL")
	if err != nil {
		return err
	}

	return nil
}

func SawKillBuilderOnly(db *sqlx.DB) error {
	_, err := db.Exec("UPDATE kills SET killerClass='builder' WHERE hitter=30")
	return err
}

func AddLastEventToPlayers(db *sqlx.DB) error {
	err := AddColumn("players", "lastEventID", "INT", "NULL", db)
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE players ADD FOREIGN KEY(lastEventID) REFERENCES events(ID)")
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE players
		INNER JOIN events as event ON event.ID = (SELECT e.ID FROM events as e WHERE e.playerID=players.ID ORDER BY e.ID DESC LIMIT 1) 
		SET lastEventID=event.ID`)
	if err != nil {
		return err
	}

	return nil
}

func AddAccolades(db *sqlx.DB) error {
	err := AddColumn("players", "gold", "INT NOT NULL", "0", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "silver", "INT NOT NULL", "0", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "bronze", "INT NOT NULL", "0", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "participation", "INT NOT NULL", "0", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "github", "BOOLEAN NOT NULL", "FALSE", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "community", "BOOLEAN NOT NULL", "FALSE", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "mapmaker", "BOOLEAN NOT NULL", "FALSE", db)
	if err != nil {
		return err
	}

	err = AddColumn("players", "moderation", "BOOLEAN NOT NULL", "FALSE", db)
	if err != nil {
		return err
	}

	return nil
}

func AddServerStatus(db *sqlx.DB) error {
	err := AddColumn("servers", "status", "BOOLEAN NOT NULL", "FALSE", db)
	if err != nil {
		return err
	}

	return nil
}
