package main

import (
	"log"
	"strings"
	"time"

	"github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/indexer"
)

type MonthlyHittersIndexer struct {
}

func (i *MonthlyHittersIndexer) Name() string {
	return "monthly_hitters"
}

func (i *MonthlyHittersIndexer) Version() int {
	return 1
}

func (i *MonthlyHittersIndexer) Keys() []IndexKey {
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

func fieldNameToSql(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func fieldNamesToSql(names []string) []string {
	for i, name := range names {
		names[i] = fieldNameToSql(name)
	}
	return names
}

func (i *MonthlyHittersIndexer) Counters() []string {
	return fieldNamesToSql(models.HitterNames)
}

func (i *MonthlyHittersIndexer) Index(kill Kill) []Index {
	var indices []Index
	if kill.KillerID != kill.VictimID && !kill.TeamKill {

		date := time.Unix(0, kill.Time*int64(time.Millisecond))
		year := date.Year()
		month := date.Month()

		hitters := map[string]int{}
		fieldNames := i.Counters()
		for _, name := range fieldNames {
			if name == fieldNameToSql(models.HitterName(kill.Hitter)) {
				hitters[name] = 1
			} else {
				hitters[name] = 0
			}
		}

		indices = append(indices, Index{
			Keys:     []interface{}{int(kill.KillerID), year, int(month)},
			Counters: hitters,
		})
	}

	return indices
}

func main() {
	indexer := MonthlyHittersIndexer{}
	err := Run(&indexer)
	if err != nil {
		log.Fatal(err)
	}
}
