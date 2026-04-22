package gofast

type HealthChecker interface {
	HealthCheck() (string, bool, error)
}
