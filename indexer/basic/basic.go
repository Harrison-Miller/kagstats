package main

import (
	"log"

	. "github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/indexer"
)

type BasicIndexer struct {
}

func (i *BasicIndexer) Name() string {
	return "basic_stats"
}

func (i *BasicIndexer) Version() int {
	return 5
}

func (i *BasicIndexer) Keys() []IndexKey {
	var keys []IndexKey
	keys = append(keys, IndexKey{
		Name:   "playerID",
		Table:  "players",
		Column: "ID",
	})
	return keys
}

func (i *BasicIndexer) Counters() []string {
	return []string{"suicides", "teamkills", "archer_kills", "archer_deaths",
		"builder_kills", "builder_deaths", "knight_kills", "knight_deaths",
		"other_kills", "other_deaths", "total_kills", "total_deaths"}
}

func OneIfEqual(a string, b string) int {
	if a == b {
		return 1
	}
	return 0
}

func ToInt(a bool) int {
	if a {
		return 1
	}

	return 0
}

func OneIfNotIn(a string, b []string) int {
	for _, v := range b {
		if a == v {
			return 0
		}
	}

	return 1
}

func (i *BasicIndexer) Index(kill Kill) []Index {
	var indices []Index

	if kill.KillerID != kill.VictimID {
		if kill.TeamKill {
			indices = append(indices, Index{
				Keys:     []int{int(kill.KillerID)},
				Counters: map[string]int{"teamkills": 1},
			})
		} else {
			indices = append(indices, Index{
				Keys: []int{int(kill.KillerID)},
				Counters: map[string]int{
					"archer_kills":  OneIfEqual(kill.KillerClass, "archer"),
					"builder_kills": OneIfEqual(kill.KillerClass, "builder"),
					"knight_kills":  OneIfEqual(kill.KillerClass, "knight"),
					"other_kills":   OneIfNotIn(kill.KillerClass, []string{"archer", "builder", "knight"}),
					"total_kills":   1,
				},
			})
		}
	}

	if !kill.TeamKill {
		indices = append(indices, Index{
			Keys: []int{int(kill.VictimID)},
			Counters: map[string]int{
				"suicides":       ToInt(kill.KillerID == kill.VictimID),
				"archer_deaths":  OneIfEqual(kill.VictimClass, "archer"),
				"builder_deaths": OneIfEqual(kill.VictimClass, "builder"),
				"knight_deaths":  OneIfEqual(kill.VictimClass, "knight"),
				"other_deaths":   OneIfNotIn(kill.VictimClass, []string{"archer", "builder", "knight"}),
				"total_deaths":   1,
			},
		})
	}

	return indices
}

func main() {
	indexer := BasicIndexer{}
	err := Run(&indexer)
	if err != nil {
		log.Fatal(err)
	}
}
