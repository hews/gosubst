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

	home := os.Getenv("HOME")
	path := os.Getenv("PATH")

	environ := os.Environ()
	os.Clearenv()

	// Panics if these are lost.
	os.Setenv("HOME", home)
	os.Setenv("PATH", path)

	return func() {
		for _, envvar := range environ {
			pair := strings.SplitN(envvar, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}
}
