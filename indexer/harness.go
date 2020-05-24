package indexer

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Harrison-Miller/kagstats/common/configs"
	"github.com/Harrison-Miller/kagstats/common/utils"
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
	version, _ := os.LookupEnv("VERSION")
	log.Printf("KAG Stats  %s\n", version)

	fmt.Printf("Starting Indexer: %s Version %d\n", indexer.Name(), indexer.Version())

	config, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("Error reading indexer configuration\n%v\n", err)
	}

	db, err := utils.ConnectToDatabase(config.DatabaseConnection, 10)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database!")

	err = Init(indexer, db)
	if err != nil {
		return fmt.Errorf("Error initing indexer tables\n%v", err)
	}

	fmt.Printf("Batch Size: %d\n", config.Indexer.BatchSize)
	fmt.Printf("Processing Interval: %s\n", config.Indexer.Interval)

	for {
		var processed int
		var err error

		processed, err = Process(indexer, config.Indexer.BatchSize, db)

		if err != nil {
			log.Println(err)
		} else if processed != 0 {
			log.Printf("Processed %d rows\n", processed)
		}

		time.Sleep(config.Indexer.IntervalDuration)
	}

}
