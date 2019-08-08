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
	DBOptions struct {
		SyncTables           []string
		IgnoreTables         []string
		MergeTables          map[string][]string
		DefaultPk            string
		DatetimeFormatFields []string
		Tables               map[string]struct {
			PK                   string
			DatetimeFormatFields []string
			IgnoreFields         []string
		}
	}
	ES struct {
		IndexPrefix string
		Urls        []string
		BaseAuth    struct {
			Username string
			Password string
		}
	}
	SizePerTime int64
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
