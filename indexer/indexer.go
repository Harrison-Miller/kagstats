package indexer

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Harrison-Miller/kagstats/common/models"
	. "github.com/Harrison-Miller/kagstats/common/models"
	"github.com/jmoiron/sqlx"
)

type IndexKey struct {
	Name   string
	Table  string
	Column string
}

type Index struct {
	Keys     []int
	Counters map[string]int
}

func Equal(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func (a Index) Add(b Index) error {
	if !Equal(a.Keys, b.Keys) {
		return fmt.Errorf("Index Key Mismatch")
	}

	for k := range b.Counters {
		a.Counters[k] += b.Counters[k]
	}

	return nil
}

type Indexer interface {
	Name() string
	Version() int

	Keys() []IndexKey
	Counters() []string
}

type KillsIndexer interface {
	Indexer
	Index(kill Kill) []Index
}

func BuildTable(indexer Indexer) string {
	keys := indexer.Keys()
	var keyNames []string

	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TABLE IF NOT EXISTS %s (", indexer.Name())

	for _, k := range keys {
		fmt.Fprintf(&b, "%s INT NOT NULL,", k.Name)
		if k.Table != "" && k.Column != "" {
			fmt.Fprintf(&b, "FOREIGN KEY(%s) REFERENCES %s(%s) ON DELETE CASCADE,", k.Name, k.Table, k.Column)
		}

		keyNames = append(keyNames, k.Name)
	}

	fmt.Fprintf(&b, "PRIMARY KEY (%s) )", strings.Join(keyNames, ","))

	return b.String()
}

func Init(indexer Indexer, db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS indexer_info (
			key_name VARCHAR(30) PRIMARY KEY,
			value INT NOT NULL
		)`)
	if err != nil {
		return err
	}

	db.Exec("INSERT INTO indexer_info (key_name, value) VALUES(?, ?)", indexer.Name()+"_version", indexer.Version())
	row := db.QueryRow("SELECT value FROM indexer_info WHERE key_name=?", indexer.Name()+"_version")
	var currentVersion int64
	err = row.Scan(&currentVersion)
	if err != nil {
		return err
	}

	if currentVersion < int64(indexer.Version()) {
		fmt.Printf("Currently deployed version %d, updating %s to newest version: %d\n", currentVersion, indexer.Name(), indexer.Version())
		_, err = db.Exec("UPDATE indexer_info SET value=? WHERE key_name=?", indexer.Version(), indexer.Name()+"_version")
		if err != nil {
			return err
		}

		_, err = db.Exec("UPDATE indexer_info SET value=? WHERE key_name=?", 0, indexer.Name())
		if err != nil {
			return err
		}

		_, err = db.Exec(fmt.Sprintf("DROP TABLE %s", indexer.Name()))
		if err != nil {
			return err
		}
	} else if currentVersion > int64(indexer.Version()) {
		panic(fmt.Sprintf("Current deployed version %d of %s, version %d is too old to be deployed", currentVersion, indexer.Name(), indexer.Version()))
	}

	db.Exec("INSERT INTO indexer_info (key_name, value) VALUES(?, ?)", indexer.Name(), 0)

	table := BuildTable(indexer)
	_, err = db.Exec(table)
	if err != nil {
		return err
	}

	for _, c := range indexer.Counters() {
		db.Exec(fmt.Sprintf("ALTER TABLE %s ADD %s INT NOT NULL DEFAULT 0", indexer.Name(), c))
	}

	return nil
}

func CurrentIndex(indexer Indexer, tx *sql.Tx) (int64, error) {
	var currentIndex int64
	row := tx.QueryRow("SELECT value from indexer_info WHERE key_name=?", indexer.Name())
	err := row.Scan(&currentIndex)
	return currentIndex, err
}

func UnprocessedRows(indexer Indexer, batchSize int, tx *sql.Tx) (int64, *sql.Rows, error) {
	currentIndex, err := CurrentIndex(indexer, tx)
	if err != nil {
		return 0, nil, err
	}

	upperBound := currentIndex + int64(batchSize) + 1
	rows, err := tx.Query("SELECT * from kills WHERE ID>? AND ID<?", currentIndex, upperBound)
	return currentIndex, rows, err
}

func BuildInsert(indexer Indexer, tx *sql.Tx) (*sql.Stmt, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "INSERT INTO %s (", indexer.Name())

	for _, k := range indexer.Keys() {
		fmt.Fprintf(&b, "%s,", k.Name)
	}
	fmt.Fprintf(&b, "%s) VALUES (", strings.Join(indexer.Counters(), ","))

	count := len(indexer.Keys()) + len(indexer.Counters())
	fmt.Fprintf(&b, "%s) ON DUPLICATE KEY UPDATE ", strings.TrimRight(strings.Repeat("?,", count), ","))

	for _, c := range indexer.Counters() {
		fmt.Fprintf(&b, "%s=%s+?,", c, c)
	}

	str := strings.TrimRight(b.String(), ",")
	stmnt, err := tx.Prepare(str)
	return stmnt, err
}

var players = make(map[int64]models.Player)

func SkipKill(kill models.Kill, db *sqlx.DB) bool {
	var victim Player
	if cacheVictim, ok := players[kill.VictimID]; ok {
		victim = cacheVictim
	} else {
		err := db.Get(&victim, "SELECT * FROM players WHERE ID=?", kill.VictimID)
		if err != nil {
			return false
		}

		players[kill.VictimID] = victim
	}

	if victim.StatsBan {
		return true
	}

	var killer Player
	if cacheKiller, ok := players[kill.KillerID]; ok {
		killer = cacheKiller
	} else {
		err := db.Get(&killer, "SELECT * FROM players WHERE ID=?", kill.KillerID)
		if err != nil {
			return false
		}

		players[kill.KillerID] = killer
	}

	if killer.StatsBan {
		return true
	}

	return false
}

func Process(indexer KillsIndexer, batchSize int, db *sqlx.DB) (int, error) {
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return 0, err
	}

	currentIndex, rows, err := UnprocessedRows(indexer, batchSize, tx)
	defer rows.Close()
	if err != nil {
		return 0, err
	}

	newIndex := currentIndex
	updates := make(map[string]Index)
	for rows.Next() {
		var kill Kill
		if err := rows.Scan(&kill.ID, &kill.KillerID, &kill.VictimID, &kill.KillerClass,
			&kill.VictimClass, &kill.Hitter, &kill.Time, &kill.ServerID, &kill.TeamKill); err != nil {
			return 0, err
		}

		if SkipKill(kill, db) {
			continue
		}

		indices := indexer.Index(kill)
		for _, index := range indices {
			if len(index.Keys) != len(indexer.Keys()) {
				return 0, fmt.Errorf("Indexer failed to return the correct number of keys\n\texpected: %d got %d", len(indexer.Keys()), len(index.Keys))
			}

			mapKey := fmt.Sprintf("%v", index.Keys)

			if _, ok := updates[mapKey]; ok {
				err := updates[mapKey].Add(index)
				if err != nil {
					return 0, err
				}

			} else {
				updates[mapKey] = index
			}

			if kill.ID > newIndex {
				newIndex = kill.ID
			}

		}

	}

	stmnt, err := BuildInsert(indexer, tx)
	defer stmnt.Close()
	if err != nil {
		return 0, err
	}

	for _, u := range updates {
		keysCount := len(indexer.Keys())
		countersCount := len(indexer.Counters())
		args := make([]interface{}, keysCount+countersCount*2)
		for i, v := range u.Keys {
			args[i] = v
		}

		for i, v := range indexer.Counters() {
			args[keysCount+i] = u.Counters[v]
		}

		for i, v := range indexer.Counters() {
			args[keysCount+countersCount+i] = u.Counters[v]
		}

		_, err = stmnt.Exec(args...)
		if err != nil {
			return 0, err
		}
	}

	if currentIndex != newIndex {
		_, err = tx.Exec("UPDATE indexer_info SET value=? WHERE key_name=? AND value=?", newIndex, indexer.Name(), currentIndex)
		if err != nil {
			return 0, err
		}

		err = tx.Commit()
		return int(newIndex - currentIndex), err
	}

	return 0, nil
}
