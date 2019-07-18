package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/rcon"
	"github.com/pkg/errors"
)

type Collector struct {
	config configs.ServerConfig
	logger *log.Logger
	server models.Server
}

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
		if err != nil {
			return err
		}
		UpdateJoinTime(player.ID, c.server.ID)
		c.logger.Printf("%s (%d) joined the game", player.Username, player.ID)
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

		}
		c.logger.Printf("%s (%d) left the game", player.Username, player.ID)
	}
	return nil
}

func (c *Collector) OnPlayerDie(m rcon.Message, r *rcon.Client) error {
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
		kill.Time = time.Now().Unix()
		kill.ServerID = c.server.ID
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
			if err != nil {
				return err
			}
			err = UpdateJoinTime(p.ID, c.server.ID)
			if err != nil {
				return err
			}
			c.logger.Printf("%s (%d) was in the game", p.Username, p.ID)
		}
	}
	return nil
}

func (c *Collector) ServerInfo(m rcon.Message, r *rcon.Client) error {
	var server models.Server
	err := json.Unmarshal([]byte(m.Args["object"]), &server)
	if err != nil {
		c.logger.Println(err)
		return nil
	}
	server.Tags = c.config.TagsString()

	err = UpdateServerInfo(&server)
	if err != nil {
		return errors.Wrap(err, "can't start collector without server info")
	}

	c.logger.Printf("%+v", server)

	c.server = server
	c.logger.SetPrefix(fmt.Sprintf("[%s] ", server.Name))

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

		client.RunScript("scripts/ServerInfo.as")
		client.HandleFunc("^ServerInfo (?P<object>.*)", collector.ServerInfo).RemoveTimestamp()

		err := client.Handle()
		if err != nil {
			logger.Println(err)
		}
	}
}
