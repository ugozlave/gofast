package gofast

import "os"

/*
** AppConfig
 */

type AppConfig struct {
	Name   string `json:"Name"`
	Env    string `json:"Environment"`
	Server struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	} `json:"Server"`
}

func (c *AppConfig) Default() *AppConfig {
	if c.Name == "" {
		name, err := os.Executable()
		if err != nil {
			panic(err)
		}
		c.Name = name
	}
	if c.Env == "" {
		c.Env = "development"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	return c
}
