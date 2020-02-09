package main_test

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	gosubst "github.com/hews/gosubst"
)

func TestEnvironment(t *testing.T) {
	environ := os.Environ()
	env := gosubst.Environment()

	if len(env) != len(environ) {
		t.Errorf("len(Environment()) == %q; expected %q", len(env), len(environ))
	}

	names := make([]string, 0, len(env))
	for name := range env {
		names = append(names, name)
	}
	sort.Strings(names)
	sort.Strings(environ)

	for i, name := range names {
		if !strings.Contains(environ[i], name) || !strings.Contains(environ[i], env[name]) {
			t.Errorf("Environment() contains %q => %q, but actual is %q", name, env[name], environ[i])
		}
	}
}

func TestProcess(t *testing.T) {
	proc := gosubst.Process()

	proc_rv := reflect.ValueOf(proc)
	proc_rt := reflect.TypeOf(proc)
	for i := 0; i < proc_rv.NumField(); i++ {
		if proc_rv.Field(i).Interface() == nil {
			t.Errorf("Process().%s == nil, should have a value", proc_rt.Field(i).Name)
		}
	}
}

var shTests = []struct {
	cmd, out string
}{
	{"", ""},
	{"echo 'SOME TEXT'", "SOME TEXT\n"},
	{"printf 'MORE %s' 'TEXT\n'", "MORE TEXT\n"},
	{"pwd", must(os.Getwd()) + "\n"},
}

func TestSh(t *testing.T) {
	for _, test := range shTests {
		out, err := gosubst.Sh(test.cmd)
		if err != nil {
			t.Errorf("Sh(%q) returned error %q; expected nil", test.cmd, err)
		}
		if test.out != out {
			t.Errorf("Sh(%q) == %q; expected %q", test.cmd, out, test.out)
		}
	}
}

var shErrorTests = []struct {
	cmd, out string
	err      error
}{
	// Are errors.
	{"$(return 1)", "", fmt.Errorf("exit status 1")},
	{"xBBB/argh", "", fmt.Errorf("exit status 127")},
	// Are NOT errors.
	{"", "", nil},
	{"echo 'IM NOT AN ERR' >2", "", nil},
	{"echo 'IM NOT AN ERR' | tee /dev/stderr", "IM NOT AN ERR\n", nil},
}

func TestShErrors(t *testing.T) {
	for _, test := range shErrorTests {
		out, err := gosubst.Sh(test.cmd)
		if test.out != out {
			t.Errorf("Sh(%q) == %q; expected %q", test.cmd, out, test.out)
		}
		if (test.err == nil && err != nil) ||
			(test.err != nil && err == nil) ||
			(test.err != nil && err != nil && (test.err.Error() != err.Error())) {
			t.Errorf("Sh(%q) has error %q; expected %q", test.cmd, err, test.err)
		}
	}
}

func must(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}
