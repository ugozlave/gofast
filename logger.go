package gofast

const (
	LogApplication string = "application"
	LogEnvironment string = "environment"
	LogService     string = "service"
	LogRequestId   string = "requestId"
	LogHttp        string = "http"
	LogMethod      string = "method"
	LogHost        string = "host"
	LogUrl         string = "url"
	LogRemote      string = "remote"
	LogAgent       string = "agent"
	LogStatus      string = "status"
	LogDuration    string = "duration"
)

type Logger interface {
	Dbg(msg string, args ...any)
	Inf(msg string, args ...any)
	Wrn(msg string, args ...any)
	Err(msg string, args ...any)
	With(args ...any) Logger
	WithGroup(name string) Logger
}
