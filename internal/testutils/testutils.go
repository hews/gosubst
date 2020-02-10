package testutils

import (
	"os"
	"strings"
	"testing"
)

// ClearEnvironment clears out the environment for testing, returning a
// function you can run on defer to restore the environment.
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
