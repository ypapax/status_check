package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Bind             string        `yaml:"bind"`
	CheckPeriod      time.Duration `yaml:"check_period"`
	DbType           string        `yaml:"db_type"`
	ConnectionString string        `yaml:"connection_string"`
	Workers          int           `yaml:"workers"`
	Schemas          []string      `yaml:"schemas"`
	Addresses        []string      `yaml:"addresses"`

	PipeType string      `yaml:"pipe_type"`
	Kafka    KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers     []string `yaml:"brokers"`
	StatusTopic string   `yaml:"status_topic"`
	ClientID    string   `yaml:"client_id"`
}

func (c Config) Validate() error {
	if len(c.Addresses) == 0 {
		return fmt.Errorf("Missing addresses list")
	}
	if len(c.Schemas) == 0 {
		return fmt.Errorf("Missing schemas list")
	}
	if c.Workers <= 0 {
		return fmt.Errorf("Workers amount should be positive")
	}
	if c.CheckPeriod == 0 {
		return fmt.Errorf("Empty check period")
	}
	return nil
}

func Parse(reader io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal([]byte(b), &c); err != nil {
		return nil, err
	}
	return &c, nil
}
