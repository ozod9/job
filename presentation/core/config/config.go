package config

import (
	"io/ioutil"
	"job/domain/models"

	"github.com/BurntSushi/toml"
)

var config models.Config

func Read(configFile string) (*models.Config, error) {

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}

	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		return nil, err
	}

	config.Application.Version = string(version)
	return &config, nil
}

func Get() *models.Config {
	return &config
}
