package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	BugNetConnectionString string
}

// Load config from json file
func (c *Config) Load() error {
	file, err := os.Open("Config.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return err
	}
	return nil
}
