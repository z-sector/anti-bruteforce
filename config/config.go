package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const prefEnvVar = "app"

type Config struct {
	Port                int           `split_words:"true" default:"9000"`
	JSONFormat          bool          `envconfig:"json_format" split_words:"true" default:"false"`
	PgDSN               string        `split_words:"true" required:"true"`
	RedisDSN            string        `split_words:"true" required:"true"`
	LimitBurstLogin     int           `split_words:"true" default:"10"`
	LimitPeriodLogin    time.Duration `split_words:"true" default:"1m"`
	LimitBurstPassword  int           `split_words:"true" default:"100"`
	LimitPeriodPassword time.Duration `split_words:"true" default:"1m"`
	LimitBurstIP        int           `split_words:"true" default:"1000"`
	LimitPeriodIP       time.Duration `split_words:"true" default:"1m"`
}

func MustConfig() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewConfig() (Config, error) {
	var config Config
	err := envconfig.Process(prefEnvVar, &config)
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return config, nil
}
