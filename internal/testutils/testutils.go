package testutils

import (
	"os"
	"strings"
	"testing"
)

// ClearEnvironment ...
func ClearEnvironment(t *testing.T) func() {
	t.Helper()

	environ := os.Environ()
	home := os.Getenv("HOME")

	os.Clearenv()
	os.Setenv("HOME", home) // Panics if this is lost.

	return func() {
		for _, envvar := range environ {
			pair := strings.SplitN(envvar, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}
}
