package config

import "time"

type Config struct {
	Bind             string        `yaml:"bind"`
	CheckPeriod      time.Duration `yaml:"check_period"`
	DbType           string        `yaml:"db_type"`
	ConnectionString string        `yaml:"connection_string"`
	Workers          int           `yaml:"workers"`
	Schemas          []string      `yaml:"schemas"`
	Addresses        []string      `yaml:"addresses"`
}

