package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Debug bool
	DB    struct {
		Host         string
		Port         string
		DatabaseName string
		Username     string
		Password     string
	}
}

func NewConfig() *Config {
	config := &Config{}

	if file, err := os.Open("./config/config.json"); err == nil {
		defer file.Close()
		jsonByte, err := ioutil.ReadAll(file)
		if err == nil {
			json.Unmarshal(jsonByte, &config)
		}
	}

	return config
}
