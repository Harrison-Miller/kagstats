package config

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	Name                   string           `json:"name"`
	RCONTimeout            string           `json:"rconTimeout"`
	RCONMaxAttempts        int              `json:"rconMaxattempts"`
	BatchSize              int              `json:"batchSize"`
	CommitInterval         string           `json:"CommitInterval"`
	CommitIntervalDuration time.Duration    `json:"-"`
	DatabaseConnection     string           `json:"databaseConnection"`
	Monitoring             MonitoringConfig `json:"monitoring"`
	Servers                []ServerConfig   `json:"servers"`
	Indexer                IndexerConfig    `json:"indexer"`
	API                    APIConfig        `json:"api"`
}

type MonitoringConfig struct {
	Enabled             bool          `json:"enabled"`
	RefreshRate         string        `json:"refreshRate"`
	RefreshRateDuration time.Duration `json:"-"`
	Host                string        `json:"host"`
}

type ServerConfig struct {
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
	Address  string   `json:"address"`
	Port     string   `json:"port"`
	Password string   `json:"password"`
}

type IndexerConfig struct {
	BatchSize        int           `json:"batchSize"`
	Interval         string        `json:"interval"`
	IntervalDuration time.Duration `json:"-"`
}

type APIConfig struct {
	Host string `json:"host"`
}

func NewConfig() Config {
	c := Config{
		Name:               "Default Collector",
		RCONTimeout:        "10s",
		RCONMaxAttempts:    100,
		BatchSize:          20,
		CommitInterval:     "1m",
		DatabaseConnection: "root:password@tcp(127.0.0.1:3306)/kagstats",
	}
	c.Monitoring = NewMonitoringConfig()
	c.Indexer = NewIndexerConfig()
	c.API = NewAPIConfig()
	return c
}

func NewMonitoringConfig() MonitoringConfig {
	c := MonitoringConfig{
		Enabled:     true,
		RefreshRate: "30s",
		Host:        ":8080",
	}
	return c
}

func NewIndexerConfig() IndexerConfig {
	c := IndexerConfig{
		BatchSize: 100,
		Interval:  "30s",
	}
	return c
}

func NewAPIConfig() APIConfig {
	c := APIConfig{
		Host: ":80",
	}
	return c
}

func ParseDuration(value string, target *time.Duration) error {
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*target = d
	return nil
}

func Decode(r io.Reader) (Config, error) {
	c := NewConfig()

	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return c, errors.Wrap(err, "could not parse config json")
	}

	if err := ParseDurations(&c); err != nil {
		return c, err
	}

	return c, nil
}

func ParseDurations(c *Config) error {
	err := ParseDuration(c.CommitInterval, &c.CommitIntervalDuration)
	if err != nil {
		return errors.Wrap(err, "error parsing commit interval")
	}

	err = ParseDuration(c.Monitoring.RefreshRate, &c.Monitoring.RefreshRateDuration)
	if err != nil {
		return errors.Wrap(err, "error parsing monitoring refresh rate")
	}

	err = ParseDuration(c.Indexer.Interval, &c.Indexer.IntervalDuration)
	if err != nil {
		return errors.Wrap(err, "error parsing indexer interval")
	}
	return nil
}

func Read(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return NewConfig(), errors.Wrap(err, "could not open config")
	}

	config, err := Decode(bytes.NewReader(file))
	if err != nil {
		return config, err
	}

	return config, nil
}

func Get() (Config, error) {
	path := "settings.json"
	if v, ok := os.LookupEnv("KAGSTATS_CONFIG"); ok {
		path = v
	}

	return Read(path)
}
