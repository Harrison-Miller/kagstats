package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func ConnectToDatabase(connection string, maxAttempts int) (*sqlx.DB, error) {
	wait := 1
	attempts := 0
	for {
		db, err := sqlx.Connect("mysql", connection)
		if err != nil {
			log.Printf("%v\n", errors.Wrap(err, "could not connect to database"))
		} else {
			return db, nil
		}

		attempts = attempts + 1
		if attempts > maxAttempts {
			return nil, errors.Wrap(err, fmt.Sprintf("could not connect to database after %d attempts", attempts))
		}

		time.Sleep(time.Duration(wait) * time.Second)
		wait = wait * 2

	}
}

func NowAsUnixMilliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}
