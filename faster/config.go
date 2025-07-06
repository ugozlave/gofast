package faster

import (
	"encoding/json"
	"os"

	"github.com/ugozlave/gofast"
)

type Config[T any] struct {
	value T
}

/*
** AppConfig
 */

func NewConfig[T any](v T) *Config[T] {
	base := struct {
		Env string `json:"Environment"`
	}{}
	data, err := os.ReadFile("config.json")
	if err == nil {
		_ = json.Unmarshal(data, &base)
		_ = json.Unmarshal(data, &v)
	}
	env, ok := os.LookupEnv("ENVIRONMENT")
	if ok {
		base.Env = env
	}
	if base.Env != "" {
		data, err = os.ReadFile("config." + base.Env + ".json")
		if err == nil {
			_ = json.Unmarshal(data, &v)
		}
	}
	return &Config[T]{value: v}
}

func (c *Config[T]) Value() T {
	return c.value
}

func NewAppConfig() *Config[gofast.AppConfig] {
	var v gofast.AppConfig
	v.Env = "development"
	v.Log.Level = "debug"
	v.Server.Port = 8080
	return NewConfig(v)
}
