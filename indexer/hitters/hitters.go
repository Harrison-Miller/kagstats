package main

import (
	"log"

	. "github.com/Harrison-Miller/kagstats/indexer"
	. "github.com/Harrison-Miller/kagstats/models"
)

type HittersIndexer struct {
}

func (i *HittersIndexer) Name() string {
	return "top_hitters"
}

func (i *HittersIndexer) Version() int {
	return 1
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
	if kill.KillerID.Valid && kill.KillerID.Int64 != kill.VictimID {
		indices = append(indices, Index{
			Keys:     []int{int(kill.KillerID.Int64), int(kill.Hitter)},
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
