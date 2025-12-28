package faster

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ugozlave/gofast"
)

/*
** Config
 */

type Config[T any] struct {
	value T
}

func NewConfig[T any](v T, keys ...string) *Config[T] {
	base := struct {
		Env string `json:"Environment"`
	}{}
	data, err := os.ReadFile(gofast.SETTINGS.CONFIG_FILE_NAME + "." + gofast.SETTINGS.CONFIG_FILE_EXT)
	if err == nil {
		_ = GetNestedConfig(data, &base, gofast.SETTINGS.CONFIG_APPLICATION_KEY)
		_ = GetNestedConfig(data, &v, keys...)
	}
	data, err = ReadEnv(gofast.SETTINGS.ENV_PREFIX)
	if err == nil {
		_ = GetNestedConfig(data, &base, gofast.SETTINGS.CONFIG_APPLICATION_KEY)
	}
	if base.Env != "" {
		data, err := os.ReadFile(gofast.SETTINGS.CONFIG_FILE_NAME + "." + base.Env + "." + gofast.SETTINGS.CONFIG_FILE_EXT)
		if err == nil {
			_ = GetNestedConfig(data, &v, keys...)
		}
	}
	if err == nil {
		_ = GetNestedConfig(data, &v, keys...)
	}
	return &Config[T]{value: v}
}

func (c *Config[T]) Value() T {
	return c.value
}

func GetNestedConfig[T any](data []byte, v *T, keys ...string) error {
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		return err
	}

	current := root

	for _, key := range keys {
		value, ok := current[key]
		if !ok {
			return errors.New("config key not found: " + key)
		}
		current, ok = value.(map[string]any)
		if !ok {
			return errors.New("config key is not a map: " + key)
		}
	}

	data, err := json.Marshal(current)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ReadEnv(prefix string) ([]byte, error) {
	root := make(map[string]any)
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix+"_") {
			continue
		}
		parts := strings.SplitN(env, "=", 2)
		k := parts[0]
		v := parts[1]
		keys := strings.Split(k[len(prefix)+1:], "_")
		current := root
		for i, key := range keys {
			if i == len(keys)-1 {
				if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
					var arr []any
					if err := json.Unmarshal([]byte(v), &arr); err != nil {
						return nil, errors.New("config value is not an array: " + v)
					}
					current[key] = arr
				} else {
					current[key] = v
				}
				break
			}
			_, ok := current[key]
			if !ok {
				current[key] = make(map[string]any)
			}
			current, ok = current[key].(map[string]any)
			if !ok {
				return nil, errors.New("config key is not a map: " + key)
			}
		}
	}
	return json.Marshal(root)
}

/*
** AppConfig
 */

func NewAppConfig() *Config[gofast.AppConfig] {
	var v gofast.AppConfig
	v.Name = "gofast"
	v.Env = "development"
	v.Log.Level = "debug"
	v.Server.Host = ""
	v.Server.Port = 8080
	return NewConfig(v, gofast.SETTINGS.CONFIG_APPLICATION_KEY)
}
