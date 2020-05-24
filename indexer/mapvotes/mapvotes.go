package main

import (
	"log"

	. "github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/indexer"
)

type MapVoteStatsIndexer struct {
}

func (i *MapVoteStatsIndexer) Name() string {
	return "map_vote_stats"
}

func (i *MapVoteStatsIndexer) Version() int {
	return 1
}

func (i *MapVoteStatsIndexer) Keys() []IndexKey {
	var keys []IndexKey
	keys = append(keys, IndexKey{
		Name: "mapName",
		Type: "varchar(255)",
	})
	return keys
}

func (i *MapVoteStatsIndexer) Counters() []string {
	return []string{"ballots", "votes", "wins"}
}

func OneIfGreater(a int64, b int64, c int64) int {
	if a > b && a > c {
		return 1
	}

	return 0
}

func (i *MapVoteStatsIndexer) Index(votes MapVotes) []Index {
	var indices []Index

	indices = append(indices, Index{
		Keys: []interface{}{votes.Map1Name},
		Counters: map[string]int{"ballots": 1,
			"votes": int(votes.Map1Votes),
			"wins":  OneIfGreater(votes.Map1Votes, votes.Map2Votes, votes.RandomVotes)},
	})

	indices = append(indices, Index{
		Keys: []interface{}{votes.Map2Name},
		Counters: map[string]int{"ballots": 1,
			"votes": int(votes.Map2Votes),
			"wins":  OneIfGreater(votes.Map2Votes, votes.Map1Votes, votes.RandomVotes)},
	})

	return indices
}

func main() {
	indexer := MapVoteStatsIndexer{}
	err := Run(&indexer)
	if err != nil {
		log.Fatal(err)
	}
}
