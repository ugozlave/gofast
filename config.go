package gofast

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Config[T any] interface {
	Value() T
}

type ConfigResolver interface {
	Path() []string
}

/*
** Config
 */

type FastConfig[T any] struct {
	value T
}

func ConfigBuilder[T any](v T, keys ...string) Builder[*FastConfig[T]] {
	return func(*BuilderContext) *FastConfig[T] {
		return NewConfig(v, keys...)
	}
}

func NewConfig[T any](v T, keys ...string) *FastConfig[T] {
	if len(keys) == 0 {
		if resolver, ok := any(v).(ConfigResolver); ok {
			keys = resolver.Path()
		}
	}
	for _, file := range ConfigFiles.files {
		data, err := os.ReadFile(file)
		if err == nil {
			_ = get_nested_config(data, &v, keys...)
		}
	}
	if ConfigFiles.env {
		data, err := read_env(CONFIG.ENV_PREFIX)
		if err == nil {
			_ = get_nested_config(data, &v, keys...)
		}
	}
	return &FastConfig[T]{value: v}
}

func (c *FastConfig[T]) Value() T {
	return c.value
}

func get_nested_config[T any](data []byte, v *T, keys ...string) error {
	root := map[string]any{}
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
		return err
	}

	return nil
}

func read_env(prefix string) ([]byte, error) {
	root := map[string]any{}
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		val := parts[1]

		if !strings.HasPrefix(key, prefix) {
			continue
		}

		path := strings.Split(strings.TrimPrefix(key, prefix+"_"), "__")

		insert(root, path, val)
	}

	return json.Marshal(root)
}

func insert(m map[string]any, path []string, val string) {
	if len(path) == 1 {
		var v any
		if err := json.Unmarshal([]byte(val), &v); err != nil {
			v = val
		}
		m[path[0]] = v
		return
	}

	key := path[0]

	next, ok := m[key]
	if !ok {
		child := map[string]any{}
		m[key] = child
		insert(child, path[1:], val)
		return
	}

	insert(next.(map[string]any), path[1:], val)
}

/*
** AppConfig
 */

type AppConfig struct {
	Name   string `json:"Name"`
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
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	return c
}

/*
** ConfigHelper
 */

type ConfigHelper struct {
	files []string
	env   bool
}

func (c *ConfigHelper) Add(file string) {
	c.files = append(c.files, file)
}

func (c *ConfigHelper) Env(b bool) {
	c.env = b
}

var ConfigFiles = &ConfigHelper{}

/*
** Settings
 */

type ConfigSettings struct {
	APPLICATION_PATH []string
	ENV_PREFIX       string
}

var CONFIG = &ConfigSettings{
	APPLICATION_PATH: []string{},
	ENV_PREFIX:       "GOFAST",
}
