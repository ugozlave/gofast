package gofast

import (
	"encoding/json"
	"os"
)

type ConfigProvider[T any] interface {
	Value() T
}

type Config[T any] struct {
	value T
}

func NewConfig[T any](v T) *Config[T] {
	data, _ := os.ReadFile("config.json")
	_ = json.Unmarshal(data, &v)
	return &Config[T]{value: v}
}

func (c *Config[T]) Value() T {
	return c.value
}

/*
** AppConfig
 */

type AppConfig struct {
	App struct {
		Name string `json:"Name"`
	} `json:"App"`
	Env string `json:"Environment"`
	Log struct {
		Level string `json:"Level"`
	} `json:"Logging"`
	Server struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	} `json:"Server"`
}

func NewAppConfig() *Config[AppConfig] {
	var v AppConfig
	v.Env = "development"
	v.Log.Level = "debug"
	v.Server.Port = 8080
	return NewConfig(v)
}
