package main

import (
	"context"
	"fmt"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/kagstats/common/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

type PlayerInfoUpdater struct {
	db *sqlx.DB
	incoming chan models.Player
	ctx context.Context
	cancelFunc context.CancelFunc
}

func NewPlayerInfoUpdater(db *sqlx.DB) *PlayerInfoUpdater {
	p := &PlayerInfoUpdater{
		db: db,
		incoming: make(chan models.Player, 10240),
	}
	p.ctx, p.cancelFunc = context.WithCancel(context.Background())
	go p.process()
	return p
}

func (p *PlayerInfoUpdater) close() {
	p.cancelFunc()
}

func (p *PlayerInfoUpdater) process() {
	for {
		select {
		case <- p.ctx.Done():
			return
		case player := <- p.incoming:
			err := utils.GetPlayerAvatar(&player)
			if err != nil {
				log.Println("err getting avatar: ", err)
			}

			err = utils.GetPlayerTier(&player)
			if err != nil {
				log.Println("err getting tier: ", err)
			}

			err = utils.GetPlayerInfo(&player)
			if err != nil {
				log.Println(err)
			}

			err = p.commitPlayer(&player)
			if err != nil {
				log.Println("err committing player: ", err)
			}

			time.Sleep(time.Second * 1)
		}
	}
}

func (p *PlayerInfoUpdater) commitPlayer(player *models.Player) error {
	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	defer tx.Rollback()


	_, err = db.Exec("UPDATE players SET oldgold=?,registered=?,role=?,avatar=?,tier=? WHERE ID=?",
		player.OldGold, player.Registered, player.Role, player.Avatar, player.Tier, player.ID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error updating player %s with api info", player.Username))
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error commiting updated player info")
	}
	return nil
}