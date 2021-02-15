package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Harrison-Miller/kagstats/common/utils"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/rcon"
	"github.com/pkg/errors"
)

// Filter multiple clients of the same account
var clientAltPattern = regexp.MustCompile("[^ ]+~[0-9]+")

func isClientAlt(username string) bool {
	if clientAltPattern.MatchString(username) {
		return true
	}

	return false
}

func isNewPlayer(player *models.Player) bool {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, player.Registered)
	if err != nil {
		return true
	}

	dur := time.Now().Sub(t)
	days := int64(dur.Hours()) / 24

	if days <= 7 {
		return true
	}

	return false
}

type Collector struct {
	config      configs.ServerConfig
	logger      *log.Logger
	server      models.Server
	playerCount int
}

const TDM_MINIMUM_PLAYERS = 4
const CTF_MINIMUM_PLAYERS = 8

func UpdatePlayer(p *models.Player) error {
	if cache, ok := players[p.Username]; ok {
		p.ID = cache.ID
		p.StatsBan = cache.StatsBan
		cache.ServerID = p.ServerID
		p.Registered = cache.Registered

		if p.Charactername != cache.Charactername || p.Clantag != cache.Clantag {
			cache.Charactername = p.Charactername
			cache.Clantag = p.Clantag

			err := UpdatePlayerInfo(p)
			if err != nil {
				return errors.Wrap(err, "error updating player info")
			}
		}
	} else {
		err := UpdatePlayerInfo(p)
		if err != nil {
			return errors.Wrap(err, "error creating player")
		}

		err = utils.GetPlayerInfo(p)
		if err != nil {
			return errors.Wrap(err, "error getting player info")
		}

		players[p.Username] = *p
	}

	return nil
}

func (c *Collector) OnPlayerJoined(m rcon.Message, r *rcon.Client) error {
	var player models.Player
	err := json.Unmarshal([]byte(m.Args["object"]), &player)
	if err != nil {
		c.logger.Println(err)
		return nil
	}

	if isClientAlt(player.Username) {
		return nil
	}

	if (player != models.Player{}) {
		player.ServerID = c.server.ID
		err = UpdatePlayer(&player)
		if err != nil {
			return err
		}

		if player.StatsBan || isNewPlayer(&player) {
			return nil
		}

		/*
			err = UpdateJoinTime(player.ID, c.server.ID)
			if err != nil {
				return err
			}
		*/
		c.logger.Printf("%s (%d) joined the game", player.Username, player.ID)
		c.playerCount++
	}
	return nil
}

func (c *Collector) OnPlayerLeave(m rcon.Message, r *rcon.Client) error {
	var player models.Player
	err := json.Unmarshal([]byte(m.Args["object"]), &player)
	if err != nil {
		c.logger.Println(err)
		return nil
	}

	if isClientAlt(player.Username) {
		return nil
	}

	if (player != models.Player{}) {
		err = UpdatePlayer(&player)
		if err != nil {
			return err
		}

		if player.StatsBan || isNewPlayer(&player) {
			return nil
		}

		/*
			err = UpdateLeaveTime(player.ID, c.server.ID)
			if err != nil {
				return err
			}
		*/
		c.logger.Printf("%s (%d) left the game", player.Username, player.ID)
		if c.playerCount > 0 {
			c.playerCount--
		}
	}
	return nil
}

func (c *Collector) OnPlayerDie(m rcon.Message, r *rcon.Client) error {
	if c.server.Gamemode == "TDM" && c.playerCount < TDM_MINIMUM_PLAYERS {
		c.logger.Println("not enough players, not adding kill to db")
		return nil
	}

	if (c.server.Gamemode == "CTF" || c.server.Gamemode == "WAR") && c.playerCount < CTF_MINIMUM_PLAYERS {
		c.logger.Println("not enough players, not adding kill to db")
		return nil
	}

	var kill models.Kill
	err := json.Unmarshal([]byte(m.Args["object"]), &kill)
	if err != nil {
		c.logger.Println(err)
		return nil
	}

	if isClientAlt(kill.Killer.Username) || isClientAlt(kill.Player.Username) {
		return nil
	}

	if (kill != models.Kill{}) {
		err = UpdatePlayer(&kill.Player)
		if err != nil {
			return errors.Wrap(err, "error getting victim id")
		}
		UpdatePlayer(&kill.Killer)
		if err != nil {
			return errors.Wrap(err, "error getting killer id")
		}

		if kill.Player.StatsBan || kill.Killer.StatsBan || isNewPlayer(&kill.Player) || isNewPlayer(&kill.Killer) {
			return nil
		}

		kill.VictimID = kill.Player.ID
		kill.KillerID = kill.Killer.ID
		kill.Time = utils.NowAsUnixMilliseconds()
		kill.ServerID = c.server.ID

		// players can't team kill themselves
		if kill.KillerID == kill.VictimID {
			kill.TeamKill = false
		} else if kill.Hitter == 29 || kill.Hitter == 30 {
			// only builders can get spike/saw kills
			kill.KillerClass = "builder"
		}

		c.logger.Printf("%+v", kill)
		kills <- kill
	}
	return nil
}

