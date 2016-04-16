package main

import "gopkg.in/gcfg.v1"

// Config is config
type Config struct {
	Telegram struct {
		Token    string
		Username string
	}
}

func getConfig() (Config, error) {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "./config.cfg")
	return cfg, err
}
