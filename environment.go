package gofast

import (
	"os"
)

/*
** EnvironmentHelper
 */

type EnvironmentHelper struct {
	value string
}

func (e *EnvironmentHelper) Read() {
	value, ok := os.LookupEnv(CONFIG.ENV_PREFIX + "_ENVIRONMENT")
	if !ok {
		value = "production"
	}
	e.value = value
}

func (e *EnvironmentHelper) Get() string {
	return e.value
}

var Environment = &EnvironmentHelper{}
