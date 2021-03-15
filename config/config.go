package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Feeds []string
}

func New(filePath string) *Config {
	config := &Config{}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		panic("Failed to parse config JSON.")
	}

	return config
}
