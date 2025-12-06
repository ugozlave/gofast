package gofast

type Config[T any] interface {
	Value() T
}

/*
** AppConfig
 */

type AppConfig struct {
	Name string `json:"Name"`
	Env  string `json:"Environment"`
	Log  struct {
		Level string `json:"Level"`
	} `json:"Logging"`
	Server struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	} `json:"Server"`
}
