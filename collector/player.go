package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/mattn/go-sqlite3"
	. "github.com/Harrison-Miller/kagstats/models"
)

type PlayerDatabase struct {
	db          *sql.DB
	playerCache map[string]*Player
	uncommited  []Kill
	Total       int
}

func CreatePlayerDatabase(connection string) (PlayerDatabase, error) {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return PlayerDatabase{}, nil
	}
	return PlayerDatabase{
		db,
		make(map[string]*Player),
		make([]Kill, 0, 100),
		0,
	}, nil
}

func (pdb *PlayerDatabase) Init() {
	_, err := pdb.db.Exec("CREATE DATABASE IF NOT EXISTS kagstats CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin")
	if err != nil {
		panic(err)
	}

	_, err = pdb.db.Exec("USE kagstats")
	if err != nil {
		panic(err)
	}

	_, err = pdb.db.Exec(`CREATE TABLE IF NOT EXISTS players (
		ID INT PRIMARY KEY AUTO_INCREMENT,
		username varchar(30) NOT NULL,
		charactername varchar(30) NOT NULL,
		clantag varchar(30) NOT NULL
	)`)
	if err != nil {
		panic(err)
	}

	_, err = pdb.db.Exec(`CREATE TABLE IF NOT EXISTS servers (
			ID INTEGER PRIMARY KEY AUTO_INCREMENT,
			name varchar(255),
			tags varchar(1000)
	)`)
	if err != nil {
		panic(err)
	}

	_, err = pdb.db.Exec(`CREATE TABLE IF NOT EXISTS kills (
			ID INTEGER PRIMARY KEY AUTO_INCREMENT,
			killerID INT,
			victimID INT NOT NULL,
			assistID INT,
			killerClass ENUM('archer', 'builder', 'knight', 'other', 'none') DEFAULT 'none',
			victimClass ENUM('archer', 'builder', 'knight', 'other') DEFAULT 'archer' NOT NULL,
			hitter INT DEFAULT 0,
			epoch INT NOT NULL,
			serverID INT NOT NULL,
			teamKill BOOLEAN DEFAULT false,
			FOREIGN KEY(killerID) REFERENCES players(ID),
			FOREIGN KEY(victimID) REFERENCES players(ID),
			FOREIGN KEY(assistID) REFERENCES players(ID),
			FOREIGN KEY(serverID) REFERENCES servers(ID)
	)`)
	if err != nil {
		panic(err)
	}
}

func (pdb *PlayerDatabase) UpdatePlayerInfo(p *Player, charactername string, clantag string) error {
	tx, err := pdb.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE kagstats.players SET charactername=?, clantag=? WHERE id=?", charactername, clantag, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	p.Charactername = charactername
	p.Clantag = clantag
	return nil
}

func (pdb *PlayerDatabase) GetOrCreatePlayer(username string, charactername string, clantag string) (*Player, error) {
	tx, err := pdb.db.Begin()
	if err != nil {
		return nil, err
	}

	var p Player
	row := tx.QueryRow("SELECT * FROM kagstats.players WHERE username=?", username)
	err = row.Scan(&p.ID, &p.Username, &p.Charactername, &p.Clantag)

	// at this point if we don't have an error we should have a valid id
	if err != nil {
		// player doesn't exist we need to insert them
		res, err := tx.Exec("INSERT INTO kagstats.players (username, charactername, clantag) VALUES(?, ?, ?)", username, charactername, clantag)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		ID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		p.ID = ID
		p.Username = username
		p.Charactername = charactername
		p.Clantag = clantag

		if err = tx.Commit(); err != nil {
			// the insert might not have succeeded
			return nil, err
		}

		return &p, nil
	} else if p.Charactername != charactername || p.Clantag != clantag {
		// update charactername or clantag
		_, err = tx.Exec("UPDATE kagstats.players SET charactername=?, clantag=? WHERE id=?", charactername, clantag, p.ID)
		if err != nil {
			tx.Rollback()
			// id is valid
			return &p, err
		}

		p.Charactername = charactername
		p.Clantag = clantag
	}

	if err = tx.Commit(); err != nil {
		// id is valid but charactername or clantag could be invalid
		return &p, err
	}

	return &p, nil
}

func (pdb *PlayerDatabase) GetOrUpdatePlayer(username string, charactername string, clantag string) (*Player, error) {
	if username == "" {
		return nil, fmt.Errorf("not a valid username")
	}

	// first check the player cache
	if p, ok := pdb.playerCache[username]; ok {
		if p.Charactername != charactername || p.Clantag != clantag {
			err := pdb.UpdatePlayerInfo(p, charactername, clantag)
			if err != nil {
				// unable to update player name but it was retrieved
				return p, err
			}
		}
		return p, nil
	}

	p, err := pdb.GetOrCreatePlayer(username, charactername, clantag)
	if p != nil {
		pdb.playerCache[username] = p
	}
	return p, err
}

func (pdb *PlayerDatabase) Commit() error {
	tx, err := pdb.db.Begin()
	if err != nil {
		return err
	}

	//https://stackoverflow.com/questions/21108084/how-to-insert-multiple-data-at-once
	values := []interface{}{}
	sqlStr := "INSERT INTO kagstats.kills (killerID, victimID, assistID, killerClass, victimClass, hitter, epoch, serverID, teamKill) VALUES "
	for _, v := range pdb.uncommited {
		sqlStr += "(?,?,?,?,?,?,?,?,?),"
		values = append(values, v.KillerID, v.VictimID, v.AssistID, v.KillerClass, v.VictimClass, v.Hitter, v.Time, v.ServerID, v.TeamKill)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	stmnt, err := tx.Prepare(sqlStr)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmnt.Close()

	_, err = stmnt.Exec(values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	pdb.uncommited = make([]Kill, 0, 100)
	return nil
}

func (pdb *PlayerDatabase) GetOrUpdateServer(name string, tags []string) (Server, error) {
	tagsStr := strings.Join(tags, ",")
	tx, err := pdb.db.Begin()
	if err != nil {
		return Server{}, err
	}

	var s Server
	row := tx.QueryRow("SELECT * FROM kagstats.servers WHERE name=?", name)
	err = row.Scan(&s.ID, &s.Name, &s.Tags)

	if err != nil {
		// time to create the new server
		res, err := tx.Exec("INSERT INTO kagstats.servers (name, tags) VALUES(?, ?)", name, tagsStr)
		if err != nil {
			tx.Rollback()
			return Server{}, err
		}

		ID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return Server{}, err
		}

		s.ID = ID
		s.Name = name
		s.Tags = tagsStr

		if err = tx.Commit(); err != nil {
			return Server{}, err
		}

		return s, nil
	}

	// check if we need to update the tags
	if s.Tags != tagsStr {
		_, err = tx.Exec("UPDATE kagstats.servers SET tags=?", tagsStr)
		if err != nil {
			tx.Rollback()
			return s, err
		}
		s.Tags = tagsStr
	}

	if err = tx.Commit(); err != nil {
		return s, err
	}

	return s, nil
}