func (c *Collector) PlayerList(m rcon.Message, r *rcon.Client) error {
	var players []models.Player
	err := json.Unmarshal([]byte(m.Args["object"]), &players)
	if err != nil {
		c.logger.Println(err)
		return nil
	}

	if len(players) > 0 {
		altPlayers := 0
		for _, p := range players {
			p.ServerID = c.server.ID

			if isClientAlt(p.Username) {
				altPlayers++
				continue
			}

			err = UpdatePlayer(&p)
			if err != nil {
				return err
			}

			if p.StatsBan || isNewPlayer(&p) {
				altPlayers++
				continue
			}

			/*
				err = UpdateJoinTime(p.ID, c.server.ID)
				if err != nil {
					return err
				}
			*/
			c.logger.Printf("%s (%d) was in the game", p.Username, p.ID)
		}
		c.playerCount = len(players) - altPlayers
	}

	r.RemoveHandler("^PlayerList (?P<object>.*)")
	return nil
}

func (c *Collector) ServerInfo(m rcon.Message, r *rcon.Client) error {
	var server models.Server
	err := json.Unmarshal([]byte(m.Args["object"]), &server)
	if err != nil {
		c.logger.Println(err)
		return nil
	}
	server.Address = c.config.Address
	server.Port = c.config.Port
	server.Tags = c.config.TagsString()

	err = UpdateServerInfo(&server)
	if err != nil {
		return errors.Wrap(err, "can't start collector without server info")
	}

	c.logger.Printf("%+v", server)

	c.server = server
	c.logger.SetPrefix(fmt.Sprintf("[%s] ", server.Name))

	r.RemoveHandler("^ServerInfo (?P<object>.*)")

	r.RunScript("scripts/AddMod.as")
	r.RunScript("scripts/PlayerList.as")
	AddHandlers(r, c)

	return nil
}

func (c *Collector) FlagCaptured(m rcon.Message, r *rcon.Client) error {
	var capture models.FlagCapture
	err := json.Unmarshal([]byte(m.Args["object"]), &capture)
	if err != nil {
		c.logger.Println(err)
		return nil
	}

	if (c.server.Gamemode == "CTF") && c.playerCount < CTF_MINIMUM_PLAYERS {
		c.logger.Println("not enough players, not flag capture to db")
		return nil
	}

	if isClientAlt(capture.Player) {
		return nil
	}

	if (capture != models.FlagCapture{}) {
		var player models.Player
		if cache, ok := players[capture.Player]; ok {
			player.ID = cache.ID
			player.StatsBan = cache.StatsBan
			player.Registered = cache.Registered
		}

		if player.StatsBan || isNewPlayer(&player) {
			return nil
		}

		capture.PlayerID = player.ID

		err = CommitFlagCapture(capture)
		if err != nil {
			return errors.Wrap(err, "can't commit flag capture")
		}

		c.logger.Printf("%+v", capture)
	}

	return nil
}

func (c *Collector) MapStats(m rcon.Message, r *rcon.Client) error {
	var stats models.MapStats
	err := json.Unmarshal([]byte(m.Args["object"]), &stats)
	if err != nil {
		c.logger.Println(err)
		return nil
	}

	// For now lets just ignore TDM since game duration isn't to important
	if c.server.Gamemode == "TDM" {
		return nil
	}
	// if c.server.Gamemode == "TDM" && c.playerCount < TDM_MINIMUM_PLAYERS {
	// 	c.logger.Println("not enough players, not adding kill to db")
	// 	return nil
	// }

	if (c.server.Gamemode == "CTF" || c.server.Gamemode == "WAR") && c.playerCount < CTF_MINIMUM_PLAYERS {
		c.logger.Println("not enough players, not adding kill to db")
		return nil
	}

	c.logger.Printf("%+v", stats)

	if (stats != models.MapStats{}) {
		err = CommitMapStats(stats)
		if err != nil {
			return errors.Wrap(err, "can't commit map stats")
		}
	}

	return nil
}

