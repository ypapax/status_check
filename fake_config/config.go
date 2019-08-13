package fake_config

import (
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Ports []Port `yaml:"ports"`
}

func Parse(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &c, nil
}

type Port struct {
	From        int   `yaml:"from"`
	To          int   `yaml:"to"`
	StatusCodes []int `yaml:"codes"`
	DelayMS     []int `yaml:"delay_ms"`
}
