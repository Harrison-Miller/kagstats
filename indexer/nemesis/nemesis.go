package main

import (
	"log"

	. "github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/indexer"
)

type NemesisIndexer struct {
}

func (i *NemesisIndexer) Name() string {
	return "nemesis"
}

func (i *NemesisIndexer) Version() int {
	return 1
}

func (i *NemesisIndexer) Keys() []IndexKey {
	var keys []IndexKey
	keys = append(keys, IndexKey{
		Name:   "playerID",
		Table:  "players",
		Column: "ID",
	}, IndexKey{
		Name:   "nemesisID",
		Table:  "players",
		Column: "ID",
	})
	return keys
}

func (i *NemesisIndexer) Counters() []string {
	return []string{"deaths"}
}

func (i *NemesisIndexer) Index(kill Kill) []Index {
	var indices []Index

	if kill.KillerID.Valid && kill.KillerID.Int64 != kill.VictimID {
		indices = append(indices, Index{
			Keys:     []int{int(kill.VictimID), int(kill.KillerID.Int64)},
			Counters: map[string]int{"deaths": 1},
		})
	}

	return indices
}

func main() {
	indexer := NemesisIndexer{}
	err := Run(&indexer)
	if err != nil {
		log.Fatal(err)
	}
}
