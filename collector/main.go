package main

import (
	"bufio"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

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
	Port          int      `json:"port"`
	Password      string   `json:"password"`
	Connected     bool
	ConnectedTime time.Time
}

func (sc *ServerConfig) GetID() string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
		return sc.Name
	}
	id := reg.ReplaceAllString(sc.Name, "")
	return id
}

type Config struct {
	Name           string         `json:"name"`
	Timeout        string         `json:"timeout"`
	MaxAttempts    int            `json:"maxAttempts"`
	BulkLoadMax    int            `json:"bulkLoadMax"`
	BulkLoadWait   string         `json:"bulkLoadWait"`
	Database       DatabaseConfig `json:"database"`
	Servers        []ServerConfig `json:"servers"`
	Monitoring     bool           `json:"monitoring"`
	RefreshRate    int            `json:"refreshRate"`
	MonitoringPort int            `json:"monitoringPort"`
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
	TeamKill            bool
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
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
		msg.TeamKill,
	}

	return kill, nil

}

func RCONLogin(address string, password string, timeout time.Duration, maxAttempts int, context string) (net.Conn, error) {
	var conn net.Conn
	var attempts int
	var shortCircuit = 1

	for ok := true; ok; {
		logOut(context, fmt.Sprintf("Attempt %d/%d...", attempts+1, maxAttempts))
		c, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			attempts++
			if attempts >= maxAttempts {
				return conn, err
			}

			logOut(context, fmt.Sprintf("%v", err))

			logOut(context, fmt.Sprintf("Waiting %d seconds before next attempt", shortCircuit))
			time.Sleep(time.Duration(shortCircuit) * time.Second)
			if shortCircuit < 30 {
				shortCircuit = shortCircuit * 2
			}
		} else {
			fmt.Fprintf(c, "%s\n", password)
			_, err = bufio.NewReader(c).ReadString('\n')
			if err != nil {
				c.Close()
				logOut(context, fmt.Sprintf("%v", err))
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
		logOut(sconfig.Name, fmt.Sprintf("%v", err))
	}

	address := fmt.Sprintf("%s:%d", sconfig.Address, sconfig.Port)

	logOut(sconfig.Name, fmt.Sprintf("Connecting to %s...", address))
	conn, err := RCONLogin(address, sconfig.Password, timeout, config.MaxAttempts, sconfig.Name)
	if err != nil {
		logOut(sconfig.Name, fmt.Sprintf("%v", err))
		return
	}
	logOut(sconfig.Name, fmt.Sprintf("Authenticated with %s!", address))
	sconfig.Connected = true
	sconfig.ConnectedTime = time.Now()

	fmt.Fprintf(conn, "CRules@ r = getRules();string myId = '%s';if(!r.exists('stats_lock')) {r.set_s32('stats_lock', getGameTime() + getTicksASecond() * 10);r.set_string('stats_lock_holder', myId);} else {s32 timeleft = r.get_s32('stats_lock') - getGameTime();if(timeleft <= 0) {r.set_s32('stats_lock', getGameTime() + getTicksASecond() * 10);r.set_string('stats_lock_holder', myId);}}\n", collectorID)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			sconfig.Connected = false
			sconfig.ConnectedTime = time.Now()
			logOut(sconfig.Name, fmt.Sprintf("Disconnected from %s...", address))
			conn, err = RCONLogin(address, sconfig.Password, timeout, config.MaxAttempts, sconfig.Name)
			if err != nil {
				logOut(sconfig.Name, fmt.Sprintf("Failed to reconnect to %s...", address))
				logOut(sconfig.Name, fmt.Sprintf("%v", err))
				return
			}
			logOut(sconfig.Name, fmt.Sprintf("Reconnected to %s!\n", address))
			sconfig.Connected = true
		} else if strings.Contains(message, "*STATS") {
			kill, err := process(pdb, message)
			logOut(sconfig.Name, message)
			if err != nil {
				logOut(sconfig.Name, fmt.Sprintf("%v", err))
			} else {
				kill.serverID = info.ID
				entires <- kill
			}
		}
	}

}

var logFilePath string
var logFile *os.File
var collectorID string

func logOut(context string, text string) {
	logMsg := fmt.Sprintf("[%s] [%s] %s\n", time.Now().Format("15:04:05-01-02-2006"), context, text)
	fmt.Print(logMsg)
	logFile.WriteString(logMsg)
	logFile.Sync()
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func main() {
	// Read in the configuration file
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

	timeout, err = time.ParseDuration(config.Timeout)
	if err != nil {
		timeout = defaultTimeout
	}

	CreateDirIfNotExist("logs")
	logFilePath = fmt.Sprintf("logs/collector-%s.log", time.Now().Format("15-04-05-01-02-2006"))
	fmt.Println(logFilePath)
	logFile, err = os.Create(logFilePath)
	if err != nil {
		log.Fatal(err)
	}

	collectorID, _ = randomHex(16)
	logOut("main", fmt.Sprintf("Collector ID: %s", collectorID))

	logOut("main", fmt.Sprintf("Connecting to the database at %s...", config.Database.Address))
	pdb, err := CreatePlayerDatabase(config.Database.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	logOut("main", "Connected to the database!")
	pdb.Init()

	startTime := time.Now()
	if config.Monitoring {
		logOut("main", "Starting monitoring webserver...")
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
		if len(pdb.uncommited) >= config.BulkLoadMax {
			logOut("main", fmt.Sprintf("Commit %d new entries", len(pdb.uncommited)))
			pdb.Total += len(pdb.uncommited)
			err = pdb.Commit()
			if err != nil {
				panic(err)
			}
		}
	}

}