func (c *Collector) MapVotes(m rcon.Message, r *rcon.Client) error {
	var votes models.MapVotes
	votes.Map1Name = m.Args["map1_name"]
	votes.Map1Votes, _ = strconv.ParseInt(m.Args["map1_votes"], 10, 64)
	votes.Map2Name = m.Args["map2_name"]
	votes.Map2Votes, _ = strconv.ParseInt(m.Args["map2_votes"], 10, 64)
	votes.RandomVotes, _ = strconv.ParseInt(m.Args["random_votes"], 10, 64)

	if (c.server.Gamemode == "CTF" || c.server.Gamemode == "WAR") && c.playerCount < CTF_MINIMUM_PLAYERS {
		c.logger.Println("not enough players, not adding kill to db")
		return nil
	}

	err := CommitMapVotes(votes)
	if err != nil {
		return errors.Wrap(err, "can't commit map votes")
	}

	return nil
}

func AddHandlers(client *rcon.Client, collector *Collector) {
	client.HandleFunc("^PlayerJoined (?P<object>.*)", collector.OnPlayerJoined).RemoveTimestamp()
	client.HandleFunc("^PlayerLeft (?P<object>.*)", collector.OnPlayerLeave).RemoveTimestamp()
	client.HandleFunc("^PlayerDied (?P<object>.*)", collector.OnPlayerDie).RemoveTimestamp()
	client.HandleFunc("^PlayerList (?P<object>.*)", collector.PlayerList).RemoveTimestamp()
	client.HandleFunc("^FlagCaptured (?P<object>.*)", collector.FlagCaptured).RemoveTimestamp()
	client.HandleFunc("^MapStats (?P<object>.*)", collector.MapStats).RemoveTimestamp()
	client.HandleFunc("^\\(MapVotes\\) Map1: (?P<map1_name>.*) = (?P<map1_votes>\\d+) Map2: (?P<map2_name>.*) = (?P<map2_votes>\\d+) Random = (?P<random_votes>\\d+)", collector.MapVotes).RemoveTimestamp()
}

func Collect(sconfig configs.ServerConfig) {
	if addrs, err := net.LookupHost(sconfig.Address); err == nil {
		sconfig.Address = addrs[0]
	}

	address := net.JoinHostPort(sconfig.Address, sconfig.Port)
	logger := log.New(os.Stdout, fmt.Sprintf("[%s] ", address), log.LstdFlags)

	collector := Collector{
		config: sconfig,
		logger: logger,
	}

	for {
		delay := 1
		var client rcon.Client
		for {
			var err error
			client, err = rcon.DialRcon(address, sconfig.Password, 2*time.Second)
			if err == nil {
				break
			}
			logger.Println(errors.Wrap(err, "error connecting to "+address))
			logger.Printf("waiting %d seconds before next attempt\n", delay)
			time.Sleep(time.Duration(delay) * time.Second)
			delay = delay * 2
			if delay > 60 {
				delay = 60
			}
		}

		logger.Printf("Connected to %s!\n", address)

		collector.playerCount = 0

		client.RunScript("scripts/ServerInfo.as")
		client.HandleFunc("^ServerInfo (?P<object>.*)", collector.ServerInfo).RemoveTimestamp()

		stopMOTD := make(chan bool)

		go func() {
			if config.MOTD == "" {
				return
			}
			client.Message(config.MOTD)
			for {
				select {
				case <-stopMOTD:
					return
				case <-time.After(config.MOTDIntervalDuration):
					client.Message(config.MOTD)
				}
			}
		}()

		// Reset connection because of rcon server hiccups
		// go func() {
		// 	time.Sleep(2 * time.Hour)
		// 	logger.Printf("Disconnecting to refresh connection")
		// 	client.Close()
		// }()

		err := client.Handle()
		stopMOTD <- true
		client.Close()
		if err != nil {
			logger.Println(err)
		}

		// Add a disconnect event if the server disconects from the collector
		for _, player := range players {
			if player.ServerID == collector.server.ID {
				/*
					err = UpdateLeaveTime(player.ID, player.ServerID)
					if err != nil {
						logger.Println(err)
					}
				*/
				logger.Printf("%s (%d) left the game", player.Username, player.ID)
			}
		}

		UpdateServerStatus(collector.server, false)

		logger.Println("Waiting 1 minute before attempting to connect again")
		time.Sleep(1 * time.Minute)
	}
}
