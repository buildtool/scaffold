// +build !prod

package config

import (
	"os"
)

func SetEnv(key, value string) func() {
	_ = os.Setenv(key, value)
	return func() { _ = os.Unsetenv(key) }
}
