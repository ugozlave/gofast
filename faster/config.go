package faster

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ugozlave/gofast"
)

type ConfigProvider[T any] interface {
	Value() T
}

/*
** Config
 */

type Config[T any] struct {
	value T
}

func NewConfigBuilder[T any](v T, keys ...string) Builder[*Config[T]] {
	return func(*gofast.BuilderContext) *Config[T] {
		return NewConfig(v, keys...)
	}
}

func NewConfig[T any](v T, keys ...string) *Config[T] {
	base := struct {
		Env string `json:"Environment"`
	}{}
	data, err := os.ReadFile(CONFIG.FILE_NAME + "." + CONFIG.FILE_EXT)
	if err == nil {
		_ = GetNestedConfig(data, &base, CONFIG.APPLICATION_PATH)
		_ = GetNestedConfig(data, &v, keys...)
	}
	data, err = ReadEnv(CONFIG.ENV_PREFIX)
	if err == nil {
		_ = GetNestedConfig(data, &base, CONFIG.APPLICATION_PATH)
	}
	if base.Env != "" {
		data, err := os.ReadFile(CONFIG.FILE_NAME + "." + base.Env + "." + CONFIG.FILE_EXT)
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

func NewAppConfig() *gofast.AppConfig {
	var v gofast.AppConfig
	v.Name = "gofast"
	v.Env = "development"
	v.Server.Host = ""
	v.Server.Port = 8080
	v = NewConfig(v, CONFIG.APPLICATION_PATH).Value()
	return &v
}

/*
** Settings
 */

type ConfigSettings struct {
	FILE_NAME        string
	FILE_EXT         string
	APPLICATION_PATH string
	ENV_PREFIX       string
}

var CONFIG = &ConfigSettings{
	FILE_NAME:        "config",
	FILE_EXT:         "json",
	APPLICATION_PATH: "Application",
	ENV_PREFIX:       "GOFAST",
}
