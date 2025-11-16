package faster

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ugozlave/gofast"
)

type FastConfig[T any] struct {
	value T
}

/*
** AppConfig
 */

func NewConfig[T any](v T, keys ...string) *FastConfig[T] {
	base := struct {
		Env string `json:"Environment"`
	}{}
	data, err := os.ReadFile(gofast.SETTINGS.CONFIG_FILE_NAME + "." + gofast.SETTINGS.CONFIG_FILE_EXT)
	if err == nil {
		_ = json.Unmarshal(data, &base)
		_ = GetNestedConfig(data, &v, keys...)
	}
	env, ok := os.LookupEnv("ENVIRONMENT")
	if ok {
		base.Env = env
	}
	if base.Env != "" {
		data, err = os.ReadFile(gofast.SETTINGS.CONFIG_FILE_NAME + "." + base.Env + "." + gofast.SETTINGS.CONFIG_FILE_EXT)
		if err == nil {
			_ = GetNestedConfig(data, &v, keys...)
		}
	}
	return &FastConfig[T]{value: v}
}

func (c *FastConfig[T]) Value() T {
	return c.value
}

func NewDefaultAppConfig() *FastConfig[gofast.AppConfig] {
	var v gofast.AppConfig
	v.App.Name = "gofast"
	v.Env = "development"
	v.Log.Level = "debug"
	v.Server.Host = ""
	v.Server.Port = 8080
	return NewConfig(v)
}

func GetNestedConfig[T any](data []byte, v *T, keys ...string) error {
	var root map[string]interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return err
	}

	current := root

	for _, key := range keys {
		value, ok := current[key]
		if !ok {
			return errors.New("config key not found: " + key)
		}
		current, ok = value.(map[string]interface{})
		if !ok {
			return errors.New("config key is not a map: " + key)
		}
	}

	data, err := json.Marshal(current)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}
