package indexer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type IndexerConfig struct {
	BatchSize    int    `json:"batchSize"`
	Interval     string `json:"interval"`
	IntervalTime time.Duration
}

type Config struct {
	DatabaseConnection string        `json:"databaseConnection"`
	Indexer            IndexerConfig `json:"indexer"`
}

func ReadConfig() (Config, error) {
	path := "settings.json"
	if value, ok := os.LookupEnv("KAGSTATS_CONFIG"); ok {
		path = value
	}

	config := Config{
		Indexer: IndexerConfig{
			BatchSize: 100,
			Interval:  "30s",
		},
	}

	file, err := ioutil.ReadFile(path)
	if err == nil {
		err = json.Unmarshal([]byte(file), &config)
		if err != nil {
			return config, fmt.Errorf("Error parsing indexer json settings %v", err)
		}
	}

	if value, ok := os.LookupEnv("INDEXER_DB"); ok {
		config.DatabaseConnection = value
	}

	if value, ok := os.LookupEnv("INDEXER_BATCHSIZE"); ok {
		config.Indexer.BatchSize, err = strconv.Atoi(value)
		if err != nil {
			return config, fmt.Errorf("INDEXER_BATCHSIZE set to %s could not be parsed", value)
		}
	}

	if value, ok := os.LookupEnv("INDEXER_INTERVAL"); ok {
		config.Indexer.Interval = value
	}

	config.Indexer.IntervalTime, err = time.ParseDuration(config.Indexer.Interval)
	if err != nil {
		return config, fmt.Errorf("Indexer interval set to %s could not be parsed", config.Indexer.Interval)
	}

	return config, nil
}

func Run(indexer Indexer) error {
	config, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("Error reading indexer configuration\n%v", err)
	}

	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("Error connecting to database %s\n%v", config.DatabaseConnection, err)
	}

	err = Init(indexer, db)
	if err != nil {
		return fmt.Errorf("Error initing indexer tables\n%v", err)
	}

	fmt.Printf("Starting Indexer: %s\n", indexer.Name())
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

		time.Sleep(config.Indexer.IntervalTime)
	}

	return nil
}
