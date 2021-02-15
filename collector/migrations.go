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

	err = RunMigration(8, AddLeaderboardBan, db)
	if err != nil {
		return err
	}

	err = RunMigration(9, AddStatsBan, db)
	if err != nil {
		return err
	}

	err = RunMigration(10, OnDeleteCascade, db)
	if err != nil {
		return err
	}

	err = RunMigration(11, RemoveAltAccounts, db)
	if err != nil {
		return err
	}

	err = RunMigration(12, AddNotes, db)
	if err != nil {
		return err
	}

	err = RunMigration(13, RemoveWeekOldAccounts, db)
	if err != nil {
		return err
	}

	err = RunMigration(14, AddFlagCaptures, db)
	if err != nil {
		return err
	}

	err = RunMigration(15, AddMapStats, db)
	if err != nil {
		return err
	}

	err = RunMigration(16, AddMapVotes, db)
	if err != nil {
		return err
	}

	err = RunMigration(17, RemoveEvents, db)
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

func DelColumn(table string, name string, db *sqlx.DB) error {
	stmnt := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table, name)
	_, err := db.Exec(stmnt)
	return err
}

func DelFK(table string, name string, db *sqlx.DB) error {
	stmnt := fmt.Sprintf("alter table %s drop foreign key %s", table, name)
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

func AddLeaderboardBan(db *sqlx.DB) error {
	err := AddColumn("players", "leaderboardBan", "BOOLEAN NOT NULL", "FALSE", db)
	if err != nil {
		return err
	}

	return nil
}

func AddStatsBan(db *sqlx.DB) error {
	err := AddColumn("players", "statsBan", "BOOLEAN NOT NULL", "FALSE", db)
	if err != nil {
		return err
	}

	return nil
}

func ModifyFkOnDelete(constraint string, table string, key string, reference string, db *sqlx.DB) error {
	_, err := db.Exec("ALTER TABLE " + table + " DROP FOREIGN KEY " + constraint)
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE " + table + " ADD FOREIGN KEY(" + key + ") REFERENCES " + reference + " ON DELETE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

func OnDeleteCascade(db *sqlx.DB) error {
	err := ModifyFkOnDelete("players_ibfk_1", "players", "lastEventID", "events(ID)", db)
	if err != nil {
		return err
	}

	err = ModifyFkOnDelete("kills_ibfk_1", "kills", "killerID", "players(ID)", db)
	if err != nil {
		return err
	}

	err = ModifyFkOnDelete("kills_ibfk_2", "kills", "victimID", "players(ID)", db)
	if err != nil {
		return err
	}

	err = ModifyFkOnDelete("kills_ibfk_3", "kills", "serverID", "servers(ID)", db)
	if err != nil {
		return err
	}

	err = ModifyFkOnDelete("events_ibfk_1", "events", "playerID", "players(ID)", db)
	if err != nil {
		return err
	}

	err = ModifyFkOnDelete("events_ibfk_2", "events", "serverID", "servers(ID)", db)
	if err != nil {
		return err
	}

	// known indexers
	err = ModifyFkOnDelete("basic_stats_ibfk_1", "basic_stats", "playerID", "players(ID)", db)
	if err != nil {
		return err
	}

	err = ModifyFkOnDelete("nemesis_ibfk_1", "nemesis", "playerID", "players(ID)", db)
	if err != nil {
		return err
	}

	err = ModifyFkOnDelete("nemesis_ibfk_2", "nemesis", "nemesisID", "players(ID)", db)
	if err != nil {
		return err
	}

	err = ModifyFkOnDelete("top_hitters_ibfk_1", "top_hitters", "playerID", "players(ID)", db)
	if err != nil {
		return err
	}

	return nil
}

func RemoveAltAccounts(db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM players WHERE username REGEXP '^.*~[0-9]+'")
	if err != nil {
		return err
	}

	return nil
}

func AddNotes(db *sqlx.DB) error {
	err := AddColumn("players", "notes", "varchar(255) NOT NULL", "''", db)
	if err != nil {
		return err
	}

	return nil
}

func RemoveWeekOldAccounts(db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM players WHERE registered = ''")
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM players WHERE registered BETWEEN NOW() - INTERVAL 7 DAY AND NOW()")
	if err != nil {
		return err
	}

	return nil
}

func AddFlagCaptures(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS flag_captures (
		ID INTEGER PRIMARY KEY AUTO_INCREMENT,
		playerID INT NOT NULL,
		ticks INT NOT NULL,
		FOREIGN KEY(playerID) REFERENCES players(ID) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}

	return nil
}

func AddMapStats(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS map_stats (
		ID INTEGER PRIMARY KEY AUTO_INCREMENT,
		mapName varchar(255) NOT NULL,
		ticks INT NOT NULL
	)`)
	if err != nil {
		return err
	}

	return nil
}

func AddMapVotes(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS map_votes (
		ID INTEGER PRIMARY KEY AUTO_INCREMENT,
		map1Name varchar(255) NOT NULL,
		map1Votes INT NOT NULL,
		map2Name varchar(255) NOT NULL,
		map2Votes INT NOT NULL,
		randomVotes INT NOT NULL
	)`)
	if err != nil {
		return err
	}

	return nil
}

func RemoveEvents(db *sqlx.DB) error {
	err := DelFK("players", "players_ibfk_1", db)
	if err != nil {
		return errd
	}

	err = DelColumn("players", "lastEventID", db)
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS events")
	if err != nil {
		return err
	}

	return nil
}
