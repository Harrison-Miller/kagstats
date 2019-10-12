package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Harrison-Miller/kagstats/common/utils"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/rcon"
	"github.com/pkg/errors"
)

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

	if (player != models.Player{}) {
		err = UpdatePlayer(&player)
		player.ServerID = c.server.ID
		if err != nil {
			return err
		}
		err = UpdateJoinTime(player.ID, c.server.ID)
		if err != nil {
			return err
		}
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

	if (player != models.Player{}) {
		err = UpdatePlayer(&player)
		if err != nil {
			return err
		}
		err = UpdateLeaveTime(player.ID, c.server.ID)
		if err != nil {
			return err
		}
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

	if (kill != models.Kill{}) {
		err = UpdatePlayer(&kill.Player)
		if err != nil {
			return errors.Wrap(err, "error getting victim id")
		}
		UpdatePlayer(&kill.Killer)
		if err != nil {
			return errors.Wrap(err, "error getting killer id")
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
		for _, p := range players {
			err = UpdatePlayer(&p)
			p.ServerID = c.server.ID
			if err != nil {
				return err
			}
			err = UpdateJoinTime(p.ID, c.server.ID)
			if err != nil {
				return err
			}
			c.logger.Printf("%s (%d) was in the game", p.Username, p.ID)
		}
		c.playerCount = len(players)
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

func AddHandlers(client *rcon.Client, collector *Collector) {
	client.HandleFunc("^PlayerJoined (?P<object>.*)", collector.OnPlayerJoined).RemoveTimestamp()
	client.HandleFunc("^PlayerLeft (?P<object>.*)", collector.OnPlayerLeave).RemoveTimestamp()
	client.HandleFunc("^PlayerDied (?P<object>.*)", collector.OnPlayerDie).RemoveTimestamp()
	client.HandleFunc("^PlayerList (?P<object>.*)", collector.PlayerList).RemoveTimestamp()
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

		defer client.Close()
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

		err := client.Handle()
		stopMOTD <- true
		if err != nil {
			logger.Println(err)
		}

		// Add a disconnect event if the server disconects from the collector
		for _, player := range players {
			if player.ServerID == collector.server.ID {
				err = UpdateLeaveTime(player.ID, player.ServerID)
				if err != nil {
					logger.Println(err)
				}
				logger.Printf("%s (%d) left the game", player.Username, player.ID)
			}
		}

		logger.Println("Waiting 1 minute before attempting to connect again")
		time.Sleep(1 * time.Minute)
	}
}
