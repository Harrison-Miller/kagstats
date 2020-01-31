package main

import (
	"log"
	"time"

	. "github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/indexer"
)

type MonthlyIndexer struct {
}

func (i *MonthlyIndexer) Name() string {
	return "monthly_stats"
}

func (i *MonthlyIndexer) Version() int {
	return 3
}

func (i *MonthlyIndexer) Keys() []IndexKey {
	var keys []IndexKey
	keys = append(keys, IndexKey{
		Name:   "playerID",
		Table:  "players",
		Column: "ID",
	}, IndexKey{
		Name: "year",
	}, IndexKey{
		Name: "month",
	})
	return keys
}

func (i *MonthlyIndexer) Counters() []string {
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

func (i *MonthlyIndexer) Index(kill Kill) []Index {
	var indices []Index

	date := time.Unix(0, kill.Time*int64(time.Millisecond))
	year := date.Year()
	month := date.Month()

	if kill.KillerID != kill.VictimID {
		if kill.TeamKill {
			indices = append(indices, Index{
				Keys:     []int{int(kill.KillerID), year, int(month)},
				Counters: map[string]int{"teamkills": 1},
			})
		} else {
			indices = append(indices, Index{
				Keys: []int{int(kill.KillerID), year, int(month)},
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
			Keys: []int{int(kill.VictimID), year, int(month)},
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
	indexer := MonthlyIndexer{}
	err := Run(&indexer)
	if err != nil {
		log.Fatal(err)
	}
}
