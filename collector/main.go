package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
)

var address = "localhost:50301"
var rconPassword = "admin1234!"
var maxUncommited = 3

const defaultTimeout = 15 * time.Second

var timeout = defaultTimeout

type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/", c.User, c.Password, c.Address)
}

type ServerConfig struct {
	Name          string   `json:"name"`
	Tags          []string `json:"tags"`
	Address       string   `json:"address"`
	Password      string   `json:"password"`
	Connected     bool
	ConnectedTime time.Time
}

type Config struct {
	Name         string         `json:"name"`
	Timeout      string         `json:"timeout"`
	MaxAttempts  int            `json:"maxAttempts"`
	BulkLoadMax  int            `json:"bulkLoadMax"`
	BulkLoadWait string         `json:"bulkLoadWait"`
	Database     DatabaseConfig `json:"database"`
	Servers      []ServerConfig `json:"servers"`
	Monitoring   bool           `json:monitoring`
}

type KillMessage struct {
	VictimUsername      string
	VictimCharacterName string
	VictimClantag       string
	VictimClass         string
	Hitter              int64

	KillerUsername      string
	KillerCharacterName string
	KillerClantag       string
	KillerClass         string
}

func process(pdb *PlayerDatabase, message string) (KillEntry, error) {
	parts := strings.SplitN(message, " ", 3)
	jsonStr := parts[2]

	var msg KillMessage
	err := json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		return KillEntry{}, err
	}

	victim, err := pdb.GetOrUpdatePlayer(msg.VictimUsername, msg.VictimCharacterName, msg.VictimClantag)
	if err != nil {
		return KillEntry{}, err
	}

	var killerID sql.NullInt64
	if msg.KillerUsername != "" {
		killer, err := pdb.GetOrUpdatePlayer(msg.KillerUsername, msg.KillerCharacterName, msg.KillerClantag)
		if err != nil {
			return KillEntry{}, err
		}

		if killer != nil {
			killerID.Int64 = killer.ID
			killerID.Valid = true
		}
	}

	if msg.VictimClass != "archer" && msg.VictimClass != "builder" && msg.VictimClass != "knight" {
		msg.VictimClass = "other"
	}

	if msg.KillerUsername == "" {
		msg.KillerClass = "none"
	} else if msg.KillerClass != "archer" && msg.KillerClass != "builder" && msg.KillerClass != "knight" {
		msg.KillerClass = "other"
	}

	// TODO: submit the kill entry into a channel
	kill := KillEntry{
		killerID,
		victim.ID,
		sql.NullInt64{},
		msg.KillerClass,
		msg.VictimClass,
		msg.Hitter,
		time.Now().Unix(),
		0,
	}

	return kill, nil

}

func RCONLogin(address string, password string, timeout time.Duration, maxAttempts int) (net.Conn, error) {
	var conn net.Conn
	var attempts int

	for ok := true; ok; {
		c, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			attempts++
			if attempts >= maxAttempts {
				return conn, err
			}
		} else {
			fmt.Fprintf(c, "%s\n", password)
			_, err = bufio.NewReader(c).ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("Wrong password")
			}
			conn = c
			ok = false
		}
	}

	return conn, nil
}

func Collect(pdb *PlayerDatabase, config *Config, index int, entires chan KillEntry) {
	sconfig := &config.Servers[index]
	info, err := pdb.GetOrUpdateServer(sconfig.Name, sconfig.Tags)
	if err != nil {
		log.Printf("[%d] %v", err)
	}

	fmt.Printf("[%d] Connecting to %s...\n", index, sconfig.Address)
	conn, err := RCONLogin(sconfig.Address, sconfig.Password, timeout, config.MaxAttempts)
	if err != nil {
		log.Printf("[%d] %v", index, err)
		return
	}
	fmt.Printf("[%d] Authenticated with %s!\n", index, sconfig.Address)
	sconfig.Connected = true
	sconfig.ConnectedTime = time.Now()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			sconfig.Connected = false
			sconfig.ConnectedTime = time.Now()
			fmt.Printf("[%d] Disconnected from %s...\n", index, sconfig.Address)
			conn, err = RCONLogin(address, rconPassword, 15*time.Second, 20)
			if err != nil {
				fmt.Printf("[%d] Failed to reconnect to %s...\n", index, sconfig.Address)
				log.Println(err)
				return
			}
			fmt.Printf("[%d] Reconnected to %s!\n", index, sconfig.Address)
			sconfig.Connected = true
		} else if strings.Contains(message, "*STATS") {
			kill, err := process(pdb, message)
			if err != nil {
				log.Printf("[%d] %v", index, err)
			} else {
				kill.serverID = info.ID
				entires <- kill
			}
		}
	}

}

func main() {
	// Read in the configuration file
	file, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		log.Fatal(err)
	}

	timeout, err = time.ParseDuration(config.Timeout)
	if err != nil {
		timeout = defaultTimeout
	}

	fmt.Printf("Connecting to the database at %s...\n", config.Database.Address)
	pdb, err := CreatePlayerDatabase(config.Database.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")
	pdb.Init()

	startTime := time.Now()
	if config.Monitoring {
		go startMonitoringServer(&config, &pdb, startTime)
	}

	// Start each collector
	entries := make(chan KillEntry)
	for i := range config.Servers {
		go Collect(&pdb, &config, i, entries)
		time.Sleep(1 * time.Second)
	}

	// collate a the kill entries from the channel here\
	for kill := range entries {
		pdb.uncommited = append(pdb.uncommited, kill)
		if len(pdb.uncommited) > maxUncommited {
			pdb.Total += len(pdb.uncommited)
			err = pdb.Commit()
			if err != nil {
				panic(err)
			}
		}
	}

}
