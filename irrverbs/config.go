package main

import "code.google.com/p/gcfg"

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
