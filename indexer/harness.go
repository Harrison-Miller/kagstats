package indexer

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/pkg/errors"
	// The harness run is used to wrap the entire functionality of an indexer including connecting to the database
	_ "github.com/go-sql-driver/mysql"
)

func ReadConfig() (configs.Config, error) {
	config, _ := configs.Get()

	if value, ok := os.LookupEnv("INDEXER_DB"); ok {
		config.DatabaseConnection = value
	}

	var err error
	if value, ok := os.LookupEnv("INDEXER_BATCHSIZE"); ok {
		config.Indexer.BatchSize, err = strconv.Atoi(value)
		if err != nil {
			return config, errors.Wrap(err, "could convert INDEXER_BATCHSIZE to int")
		}
	}

	if value, ok := os.LookupEnv("INDEXER_INTERVAL"); ok {
		config.Indexer.Interval = value
		err := configs.ParseDuration(config.Indexer.Interval, &config.Indexer.IntervalDuration)
		if err != nil {
			return config, errors.Wrap(err, "error parsing indexer interval")
		}
	}

	return config, nil
}

func Run(indexer Indexer) error {
	fmt.Printf("Starting Indexer: %s Version %d\n", indexer.Name(), indexer.Version())

	config, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("Error reading indexer configuration\n%v\n", err)
	}

	var db *sql.DB
	attempts := 0
	for {
		db, err = sql.Open("mysql", config.DatabaseConnection)
		if err != nil {
			log.Printf("%v", err)
		}

		err = db.Ping()
		if err != nil {
			fmt.Printf("Error connecting to database %s\n%v\n", config.DatabaseConnection, err)
		} else {
			break
		}

		attempts = attempts + 1
		if attempts > 10 {
			return fmt.Errorf("Could not connect to database after 10 attempts")
		}

		time.Sleep(5 * time.Second)
	}

	err = Init(indexer, db)
	if err != nil {
		return fmt.Errorf("Error initing indexer tables\n%v", err)
	}

	fmt.Printf("Batch Size: %d\n", config.Indexer.BatchSize)
	fmt.Printf("Processing Interval: %s\n", config.Indexer.Interval)

	for {
		var processed int
		var err error

		switch indexer.(type) {
		case KillsIndexer:
			processed, err = Process(indexer.(KillsIndexer), config.Indexer.BatchSize, db)
		}

		if err != nil {
			log.Println(err)
		} else if processed != 0 {
			log.Printf("Processed %d rows\n", processed)
		}

		time.Sleep(config.Indexer.IntervalDuration)
	}

}
