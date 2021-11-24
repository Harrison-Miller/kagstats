package main

import (
	"context"
	"github.com/Harrison-Miller/kagstats/collector/database"
	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/Harrison-Miller/kagstats/common/utils"
	"log"
	"time"
)

type PlayerInfoUpdater struct {
	db database.Database
	incoming chan models.Player
	ctx context.Context
	cancelFunc context.CancelFunc
}

func NewPlayerInfoUpdater(db database.Database) *PlayerInfoUpdater {
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
	return p.db.CommitPlayer(player)
}