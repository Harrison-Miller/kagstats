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
	server models.Server
	logger *log.Logger
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
	}

	if (player != models.Player{}) {
		err = UpdatePlayer(&player)
		if err != nil {
			return err
		}
		UpdateJoinTime(player)
		c.logger.Printf("%s (%d) joined the game", player.Username, player.ID)
	}
	return nil
}

func (c *Collector) OnPlayerLeave(m rcon.Message, r *rcon.Client) error {
	var player models.Player
	err := json.Unmarshal([]byte(m.Args["object"]), &player)
	if err != nil {
		c.logger.Println(err)
	}

	if (player != models.Player{}) {
		err = UpdatePlayer(&player)
		if err != nil {
			return err
		}
		err = UpdateLeaveTime(player)
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
	}

	if len(players) > 0 {
		for _, p := range players {
			err = UpdatePlayer(&p)
			if err != nil {
				return err
			}
			err = UpdateJoinTime(p)
			if err != nil {
				return err
			}
			c.logger.Printf("%s (%d) was in the game", p.Username, p.ID)
		}
	}
	return nil
}

func AddHandlers(client *rcon.Client, collector *Collector) {
	client.HandleFunc("^PlayerJoined (?P<object>.*)", collector.OnPlayerJoined).RemoveTimestamp()
	client.HandleFunc("^PlayerLeft (?P<object>.*)", collector.OnPlayerLeave).RemoveTimestamp()
	client.HandleFunc("^PlayerDied (?P<object>.*)", collector.OnPlayerDie).RemoveTimestamp()
	client.HandleFunc("^PlayerList (?P<object>.*)", collector.PlayerList).RemoveTimestamp()
}

func Collect(sconfig configs.ServerConfig) {
	logger := log.New(os.Stdout, fmt.Sprintf("[%s] ", sconfig.Name), log.LstdFlags)
	address := net.JoinHostPort(sconfig.Address, sconfig.Port)

	serverInfo, err := UpdateServerInfo(sconfig)
	if err != nil {
		log.Println("can't get server ID not starting collector", err)
		return
	}
	collector := Collector{sconfig, serverInfo, logger}

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

		err := client.RunScript("scripts/PlayerList.as")
		if err != nil {
			logger.Println("wat", err)
		}

		AddHandlers(&client, &collector)

		err = client.Handle()
		if err != nil {
			logger.Println(err)
		}
	}
}
