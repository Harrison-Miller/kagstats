package main

import (
	"log"

	"github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/indexer"
)

type HittersIndexer struct {
}

func (i *HittersIndexer) Name() string {
	return "top_hitters"
}

func (i *HittersIndexer) Version() int {
	return 3
}

func (i *HittersIndexer) Keys() []IndexKey {
	var keys []IndexKey
	keys = append(keys, IndexKey{
		Name:   "playerID",
		Table:  "players",
		Column: "ID",
	}, IndexKey{
		Name: "hitter",
	})
	return keys
}

func (i *HittersIndexer) Counters() []string {
	return []string{"kills"}
}

func (i *HittersIndexer) Index(kill Kill) []Index {
	var indices []Index
	if kill.KillerID != kill.VictimID && !kill.TeamKill {
		// remove duplicate hitters
		if kill.Hitter == models.Burn {
			kill.Hitter = models.Fire
		} else if kill.Hitter == models.Mine_special {
			kill.Hitter = models.Mine_special
		} else if kill.Hitter == models.Cata_boulder {
			kill.Hitter = models.Cata_stones
		}

		indices = append(indices, Index{
			Keys:     []interface{}{int(kill.KillerID), int(kill.Hitter)},
			Counters: map[string]int{"kills": 1},
		})
	}

	return indices
}

func main() {
	indexer := HittersIndexer{}
	err := Run(&indexer)
	if err != nil {
		log.Fatal(err)
	}
}
